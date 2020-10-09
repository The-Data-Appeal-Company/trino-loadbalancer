package healthcheck

import (
	"database/sql"
	"fmt"
	_ "github.com/prestodb/presto-go-client/presto"
	"net/url"
	"time"
)

type PrestoClusterHealth struct {
}

func NewPrestoHealth() PrestoClusterHealth {
	return PrestoClusterHealth{}
}

func (p PrestoClusterHealth) Check(u *url.URL) (Health, error) {
	urlWithName := fmt.Sprintf("%s://hc@%s", u.Scheme, u.Host)
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
