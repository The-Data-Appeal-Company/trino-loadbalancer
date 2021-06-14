package routing

import (
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestUserAwareRouter(t *testing.T) {
	uar := NewUserAwareRouter(UserAwareRoutingConf{
		Default: UserAwareDefault{
			Behaviour: NoMatchBehaviourDefault,
			Cluster: UserAwareClusterMatchRule{
				Name: mustRegex(t, "cluster-00"),
				Tags: map[string]string{
					"test": "true",
				},
			},
		},
		Rules: []UserAwareRoutingRule{
			{
				User: mustRegex(t, "test-app-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				User: mustRegex(t, "test-user-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	})

	request := Request{
		User: "test-user-00",
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models.Coordinator{
					Name: "cluster-00",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				Coordinator: models.Coordinator{
					Name: "cluster-01",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	}
	coords, err := uar.Route(request)

	require.NoError(t, err)

	require.Equal(t, coords.User, request.User)
	require.Equal(t, coords.Coordinators[0], request.Coordinators[0])
	require.Len(t, coords.Coordinators, 1)
}

func TestUserAwareRoutingFallbackOnDefault(t *testing.T) {
	uar := NewUserAwareRouter(UserAwareRoutingConf{
		Default: UserAwareDefault{
			Behaviour: NoMatchBehaviourDefault,
			Cluster: UserAwareClusterMatchRule{
				Name: mustRegex(t, "cluster-01"),
				Tags: map[string]string{
					"test": "true",
				},
			},
		},
		Rules: []UserAwareRoutingRule{
			{
				User: mustRegex(t, "test-app-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				User: mustRegex(t, "test-user-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	})

	request := Request{
		User: "non-matching-user",
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models.Coordinator{
					Name: "cluster-00",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				Coordinator: models.Coordinator{
					Name: "cluster-01",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	}
	coords, err := uar.Route(request)

	require.NoError(t, err)

	require.Equal(t, coords.User, request.User)
	require.Equal(t, coords.Coordinators[0], request.Coordinators[1])
	require.Len(t, coords.Coordinators, 1)
}

func TestRoutingAwareFilterByTags(t *testing.T) {
	uar := NewUserAwareRouter(UserAwareRoutingConf{
		Default: UserAwareDefault{
			Behaviour: NoMatchBehaviourForbid,
		},
		Rules: []UserAwareRoutingRule{
			{
				User: mustRegex(t, "test-user-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	})

	request := Request{
		User: "test-user-00",
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models.Coordinator{
					Name: "cluster-00",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				Coordinator: models.Coordinator{
					Name: "cluster-01",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	}
	coords, err := uar.Route(request)

	require.NoError(t, err)

	require.Equal(t, coords.User, request.User)
	require.Len(t, coords.Coordinators, 2)
}

func TestRoutingAwareFilterByUserForbidOnError(t *testing.T) {
	uar := NewUserAwareRouter(UserAwareRoutingConf{
		Default: UserAwareDefault{
			Behaviour: NoMatchBehaviourForbid,
			Cluster: UserAwareClusterMatchRule{
				Name: mustRegex(t, "cluster-01"),
				Tags: map[string]string{
					"test": "true",
				},
			},
		},
		Rules: []UserAwareRoutingRule{
			{
				User: mustRegex(t, "test-app-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				User: mustRegex(t, "test-user-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	})

	request := Request{
		User: "non-matching-user",
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models.Coordinator{
					Name: "cluster-00",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				Coordinator: models.Coordinator{
					Name: "cluster-01",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	}
	_, err := uar.Route(request)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrForbiddenRouting))
}

func TestUserAwareRoutingSelectAllCoordinatorsWhenFallbackOnEmptyDefault(t *testing.T) {
	uar := NewUserAwareRouter(UserAwareRoutingConf{
		Default: UserAwareDefault{
			Behaviour: NoMatchBehaviourDefault,
			Cluster: UserAwareClusterMatchRule{},
		},
		Rules: []UserAwareRoutingRule{
			{
				User: mustRegex(t, "test-app-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				User: mustRegex(t, "test-user-(.+)"),
				Cluster: UserAwareClusterMatchRule{
					Name: mustRegex(t, "cluster-00"),
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	})

	request := Request{
		User: "non-matching-user",
		Coordinators: []CoordinatorWithStatistics{
			{
				Coordinator: models.Coordinator{
					Name: "cluster-00",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
			{
				Coordinator: models.Coordinator{
					Name: "cluster-01",
					Tags: map[string]string{
						"test": "true",
					},
				},
			},
		},
	}
	coords, err := uar.Route(request)

	require.NoError(t, err)

	require.Equal(t, coords.User, request.User)
	require.Len(t, coords.Coordinators, 2)
}

func mustRegex(t *testing.T, raw string) *regexp.Regexp {
	re, err := regexp.Compile(raw)
	require.NoError(t, err)
	return re
}
