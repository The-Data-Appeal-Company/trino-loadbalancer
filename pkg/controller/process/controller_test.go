package process

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/components"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestController_Run(t *testing.T) {
	c := Controller{
		api:          trino.MockApi{},
		discovery:    discovery.NewMemoryStorage(),
		healthCheck:  healthcheck.NoOp(),
		state:        NewInMemoryState(),
		queryHandler: components.NewMultiQueryHandler(),
		logger:       logging.Noop(),
	}

	err := c.Run(context.TODO())
	require.NoError(t, err)
}

func TestController_RunUpdateState(t *testing.T) {
	ctx := context.TODO()

	coordinator0 := models.Coordinator{
		Name: "cluster-00",
	}

	state := NewInMemoryState()
	storage := discovery.NewMemoryStorage()
	require.NoError(t, storage.Add(ctx, coordinator0))

	c := Controller{
		api:          trino.MockApi{},
		discovery:    storage,
		healthCheck:  healthcheck.NoOp(),
		state:        state,
		queryHandler: components.NewMultiQueryHandler(),
		logger:       logging.Noop(),
	}

	err := c.Run(ctx)

	require.NoError(t, err)
	run0, err := state.Get(ctx, coordinator0)
	require.NoError(t, err)
	require.NotEqual(t, run0, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))

	err = c.Run(ctx)

	require.NoError(t, err)
	run1, err := state.Get(ctx, coordinator0)
	require.NoError(t, err)
	require.NotEqual(t, run0, run1)
}

func TestController_RunWithQueryHandler(t *testing.T) {
	ctx := context.TODO()

	coordinator0 := models.Coordinator{
		Name: "cluster-00",
	}

	state := NewInMemoryState()
	storage := discovery.NewMemoryStorage()
	require.NoError(t, storage.Add(ctx, coordinator0))

	queriesTs := time.Now()

	api := trino.MockApi{
		QueryDetailFn: func(coordinator models.Coordinator, queryID string) (trino.QueryDetail, error) {
			switch queryID {
			case "query-00":
				return trino.QueryDetail{QueryID: "query-00"}, nil
			case "query-01":
				return trino.QueryDetail{QueryID: "query-01"}, nil
			default:
				return trino.QueryDetail{}, errors.New("query not found")
			}
		},
		QueryListFn: func(coordinator models.Coordinator) (trino.QueryList, error) {
			return trino.QueryList{
				{
					QueryId:    "query-00",
					State:      trino.QueryFinished,
					QueryStats: trino.QueryStats{CreateTime: queriesTs},
				},
				{
					QueryId:    "query-01",
					State:      trino.QueryFinished,
					QueryStats: trino.QueryStats{CreateTime: queriesTs},
				},
				{
					QueryId:    "query-02",
					State:      "RUNNING",
					QueryStats: trino.QueryStats{CreateTime: queriesTs},
				},
			}, nil
		},
	}

	queryHandler := components.NewMockQueryHandler(nil)

	c := Controller{
		api:          api,
		discovery:    storage,
		healthCheck:  healthcheck.NoOp(),
		state:        state,
		queryHandler: queryHandler,
		logger: logging.Noop(),
	}

	err := c.Run(ctx)
	require.NoError(t, err)

	queryHandlerCalls := queryHandler.Calls()

	require.Len(t, queryHandlerCalls, 2)

	queryIds := make([]string, len(queryHandlerCalls))
	for i := range queryHandlerCalls {
		queryIds[i] = queryHandlerCalls[i].QueryID
	}

	require.Contains(t, queryIds, "query-00")
	require.Contains(t, queryIds, "query-01")

	err = c.Run(ctx)
	require.NoError(t, err)
	require.Len(t, queryHandler.Calls(), 2)
}
