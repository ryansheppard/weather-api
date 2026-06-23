use crate::{
    error::AppError,
    nws,
    state::AppState,
    types::{AlertResponse, ForecastPeriod},
};
use anyhow::Result;
use axum::{
    Json,
    extract::{Path, Query, State},
    response::{Html, IntoResponse},
};
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Default)]
#[serde(rename_all = "lowercase")]
enum ResponseFormat {
    #[default]
    Html,
    Json,
}

#[derive(Deserialize)]
pub struct ForecastParams {
    #[serde(default, deserialize_with = "serde_bool_param")]
    short: bool,
    #[serde(rename = "hidealerts", default, deserialize_with = "serde_bool_param")]
    hide_alerts: bool,
    limit: Option<u8>,
    #[serde(default)]
    format: ResponseFormat,
}

#[derive(Serialize)]
pub struct ForecastResponse {
    forecasts: Vec<String>,
    alerts: Vec<String>,
}

fn serde_bool_param<'de, D>(deserializer: D) -> Result<bool, D::Error>
where
    D: serde::Deserializer<'de>,
{
    let s = String::deserialize(deserializer)?;
    if s == "1" { Ok(true) } else { Ok(false) }
}

pub async fn forecast(
    State(state): State<AppState>,
    Path(coords): Path<String>,
    Query(params): Query<ForecastParams>,
) -> Result<axum::response::Response, AppError> {
    let (lat, long) = parse_coordinates(coords)?;

    let points = nws::get_points(&state.client, &state.redis, &state.base_url, lat, long).await?;
    let points_properties = points.properties;

    let (alerts, forecast) = tokio::try_join!(
        nws::get_alerts(
            &state.client,
            &state.redis,
            &state.base_url,
            lat,
            long,
            params.hide_alerts
        ),
        nws::get_forecast(
            &state.client,
            &state.redis,
            &state.base_url,
            points_properties.grid_id,
            points_properties.grid_x,
            points_properties.grid_y,
        )
    )?;

    let forecasts = format_forecast(forecast.properties.periods, params.short, params.limit);
    let alerts = format_alerts(alerts, params.short);

    match params.format {
        ResponseFormat::Json => Ok(Json(ForecastResponse { forecasts, alerts }).into_response()),
        ResponseFormat::Html => {
            Ok(Html(build_forecast_html(lat, long, forecasts, alerts)).into_response())
        }
    }
}

fn format_forecast(periods: Vec<ForecastPeriod>, short: bool, limit: Option<u8>) -> Vec<String> {
    periods
        .into_iter()
        .take(limit.map_or(usize::MAX, |n| n as usize))
        .map(|p| {
            if short {
                format!(
                    "{}: {}, {}{}, {}% precip",
                    p.name,
                    p.short_forecast,
                    p.temperature,
                    p.temperature_unit,
                    p.prob_of_precip.value.unwrap_or(0.0) as i32
                )
            } else {
                format!("{}: {}", p.name, p.detailed_forecast)
            }
        })
        .collect::<Vec<String>>()
}

fn format_alerts(alerts: AlertResponse, short: bool) -> Vec<String> {
    alerts
        .features
        .into_iter()
        .map(|a| {
            if short {
                a.properties.headline
            } else {
                format!("{}: {}", a.properties.headline, a.properties.description)
            }
        })
        .collect::<Vec<String>>()
}

