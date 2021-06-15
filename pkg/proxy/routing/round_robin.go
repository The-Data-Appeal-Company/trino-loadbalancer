package routing

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"sync"
)

func RoundRobin() *RoundRobinRule {
	return &RoundRobinRule{
		index:        0,
		coordinators: 0,
		mutex:        &sync.Mutex{},
	}
}

type RoundRobinRule struct {
	index        int
	coordinators int
	mutex        *sync.Mutex
}

func (r *RoundRobinRule) Route(request Request) (models.Coordinator, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// if the routing state has changed reset the index
	if r.coordinators != len(request.Coordinators) {
		r.index = 0
		r.coordinators = len(request.Coordinators)
	}

	selected := request.Coordinators[r.index].Coordinator

	r.index = (r.index + 1) % r.coordinators
	return selected, nil
}
