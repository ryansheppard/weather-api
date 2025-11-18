use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct PointsProperties {
    #[serde(rename = "gridId")]
    pub grid_id: String,
    #[serde(rename = "gridX")]
    pub grid_x: u8,
    #[serde(rename = "gridY")]
    pub grid_y: u8,
}

#[derive(Deserialize, Debug)]
pub struct PointsResponse {
    pub properties: PointsProperties,
}

#[derive(Deserialize, Debug)]
pub struct ForecastPeriod {
    #[serde(rename = "detailedForecast")]
    pub detailed_forecast: String,
    pub name: String,
}

#[derive(Deserialize, Debug)]
pub struct ForecastProperties {
    pub periods: Vec<ForecastPeriod>,
}

#[derive(Deserialize, Debug)]
pub struct ForecastResponse {
    pub properties: ForecastProperties,
}

#[derive(Deserialize, Debug)]
pub struct AlertResponse {
    pub features: Vec<AlertFeature>,
}

#[derive(Deserialize, Debug)]
pub struct AlertFeature {
    pub properties: FeatureProperties,
}

#[derive(Deserialize, Debug)]
pub struct FeatureProperties {
    pub headline: String,
    pub description: String,
}
