package web

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func initHealthRoutes(a *Web, e *echo.Echo) {
	e.GET("/health", a.HealthGet)
}

func (a *Web) HealthGet(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
