package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

func OpenDatabase(
	ctx context.Context,
	conf config.DatabaseConfig,
) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", conf.DSN)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(conf.MaxConnectionIdleTime)
	db.SetConnMaxLifetime(conf.MaxConnectionLifeTime)
	db.SetMaxOpenConns(int(conf.MaxOpenConnections))

	timeoutContext, cancel := context.WithTimeout(ctx, conf.ConnectionTimeout)
	defer cancel()

	if err := db.PingContext(timeoutContext); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := verifySQLiteSupportsReturning(timeoutContext, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA wal_autocheckpoint=1000;",
	}
	for _, pragma := range pragmas {
		if _, err := db.ExecContext(timeoutContext, pragma); err != nil {
			_ = db.Close()
			return nil, err
		}
	}

	return db, nil
}
