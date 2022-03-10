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
	req, err := http.NewRequest(http.MethodGet, statusUrl, nil)
	if err != nil {
		return healthFromErr(fmt.Errorf("error creating http request: %w", err)), nil
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return healthFromErr(fmt.Errorf("error executing http request: %w", err)), nil
	}

	if resp.StatusCode != http.StatusOK {
		return healthFromErr(fmt.Errorf("http request return %d status code", resp.StatusCode)), nil
	}

	return Health{
		Status:    StatusHealthy,
		Message:   "all checks passed",
		Timestamp: time.Now(),
	}, nil
}
