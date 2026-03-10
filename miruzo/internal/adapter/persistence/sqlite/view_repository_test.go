package sqlite_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestViewRepositoryGetImageWithStats(t *testing.T) {
	testutilSQLite.NewViewSuite(t).RunTestGetImageWithStats(t)
}
