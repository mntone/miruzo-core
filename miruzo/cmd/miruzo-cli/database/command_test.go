package database

import (
	"errors"
	"testing"

	cliio "github.com/mntone/miruzo-core/miruzo/internal/cli/io"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func resetDatabaseFlags() {
	adminDatabase = ""
	adminUserName = ""
	adminPassword = ""
	adminPasswordStdin = false
	adminPasswordEnv = ""
}

func TestValidateSQLiteAdminOverrides(t *testing.T) {
	baseConfig := config.DatabaseConfig{
		Backend: backend.SQLite,
	}

	t.Run("NoOverrides", func(t *testing.T) {
		resetDatabaseFlags()
		err := validateSQLiteAdminOverrides(baseConfig)
		assert.NilError(t, "validateSQLiteAdminOverrides() error", err)
	})

	t.Run("CLIAdminDatabase", func(t *testing.T) {
		resetDatabaseFlags()
		adminDatabase = "sqlite"
		err := validateSQLiteAdminOverrides(baseConfig)
		assert.Error(t, "validateSQLiteAdminOverrides() error", err)
		assert.Equal(
			t,
			"validateSQLiteAdminOverrides() error",
			err.Error(),
			"sqlite backend does not support --admin-database",
		)
	})

	t.Run("CLIAdminUserName", func(t *testing.T) {
		resetDatabaseFlags()
		adminUserName = "root"
		err := validateSQLiteAdminOverrides(baseConfig)
		assert.Error(t, "validateSQLiteAdminOverrides() error", err)
		assert.Equal(
			t,
			"validateSQLiteAdminOverrides() error",
			err.Error(),
			"sqlite backend does not support --admin-username",
		)
	})

	t.Run("ConfigAdminDatabase", func(t *testing.T) {
		resetDatabaseFlags()
		cfg := baseConfig
		cfg.AdminDatabase = "sqlite"
		err := validateSQLiteAdminOverrides(cfg)
		assert.Error(t, "validateSQLiteAdminOverrides() error", err)
		assert.Equal(
			t,
			"validateSQLiteAdminOverrides() error",
			err.Error(),
			"sqlite backend does not support database.admin_database",
		)
	})

	t.Run("ConfigAdminUserName", func(t *testing.T) {
		resetDatabaseFlags()
		cfg := baseConfig
		cfg.AdminUserName = "root"
		err := validateSQLiteAdminOverrides(cfg)
		assert.Error(t, "validateSQLiteAdminOverrides() error", err)
		assert.Equal(
			t,
			"validateSQLiteAdminOverrides() error",
			err.Error(),
			"sqlite backend does not support database.admin_username",
		)
	})
}

func TestResolveAdminDatabase(t *testing.T) {
	tests := []struct {
		name    string
		backend backend.Backend
		cli     string
		config  string
		want    string
	}{
		{
			name:    "ResolveFromCLIOverride",
			backend: backend.MySQL, // fallback: mysql
			cli:     "cli_admin",
			config:  "cfg_admin",
			want:    "cli_admin",
		},
		{
			name:    "ResolveFromConfigWhenCLIEmpty",
			backend: backend.MySQL, // fallback: mysql
			cli:     "",
			config:  "cfg_admin",
			want:    "cfg_admin",
		},
		{
			name:    "ResolveDefaultForMySQLWhenUnset",
			backend: backend.MySQL, // fallback: mysql
			cli:     "",
			config:  "",
			want:    "mysql",
		},
		{
			name:    "ResolveDefaultForPostgreSQLWhenUnset",
			backend: backend.PostgreSQL, // fallback: postgres
			cli:     "",
			config:  "",
			want:    "postgres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDatabase = tt.cli
			config := config.DatabaseConfig{
				Backend:       tt.backend,
				AdminDatabase: tt.config,
			}

			got := resolveAdminDatabase(config)
			assert.Equal(
				t,
				"resolveAdminDatabase()",
				got,
				tt.want,
			)
		})
	}
}

