package discovery

import (
	"context"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) Discover(ctx context.Context) ([]models.Coordinator, error) {
	return []models.Coordinator{}, nil
}
