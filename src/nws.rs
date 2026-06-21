use crate::types::{AlertResponse, ForecastResponse, PointsResponse};
use anyhow::Result;
use redis::aio::MultiplexedConnection;
use redis::{AsyncCommands, SetExpiry, SetOptions};
use reqwest::{Client, Error};
use url::Url;

pub async fn get_points(
    client: &Client,
    redis: &Option<MultiplexedConnection>,
    base_url: &Url,
    lat: f64,
    long: f64,
) -> Result<PointsResponse, Error> {
    let endpoint = base_url
        .join(&format!("points/{},{}", lat, long))
        .expect("Failed to construct URL");

    let response: PointsResponse = get_as_json(client, redis, endpoint).await?;

    Ok(response)
}

pub async fn get_forecast(
    client: &Client,
    redis: &Option<MultiplexedConnection>,
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

    let response: ForecastResponse =
        get_as_json::<ForecastResponse>(client, redis, endpoint).await?;

    Ok(response)
}

pub async fn get_alerts(
    client: &Client,
    redis: &Option<MultiplexedConnection>,
    base_url: &Url,
    lat: f64,
    long: f64,
    hide_alerts: bool,
) -> Result<AlertResponse, Error> {
    if hide_alerts {
        return Ok(AlertResponse { features: vec![] });
    }
    let endpoint = base_url
        .join(&format!("alerts/active?point={},{}", lat, long))
        .expect("Failed to construct URL");

    let response: AlertResponse = get_as_json(client, redis, endpoint).await?;

    Ok(response)
}

async fn get_as_json<T: serde::de::DeserializeOwned + serde::Serialize>(
    client: &Client,
    redis: &Option<MultiplexedConnection>,
    endpoint: Url,
) -> Result<T, Error> {
    if let Some(ref mut con) = redis.clone()
        && let Ok(Some(cached)) = con.get::<_, Option<String>>(endpoint.as_str()).await
        && let Ok(parsed) = serde_json::from_str::<T>(&cached)
    {
        return Ok(parsed);
    }

    let response = client
        .get(endpoint.as_str())
        .send()
        .await?
        .json::<T>()
        .await?;

    if let Some(ref mut con) = redis.clone()
        && let Ok(json) = serde_json::to_string(&response)
    {
        let options = SetOptions::default().with_expiration(SetExpiry::EX(600));
        let _: () = con
            .set_options(endpoint.as_str(), json, options)
            .await
            .unwrap_or(());
    }

    Ok(response)
}
