package autoscaler

import (
	"time"
)

type State interface {
	LastQueryExecution(clusterID string) (time.Time, error)
	SetLastQueryExecution(clusterID string, t time.Time) error

	CurrentWorker(clusterID string) (int, error)
	SetCurrentWorker(clusterID string, currentWorker int) error
}

type InMemory struct {
	stateLastQuery     map[string]time.Time
	stateCurrentWorker map[string]int
}

func MemoryState() *InMemory {
	return &InMemory{
		stateLastQuery:     make(map[string]time.Time),
		stateCurrentWorker: make(map[string]int),
	}
}

func (i *InMemory) LastQueryExecution(clusterID string) (time.Time, error) {
	return i.stateLastQuery[clusterID], nil
}

func (i *InMemory) SetLastQueryExecution(clusterID string, t time.Time) error {
	i.stateLastQuery[clusterID] = t
	return nil
}

func (i *InMemory) CurrentWorker(clusterID string) (int, error) {
	return i.stateCurrentWorker[clusterID], nil
}

func (i *InMemory) SetCurrentWorker(clusterID string, currentWorker int) error {
	i.stateCurrentWorker[clusterID] = currentWorker
	return nil
}

type mockState struct {
	set func(clusterID string, t time.Time) error
	get func(clusterID string) (time.Time, error)

	setCW func(clusterID string, currentWorker int) error
	getCW func(clusterID string) (int, error)
}

func (m mockState) LastQueryExecution(clusterID string) (time.Time, error) {
	return m.get(clusterID)
}

func (m mockState) SetLastQueryExecution(clusterID string, t time.Time) error {
	return m.set(clusterID, t)
}

func (m mockState) CurrentWorker(clusterID string) (int, error) {
	return m.getCW(clusterID)
}

func (m mockState) SetCurrentWorker(clusterID string, currentWorker int) error {
	return m.setCW(clusterID, currentWorker)
}
