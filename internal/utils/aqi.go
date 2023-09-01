package utils

type aqiBreakpoint struct {
	LowConc  float64
	HighConc float64
	LowAQI   int
	HighAQI  int
}

func CalculateAQI(concentration float64) int {
	breakpoints := []aqiBreakpoint{
		{0.0, 12.0, 0, 50},
		{12.1, 35.4, 51, 100},
		{35.5, 55.4, 101, 150},
		{55.5, 150.4, 151, 200},
		{150.5, 250.4, 201, 300},
		{250.5, 350.4, 301, 400},
		{350.5, 500.4, 401, 500},
	}

	for _, bp := range breakpoints {
		if concentration >= bp.LowConc && concentration <= bp.HighConc {
			return int((float64(bp.HighAQI-bp.LowAQI)/(bp.HighConc-bp.LowConc))*(concentration-bp.LowConc) + float64(bp.LowAQI))
		}
	}

	// Return -1 if concentration doesn't fit into any range (this shouldn't generally happen unless the input is out of expected bounds)
	return -1
}
