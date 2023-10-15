package handlers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"go_racoon/models"
	"log"
	"net/http"
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

		// Insert the new game into the database
		insertSQL := "INSERT INTO game (uuid, service_name, server_ip, player_quantity, game_state, state_time) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err := db.Exec(insertSQL, game.UUID, game.ServiceName, game.ServerIP, game.PlayerQuantity, game.GameState, game.StateTime)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to create the game")
		}

		return c.JSON(http.StatusCreated, "Game created successfully")
	}
}
