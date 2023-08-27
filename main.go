package main

import (
	"log"
	"net/http"

	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryansheppard/weather/internal/noaa"
	"github.com/ryansheppard/weather/internal/utils"
)

const baseurl = "https://api.weather.gov"

func getForecast(c echo.Context) error {
	rawCoords := c.Param("coords")
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

	e.GET("/f/:coords", getForecast)

	e.Logger.Fatal(e.Start(":1323"))
}
