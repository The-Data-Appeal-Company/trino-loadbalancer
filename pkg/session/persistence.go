package session

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

var (
	ErrLinkNotFound = errors.New("no link found for query")
)

type Writer interface {
	Link(context.Context, models.QueryInfo, string) error
	Unlink(context.Context, models.QueryInfo, string) error
}

type Reader interface {
	Get(context.Context, models.QueryInfo) (string, error)
}

type Storage interface {
	Link(context.Context, models.QueryInfo, string) error
	Unlink(context.Context, models.QueryInfo) error
	Get(context.Context, models.QueryInfo) (string, error)
}
