package discovery

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

var ErrClusterNotFound = errors.New("cluster not found")

type Storage interface {
	Remove(context.Context, string) error
	Add(context.Context, models.Coordinator) error
	Update(ctx context.Context, name string, req UpdateRequest) error
	Get(context.Context, string) (models.Coordinator, error)
	All(context.Context) ([]models.Coordinator, error)
}

type UpdateRequest struct {
	Enabled *bool
	Tags    map[string]string
}
