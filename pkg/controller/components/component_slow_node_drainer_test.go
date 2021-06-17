package components

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestSlowNodeDrainerComponent(t *testing.T) {
	const nodeToDrain = "node-00"

	analyzer := newMockAnalyzer(func(detail trino.QueryDetail) ([]SlowNodeRef, error) {
		return []SlowNodeRef{
			{
				NodeID: nodeToDrain,
			},
		}, nil
	})

	slowNodeMarker := NewInMemorySlowNodeMarker()
	nodeDrainer := newMockDrainer(nil)
	conf := SlowNodeDrainerConf{
		DrainThreshold: 1,
	}

	drainer := NewSlowNodeDrainer(analyzer, nodeDrainer, slowNodeMarker, conf, logging.Noop())

	ctx := context.TODO()
	err := drainer.Execute(ctx, trino.QueryDetail{})
	require.NoError(t, err)

	marked, present := nodeDrainer.drained[nodeToDrain]

	require.True(t, present)
	require.Equal(t, marked, 1)

}

func TestSlowNodeDrainerComponentReturnErrOnAnalyzerError(t *testing.T) {
	analyzer := newMockAnalyzer(func(detail trino.QueryDetail) ([]SlowNodeRef, error) {
		return nil, errors.New("error during query analyze")
	})

	slowNodeMarker := NewInMemorySlowNodeMarker()
	nodeDrainer := newMockDrainer(nil)
	conf := SlowNodeDrainerConf{}

	drainer := NewSlowNodeDrainer(analyzer, nodeDrainer, slowNodeMarker, conf, logging.Noop())

	ctx := context.TODO()
	err := drainer.Execute(ctx, trino.QueryDetail{})
	require.Error(t, err)
}

func TestSlowNodeDrainerComponentReturnErrOnDrainerError(t *testing.T) {
	analyzer := newMockAnalyzer(func(detail trino.QueryDetail) ([]SlowNodeRef, error) {
		return []SlowNodeRef{
			{
				NodeID: "node-00",
			},
		}, nil
	})

	slowNodeMarker := NewInMemorySlowNodeMarker()
	nodeDrainer := newMockDrainer(errors.New("node drain error"))
	conf := SlowNodeDrainerConf{}

	drainer := NewSlowNodeDrainer(analyzer, nodeDrainer, slowNodeMarker, conf, logging.Noop())

	ctx := context.TODO()
	err := drainer.Execute(ctx, trino.QueryDetail{})
	require.Error(t, err)
}

func TestSlowNodeDrainerComponentNoActionOnEmptySlowNodes(t *testing.T) {
	analyzer := newMockAnalyzer(func(detail trino.QueryDetail) ([]SlowNodeRef, error) {
		return []SlowNodeRef{}, nil
	})

	slowNodeMarker := NewInMemorySlowNodeMarker()
	nodeDrainer := newMockDrainer(nil)
	conf := SlowNodeDrainerConf{}

	drainer := NewSlowNodeDrainer(analyzer, nodeDrainer, slowNodeMarker, conf, logging.Noop())

	ctx := context.TODO()
	err := drainer.Execute(ctx, trino.QueryDetail{})

	require.NoError(t, err)
	require.Empty(t, nodeDrainer.drained)
}

type mockNodeDrainer struct {
	drained map[string]int
	l       *sync.Mutex
	err     error
}

func newMockDrainer(err error) *mockNodeDrainer {
	return &mockNodeDrainer{
		drained: make(map[string]int),
		l:       &sync.Mutex{},
		err:     err,
	}
}

func (m *mockNodeDrainer) Drain(ctx context.Context, nodeID string) error {
	m.l.Lock()
	defer m.l.Unlock()

	val := m.drained[nodeID]
	m.drained[nodeID] = val + 1

	return m.err
}

type mockAnalyzer struct {
	fn func(detail trino.QueryDetail) ([]SlowNodeRef, error)
}

func newMockAnalyzer(fn func(detail trino.QueryDetail) ([]SlowNodeRef, error)) *mockAnalyzer {
	return &mockAnalyzer{fn: fn}
}

func (m mockAnalyzer) Analyze(detail trino.QueryDetail) ([]SlowNodeRef, error) {
	return m.fn(detail)
}
