package stats_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

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
