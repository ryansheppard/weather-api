// Package noaa holds the types for the NOAA API
package noaa

import "fmt"

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

func (f *ForecastResponse) GetPeriods(n int) []string {
	// TODO: handle array out of bounds
	periods := f.Properties.Periods[:n]
	var periodsAsString []string
	for _, period := range periods {
		periodsAsString = append(periodsAsString, period.String())
	}

	return periodsAsString
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

func (fp *ForecastPeriod) String() string {
	var precipProbString string
	precipProb := fp.ProbabilityOfPrecipitation.Value
	if precipProb > 0 {
		precipProbString = fmt.Sprintf(", %.0f%% precip", precipProb)
	}
	return fmt.Sprintf("%s %d%s, %s %s%s", fp.Name, fp.Temperature, fp.TemperatureUnit, fp.WindDirection, fp.WindSpeed, precipProbString)
}

// DetailedUnit is a unit and value
type DetailedUnit struct {
	UnitCode string  `json:"unitCode"`
	Value    float64 `json:"value"`
}

type Forecast struct {
	Periods []string `json:"periods"`
}
