package routing

import (
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
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
				Coordinator: models.Coordinator{
					Name: values[0],
				},
				Statistics: models.ClusterStatistics{},
			},
			{
				Coordinator: models.Coordinator{
					Name: values[1],
				},
				Statistics: models.ClusterStatistics{},
			},
			{
				Coordinator: models.Coordinator{
					Name: values[2],
				},
				Statistics: models.ClusterStatistics{},
			},
		},
	}

	for _, coord := range values {
		route, err := random.Route(request)
		require.NoError(t, err)
		require.Equal(t, route.Name, coord)
	}

}
