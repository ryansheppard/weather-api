package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	pa "github.com/ryansheppard/weather/internal/purpleair"
	"github.com/ryansheppard/weather/internal/utils"
)

func AQIByID(c echo.Context) error {
	sensorId := c.Param("sensorId")

	cc := c.(*ContextWithAPIs)
	purpleair := cc.PurpleAir

	sensor, err := purpleair.GetSensor(sensorId)

	pmtwofive := sensor.Sensor.PmTwoFiveAtm
	coords := fmt.Sprintf("%f, %f", sensor.Sensor.Latitude, sensor.Sensor.Longitude)

	aqi, desc := utils.CalculateAQI(pmtwofive)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	resp := echo.Map{
		"aqi":    strconv.Itoa(aqi),
		"desc":   desc,
		"coords": coords,
		"name":   sensor.Sensor.Name,
	}

	return c.Render(http.StatusOK, "aqi.html", resp)
}

func AQIByCoords(c echo.Context) error {
	rawCoords := c.Param("coords")
	coords, err := utils.ParseCoordinates(rawCoords)

	cc := c.(*ContextWithAPIs)
	purpleair := cc.PurpleAir

	sensors, err := purpleair.ListSensors(coords.Latitude, coords.Longitude, 15)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if len(sensors.Data) == 0 {
		resp := echo.Map{
			"coords": rawCoords,
		}
		return c.Render(http.StatusNotFound, "aqi_err.html", resp)
	}

	sensorsByDistance := make(map[float64]*pa.Sensor)
	for _, sensor := range sensors.Data {
		s := pa.NewSensor(
			sensor.([]interface{})[0].(float64),
			sensor.([]interface{})[1].(string),
			sensor.([]interface{})[2].(float64),
			sensor.([]interface{})[3].(float64),
			sensor.([]interface{})[4].(float64),
		)

		distance := utils.HaversineDistance(coords.Latitude, coords.Longitude, s.Latitude, s.Longitude)
		sensorsByDistance[distance] = s
	}

	keys := make([]float64, 0, len(sensorsByDistance))
	for k := range sensorsByDistance {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	closestSensor := sensorsByDistance[keys[0]]

	pmtwofive := closestSensor.PmTwoFiveAtm
	sensorCoords := fmt.Sprintf("%f, %f", closestSensor.Latitude, closestSensor.Longitude)

	aqi, desc := utils.CalculateAQI(pmtwofive)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	resp := echo.Map{
		"aqi":    strconv.Itoa(aqi),
		"desc":   desc,
		"coords": sensorCoords,
		"name":   closestSensor.Name,
	}

	return c.Render(http.StatusOK, "aqi.html", resp)
}
