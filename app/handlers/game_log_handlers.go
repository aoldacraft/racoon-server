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
		// Get the game_uuid from the URL parameter
		gameUUID := c.Param("game_uuid")

		// Query to select all logs for the specified game_uuid
		query := "SELECT log_id, game_uuid, log_text FROM game_log WHERE game_uuid = $1"

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
			err := rows.Scan(&logEntry.LogID, &logEntry.GameUUID, &logEntry.LogText)
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
		// Get the game_uuid from the URL parameter
		gameUUID := c.Param("game_uuid")

		// Parse the request body to extract the log text
		type LogRequest struct {
			LogText string `json:"log_text"`
		}

		logRequest := new(LogRequest)
		if err := c.Bind(logRequest); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Insert the new game log into the database
		insertSQL := "INSERT INTO game_log (game_uuid, log_text) VALUES ($1, $2)"
		_, err := db.Exec(insertSQL, gameUUID, logRequest.LogText)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to create the game log")
		}

		return c.JSON(http.StatusCreated, "Game log created successfully")
	}
}
