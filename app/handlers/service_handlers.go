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
		query := "SELECT service_id, service_name FROM service"
		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to retrieve services")
		}
		defer rows.Close()

		// Create a slice to store the retrieved services
		var services []models.Service

		for rows.Next() {
			var service models.Service
			if err := rows.Scan(&service.ServiceID, &service.ServiceName); err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to retrieve services")
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
		query := "SELECT service_id, service_name FROM service WHERE service_name = $1"
		var service models.Service
		err := db.QueryRow(query, serviceName).Scan(&service.ServiceID, &service.ServiceName)
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
		// Parse the request body to extract service_name
		service := new(models.Service)
		if err := c.Bind(service); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Insert the new service into the database
		insertSQL := "INSERT INTO service (service_name) VALUES ($1) RETURNING service_id"
		var serviceID int
		err := db.QueryRow(insertSQL, service.ServiceName).Scan(&serviceID)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to create the service")
		}

		// Set the generated service_id and return the created service
		service.ServiceID = serviceID
		return c.JSON(http.StatusCreated, service)
	}
}

// UpdateService updates a service's information in the database.
func UpdateService(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the service ID from the URL parameter
		serviceID := c.Param("id")

		// Parse the request body to extract the updated service name
		service := new(models.Service)
		if err := c.Bind(service); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Update the service in the database
		updateSQL := "UPDATE service SET service_name = $1 WHERE service_id = $2"
		_, err := db.Exec(updateSQL, service.ServiceName, serviceID)
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
		// Get the service ID from the URL parameter
		serviceID := c.Param("id")

		// Delete the service from the database
		deleteSQL := "DELETE FROM service WHERE service_id = $1"
		_, err := db.Exec(deleteSQL, serviceID)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to delete the service")
		}

		return c.JSON(http.StatusOK, "Service deleted successfully")
	}
}

func GetServiceGame(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the service name from the URL parameter
		serviceName := c.Param("service_name")

		// Query the database to fetch all games for the specified service_name
		query := `
		SELECT game.game_uuid, game.service_id
		FROM service
		INNER JOIN game ON service.service_id = game.service_id
		WHERE service.service_name = $1
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
			if err := rows.Scan(&game.GameUUID, &game.ServiceID); err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to retrieve games for the service")
			}

			games = append(games, game)
		}

		return c.JSON(http.StatusOK, games)
	}
}
