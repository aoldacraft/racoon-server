package routes

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"rc_app/app/handlers"
)

func RegisterServiceRoutes(e *echo.Echo, db *sql.DB) {
	serviceGroup := e.Group("/api/service")

	// Define routes for service endpoints
	serviceGroup.GET("", handlers.GetAllServices(db))
	serviceGroup.GET("/:service_name", handlers.GetService(db))
	serviceGroup.GET("/:service_name/game", handlers.GetServiceGame(db))
	serviceGroup.POST("", handlers.CreateService(db))
	serviceGroup.PUT("/:service_name", handlers.UpdateService(db))
	serviceGroup.DELETE("/:service_name", handlers.DeleteService(db))
}
