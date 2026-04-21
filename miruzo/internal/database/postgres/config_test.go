package postgres

import (
	"testing"
	"time"

	sharedDB "github.com/mntone/miruzo-core/miruzo/internal/database/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNewConnectConfigFromDSNAppliesOptions(t *testing.T) {
	t.Parallel()

	options := ConnectOptions{
		ConnectionTuning: sharedDB.ConnectionTuning{
			ConnectionTimeout:     7 * time.Second,
			PoolWarmConnections:   2,
			MaxOpenConnections:    4,
			MaxConnectionIdleTime: 11 * time.Second,
			MaxConnectionLifeTime: 13 * time.Second,
		},
	}

	cfg, err := NewConnectConfigFromDSN(
		"postgres://app:secret@localhost:5432/appdb?sslmode=disable",
		options,
	)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	assert.Equal(t, "Database()", cfg.Database(), "appdb")
	assert.Equal(
		t,
		"ConnectTimeout",
		cfg.baseConfig.ConnConfig.ConnectTimeout,
		options.ConnectionTimeout,
	)
	assert.Equal(
		t,
		"timezone runtime param",
		cfg.baseConfig.ConnConfig.RuntimeParams["timezone"],
		"UTC",
	)
	assert.Equal(t, "MaxConnIdleTime", cfg.baseConfig.MaxConnIdleTime, options.MaxConnectionIdleTime)
	assert.Equal(t, "MaxConnLifetime", cfg.baseConfig.MaxConnLifetime, options.MaxConnectionLifeTime)
	assert.Equal(t, "MinConns", cfg.baseConfig.MinConns, options.PoolWarmConnections)
	assert.Equal(t, "MaxConns", cfg.baseConfig.MaxConns, options.MaxOpenConnections)
}

func TestNewConnectConfigFromDSNReturnsErrorOnInvalidDSN(t *testing.T) {
	t.Parallel()

	_, err := NewConnectConfigFromDSN(":::", ConnectOptions{})
	assert.Error(t, "NewConnectConfigFromDSN() error", err)
}

func TestWithCredentialsDoesNotMutateOriginal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConnectConfigFromDSN(
		"postgres://app:secret@localhost:5432/appdb?sslmode=disable",
		ConnectOptions{},
	)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	updated := cfg.WithCredentials("root", "admin_password")

	if cfg.baseConfig.ConnConfig.User != "app" || cfg.baseConfig.ConnConfig.Password != "secret" {
		t.Fatalf(
			"original credentials = (%q, %q), want (%q, %q)",
			cfg.baseConfig.ConnConfig.User,
			cfg.baseConfig.ConnConfig.Password,
			"app",
			"secret",
		)
	}
	if updated.baseConfig.ConnConfig.User != "root" ||
		updated.baseConfig.ConnConfig.Password != "admin_password" {
		t.Fatalf(
			"updated credentials = (%q, %q), want (%q, %q)",
			updated.baseConfig.ConnConfig.User,
			updated.baseConfig.ConnConfig.Password,
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

	original, err := NewConnectConfigFromDSN(
		"postgres://app:secret@localhost:5432/appdb?sslmode=disable",
		ConnectOptions{},
	)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	updated := original.WithDatabase("postgres")
	assert.Equal(t, "original Database()", original.Database(), "appdb")
	assert.Equal(t, "updated Database()", updated.Database(), "postgres")
	if updated.baseConfig == original.baseConfig {
		t.Fatalf("updated baseConfig pointer must differ from original")
	}
}
