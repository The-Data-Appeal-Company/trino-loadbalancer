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
	queryDefaultTimeout = 15 * time.Second
)

type TrinoQueryClusterHealth struct {
	client *http.Client
}

func NewTrinoQueryHealth() *TrinoQueryClusterHealth {
	return NewTrinoQueryHealthWithTimeout(queryDefaultTimeout)
}

func NewTrinoQueryHealthWithTimeout(timeout time.Duration) *TrinoQueryClusterHealth {
	return &TrinoQueryClusterHealth{
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

func (p *TrinoQueryClusterHealth) Check(u *url.URL) (Health, error) {
	if err := trino.RegisterCustomClient("hc", p.client); err != nil {
		return Health{}, err
	}

	urlWithName := fmt.Sprintf("%s://hc@%s?custom_client=hc", u.Scheme, u.Host)
	db, err := sql.Open("trino", urlWithName)
	if err != nil {
		return healthFromErr(fmt.Errorf("error opening sql connection: %w", err)), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.client.Timeout)
	defer cancel()

	row, err := db.QueryContext(ctx, "select 1")
	if err != nil {
		return healthFromErr(fmt.Errorf("error executing query: %w", err)), nil
	}

	defer row.Close()

	row.Next()
	var r int
	if err := row.Scan(&r); err != nil {
		return healthFromErr(fmt.Errorf("error reading query results: %w", err)), nil
	}

	return Health{
		Status:    StatusHealthy,
		Message:   "all checks passed",
		Timestamp: time.Now(),
	}, nil
}
