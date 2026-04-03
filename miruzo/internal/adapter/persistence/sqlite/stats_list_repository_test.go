package sqlite_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestStatsListRepositoryIterateStatsPaginatesWithoutDuplicates(t *testing.T) {
	testutilSQLite.NewStatsListSuite(t).RunTestIterateStatsPaginatesWithoutDuplicates(t)
}

func TestStatsListRepositoryIterateStatsReturnsEmpty(t *testing.T) {
	testutilSQLite.NewStatsListSuite(t).RunTestIterateStatsReturnsEmpty(t)
}
