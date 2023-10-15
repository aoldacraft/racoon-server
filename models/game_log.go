package models

// GameLog represents a game log in the application.
type GameLog struct {
	LogID    int    `json:"log_id"`
	GameUUID string `json:"game_uuid"`
	LogText  string `json:"log_text"`
}
