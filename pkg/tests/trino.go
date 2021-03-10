package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	_ "github.com/prestodb/presto-go-client/presto"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/url"
	"time"
)

const (
	prestoImage = "trinodb/trino:353"
	prestoPort  = 8080
)

func CreatePrestoDatabase(ctx context.Context, opts ...InitOpt) (testcontainers.Container, *url.URL, error) {
	presto, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        prestoImage,
			ExposedPorts: []string{fmt.Sprintf("%d/tcp", prestoPort)},
			WaitingFor:   wait.ForLog("SERVER STARTED"),
		},
		Started: false,
	})

	if err != nil {
		return presto, nil, err
	}

	err = presto.Start(ctx)

	time.Sleep(1 * time.Second)
	if err != nil {
		return presto, nil, err
	}

	// At the end of the test remove the container
	// defer presto.Terminate(ctx)
	// Retrieve the container IP

	ip, err := presto.Host(ctx)

	if err != nil {
		return presto, nil, err
	}

	port, err := presto.MappedPort(ctx, nat.Port(fmt.Sprintf("%d", prestoPort)))
	if err != nil {
		return presto, nil, err
	}

	conn := fmt.Sprintf("http://%s:%d", ip, port.Int())

	uri, err := url.Parse(conn)
	if err != nil {
		return presto, nil, err
	}

	db, err := sql.Open("presto", conn)

	if err != nil {
		return presto, nil, err
	}

	for _, opt := range opts {
		if err := opt.Init(db); err != nil {
			return presto, uri, err
		}
	}

	return presto, uri, nil
}
