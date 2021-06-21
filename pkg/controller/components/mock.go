package components

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"sync"
)

type MockHandler struct {
	l     *sync.Mutex
	calls []trino.QueryDetail
	err   error
}

func NewMockQueryHandler(err error) *MockHandler {
	return &MockHandler{
		l:     &sync.Mutex{},
		calls: make([]trino.QueryDetail, 0),
		err:   err,
	}
}

func (m *MockHandler) Execute(ctx context.Context, detail trino.QueryDetail) error {
	m.l.Lock()
	defer m.l.Unlock()

	m.calls = append(m.calls, detail)
	return m.err
}

func (m *MockHandler) Calls() []trino.QueryDetail {
	m.l.Lock()
	defer m.l.Unlock()
	return m.calls
}
