package routing

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type RunningQueriesRouter struct {
}

func LessRunningQueries() RunningQueriesRouter {
	return RunningQueriesRouter{}
}

func (r RunningQueriesRouter) Route(request Request) (models2.Coordinator, error) {
	var selected = request.Coordinators[0]

	for i := 1; i < len(request.Coordinators); i++ {
		current := request.Coordinators
		if current[i].Statistics.RunningQueries < selected.Statistics.RunningQueries {
			selected = request.Coordinators[i]
		}
	}

	return selected.Coordinator, nil
}
