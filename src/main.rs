use axum::{Router, routing::get};
use reqwest::Client;
use reqwest::Error;
use std::env;
use url::Url;

mod handlers;
mod nws;
mod state;
mod types;

#[tokio::main]
async fn main() -> Result<(), Error> {
    let user_agent = env::var("USER_AGENT").unwrap();

    let base_url =
        env::var("NWS_BASE_URL").unwrap_or_else(|_| "https://api.weather.gov".to_string());

    let redis_url = env::var("REDIS_URL").unwrap_or_else(|_| "redis://127.0.0.1/".to_string());
    let client = redis::Client::open(redis_url).unwrap();
    let redis_con = client.get_multiplexed_async_connection().await.unwrap();

    let state = state::AppState {
        client: Client::builder().user_agent(user_agent).build()?,
        base_url: Url::parse(&base_url).unwrap(),
        redis: redis_con,
    };

    let app = Router::new()
        .route("/f/{coords}", get(handlers::forecast))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();

    Ok(())
}
