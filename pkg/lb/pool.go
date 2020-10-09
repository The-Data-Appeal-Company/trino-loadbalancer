package lb

import (
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/proxy"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/session"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/statistics"
	"sync"
	"time"
)

var (
	ErrNoBackendsAvailable = errors.New("no backend available")
)

type CoordinatorConnection struct {
	Proxy      proxy.HttpProxy
	Backend    models.Coordinator
	health     healthcheck.Health
	statistics models.ClusterStatistics
	termHc     chan bool
	termStats  chan bool
	stateMutex *sync.Mutex
}

type HttpPool interface {
	AllBackends() []*CoordinatorConnection
	AvailableBackends() ([]*CoordinatorConnection, error)
	GetByName(name string, unhealthy healthcheck.HealthStatus) (*CoordinatorConnection, error)
	Add(coordinator models.Coordinator) error
	Remove(name string) error
	UpdateStatus() error
	Update(name string, state models.Coordinator) error
}

type PoolConfig struct {
	HealthCheckDelay time.Duration
	StatisticsDelay  time.Duration
}

type Pool struct {
	conf               PoolConfig
	logger             logging.Logger
	sessionStore       session.Storage
	coordinators       []*CoordinatorConnection
	healthChecker      healthcheck.HealthCheck
	statisticRetriever statistics.Retriever
	rwLock             *sync.RWMutex
}

func NewPool(conf PoolConfig, sessionStore session.Storage, hc healthcheck.HealthCheck, statisticRetriever statistics.Retriever, logger logging.Logger) *Pool {
	return &Pool{
		conf:               conf,
		statisticRetriever: statisticRetriever,
		sessionStore:       sessionStore,
		logger:             logger,
		healthChecker:      hc,
		coordinators:       make([]*CoordinatorConnection, 0),
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

func (p *Pool) AllBackends() []*CoordinatorConnection {
	return p.coordinators
}

func (p *Pool) AvailableBackends() ([]*CoordinatorConnection, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()

	available := make([]*CoordinatorConnection, 0)
	for _, conn := range p.coordinators {
		if conn.Backend.Enabled && conn.health.IsAvailable() {
			available = append(available, conn)
		}
	}

	if len(available) == 0 {
		return nil, ErrNoBackendsAvailable
	}

	return available, nil
}

func (p *Pool) GetByName(name string, includedStatus healthcheck.HealthStatus) (*CoordinatorConnection, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()

	for _, conn := range p.coordinators {
		if name == conn.Backend.Name && conn.health.Status >= includedStatus {
			return conn, nil
		}
	}

	return nil, ErrNoBackendsAvailable
}

func (p *Pool) Update(name string, state models.Coordinator) error {
	target, err := p.GetByName(name, healthcheck.StatusUnknown)

	if err != nil {
		return err
	}

	p.rwLock.Lock()

	target.Backend.Tags = state.Tags
	target.Backend.Enabled = state.Enabled
	target.Backend.Distribution = state.Distribution // this is very unlikely

	p.rwLock.Unlock()
	return nil
}

func (p *Pool) Add(coordinator models.Coordinator) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	for _, c := range p.coordinators {
		if c.Backend.Name == coordinator.Name {
			return fmt.Errorf("duplicated backend name: %s", coordinator.Name)
		}
	}

	backendConn := &CoordinatorConnection{
		Backend:    coordinator,
		Proxy:      proxy.NewReverseProxy(coordinator.URL, NewQueryClusterLinker(p.sessionStore, coordinator.Name)),
		termHc:     make(chan bool),
		termStats:  make(chan bool),
		stateMutex: &sync.Mutex{},
	}

	p.coordinators = append(p.coordinators, backendConn)
	p.updateBackendHealth(backendConn)

	go p.healthCheck(backendConn)
	go p.clusterStatistics(backendConn)

	p.logger.Info("added backend to pool: %s ( %s )", coordinator.Name, coordinator.URL.String())
	return nil
}

func (p *Pool) Remove(name string) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	for i, conn := range p.coordinators {
		if conn.Backend.Name == name {
			conn.termHc <- true
			conn.termStats <- true

			conn.health = healthcheck.Health{
				Timestamp: time.Now(),
				Status:    healthcheck.StatusUnhealthy,
				Message:   "backend has been removed from the pool",
			}

			p.logger.Info("removed backend from pool: %s ( %s )", name, conn.Backend)

			p.coordinators = remove(p.coordinators, i)
			return nil
		}
	}

	return errors.New("backend not found")
}

func remove(s []*CoordinatorConnection, i int) []*CoordinatorConnection {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (p *Pool) clusterStatistics(b *CoordinatorConnection) {
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

func (p *Pool) healthCheck(b *CoordinatorConnection) {
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

func (p *Pool) updateBackendHealth(b *CoordinatorConnection) {
	b.stateMutex.Lock()
	defer b.stateMutex.Unlock()
	result, err := p.healthChecker.Check(b.Backend.URL)
	if err != nil {
		result = healthcheck.Health{
			Timestamp: time.Now(),
			Status:    healthcheck.StatusUnhealthy,
			Message:   err.Error(),
		}
	}

	if b.health.Status != result.Status {
		p.logger.Warn("%s health status changed %s -> %s", b.Backend.Name, b.health.Status.String(), result.Status.String())
	}

	b.health = result
}

func (p *Pool) updateBackendStatistics(b *CoordinatorConnection) {
	b.stateMutex.Lock()
	defer b.stateMutex.Unlock()

	if !b.health.IsAvailable() {
		return
	}

	stats, err := p.statisticRetriever.GetStatistics(b.Backend)
	if err != nil {
		p.logger.Warn("no statistics available for %s: %s", b.Backend.Name, err.Error())
		return
	}

	b.statistics = stats
}
