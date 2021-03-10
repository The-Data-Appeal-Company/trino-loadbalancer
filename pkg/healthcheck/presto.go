package healthcheck

import (
	"database/sql"
	"fmt"
	"github.com/prestodb/presto-go-client/presto"
	_ "github.com/prestodb/presto-go-client/presto"
	"net"
	"net/http"
	"net/url"
	"time"
)

type PrestoClusterHealth struct {
	client *http.Client
}

func NewPrestoHealth() PrestoClusterHealth {
	return PrestoClusterHealth{
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

func (p PrestoClusterHealth) Check(u *url.URL) (Health, error) {
	if err := presto.RegisterCustomClient("hc", p.client); err != nil {
		return Health{}, err
	}

	urlWithName := fmt.Sprintf("%s://hc@%s?custom_client=hc", u.Scheme, u.Host)
	db, err := sql.Open("presto", urlWithName)
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
