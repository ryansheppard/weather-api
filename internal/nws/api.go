package nws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NWS struct {
	baseURL string
}

func NewNWS(baseURL string) *NWS {
	n := NWS{baseURL: baseURL}
	return &n
}

func (n *NWS) get(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", n.baseURL, path)

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

	return ioutil.ReadAll(resp.Body)
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

func (n *NWS) GetAlerts(state string) {
	path := fmt.Sprintf("/alerts/active?area=%s", state)
	body, err := n.get(path)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}