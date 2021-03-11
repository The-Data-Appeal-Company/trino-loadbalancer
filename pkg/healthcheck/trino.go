package healthcheck

import (
	"database/sql"
	"fmt"
	"github.com/trinodb/trino-go-client/trino"
	_ "github.com/trinodb/trino-go-client/trino"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ClusterHealth struct {
	client *http.Client
}

func NewHttpHealth() ClusterHealth {
	return ClusterHealth{
		client: &http.Client{
			Transport: &http.Transport{
				IdleConnTimeout:       15 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DialContext: (&net.Dialer{
					Timeout:   15 * time.Second,
					KeepAlive: 15 * time.Second,
				}).DialContext,
			},
			Timeout: 10 * time.Second,
		},
	}
}

func (p ClusterHealth) Check(u *url.URL) (Health, error) {
	if err := trino.RegisterCustomClient("hc", p.client); err != nil {
		return Health{}, err
	}

	urlWithName := fmt.Sprintf(	"%s://hc@%s?custom_client=hc", u.Scheme, u.Host)
	db, err := sql.Open("trino", urlWithName)
	if err != nil {
		return Health{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("error connecting to cluster: %s", err.Error()),
			Timestamp: time.Now(),
		}, nil
	}

	row, err := db.Query("select 1")
	if err != nil {
		return Health{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("error executing query: %s", err.Error()),
			Timestamp: time.Now(),
		}, nil
	}

	defer row.Close()

	row.Next()

	var r int
	if err := row.Scan(&r); err != nil {
		return Health{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("error reading query results: %s", err.Error()),
			Timestamp: time.Now(),
		}, nil
	}

	return Health{
		Status:    StatusHealthy,
		Message:   "all checks passed",
		Timestamp: time.Now(),
	}, nil
}
