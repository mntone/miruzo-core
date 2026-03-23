package stats_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestStatsRepositoryStatsSchemaRejectsInvalidScore(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestStatsSchemaRejectsInvalidScore(t)
}

func TestStatsRepositoryStatsSchemaRejectsInvalidScoreEvaluated(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestStatsSchemaRejectsInvalidScoreEvaluated(t)
}

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

func TestStatsRepositoryApplyLoveRejectsCurrentPeriod(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyLoveRejectsCurrentPeriod(t)
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

func TestStatsRepositoryApplyView(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyView(t)
}

func TestStatsRepositoryApplyViewNotFound(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyViewNotFound(t)
}

func TestStatsRepositoryApplyViewWithMilestone(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyViewWithMilestone(t)
}

func TestStatsRepositoryApplyViewWithMilestoneNotFound(t *testing.T) {
	testutilSQLite.NewStatsSuite(t).RunTestApplyViewWithMilestoneNotFound(t)
}
