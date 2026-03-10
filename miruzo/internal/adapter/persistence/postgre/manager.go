package postgre

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/user"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/postgre"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type persistenceManager struct {
	pool  *pgxpool.Pool
	repos persist.Repositories
}

func newPersistenceManager(pool *pgxpool.Pool) persistenceManager {
	queries := gen.New(pool)
	return persistenceManager{
		pool: pool,
		repos: persist.Repositories{
			Action:    action.NewRepository(queries),
			ImageList: imagelist.NewRepository(queries),
			User:      user.NewRepository(queries),
			View:      NewViewRepository(queries),
		},
	}
}

func NewPersistenceManager(
	ctx context.Context,
	conf config.DatabaseConfig,
) (persistenceManager, error) {
	pool, err := database.OpenDatabase(ctx, conf)
	if err != nil {
		return persistenceManager{}, err
	}

	return newPersistenceManager(pool), nil
}

func (manager persistenceManager) Close() error {
	manager.pool.Close()
	return nil
}

func (manager persistenceManager) Repos() persist.Repositories {
	return manager.repos
}
