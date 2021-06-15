package routing

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"math/rand"
)

func Random() RandomRouter {
	return RandomRouter{}
}

type RandomRouter struct {
}

func (r RandomRouter) Route(request Request) (models2.Coordinator, error) {
	max := len(request.Coordinators)
	n := rand.Int31n(int32(max))
	return request.Coordinators[n].Coordinator, nil
}
