package routing

import (
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Rule interface {
	Route(Request) (models2.Coordinator, error)
}

type CoordinatorWithStatistics struct {
	Coordinator models2.Coordinator
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

func (r Router) Route(req Request) (models2.Coordinator, error) {
	if len(req.Coordinators) == 0 {
		return models2.Coordinator{}, errors.New("unable to handle routing with no available coordinators")
	}

	req, err := r.UserAwareRouter.Route(req)
	if err != nil {
		return models2.Coordinator{}, fmt.Errorf("error routing request: %w", err)
	}

	return r.Rule.Route(req)
}
