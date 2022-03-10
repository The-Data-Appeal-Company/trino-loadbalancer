package healthcheck

import (
	"fmt"
	_ "github.com/trinodb/trino-go-client/trino"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultTimeout = 15 * time.Second
)

type HttpClusterHealth struct {
	client *http.Client
}

func NewHttpHealth() *HttpClusterHealth {
	return NewHttpHealthWithTimeout(defaultTimeout)
}

func NewHttpHealthWithTimeout(timeout time.Duration) *HttpClusterHealth {
	return &HttpClusterHealth{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   timeout,
					KeepAlive: timeout,
				}).DialContext,
				IdleConnTimeout:       timeout,
				TLSHandshakeTimeout:   timeout,
				ExpectContinueTimeout: timeout,
			},
		},
	}
}

func (p *HttpClusterHealth) Check(u *url.URL) (Health, error) {

	statusUrl := fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, "v1/status")
	req, err := http.NewRequest("GET", statusUrl, nil)
	if err != nil {
		return healthFromErr("error creating http request", err), nil
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return healthFromErr("error executing http request", err), nil
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("http request return %d status code", resp.StatusCode)
		return healthFromErr(errorMsg, fmt.Errorf(errorMsg)), nil
	}

	return Health{
		Status:    StatusHealthy,
		Message:   "all checks passed",
		Timestamp: time.Now(),
	}, nil
}
