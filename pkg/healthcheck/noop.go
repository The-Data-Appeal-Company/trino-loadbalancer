package healthcheck

import (
	"net/url"
	"time"
)

type NoOpCheck struct {
}

func NoOp() NoOpCheck {
	return NoOpCheck{}
}

func (n NoOpCheck) Check(*url.URL) (Health, error) {
	return Health{
		Status:    StatusHealthy,
		Message:   "noop hc",
		Timestamp: time.Now(),
	}, nil
}
