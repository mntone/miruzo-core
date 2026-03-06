package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgreImageName = "postgres:18-alpine"
	postgreDatabase  = "miruzo"
	postgreUsername  = "miruzo"
	postgrePassword  = "miruzo1234"
)

func startPostgreContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(
		ctx,
		postgreImageName,
		postgres.WithDatabase(postgreDatabase),
		postgres.WithUsername(postgreUsername),
		postgres.WithPassword(postgrePassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		if container != nil {
			_ = container.Terminate(ctx)
		}
		return nil, fmt.Errorf("run postgre container: %w", err)
	}

	return container, nil
}
