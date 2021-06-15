package session

import (
	"context"
	"errors"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

var (
	ErrLinkNotFound = errors.New("no link found for query")
)

type Writer interface {
	Link(context.Context, models2.QueryInfo, string) error
	Unlink(context.Context, models2.QueryInfo, string) error
}

type Reader interface {
	Get(context.Context, models2.QueryInfo) (string, error)
}

type Storage interface {
	Link(context.Context, models2.QueryInfo, string) error
	Unlink(context.Context, models2.QueryInfo) error
	Get(context.Context, models2.QueryInfo) (string, error)
}
