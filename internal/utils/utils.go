package utils

import (
	"fmt"
	"log"
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

func parseCoordinates(raw string) *coordinates {
	rawCoords := strings.Split(raw, ",")
	lat, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[0]), 64)
	if err != nil {
		log.Fatal(err)
	}
	long, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[1]), 64)
	if err != nil {
		log.Fatal(err)
	}
	coords := &coordinates{
		Latitude:  lat,
		Longitude: long,
	}

	return coords
}

func GetForecast(c echo.Context) error {
	var p Params
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	coords := parseCoordinates(p.Coords)

	cc := c.(*ContextWithCache)
	n := nws.NewNWS(baseurl, cc.Cache)
	point, err := n.GetPoints(coords.Latitude, coords.Longitude)
	if err != nil {
		log.Fatal(err)
	}

	forecast, err := n.GetForecast(point)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
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

func GetHelp(c echo.Context) error {
	return c.Render(http.StatusOK, "help.html", nil)
}
