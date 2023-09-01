package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ryansheppard/weather/internal/nws"
	"github.com/ryansheppard/weather/internal/purpleair"
)

type ContextWithAPIs struct {
	echo.Context
	NWS       *nws.NWS
	PurpleAir *purpleair.PurpleAir
}
