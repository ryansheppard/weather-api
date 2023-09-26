package main

import (
	"embed"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryansheppard/weather/internal/cache"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	cache.New(redisAddr)

	nws.New(config.NWSBaseURL, config.UserAgent)
	purpleair.New(config.PurpleAirBaseURL, config.PurpleAirAPIKey)

	e := echo.New()

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
	e.GET("/aqi/s/:sensorId", handlers.AQIByID)
	e.GET("/aqi/c/:coords", handlers.AQIByCoords)
	e.GET("/f/help", handlers.Help)

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
