package lb

import (
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

type MockPool struct {
	coordinators []*CoordinatorConnection
}

func NewMockPool() *MockPool {
	return &MockPool{
		coordinators: make([]*CoordinatorConnection, 0),
	}
}

func (m *MockPool) AllBackends() []*CoordinatorConnection {
	return m.coordinators
}

func (m *MockPool) AvailableBackends() ([]*CoordinatorConnection, error) {
	return m.coordinators, nil
}

func (m *MockPool) Add(coordinator models.Coordinator) error {
	m.coordinators = append(m.coordinators, &CoordinatorConnection{
		Proxy:   nil,
		Backend: coordinator,
	})
	return nil
}

func (m *MockPool) GetByName(name string, unhealthy healthcheck.HealthStatus) (*CoordinatorConnection, error) {
	for _, c := range m.coordinators {
		if c.Backend.Name == name {
			return c, nil
		}
	}

	return nil, ErrNoBackendsAvailable
}

func (m *MockPool) Remove(name string) error {
	for i, c := range m.coordinators {
		if c.Backend.Name == name {
			m.coordinators = remove(m.coordinators, i)
		}
	}
	return nil
}

func (m *MockPool) Update(name string, state models.Coordinator) error {
	for _, c := range m.coordinators {
		if c.Backend.Name == name {
			c.Backend.Tags = state.Tags
			c.Backend.Enabled = state.Enabled
		}
	}
	return nil
}

func (m *MockPool) UpdateStatus() error {
	return nil
}
