package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, cfg ConnectConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnectionTimeout
	poolConfig.MaxConnIdleTime = cfg.MaxConnectionIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnectionLifeTime
	poolConfig.MinConns = cfg.PoolWarmConnections
	poolConfig.MaxConns = cfg.MaxOpenConnections
	poolConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
