package lb

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	session2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	TrinoHeaderUser               = "X-Trino-User"
	TrinoHeaderTransaction        = "X-Trino-Transaction-Id"
	TrinoHeaderTransactionStarted = "X-Trino-Started-Transaction-Id"
)

const (
	TrinoDefaultTransactionID = "NONE"
	TrinoQueryStatusFinished  = "FINISHED"
)

type QueryClusterLinker struct {
	coordinatorName string
	storage         session2.Storage
}

func NewQueryClusterLinker(storage session2.Storage, coordinatorName string) QueryClusterLinker {
	return QueryClusterLinker{
		storage:         storage,
		coordinatorName: coordinatorName,
	}
}

// Handle Intercepts call to HttpProxy, when a response to POST /v1/statement request is detected it will create a link
// between the user/query/tx and coordinator that has provided the response to the http request.
// All the other requests are ignored, no request/response object modification should be performed.
func (q QueryClusterLinker) Handle(request *http.Request, response *http.Response) error {
	if isStatementRequest(request.URL) && response.StatusCode == http.StatusOK && request.Method == http.MethodPost {
		queryInfo, err := QueryInfoFromResponse(request, response)
		if err != nil {
			return err
		}

		return q.storage.Link(request.Context(), queryInfo, q.coordinatorName)
	}

	// todo we can simplify this nested if with a fail fast
	if isStatementRequest(request.URL) && request.Method == http.MethodGet {
		state, err := queryStateFromResponse(response)
		if err != nil {
			return err
		}

		if state.NextURI == nil {
			queryInfo, err := QueryInfoFromResponse(request, response)
			if err != nil {
				return err
			}

			if state.Stats.State == TrinoQueryStatusFinished {
				err = q.storage.Unlink(request.Context(), queryInfo)
				return err
			}
			return nil
		}
	}

	return nil
}

func QueryInfoFromResponse(req *http.Request, res *http.Response) (models2.QueryInfo, error) {
	user := req.Header.Get(TrinoHeaderUser)
	tx := req.Header.Get(TrinoHeaderTransaction)

	if len(tx) == 0 {
		tx = TrinoDefaultTransactionID
	}

	queryState, err := queryStateFromResponse(res)
	if err != nil {
		return models2.QueryInfo{}, err
	}

	return models2.QueryInfo{
		QueryID:       queryState.ID,
		User:          user,
		TransactionID: tx,
	}, nil
}

func queryStateFromResponse(res *http.Response) (models2.TrinoQueryState, error) {
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return models2.TrinoQueryState{}, err
	}

	var jsonBody = resp
	if isGzip(jsonBody) {
		reader, err := gzip.NewReader(bytes.NewBuffer(jsonBody))
		if err != nil {
			return models2.TrinoQueryState{}, err
		}

		jsonBody, err = ioutil.ReadAll(reader)
		if err != nil {
			return models2.TrinoQueryState{}, err
		}
	}

	var queryState models2.TrinoQueryState
	if err := json.Unmarshal(jsonBody, &queryState); err != nil {
		return models2.TrinoQueryState{}, err
	}

	// Reset response Body ReadCloser to be read again, this should be transparent but we may need
	// to do further checks or better way to reset a ReadCloser
	res.Body = ioutil.NopCloser(bytes.NewBuffer(resp))
	return queryState, nil
}

func queryInfoFromRequest(req *http.Request) (models2.QueryInfo, error) {
	user := req.Header.Get(TrinoHeaderUser)
	tx := req.Header.Get(TrinoHeaderTransaction)

	if len(tx) == 0 {
		tx = TrinoDefaultTransactionID
	}

	path := strings.Split(req.URL.Path, "/")
	var queryID = path[3]
	if queryID == "queued" || queryID == "executing" {
		queryID = path[4]
	}

	return models2.QueryInfo{
		QueryID:       queryID,
		User:          user,
		TransactionID: tx,
	}, nil
}

func isStatementRequest(url *url.URL) bool {
	return strings.HasPrefix(url.Path, "/v1/statement")
}

func isGzip(content []byte) bool {
	if len(content) < 2 {
		return false
	}

	return content[0] == 31 && content[1] == 139
}
