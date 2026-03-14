package config

import "time"

type DatabaseBackend string

const (
	DatabaseBackendSQLite   DatabaseBackend = "sqlite"
	DatabaseBackendPostgres DatabaseBackend = "postgres"
)

type DatabaseConfig struct {
	Backend DatabaseBackend `mapstructure:"backend"`
	DSN     string          `mapstructure:"dsn"`

	ConnectionTimeout     time.Duration `mapstructure:"conn_timeout"`
	MaxOpenConnections    int32         `mapstructure:"max_open_connections"`
	MaxConnectionIdleTime time.Duration `mapstructure:"max_conn_idletime"`
	MaxConnectionLifeTime time.Duration `mapstructure:"max_conn_lifetime"`
}

func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		ConnectionTimeout:     30 * time.Second,
		MaxOpenConnections:    1,
		MaxConnectionIdleTime: 5 * time.Minute,
		MaxConnectionLifeTime: 30 * time.Minute,
	}
}
