package utils

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ryansheppard/weather/internal/noaa"
)

const baseurl = "https://api.weather.gov"

type coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
	var err error

	rawCoords := c.Param("coords")
	limit := c.QueryParam("limit")

	maxPeriods := 0 // Ignore if 0
	if limit != "" {
		maxPeriods, err = strconv.Atoi(limit)
		if err != nil {
			log.Fatal(err)
		}
	}

	coords := parseCoordinates(rawCoords)

	n := noaa.NewNOAA(baseurl)
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
		forecastString := fmt.Sprintf("%s: %s", period.Name, period.DetailedForecast)
		forecasts = append(forecasts, forecastString)
	}

	if maxPeriods > 0 && len(forecasts) > maxPeriods {
		forecasts = forecasts[:maxPeriods]
	}

	resp := echo.Map{
		"coords":    rawCoords,
		"forecasts": forecasts,
	}

	return c.Render(http.StatusOK, "weather.html", resp)
}
