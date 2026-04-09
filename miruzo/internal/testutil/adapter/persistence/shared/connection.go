package shared

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

func NewTestConnectionTuning() shared.ConnectionTuning {
	return shared.ConnectionTuning{
		ConnectionTimeout:     20 * time.Second,
		PoolWarmConnections:   1,
		MaxOpenConnections:    1,
		MaxConnectionIdleTime: 0,
		MaxConnectionLifeTime: 1 * time.Minute,
	}
}
