package models

// Game represents a game in the application.
type Game struct {
	UUID           string `json:"uuid"`
	ServiceName    string `json:"service_name"`
	ServerIP       string `json:"server_ip"`
	PlayerQuantity int    `json:"player_quantity"`
	GameState      string `json:"game_state"`
	StateTime      string `json:"state_time"`
}
