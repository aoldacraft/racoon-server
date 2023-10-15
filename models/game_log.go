package models

// GameLog represents a game log in the application.
type GameLog struct {
	LogID       int    `json:"log_id"`
	ServiceName string `json:"service_name"`
	UUID        string `json:"uuid"`
	LogText     string `json:"log_text"`
	LogTime     string `json:"log_time"`
}
