package mysql

import (
	"testing"
	"time"

	sharedDB "github.com/mntone/miruzo-core/miruzo/internal/database/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNewConnectConfigFromDSNAppliesOptions(t *testing.T) {
	t.Parallel()

	options := ConnectOptions{
		MultiStatements: true,
		ConnectionTuning: sharedDB.ConnectionTuning{
			ConnectionTimeout:     7 * time.Second,
			PoolWarmConnections:   2,
			MaxOpenConnections:    4,
			MaxConnectionIdleTime: 11 * time.Second,
			MaxConnectionLifeTime: 13 * time.Second,
		},
	}

	cfg, err := NewConnectConfigFromDSN(
		"app:secret@tcp(localhost:3306)/appdb",
		options,
	)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	assert.Equal(t, "Database()", cfg.Database(), "appdb")
	assert.Equal(t, "MultiStatements", cfg.baseConfig.MultiStatements, true)
	assert.Equal(t, "ParseTime", cfg.baseConfig.ParseTime, true)
	assert.Equal(t, "Params[\"sql_mode\"]", cfg.baseConfig.Params["sql_mode"], mysqlSQLMode)
	assert.Equal(t, "Params[\"time_zone\"]", cfg.baseConfig.Params["time_zone"], mysqlTimezone)
	if cfg.connTune != options.ConnectionTuning {
		t.Fatalf("connTune = %#v, want %#v", cfg.connTune, options.ConnectionTuning)
	}
}

func TestNewConnectConfigFromDSNReturnsErrorOnInvalidDSN(t *testing.T) {
	t.Parallel()

	_, err := NewConnectConfigFromDSN(":::", ConnectOptions{})
	assert.Error(t, "NewConnectConfigFromDSN() error", err)
}

func TestWithCredentialsDoesNotMutateOriginal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConnectConfigFromDSN("app:secret@tcp(localhost:3306)/appdb", ConnectOptions{})
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	updated := cfg.WithCredentials("root", "admin_password")

	if cfg.baseConfig.User != "app" || cfg.baseConfig.Passwd != "secret" {
		t.Fatalf(
			"original credentials = (%q, %q), want (%q, %q)",
			cfg.baseConfig.User,
			cfg.baseConfig.Passwd,
			"app",
			"secret",
		)
	}
	if updated.baseConfig.User != "root" || updated.baseConfig.Passwd != "admin_password" {
		t.Fatalf(
			"updated credentials = (%q, %q), want (%q, %q)",
			updated.baseConfig.User,
			updated.baseConfig.Passwd,
			"root",
			"admin_password",
		)
	}
	if updated.baseConfig == cfg.baseConfig {
		t.Fatalf("updated baseConfig pointer must differ from original")
	}
}

func TestWithDatabaseDoesNotMutateOriginal(t *testing.T) {
	t.Parallel()

	original, err := NewConnectConfigFromDSN("app:secret@tcp(localhost:3306)/appdb", ConnectOptions{})
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	updated := original.WithDatabase("mysql")
	assert.Equal(t, "original Database()", original.Database(), "appdb")
	assert.Equal(t, "updated Database()", updated.Database(), "mysql")
	if updated.baseConfig == original.baseConfig {
		t.Fatalf("updated baseConfig pointer must differ from original")
	}
}
