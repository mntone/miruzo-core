package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/postgres"
)

const (
	adminVersionStmt  = "SELECT current_setting('server_version_num')::int >= 170000"
	adminCreateStmt   = "CREATE DATABASE %s ENCODING 'UTF8' LC_COLLATE 'C' LC_CTYPE 'C' TEMPLATE template0"
	adminCreateStmt17 = "CREATE DATABASE %s LOCALE_PROVIDER builtin BUILTIN_LOCALE 'C.UTF-8' TEMPLATE template0"
	adminDropStmt     = "DROP DATABASE %s"
	adminExistsStmt   = "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)"
)

type postgresAdminHandle struct {
	pool *pgxpool.Pool

	appDatabaseName string
}

func OpenAdminHandle(
	ctx context.Context,
	appConfig config.DatabaseConfig,
	options persistshared.DatabaseAdminOptions,
) (postgresAdminHandle, error) {
	connOptions := database.ConnectOptions{
		UseSimpleProtocol: true,
		ConnectionTuning:  persistshared.NewConnectionTuningFromConfig(appConfig),
	}
	connOptions.PoolWarmConnections = 1
	connOptions.MaxOpenConnections = 1

	cfg, err := database.NewConnectConfigFromDSN(appConfig.DSN, connOptions)
	if err != nil {
		return postgresAdminHandle{}, err
	}

	adminUserName, adminPassword := options.ResolveCredentials(
		cfg.UserName(),
		cfg.Password(),
	)
	adminDatabaseName := options.ResolveDatabaseName(
		appConfig.AdminDatabaseName,
		"postgres",
	)
	appDatabaseName := cfg.Database()

	adminConfig := cfg.
		WithCredentials(adminUserName, adminPassword).
		WithDatabase(adminDatabaseName)
	pool, err := database.Open(ctx, adminConfig)
	if err != nil {
		return postgresAdminHandle{}, err
	}

	return postgresAdminHandle{
		pool: pool,

		appDatabaseName: appDatabaseName,
	}, nil
}

func (hdl postgresAdminHandle) Close() error {
	hdl.pool.Close()
	return nil
}

func postgresQuoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}

func (hdl postgresAdminHandle) Create(ctx context.Context) error {
	row := hdl.pool.QueryRow(ctx, adminVersionStmt)

	var supportsBuiltinLocaleProvider bool
	err := row.Scan(&supportsBuiltinLocaleProvider)
	if err != nil {
		return fmt.Errorf("postgres admin query server version failed: %w", err)
	}

	var createStmt string
	if supportsBuiltinLocaleProvider {
		createStmt = adminCreateStmt17
	} else {
		createStmt = adminCreateStmt
	}

	// PostgreSQL may return pgerrcode.DuplicateDatabase ("42P04").
	_, err = hdl.pool.Exec(
		ctx,
		fmt.Sprintf(createStmt, postgresQuoteIdentifier(hdl.appDatabaseName)),
	)
	if err != nil {
		return fmt.Errorf(
			"postgres admin create database %q failed: %w",
			hdl.appDatabaseName,
			err,
		)
	}
	return nil
}

func (hdl postgresAdminHandle) Drop(ctx context.Context) error {
	// PostgreSQL may return pgerrcode.InvalidCatalogName ("3D000").
	_, err := hdl.pool.Exec(
		ctx,
		fmt.Sprintf(adminDropStmt, postgresQuoteIdentifier(hdl.appDatabaseName)),
	)
	if err != nil {
		return fmt.Errorf(
			"postgres admin drop database %q failed: %w",
			hdl.appDatabaseName,
			err,
		)
	}
	return nil
}

func (hdl postgresAdminHandle) Exists(ctx context.Context) (bool, error) {
	row := hdl.pool.QueryRow(ctx, adminExistsStmt, hdl.appDatabaseName)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf(
			"postgres admin check database %q exists failed: %w",
			hdl.appDatabaseName,
			err,
		)
	}

	return exists, nil
}
