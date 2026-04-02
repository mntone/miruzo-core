package persistence

import (
	"context"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func NewMigrationRunner(
	ctx context.Context,
	conf config.DatabaseConfig,
) (persist.MigrationRunner, error) {
	switch conf.Backend {
	case config.DatabaseBackendPostgres:
		return postgres.NewMigrationRunner(ctx, conf)
	case config.DatabaseBackendSQLite:
		return sqlite.NewMigrationRunner(ctx, conf)
	default:
		return nil, fmt.Errorf("unsupported database backend: %s", conf.Backend)
	}
}

func NewPersistenceManager(
	ctx context.Context,
	conf config.DatabaseConfig,
) (persist.PersistenceManager, error) {
	switch conf.Backend {
	case config.DatabaseBackendPostgres:
		return postgres.NewPersistenceManager(ctx, conf)
	case config.DatabaseBackendSQLite:
		return sqlite.NewPersistenceManager(ctx, conf)
	default:
		return nil, fmt.Errorf("unsupported database backend: %s", conf.Backend)
	}
}
