package handlers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"rc_app/models"
	"strconv"
	"strings"
)

// GetAllGames fetches all games from the database.
func GetAllGames(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Query to select all games
		query := "SELECT uuid, service_name, server_ip, player_quantity, game_state, state_time FROM game"

		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to fetch games")
		}
		defer rows.Close()

		// Create a slice to store the results
		games := []models.Game{}

		// Iterate through the rows and scan the results
		for rows.Next() {
			game := models.Game{}
			err := rows.Scan(&game.UUID, &game.ServiceName, &game.ServerIP, &game.PlayerQuantity, &game.GameState, &game.StateTime)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to fetch games")
			}
			games = append(games, game)
		}

		// Check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to fetch games")
		}

		// Return the list of games
		return c.JSON(http.StatusOK, games)
	}
}

// GetGame fetches a game by its UUID from the database.
func GetGame(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the game UUID from the URL parameter
		gameUUID := c.Param("uuid")

		// Query to select the game by game UUID
		query := "SELECT uuid, service_name, server_ip, player_quantity, game_state, state_time FROM game WHERE uuid = $1"

		row := db.QueryRow(query, gameUUID)

		game := models.Game{}
		err := row.Scan(&game.UUID, &game.ServiceName, &game.ServerIP, &game.PlayerQuantity, &game.GameState, &game.StateTime)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, "Game not found")
			} else {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to fetch the game")
			}
		}

		// Return the game
		return c.JSON(http.StatusOK, game)
	}
}

// CreateGame creates a new game in the database.
func CreateGame(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Parse the request body to extract the game data
		game := new(models.Game)
		if err := c.Bind(game); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Start a database transaction
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to start a transaction")
		}
		defer tx.Rollback() // Rollback the transaction in case of an error

		// Insert the new game into the database
		insertSQL := "INSERT INTO game (uuid, service_name, server_ip, player_quantity, game_state, state_time) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err = tx.Exec(insertSQL, game.UUID, game.ServiceName, game.ServerIP, game.PlayerQuantity, game.GameState, game.StateTime)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to create the game")
		}

		// Update the service table if the game state is "ACTIVE" or "PENDING"
		if game.GameState == "ACTIVE" || game.GameState == "PENDING" {
			updateSQL := "UPDATE service SET server_quantity = server_quantity + 1, total_player = total_player + $1 WHERE service_name = $2"
			_, err = tx.Exec(updateSQL, game.PlayerQuantity, game.ServiceName)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to update the service table")
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to commit the transaction")
		}

		return c.JSON(http.StatusCreated, "Game created successfully")
	}
}

// UpdateGame updates a game in the database.
func UpdateGame(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Get the game UUID from the URL parameter
		gameUUID := c.Param("uuid")

		// Parse the request body to extract the updated game data
		game := new(models.Game)
		if err := c.Bind(game); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request body")
		}

		// Get the current game state before the update
		currentState := ""
		err := db.QueryRow("SELECT game_state FROM game WHERE uuid = $1", gameUUID).Scan(&currentState)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to get current game state")
		}

		// Construct the dynamic SQL query
		updateSQL := "UPDATE game SET "
		args := []interface{}{}
		argCount := 1

		// Check each field for updates and add it to the query
		if game.ServerIP != "" {
			updateSQL += "server_ip = $" + strconv.Itoa(argCount) + ", "
			args = append(args, game.ServerIP)
			argCount++
		}
		if game.PlayerQuantity > 0 {
			updateSQL += "player_quantity = $" + strconv.Itoa(argCount) + ", "
			args = append(args, game.PlayerQuantity)
			argCount++
		}
		if game.GameState != "" {
			updateSQL += "game_state = $" + strconv.Itoa(argCount) + ", "
			args = append(args, game.GameState)
			argCount++
		}
		if game.StateTime != "" {
			updateSQL += "state_time = $" + strconv.Itoa(argCount) + ", "
			args = append(args, game.StateTime)
		}

		// Remove the trailing comma and add the WHERE clause
		updateSQL = strings.TrimSuffix(updateSQL, ", ") + " WHERE uuid = $" + strconv.Itoa(argCount)
		args = append(args, gameUUID)

		// Execute the dynamic SQL query
		_, err = db.Exec(updateSQL, args...)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to update the game")
		}

		// Update the service table if the game state changes
		if currentState != "ACTIVE" && currentState != "PENDING" && (game.GameState == "ACTIVE" || game.GameState == "PENDING") {
			_, err := db.Exec("UPDATE service SET server_quantity = server_quantity + 1, total_player = total_player + $1 WHERE service_name = (SELECT service_name FROM game WHERE uuid = $2)", game.PlayerQuantity, gameUUID)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to update the service table")
			}
		} else if currentState == "ACTIVE" || currentState == "PENDING" {
			// Game state changed from "ACTIVE" or "PENDING" to another state
			_, err := db.Exec("UPDATE service SET server_quantity = server_quantity - 1, total_player = total_player - $1 WHERE service_name = (SELECT service_name FROM game WHERE uuid = $2)", game.PlayerQuantity, gameUUID)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusInternalServerError, "Failed to update the service table")
			}
		}

		return c.JSON(http.StatusOK, "Game updated successfully")
	}
}
