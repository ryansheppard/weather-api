package purpleair

type SensorResponse struct {
	Sensor Sensor `json:"sensor"`
}

type Sensor struct {
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	PmTwoFiveAtm float64 `json:"pm2.5_atm"`
}
