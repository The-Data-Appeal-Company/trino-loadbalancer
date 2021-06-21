package components

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/notifier"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/go-redis/redis/v8"
	"sync"
)

type SlowNodeDrainerConf struct {
	DrainThreshold int64
}

type SlowNodeDrainer struct {
	analyzer       SlowNodeAnalyzer
	nodeDrainer    NodeDrainer
	slowNodeMarker SlowNodeMarker
	conf           SlowNodeDrainerConf
	notifier       notifier.Notifier
	logger         logging.Logger
}

func NewSlowNodeDrainer(analyzer SlowNodeAnalyzer, nodeDrainer NodeDrainer, slowNodeMarker SlowNodeMarker, conf SlowNodeDrainerConf, logger logging.Logger, notifier notifier.Notifier) *SlowNodeDrainer {
	return &SlowNodeDrainer{analyzer: analyzer, nodeDrainer: nodeDrainer, slowNodeMarker: slowNodeMarker, conf: conf, logger: logger, notifier: notifier}
}

func (s SlowNodeDrainer) Execute(ctx context.Context, detail trino.QueryDetail) error {
	slowNodes, err := s.analyzer.Analyze(detail)
	if err != nil {
		return err
	}

	for _, node := range slowNodes {
		s.logger.Info("%s marked as slow", node)
		isDrainable, err := s.markSlowNode(ctx, node)
		if err != nil {
			return err
		}

		if !isDrainable {
			continue
		}

		state, err := s.slowNodeMarker.State(ctx, node.NodeID)
		if err != nil {
			// if the node status in unknown that means that is running
			if err != ErrNodeStateNotFound {
				return err
			}
			state = SlowNodeState{
				Status: SlowNodeStateRunning,
			}
		}

		if state.Status == SlowNodeStateDraining {
			s.logger.Info("throttling drain request for %s", node)
			continue
		}

		s.logger.Info("%s selected for drain", node)

		// TODO Check if this is blocking
		err = s.nodeDrainer.Drain(ctx, node.NodeID)
		if err != nil {
			return err
		}

		if err := s.notifier.Notify(notifier.Request{
			Title:   "Slow worker node",
			Message: "draining slow worker node",
			Metadata: map[string]string{
				"node":  node.NodeID,
				"query": detail.QueryID,
			},
		}); err != nil {
			s.logger.Warn("error notifying node drain: %s", err.Error())
		}

		// set node status to Drain to prevent multiple drain requests
		if err := s.slowNodeMarker.Update(ctx, node.NodeID, SlowNodeState{Status: SlowNodeStateDraining}); err != nil {
			return err
		}

	}

	return nil
}

func (s SlowNodeDrainer) markSlowNode(ctx context.Context, node SlowNodeRef) (bool, error) {
	val, err := s.slowNodeMarker.Mark(ctx, node.NodeID)
	return val >= s.conf.DrainThreshold, err
}

var (
	ErrNodeStateNotFound = errors.New("node state not present")
)

const (
	SlowNodeStateRunning  = "Running"
	SlowNodeStateDraining = "Draining"
)

type SlowNodeState struct {
	Status string
}

type SlowNodeMarker interface {
	Update(ctx context.Context, nodeName string, state SlowNodeState) error
	State(ctx context.Context, nodeName string) (SlowNodeState, error)
	Mark(ctx context.Context, nodeName string) (int64, error)
	Delete(ctx context.Context, nodeName string) error
}

type InMemorySlowNodeMarker struct {
	marks  map[string]int64
	status map[string]SlowNodeState
	l      *sync.Mutex
}

func NewInMemorySlowNodeMarker() *InMemorySlowNodeMarker {
	return &InMemorySlowNodeMarker{
		marks:  make(map[string]int64),
		status: make(map[string]SlowNodeState),
		l:      &sync.Mutex{},
	}
}

func (i *InMemorySlowNodeMarker) Mark(ctx context.Context, nodeName string) (int64, error) {
	i.l.Lock()
	defer i.l.Unlock()

	val := i.marks[nodeName]
	val = val + 1
	i.marks[nodeName] = val

	return val, nil
}

func (i *InMemorySlowNodeMarker) Delete(ctx context.Context, nodeName string) error {
	i.l.Lock()
	defer i.l.Unlock()

	delete(i.marks, nodeName)
	return nil
}

func (i *InMemorySlowNodeMarker) Update(ctx context.Context, nodeName string, state SlowNodeState) error {
	i.l.Lock()
	defer i.l.Unlock()
	i.status[nodeName] = state
	return nil
}

func (i *InMemorySlowNodeMarker) State(ctx context.Context, nodeName string) (SlowNodeState, error) {
	i.l.Lock()
	defer i.l.Unlock()
	state, present := i.status[nodeName]
	if !present {
		return SlowNodeState{}, ErrNodeStateNotFound
	}
	return state, nil
}

type RedisSlowNodeMarker struct {
	redis redis.UniversalClient
}

func NewRedisSlowNodeMarker(redis redis.UniversalClient) *RedisSlowNodeMarker {
	return &RedisSlowNodeMarker{redis: redis}
}

func (r RedisSlowNodeMarker) Mark(ctx context.Context, nodeName string) (int64, error) {
	var key = fmt.Sprintf("component:%s:%s", "slow-node-marker", nodeName)
	return r.redis.Incr(ctx, key).Result()
}

func (r RedisSlowNodeMarker) Delete(ctx context.Context, nodeName string) error {
	return r.redis.Del(ctx, nodeName).Err()
}

func (r RedisSlowNodeMarker) Update(ctx context.Context, nodeName string, state SlowNodeState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	var key = fmt.Sprintf("component:%s:%s", "slow-node-marker-node-status", nodeName)
	return r.redis.Set(ctx, key, data, -1).Err()
}

func (r RedisSlowNodeMarker) State(ctx context.Context, nodeName string) (SlowNodeState, error) {
	var key = fmt.Sprintf("component:%s:%s", "slow-node-marker-node-status", nodeName)

	data, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return SlowNodeState{}, ErrNodeStateNotFound
		}
		return SlowNodeState{}, err
	}

	var state SlowNodeState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return SlowNodeState{}, err
	}

	return state, nil
}
