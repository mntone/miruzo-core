package config

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

func newValidDatabaseConfig(dbBackend backend.Backend, dsn string) DatabaseConfig {
	cfg := DefaultDatabaseConfig()
	cfg.Backend = dbBackend
	cfg.DSN = dsn
	return cfg
}

func TestDatabaseConfigValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cfg         DatabaseConfig
		wantErrText string
	}{
		{
			name: "mysql_with_dsn_is_valid",
			cfg: newValidDatabaseConfig(
				backend.MySQL,
				"user:pass@tcp(127.0.0.1:3306)/miruzo",
			),
		},
		{
			name: "postgresql_with_supported_prefix_is_valid",
			cfg: newValidDatabaseConfig(
				backend.PostgreSQL,
				"postgresql://user:pass@127.0.0.1:5432/miruzo",
			),
		},
		{
			name: "sqlite_with_file_prefix_is_valid",
			cfg: newValidDatabaseConfig(
				backend.SQLite,
				"file:../var/miruzo.sqlite",
			),
		},
		{
			name: "dsn_is_required",
			cfg: func() DatabaseConfig {
				cfg := newValidDatabaseConfig(backend.MySQL, "")
				return cfg
			}(),
			wantErrText: "dsn must not be empty",
		},
		{
			name: "postgresql_requires_postgresql_prefix",
			cfg: newValidDatabaseConfig(
				backend.PostgreSQL,
				"postgres://user:pass@127.0.0.1:5432/miruzo",
			),
			wantErrText: "dsn must start with 'postgresql://' prefix",
		},
		{
			name: "sqlite_requires_file_prefix",
			cfg: newValidDatabaseConfig(
				backend.SQLite,
				"sqlite:../var/miruzo.sqlite",
			),
			wantErrText: "dsn must start with 'file:' prefix",
		},
		{
			name: "unsupported_backend_is_rejected",
			cfg: newValidDatabaseConfig(
				backend.Backend("postgres"),
				"postgresql://user:pass@127.0.0.1:5432/miruzo",
			),
			wantErrText: "backend must be one of 'mysql', 'postgresql' or 'sqlite'",
		},
		{
			name: "max_open_conns_must_be_positive",
			cfg: func() DatabaseConfig {
				cfg := newValidDatabaseConfig(
					backend.MySQL,
					"user:pass@tcp(127.0.0.1:3306)/miruzo",
				)
				cfg.MaxOpenConnections = 0
				return cfg
			}(),
			wantErrText: "max_open_conns must be >= 1",
		},
		{
			name: "pool_warm_conns_must_be_positive",
			cfg: func() DatabaseConfig {
				cfg := newValidDatabaseConfig(
					backend.MySQL,
					"user:pass@tcp(127.0.0.1:3306)/miruzo",
				)
				cfg.PoolWarmConnections = 0
				return cfg
			}(),
			wantErrText: "pool_warm_conns must be between 1 and max_open_conns",
		},
		{
			name: "pool_warm_conns_must_not_exceed_max_open_conns",
			cfg: func() DatabaseConfig {
				cfg := newValidDatabaseConfig(
					backend.MySQL,
					"user:pass@tcp(127.0.0.1:3306)/miruzo",
				)
				cfg.MaxOpenConnections = 2
				cfg.PoolWarmConnections = 3
				return cfg
			}(),
			wantErrText: "pool_warm_conns must be between 1 and max_open_conns",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()
			if tt.wantErrText == "" {
				if err != nil {
					t.Fatalf("Validate() unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Validate() expected error %q, got nil", tt.wantErrText)
			}
			if err.Error() != tt.wantErrText {
				t.Fatalf("Validate() error = %q, want %q", err.Error(), tt.wantErrText)
			}
		})
	}
}
