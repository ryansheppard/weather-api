package main

import (
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryansheppard/weather/internal/utils"
)

func main() {
	e := echo.New()

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(2)))

	e.Renderer = echoview.Default()

	e.GET("/f/:coords", utils.GetForecast)
	e.GET("f/help", utils.GetHelp)

	e.Logger.Fatal(e.Start(":1323"))
}
