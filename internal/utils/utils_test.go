package utils

import "testing"

func TestParseCoordinates(t *testing.T) {
	t.Parallel()
	want := &coordinates{
		Latitude:  1.0,
		Longitude: -1.0,
	}
	got := parseCoordinates("1.0, -1.0")
	if got.Latitude != want.Latitude || got.Longitude != want.Longitude {
		t.Errorf("parseCoordinates() = %v, want %v", got, want)
	}
}

func TestParseCoordinatesNoSpace(t *testing.T) {
	t.Parallel()
	want := &coordinates{
		Latitude:  1.0,
		Longitude: -1.0,
	}
	got := parseCoordinates("1.0,-1.0")
	if got.Latitude != want.Latitude || got.Longitude != want.Longitude {
		t.Errorf("parseCoordinates() = %v, want %v", got, want)
	}
}
func TestParseCoordinatesLotsOfSpace(t *testing.T) {
	t.Parallel()
	want := &coordinates{
		Latitude:  1.0,
		Longitude: -1.0,
	}
	got := parseCoordinates("  1.0 ,     -1.0")
	if got.Latitude != want.Latitude || got.Longitude != want.Longitude {
		t.Errorf("parseCoordinates() = %v, want %v", got, want)
	}
}
