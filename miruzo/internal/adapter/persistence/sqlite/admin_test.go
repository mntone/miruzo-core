package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestOpenAdminHandleRejectsAdminDatabaseName(t *testing.T) {
	t.Parallel()

	databasePath := filepath.Join(t.TempDir(), "admin-name.sqlite")
	_, err := OpenAdminHandle(config.DatabaseConfig{
		DSN:               "file:" + databasePath,
		AdminDatabaseName: "ignored",
	}, "")
	assert.Error(t, "OpenAdminHandle() error", err)

	if err.Error() !=
		"sqlite backend does not support admin database override: \"ignored\"" {
		t.Fatalf(
			"OpenAdminHandle() error = %q, want %q",
			err.Error(),
			"sqlite backend does not support admin database override: \"ignored\"",
		)
	}
}

func TestOpenAdminHandleRejectsAdminDatabaseNameOverride(t *testing.T) {
	t.Parallel()

	databasePath := filepath.Join(t.TempDir(), "admin-name-override.sqlite")
	_, err := OpenAdminHandle(config.DatabaseConfig{
		DSN:               "file:" + databasePath,
		AdminDatabaseName: "config-admin",
	}, "cli-admin")
	assert.Error(t, "OpenAdminHandle() error", err)

	if err.Error() !=
		"sqlite backend does not support admin database override: \"cli-admin\"" {
		t.Fatalf(
			"OpenAdminHandle() error = %q, want %q",
			err.Error(),
			"sqlite backend does not support admin database override: \"cli-admin\"",
		)
	}
}

func TestDatabaseAdminCreateCreatesFile(t *testing.T) {
	t.Parallel()

	databasePath := filepath.Join(t.TempDir(), "create.sqlite")
	admin, err := OpenAdminHandle(config.DatabaseConfig{
		DSN: "file:" + databasePath,
	}, "")
	assert.NilError(t, "OpenAdminHandle() error", err)

	err = admin.Create(context.Background())
	assert.NilError(t, "Create() error", err)

	_, err = os.Stat(databasePath)
	assert.NilError(t, "Stat() error", err)
}

func TestDatabaseAdminCreateReturnsErrExistWhenFileExists(t *testing.T) {
	t.Parallel()

	databasePath := filepath.Join(t.TempDir(), "exists.sqlite")
	err := os.WriteFile(databasePath, []byte("x"), 0o644)
	assert.NilError(t, "WriteFile() error", err)

	admin, err := OpenAdminHandle(config.DatabaseConfig{
		DSN: "file:" + databasePath,
	}, "")
	assert.NilError(t, "OpenAdminHandle() error", err)

	err = admin.Create(context.Background())
	assert.ErrorIs(t, "Create() error", err, os.ErrExist)
}

func TestDatabaseAdminDropRemovesFile(t *testing.T) {
	t.Parallel()

	databasePath := filepath.Join(t.TempDir(), "drop.sqlite")
	err := os.WriteFile(databasePath, []byte("x"), 0o644)
	assert.NilError(t, "WriteFile() error", err)

	admin, err := OpenAdminHandle(config.DatabaseConfig{
		DSN: "file:" + databasePath,
	}, "")
	assert.NilError(t, "OpenAdminHandle() error", err)

	err = admin.Drop(context.Background())
	assert.NilError(t, "Drop() error", err)

	_, err = os.Stat(databasePath)
	assert.ErrorIs(t, "Stat() error", err, os.ErrNotExist)
}

func TestDatabaseAdminExistsReturnsTrueWhenFileExists(t *testing.T) {
	t.Parallel()

	databasePath := filepath.Join(t.TempDir(), "exists-check.sqlite")
	err := os.WriteFile(databasePath, []byte("x"), 0o644)
	assert.NilError(t, "WriteFile() error", err)

	admin, err := OpenAdminHandle(config.DatabaseConfig{
		DSN: "file:" + databasePath,
	}, "")
	assert.NilError(t, "OpenAdminHandle() error", err)

	exists, err := admin.Exists(context.Background())
	assert.NilError(t, "Exists() error", err)
	assert.Equal(t, "Exists()", exists, true)
}

func TestDatabaseAdminExistsReturnsFalseWhenFileDoesNotExist(t *testing.T) {
	t.Parallel()

	databasePath := filepath.Join(t.TempDir(), "not-found-check.sqlite")
	admin, err := OpenAdminHandle(config.DatabaseConfig{
		DSN: "file:" + databasePath,
	}, "")
	assert.NilError(t, "OpenAdminHandle() error", err)

	exists, err := admin.Exists(context.Background())
	assert.NilError(t, "Exists() error", err)
	assert.Equal(t, "Exists()", exists, false)
}
