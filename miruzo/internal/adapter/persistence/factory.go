package persistence

import (
	"context"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func NewRepository(
	ctx context.Context,
	conf config.DatabaseConfig,
) (persist.RepositoryFactory, error) {
	switch conf.Backend {
	case config.DatabaseBackendPostgre:
		return postgre.NewRepositoryFactory(ctx, conf)
	case config.DatabaseBackendSQLite:
		return sqlite.NewRepositoryFactory(ctx, conf)
	default:
		return nil, fmt.Errorf("unsupported database backend: %s", conf.Backend)
	}
}
