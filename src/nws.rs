use crate::types::{AlertResponse, ForecastResponse, PointsResponse};
use reqwest::{Client, Error};
use url::Url;

pub async fn get_points(
    client: &Client,
    base_url: &Url,
    lat: f64,
    long: f64,
) -> Result<PointsResponse, Error> {
    let endpoint = base_url
        .join(&format!("points/{},{}", lat, long))
        .expect("Failed to construct URL");

    let response: PointsResponse = get_as_json(client, endpoint).await?;

    Ok(response)
}

pub async fn get_forecast(
    client: &Client,
    base_url: &Url,
    grid_id: String,
    grid_x: u8,
    grid_y: u8,
) -> Result<ForecastResponse, Error> {
    let endpoint = base_url
        .join(&format!(
            "gridpoints/{}/{},{}/forecast",
            grid_id, grid_x, grid_y
        ))
        .expect("Failed to construct URL");

    let response: ForecastResponse = get_as_json::<ForecastResponse>(client, endpoint).await?;

    Ok(response)
}

pub async fn get_alerts(
    client: &Client,
    base_url: &Url,
    lat: f64,
    long: f64,
) -> Result<AlertResponse, Error> {
    let endpoint = base_url
        .join(&format!("alerts/active?point={},{}", lat, long))
        .expect("Failed to construct URL");

    let response: AlertResponse = get_as_json(client, endpoint).await?;

    Ok(response)
}

async fn get_as_json<T: serde::de::DeserializeOwned>(
    client: &Client,
    endpoint: Url,
) -> Result<T, Error> {
    let response = client
        .get(endpoint.as_str())
        .send()
        .await?
        .json::<T>()
        .await?;

    Ok(response)
}
