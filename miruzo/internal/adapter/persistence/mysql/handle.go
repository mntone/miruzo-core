package mysql

import (
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type mysqlHandle struct {
	db *sql.DB
}

func (hdl mysqlHandle) Close() error {
	return hdl.db.Close()
}

func (hdl mysqlHandle) PersistenceProvider() persist.PersistenceProvider {
	return newProvider(hdl.db)
}

func (hdl mysqlHandle) MigrationRunner() persist.MigrationRunner {
	return NewMigrationRunnerFromDB(hdl.db)
}
