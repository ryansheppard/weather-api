package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryansheppard/weather/internal/noaa"
)

type coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

const baseurl = "https://api.weather.gov"

func parseCooridinates(raw string) *coordinates {
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

func getForecast(c echo.Context) error {
	rawCoords := c.Param("coords")
	coords := parseCooridinates(rawCoords)

	n := noaa.NewNOAA(baseurl)
	point, err := n.GetPoints(coords.Latitude, coords.Longitude)
	if err != nil {
		log.Fatal(err)
	}

	forecast, err := n.GetForecast(point)
	if err != nil {
		log.Fatal(err)
	}

	forecasts := make(map[string]string)
	for _, period := range forecast.Properties.Periods {
		forecasts[period.Name] = period.DetailedForecast
	}

	resp := echo.Map{
		"coords":    rawCoords,
		"forecasts": forecasts,
	}

	return c.Render(http.StatusOK, "weather.html", resp)
}
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Renderer = echoview.Default()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/forecast/:coords", getForecast)

	e.Logger.Fatal(e.Start(":1323"))
}
