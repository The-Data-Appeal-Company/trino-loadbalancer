package routing

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRouterHandleRequest(t *testing.T) {
	router := New(NewUserAwareRouter(UserAwareRoutingConf{}), RoundRobin())
	route, err := router.Route(Request{
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models.Coordinator{
					Name: "test",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, route.Name, "test")
}

func TestRouterHandleEmptyCoordinators(t *testing.T) {
	router := New(NewUserAwareRouter(UserAwareRoutingConf{}), RoundRobin())
	_, err := router.Route(Request{})
	require.Error(t, err)
}
