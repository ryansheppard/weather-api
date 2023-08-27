package main

import (
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryansheppard/weather/internal/utils"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Renderer = echoview.Default()

	e.GET("/f/:coords", utils.GetForecast)

	e.Logger.Fatal(e.Start(":1323"))
}
