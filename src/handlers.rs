use crate::{
    error::AppError,
    nws,
    state::AppState,
    types::{AlertResponse, ForecastPeriod},
};
use anyhow::Result;
use axum::{
    extract::{Path, Query, State},
    response::Html,
};
use serde::Deserialize;

#[derive(Deserialize)]
pub struct ForecastParams {
    #[serde(default, deserialize_with = "serde_bool_param")]
    short: bool,
}

fn serde_bool_param<'de, D>(deserializer: D) -> Result<bool, D::Error>
where
    D: serde::Deserializer<'de>,
{
    let s = String::deserialize(deserializer)?;
    if s.is_empty() || s == "true" {
        Ok(true)
    } else {
        Ok(false)
    }
}

pub async fn forecast(
    State(state): State<AppState>,
    Path(coords): Path<String>,
    Query(params): Query<ForecastParams>,
) -> Result<Html<String>, AppError> {
    let (lat, long) = parse_coordinates(coords)?;

    let points = nws::get_points(&state.client, &state.redis, &state.base_url, lat, long).await?;
    let points_properties = points.properties;

    let (alerts, forecast) = tokio::try_join!(
        nws::get_alerts(&state.client, &state.redis, &state.base_url, lat, long),
        nws::get_forecast(
            &state.client,
            &state.redis,
            &state.base_url,
            points_properties.grid_id,
            points_properties.grid_x,
            points_properties.grid_y,
        )
    )?;

    let forecasts = format_forecast(forecast.properties.periods, params.short);
    let alerts = format_alerts(alerts, params.short);

    let resp = build_forecast_html(lat, long, forecasts, alerts);

    Ok(Html(resp))
}

fn format_forecast(periods: Vec<ForecastPeriod>, short: bool) -> String {
    periods
        .into_iter()
        .map(|p| {
            if short {
                format!(
                    "<p>{}: {}, {}{}, {}% precip</p>",
                    p.name,
                    p.short_forecast,
                    p.temperature,
                    p.temperature_unit,
                    p.prob_of_precip.value
                )
            } else {
                format!("<p>{}: {}</p>", p.name, p.detailed_forecast)
            }
        })
        .collect::<Vec<String>>()
        .join(" ")
}

fn format_alerts(alerts: AlertResponse, short: bool) -> String {
    alerts
        .features
        .into_iter()
        .map(|a| {
            if short {
                format!("<p>{}</p>", a.properties.headline)
            } else {
                format!(
                    "<p>{}: {}</p>",
                    a.properties.headline, a.properties.description
                )
            }
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
    use crate::types::{AlertFeature, DetailedUnit, FeatureProperties};

    fn make_period(name: &str, detailed: &str, short: &str) -> ForecastPeriod {
        ForecastPeriod {
            name: name.to_string(),
            detailed_forecast: detailed.to_string(),
            short_forecast: short.to_string(),
            temperature: 72,
            temperature_unit: "F".to_string(),
            prob_of_precip: DetailedUnit { value: 10.0 },
        }
    }

    #[test]
    fn test_format_forecast() {
        let periods = vec![
            make_period("today", "Sunny", "sun"),
            make_period("tonight", "not sunny", "cloudy"),
            make_period("tomorrow", "cloudy", "cloudy"),
        ];

        let result = format_forecast(periods, false);
        assert!(result.contains("<p>today: Sunny</p>"));
        assert!(result.contains("<p>tonight: not sunny</p>"));
        assert!(result.contains("<p>tomorrow: cloudy</p>"));
    }

    #[test]
    fn test_format_short_forecast() {
        let periods = vec![
            make_period("today", "Sunny", "sun"),
            make_period("tonight", "not sunny", "cloudy"),
        ];

        let result = format_forecast(periods, true);
        assert!(result.contains("<p>today: sun, 72F, 10% precip</p>"));
        assert!(result.contains("<p>tonight: cloudy, 72F, 10% precip</p>"));
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

        let result = format_alerts(alerts, false);
        assert!(result.contains("<p>Winter Storm Warning: Heavy snow expected</p>"));
        assert!(result.contains("<p>Wind Advisory: Gusts up to 50mph</p>"));
    }

    #[test]
    fn test_format_short_alerts() {
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

        let result = format_alerts(alerts, true);
        assert!(result.contains("<p>Winter Storm Warning</p>"));
        assert!(result.contains("<p>Wind Advisory</p>"));
    }

    #[test]
    fn test_format_alerts_empty() {
        let alerts = AlertResponse { features: vec![] };
        let result = format_alerts(alerts, false);
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
