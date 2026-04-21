package mysql

import (
	"context"
	"database/sql"
)

func Open(ctx context.Context, cfg *ConnectConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.baseConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(cfg.connTune.MaxConnectionIdleTime)
	db.SetConnMaxLifetime(cfg.connTune.MaxConnectionLifeTime)
	db.SetMaxIdleConns(int(cfg.connTune.PoolWarmConnections))
	db.SetMaxOpenConns(int(cfg.connTune.MaxOpenConnections))

	timeoutContext, cancel := context.WithTimeout(ctx, cfg.connTune.ConnectionTimeout)
	defer cancel()

	if err := db.PingContext(timeoutContext); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := verifyMySQLSupportsCheck(timeoutContext, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
