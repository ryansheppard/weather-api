package utils

import (
	"math"
	"strconv"
	"strings"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Taken from https://gosamples.dev/round-float/
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func ParseCoordinates(raw string) (coords *Coordinates, err error) {
	rawCoords := strings.Split(raw, ",")

	lat, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[0]), 64)
	if err != nil {
		return
	}

	long, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[1]), 64)
	if err != nil {
		return
	}

	latRounded := roundFloat(lat, 3)
	longRounded := roundFloat(long, 3)

	coords = &Coordinates{
		Latitude:  latRounded,
		Longitude: longRounded,
	}

	return
}
