package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	sharedSQLite "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/user"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type sqlitePersistenceManager struct {
	db    *sql.DB
	repos persist.Repositories
}

func newRepositories(queries *gen.Queries) persist.Repositories {
	return persist.Repositories{
		Action:    action.NewRepository(queries),
		ImageList: imagelist.NewRepository(queries),
		Job:       NewJobRepository(queries),
		Settings:  NewSettingsRepository(queries),
		Stats:     stats.NewRepository(queries),
		StatsList: NewStatsListRepository(queries),
		User:      user.NewRepository(queries),
		View:      NewViewRepository(queries),
	}
}

func newPersistenceManager(db *sql.DB) sqlitePersistenceManager {
	return sqlitePersistenceManager{
		db:    db,
		repos: newRepositories(gen.New(db)),
	}
}

func (manager sqlitePersistenceManager) Repos() persist.Repositories {
	return manager.repos
}

func (manager sqlitePersistenceManager) Session(
	ctx context.Context,
	callback persist.SessionCallback,
) error {
	tx, err := manager.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf(
			"begin sqlite transaction: %w",
			sharedSQLite.MapSQLiteError("Session()", err),
		)
	}

	repos := newRepositories(gen.New(tx))
	err = callback(ctx, repos)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf(
				"rollback sqlite: %w",
				shared.JoinErrors(
					err,
					sharedSQLite.MapSQLiteError("Session()", rollbackErr),
				),
			)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf(
			"commit sqlite: %w",
			sharedSQLite.MapSQLiteError("Session()", err),
		)
	}

	return nil
}
