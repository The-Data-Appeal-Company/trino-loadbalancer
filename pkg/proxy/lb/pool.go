package lb

import (
	"errors"
	"fmt"
	trino2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	logging2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	healthcheck2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	http2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/http"
	session2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

var (
	ErrNoBackendsAvailable = errors.New("no backend available")
)

type CoordinatorConnectionID string

type coordinatorConnection struct {
	proxy       http2.HttpProxy
	coordinator models2.Coordinator
	health     healthcheck2.Health
	statistics trino2.ClusterStatistics
	termHc     chan bool
	termStats   chan bool
	stateMutex  *sync.Mutex
}

type CoordinatorRef struct {
	ID         CoordinatorConnectionID
	Statistics trino2.ClusterStatistics

	models2.Coordinator
}

type TrinoPool interface {
	Handle(coordinator CoordinatorRef, writer http.ResponseWriter, request *http.Request) error
	Fetch(FetchRequest) []CoordinatorRef
	Add(models2.Coordinator) error
	Remove(CoordinatorConnectionID) error
	Update(CoordinatorConnectionID, models2.Coordinator) error
}

type PoolConfig struct {
	HealthCheckDelay time.Duration
	StatisticsDelay  time.Duration
}

type Pool struct {
	conf         PoolConfig
	logger       logging2.Logger
	sessionStore session2.Storage
	coordinators       map[CoordinatorConnectionID]*coordinatorConnection
	healthChecker      healthcheck2.HealthCheck
	statisticRetriever trino2.Api
	rwLock             *sync.RWMutex
}

func NewPool(conf PoolConfig, sessionStore session2.Storage, hc healthcheck2.HealthCheck, statisticRetriever trino2.Api, logger logging2.Logger) *Pool {
	return &Pool{
		conf:               conf,
		statisticRetriever: statisticRetriever,
		sessionStore:       sessionStore,
		logger:             logger,
		healthChecker:      hc,
		coordinators:       make(map[CoordinatorConnectionID]*coordinatorConnection),
		rwLock:             &sync.RWMutex{},
	}
}

func (p *Pool) UpdateStatus() error {
	for _, c := range p.coordinators {
		p.updateBackendHealth(c)
		p.updateBackendStatistics(c)
	}

	return nil
}

type EnabledStatus int

const (
	ClusterStatusAll      EnabledStatus = 0
	ClusterStatusEnabled  EnabledStatus = 1
	ClusterStatusDisabled EnabledStatus = 2
)

type FetchRequest struct {
	Name   string
	Tags   map[string]string
	Health healthcheck2.HealthStatus
	Status EnabledStatus
}

func (p *Pool) Fetch(req FetchRequest) []CoordinatorRef {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()

	selected := make([]CoordinatorRef, 0)
	for id, cc := range p.coordinators {
		if len(req.Name) != 0 && cc.coordinator.Name != req.Name {
			continue
		}

		if req.Health != 0 && cc.health.Status < req.Health {
			continue
		}

		if len(req.Tags) != 0 && !matchTags(cc.coordinator.Tags, req.Tags) {
			continue
		}

		if req.Status == ClusterStatusEnabled && !cc.coordinator.Enabled || req.Status == ClusterStatusDisabled && cc.coordinator.Enabled {
			continue
		}

		selected = append(selected, CoordinatorRef{
			ID:          id,
			Coordinator: cc.coordinator,
			Statistics:  cc.statistics,
		})
	}

	return selected
}

func (p *Pool) Update(id CoordinatorConnectionID, state models2.Coordinator) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	target, err := p.connectionByID(id)
	if err != nil {
		return err
	}

	target.coordinator.Tags = state.Tags
	target.coordinator.Enabled = state.Enabled
	return nil
}

func (p *Pool) connectionByID(id CoordinatorConnectionID) (*coordinatorConnection, error) {
	conn, present := p.coordinators[id]
	if !present {
		return nil, ErrNoBackendsAvailable
	}
	return conn, nil
}

func (p *Pool) Add(coordinator models2.Coordinator) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	for _, c := range p.coordinators {
		if c.coordinator.Name == coordinator.Name {
			return fmt.Errorf("duplicated backend name: %s", coordinator.Name)
		}
	}

	connectionID := CoordinatorConnectionID(uuid.New().String())
	backendConn := &coordinatorConnection{
		coordinator: coordinator,
		proxy: http2.NewReverseProxy(coordinator.URL, http2.NewCompositeInterceptor(
			NewQueryClusterLinker(p.sessionStore, coordinator.Name),
		)),
		termHc:     make(chan bool),
		termStats:  make(chan bool),
		stateMutex: &sync.Mutex{},
	}

	p.coordinators[connectionID] = backendConn

	p.updateBackendHealth(backendConn)

	go p.healthCheck(backendConn)
	go p.clusterStatistics(backendConn)

	p.logger.Info("added backend to pool: %s ( %s )", coordinator.Name, coordinator.URL.String())
	return nil
}

func (p *Pool) Remove(id CoordinatorConnectionID) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	conn, present := p.coordinators[id]

	if !present {
		return errors.New("backend not found")
	}

	conn.termHc <- true
	conn.termStats <- true

	conn.health = healthcheck2.Health{
		Timestamp: time.Now(),
		Status:    healthcheck2.StatusUnhealthy,
		Message:   "backend has been removed from the pool",
	}

	p.logger.Info("removed backend from pool: %s ( %s )", conn.coordinator.Name, conn.coordinator)

	delete(p.coordinators, id)
	return nil
}

func (p *Pool) clusterStatistics(b *coordinatorConnection) {
	ticker := time.NewTicker(p.conf.StatisticsDelay)
	for {
		select {
		case <-ticker.C:
			p.updateBackendStatistics(b)
		case <-b.termStats:
			return
		}
	}
}

func (p *Pool) healthCheck(b *coordinatorConnection) {
	ticker := time.NewTicker(p.conf.HealthCheckDelay)
	for {
		select {
		case <-ticker.C:
			p.updateBackendHealth(b)
		case <-b.termHc:
			return
		}
	}
}

func (p *Pool) updateBackendHealth(b *coordinatorConnection) {
	b.stateMutex.Lock()
	defer b.stateMutex.Unlock()
	result, err := p.healthChecker.Check(b.coordinator.URL)
	if err != nil {
		result = healthcheck2.Health{
			Timestamp: time.Now(),
			Status:    healthcheck2.StatusUnhealthy,
			Message:   err.Error(),
		}
	}

	if b.health.Status != result.Status {
		p.logger.Warn("%s health status changed %s -> %s", b.coordinator.Name, b.health.Status.String(), result.Status.String())
	}

	b.health = result
}

func (p *Pool) updateBackendStatistics(b *coordinatorConnection) {
	b.stateMutex.Lock()
	defer b.stateMutex.Unlock()

	if !b.health.IsAvailable() {
		return
	}

	stats, err := p.statisticRetriever.ClusterStatistics(b.coordinator)
	if err != nil {
		p.logger.Warn("no statistics available for %s: %s", b.coordinator.Name, err.Error())
		return
	}

	b.statistics = stats
}

func (p *Pool) Handle(coordinator CoordinatorRef, writer http.ResponseWriter, request *http.Request) error {
	conn, err := p.connectionByID(coordinator.ID)
	if err != nil {
		return err
	}
	return conn.proxy.Handle(writer, request)
}

func matchTags(source map[string]string, match map[string]string) bool {
	for k, v := range match {
		if source[k] != v {
			return false
		}
	}
	return true
}
