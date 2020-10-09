package session

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

type Memory struct {
	sessions map[string]string
}

func NewMemoryStorage() *Memory {
	return &Memory{
		sessions: make(map[string]string),
	}
}

func (m *Memory) Link(ctx context.Context, info models.QueryInfo, s string) error {
	hash := m.queryHash(info)
	m.sessions[hash] = s
	return nil
}

func (m *Memory) Unlink(ctx context.Context, info models.QueryInfo) error {
	hash := m.queryHash(info)
	delete(m.sessions, hash)
	return nil
}

func (m *Memory) Get(ctx context.Context, info models.QueryInfo) (string, error) {
	hash := m.queryHash(info)
	val, present := m.sessions[hash]

	if !present {
		return "", ErrLinkNotFound
	}

	return val, nil
}

func (m *Memory) queryHash(info models.QueryInfo) string {
	return fmt.Sprintf("%s::%s", info.TransactionID, info.QueryID)
}
