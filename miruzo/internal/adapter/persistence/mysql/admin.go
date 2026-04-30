package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/mysql"
)

const (
	adminCreateStmt = "CREATE DATABASE %s CHARSET utf8mb4 COLLATE utf8mb4_0900_bin"
	adminDropStmt   = "DROP DATABASE %s"
	adminExistsStmt = "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name=?)"
)

type mysqlAdminHandle struct {
	db *sql.DB

	appDatabase string
	appUserName string
}

func OpenAdminHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
	options persistshared.DatabaseAdminOptions,
) (mysqlAdminHandle, error) {
	connOptions := database.ConnectOptions{
		MultiStatements:  false,
		ConnectionTuning: persistshared.NewConnectionTuningFromConfig(appConfig),
	}
	connOptions.PoolWarmConnections = 1
	connOptions.MaxOpenConnections = 1

	cfg, err := database.NewConnectConfigFromDSN(appConfig.DSN, connOptions)
	if err != nil {
		return mysqlAdminHandle{}, err
	}

	adminUserName, adminPassword := options.ResolveCredentials(
		cfg.UserName(),
		cfg.Password(),
	)

	adminConfig := cfg.
		WithCredentials(adminUserName, adminPassword).
		WithDatabase(options.Database)
	db, err := database.Open(ctx, adminConfig)
	if err != nil {
		return mysqlAdminHandle{}, err
	}

	return mysqlAdminHandle{
		db: db,

		appDatabase: cfg.Database(),
		appUserName: cfg.UserName(),
	}, nil
}

func (hdl mysqlAdminHandle) Close() error {
	return hdl.db.Close()
}

func mysqlQuoteIdentifier(identifier string) string {
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`"
}

func (hdl mysqlAdminHandle) Create(ctx context.Context) error {
	// MySQL may return ER_DB_CREATE_EXISTS (1007, SQLSTATE HY000).
	_, err := hdl.db.ExecContext(
		ctx,
		fmt.Sprintf(adminCreateStmt, mysqlQuoteIdentifier(hdl.appDatabase)),
	)
	if err != nil {
		return fmt.Errorf(
			"mysql admin create database %q failed: %w",
			hdl.appDatabase,
			err,
		)
	}
	return nil
}

func (hdl mysqlAdminHandle) Drop(ctx context.Context) error {
	// MySQL may return ER_DB_DROP_EXISTS (1008, SQLSTATE HY000).
	_, err := hdl.db.ExecContext(
		ctx,
		fmt.Sprintf(adminDropStmt, mysqlQuoteIdentifier(hdl.appDatabase)),
	)
	if err != nil {
		return fmt.Errorf(
			"mysql admin drop database %q failed: %w",
			hdl.appDatabase,
			err,
		)
	}
	return nil
}

func (hdl mysqlAdminHandle) Exists(ctx context.Context) (bool, error) {
	row := hdl.db.QueryRowContext(ctx, adminExistsStmt, hdl.appDatabase)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf(
			"mysql admin check database %q exists failed: %w",
			hdl.appDatabase,
			err,
		)
	}

	return exists, nil
}
