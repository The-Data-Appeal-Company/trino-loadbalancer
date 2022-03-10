package trino

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const userName = "trinolb"

type ClusterApi struct {
	client         *http.Client
	trinoAuthState *trinoAuthState
}

var (
	ErrAuthFailed    = errors.New("trino ui auth failed")
	ErrQueryNotFound = errors.New("query not found")
)

func NewClusterApi() *ClusterApi {
	return &ClusterApi{
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   15 * time.Second,
				IdleConnTimeout:       15 * time.Second,
				ResponseHeaderTimeout: 15 * time.Second,
				ExpectContinueTimeout: 15 * time.Second,
			},
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

func (p *ClusterApi) QueryList(url *url.URL) (QueryList, error) {
	queryStatsUrl := fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, "/ui/api/query/")
	req, err := http.NewRequest("GET", queryStatsUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.authenticatedRequest(url, req)
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

func (p *ClusterApi) QueryDetail(coord *url.URL, queryID string) (QueryDetail, error) {
	queryStatsUrl := fmt.Sprintf("%s://%s%s%s", coord.Scheme, coord.Host, "/ui/api/query/", queryID)
	req, err := http.NewRequest("GET", queryStatsUrl, nil)
	if err != nil {
		return QueryDetail{}, err
	}

	resp, err := p.authenticatedRequest(coord, req)
	if err != nil {
		return QueryDetail{}, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusGone {
			return QueryDetail{}, ErrQueryNotFound
		}
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

func (p *ClusterApi) ClusterStatistics(coord *url.URL) (ClusterStatistics, error) {
	apiStatsUrl := fmt.Sprintf("%s://%s%s", coord.Scheme, coord.Host, "/ui/api/stats")
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

func (p *ClusterApi) authenticatedRequest(coordinatorUrl *url.URL, req *http.Request) (*http.Response, error) {
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		auth, _ := p.trinoAuthState.GetAuth(coordinatorUrl.String())

		req.Header.Set("Cookie", auth)
		resp, err := p.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusUnauthorized {
			_, err = p.performLogin(coordinatorUrl, true)
			if err != nil {
				return nil, err
			}
			continue
		}

		return resp, nil
	}

	return nil, ErrAuthFailed
}

func (p *ClusterApi) performLogin(coordinatoUrl *url.URL, force bool) (string, error) {
	coordinatorID := coordinatoUrl.String()
	auth, hasAuth := p.trinoAuthState.GetAuth(coordinatorID)
	if !hasAuth || force {
		login, err := p.login(coordinatoUrl)

		if err != nil {
			return "", err
		}

		auth = login
		p.trinoAuthState.SetAuth(coordinatorID, login)
	}
	return auth, nil
}

func (p *ClusterApi) login(coord *url.URL) (string, error) {
	loginUrl := fmt.Sprintf("%s://%s%s", coord.Scheme, coord.Host, "/ui/login")
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
