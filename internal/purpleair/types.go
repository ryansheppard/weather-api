package purpleair

type SensorResponse struct {
	Sensor Sensor `json:"sensor"`
}

type ListSensorsResponse struct {
	Data []interface{} `json:"data"`
}

type Sensor struct {
	SensorIndex  float64 `json:"sensor_index"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	PmTwoFiveAtm float64 `json:"pm2.5_atm"`
}

func NewSensor(sensorIndex float64, name string, latitude float64, longitude float64, pmTwoFiveAtm float64) *Sensor {
	s := Sensor{
		SensorIndex:  sensorIndex,
		Name:         name,
		Latitude:     latitude,
		Longitude:    longitude,
		PmTwoFiveAtm: pmTwoFiveAtm,
	}
	return &s
}
