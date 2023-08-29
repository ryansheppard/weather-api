package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Help(c echo.Context) error {
	return c.Render(http.StatusOK, "help.html", nil)
}
