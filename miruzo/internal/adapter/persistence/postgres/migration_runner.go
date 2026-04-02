package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/postgres"
	migrations "github.com/mntone/miruzo-core/miruzo/internal/database/postgres/migrations_min"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func NewMigrationRunnerFromPool(pool *pgxpool.Pool) shared.MigrationRunner {
	return shared.MigrationRunner{
		Spec: migrations.NewSpec(pool),
	}
}

func NewMigrationRunner(
	ctx context.Context,
	conf config.DatabaseConfig,
) (persist.MigrationRunner, error) {
	pool, err := database.OpenDatabase(ctx, conf)
	if err != nil {
		return nil, err
	}

	return shared.MigrationRunner{
		Spec: migrations.NewSpec(pool),
		CloseDatabase: func() error {
			pool.Close()
			return nil
		},
	}, nil
}
