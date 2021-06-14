package factory

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/routing"
	"regexp"
	"strings"
)

type RoutingUsersConf struct {
	Default struct {
		Behaviour string `json:"behaviour" yaml:"behaviour" mapstructure:"behaviour"`
		Cluster   struct {
			Name string            `json:"name" yaml:"name" mapstructure:"name"`
			Tags map[string]string `json:"tags" yaml:"tags" mapstructure:"tags"`
		} `json:"cluster" yaml:"cluster" mapstructure:"cluster"`
	} `json:"default" yaml:"default" mapstructure:"default"`
	Rules []struct {
		User    string `json:"user" yaml:"user" mapstructure:"user"`
		Cluster struct {
			Name string            `json:"name" yaml:"name" mapstructure:"name"`
			Tags map[string]string `json:"tags" yaml:"tags" mapstructure:"tags"`
		} `json:"cluster" yaml:"cluster" mapstructure:"cluster"`
	} `json:"rules" yaml:"rules" mapstructure:"rules"`
}

type RoutingConf struct {
	Routing struct {
		Mode  string           `json:"mode" yaml:"mode" mapstructure:"mode"`
		Users RoutingUsersConf `json:"users" yaml:"users" mapstructure:"users"`
	} `json:"routing" yaml:"routing" mapstructure:"routing"`
}

func CreateQueryRouter(conf RoutingConf) (routing.Router, error) {
	userAwareRouter, err := createUserAwareRouter(conf.Routing.Users)
	if err != nil {
		return routing.Router{}, err
	}

	rule, err := createRouterRule(conf.Routing.Mode)
	if err != nil {
		return routing.Router{}, err
	}

	return routing.New(userAwareRouter, rule), nil
}

func createUserAwareRouter(users RoutingUsersConf) (routing.UserAwareRouter, error) {
	defaultBehaviour, err := createBehaviour(users.Default.Behaviour)
	if err != nil {
		return routing.UserAwareRouter{}, nil
	}

	defaultNameRe, err := regexp.Compile(users.Default.Cluster.Name)
	if err != nil {
		return routing.UserAwareRouter{}, nil
	}

	var conf routing.UserAwareRoutingConf
	conf.Default = routing.UserAwareDefault{
		Behaviour: defaultBehaviour,
		Cluster: routing.UserAwareClusterMatchRule{
			Name: defaultNameRe,
			Tags: users.Default.Cluster.Tags,
		},
	}

	rules := make([]routing.UserAwareRoutingRule, 0)
	for i, r := range users.Rules {
		userRe, err := regexp.Compile(r.User)
		if err != nil {
			return routing.UserAwareRouter{}, nil
		}

		clusterNameRe, err := regexp.Compile(r.Cluster.Name)
		if err != nil {
			return routing.UserAwareRouter{}, nil
		}

		rules[i] = routing.UserAwareRoutingRule{
			User: userRe,
			Cluster: routing.UserAwareClusterMatchRule{
				Name: clusterNameRe,
				Tags: r.Cluster.Tags,
			},
		}
	}

	conf.Rules = rules
	return routing.NewUserAwareRouter(conf), nil
}

func createBehaviour(raw string) (routing.NoMatchBehaviour, error) {
	switch strings.ToLower(raw) {
	case "forbid":
		return routing.NoMatchBehaviourForbid, nil
	case "default":
		return routing.NoMatchBehaviourDefault, nil
	default:
		return routing.NoMatchBehaviourForbid, fmt.Errorf("invalid behaviour type: %s", raw)
	}
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
