package models

// Service represents a service in the application.
type Service struct {
	ServiceName    string `json:"service_name"`
	ServerQuantity int    `json:"server_quantity"`
	TotalPlayer    int    `json:"total_player"`
}
