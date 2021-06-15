package discovery

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"net/url"
)

const (
	DefaultDatabaseTableName = "trino_clusters"
)

type DatabaseStorage struct {
	db    *sql.DB
	table string
}

func NewDatabaseStorage(db *sql.DB, table string) *DatabaseStorage {
	return &DatabaseStorage{
		db:    db,
		table: table,
	}
}

func (d DatabaseStorage) Remove(ctx context.Context, name string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE name = $1", d.table)
	_, err := d.db.ExecContext(ctx, query, name)
	return err
}

func (d DatabaseStorage) Add(ctx context.Context, coordinator models2.Coordinator) error {
	query := fmt.Sprintf(`
INSERT INTO %s (name, url, tags, enabled) VALUES ($1, $2, $3, $4) 
ON CONFLICT (name) DO UPDATE 
	SET tags = excluded.tags, enabled = excluded.enabled 
`, d.table)

	tags, err := json.Marshal(coordinator.Tags)
	if err != nil {
		return fmt.Errorf("error serializing tags: %w", err)
	}

	_, err = d.db.ExecContext(ctx, query, coordinator.Name, coordinator.URL.String(), tags, coordinator.Enabled)

	return err
}

func (d DatabaseStorage) Get(ctx context.Context, name string) (models2.Coordinator, error) {
	query := fmt.Sprintf("SELECT name, url, tags, enabled FROM %s WHERE name = $1", d.table)
	rows, err := d.db.QueryContext(ctx, query, name)
	if err != nil {
		return models2.Coordinator{}, err
	}

	defer rows.Close()

	coordinators := make([]models2.Coordinator, 0)
	for rows.Next() {
		coordinator, err := coordinatorFromRow(rows)
		if err != nil {
			return models2.Coordinator{}, err
		}
		coordinators = append(coordinators, coordinator)
	}

	if len(coordinators) == 0 {
		return models2.Coordinator{}, ErrClusterNotFound
	}

	if len(coordinators) > 1 {
		return models2.Coordinator{}, fmt.Errorf("multiple clusters found with name %s", name)
	}

	return coordinators[0], nil
}

func (d DatabaseStorage) All(ctx context.Context) ([]models2.Coordinator, error) {
	query := fmt.Sprintf("SELECT name, url, tags, enabled FROM %s", d.table)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	coordinators := make([]models2.Coordinator, 0)
	for rows.Next() {
		coordinator, err := coordinatorFromRow(rows)
		if err != nil {
			return nil, err
		}
		coordinators = append(coordinators, coordinator)
	}

	return coordinators, nil
}

func coordinatorFromRow(rows *sql.Rows) (models2.Coordinator, error) {
	var name string
	var urlRaw string
	var tagsRaw string
	var enabled bool

	if err := rows.Scan(&name, &urlRaw, &tagsRaw, &enabled); err != nil {
		return models2.Coordinator{}, err
	}

	uri, err := url.Parse(urlRaw)
	if err != nil {
		return models2.Coordinator{}, err
	}

	var tags map[string]string
	if err := json.Unmarshal([]byte(tagsRaw), &tags); err != nil {
		return models2.Coordinator{}, err
	}

	var coordinator = models2.Coordinator{
		Name:    name,
		URL:     uri,
		Tags:    tags,
		Enabled: enabled,
	}

	return coordinator, nil
}
