package lb

import (
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/google/uuid"
	"net/http"
)

type MockPool struct {
	coordinators []*CoordinatorRef
}

func NewMockPool() *MockPool {
	return &MockPool{
		coordinators: make([]*CoordinatorRef, 0),
	}
}

func (m *MockPool) Handle(coordinator CoordinatorRef, writer http.ResponseWriter, request *http.Request) error {
	return errors.New("mock doesn't implement Handle yet")
}

func (m *MockPool) Fetch(req FetchRequest) []CoordinatorRef {
	selected := make([]CoordinatorRef, 0)
	for _, cc := range m.coordinators {
		if len(req.Name) != 0 && cc.Name != req.Name {
			continue
		}

		//if req.Health != 0 && cc.health.Status < req.Health {
		//	continue
		//}

		if len(req.Tags) != 0 && !matchTags(cc.Tags, req.Tags) {
			continue
		}

		selected = append(selected, *cc)
	}

	return selected

}

func (m *MockPool) Add(coordinator models2.Coordinator) error {
	m.coordinators = append(m.coordinators, &CoordinatorRef{
		ID:          CoordinatorConnectionID(uuid.New().String()),
		Statistics:  trino.ClusterStatistics{},
		Coordinator: coordinator,
	})
	return nil
}

func (m *MockPool) Remove(id CoordinatorConnectionID) error {
	for i, c := range m.coordinators {
		if c.ID == id {
			m.coordinators = removeCoordRef(m.coordinators, i)
		}
	}
	return nil
}

func (m *MockPool) Update(id CoordinatorConnectionID, state models2.Coordinator) error {
	for _, c := range m.coordinators {
		if c.ID == id {
			c.Coordinator.Tags = state.Tags
			c.Coordinator.Enabled = state.Enabled
		}
	}
	return nil
}

func removeCoordRef(s []*CoordinatorRef, i int) []*CoordinatorRef {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
