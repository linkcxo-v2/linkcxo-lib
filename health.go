package linkcxo

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthController interface
type HealthController struct{}

// Status - return status
func (h HealthController) Status() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(http.StatusOK, "Working!")
	}
}
