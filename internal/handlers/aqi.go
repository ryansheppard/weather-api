package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryansheppard/weather/internal/utils"
)

func AQI(c echo.Context) error {
	sensorId := "184941"
	cc := c.(*ContextWithAPIs)
	purpleair := cc.PurpleAir
	resp, err := purpleair.GetSensor(sensorId)
	pmtwofive := resp.Sensor.PmTwoFiveAtm
	fmt.Println(utils.CalculateAQI(pmtwofive))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	return c.String(http.StatusOK, "Hello, World!")
}
