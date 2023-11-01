package routes

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"rc_app/app/handlers"
)

func RegisterGameRoutes(e *echo.Echo, db *sql.DB) {
	// Define routes related to games
	gameGroup := e.Group("/api/game")
	gameGroup.GET("", handlers.GetAllGames(db))
	gameGroup.GET("/:uuid", handlers.GetGame(db))
	gameGroup.POST("", handlers.CreateGame(db))
	gameGroup.POST("/:uuid", handlers.UpdateGame(db))
}
