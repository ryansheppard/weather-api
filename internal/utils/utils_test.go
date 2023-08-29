package utils

import "testing"

func TestParseCoordinates(t *testing.T) {
	t.Parallel()
	want := &Coordinates{
		Latitude:  1.0,
		Longitude: -1.0,
	}
	got, _ := ParseCoordinates("1.0, -1.0")
	if got.Latitude != want.Latitude || got.Longitude != want.Longitude {
		t.Errorf("ParseCoordinates() = %v, want %v", got, want)
	}
}

func TestParseCoordinatesRound(t *testing.T) {
	t.Parallel()
	want := &Coordinates{
		Latitude:  1.123,
		Longitude: -1.152,
	}
	got, _ := ParseCoordinates("1.12345, -1.15235")
	if got.Latitude != want.Latitude || got.Longitude != want.Longitude {
		t.Errorf("ParseCoordinates() = %v, want %v", got, want)
	}
}

func TestParseCoordinatesNoSpace(t *testing.T) {
	t.Parallel()
	want := &Coordinates{
		Latitude:  1.0,
		Longitude: -1.0,
	}
	got, _ := ParseCoordinates("1.0,-1.0")
	if got.Latitude != want.Latitude || got.Longitude != want.Longitude {
		t.Errorf("ParseCoordinates() = %v, want %v", got, want)
	}
}
func TestParseCoordinatesLotsOfSpace(t *testing.T) {
	t.Parallel()
	want := &Coordinates{
		Latitude:  1.0,
		Longitude: -1.0,
	}
	got, _ := ParseCoordinates("  1.0 ,     -1.0")
	if got.Latitude != want.Latitude || got.Longitude != want.Longitude {
		t.Errorf("ParseCoordinates() = %v, want %v", got, want)
	}
}
