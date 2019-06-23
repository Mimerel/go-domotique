package healthCheck

type HealthDetail struct {
	Name   string `json:"name,omitempty"`
	Health string `json:"health,omitempty"`
}
