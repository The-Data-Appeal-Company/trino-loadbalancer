package routing

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoundRobinRouting(t *testing.T) {
	values := []string{"test-0", "test-1", "test-2"}

	random := RoundRobin()
	request := Request{
		User: "test",
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models2.Coordinator{
					Name: values[0],
				},
				Statistics: trino.ClusterStatistics{},
			},
			{
				Coordinator: models2.Coordinator{
					Name: values[1],
				},
				Statistics: trino.ClusterStatistics{},
			},
			{
				Coordinator: models2.Coordinator{
					Name: values[2],
				},
				Statistics: trino.ClusterStatistics{},
			},
		},
	}

	for _, coord := range values {
		route, err := random.Route(request)
		require.NoError(t, err)
		require.Equal(t, route.Name, coord)
	}

}
