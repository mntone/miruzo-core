package sqlite

import (
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type sqliteHandle struct {
	db *sql.DB
}

func (hdl sqliteHandle) Close() error {
	return hdl.db.Close()
}

func (hdl sqliteHandle) PersistenceManager() persist.PersistenceManager {
	return newPersistenceManager(hdl.db)
}

func (hdl sqliteHandle) MigrationRunner() persist.MigrationRunner {
	return NewMigrationRunnerFromDB(hdl.db)
}
