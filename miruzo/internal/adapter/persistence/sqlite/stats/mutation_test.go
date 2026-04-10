package stats_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestStatsRepositoryApplyHallOfFameGrantedUpdates(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyHallOfFameGrantedUpdates(t)
}

func TestStatsRepositoryApplyHallOfFameGrantedReturnsConflict(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyHallOfFameGrantedReturnsConflict(t)
}

func TestStatsRepositoryApplyHallOfFameGrantedReturnsConflictWithoutStats(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyHallOfFameGrantedReturnsConflictWithoutStats(t)
}

func TestStatsRepositoryApplyHallOfFameRevokedUpdates(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyHallOfFameRevokedUpdates(t)
}

func TestStatsRepositoryApplyHallOfFameRevokedReturnsConflict(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyHallOfFameRevokedReturnsConflict(t)
}

func TestStatsRepositoryApplyHallOfFameRevokedReturnsConflictWithoutStats(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyHallOfFameRevokedReturnsConflictWithoutStats(t)
}

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
