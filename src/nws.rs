use crate::types::{AlertResponse, ForecastResponse, PointsResponse};
use anyhow::{Context, Result};
use reqwest::Client;
use serde::de::DeserializeOwned;
use std::time::Duration;
use url::Url;
use worker::kv::KvStore;

const FORECAST_TTL: u64 = 600;
const ALERTS_TTL: u64 = 60;
const NWS_REQUEST_TIMEOUT: Duration = Duration::from_secs(10);

pub struct Nws<'a> {
    client: &'a Client,
    kv: &'a KvStore,
    base_url: &'a Url,
}

impl<'a> Nws<'a> {
    pub fn new(client: &'a Client, kv: &'a KvStore, base_url: &'a Url) -> Self {
        Self {
            client,
            kv,
            base_url,
        }
    }

    pub async fn get_points(&self, lat: f64, long: f64) -> Result<PointsResponse> {
        let endpoint = self.base_url.join(&format!("points/{lat},{long}"))?;
        self.get_as_json(endpoint, FORECAST_TTL).await
    }

    pub async fn get_forecast(
        &self,
        grid_id: String,
        grid_x: u16,
        grid_y: u16,
    ) -> Result<ForecastResponse> {
        let endpoint = self
            .base_url
            .join(&format!("gridpoints/{grid_id}/{grid_x},{grid_y}/forecast"))?;
        self.get_as_json(endpoint, FORECAST_TTL).await
    }

    pub async fn get_alerts(
        &self,
        lat: f64,
        long: f64,
        hide_alerts: bool,
    ) -> Result<AlertResponse> {
        if hide_alerts {
            return Ok(AlertResponse { features: vec![] });
        }
        let endpoint = self
            .base_url
            .join(&format!("alerts/active?point={lat},{long}"))?;
        self.get_as_json(endpoint, ALERTS_TTL).await
    }

    async fn get_as_json<T: DeserializeOwned>(&self, endpoint: Url, ttl: u64) -> Result<T> {
        if let Ok(Some(cached)) = self.kv.get(endpoint.as_str()).text().await
            && let Ok(parsed) = serde_json::from_str::<T>(&cached)
        {
            return Ok(parsed);
        }

        let body = self
            .client
            .get(endpoint.as_str())
            .timeout(NWS_REQUEST_TIMEOUT)
            .send()
            .await?
            .error_for_status()?
            .text()
            .await?;
        let parsed = serde_json::from_str::<T>(&body)
            .with_context(|| format!("Invalid NWS JSON from {endpoint}"))?;

        if let Ok(put) = self.kv.put(endpoint.as_str(), body) {
            let _ = put.expiration_ttl(ttl).execute().await;
        }

        Ok(parsed)
    }
}
