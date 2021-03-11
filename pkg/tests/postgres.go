package tests

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"net/url"
	"time"
)

func CreatePostgresDatabase(ctx context.Context, opts ...InitOpt) (testcontainers.Container, *sql.DB, error) {
	const password = "password"
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:10",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForLog("server stopped"),
			Env: map[string]string{
				"POSTGRES_PASSWORD": password,
			},
		},
		Started: false,
	})

	if err != nil {
		return pg, nil, err
	}

	err = pg.Start(ctx)

	time.Sleep(1 * time.Second)
	if err != nil {
		return pg, nil, err
	}

	// At the end of the test remove the container
	// defer pg.Terminate(ctx)
	// Retrieve the container IP

	ip, err := pg.Host(ctx)

	if err != nil {
		return pg, nil, err
	}

	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		return pg, nil, err
	}

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", "postgres", url.QueryEscape(password), ip, port.Int(), "postgres", "disable")

	db, err := sql.Open("postgres", conn)

	if err != nil {
		return pg, nil, err
	}

	for _, opt := range opts {
		if err := opt.Init(db); err != nil {
			return pg, db, err
		}
	}

	return pg, db, nil
}

type InitOpt interface {
	Init(db *sql.DB) error
}

type InitScript struct {
	script string
}

func WithInitScript(filePath string) InitScript {
	return InitScript{script: filePath}
}

func (i InitScript) Init(db *sql.DB) error {

	script, err := ioutil.ReadFile(i.script)

	if err != nil {
		return err
	}

	if _, err = db.Exec(string(script)); err != nil {
		return err
	}

	return nil
}
