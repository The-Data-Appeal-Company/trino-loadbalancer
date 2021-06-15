package discovery

import (
	"context"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) Discover(ctx context.Context) ([]models2.Coordinator, error) {
	return []models2.Coordinator{}, nil
}
