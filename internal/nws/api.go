package nws

import (
	"fmt"

	"github.com/ryansheppard/weather/internal/cache"
	"github.com/ryansheppard/weather/internal/utils"
)

type NWS struct {
	baseURL   string
	userAgent string
	cache     *cache.Cache
}

func New(baseURL string, userAgent string, cache *cache.Cache) *NWS {
	nws := NWS{
		baseURL:   baseURL,
		userAgent: userAgent,
		cache:     cache,
	}
	return &nws
}

// Gets points from NWS weather API
func (n *NWS) GetPoints(lat float64, long float64) (*PointResponse, error) {
	var point *PointResponse
	endpoint := fmt.Sprintf("%s/points/%f,%f", n.baseURL, lat, long)
	err := n.getAndReturn(endpoint, &point)
	if err != nil {
		return nil, err
	}

	return point, nil
}

func (n *NWS) GetForecast(gridId string, gridX int, gridY int) (*ForecastResponse, error) {
	var forecast *ForecastResponse
	endpoint := fmt.Sprintf("%s/gridpoints/%s/%d,%d/forecast", n.baseURL, gridId, gridX, gridY)
	err := n.getAndReturn(endpoint, &forecast)
	if err != nil {
		return nil, err
	}

	return forecast, nil
}

func (n *NWS) GetAlerts(lat float64, long float64) (*AlertResponse, error) {
	var alerts *AlertResponse
	endpoint := fmt.Sprintf("%s/alerts/active?point=%f,%f", n.baseURL, lat, long)
	err := n.getAndReturn(endpoint, &alerts)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

func (n *NWS) getAndReturn(endpoint string, data interface{}) error {
	r := utils.NewHttpRequest(
		endpoint,
		utils.WithUserAgent(n.userAgent),
		utils.WithCaller("NWS"),
		utils.WithCache(n.cache),
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
