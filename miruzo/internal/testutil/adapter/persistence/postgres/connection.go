package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/shared"
)

func openTestPoolFromDSN(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg := database.ConnectConfig{
		DSN:              dsn,
		ConnectionTuning: shared.NewTestConnectionTuning(),
	}
	pool, err := database.Open(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open postgres database: %w", err)
	}

	return pool, nil
}
