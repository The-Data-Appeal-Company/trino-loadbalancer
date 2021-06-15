package discovery

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Static struct {
	Coordinators []models.Coordinator
}

func NewStatic(coords ...models.Coordinator) Static {
	return Static{Coordinators: coords}
}

func (s Static) Discover(ctx context.Context) ([]models.Coordinator, error) {
	return s.Coordinators, nil
}
