package trino

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const userName = "trinolb"

type ClusterApi struct {
	client         *http.Client
	trinoAuthState *trinoAuthState
}

var (
	ErrAuthFailed = errors.New("trino ui auth failed")
)

func NewClusterApi() *ClusterApi {
	return &ClusterApi{
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		trinoAuthState: &trinoAuthState{
			auth:  make(map[string]string),
			mutex: &sync.Mutex{},
		},
	}
}


func (p *ClusterApi) QueryList(coord models2.Coordinator) (QueryList, error) {
	queryStatsUrl := fmt.Sprintf("%s://%s%s", coord.URL.Scheme, coord.URL.Host, "/ui/api/query/")
	req, err := http.NewRequest("GET", queryStatsUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.authenticatedRequest(coord, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var queries QueryList
	err = json.NewDecoder(resp.Body).Decode(&queries)

	if err != nil {
		return nil, err
	}

	return queries, nil

}

func (p *ClusterApi) QueryDetail(coord models2.Coordinator, queryID string) (QueryDetail, error) {
	queryStatsUrl := fmt.Sprintf("%s://%s%s%s", coord.URL.Scheme, coord.URL.Host, "/ui/api/query/", queryID)
	req, err := http.NewRequest("GET", queryStatsUrl, nil)
	if err != nil {
		return QueryDetail{}, err
	}

	resp, err := p.authenticatedRequest(coord, req)
	if err != nil {
		return QueryDetail{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return QueryDetail{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var stats QueryDetail
	err = json.NewDecoder(resp.Body).Decode(&stats)

	if err != nil {
		return QueryDetail{}, err
	}

	return stats, nil

}

func (p *ClusterApi) ClusterStatistics(coord models2.Coordinator) (ClusterStatistics, error) {
	apiStatsUrl := fmt.Sprintf("%s://%s%s", coord.URL.Scheme, coord.URL.Host, "/ui/api/stats")
	req, err := http.NewRequest("GET", apiStatsUrl, nil)
	if err != nil {
		return ClusterStatistics{}, err
	}

	resp, err := p.authenticatedRequest(coord, req)
	if err != nil {
		return ClusterStatistics{}, err
	}

	if resp.StatusCode != 200 {
		return ClusterStatistics{}, fmt.Errorf("unexpected status code %d != 200", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ClusterStatistics{}, err
	}

	var response ClusterStatistics
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ClusterStatistics{}, err
	}

	return response, nil
}

func (p *ClusterApi) authenticatedRequest(coordinator models2.Coordinator, req *http.Request) (*http.Response, error) {
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		auth, _ := p.trinoAuthState.GetAuth(coordinator.Name)

		req.Header.Set("Cookie", auth)
		resp, err := p.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusUnauthorized {
			_, err = p.performLogin(coordinator, true)
			if err != nil {
				return nil, err
			}
			continue
		}

		return resp, nil
	}

	return nil, ErrAuthFailed
}

func (p *ClusterApi) performLogin(coord models2.Coordinator, force bool) (string, error) {
	auth, hasAuth := p.trinoAuthState.GetAuth(coord.Name)
	if !hasAuth || force {
		login, err := p.login(coord)

		if err != nil {
			return "", err
		}

		auth = login
		p.trinoAuthState.SetAuth(coord.Name, login)
	}
	return auth, nil
}

func (p *ClusterApi) login(coord models2.Coordinator) (string, error) {
	loginUrl := fmt.Sprintf("%s://%s%s", coord.URL.Scheme, coord.URL.Host, "/ui/login")
	const contentType = "application/x-www-form-urlencoded"

	body := bytes.NewBuffer([]byte(fmt.Sprintf("username=%s&password=&redirectPath=", userName)))
	resp, err := p.client.Post(loginUrl, contentType, body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusSeeOther {
		return "", fmt.Errorf("%w unexpected status code: %d for uri %s", ErrAuthFailed, resp.StatusCode, loginUrl)
	}

	cookie := resp.Header.Get("Set-Cookie")

	if cookie == "" {
		return "", fmt.Errorf("%w no Set-Cookie header present in response", ErrAuthFailed)
	}

	return cookie, nil
}

type trinoAuthState struct {
	auth  map[string]string
	mutex *sync.Mutex
}

func (p *trinoAuthState) DelAuth(id string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.auth, id)
}

func (p *trinoAuthState) SetAuth(id string, val string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.auth[id] = val
}

func (p *trinoAuthState) GetAuth(id string) (string, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	val, present := p.auth[id]
	return val, present
}
