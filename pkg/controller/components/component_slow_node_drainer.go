package components

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/go-redis/redis/v8"
	"sync"
)

type SlowNodeDrainerConf struct {
	DrainThreshold int
}

type SlowNodeDrainer struct {
	analyzer       SlowNodeAnalyzer
	nodeDrainer    NodeDrainer
	slowNodeMarker SlowNodeMarker
	conf           SlowNodeDrainerConf
	logger         logging.Logger
}

func NewSlowNodeDrainer(analyzer SlowNodeAnalyzer, nodeDrainer NodeDrainer, slowNodeMarker SlowNodeMarker, conf SlowNodeDrainerConf, logger logging.Logger) *SlowNodeDrainer {
	return &SlowNodeDrainer{analyzer: analyzer, nodeDrainer: nodeDrainer, slowNodeMarker: slowNodeMarker, conf: conf, logger: logger}
}

func (s SlowNodeDrainer) Execute(ctx context.Context, detail trino.QueryDetail) error {
	slowNodes, err := s.analyzer.Analyze(detail)
	if err != nil {
		return err
	}

	for _, node := range slowNodes {
		isDrainable, err := s.markSlowNode(ctx, node)
		if err != nil {
			return err
		}

		if !isDrainable {
			continue
		}

		// TODO Check if this is blocking
		err = s.nodeDrainer.Drain(ctx, node.NodeID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SlowNodeDrainer) markSlowNode(ctx context.Context, node SlowNodeRef) (bool, error) {
	val, err := s.slowNodeMarker.Mark(ctx, node.NodeID)
	return val > s.conf.DrainThreshold, err
}

type SlowNodeMarker interface {
	Mark(ctx context.Context, nodeName string) (int, error)
}

type InMemorySlowNodeMarker struct {
	status map[string]int
	l      *sync.Mutex
}

func NewInMemoryMarker() *InMemorySlowNodeMarker {
	return &InMemorySlowNodeMarker{
		status: make(map[string]int),
		l:      &sync.Mutex{},
	}
}

func (i InMemorySlowNodeMarker) Mark(ctx context.Context, nodeName string) (int, error) {
	i.l.Lock()
	defer i.l.Unlock()

	val := i.status[nodeName]
	val = val + 1
	i.status[nodeName] = val

	return val, nil
}

type RedisSlowNodeMarker struct {
	redis redis.UniversalClient
}

func NewRedisSlowNodeMarker(redis redis.UniversalClient) *RedisSlowNodeMarker {
	return &RedisSlowNodeMarker{redis: redis}
}

func (r RedisSlowNodeMarker) Mark(ctx context.Context, nodeName string) (int, error) {
	val, err := r.redis.Get(ctx, fmt.Sprintf("component:%s:%s", "slow-node-marker", nodeName)).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
