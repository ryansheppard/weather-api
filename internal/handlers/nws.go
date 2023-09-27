package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryansheppard/weather/internal/nws"
	"github.com/ryansheppard/weather/internal/utils"
)

type NWSHandler struct {
	NWS *nws.NWS
}

func NewNWSHandler(nws *nws.NWS) *NWSHandler {
	return &NWSHandler{
		NWS: nws,
	}
}

type ForecastParams struct {
	Coords     string `param:"coords"`
	Limit      int    `query:"limit"`
	Short      bool   `query:"short"`
	HideAlerts bool   `query:"hidealerts"`
}

func (n *NWSHandler) Forecast(c echo.Context) error {
	var p ForecastParams
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	coords, err := utils.ParseCoordinates(p.Coords)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	point, err := n.NWS.GetPoints(coords.Latitude, coords.Longitude)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	forecast, err := n.NWS.GetForecast(point.Properties.GridID, point.Properties.GridX, point.Properties.GridY)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	forecasts := n.processForecast(forecast, p.Short, p.Limit)

	alertStrings := []string{}
	if !p.HideAlerts {
		alerts, err := n.NWS.GetAlerts(coords.Latitude, coords.Longitude)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		alertStrings = n.processAlerts(alerts, p.Short)
	}

	resp := echo.Map{
		"coords":    p.Coords,
		"forecasts": forecasts,
		"alerts":    alertStrings,
	}

	return c.Render(http.StatusOK, "weather.html", resp)
}

func (n *NWSHandler) processForecast(forecast *nws.ForecastResponse, short bool, limit int) (forecasts []string) {
	for _, period := range forecast.Properties.Periods {
		var forecastDesc string
		if short {
			precip := ""
			if period.ProbabilityOfPrecipitation.Value > 0 {
				precip = fmt.Sprintf(", %.0f%%", period.ProbabilityOfPrecipitation.Value)
			}

			forecastDesc = fmt.Sprintf("%s: %s, %d%s%s", period.Name, period.ShortForecast, period.Temperature, period.TemperatureUnit, precip)
		} else {
			forecastDesc = fmt.Sprintf("%s: %s", period.Name, period.DetailedForecast)
		}

		forecasts = append(forecasts, forecastDesc)
	}

	if limit > 0 && len(forecasts) > limit {
		forecasts = forecasts[:limit]
	}

	return forecasts
}

func (n *NWSHandler) processAlerts(alerts *nws.AlertResponse, short bool) (alertStrings []string) {
	for _, alert := range alerts.Features {
		var alertString string

		if short {
			alertString = fmt.Sprintf("%s", alert.Properties.Headline)
		} else {
			alertString = fmt.Sprintf("%s: %s", alert.Properties.Headline, alert.Properties.Description)
		}

		alertStrings = append(alertStrings, alertString)
	}

	return alertStrings
}
