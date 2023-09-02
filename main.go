package main

import (
	"embed"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"github.com/ryansheppard/weather/internal/config"
	"github.com/ryansheppard/weather/internal/handlers"
	"github.com/ryansheppard/weather/internal/nws"
	"github.com/ryansheppard/weather/internal/purpleair"
)

//go:embed views/*
var views embed.FS

func embeddedFH(config goview.Config, tmpl string) (string, error) {
	path := filepath.Join(config.Root, tmpl)
	bytes, err := views.ReadFile(path + config.Extension)
	return string(bytes), err
}

func main() {
	config := config.NewConfig()

	nwsCache := cache.New(5*time.Minute, 10*time.Minute)
	nws := nws.NewNWS(config.NWSBaseURL, config.UserAgent, nwsCache)

	purpleairCache := cache.New(5*time.Minute, 60*time.Minute)
	purpleair := purpleair.NewPurpleAir(config.PurpleAirBaseURL, config.PurpleAirAPIKey, purpleairCache)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handlers.ContextWithAPIs{c, nws, purpleair}
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
	e.Use(echoprometheus.NewMiddleware("weather"))

	renderer := echoview.Default()
	renderer.SetFileHandler(embeddedFH)
	e.Renderer = renderer

	e.GET("/f/:coords", handlers.Forecast)
	e.GET("f/aqi/s/:sensorId", handlers.AQIByID)
	e.GET("f/aqi/c/:coords", handlers.AQIByCoords)
	e.GET("f/help", handlers.Help)

	// Serve prometheus metrics on a different port
	go func() {
		metrics := echo.New()
		metrics.GET("/metrics", echoprometheus.NewHandler())
		if err := metrics.Start(":1324"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	e.Logger.Fatal(e.Start(":1323"))
}
