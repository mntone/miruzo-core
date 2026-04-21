package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, cfg *ConnectConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, cfg.baseConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
