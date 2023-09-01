package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ryansheppard/weather/internal/nws"
)

type ContextWithNWS struct {
	echo.Context
	NWS *nws.NWS
}
