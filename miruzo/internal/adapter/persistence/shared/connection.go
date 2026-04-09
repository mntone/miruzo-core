package shared

import (
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	dbshared "github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

func NewConnectionTuningFromConfig(appConfig config.DatabaseConfig) dbshared.ConnectionTuning {
	return dbshared.ConnectionTuning{
		ConnectionTimeout:     appConfig.ConnectionTimeout,
		PoolWarmConnections:   appConfig.PoolWarmConnections,
		MaxOpenConnections:    appConfig.MaxOpenConnections,
		MaxConnectionIdleTime: appConfig.MaxConnectionIdleTime,
		MaxConnectionLifeTime: appConfig.MaxConnectionLifeTime,
	}
}
