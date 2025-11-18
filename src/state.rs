use reqwest::Client;
use url::Url;

#[derive(Clone)]
pub struct AppState {
    pub client: Client,
    pub base_url: Url,
}
