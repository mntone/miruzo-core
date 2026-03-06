package imagelist_test

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestImageListRepositoryRunSuite(t *testing.T) {
	testutilPersistence.RunImageListSuite(t, func(t testing.TB) testutilPersistence.ImageListSetup {
		setup := testutilSQLite.SetupOperations(t)
		repo := imagelist.NewRepository(setup.DB)
		return testutilPersistence.ImageListSetup{
			Ctx:  setup.Ctx,
			Ops:  setup.Ops,
			Repo: repo,
		}
	})
}
