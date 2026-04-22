package postgres

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

type ConnectOptions struct {
	UseSimpleProtocol bool
	shared.ConnectionTuning
}

type ConnectConfig struct {
	baseConfig *pgxpool.Config
}

func NewConnectConfigFromDSN(
	dsn string,
	options ConnectOptions,
) (*ConnectConfig, error) {
	baseConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	baseConfig.ConnConfig.ConnectTimeout = options.ConnectionTimeout
	if options.UseSimpleProtocol {
		baseConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
		baseConfig.ConnConfig.StatementCacheCapacity = 0
		baseConfig.ConnConfig.DescriptionCacheCapacity = 0
	}
	baseConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"
	baseConfig.MaxConnIdleTime = options.MaxConnectionIdleTime
	baseConfig.MaxConnLifetime = options.MaxConnectionLifeTime
	baseConfig.MinConns = options.PoolWarmConnections
	baseConfig.MaxConns = options.MaxOpenConnections
	return &ConnectConfig{
		baseConfig: baseConfig,
	}, nil
}

func (c *ConnectConfig) WithCredentials(userName, password string) *ConnectConfig {
	newConfig := c.baseConfig.Copy()
	newConfig.ConnConfig.User = userName
	newConfig.ConnConfig.Password = password
	return &ConnectConfig{
		baseConfig: newConfig,
	}
}

func (c *ConnectConfig) Database() string {
	return c.baseConfig.ConnConfig.Database
}

func (c *ConnectConfig) WithDatabase(databaseName string) *ConnectConfig {
	newConfig := c.baseConfig.Copy()
	newConfig.ConnConfig.Database = databaseName
	return &ConnectConfig{
		baseConfig: newConfig,
	}
}
