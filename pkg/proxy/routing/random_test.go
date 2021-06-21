package routing

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
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
				Statistics: trino.ClusterStatistics{},
			},
			{
				Coordinator: models.Coordinator{
					Name: values[1],
				},
				Statistics: trino.ClusterStatistics{},
			},
			{
				Coordinator: models.Coordinator{
					Name: values[2],
				},
				Statistics: trino.ClusterStatistics{},
			},
		},
	})

	require.NoError(t, err)

	require.Contains(t, values, route.Name)

}