func TestResolveAdminUserName(t *testing.T) {
	tests := []struct {
		name    string
		backend backend.Backend
		cli     string
		config  string
		want    string
	}{
		{
			name:    "ResolveFromCLIOverride",
			backend: backend.MySQL,
			cli:     "cli_user",
			config:  "cfg_user",
			want:    "cli_user",
		},
		{
			name:    "ResolveFromConfigWhenCLIEmpty",
			backend: backend.MySQL,
			cli:     "",
			config:  "cfg_user",
			want:    "cfg_user",
		},
		{
			name:    "ResolveDefaultForMySQLWhenUnset",
			backend: backend.MySQL,
			cli:     "",
			config:  "",
			want:    "",
		},
		{
			name:    "ResolveDefaultForPostgreSQLWhenUnset",
			backend: backend.PostgreSQL,
			cli:     "",
			config:  "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminUserName = tt.cli
			config := config.DatabaseConfig{
				Backend:       tt.backend,
				AdminUserName: tt.config,
			}

			got := resolveAdminUserName(config)
			assert.Equal(
				t,
				"resolveAdminUserName()",
				got,
				tt.want,
			)
		})
	}
}

func TestDatabaseCommandAdminDatabaseFlag(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantDBName  string
		wantHasFlag bool
		wantErr     bool
	}{
		{
			name:        "default_empty",
			args:        []string{},
			wantDBName:  "",
			wantHasFlag: true,
		},
		{
			name:        "long_flag",
			args:        []string{"--admin-database", "postgres"},
			wantDBName:  "postgres",
			wantHasFlag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetDatabaseFlags()
			err := Command.ParseFlags(tt.args)
			if tt.wantErr {
				assert.Error(t, "Command.ParseFlags() error", err)
			} else {
				assert.NilError(t, "Command.ParseFlags() error", err)
			}

			assert.Equal(t, "adminDatabaseName", adminDatabase, tt.wantDBName)
			hasFlag := Command.Flags().Lookup("admin-database") != nil
			assert.Equal(t, "admin-database flag exists", hasFlag, tt.wantHasFlag)
		})
	}
}

func TestDatabaseCommandAdminUserNameFlag(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantUser    string
		wantHasFlag bool
		wantErr     bool
	}{
		{
			name:        "default_empty",
			args:        []string{},
			wantUser:    "",
			wantHasFlag: true,
		},
		{
			name:        "long_flag",
			args:        []string{"--admin-username", "root"},
			wantUser:    "root",
			wantHasFlag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetDatabaseFlags()
			err := Command.ParseFlags(tt.args)
			if tt.wantErr {
				assert.Error(t, "Command.ParseFlags() error", err)
			} else {
				assert.NilError(t, "Command.ParseFlags() error", err)
			}

			assert.Equal(t, "adminUserName", adminUserName, tt.wantUser)
			hasFlag := Command.Flags().Lookup("admin-username") != nil
			assert.Equal(t, "admin-username flag exists", hasFlag, tt.wantHasFlag)
		})
	}
}

