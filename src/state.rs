use reqwest::Client;
use url::Url;
use worker::kv;

#[derive(Clone)]
pub struct AppState {
    pub client: Client,
    pub base_url: Url,
    pub kv: kv::KvStore,
}
