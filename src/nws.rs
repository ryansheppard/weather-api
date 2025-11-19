use crate::types;
use reqwest::{Client, Error};
use url::Url;

pub async fn get_points(
    client: &Client,
    base_url: &Url,
    lat: f64,
    long: f64,
) -> Result<types::PointsResponse, Error> {
    let endpoint = base_url
        .join(&format!("points/{},{}", lat, long))
        .expect("Failed to construct URL");

    let response = client
        .get(endpoint.as_str())
        .send()
        .await?
        .json::<types::PointsResponse>()
        .await?;

    Ok(response)
}

pub async fn get_forecast(
    client: &Client,
    base_url: &Url,
    grid_id: String,
    grid_x: u8,
    grid_y: u8,
) -> Result<types::ForecastResponse, Error> {
    let endpoint = base_url
        .join(&format!(
            "gridpoints/{}/{},{}/forecast",
            grid_id, grid_x, grid_y
        ))
        .expect("Failed to construct URL");

    let response = client
        .get(endpoint.as_str())
        .send()
        .await?
        .json::<types::ForecastResponse>()
        .await?;

    Ok(response)
}

pub async fn get_alerts(
    client: &Client,
    base_url: &Url,
    lat: f64,
    long: f64,
) -> Result<types::AlertResponse, Error> {
    let endpoint = base_url
        .join(&format!("alerts/active?point={},{}", lat, long))
        .expect("Failed to construct URL");

    let response = client
        .get(endpoint.as_str())
        .send()
        .await?
        .json::<types::AlertResponse>()
        .await?;

    Ok(response)
}
