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
	SetLastScaleUp(clusterID string, i int32, t time.Time) error
	GetLastScaleUp(clusterID string) (int32, time.Time, error)
}
type InMemory struct {
	stateTime      map[string]time.Time
	stateInstances map[string]int32
	stateScaleUp   map[string]ScaleUpState
}

type ScaleUpState struct {
	instances int32
	time      time.Time
}

func MemoryState() *InMemory {
	return &InMemory{
		stateTime:      make(map[string]time.Time),
		stateInstances: make(map[string]int32),
		stateScaleUp:   make(map[string]ScaleUpState),
	}
}

var NoInstancesInStateError = errors.New("no instances state")
var NoLastScaleUpStateError = errors.New("no Last Scale Up state")

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

func (i *InMemory) GetLastScaleUp(clusterID string) (int32, time.Time, error) {
	value, ok := i.stateScaleUp[clusterID]
	if !ok {
		return 0, time.Now(), NoLastScaleUpStateError
	}
	return value.instances, value.time, nil
}

func (i *InMemory) SetLastScaleUp(clusterID string, ii int32, t time.Time) error {
	state := ScaleUpState{
		instances: ii,
		time:      t,
	}
	i.stateScaleUp[clusterID] = state
	return nil
}

type mockState struct {
	setTime        func(clusterID string, t time.Time) error
	getTime        func(clusterID string) (time.Time, error)
	setInstances   func(clusterID string, i int32) error
	getInstances   func(clusterID string) (int32, error)
	setLastScaleUp func(clusterID string, i int32, t time.Time) error
	getLastScaleUp func(clusterID string) (int32, time.Time, error)
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

func (m mockState) GetLastScaleUp(clusterID string) (int32, time.Time, error) {
	return m.getLastScaleUp(clusterID)
}

func (m mockState) SetLastScaleUp(clusterID string, ii int32, t time.Time) error {
	return m.setLastScaleUp(clusterID, ii, t)
}
