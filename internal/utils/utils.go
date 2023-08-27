package utils

import (
	"log"
	"strconv"
	"strings"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func ParseCoordinates(raw string) *Coordinates {
	rawCoords := strings.Split(raw, ",")
	lat, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[0]), 64)
	if err != nil {
		log.Fatal(err)
	}
	long, err := strconv.ParseFloat(strings.TrimSpace(rawCoords[1]), 64)
	if err != nil {
		log.Fatal(err)
	}
	coords := &Coordinates{
		Latitude:  lat,
		Longitude: long,
	}

	return coords
}
