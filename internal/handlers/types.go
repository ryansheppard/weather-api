package handlers

type ForecastParams struct {
	Coords     string `param:"coords"`
	Limit      int    `query:"limit"`
	Short      bool   `query:"short"`
	HideAlerts bool   `query:"hidealerts"`
}
