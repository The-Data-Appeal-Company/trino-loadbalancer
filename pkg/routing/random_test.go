package routing

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomRouter(t *testing.T) {

	values := []string{"test-0", "test-1", "test-2"}

	random := RandomRouter{}
	route, err := random.Route(Request{
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
	})

	require.NoError(t, err)

	require.Contains(t, values, route.Name)

}
