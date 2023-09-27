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
func (n *NWS) GetPoints(lat float64, long float64) (point *PointResponse, err error) {
	endpoint := fmt.Sprintf("%s/points/%f,%f", n.baseURL, lat, long)
	n.getAndReturn(endpoint, &point)
	if err != nil {
		return
	}

	return
}

func (n *NWS) GetForecast(gridId string, gridX int, gridY int) (forecast *ForecastResponse, err error) {
	endpoint := fmt.Sprintf("%s/gridpoints/%s/%d,%d/forecast", n.baseURL, gridId, gridX, gridY)
	n.getAndReturn(endpoint, &forecast)
	if err != nil {
		return
	}

	return
}

func (n *NWS) GetAlerts(lat float64, long float64) (alerts *AlertResponse, err error) {
	endpoint := fmt.Sprintf("%s/alerts/active?point=%f,%f", n.baseURL, lat, long)
	n.getAndReturn(endpoint, &alerts)
	if err != nil {
		return
	}

	return
}

func (n *NWS) getAndReturn(endpoint string, data interface{}) (body []byte, err error) {
	r := utils.NewHttpRequest(
		endpoint,
		utils.WithUserAgent(n.userAgent),
		utils.WithCaller("NWS"),
		utils.WithCache(n.cache),
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
