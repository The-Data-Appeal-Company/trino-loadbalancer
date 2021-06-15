package discovery

import (
	"context"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Static struct {
	Coordinators []models2.Coordinator
}

func NewStatic(coords ...models2.Coordinator) Static {
	return Static{Coordinators: coords}
}

func (s Static) Discover(ctx context.Context) ([]models2.Coordinator, error) {
	return s.Coordinators, nil
}
