package sqlite_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestJobRepositoryMarks(t *testing.T) {
	testutilSQLite.NewJobSuite(t).RunTestMarks(t)
}

func TestJobRepositoryMarkStartedReturnsConflict(t *testing.T) {
	testutilSQLite.NewJobSuite(t).RunTestMarkStartedReturnsConflict(t)
}

func TestJobRepositoryMarkFinishedReturnsConflict(t *testing.T) {
	testutilSQLite.NewJobSuite(t).RunTestMarkFinishedReturnsConflict(t)
}
