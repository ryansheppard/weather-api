package purpleair

import (
	"fmt"

	"github.com/ryansheppard/weather/internal/cache"
	"github.com/ryansheppard/weather/internal/utils"
)

type PurpleAir struct {
	baseURL string
	apiKey  string
	cache   *cache.Cache
}

func New(baseURL string, apiKey string, cache *cache.Cache) *PurpleAir {
	p := PurpleAir{
		baseURL: baseURL,
		apiKey:  apiKey,
		cache:   cache,
	}

	return &p
}

func (p *PurpleAir) GetSensor(sensorId string) (*SensorResponse, error) {
	var sensor *SensorResponse
	endpoint := fmt.Sprintf("%s/sensors/%s?fields=name,latitude,longitude,pm2.5_atm", p.baseURL, sensorId)
	err := p.getAndReturn(endpoint, &sensor)
	if err != nil {
		return nil, err
	}
	return sensor, nil
}

func (p *PurpleAir) ListSensors(lat, long, distanceKm float64) (*ListSensorsResponse, error) {
	var sensors *ListSensorsResponse
	northWest, southEast := utils.BoundingBox(lat, long, distanceKm)

	endpoint := fmt.Sprintf("%s/sensors?fields=name,latitude,longitude,pm2.5_atm&nwlng=%f&nwlat=%f&selng=%f&selat=%f&location_type=0", p.baseURL, northWest.Longitude, northWest.Latitude, southEast.Longitude, southEast.Latitude)
	err := p.getAndReturn(endpoint, &sensors)
	if err != nil {
		return nil, err
	}
	return sensors, err
}

func (p *PurpleAir) getAndReturn(endpoint string, data interface{}) error {
	headers := make(map[string]string)
	headers["X-API-Key"] = p.apiKey

	r := utils.NewHttpRequest(
		endpoint,
		utils.WithCaller("PurpleAir"),
		utils.WithHeaders(headers),
		utils.WithCache(p.cache),
	)
	body, err := r.Get()
	if err != nil {
		return err
	}
	err = utils.Decode(body, &data)
	if err != nil {
		return err
	}
	return nil
}
