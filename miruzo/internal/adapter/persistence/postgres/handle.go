package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type postgresHandle struct {
	pool *pgxpool.Pool
}

func (hdl postgresHandle) Close() error {
	hdl.pool.Close()
	return nil
}

func (hdl postgresHandle) PersistenceManager() persist.PersistenceManager {
	return newPersistenceManager(hdl.pool)
}

func (hdl postgresHandle) MigrationRunner() persist.MigrationRunner {
	return NewMigrationRunnerFromPool(hdl.pool)
}
