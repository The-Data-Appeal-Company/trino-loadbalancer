package configuration

import (
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/routing"
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
	Rule  string           `json:"rule" yaml:"rule" mapstructure:"rule"`
	Users RoutingUsersConf `json:"users" yaml:"users" mapstructure:"users"`
}

func CreateQueryRouter(conf RoutingConf) (routing.Router, error) {
	userAwareRouter, err := createUserAwareRouter(conf.Users)
	if err != nil {
		return routing.Router{}, err
	}

	rule, err := createRouterRule(conf.Rule)
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

	defaultNameRe, err := regexpOrNil(users.Default.Cluster.Name)
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

	rules := make([]routing.UserAwareRoutingRule, len(users.Rules))
	for i, r := range users.Rules {
		userRe, err := regexpOrNil(r.User)
		if err != nil {
			return routing.UserAwareRouter{}, nil
		}

		if userRe == nil {
			return routing.UserAwareRouter{}, errors.New("user must be specified on routing rule")
		}

		clusterNameRe, err := regexpOrNil(r.Cluster.Name)
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

func regexpOrNil(raw string) (*regexp.Regexp, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return regexp.Compile(raw)
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
