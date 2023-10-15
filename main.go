package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go_racoon/app/handlers"
	"go_racoon/app/routes"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database connection
	db, err := sql.Open("postgres", "user=racoondb password=racoondb dbname=racoondb sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	// Register routes
	routes.RegisterServiceRoutes(e, db)
	routes.RegisterGameRoutes(e, db)
	routes.RegisterGameLogRoutes(e, db)

	// Health check endpoint
	e.GET("/api/health", handlers.HealthCheck)

	// Run the server
	e.Logger.Fatal(e.Start(":1323"))
}
