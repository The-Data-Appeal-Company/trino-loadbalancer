package autoscaler

import (
	"time"
)

type State interface {
	LastQueryExecution(clusterID string) (time.Time, error)
	SetLastQueryExecution(clusterID string, t time.Time) error
}

type InMemory struct {
	state map[string]time.Time
}

func MemoryState() *InMemory {
	return &InMemory{state: make(map[string]time.Time)}
}

func (i *InMemory) LastQueryExecution(clusterID string) (time.Time, error) {
	return i.state[clusterID], nil
}

func (i *InMemory) SetLastQueryExecution(clusterID string, t time.Time) error {
	i.state[clusterID] = t
	return nil
}

type mockState struct {
	set func(clusterID string, t time.Time) error
	get func(clusterID string) (time.Time, error)
}

func (m mockState) LastQueryExecution(clusterID string) (time.Time, error) {
	return m.get(clusterID)
}

func (m mockState) SetLastQueryExecution(clusterID string, t time.Time) error {
	return m.set(clusterID, t)
}
