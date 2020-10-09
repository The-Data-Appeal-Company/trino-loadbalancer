package statistics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const userName = "prestolb"

type PrestoClusterApi struct {
	client         *http.Client
	prestoSqlState *prestoSQLState
}

func NewPrestoClusterApi() *PrestoClusterApi {
	return &PrestoClusterApi{
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		prestoSqlState: &prestoSQLState{
			auth:  make(map[string]string),
			mutex: &sync.Mutex{},
		},
	}
}

func (p *PrestoClusterApi) GetStatistics(coord models.Coordinator) (models.ClusterStatistics, error) {
	switch coord.Distribution {
	case models.PrestoDistDb:
		return p.getPrestoDBStatistics(coord)
	case models.PrestoDistSql:
		return p.getPrestoSQLStatistics(coord)
	default:
		return models.ClusterStatistics{}, fmt.Errorf("statistics not implemented for distribution: %s", string(coord.Distribution))
	}
}

func (p *PrestoClusterApi) getPrestoSQLStatistics(coord models.Coordinator) (models.ClusterStatistics, error) {

	auth, hasAuth := p.prestoSqlState.GetAuth(coord.Name)

	if !hasAuth {
		login, err := p.login(coord)

		if err != nil {
			return models.ClusterStatistics{}, err
		}

		auth = login
		p.prestoSqlState.SetAuth(coord.Name, login)
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
		p.prestoSqlState.DelAuth(coord.Name)
		return models.ClusterStatistics{}, errors.New("deauth by presto ui, trying again next update")
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

func (p *PrestoClusterApi) login(coord models.Coordinator) (string, error) {
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

func (p *PrestoClusterApi) getPrestoDBStatistics(coord models.Coordinator) (models.ClusterStatistics, error) {
	u := fmt.Sprintf("%s://%s%s", coord.URL.Scheme, coord.URL.Host, "/v1/cluster")
	resp, err := p.client.Get(u)

	if err != nil {
		return models.ClusterStatistics{}, err
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

type prestoSQLState struct {
	auth  map[string]string
	mutex *sync.Mutex
}

func (p *prestoSQLState) DelAuth(id string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.auth, id)
}

func (p *prestoSQLState) SetAuth(id string, val string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.auth[id] = val
}

func (p *prestoSQLState) GetAuth(id string) (string, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	val, present := p.auth[id]
	return val, present
}
