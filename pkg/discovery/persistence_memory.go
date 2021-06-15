package discovery

import (
	"context"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

// this is just for single node usage / testing purpose DO NOT use in production
type MemoryStorage struct {
	status []models2.Coordinator
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		status: make([]models2.Coordinator, 0),
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

func (m *MemoryStorage) Add(ctx context.Context, coordinator models2.Coordinator) error {
	m.status = append(m.status, coordinator)
	return nil
}

func (m *MemoryStorage) Get(ctx context.Context, name string) (models2.Coordinator, error) {
	for _, s := range m.status {
		if s.Name == name {
			return s, nil
		}
	}

	return models2.Coordinator{}, ErrClusterNotFound
}

func (m *MemoryStorage) All(ctx context.Context) ([]models2.Coordinator, error) {
	return m.status, nil
}

func remove(s []models2.Coordinator, i int) []models2.Coordinator {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
