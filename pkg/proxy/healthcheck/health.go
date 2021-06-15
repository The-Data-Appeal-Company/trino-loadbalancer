package healthcheck

import (
	"net/url"
	"time"
)

type HealthStatus int

const (
	StatusUnknown   HealthStatus = iota
	StatusUnhealthy HealthStatus = iota
	StatusHealthy   HealthStatus = iota
)

func (d HealthStatus) String() string {
	return [...]string{"unknown", "unhealthy", "healthy"}[d]
}

type HealthCheck interface {
	Check(*url.URL) (Health, error)
}

type Health struct {
	Status  HealthStatus
	Message string
	Timestamp time.Time
}

func (s Health) IsAvailable() bool {
	return s.Status == StatusHealthy
}
