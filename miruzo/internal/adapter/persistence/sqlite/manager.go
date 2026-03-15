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
			Settings:  NewSettingsRepository(queries),
			Stats:     stats.NewRepository(queries),
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

func (manager persistenceManager) Session(
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

	queries := gen.New(tx)
	repos := persist.Repositories{
		Action:    action.NewRepository(queries),
		ImageList: imagelist.NewRepository(queries),
		Settings:  NewSettingsRepository(queries),
		Stats:     stats.NewRepository(queries),
		User:      user.NewRepository(queries),
		View:      NewViewRepository(queries),
	}

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
