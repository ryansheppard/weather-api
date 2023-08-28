package nws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	getsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_get_calls_total",
		Help: "The total number of processed events",
	})

	cacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_cache_hits_total",
		Help: "The total number of cache hits",
	})

	cacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_cache_misses_total",
		Help: "The total number of cache misses",
	})
)

type NWS struct {
	baseURL string
	cache   *cache.Cache
}

func NewNWS(baseURL string, cache *cache.Cache) *NWS {
	n := NWS{
		baseURL: baseURL,
		cache:   cache,
	}
	return &n
}

func (n *NWS) get(path string) ([]byte, error) {
	getsProcessed.Inc()
	url := fmt.Sprintf("%s%s", n.baseURL, path)
	rawBody, found := n.cache.Get(url)
	if found {
		cacheHits.Inc()
		fmt.Printf("Cache hit for %s\n", url)
		return rawBody.([]byte), nil
	}

	cacheMisses.Inc()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []byte{}, err
	}

	req.Header.Set("User-Agent", "Weatherbot, ryandsheppard95@gmail.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	n.cache.Set(url, body, cache.DefaultExpiration)

	return body, nil
}

func decode(body []byte, v interface{}) error {
	return json.Unmarshal(body, &v)
}

// Gets points from NWS weather API
func (n *NWS) GetPoints(lat float64, long float64) (point *PointResponse, err error) {
	path := fmt.Sprintf("/points/%f,%f", lat, long)
	body, err := n.get(path)

	if err != nil {
		return
	}

	err = decode(body, &point)
	if err != nil {
		return
	}

	return
}

func (n *NWS) GetForecast(gridId string, gridX int, gridY int) (forecast *ForecastResponse, err error) {
	path := fmt.Sprintf("/gridpoints/%s/%d,%d/forecast", gridId, gridX, gridY)
	body, err := n.get(path)
	if err != nil {
		return
	}

	err = decode(body, &forecast)
	if err != nil {
		return
	}

	return
}

func (n *NWS) GetAlerts(lat float64, long float64) (alerts *AlertResponse, err error) {
	path := fmt.Sprintf("/alerts/active?point=%f,%f", lat, long)
	body, err := n.get(path)
	if err != nil {
		return
	}

	err = decode(body, &alerts)
	if err != nil {
		return
	}

	return
}
