package routes

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"go_racoon/app/handlers"
)

func RegisterServiceRoutes(e *echo.Echo, db *sql.DB) {
	serviceGroup := e.Group("/api/service")

	// Define routes for service endpoints
	serviceGroup.GET("", handlers.GetAllServices(db))
	serviceGroup.GET("/:service_name", handlers.GetService(db))
	serviceGroup.GET("/:service_name/game", handlers.GetServiceGame(db))
	serviceGroup.POST("", handlers.CreateService(db))
	serviceGroup.PUT("/:id", handlers.UpdateService(db))
	serviceGroup.DELETE("/:id", handlers.DeleteService(db))
}
