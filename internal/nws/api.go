package nws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/patrickmn/go-cache"
)

type NWS struct {
	baseURL string
	cache   *cache.Cache
}

func NewNWS(baseURL string, cache *cache.Cache) *NWS {
	n := NWS{baseURL: baseURL, cache: cache}
	return &n
}

func (n *NWS) get(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", n.baseURL, path)
	rawBody, found := n.cache.Get(url)
	if found {
		fmt.Printf("Cache hit for %s\n", url)
		return rawBody.([]byte), nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Weatherbot, ryandsheppard95@gmail.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	err = decode(body, &point)
	if err != nil {
		return nil, err
	}

	return
}

func (n *NWS) GetForecast(point *PointResponse) (forecast *ForecastResponse, err error) {
	path := fmt.Sprintf("/gridpoints/%s/%d,%d/forecast", point.Properties.GridID, point.Properties.GridX, point.Properties.GridY)
	body, err := n.get(path)
	if err != nil {
		log.Fatal(err)
	}

	err = decode(body, &forecast)
	if err != nil {
		return nil, err
	}

	return
}

func (n *NWS) GetAlerts(lat float64, long float64) (alerts *AlertResponse, err error) {
	path := fmt.Sprintf("/alerts/active?point=%f,%f", lat, long)
	body, err := n.get(path)
	if err != nil {
		log.Fatal(err)
	}

	err = decode(body, &alerts)
	if err != nil {
		return nil, err
	}

	return
}
