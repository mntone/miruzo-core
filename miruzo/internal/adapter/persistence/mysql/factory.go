package mysql

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/role"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/mysql"
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
) (mysqlHandle, error) {
	cfg := buildConnectConfig(appConfig)
	if handleRole == role.Management {
		cfg.MultiStatements = true // required golang-migrate
		cfg.PoolWarmConnections = 1
		cfg.MaxOpenConnections = 1
	}

	db, err := database.Open(ctx, cfg)
	if err != nil {
		return mysqlHandle{}, err
	}

	return mysqlHandle{
		db: db,
	}, nil
}
