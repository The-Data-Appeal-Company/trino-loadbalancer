package session

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
)

var (
	ErrLinkNotFound = errors.New("no link found for query")
)

type Writer interface {
	Link(context.Context, trino.QueryInfo, string) error
	Unlink(context.Context, trino.QueryInfo, string) error
}

type Reader interface {
	Get(context.Context, trino.QueryInfo) (string, error)
}

type Storage interface {
	Link(context.Context, trino.QueryInfo, string) error
	Unlink(context.Context, trino.QueryInfo) error
	Get(context.Context, trino.QueryInfo) (string, error)
}
