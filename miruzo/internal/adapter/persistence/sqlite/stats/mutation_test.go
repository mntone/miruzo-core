package stats_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestStatsRepositoryApplyLoveUpdatesWhenEmpty(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyLoveUpdatesWhenEmpty(t)
}

func TestStatsRepositoryApplyLoveReturnsConflict(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyLoveReturnsConflict(t)
}

func TestStatsRepositoryApplyLoveCanceledUpdatesTimestamps(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyLoveCanceledUpdatesTimestamps(t)
}

func TestStatsRepositoryApplyLoveCanceledReturnsConflict(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyLoveCanceledReturnsConflict(t)
}

func TestStatsRepositoryApplyLoveCanceledReturnsConflictWithoutStats(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyLoveCanceledReturnsConflictWithoutStats(t)
}
