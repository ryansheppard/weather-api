use axum::{Router, routing::get};
use tower_service::Service;
use worker::*;

use anyhow::Result;
use reqwest::Client;
use url::Url;

mod error;
mod handlers;
mod nws;
mod state;
mod types;

const NWS_BASE_URL: &str = "https://api.weather.gov";

#[event(fetch)]
async fn fetch(
    req: HttpRequest,
    env: Env,
    _ctx: worker::Context,
) -> Result<axum::http::Response<axum::body::Body>> {
    let user_agent = env.var("USER_AGENT")?.to_string();

    let state = state::AppState {
        client: Client::builder().user_agent(user_agent).build()?,
        base_url: Url::parse(NWS_BASE_URL)?,
        kv: env.kv("weather")?,
    };

    let mut router = Router::new()
        .route("/f/{coords}", get(handlers::forecast))
        .with_state(state);

    Ok(router.call(req).await?)
}
