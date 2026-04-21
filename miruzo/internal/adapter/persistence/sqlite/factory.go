package sqlite

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/role"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
)

func OpenHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
	handleRole role.HandleRole,
) (sqliteHandle, error) {
	options := database.ConnectOptions{
		ConnectionTuning: persistshared.NewConnectionTuningFromConfig(appConfig),
	}
	if handleRole == role.Management {
		options.PoolWarmConnections = 1
		options.MaxOpenConnections = 1
	}

	cfg, err := database.NewConnectConfigFromDSN(appConfig.DSN, options)
	if err != nil {
		return sqliteHandle{}, err
	}

	db, err := database.Open(ctx, cfg)
	if err != nil {
		return sqliteHandle{}, err
	}

	return sqliteHandle{
		db: db,
	}, nil
}
