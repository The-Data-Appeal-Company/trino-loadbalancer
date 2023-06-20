package autoscaler

import (
	"errors"
	"time"
)

type State interface {
	LastQueryExecution(clusterID string) (time.Time, error)
	SetLastQueryExecution(clusterID string, t time.Time) error
	GetClusterInstances(clusterID string) (int32, error)
	SetClusterInstances(clusterID string, i int32) error
}

type InMemory struct {
	stateTime      map[string]time.Time
	stateInstances map[string]int32
}

func MemoryState() *InMemory {
	return &InMemory{
		stateTime:      make(map[string]time.Time),
		stateInstances: make(map[string]int32),
	}
}

var NoInstancesInStateError = errors.New("no instances state")

func (i *InMemory) LastQueryExecution(clusterID string) (time.Time, error) {
	return i.stateTime[clusterID], nil
}

func (i *InMemory) SetLastQueryExecution(clusterID string, t time.Time) error {
	i.stateTime[clusterID] = t
	return nil
}

func (i *InMemory) GetClusterInstances(clusterID string) (int32, error) {
	value, ok := i.stateInstances[clusterID]
	if !ok {
		return 0, NoInstancesInStateError
	}
	return value, nil
}

func (i *InMemory) SetClusterInstances(clusterID string, instances int32) error {
	i.stateInstances[clusterID] = instances
	return nil
}

type mockState struct {
	setTime      func(clusterID string, t time.Time) error
	getTime      func(clusterID string) (time.Time, error)
	setInstances func(clusterID string, i int32) error
	getInstances func(clusterID string) (int32, error)
}

func (m mockState) LastQueryExecution(clusterID string) (time.Time, error) {
	return m.getTime(clusterID)
}

func (m mockState) SetLastQueryExecution(clusterID string, t time.Time) error {
	return m.setTime(clusterID, t)
}

func (m mockState) GetClusterInstances(clusterID string) (int32, error) {
	return m.getInstances(clusterID)
}

func (m mockState) SetClusterInstances(clusterID string, i int32) error {
	return m.setInstances(clusterID, i)
}
