package purpleair

import (
	"fmt"

	"github.com/patrickmn/go-cache"
	"github.com/ryansheppard/weather/internal/utils"
)

type PurpleAir struct {
	baseURL string
	apiKey  string
	cache   *cache.Cache
}

func NewPurpleAir(baseURL string, apiKey string, cache *cache.Cache) *PurpleAir {
	p := PurpleAir{
		baseURL: baseURL,
		apiKey:  apiKey,
		cache:   cache,
	}
	return &p
}

func (p *PurpleAir) GetSensor(sensorId string) (sensor *SensorResponse, err error) {
	endpoint := fmt.Sprintf("%s/sensors/%s?fields=name,latitude,longitude,pm2.5_atm", p.baseURL, sensorId)
	getAndReturn(endpoint, p, &sensor)
	if err != nil {
		return
	}
	return
}

func (p *PurpleAir) ListSensors(lat, long, distanceKm float64) (sensors *ListSensorsResponse, err error) {
	northWest, southEast := utils.BoundingBox(lat, long, distanceKm)

	endpoint := fmt.Sprintf("%s/sensors?fields=name,latitude,longitude,pm2.5_atm&nwlng=%f&nwlat=%f&selng=%f&selat=%f&location_type=0", p.baseURL, northWest.Longitude, northWest.Latitude, southEast.Longitude, southEast.Latitude)
	getAndReturn(endpoint, p, &sensors)

	if err != nil {
		return
	}
	return
}

func getAndReturn(endpoint string, p *PurpleAir, data interface{}) (body []byte, err error) {
	headers := make(map[string]string)
	headers["X-API-Key"] = p.apiKey

	r := utils.NewHttpRequest(
		endpoint,
		utils.WithCache(p.cache),
		utils.WithCaller("PurpleAir"),
		utils.WithHeaders(headers),
	)
	body, err = r.Get()
	if err != nil {
		return
	}
	err = utils.Decode(body, &data)
	if err != nil {
		return
	}
	return
}
