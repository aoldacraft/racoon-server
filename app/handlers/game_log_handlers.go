package handlers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"go_racoon/models"
	"log"
	"net/http"
)

func GetGameLog(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the game UUID from the URL parameter
		gameUUID := c.Param("uuid")

		// Query to select all logs for the specified game UUID, including service_name
		query := "SELECT log_id, service_name, uuid, log_text, log_time FROM game_log WHERE uuid = $1"

		rows, err := db.Query(query, gameUUID)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to fetch game logs")
		}
		defer rows.Close()

		// Create a slice to store the results
		gameLogs := []models.GameLog{}

		// Iterate through the rows and scan the results
		for rows.Next() {
			logEntry := models.GameLog{}
			err := rows.Scan(&logEntry.LogID, &logEntry.ServiceName, &logEntry.UUID, &logEntry.LogText, &logEntry.LogTime)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to fetch game logs")
			}
			gameLogs = append(gameLogs, logEntry)
		}

		// Check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to fetch game logs")
		}

		// Return the list of game logs
		return c.JSON(http.StatusOK, gameLogs)
	}
}

func CreateGameLog(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the game UUID from the URL parameter
		gameUUID := c.Param("uuid")

		// Parse the request body to extract the log text
		logRequest := new(models.GameLog)
		if err := c.Bind(logRequest); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Get the game's service name from the database
		queryServiceName := "SELECT service_name FROM game WHERE uuid = $1"
		var serviceName string
		err := db.QueryRow(queryServiceName, gameUUID).Scan(&serviceName)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to get service name")
		}

		// Insert the new game log into the database with the service name
		insertSQL := "INSERT INTO game_log (uuid, service_name, log_text, log_time) VALUES ($1, $2, $3, now())"
		_, err = db.Exec(insertSQL, gameUUID, serviceName, logRequest.LogText)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to create the game log")
		}

		return c.JSON(http.StatusCreated, "Game log created successfully")
	}
}
