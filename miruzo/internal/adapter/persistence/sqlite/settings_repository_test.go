package sqlite_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestSettingsRepositoryGetValueReturnsNotFound(t *testing.T) {
	testutilSQLite.NewSettingsSuite(t).RunTestGetValueReturnsNotFound(t)
}

func TestSettingsRepositoryUpdateAndGetValue(t *testing.T) {
	testutilSQLite.NewSettingsSuite(t).RunTestUpdateAndGetValue(t)
}

func TestSettingsRepositoryUpdateValueReturnsCheckViolation(t *testing.T) {
	testutilSQLite.NewSettingsSuite(t).RunTestUpdateValueReturnsCheckViolation(t)
}
