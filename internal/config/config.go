package config

import (
	"log"
	"os"
)

type Config struct {
	BaseURL   string
	UserAgent string
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

	return &Config{
		BaseURL:   baseUrl,
		UserAgent: userAgent,
	}
}
