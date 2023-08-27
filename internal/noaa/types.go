// Package noaa holds the types for the NOAA API
package noaa

// PointResponse is the JSON response from the NOAA point api
type PointResponse struct {
	Properties PointProperties `json:"properties"`
}

// PointProperties is the properties key in the PointResponse
type PointProperties struct {
	GridID           string           `json:"gridId"`
	GridX            int              `json:"gridX"`
	GridY            int              `json:"gridY"`
	RelativeLocation RelativeLocation `json:"relativeLocation"`
}

// RelativeLocation is the relativeLocation key in the PointProperties
type RelativeLocation struct {
	Properties LocationProperties `json:"properties"`
}

// LocationProperties is the LocationProperties key in relativeLocation
type LocationProperties struct {
	State string `json:"state"`
}

// ForecastResponse is the JSON response from the forecast API
type ForecastResponse struct {
	Properties ForecastProperties `json:"properties"`
}

// ForecastProperties is the properties key in the ForecastResponse
type ForecastProperties struct {
	Periods []ForecastPeriod `json:"periods"`
}

// ForecastPeriod is a period in the forecast
type ForecastPeriod struct {
	Number                     int          `json:"number"`
	Name                       string       `json:"name"`
	StartTime                  string       `json:"startTime"`
	EndTime                    string       `json:"endTime"`
	isDaytime                  bool         `json:"isDaytime"`
	Temperature                int          `json:"temperature"`
	TemperatureUnit            string       `json:"temperatureUnit"`
	TemperatureTrend           string       `json:"temperatureTrend"`
	ProbabilityOfPrecipitation DetailedUnit `json:"probabilityOfPrecipitation"`
	Dewpoint                   DetailedUnit `json:"dewpoint"`
	RelativeHumidity           DetailedUnit `json:"relativeHumidity"`
	WindSpeed                  string       `json:"windSpeed"`
	WindDirection              string       `json:"windDirection"`
	ShortForecast              string       `json:"shortForecast"`
	DetailedForecast           string       `json:"detailedForecast"`
}

// DetailedUnit is a unit and value
type DetailedUnit struct {
	UnitCode string  `json:"unitCode"`
	Value    float64 `json:"value"`
}
