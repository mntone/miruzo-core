package imagelist_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestImageListRepositoryListLatest(t *testing.T) {
	testutilSQLite.NewImageListSuite(t).RunTestListLatest(t)
}

func TestImageListRepositoryListChronological(t *testing.T) {
	testutilSQLite.NewImageListSuite(t).RunTestListChronological(t)
}

func TestImageListRepositoryListRecently(t *testing.T) {
	testutilSQLite.NewImageListSuite(t).RunTestListRecently(t)
}

func TestImageListRepositoryListFirstLove(t *testing.T) {
	testutilSQLite.NewImageListSuite(t).RunTestListFirstLove(t)
}

func TestImageListRepositoryListHallOfFame(t *testing.T) {
	testutilSQLite.NewImageListSuite(t).RunTestListHallOfFame(t)
}

func TestImageListRepositoryListEngaged(t *testing.T) {
	testutilSQLite.NewImageListSuite(t).RunTestListEngaged(t)
}
