package persistence

import (
	"context"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/role"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

func OpenAppHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
) (DatabaseAppHandle, error) {
	switch appConfig.Backend {
	case backend.MySQL:
		return mysql.OpenHandle(ctx, appConfig, role.App)
	case backend.PostgreSQL:
		return postgres.OpenHandle(ctx, appConfig, role.App)
	case backend.SQLite:
		return sqlite.OpenHandle(ctx, appConfig, role.App)
	default:
		return nil, fmt.Errorf("unsupported database backend: %s", appConfig.Backend)
	}
}

func OpenManagementHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
) (DatabaseManagementHandle, error) {
	switch appConfig.Backend {
	case backend.MySQL:
		return mysql.OpenHandle(ctx, appConfig, role.Management)
	case backend.PostgreSQL:
		return postgres.OpenHandle(ctx, appConfig, role.Management)
	case backend.SQLite:
		return sqlite.OpenHandle(ctx, appConfig, role.Management)
	default:
		return nil, fmt.Errorf("unsupported database backend: %s", appConfig.Backend)
	}
}
