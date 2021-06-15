package discovery

import (
	"context"
	"errors"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

var ErrClusterNotFound = errors.New("cluster not found")

type Storage interface {
	Remove(context.Context, string) error
	Add(context.Context, models2.Coordinator) error
	Get(context.Context, string) (models2.Coordinator, error)
	All(context.Context) ([]models2.Coordinator, error)
}
