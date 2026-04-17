package mysql

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

const (
	mysqlSQLMode  = "'TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY'"
	mysqlTimezone = "'+00:00'"
)

func Open(ctx context.Context, cfg ConnectConfig) (*sql.DB, error) {
	config, err := mysql.ParseDSN(cfg.DSN)
	if err != nil {
		return nil, err
	}
	config.MultiStatements = cfg.MultiStatements
	if config.Params == nil {
		config.Params = make(map[string]string, 2)
	}
	config.Params["sql_mode"] = mysqlSQLMode
	config.Params["time_zone"] = mysqlTimezone
	config.ParseTime = true

	db, err := sql.Open("mysql", config.FormatDSN())
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
	if err := verifyMySQLSupportsCheck(timeoutContext, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
