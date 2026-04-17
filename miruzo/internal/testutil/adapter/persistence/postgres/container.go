package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresImageName = "postgres:18-alpine"
	postgresDatabase  = "miruzo_test"
	postgresUsername  = "m"
	postgresPassword  = "miruzo1234"
	postgresArgs      = "--encoding=UTF8 --lc-collate=C --lc-ctype=C"
)

func startPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(
		ctx,
		postgresImageName,
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_DB":          postgresDatabase,
			"POSTGRES_USER":        postgresUsername,
			"POSTGRES_PASSWORD":    postgresPassword,
			"POSTGRES_INITDB_ARGS": postgresArgs,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		if container != nil {
			_ = container.Terminate(ctx)
		}
		return nil, fmt.Errorf("run postgres container: %w", err)
	}

	return container, nil
}
