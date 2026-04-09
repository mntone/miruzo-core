package sqlite

import (
	"context"
	"database/sql"
	"net/url"
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

func buildSQLiteDSN(dsn string) (string, error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	queries := parsed.Query()
	if !queries.Has("_txlock") {
		queries.Set("_txlock", "immediate")
	}
	if !queries.Has("_foreign_keys") {
		queries.Set("_foreign_keys", "1") // 1=ON
	}
	if !queries.Has("_journal_mode") && !queries.Has("_synchronous") {
		queries.Set("_journal_mode", "WAL")
	}

	parsed.RawQuery = queries.Encode()
	return parsed.String(), nil
}

func Open(ctx context.Context, cfg ConnectConfig) (*sql.DB, error) {
	configureSQLiteTimestampFormats()

	dsn, err := buildSQLiteDSN(cfg.DSN)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(cfg.MaxConnectionIdleTime)
	db.SetConnMaxLifetime(cfg.MaxConnectionLifeTime)
	db.SetMaxIdleConns(int(cfg.PoolWarmConnections))
	db.SetMaxOpenConns(int(cfg.MaxOpenConnections))

	timeoutContext, cancel := context.WithTimeout(ctx, cfg.ConnectionTimeout)
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
