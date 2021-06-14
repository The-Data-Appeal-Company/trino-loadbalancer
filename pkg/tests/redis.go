package tests

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateRedisServer(ctx context.Context) (testcontainers.Container, redis.UniversalClient, error) {
	const password = ""

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:6.2-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp"),
		},
		Started: false,
	})

	if err != nil {
		return redisContainer, nil, err
	}

	err = redisContainer.Start(ctx)

	if err != nil {
		return redisContainer, nil, err
	}

	ip, err := redisContainer.Host(ctx)

	if err != nil {
		return redisContainer, nil, err
	}

	port, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return redisContainer, nil, err
	}

	return redisContainer, redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", ip, port.Int()),
		Password: password,
		DB:       0, // use default DB
	}), nil
}
