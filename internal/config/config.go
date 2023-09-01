package config

import (
	"log"
	"os"
)

type Config struct {
	BaseURL         string
	UserAgent       string
	PurpleAirAPIKey string
}

func NewConfig() *Config {
	baseUrl, ok := os.LookupEnv("BASEURL")
	if !ok {
		baseUrl = "https://api.weather.gov"
	}
	userAgent, ok := os.LookupEnv("USERAGENT")
	if !ok {
		log.Fatal("USERAGENT environment variable not set")
	}
	purpleAirAPIKey, ok := os.LookupEnv("PURPLE_AIR_API_KEY")
	if !ok {
		log.Fatal("PURPLEAIRAPIKEY environment variable not set")
	}

	return &Config{
		BaseURL:         baseUrl,
		UserAgent:       userAgent,
		PurpleAirAPIKey: purpleAirAPIKey,
	}
}
