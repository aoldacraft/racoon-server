package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go_racoon/app/handlers"
	"go_racoon/app/routes"
	"os"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database connection
	// 환경 변수에서 데이터베이스 연결 정보 가져오기
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// PostgreSQL 연결 문자열 생성
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// 데이터베이스에 연결
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
	}

	// 데이터베이스에 테이블 생성
	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS service (
			service_name VARCHAR PRIMARY KEY,
			server_quantity INTEGER,
			total_player INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS game (
			uuid VARCHAR PRIMARY KEY,
			service_name VARCHAR REFERENCES service(service_name),
			server_ip VARCHAR,
			player_quantity INTEGER,
			game_state VARCHAR,
			state_time TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS game_log (
			log_id SERIAL PRIMARY KEY,
			service_name VARCHAR REFERENCES service(service_name),
			uuid VARCHAR REFERENCES game(uuid),
			log_text VARCHAR,
			log_time TIMESTAMP
		);`,
	}

	for _, query := range createTableQueries {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Table creation complete!")

	// Register routes
	routes.RegisterServiceRoutes(e, db)
	routes.RegisterGameRoutes(e, db)
	routes.RegisterGameLogRoutes(e, db)

	// Health check endpoint
	e.GET("/api/health", handlers.HealthCheck)

	// Run the server
	e.Logger.Fatal(e.Start(":1323"))
}
