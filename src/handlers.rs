use crate::{error::AppError, nws, state::AppState};
use axum::{
    extract::{Path, State},
    response::Html,
};
use redis::{AsyncCommands, SetExpiry, SetOptions};

pub async fn forecast(
    State(state): State<AppState>,
    Path(coords): Path<String>,
) -> Result<Html<String>, AppError> {
    let (lat, long) = coords
        .split_once(",")
        .ok_or_else(|| anyhow::anyhow!("Invalid coords format"))?;
    let lat = lat.trim().parse::<f64>()?;
    let long = long.trim().parse::<f64>()?;

    let lat = (lat * 1000.0).round() / 1000.0;
    let long = (long * 1000.0).round() / 1000.0;

    let mut con = state.redis.clone();

    let redis_key = format!("{},{}", lat, long);

    if let Ok(Some(cached)) = con.get(&redis_key).await {
        return Ok(Html(cached));
    }

    let points = nws::get_points(&state.client, &state.base_url, lat, long).await?;
    let points_properties = points.properties;

    let (alerts, forecast) = tokio::try_join!(
        nws::get_alerts(&state.client, &state.base_url, lat, long),
        nws::get_forecast(
            &state.client,
            &state.base_url,
            points_properties.grid_id,
            points_properties.grid_x,
            points_properties.grid_y,
        )
    )?;

    let parts = forecast.properties.periods;
    let forecasts = parts
        .into_iter()
        .map(|p| format!("<p>{}: {}</p>", p.name, p.detailed_forecast))
        .collect::<Vec<String>>()
        .join(" ");

    let parsed_alerts = alerts
        .features
        .into_iter()
        .map(|a| {
            format!(
                "<p>{}: {}</p>",
                a.properties.headline, a.properties.description
            )
        })
        .collect::<Vec<String>>()
        .join(" ");

    let mut resp = String::new();
    resp.push_str(format!("<h3>Forecast for {}, {}</h3>", lat, long).as_ref());
    resp.push_str(forecasts.as_ref());
    if !parsed_alerts.is_empty() {
        resp.push_str("<h3>Alerts</h3>");
        resp.push_str(parsed_alerts.as_ref());
    }

    let options = SetOptions::default().with_expiration(SetExpiry::EX(600));
    let _: () = con
        .set_options(&redis_key, &resp, options)
        .await
        .unwrap_or_default();

    Ok(Html(resp))
}
