package imagelist_test

import (
	"context"
	"log"
	"os"
	"testing"

	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
	testutilPostgre "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/postgre"
)

var suite *testutilPostgre.Suite

func TestMain(m *testing.M) {
	ctx := context.Background()

	// setup
	testSuite, err := testutilPostgre.NewSuite(ctx)
	if err != nil {
		log.Printf("setup postgre test suite: %v", err)
		os.Exit(1)
	}
	suite = testSuite

	exitCode := m.Run()

	// teardown
	if suite != nil {
		if err := suite.Close(); err != nil {
			log.Printf("teardown postgre test suite: %v", err)
			if exitCode == 0 {
				exitCode = 1
			}
		}
	}

	os.Exit(exitCode)
}

func TestImageListRepositoryRunSuite(t *testing.T) {
	if suite == nil {
		t.Fatal("suite is nil")
	}

	testutilPersistence.RunImageListSuite(t, func(tb testing.TB) testutilPersistence.ImageListSetup {
		tb.Helper()

		ctx := context.Background()
		suite.MustReset(tb, ctx)
		return suite.NewImageList(ctx)
	})
}
