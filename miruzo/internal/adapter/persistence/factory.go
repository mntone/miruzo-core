package persistence

import (
	"context"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func NewPersistenceManager(
	ctx context.Context,
	conf config.DatabaseConfig,
) (persist.PersistenceManager, error) {
	switch conf.Backend {
	case config.DatabaseBackendPostgre:
		return postgre.NewPersistenceManager(ctx, conf)
	case config.DatabaseBackendSQLite:
		return sqlite.NewPersistenceManager(ctx, conf)
	default:
		return nil, fmt.Errorf("unsupported database backend: %s", conf.Backend)
	}
}
