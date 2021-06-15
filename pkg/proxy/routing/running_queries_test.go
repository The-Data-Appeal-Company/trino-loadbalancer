package routing

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
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
				Coordinator: models2.Coordinator{
					Name: values[0],
				},
				Statistics: models2.ClusterStatistics{
					RunningQueries: 100,
				},
			},
			{
				Coordinator: models2.Coordinator{
					Name: values[1],
				},
				Statistics: models2.ClusterStatistics{
					RunningQueries: 50,
				},
			},
			{
				Coordinator: models2.Coordinator{
					Name: values[2],
				},
				Statistics: models2.ClusterStatistics{
					RunningQueries: 20,
				},
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, values[2], route.Name)
}
