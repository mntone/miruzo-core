package stats_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestStatsListRepositoryIterateStatsPaginatesWithoutDuplicates(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestIterateStatsPaginatesWithoutDuplicates(t)
}

func TestStatsListRepositoryIterateStatsReturnsEmpty(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestIterateStatsReturnsEmpty(t)
}
