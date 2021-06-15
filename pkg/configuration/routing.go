package configuration

import (
	"errors"
	"fmt"
	routing2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/routing"
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

func CreateQueryRouter(conf RoutingConf) (routing2.Router, error) {
	userAwareRouter, err := createUserAwareRouter(conf.Users)
	if err != nil {
		return routing2.Router{}, err
	}

	rule, err := createRouterRule(conf.Rule)
	if err != nil {
		return routing2.Router{}, err
	}

	return routing2.New(userAwareRouter, rule), nil
}

func createUserAwareRouter(users RoutingUsersConf) (routing2.UserAwareRouter, error) {
	defaultBehaviour, err := createBehaviour(users.Default.Behaviour)
	if err != nil {
		return routing2.UserAwareRouter{}, nil
	}

	defaultNameRe, err := regexpOrNil(users.Default.Cluster.Name)
	if err != nil {
		return routing2.UserAwareRouter{}, nil
	}

	var conf routing2.UserAwareRoutingConf
	conf.Default = routing2.UserAwareDefault{
		Behaviour: defaultBehaviour,
		Cluster: routing2.UserAwareClusterMatchRule{
			Name: defaultNameRe,
			Tags: users.Default.Cluster.Tags,
		},
	}

	rules := make([]routing2.UserAwareRoutingRule, len(users.Rules))
	for i, r := range users.Rules {
		userRe, err := regexpOrNil(r.User)
		if err != nil {
			return routing2.UserAwareRouter{}, nil
		}

		if userRe == nil {
			return routing2.UserAwareRouter{}, errors.New("user must be specified on routing rule")
		}

		clusterNameRe, err := regexpOrNil(r.Cluster.Name)
		if err != nil {
			return routing2.UserAwareRouter{}, nil
		}

		rules[i] = routing2.UserAwareRoutingRule{
			User: userRe,
			Cluster: routing2.UserAwareClusterMatchRule{
				Name: clusterNameRe,
				Tags: r.Cluster.Tags,
			},
		}
	}

	conf.Rules = rules
	return routing2.NewUserAwareRouter(conf), nil
}

func regexpOrNil(raw string) (*regexp.Regexp, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return regexp.Compile(raw)
}

func createBehaviour(raw string) (routing2.NoMatchBehaviour, error) {
	switch strings.ToLower(raw) {
	case "forbid":
		return routing2.NoMatchBehaviourForbid, nil
	case "default":
		return routing2.NoMatchBehaviourDefault, nil
	default:
		return routing2.NoMatchBehaviourForbid, fmt.Errorf("invalid behaviour type: %s", raw)
	}
}

func createRouterRule(t string) (routing2.Rule, error) {
	switch t {
	case "random":
		return routing2.Random(), nil
	case "round-robin":
		return routing2.RoundRobin(), nil
	case "less-running-queries":
		return routing2.LessRunningQueries(), nil
	default:
		return nil, fmt.Errorf("no router rule for value: %s", t)
	}
}
