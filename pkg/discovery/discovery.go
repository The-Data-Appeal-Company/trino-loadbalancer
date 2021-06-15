package discovery

import (
	"context"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Discovery interface {
	Discover(ctx context.Context) ([]models2.Coordinator, error)
}

type CrossProviderDiscovery struct {
	discoveryProviders []Discovery
}

func NewCrossProviderDiscovery(discoveryProviders []Discovery) *CrossProviderDiscovery {
	return &CrossProviderDiscovery{discoveryProviders: discoveryProviders}
}

func (c *CrossProviderDiscovery) Discover(ctx context.Context) ([]models2.Coordinator, error) {

	coordinators := make([]models2.Coordinator, 0)

	for _, dProvider := range c.discoveryProviders {
		currentProviderCoordinators, err := dProvider.Discover(ctx)

		if err != nil {
			return nil, err
		}
		coordinators = append(coordinators, currentProviderCoordinators...)
	}

	return coordinators, nil

}
