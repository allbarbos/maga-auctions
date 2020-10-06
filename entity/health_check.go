package entity

// HealthCheck entity
type HealthCheck struct {
	Status       string            `json:"status"`
	Dependencies map[string]string `json:"dependencies"`
}
