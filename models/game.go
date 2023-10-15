package models

// Game represents a game in the application.
type Game struct {
	GameUUID  string `json:"game_uuid"`
	ServiceID int    `json:"service_id"`
}
