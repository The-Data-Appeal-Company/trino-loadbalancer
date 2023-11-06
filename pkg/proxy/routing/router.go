package routing

import (
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Rule interface {
	Route(Request) (models.Coordinator, error)
}

type CoordinatorWithStatistics struct {
	Coordinator models.Coordinator
	Statistics  trino.ClusterStatistics
}

type Request struct {
	User         string
	Coordinators []CoordinatorWithStatistics
}

type Router struct {
	UserAwareRouter UserAwareRouter
	Rule            Rule
}

func New(userAwareRouter UserAwareRouter, rule Rule) Router {
	return Router{
		UserAwareRouter: userAwareRouter,
		Rule:            rule,
	}
}

func (r Router) Route(req Request) (models.Coordinator, error) {
	if len(req.Coordinators) == 0 {
		return models.Coordinator{}, errors.New("unable to handle routing with no available coordinators")
	}

	req, err := r.UserAwareRouter.Route(req)
	if err != nil {
		return models.Coordinator{}, fmt.Errorf("error routing request: %w", err)
	}
	 
	if len(req.Coordinators) == 0 {
		return models.Coordinator{}, ErrRoutingNotFound
	}

	return r.Rule.Route(req)
}
