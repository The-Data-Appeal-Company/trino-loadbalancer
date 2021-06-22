package healthcheck

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/trinodb/trino-go-client/trino"
	_ "github.com/trinodb/trino-go-client/trino"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultTimeout = 15 * time.Second
)

type ClusterHealth struct {
	client *http.Client
}

func NewHttpHealth() *ClusterHealth {
	return NewHttpHealthWithTimeout(defaultTimeout)
}

func NewHttpHealthWithTimeout(timeout time.Duration) *ClusterHealth {
	return &ClusterHealth{
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

func (p *ClusterHealth) Check(u *url.URL) (Health, error) {
	if err := trino.RegisterCustomClient("hc", p.client); err != nil {
		return Health{}, err
	}

	urlWithName := fmt.Sprintf("%s://hc@%s?custom_client=hc", u.Scheme, u.Host)
	db, err := sql.Open("trino", urlWithName)
	if err != nil {
		return healthFromErr("error opening sql connection", err), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.client.Timeout)
	defer cancel()

	row, err := db.QueryContext(ctx, "select 1")
	if err != nil {
		return healthFromErr("error executing query", err), nil
	}

	defer row.Close()

	row.Next()
	var r int
	if err := row.Scan(&r); err != nil {
		return healthFromErr("error reading query results", err), nil
	}

	return Health{
		Status:    StatusHealthy,
		Message:   "all checks passed",
		Timestamp: time.Now(),
	}, nil
}

func healthFromErr(message string, err error) Health {
	return Health{
		Status:    StatusUnhealthy,
		Message:   fmt.Sprintf("%s: %s", message, err.Error()),
		Error:     err,
		Timestamp: time.Now(),
	}
}
