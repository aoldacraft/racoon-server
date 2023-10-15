package routes

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"go_racoon/app/handlers"
)

func RegisterGameLogRoutes(e *echo.Echo, db *sql.DB) {
	// Define routes related to game logs
	gameLogGroup := e.Group("/api/game/:game_uuid/log")
	gameLogGroup.GET("", handlers.GetGameLog(db))
	gameLogGroup.POST("", handlers.CreateGameLog(db))
}
