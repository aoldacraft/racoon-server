package models

// Service represents a service in the application.
type Service struct {
	ServiceID   int    `json:"service_id"`
	ServiceName string `json:"service_name"`
}
