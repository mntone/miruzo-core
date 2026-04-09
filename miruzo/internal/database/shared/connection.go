package shared

import "time"

type ConnectionTuning struct {
	ConnectionTimeout     time.Duration
	PoolWarmConnections   int32
	MaxOpenConnections    int32
	MaxConnectionIdleTime time.Duration
	MaxConnectionLifeTime time.Duration
}
