package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// HealthCheck is a handler for checking the health of the server.
func HealthCheck(c echo.Context) error {

	// Sample health check response
	response := "Server is healthy"

	return c.String(http.StatusOK, response)
}
