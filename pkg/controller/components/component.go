package components

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
)

type QueryComponent interface {
	Execute(context.Context, trino.QueryDetail) error
}
