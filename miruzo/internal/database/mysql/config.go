package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

type ConnectOptions struct {
	MultiStatements bool
	shared.ConnectionTuning
}

type ConnectConfig struct {
	baseConfig *mysql.Config
	connTune   shared.ConnectionTuning
}

const (
	mysqlSQLMode  = "'TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY'"
	mysqlTimezone = "'+00:00'"
)

func NewConnectConfigFromDSN(
	dsn string,
	options ConnectOptions,
) (*ConnectConfig, error) {
	baseConfig, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	baseConfig.MultiStatements = options.MultiStatements
	if baseConfig.Params == nil {
		baseConfig.Params = make(map[string]string, 2)
	}
	baseConfig.Params["sql_mode"] = mysqlSQLMode
	baseConfig.Params["time_zone"] = mysqlTimezone
	baseConfig.ParseTime = true
	return &ConnectConfig{
		baseConfig: baseConfig,
		connTune:   options.ConnectionTuning,
	}, nil
}

func (c *ConnectConfig) WithCredentials(userName, password string) *ConnectConfig {
	newConfig := c.baseConfig.Clone()
	newConfig.User = userName
	newConfig.Passwd = password
	return &ConnectConfig{
		baseConfig: newConfig,
		connTune:   c.connTune,
	}
}

func (c *ConnectConfig) Database() string {
	return c.baseConfig.DBName
}

func (c *ConnectConfig) WithDatabase(databaseName string) *ConnectConfig {
	newConfig := c.baseConfig.Clone()
	newConfig.DBName = databaseName
	return &ConnectConfig{
		baseConfig: newConfig,
		connTune:   c.connTune,
	}
}
