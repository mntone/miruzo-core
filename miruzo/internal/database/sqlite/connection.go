package sqlite

import (
	"context"
	"database/sql"
	"sync"

	sqlite3 "github.com/mattn/go-sqlite3"
)

var sqliteTimestampFormatsOnce sync.Once

func configureSQLiteTimestampFormats() {
	sqliteTimestampFormatsOnce.Do(func() {
		sqlite3.SQLiteTimestampFormats = []string{
			"2006-01-02 15:04:05.999999",
		}
	})
}

func Open(ctx context.Context, cfg *ConnectConfig) (*sql.DB, error) {
	configureSQLiteTimestampFormats()

	db, err := sql.Open("sqlite3", cfg.dsn.String())
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
	if err := verifySQLiteSupportsReturningAndStrict(timeoutContext, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	pragmas := []string{
		"PRAGMA wal_autocheckpoint=100;", // Default: 1000
	}
	for _, pragma := range pragmas {
		if _, err := db.ExecContext(timeoutContext, pragma); err != nil {
			_ = db.Close()
			return nil, err
		}
	}

	return db, nil
}
