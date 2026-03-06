package postgre

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

func OpenDatabase(
	ctx context.Context,
	conf config.DatabaseConfig,
) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(conf.DSN)
	if err != nil {
		return nil, err
	}
	cfg.MaxConnIdleTime = conf.MaxConnectionIdleTime
	cfg.MaxConnLifetime = conf.MaxConnectionLifeTime
	cfg.MaxConns = conf.MaxOpenConnections
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"

	timeoutContext, cancel := context.WithTimeout(ctx, conf.ConnectionTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(timeoutContext, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(timeoutContext); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
