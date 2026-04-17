package sqlite

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/role"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
)

func buildConnectConfig(appConfig config.DatabaseConfig) database.ConnectConfig {
	return database.ConnectConfig{
		DSN:              appConfig.DSN,
		ConnectionTuning: persistshared.NewConnectionTuningFromConfig(appConfig),
	}
}

func OpenHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
	handleRole role.HandleRole,
) (sqliteHandle, error) {
	cfg := buildConnectConfig(appConfig)
	if handleRole == role.Management {
		cfg.PoolWarmConnections = 1
		cfg.MaxOpenConnections = 1
	}

	db, err := database.Open(ctx, cfg)
	if err != nil {
		return sqliteHandle{}, err
	}

	return sqliteHandle{
		db: db,
	}, nil
}
