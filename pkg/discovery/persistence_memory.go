package discovery

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

// this is just for single node usage / testing purpose DO NOT use in production
type MemoryStorage struct {
	status []models.Coordinator
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		status: make([]models.Coordinator, 0),
	}
}

func (m *MemoryStorage) Remove(ctx context.Context, name string) error {
	for i, s := range m.status {
		if s.Name == name {
			m.status = remove(m.status, i)
		}
	}
	return nil
}

func (m *MemoryStorage) Add(ctx context.Context, coordinator models.Coordinator) error {
	m.status = append(m.status, coordinator)
	return nil
}

func (m *MemoryStorage) Get(ctx context.Context, name string) (models.Coordinator, error) {
	for _, s := range m.status {
		if s.Name == name {
			return s, nil
		}
	}

	return models.Coordinator{}, ErrClusterNotFound
}

func (m *MemoryStorage) All(ctx context.Context) ([]models.Coordinator, error) {
	return m.status, nil
}

func remove(s []models.Coordinator, i int) []models.Coordinator {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
