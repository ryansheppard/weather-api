package main

import (
	"time"

	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"github.com/ryansheppard/weather/internal/utils"
)

func main() {
	memCache := cache.New(5*time.Minute, 10*time.Minute)
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &utils.ContextWithCache{c, memCache}
			return next(cc)
		}
	})

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(2)))
	e.Use(middleware.Recover())

	e.Renderer = echoview.Default()

	e.GET("/f/:coords", utils.GetForecast)
	e.GET("f/help", utils.GetHelp)

	e.Logger.Fatal(e.Start(":1323"))
}
