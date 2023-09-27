package config

import (
	"log"
	"os"
)

type Config struct {
	NWSBaseURL       string
	UserAgent        string
	PurpleAirBaseURL string
	PurpleAirAPIKey  string
	RedisAddr        string
}

func NewConfig() *Config {
	nwsBaseURL, ok := os.LookupEnv("NWS_BASE_URl")
	if !ok {
		nwsBaseURL = "https://api.weather.gov"
	}
	userAgent, ok := os.LookupEnv("USERAGENT")
	if !ok {
		log.Fatal("USERAGENT environment variable not set")
	}
	purpleAirBaseURL, ok := os.LookupEnv("PURPLE_AIR_BASE_URL")
	if !ok {
		purpleAirBaseURL = "https://api.purpleair.com/v1"
	}
	purpleAirAPIKey, ok := os.LookupEnv("PURPLE_AIR_API_KEY")
	if !ok {
		log.Fatal("PURPLEAIRAPIKEY environment variable not set")
	}
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = ""
	}

	return &Config{
		NWSBaseURL:       nwsBaseURL,
		UserAgent:        userAgent,
		PurpleAirBaseURL: purpleAirBaseURL,
		PurpleAirAPIKey:  purpleAirAPIKey,
		RedisAddr:        redisAddr,
	}
}
