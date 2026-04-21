package mysql

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/role"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/mysql"
)

func OpenHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
	handleRole role.HandleRole,
) (mysqlHandle, error) {
	options := database.ConnectOptions{
		ConnectionTuning: persistshared.NewConnectionTuningFromConfig(appConfig),
	}
	if handleRole == role.Management {
		options.MultiStatements = true // required golang-migrate
		options.PoolWarmConnections = 1
		options.MaxOpenConnections = 1
	}

	cfg, err := database.NewConnectConfigFromDSN(appConfig.DSN, options)
	if err != nil {
		return mysqlHandle{}, err
	}

	db, err := database.Open(ctx, cfg)
	if err != nil {
		return mysqlHandle{}, err
	}

	return mysqlHandle{
		db: db,
	}, nil
}
