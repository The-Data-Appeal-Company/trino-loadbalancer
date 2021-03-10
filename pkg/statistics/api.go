package statistics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const userName = "trinolb"

type TrinoClusterApi struct {
	client         *http.Client
	trinoAuthState *trinoAuthState
}

func NewPrestoClusterApi() *TrinoClusterApi {
	return &TrinoClusterApi{
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

func (p *TrinoClusterApi) GetStatistics(coord models.Coordinator) (models.ClusterStatistics, error) {

	auth, hasAuth := p.trinoAuthState.GetAuth(coord.Name)

	if !hasAuth {
		login, err := p.login(coord)

		if err != nil {
			return models.ClusterStatistics{}, err
		}

		auth = login
		p.trinoAuthState.SetAuth(coord.Name, login)
	}

	apiStatsUrl := fmt.Sprintf("%s://%s%s", coord.URL.Scheme, coord.URL.Host, "/ui/api/stats")
	req, err := http.NewRequest("GET", apiStatsUrl, nil)
	if err != nil {
		return models.ClusterStatistics{}, err
	}

	req.Header.Set("Cookie", auth)

	resp, err := p.client.Do(req)
	if err != nil {
		return models.ClusterStatistics{}, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		p.trinoAuthState.DelAuth(coord.Name)
		return models.ClusterStatistics{}, errors.New("deauthenticated from trino ui, trying again next update")
	}

	if resp.StatusCode != 200 {
		return models.ClusterStatistics{}, fmt.Errorf("unexpected status code %d != 200", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.ClusterStatistics{}, err
	}

	var response models.ClusterStatistics
	err = json.Unmarshal(body, &response)
	if err != nil {
		return models.ClusterStatistics{}, err
	}

	return response, nil
}

func (p *TrinoClusterApi) login(coord models.Coordinator) (string, error) {
	loginUrl := fmt.Sprintf("%s://%s%s", coord.URL.Scheme, coord.URL.Host, "/ui/login")
	const contentType = "application/x-www-form-urlencoded"

	body := bytes.NewBuffer([]byte(fmt.Sprintf("username=%s&password=&redirectPath=", userName)))
	resp, err := p.client.Post(loginUrl, contentType, body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusSeeOther {
		return "", fmt.Errorf("unexpected status code: %d for uri %s", resp.StatusCode, loginUrl)
	}

	cookie := resp.Header.Get("Set-Cookie")

	if cookie == "" {
		return "", errors.New("no Set-Cookie header present in response")
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