func TestResolveAdminPassword(t *testing.T) {
	orig := readAdminPasswordFn
	t.Cleanup(func() {
		readAdminPasswordFn = orig
	})

	t.Run("none", func(t *testing.T) {
		resetDatabaseFlags()
		readAdminPasswordFn = func(
			mode cliio.ReadPasswordMode,
			envName string,
		) (string, error) {
			return "", nil
		}

		got, err := resolveAdminPassword()
		assert.NilError(t, "resolveAdminPassword() error", err)
		assert.Equal(t, "resolveAdminPassword()", got, "")
	})

	t.Run("from_flag", func(t *testing.T) {
		resetDatabaseFlags()
		adminPassword = "secret"
		readAdminPasswordFn = func(
			mode cliio.ReadPasswordMode,
			envName string,
		) (string, error) {
			t.Fatalf("readAdminPasswordFn should not be called")
			return "", nil
		}

		got, err := resolveAdminPassword()
		assert.NilError(t, "resolveAdminPassword() error", err)
		assert.Equal(t, "resolveAdminPassword()", got, "secret")
	})

	t.Run("from_stdin", func(t *testing.T) {
		resetDatabaseFlags()
		adminPasswordStdin = true
		readAdminPasswordFn = func(
			mode cliio.ReadPasswordMode,
			envName string,
		) (string, error) {
			assert.Equal(t, "mode", mode, cliio.ReadPasswordModeStdin)
			assert.Equal(t, "envName", envName, "")
			return "stdin-secret", nil
		}

		got, err := resolveAdminPassword()
		assert.NilError(t, "resolveAdminPassword() error", err)
		assert.Equal(t, "resolveAdminPassword()", got, "stdin-secret")
	})

	t.Run("from_env", func(t *testing.T) {
		resetDatabaseFlags()
		adminPasswordEnv = "MIRUZO_ADMIN_PASSWORD"
		readAdminPasswordFn = func(
			mode cliio.ReadPasswordMode,
			envName string,
		) (string, error) {
			assert.Equal(t, "mode", mode, cliio.ReadPasswordModeEnv)
			assert.Equal(t, "envName", envName, "MIRUZO_ADMIN_PASSWORD")
			return "env-secret", nil
		}

		got, err := resolveAdminPassword()
		assert.NilError(t, "resolveAdminPassword() error", err)
		assert.Equal(t, "resolveAdminPassword()", got, "env-secret")
	})

	t.Run("conflict", func(t *testing.T) {
		resetDatabaseFlags()
		adminPassword = "secret"
		adminPasswordStdin = true

		_, err := resolveAdminPassword()
		assert.Error(t, "resolveAdminPassword() error", err)
		if err.Error() !=
			"--admin-password, --admin-password-stdin and --admin-password-env are mutually exclusive" {
			t.Fatalf(
				"resolveAdminPassword() error = %q, want %q",
				err.Error(),
				"--admin-password, --admin-password-stdin and --admin-password-env are mutually exclusive",
			)
		}
	})

	t.Run("conflict_with_env", func(t *testing.T) {
		resetDatabaseFlags()
		adminPasswordStdin = true
		adminPasswordEnv = "MIRUZO_ADMIN_PASSWORD"

		_, err := resolveAdminPassword()
		assert.Error(t, "resolveAdminPassword() error", err)
		if err.Error() !=
			"--admin-password, --admin-password-stdin and --admin-password-env are mutually exclusive" {
			t.Fatalf(
				"resolveAdminPassword() error = %q, want %q",
				err.Error(),
				"--admin-password, --admin-password-stdin and --admin-password-env are mutually exclusive",
			)
		}
	})

	t.Run("stdin_error", func(t *testing.T) {
		resetDatabaseFlags()
		adminPasswordStdin = true
		readAdminPasswordFn = func(
			mode cliio.ReadPasswordMode,
			envName string,
		) (string, error) {
			return "", errors.New("read password failed")
		}

		_, err := resolveAdminPassword()
		assert.Error(t, "resolveAdminPassword() error", err)
		if err.Error() != "read password failed" {
			t.Fatalf(
				"resolveAdminPassword() error = %q, want %q",
				err.Error(),
				"read password failed",
			)
		}
	})

	t.Run("stdin_not_terminal_error_adds_hint", func(t *testing.T) {
		resetDatabaseFlags()
		adminPasswordStdin = true
		readAdminPasswordFn = func(
			mode cliio.ReadPasswordMode,
			envName string,
		) (string, error) {
			return "", cliio.ErrStdinNotTerminal
		}

		_, err := resolveAdminPassword()
		assert.Error(t, "resolveAdminPassword() error", err)
		assert.ErrorIs(t, "resolveAdminPassword() error", err, cliio.ErrStdinNotTerminal)
		if err.Error() != "stdin is not a terminal; use --admin-password-stdin" {
			t.Fatalf(
				"resolveAdminPassword() error = %q, want %q",
				err.Error(),
				"stdin is not a terminal; use --admin-password-stdin",
			)
		}
	})
}
