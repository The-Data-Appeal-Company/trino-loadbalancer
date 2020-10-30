package api

import (
	"encoding/json"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Cluster struct {
	Name      string            `json:"name"`
	Host      string            `json:"host"`
	Available bool              `json:"available"`
	Enabled   bool              `json:"enabled"`
	Tags      map[string]string `json:"tags"`
}

type ClusterUpdateRequest struct {
	Enabled bool `json:"enabled"`
}

type ClustersResponse struct {
	Clusters []Cluster `json:"clusters"`
}

func (a Api) updateCluster(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req ClusterUpdateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cluster, err := a.discoveryStorage.Get(ctx, vars["name"])

	if err == discovery.ErrClusterNotFound {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cluster.Enabled = req.Enabled
	err = a.discoveryStorage.Add(ctx, cluster)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a Api) clustersList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	clusters, err := a.discoveryStorage.All(ctx)
	if err != nil {
		writer.Write([]byte(err.Error()))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	results := make([]Cluster, len(clusters))
	for i, c := range clusters {
		results[i] = Cluster{
			Name:      c.Name,
			Host:      c.URL.String(),
			Enabled:   c.Enabled,
			Tags:      c.Tags,
			Available: true, // TODO add health check
		}

	}

	body, err := json.Marshal(results)
	if err != nil {
		writer.Write([]byte(err.Error()))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(body)
}

type ClusterAddRequest struct {
	Name         string `json:"name"`
	Url          string `json:"url"`
	Distribution string `json:"distribution"`
	Enabled      bool   `json:"enabled"`
}

func (a Api) addCluster(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req ClusterAddRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedUrl, err := url.Parse(req.Url)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dist, err := a.distFromRequest(req)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.discoveryStorage.Add(ctx, models.Coordinator{
		Name:         req.Name,
		URL:          parsedUrl,
		Enabled:      req.Enabled,
		Distribution: dist,
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a Api) launchDiscover(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clusters, err := a.discover.Discover(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, cluster := range clusters {
		if err := a.discoveryStorage.Add(ctx, cluster); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (a Api) distFromRequest(req ClusterAddRequest) (models.PrestoDist, error) {
	switch req.Distribution {
	case string(models.PrestoDistSql):
		return models.PrestoDistSql, nil
	case string(models.PrestoDistDb):
		return models.PrestoDistDb, nil
	default:
		return "", fmt.Errorf("unknown distribution %s", req.Distribution)
	}
}
