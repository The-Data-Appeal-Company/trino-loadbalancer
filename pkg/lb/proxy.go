package lb

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/routing"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/session"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ProxyConf struct {
	SyncDelay time.Duration
}

type PrestoProxy struct {
	conf          ProxyConf
	logger        logging.Logger
	pool          *Pool
	sessionReader session.Reader
	router        routing.Router
	poolSync      PoolSync
	termSync      chan bool
}

func NewPrestoProxy(conf ProxyConf, pool *Pool, sync PoolSync, sessReader session.Reader, router routing.Router, logger logging.Logger) *PrestoProxy {
	return &PrestoProxy{
		conf:          conf,
		poolSync:      sync,
		router:        router,
		logger:        logger,
		pool:          pool,
		sessionReader: sessReader,
		termSync:      make(chan bool),
	}
}

func (p *PrestoProxy) Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	r.PathPrefix("/").Handler(http.HandlerFunc(p.Handle))
	return r
}

func (p *PrestoProxy) Init() error {
	// TODO We can start the proxy when we have at least 1 available cluster
	if err := p.poolSync.Sync(p.pool); err != nil {
		return err
	}

	go p.syncPoolState()
	return nil
}

func (p *PrestoProxy) Serve(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: p.Router(),
	}
	return srv.ListenAndServe()
}

func (p *PrestoProxy) Handle(writer http.ResponseWriter, request *http.Request) {
	coordinator, err := p.selectCoordinatorForRequest(request)
	if err == ErrNoBackendsAvailable {
		p.logger.Warn("no available backends for request %s", request.URL)
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	if err := coordinator.Proxy.Handle(writer, request); err != nil {
		p.logger.Error("error handling request %s: %s", request.URL, err.Error())
	}
}

func (p *PrestoProxy) selectCoordinatorForRequest(request *http.Request) (*CoordinatorConnection, error) {

	// the request is not query related OR the request is a query submission
	// we can apply the user selected request routing algorithm
	if !isStatementRequest(request.URL) || isStatementRequest(request.URL) && request.Method == http.MethodPost {
		backends, err := p.pool.AvailableBackends()

		if err != nil {
			return nil, ErrNoBackendsAvailable
		}

		selectedBackend, err := p.router.Route(routingRequest(backends, request))
		if err != nil {
			return nil, err
		}

		coordinator, err := p.pool.GetByName(selectedBackend.Name, 0)
		if err != nil {
			return nil, ErrNoBackendsAvailable
		}

		return coordinator, nil
	}

	// the request is retrieving info about a specific query we must get coordinator with planned the query
	// so we use the sessionReader to retrieve its name
	if isStatementRequest(request.URL) && request.Method == http.MethodGet { //todo check path
		queryInfo, err := QueryInfoFromRequest(request)
		if err != nil {
			return nil, err
		}

		coordinatorName, err := p.sessionReader.Get(request.Context(), queryInfo)
		if err != nil {
			return nil, err
		}

		backend, err := p.getBackendByName(coordinatorName)

		return backend, err
	}

	return nil, ErrNoBackendsAvailable
}

// retrieve backend by name, if not present force cluster status sync for the pool and then try again to fetch the request backend,
func (p *PrestoProxy) getBackendByName(name string) (*CoordinatorConnection, error) {
	backend, err := p.pool.GetByName(name, healthcheck.StatusUnhealthy)
	if err != nil {
		// If the pool doesn't have a backend with the specified name we force a state refresh to be sure that we
		// are aligned with other proxies that may have discovered new clusters that the current proxy is not aware of.
		// since sync / add perform storage lookup and health check of the added backend this may slow down the first query triggering
		// the update, but this case shouldn't happen very often
		if err == ErrNoBackendsAvailable {
			p.logger.Info("no backend with name %s found for query forcing state sync", name)
			if synErr := p.poolSync.Sync(p.pool); synErr != nil {
				return nil, fmt.Errorf("no backend found for name %s, unable to sync pool: %w", name, err)
			}
			backend, err = p.pool.GetByName(name, 0)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return backend, nil
}

func routingRequest(backends []*CoordinatorConnection, req *http.Request) routing.Request {
	coordinatorsWithStatistics := make([]routing.CoordinatorWithStatistics, len(backends))
	for i, backend := range backends {
		coordinatorsWithStatistics[i] = routing.CoordinatorWithStatistics{
			Coordinator: backend.Backend,
			Statistics:  backend.statistics,
		}
	}

	return routing.Request{
		Coordinators: coordinatorsWithStatistics,
		User:         req.Header.Get(TrinoHeaderUser),
	}
}

func (p *PrestoProxy) syncPoolState() {
	ticker := time.NewTicker(p.conf.SyncDelay)
	for {
		select {
		case <-ticker.C:
			if err := p.poolSync.Sync(p.pool); err != nil {
				p.logger.Error("error syncing state: %s", err.Error())
			}
		case <-p.termSync:
			return
		}
	}
}

func (p *PrestoProxy) Close() error {
	p.termSync <- true
	return nil
}
