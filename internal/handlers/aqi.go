package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ryansheppard/weather/internal/utils"
)

func AQI(c echo.Context) error {
	// sensorId := "184941"
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
	}

	return c.Render(http.StatusOK, "aqi.html", resp)
}
