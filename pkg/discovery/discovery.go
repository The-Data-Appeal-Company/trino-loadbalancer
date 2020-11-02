package discovery

import (
	"context"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

type Discovery interface {
	Discover(ctx context.Context) ([]models.Coordinator, error)
}

type CrossProviderDiscovery struct {
	discoveryProviders []Discovery
}

func NewCrossProviderDiscovery(discoveryProviders []Discovery) *CrossProviderDiscovery {
	return &CrossProviderDiscovery{discoveryProviders: discoveryProviders}
}

func (c *CrossProviderDiscovery) Discover(ctx context.Context) ([]models.Coordinator, error) {

	coordinators := make([]models.Coordinator, 0)

	for _, dProvider := range c.discoveryProviders {
		currentProviderCoordinators, err := dProvider.Discover(ctx)

		if err != nil {
			return nil, err
		}
		coordinators = append(coordinators, currentProviderCoordinators...)
	}

	return coordinators, nil

}
