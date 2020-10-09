package routing

import (
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunningQueriesRouter(t *testing.T) {

	values := []string{"test-0", "test-1", "test-2"}

	random := RunningQueriesRouter{}
	route, err := random.Route(Request{
		User: "test",
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models.Coordinator{
					Name: values[0],
				},
				Statistics: models.ClusterStatistics{
					RunningQueries: 100,
				},
			},
			{
				Coordinator: models.Coordinator{
					Name: values[1],
				},
				Statistics: models.ClusterStatistics{
					RunningQueries: 50,
				},
			},
			{
				Coordinator: models.Coordinator{
					Name: values[2],
				},
				Statistics: models.ClusterStatistics{
					RunningQueries: 20,
				},
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, values[2], route.Name)
}
