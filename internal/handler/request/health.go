package request

type HealthCheckRequest struct {
	IncludeDB bool `json:"include_db"`
}
