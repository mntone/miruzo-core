package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type postgresProvider struct {
	pool  *pgxpool.Pool
	repos postgresRepositories
}

func newProvider(pool *pgxpool.Pool) postgresProvider {
	return postgresProvider{
		pool:  pool,
		repos: newRepositories(gen.New(pool)),
	}
}

func (prov postgresProvider) Repos() persist.Repositories {
	return prov.repos
}

func (prov postgresProvider) Session(
	ctx context.Context,
	callback persist.SessionCallback,
) error {
	tx, err := prov.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf(
			"begin postgres transaction: %w",
			dberrors.ToPersist("Session()", err),
		)
	}

	repos := NewSessionRepositories(gen.New(tx))
	err = callback(ctx, repos)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf(
				"rollback postgres: %w",
				shared.JoinErrors(
					err,
					dberrors.ToPersist("Session()", rollbackErr),
				),
			)
		}

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf(
			"commit postgres: %w",
			dberrors.ToPersist("Session()", err),
		)
	}

	return nil
}
