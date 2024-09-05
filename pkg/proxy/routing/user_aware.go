package routing

import (
	"errors"
	"regexp"
)

var (
	ErrForbiddenRouting = errors.New("no route match for user")
	ErrRouteNotFound    = errors.New("no routable cluster")
)

type NoMatchBehaviour string

const (
	NoMatchBehaviourForbid  = "FORBID"
	NoMatchBehaviourDefault = "DEFAULT"
)

type UserAwareRoutingConf struct {
	Default UserAwareDefault
	Rules   []UserAwareRoutingRule
}

type UserAwareRoutingRule struct {
	User    *regexp.Regexp
	Cluster UserAwareClusterMatchRule
}

type UserAwareDefault struct {
	Behaviour NoMatchBehaviour
	Cluster   UserAwareClusterMatchRule
}

type UserAwareClusterMatchRule struct {
	Name                  *regexp.Regexp
	Tags                  map[string]string
	UseDefaultIfUnhealthy bool
}

func NewUserAwareRouter(conf UserAwareRoutingConf) UserAwareRouter {
	return UserAwareRouter{conf}
}

type UserAwareRouter struct {
	conf UserAwareRoutingConf
}

func (u UserAwareRouter) Route(req Request) (Request, error) {
	if len(u.conf.Rules) == 0 {
		return req, nil
	}

	// test matching rule for the requester user
	rule, matched := u.matchRule(req)

	// if no rule is found we apply the configured default behaviour
	if !matched {
		if u.conf.Default.Behaviour == NoMatchBehaviourForbid {
			return Request{}, ErrForbiddenRouting
		}
		rule = u.conf.Default.Cluster
	}

	// filter request's coordinator using the rule configuration
	coordinators := filterByRule(rule, req.Coordinators)

	if (coordinators == nil || len(coordinators) == 0) && rule.UseDefaultIfUnhealthy {
		coordinators = filterByRule(u.conf.Default.Cluster, req.Coordinators)
	}

	req.Coordinators = coordinators
	return req, nil
}

func filterByRule(rule UserAwareClusterMatchRule, coordinators []CoordinatorWithStatistics) []CoordinatorWithStatistics {
	coords := make([]CoordinatorWithStatistics, 0)
	for _, coord := range coordinators {

		if rule.Name != nil && !rule.Name.MatchString(coord.Coordinator.Name) {
			continue
		}

		if rule.Tags != nil && !matchTags(coord.Coordinator.Tags, rule.Tags) {
			continue
		}

		coords = append(coords, coord)
	}

	return coords
}

func (u UserAwareRouter) matchRule(req Request) (UserAwareClusterMatchRule, bool) {
	for _, r := range u.conf.Rules {
		if r.User.MatchString(req.User) {
			return r.Cluster, true
		}
	}
	return UserAwareClusterMatchRule{}, false
}

func matchTags(source map[string]string, match map[string]string) bool {
	for k, v := range match {
		if source[k] != v {
			return false
		}
	}
	return true
}
