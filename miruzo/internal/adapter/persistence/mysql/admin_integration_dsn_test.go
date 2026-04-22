package mysql_test

import (
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func withMySQLDatabaseName(dsn string, databaseName string) (string, error) {
	cfg, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		return "", err
	}

	cfg.DBName = databaseName
	return cfg.FormatDSN(), nil
}

func TestWithMySQLDatabaseName(t *testing.T) {
	t.Parallel()

	const srcDSN = "app:secret@tcp(localhost:3306)/miruzo_test?parseTime=true"
	gotDSN, err := withMySQLDatabaseName(srcDSN, "miruzo_tmp")
	assert.NilError(t, "withMySQLDatabaseName() error", err)

	gotCfg, err := mysqlDriver.ParseDSN(gotDSN)
	assert.NilError(t, "ParseDSN(gotDSN) error", err)
	assert.Equal(t, "database", gotCfg.DBName, "miruzo_tmp")
	assert.Equal(t, "user", gotCfg.User, "app")
	assert.Equal(t, "address", gotCfg.Addr, "localhost:3306")
	assert.Equal(t, "parseTime", gotCfg.ParseTime, true)
}

func TestWithMySQLDatabaseNameReturnsErrorOnInvalidDSN(t *testing.T) {
	t.Parallel()

	_, err := withMySQLDatabaseName(":::", "miruzo_tmp")
	assert.Error(t, "withMySQLDatabaseName() error", err)
}
