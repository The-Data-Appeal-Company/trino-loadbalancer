package discovery

import (
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) Discover() ([]models.Coordinator, error) {
	return []models.Coordinator{}, nil
}
