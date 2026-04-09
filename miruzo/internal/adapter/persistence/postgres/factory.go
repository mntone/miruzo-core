package postgres

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/role"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/postgres"
)

func buildConnectConfig(appConfig config.DatabaseConfig) database.ConnectConfig {
	return database.ConnectConfig{
		DSN:              appConfig.DSN,
		ConnectionTuning: shared.NewConnectionTuningFromConfig(appConfig),
	}
}

func OpenHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
	handleRole role.HandleRole,
) (postgresHandle, error) {
	cfg := buildConnectConfig(appConfig)
	if handleRole == role.Management {
		cfg.PoolWarmConnections = 1
		cfg.MaxOpenConnections = 1
	}

	pool, err := database.Open(ctx, cfg)
	if err != nil {
		return postgresHandle{}, err
	}

	return postgresHandle{
		pool: pool,
	}, nil
}
