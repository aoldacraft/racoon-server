package handlers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"go_racoon/models"
	"log"
	"net/http"
)

// GetAllServices returns all services from the database.
func GetAllServices(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Fetch all services from the database
		query := "SELECT service_name, server_quantity, total_player FROM service"
		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to retrieve services1")
		}
		defer rows.Close()

		// Create a slice to store the retrieved services
		var services []models.Service

		for rows.Next() {
			var service models.Service
			if err := rows.Scan(&service.ServiceName, &service.ServerQuantity, &service.TotalPlayer); err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to retrieve services2")
			}
			services = append(services, service)
		}

		return c.JSON(http.StatusOK, services)
	}
}

// GetService returns a service by its service_name.
func GetService(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the service name from the URL parameter
		serviceName := c.Param("service_name")

		// Query the database to fetch the service by its service_name
		query := "SELECT service_name, server_quantity, total_player FROM service WHERE service_name = $1"
		var service models.Service
		err := db.QueryRow(query, serviceName).Scan(&service.ServiceName, &service.ServerQuantity, &service.TotalPlayer)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, "Service not found")
			} else {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to retrieve the service")
			}
		}

		return c.JSON(http.StatusOK, service)
	}
}

// CreateService creates a new service in the database.
func CreateService(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Parse the request body to extract the data
		service := new(models.Service)
		if err := c.Bind(service); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Insert the new service into the database
		insertSQL := "INSERT INTO service (service_name, server_quantity, total_player) VALUES ($1, $2, $3) RETURNING service_name"
		err := db.QueryRow(insertSQL, service.ServiceName, service.ServerQuantity, service.TotalPlayer).Scan(&service.ServiceName)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to create the service")
		}

		return c.JSON(http.StatusCreated, service)
	}
}

// UpdateService updates a service's information in the database.
func UpdateService(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the service name from the URL parameter
		serviceName := c.Param("service_name")

		// Parse the request body to extract the updated service data
		service := new(models.Service)
		if err := c.Bind(service); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Update the service in the database
		updateSQL := "UPDATE service SET service_name = $1, server_quantity = $2, total_player = $3 WHERE service_name = $4"
		_, err := db.Exec(updateSQL, service.ServiceName, service.ServerQuantity, service.TotalPlayer, serviceName)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to update the service")
		}

		return c.JSON(http.StatusOK, "Service updated successfully")
	}
}

// DeleteService deletes a service from the database.
func DeleteService(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the service name from the URL parameter
		serviceName := c.Param("service_name")

		// Delete the service from the database
		deleteSQL := "DELETE FROM service WHERE service_name = $1"
		_, err := db.Exec(deleteSQL, serviceName)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to delete the service")
		}

		return c.JSON(http.StatusOK, "Service deleted successfully")
	}
}

// GetServiceGame fetches games for a specific service from the database.
func GetServiceGame(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the service name from the URL parameter
		serviceName := c.Param("service_name")

		// Query the database to fetch all games for the specified service_name
		query := `
		SELECT game.uuid, game.service_name, game.server_ip, game.player_quantity, game.game_state, game.state_time
		FROM game
		WHERE game.service_name = $1
		`

		rows, err := db.Query(query, serviceName)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to retrieve games for the service")
		}
		defer rows.Close()

		// Create a data structure to store the games
		var games []models.Game

		for rows.Next() {
			var game models.Game
			if err := rows.Scan(&game.UUID, &game.ServiceName, &game.ServerIP, &game.PlayerQuantity, &game.GameState, &game.StateTime); err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to retrieve games for the service")
			}

			games = append(games, game)
		}

		return c.JSON(http.StatusOK, games)
	}
}
