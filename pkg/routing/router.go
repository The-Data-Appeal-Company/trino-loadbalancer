package routing

import (
	"errors"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

type Rule interface {
	Route(Request) (models.Coordinator, error)
}

type CoordinatorWithStatistics struct {
	Coordinator models.Coordinator
	Statistics  models.ClusterStatistics
}

type Request struct {
	User         string
	Coordinators []CoordinatorWithStatistics
}

type Router struct {
	Rule Rule
}

func New(rule Rule) Router {
	return Router{
		Rule: rule,
	}
}

func (r Router) Route(req Request) (models.Coordinator, error) {
	if len(req.Coordinators) == 0 {
		return models.Coordinator{}, errors.New("unable to handle routing with no available coordinators")
	}

	return r.Rule.Route(req)
}
