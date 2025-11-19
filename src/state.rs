use redis::aio::MultiplexedConnection;
use reqwest::Client;
use url::Url;

#[derive(Clone)]
pub struct AppState {
    pub client: Client,
    pub base_url: Url,
    pub redis: Option<MultiplexedConnection>,
}
