package sqlite

import (
	"context"
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/user"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type persistenceManager struct {
	db    *sql.DB
	repos persist.Repositories
}

func newPersistenceManager(db *sql.DB) persistenceManager {
	queries := gen.New(db)
	return persistenceManager{
		db: db,
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
	db, err := database.OpenDatabase(ctx, conf)
	if err != nil {
		return persistenceManager{}, err
	}

	return newPersistenceManager(db), nil
}

func (manager persistenceManager) Close() error {
	return manager.db.Close()
}

func (manager persistenceManager) Repos() persist.Repositories {
	return manager.repos
}
