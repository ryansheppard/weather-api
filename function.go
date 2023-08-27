// Package function does stuff
package function

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/ryansheppard/weather/internal/noaa"
)

type coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

const baseurl = "https://api.weather.gov"

func init() {
	functions.HTTP("NOAA", getForecast)
}

func parseCooridinates(raw string) *coordinates {
	rawCoords := strings.Split(raw, ",")
	lat, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[0]), 64)
	if err != nil {
		log.Fatal(err)
	}
	long, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[1]), 64)
	if err != nil {
		log.Fatal(err)
	}
	coords := &coordinates{
		Latitude:  lat,
		Longitude: long,
	}

	return coords
}

func getForecast(w http.ResponseWriter, r *http.Request) {
	var coords coordinates
	err = json.NewDecoder(r.Body).Decode(&coords)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	n := noaa.NewNOAA(baseurl)
	point, err := n.GetPoints(coords.Latitude, coords.Longitude)
	if err != nil {
		log.Fatal(err)
	}

	forecast, err := n.GetForecast(point)
	if err != nil {
		log.Fatal(err)
	}

	periods := forecast.GetPeriods(3)
	response := &noaa.Forecast{
		Periods: periods,
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(encoded)
}
