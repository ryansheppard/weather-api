use anyhow::{Context, Result};
use axum::{Router, routing::get};
use log::info;
use reqwest::Client;
use std::env;
// use tower_http::compression::CompressionLayer;
use tower_http::compression::CompressionLayer;
use url::Url;

mod error;
mod handlers;
mod nws;
mod state;
mod types;

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();

    let user_agent = env::var("USER_AGENT").context("USER_AGENT env var must be set")?;

    let base_url =
        env::var("NWS_BASE_URL").unwrap_or_else(|_| "https://api.weather.gov".to_string());

    let redis_con = if let Ok(redis_url) = env::var("REDIS_URL") {
        let client = redis::Client::open(redis_url).context("Failed to create redis client")?;
        info!("Using redis");
        Some(
            client
                .get_multiplexed_async_connection()
                .await
                .context("Failed to create redis connection")?,
        )
    } else {
        info!("REDIS_URL not set, skipping caching");
        None
    };

    let state = state::AppState {
        client: Client::builder().user_agent(user_agent).build()?,
        base_url: Url::parse(&base_url).context("Failed to parse base URL")?,
        redis: redis_con,
    };

    let app = Router::new()
        .route("/f/{coords}", get(handlers::forecast))
        .layer(CompressionLayer::new())
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .context("Failed to create tokio listener")?;
    axum::serve(listener, app)
        .await
        .context("axum failed to serve")?;

    Ok(())
}
