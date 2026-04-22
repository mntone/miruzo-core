package database

import (
	"errors"
	"testing"

	cliio "github.com/mntone/miruzo-core/miruzo/internal/cli/io"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func resetDatabaseFlags() {
	adminDatabaseName = ""
	adminUserName = ""
	adminPassword = ""
	adminPasswordStdin = false
	adminPasswordEnv = ""
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

			assert.Equal(t, "adminDatabaseName", adminDatabaseName, tt.wantDBName)
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
