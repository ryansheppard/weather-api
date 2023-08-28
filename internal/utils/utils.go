package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ryansheppard/weather/internal/nws"
)

const baseurl = "https://api.weather.gov"

type coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Params struct {
	Coords     string `param:"coords"`
	Limit      int    `query:"limit"`
	Short      bool   `query:"short"`
	HideAlerts bool   `query:"hidealerts"`
}

func parseCoordinates(raw string) (*coordinates, error) {
	rawCoords := strings.Split(raw, ",")
	lat, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[0]), 64)
	if err != nil {
		return &coordinates{}, err
	}
	long, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[1]), 64)
	if err != nil {
		return &coordinates{}, err
	}
	coords := &coordinates{
		Latitude:  lat,
		Longitude: long,
	}

	return coords, nil
}

func RenderForecast(c echo.Context) error {
	var p Params
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	coords, err := parseCoordinates(p.Coords)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	cc := c.(*ContextWithCache)
	n := nws.NewNWS(baseurl, cc.Cache)
	point, err := n.GetPoints(coords.Latitude, coords.Longitude)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	forecast, err := n.GetForecast(point.Properties.GridID, point.Properties.GridX, point.Properties.GridY)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	forecasts := []string{}
	for _, period := range forecast.Properties.Periods {
		var forecastDesc string
		if p.Short {
			precip := ""
			if period.ProbabilityOfPrecipitation.Value > 0 {
				precip = fmt.Sprintf(", %.0f%%", period.ProbabilityOfPrecipitation.Value)
			}
			forecastDesc = fmt.Sprintf("%s, %d%s%s", period.ShortForecast, period.Temperature, period.TemperatureUnit, precip)
		} else {
			forecastDesc = period.DetailedForecast
		}
		forecastString := fmt.Sprintf("%s: %s", period.Name, forecastDesc)
		forecasts = append(forecasts, forecastString)
	}

	if p.Limit > 0 && len(forecasts) > p.Limit {
		forecasts = forecasts[:p.Limit]
	}

	alertMap := make(map[string]string)
	if !p.HideAlerts {
		alerts, err := n.GetAlerts(coords.Latitude, coords.Longitude)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		for _, alert := range alerts.Features {
			alertMap[alert.Properties.Headline] = alert.Properties.Description
		}
	}

	resp := echo.Map{
		"coords":    p.Coords,
		"forecasts": forecasts,
		"alerts":    alertMap,
	}

	return c.Render(http.StatusOK, "weather.html", resp)
}

func RenderHelp(c echo.Context) error {
	return c.Render(http.StatusOK, "help.html", nil)
}
