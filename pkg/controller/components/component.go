package components

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
)

type QueryHandler interface {
	Execute(context.Context, trino.QueryDetail) error
}

type MultiQueryComponent struct {
	Handlers []QueryHandler
}

func NewMultiQueryComponent(handlers ...QueryHandler) MultiQueryComponent {
	return MultiQueryComponent{Handlers: handlers}
}

func (m MultiQueryComponent) Execute(ctx context.Context, detail trino.QueryDetail) error {
	for _, c := range m.Handlers {
		if err := c.Execute(ctx, detail); err != nil {
			return err
		}
	}
	return nil
}
