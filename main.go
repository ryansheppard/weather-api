package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryansheppard/weather/internal/noaa"
	"github.com/ryansheppard/weather/internal/utils"
)

const baseurl = "https://api.weather.gov"

func getForecast(c echo.Context) error {
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

	coords := utils.ParseCoordinates(rawCoords)

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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Renderer = echoview.Default()

	e.GET("/f/:coords", getForecast)

	e.Logger.Fatal(e.Start(":1323"))
}
