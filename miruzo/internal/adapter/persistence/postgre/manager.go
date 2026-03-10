package postgre

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/imagelist"
	sharedPostgre "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/user"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
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

func (manager persistenceManager) Session(
	ctx context.Context,
	callback persist.SessionCallback,
) error {
	tx, err := manager.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	})
	if err != nil {
		return fmt.Errorf(
			"begin postgre transaction: %w",
			sharedPostgre.MapPostgreError("Session()", err),
		)
	}

	queries := gen.New(tx)
	repos := persist.Repositories{
		Action:    action.NewRepository(queries),
		ImageList: imagelist.NewRepository(queries),
		Stats:     stats.NewRepository(queries),
		User:      user.NewRepository(queries),
		View:      NewViewRepository(queries),
	}

	err = callback(ctx, repos)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return fmt.Errorf(
			"rollback postgre: %w",
			shared.JoinErrors(
				sharedPostgre.MapPostgreError("Session()", err),
				sharedPostgre.MapPostgreError("Session()", rollbackErr),
			),
		)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf(
			"commit postgre: %w",
			sharedPostgre.MapPostgreError("Session()", err),
		)
	}

	return nil
}
