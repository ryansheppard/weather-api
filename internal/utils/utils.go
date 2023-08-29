package utils

import (
	"strconv"
	"strings"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
	coords = &Coordinates{
		Latitude:  lat,
		Longitude: long,
	}

	return
}
