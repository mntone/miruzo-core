package config

import (
	"errors"
	"strings"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type DatabaseConfig struct {
	Backend       backend.Backend `mapstructure:"backend"`
	DSN           string          `mapstructure:"dsn"`
	AdminDatabase string          `mapstructure:"admin_database"`
	AdminUserName string          `mapstructure:"admin_username"`

	ConnectionTimeout     time.Duration `mapstructure:"conn_timeout"`
	PoolWarmConnections   int32         `mapstructure:"pool_warm_conns"`
	MaxOpenConnections    int32         `mapstructure:"max_open_conns"`
	MaxConnectionIdleTime time.Duration `mapstructure:"max_conn_idletime"`
	MaxConnectionLifeTime time.Duration `mapstructure:"max_conn_lifetime"`
}

func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		ConnectionTimeout:     20 * time.Second,
		PoolWarmConnections:   1,
		MaxOpenConnections:    2,
		MaxConnectionIdleTime: 5 * time.Minute,
		MaxConnectionLifeTime: 30 * time.Minute,
	}
}

func (c *DatabaseConfig) Validate() error {
	if c.DSN == "" {
		return errors.New("dsn must not be empty")
	}
	switch c.Backend {
	case backend.MySQL:
		// no operation
	case backend.PostgreSQL:
		if !strings.HasPrefix(c.DSN, "postgresql://") {
			return errors.New("dsn must start with 'postgresql://' prefix")
		}
	case backend.SQLite:
		if !strings.HasPrefix(c.DSN, "file:") {
			return errors.New("dsn must start with 'file:' prefix")
		}
		if c.AdminDatabase != "" {
			return errors.New("admin_database is not supported for sqlite backend")
		}
		if c.AdminUserName != "" {
			return errors.New("admin_username is not supported for sqlite backend")
		}
	default:
		return errors.New("backend must be one of 'mysql', 'postgresql' or 'sqlite'")
	}
	if c.MaxOpenConnections < 1 {
		return errors.New("max_open_conns must be >= 1")
	}
	if c.PoolWarmConnections < 1 || c.PoolWarmConnections > c.MaxOpenConnections {
		return errors.New("pool_warm_conns must be between 1 and max_open_conns")
	}
	return nil
}
