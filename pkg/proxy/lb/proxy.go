package lb

import (
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/routing"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ProxyConf struct {
	SyncDelay time.Duration
}

type Proxy struct {
	conf            ProxyConf
	logger          logging.Logger
	pool            *Pool
	sessionReader   session.Reader
	router          routing.Router
	poolSync        PoolSync
	termSync        chan bool
	requestRewriter RequestRewriter
}

func NewProxy(conf ProxyConf, pool *Pool, sync PoolSync, sessReader session.Reader, router routing.Router, requestRewriter RequestRewriter, logger logging.Logger) *Proxy {
	return &Proxy{
		conf:            conf,
		poolSync:        sync,
		router:          router,
		logger:          logger,
		pool:            pool,
		sessionReader:   sessReader,
		termSync:        make(chan bool),
		requestRewriter: requestRewriter,
	}
}

func (p *Proxy) Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	r.PathPrefix("/").Handler(http.HandlerFunc(p.Handle))
	return r
}

func (p *Proxy) Init() error {
	// TODO We can start the proxy when we have at least 1 available cluster
	if err := p.poolSync.Sync(p.pool); err != nil {
		return err
	}

	go p.syncPoolState()
	return nil
}

func (p *Proxy) Serve(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: p.Router(),
	}
	return srv.ListenAndServe()
}

func (p *Proxy) Handle(writer http.ResponseWriter, request *http.Request) {
	coordinator, err := p.selectCoordinatorForRequest(request)
	if errors.Is(err, ErrNoBackendsAvailable) {
		p.logger.Warn("no available backends for request %s", request.URL)
		writer.WriteHeader(http.StatusServiceUnavailable)
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			p.logger.Error("error writing response: %w", err)
		}
		return
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			p.logger.Error("error writing response: %w", err)
		}
		return
	}

	if err := p.pool.Handle(coordinator, writer, request); err != nil {
		p.logger.Error("error handling request %s: %s", request.URL, err.Error())
	}
}

func (p *Proxy) selectCoordinatorForRequest(request *http.Request) (CoordinatorRef, error) {
	// the request is not query related OR the request is a query submission
	// we can apply the user selected request routing algorithm
	if !isStatementRequest(request.URL) || isStatementRequest(request.URL) && request.Method == http.MethodPost {
		request, err := p.requestRewriter.Rewrite(request)
		if err != nil {
			return CoordinatorRef{}, err
		}

		healthyCoordinators := p.pool.Fetch(FetchRequest{
			Health: healthcheck.StatusHealthy,
		})

		if len(healthyCoordinators) == 0 {
			return CoordinatorRef{}, ErrNoBackendsAvailable
		}

		targetCoordinator, err := p.router.Route(routingRequest(healthyCoordinators, request))
		if err != nil {
			return CoordinatorRef{}, err
		}

		return p.coordinatorRefByName(targetCoordinator.Name)
	}

	// the request is retrieving info about a specific query we must get coordinator with planned the query
	// so we use the sessionReader to retrieve its name
	if isStatementRequest(request.URL) && request.Method == http.MethodGet {
		queryInfo, err := queryInfoFromRequest(request)
		if err != nil {
			return CoordinatorRef{}, err
		}

		coordinatorName, err := p.sessionReader.Get(request.Context(), queryInfo)
		if err != nil {
			return CoordinatorRef{}, err
		}

		return p.coordinatorRefByName(coordinatorName)
	}

	return CoordinatorRef{}, ErrNoBackendsAvailable
}

// retrieve backend by name, if not present force cluster status sync for the pool and then try again to fetch the request backend,
func (p *Proxy) coordinatorRefByName(name string) (CoordinatorRef, error) {
	coordinator := p.pool.Fetch(FetchRequest{
		Name: name,
	})

	if len(coordinator) == 0 {
		// If the pool doesn't have a coordinator with the specified name we force a state refresh to be sure that we
		// are aligned with other proxies that may have discovered new clusters that the current proxy is not aware of.
		// since sync / add perform storage lookup and health check of the added coordinator this may slow down the first query triggering
		// the update, but this case shouldn't happen very often

		p.logger.Info("no coordinator with name %s found for query forcing state sync", name)
		if synErr := p.poolSync.Sync(p.pool); synErr != nil {
			return CoordinatorRef{}, fmt.Errorf("no coordinator found for name %s, unable to sync pool: %w", name, synErr)
		}

		coordinator = p.pool.Fetch(FetchRequest{
			Name: name,
		})
	}

	if len(coordinator) != 1 {
		return CoordinatorRef{}, fmt.Errorf("unexpected number of target coordinators for request %d", len(coordinator))
	}

	return coordinator[0], nil
}

func routingRequest(backends []CoordinatorRef, req *http.Request) routing.Request {
	coordinatorsWithStatistics := make([]routing.CoordinatorWithStatistics, len(backends))
	for i, backend := range backends {
		coordinatorsWithStatistics[i] = routing.CoordinatorWithStatistics{
			Coordinator: backend.Coordinator,
			Statistics:  backend.Statistics,
		}
	}

	return routing.Request{
		Coordinators: coordinatorsWithStatistics,
		User:         req.Header.Get(TrinoHeaderUser),
	}
}

func (p *Proxy) syncPoolState() {
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

func (p *Proxy) Close() error {
	p.termSync <- true
	return nil
}
