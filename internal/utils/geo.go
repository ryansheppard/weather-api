package utils

import (
	"math"
)

const earthRadiusKm = 6371.0

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func BoundingBox(lat, long, distanceKm float64) (Coordinates, Coordinates) {
	latRadians := toRadians(lat)

	deltaLat := distanceKm / earthRadiusKm
	deltaLong := distanceKm / (earthRadiusKm * math.Cos(latRadians))

	deltaLat = deltaLat * (180.0 / math.Pi)
	deltaLong = deltaLong * (180.0 / math.Pi)

	northWest := Coordinates{
		Latitude:  lat + deltaLat,
		Longitude: long - deltaLong,
	}

	southEast := Coordinates{
		Latitude:  lat - deltaLat,
		Longitude: long + deltaLong,
	}

	return northWest, southEast
}

func HaversineDistance(lat1, long1, lat2, long2 float64) float64 {
	lat1Radians := toRadians(lat1)
	lat2Radians := toRadians(lat2)

	deltaLat := toRadians(lat2 - lat1)
	deltaLong := toRadians(long2 - long1)

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Radians)*math.Cos(lat2Radians)*math.Pow(math.Sin(deltaLong/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}
