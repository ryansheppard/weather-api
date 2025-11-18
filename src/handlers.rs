use crate::{nws, state::AppState};
use axum::{
    extract::{Path, State},
    response::Html,
};

pub async fn forecast(State(state): State<AppState>, Path(coords): Path<String>) -> Html<String> {
    let (lat, long) = coords.split_once(",").unwrap();
    let lat = lat.trim().parse::<f64>().unwrap();
    let long = long.trim().parse::<f64>().unwrap();

    let points = nws::get_points(&state.client, &state.base_url, lat, long)
        .await
        .unwrap();
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
    )
    .unwrap();

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

    Html(resp)
}
