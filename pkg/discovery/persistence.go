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
	Get(context.Context, string) (models.Coordinator, error)
	All(context.Context) ([]models.Coordinator, error)
}
