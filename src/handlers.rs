use crate::{
    error::AppError,
    nws,
    state::AppState,
    types::{AlertResponse, ForecastPeriod},
};
use axum::{
    extract::{Path, State},
    response::Html,
};
use redis::{AsyncCommands, SetExpiry, SetOptions};

pub async fn forecast(
    State(state): State<AppState>,
    Path(coords): Path<String>,
) -> Result<Html<String>, AppError> {
    let (lat, long) = parse_coordinates(coords)?;

    let redis_key = format!("{},{}", lat, long);
    if let Some(ref mut con) = state.redis.clone()
        && let Ok(Some(cached)) = con.get(&redis_key).await
    {
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

    let forecasts = format_forecast(forecast.properties.periods);
    let alerts = format_alerts(alerts);

    let resp = build_forecast_html(lat, long, forecasts, alerts);

    if let Some(ref mut con) = state.redis.clone() {
        let options = SetOptions::default().with_expiration(SetExpiry::EX(600));
        let _: () = con
            .set_options(&redis_key, &resp, options)
            .await
            .unwrap_or_default();
    }

    Ok(Html(resp))
}

fn format_forecast(periods: Vec<ForecastPeriod>) -> String {
    periods
        .into_iter()
        .map(|p| format!("<p>{}: {}</p>", p.name, p.detailed_forecast))
        .collect::<Vec<String>>()
        .join(" ")
}

fn format_alerts(alerts: AlertResponse) -> String {
    alerts
        .features
        .into_iter()
        .map(|a| {
            format!(
                "<p>{}: {}</p>",
                a.properties.headline, a.properties.description
            )
        })
        .collect::<Vec<String>>()
        .join(" ")
}

fn build_forecast_html(lat: f64, long: f64, forecasts: String, alerts: String) -> String {
    let mut resp = String::new();
    resp.push_str(format!("<h3>Forecast for {}, {}</h3>", lat, long).as_ref());
    resp.push_str(forecasts.as_ref());
    if !alerts.is_empty() {
        resp.push_str("<h3>Alerts</h3>");
        resp.push_str(alerts.as_ref());
    }

    resp
}

fn parse_coordinates(coords: String) -> Result<(f64, f64), anyhow::Error> {
    let (lat, long) = coords
        .split_once(",")
        .ok_or_else(|| anyhow::anyhow!("Invalid coords format"))?;
    let lat = lat.trim().parse::<f64>()?;
    let long = long.trim().parse::<f64>()?;

    // Round to 3 decimals to cut down on forecast areas
    let lat = (lat * 1000.0).round() / 1000.0;
    let long = (long * 1000.0).round() / 1000.0;

    Ok((lat, long))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::types::{AlertFeature, FeatureProperties};

    #[test]
    fn test_format_forecast() {
        let periods = vec![
            ForecastPeriod {
                name: "today".to_string(),
                detailed_forecast: "Sunny".to_string(),
            },
            ForecastPeriod {
                name: "tonight".to_string(),
                detailed_forecast: "not sunny".to_string(),
            },
            ForecastPeriod {
                name: "tomorrow".to_string(),
                detailed_forecast: "cloudy".to_string(),
            },
        ];

        let result = format_forecast(periods);
        assert!(result.contains("<p>today: Sunny</p>"));
        assert!(result.contains("<p>tonight: not sunny</p>"));
        assert!(result.contains("<p>tomorrow: cloudy</p>"));
    }

    #[test]
    fn test_parse_coordinates() {
        let coords = "40.77502955625229, -73.97036692492003";
        let (lat, long) = parse_coordinates(coords.to_string()).unwrap();
        assert_eq!(lat, 40.775);
        assert_eq!(long, -73.970);
    }

    #[test]
    fn test_parse_coordinates_with_less_precision() {
        let coords = "40.77, -73.97";
        let (lat, long) = parse_coordinates(coords.to_string()).unwrap();
        assert_eq!(lat, 40.77);
        assert_eq!(long, -73.97);
    }

    #[test]
    #[should_panic]
    fn test_parse_coordinates_fails() {
        let coords = "40.77502955625229 .97036692492003";
        let (_, _) = parse_coordinates(coords.to_string()).unwrap();
    }

    #[test]
    fn test_format_alerts() {
        let alerts = AlertResponse {
            features: vec![
                AlertFeature {
                    properties: FeatureProperties {
                        headline: "Winter Storm Warning".to_string(),
                        description: "Heavy snow expected".to_string(),
                    },
                },
                AlertFeature {
                    properties: FeatureProperties {
                        headline: "Wind Advisory".to_string(),
                        description: "Gusts up to 50mph".to_string(),
                    },
                },
            ],
        };

        let result = format_alerts(alerts);
        assert!(result.contains("<p>Winter Storm Warning: Heavy snow expected</p>"));
        assert!(result.contains("<p>Wind Advisory: Gusts up to 50mph</p>"));
    }

    #[test]
    fn test_format_alerts_empty() {
        let alerts = AlertResponse { features: vec![] };
        let result = format_alerts(alerts);
        assert_eq!(result, "");
    }

    #[test]
    fn test_build_forecast_html_with_alerts() {
        let forecasts = "<p>Today: Sunny</p>".to_string();
        let alerts = "<p>Heat Advisory: Stay hydrated</p>".to_string();

        let result = build_forecast_html(40.775, -73.970, forecasts, alerts);

        assert!(result.contains("<h3>Forecast for 40.775, -73.97</h3>"));
        assert!(result.contains("<p>Today: Sunny</p>"));
        assert!(result.contains("<h3>Alerts</h3>"));
        assert!(result.contains("<p>Heat Advisory: Stay hydrated</p>"));
    }

    #[test]
    fn test_build_forecast_html_without_alerts() {
        let forecasts = "<p>Today: Sunny</p>".to_string();
        let alerts = "".to_string();

        let result = build_forecast_html(40.775, -73.970, forecasts, alerts);

        assert!(result.contains("<h3>Forecast for 40.775, -73.97</h3>"));
        assert!(result.contains("<p>Today: Sunny</p>"));
        assert!(!result.contains("<h3>Alerts</h3>"));
    }
}
