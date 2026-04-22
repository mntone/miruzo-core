package persistence

import "github.com/mntone/miruzo-core/miruzo/internal/persist"

type DatabaseCloser interface {
	Close() error
}

type DatabaseAdminHandle interface {
	DatabaseCloser
	DatabaseAdmin
}

type DatabaseAppHandle interface {
	DatabaseCloser
	PersistenceProvider() persist.PersistenceProvider
}

type DatabaseManagementHandle interface {
	DatabaseAppHandle
	MigrationRunner() persist.MigrationRunner
}
