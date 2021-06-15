package ui

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/gorilla/mux"
	"net/http"
)

type Api struct {
	statsRetriever   trino.Api
	discoveryStorage discovery.Storage
	discover         discovery.Discovery
	logger           logging.Logger
}

func NewApi(statsRetriever trino.Api, discover discovery.Discovery, discoverStorage discovery.Storage, logger logging.Logger) Api {
	return Api{
		statsRetriever:   statsRetriever,
		discoveryStorage: discoverStorage,
		discover:         discover,
		logger:           logger,
	}
}

func (a *Api) Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/health", healthProbe)
	r.HandleFunc("/api/stats", a.statistics)
	r.Methods(http.MethodGet).Path("/api/clusters").HandlerFunc(a.clustersList)
	r.Methods(http.MethodPatch).Path("/api/cluster/{name}").HandlerFunc(a.updateCluster)
	r.Methods(http.MethodPost).Path("/api/cluster").HandlerFunc(a.addCluster)
	r.Methods(http.MethodPost).Path("/api/cluster/discover").HandlerFunc(a.launchDiscover)

	return r
}

func (a *Api) Serve(addr string) error {
	return http.ListenAndServe(addr, a.Router())
}
