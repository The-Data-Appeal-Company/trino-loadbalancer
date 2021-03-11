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
	trinoImage = "trinodb/trino:353"
	trinoPort  = 8080
)

func CreateTrinoCluster(ctx context.Context, opts ...InitOpt) (testcontainers.Container, *url.URL, error) {
	trino, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        trinoImage,
			ExposedPorts: []string{fmt.Sprintf("%d/tcp", trinoPort)},
			WaitingFor:   wait.ForLog("SERVER STARTED"),
		},
		Started: false,
	})

	if err != nil {
		return trino, nil, err
	}

	err = trino.Start(ctx)

	time.Sleep(1 * time.Second)
	if err != nil {
		return trino, nil, err
	}

	// At the end of the test remove the container
	// defer trino.Terminate(ctx)
	// Retrieve the container IP

	ip, err := trino.Host(ctx)

	if err != nil {
		return trino, nil, err
	}

	port, err := trino.MappedPort(ctx, nat.Port(fmt.Sprintf("%d", trinoPort)))
	if err != nil {
		return trino, nil, err
	}

	conn := fmt.Sprintf("http://%s:%d", ip, port.Int())

	uri, err := url.Parse(conn)
	if err != nil {
		return trino, nil, err
	}

	db, err := sql.Open("trino", conn)

	if err != nil {
		return trino, nil, err
	}

	for _, opt := range opts {
		if err := opt.Init(db); err != nil {
			return trino, uri, err
		}
	}

	return trino, uri, nil
}
