use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct PointsProperties {
    #[serde(rename = "gridId")]
    pub grid_id: String,
    #[serde(rename = "gridX")]
    pub grid_x: u16,
    #[serde(rename = "gridY")]
    pub grid_y: u16,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct PointsResponse {
    pub properties: PointsProperties,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ForecastPeriod {
    #[serde(rename = "detailedForecast")]
    pub detailed_forecast: String,
    #[serde(rename = "shortForecast")]
    pub short_forecast: String,
    #[serde(rename = "probabilityOfPrecipitation")]
    pub prob_of_precip: DetailedUnit,
    pub temperature: i16,
    #[serde(rename = "temperatureUnit")]
    pub temperature_unit: String,
    pub name: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct DetailedUnit {
    pub value: Option<f64>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ForecastProperties {
    pub periods: Vec<ForecastPeriod>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ForecastResponse {
    pub properties: ForecastProperties,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct AlertResponse {
    pub features: Vec<AlertFeature>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct AlertFeature {
    pub properties: FeatureProperties,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct FeatureProperties {
    pub headline: String,
    pub description: String,
}