fn build_forecast_html(lat: f64, long: f64, forecasts: Vec<String>, alerts: Vec<String>) -> String {
    let mut resp = String::new();
    resp.push_str(&format!("<h3>Forecast for {}, {}</h3>", lat, long));
    let forecasts = forecasts
        .into_iter()
        .map(|f| format!("<p>{}</p>", f))
        .collect::<String>();
    resp.push_str(&forecasts);
    if !alerts.is_empty() {
        resp.push_str("<h3>Alerts</h3>");
        let alerts = alerts
            .into_iter()
            .map(|a| format!("<p>{}</p>", a))
            .collect::<String>();
        resp.push_str(&alerts);
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
            prob_of_precip: DetailedUnit { value: Some(10.0) },
        }
    }

    #[test]
    fn test_format_forecast() {
        let periods = vec![
            make_period("today", "Sunny", "sun"),
            make_period("tonight", "not sunny", "cloudy"),
            make_period("tomorrow", "cloudy", "cloudy"),
        ];

        let result = format_forecast(periods, false, None);
        assert!(result.contains(&"today: Sunny".to_string()));
        assert!(result.contains(&"tonight: not sunny".to_string()));
        assert!(result.contains(&"tomorrow: cloudy".to_string()));
    }

    #[test]
    fn test_format_limit_forecast() {
        let periods = vec![
            make_period("today", "Sunny", "sun"),
            make_period("tonight", "not sunny", "cloudy"),
            make_period("tomorrow", "cloudy", "cloudy"),
        ];

        let result = format_forecast(periods, false, Some(1));
        assert_eq!(result, vec!["today: Sunny".to_string()]);
    }

    #[test]
    fn test_format_short_forecast() {
        let periods = vec![
            make_period("today", "Sunny", "sun"),
            make_period("tonight", "not sunny", "cloudy"),
        ];

        let result = format_forecast(periods, true, None);
        assert!(result.contains(&"today: sun, 72F, 10% precip".to_string()));
        assert!(result.contains(&"tonight: cloudy, 72F, 10% precip".to_string()));
    }

    #[test]
    fn test_format_short_forecast_none_precip() {
        let periods = vec![ForecastPeriod {
            name: "today".to_string(),
            detailed_forecast: "Sunny".to_string(),
            short_forecast: "sun".to_string(),
            temperature: 72,
            temperature_unit: "F".to_string(),
            prob_of_precip: DetailedUnit { value: None },
        }];

        let result = format_forecast(periods, true, None);
        assert!(result.contains(&"today: sun, 72F, 0% precip".to_string()));
    }

    #[test]
    fn test_format_short_forecast_zero_precip() {
        let periods = vec![ForecastPeriod {
            name: "today".to_string(),
            detailed_forecast: "Sunny".to_string(),
            short_forecast: "sun".to_string(),
            temperature: 72,
            temperature_unit: "F".to_string(),
            prob_of_precip: DetailedUnit { value: Some(0.0) },
        }];

        let result = format_forecast(periods, true, None);
        assert!(result.contains(&"today: sun, 72F, 0% precip".to_string()));
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
        assert!(result.contains(&"Winter Storm Warning: Heavy snow expected".to_string()));
        assert!(result.contains(&"Wind Advisory: Gusts up to 50mph".to_string()));
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
        assert!(result.contains(&"Winter Storm Warning".to_string()));
        assert!(result.contains(&"Wind Advisory".to_string()));
    }

    #[test]
    fn test_format_alerts_empty() {
        let alerts = AlertResponse { features: vec![] };
        let result = format_alerts(alerts, false);
        assert!(result.is_empty());
    }

    #[test]
    fn test_build_forecast_html_with_alerts() {
        let forecasts = vec!["Today: Sunny".to_string()];
        let alerts = vec!["Heat Advisory: Stay hydrated".to_string()];

        let result = build_forecast_html(40.775, -73.970, forecasts, alerts);

        assert!(result.contains("<h3>Forecast for 40.775, -73.97</h3>"));
        assert!(result.contains("<p>Today: Sunny</p>"));
        assert!(result.contains("<h3>Alerts</h3>"));
        assert!(result.contains("<p>Heat Advisory: Stay hydrated</p>"));
    }

    #[test]
    fn test_build_forecast_html_without_alerts() {
        let forecasts = vec!["Today: Sunny".to_string()];
        let alerts = vec![];

        let result = build_forecast_html(40.775, -73.970, forecasts, alerts);

        assert!(result.contains("<h3>Forecast for 40.775, -73.97</h3>"));
        assert!(result.contains("<p>Today: Sunny</p>"));
        assert!(!result.contains("<h3>Alerts</h3>"));
    }
}
