package testutil

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	db "github.com/mntone/miruzo-core/miruzo/internal/database/postgre"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func openTestPoolFromContainer(
	ctx context.Context,
	container *postgres.PostgresContainer,
) (*pgxpool.Pool, error) {
	conf := config.DefaultDatabaseConfig()
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("get postgre container dsn: %w", err)
	}
	conf.DSN = dsn
	conf.MaxOpenConnections = 1

	pool, err := db.OpenDatabase(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("open postgre database: %w", err)
	}

	return pool, nil
}
