package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type sqliteProvider struct {
	db    *sql.DB
	repos sqliteRepositories
}

func newProvider(db *sql.DB) sqliteProvider {
	return sqliteProvider{
		db:    db,
		repos: newRepositories(gen.New(db)),
	}
}

func (prov sqliteProvider) Repos() persist.Repositories {
	return prov.repos
}

func (prov sqliteProvider) Session(
	ctx context.Context,
	callback persist.SessionCallback,
) error {
	tx, err := prov.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"begin sqlite transaction: %w",
			dberrors.ToPersist("Session()", err),
		)
	}

	repos := NewSessionRepositories(gen.New(tx))
	err = callback(ctx, repos)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf(
				"rollback sqlite: %w",
				persistshared.JoinErrors(
					err,
					dberrors.ToPersist("Session()", rollbackErr),
				),
			)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf(
			"commit sqlite: %w",
			dberrors.ToPersist("Session()", err),
		)
	}

	return nil
}
