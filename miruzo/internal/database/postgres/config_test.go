package postgres

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	sharedDB "github.com/mntone/miruzo-core/miruzo/internal/database/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func assertCredentials(
	t *testing.T,
	name string,
	got *pgconn.Config,
	wantUserName string,
	wantPassword string,
) {
	t.Helper()

	if got.User != wantUserName || got.Password != wantPassword {
		t.Fatalf(
			"%s = (%q, %q), want (%q, %q)",
			name,
			got.User, got.Password,
			wantUserName, wantPassword,
		)
	}
}

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

	assertCredentials(
		t,
		"credentials",
		&cfg.baseConfig.ConnConfig.Config,
		"app", "secret",
	)

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

	original, err := NewConnectConfigFromDSN(
		"postgres://app:secret@localhost:5432/appdb?sslmode=disable",
		ConnectOptions{},
	)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	updated := original.WithCredentials("root", "admin_password")
	assertCredentials(
		t,
		"original credentials",
		&original.baseConfig.ConnConfig.Config,
		"app", "secret",
	)
	assertCredentials(
		t,
		"updated credentials",
		&updated.baseConfig.ConnConfig.Config,
		"root", "admin_password",
	)

	if updated.baseConfig == original.baseConfig {
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

func TestWithMutatorPreservesConnectionOptions(t *testing.T) {
	t.Parallel()

	original, err := NewConnectConfigFromDSN(
		"postgres://app:secret@localhost:5432/appdb?sslmode=disable",
		ConnectOptions{
			UseSimpleProtocol: true,
			ConnectionTuning: sharedDB.ConnectionTuning{
				ConnectionTimeout:     7 * time.Second,
				PoolWarmConnections:   2,
				MaxOpenConnections:    4,
				MaxConnectionIdleTime: 11 * time.Second,
				MaxConnectionLifeTime: 13 * time.Second,
			},
		},
	)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	assertPreservedOptions := func(t *testing.T, updated *ConnectConfig) {
		t.Helper()

		assert.Equal(
			t,
			"ConnectTimeout",
			updated.baseConfig.ConnConfig.ConnectTimeout,
			7*time.Second,
		)
		assert.Equal(
			t,
			"DefaultQueryExecMode",
			updated.baseConfig.ConnConfig.DefaultQueryExecMode,
			pgx.QueryExecModeSimpleProtocol,
		)
		assert.Equal(
			t,
			"StatementCacheCapacity",
			updated.baseConfig.ConnConfig.StatementCacheCapacity,
			0,
		)
		assert.Equal(
			t,
			"DescriptionCacheCapacity",
			updated.baseConfig.ConnConfig.DescriptionCacheCapacity,
			0,
		)
		assert.Equal(
			t,
			"timezone runtime param",
			updated.baseConfig.ConnConfig.RuntimeParams["timezone"],
			"UTC",
		)
		assert.Equal(
			t,
			"MaxConnIdleTime",
			updated.baseConfig.MaxConnIdleTime,
			11*time.Second,
		)
		assert.Equal(
			t,
			"MaxConnLifetime",
			updated.baseConfig.MaxConnLifetime,
			13*time.Second,
		)
		assert.Equal(t, "MinConns", updated.baseConfig.MinConns, int32(2))
		assert.Equal(t, "MaxConns", updated.baseConfig.MaxConns, int32(4))
	}

	tests := []struct {
		name   string
		mutate func(*ConnectConfig) *ConnectConfig
		assert func(*testing.T, *ConnectConfig)
	}{
		{
			name: "WithCredentials",
			mutate: func(cfg *ConnectConfig) *ConnectConfig {
				return cfg.WithCredentials("root", "admin_password")
			},
			assert: func(t *testing.T, updated *ConnectConfig) {
				t.Helper()
				assert.Equal(
					t,
					"updated user",
					updated.baseConfig.ConnConfig.User,
					"root",
				)
				assert.Equal(
					t,
					"updated password",
					updated.baseConfig.ConnConfig.Password,
					"admin_password",
				)
				assert.Equal(
					t,
					"database",
					updated.baseConfig.ConnConfig.Database,
					"appdb",
				)
			},
		},
		{
			name: "WithDatabase",
			mutate: func(cfg *ConnectConfig) *ConnectConfig {
				return cfg.WithDatabase("postgres")
			},
			assert: func(t *testing.T, updated *ConnectConfig) {
				t.Helper()
				assert.Equal(
					t,
					"database",
					updated.baseConfig.ConnConfig.Database,
					"postgres",
				)
				assert.Equal(t, "user", updated.baseConfig.ConnConfig.User, "app")
				assert.Equal(
					t,
					"password",
					updated.baseConfig.ConnConfig.Password,
					"secret",
				)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			updated := tt.mutate(original)
			assert.NotNil(t, tt.name, updated)
			tt.assert(t, updated)
			assertPreservedOptions(t, updated)
		})
	}
}
