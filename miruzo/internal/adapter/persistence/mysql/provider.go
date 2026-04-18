package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/dberrors"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/mysql/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type mysqlProvider struct {
	db    *sql.DB
	repos mysqlRepositories
}

func newProvider(db *sql.DB) mysqlProvider {
	return mysqlProvider{
		db:    db,
		repos: newRepositories(gen.New(db)),
	}
}

func (prov mysqlProvider) Repos() persist.Repositories {
	return prov.repos
}

func (prov mysqlProvider) Session(
	ctx context.Context,
	callback persist.SessionCallback,
) error {
	tx, err := prov.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"begin mysql transaction: %w",
			dberrors.ToPersist("Session()", err),
		)
	}

	repos := NewSessionRepositories(gen.New(tx))
	err = callback(ctx, repos)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf(
				"rollback mysql: %w",
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
			"commit mysql: %w",
			dberrors.ToPersist("Session()", err),
		)
	}

	return nil
}
