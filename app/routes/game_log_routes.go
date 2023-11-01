package routes

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"rc_app/app/handlers"
)

func RegisterGameLogRoutes(e *echo.Echo, db *sql.DB) {
	// Define routes related to game logs
	gameLogGroup := e.Group("/api/game/:uuid/log")
	gameLogGroup.GET("", handlers.GetGameLog(db))
	gameLogGroup.POST("", handlers.CreateGameLog(db))
}
