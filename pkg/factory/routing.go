package factory

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/routing"
)

type QueryRouterConfiguration struct {
	Type string
}

func CreateQueryRouter(conf QueryRouterConfiguration) (routing.Router, error) {
	rule, err := createRouterRule(conf.Type)
	if err != nil {
		return routing.Router{}, err
	}
	return routing.New(rule), nil
}

func createRouterRule(t string) (routing.Rule, error) {
	switch t {
	case "random":
		return routing.Random(), nil
	case "round-robin":
		return routing.RoundRobin(), nil
	case "less-running-queries":
		return routing.LessRunningQueries(), nil
	default:
		return nil, fmt.Errorf("no router rule for value: %s", t)
	}
}
