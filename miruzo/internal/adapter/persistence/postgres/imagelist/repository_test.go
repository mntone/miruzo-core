package imagelist_test

import (
	"context"
	"log"
	"os"
	"testing"

	testutilPostgres "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/postgres"
)

var factory *testutilPostgres.SuiteFactory

func TestMain(m *testing.M) {
	ctx := context.Background()

	// setup
	testFactory, err := testutilPostgres.NewSuiteFactory(ctx)
	if err != nil {
		log.Printf("setup postgres test suite: %v", err)
		os.Exit(1)
	}
	factory = testFactory

	exitCode := m.Run()

	// teardown
	if factory != nil {
		if err := factory.Close(); err != nil {
			log.Printf("teardown postgres test suite: %v", err)
			if exitCode == 0 {
				exitCode = 1
			}
		}
	}

	os.Exit(exitCode)
}

func TestImageListRepositoryListLatest(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewImageList(t, ctx).RunTestListLatest(t)
}

func TestImageListRepositoryListChronological(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewImageList(t, ctx).RunTestListChronological(t)
}

func TestImageListRepositoryListRecently(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewImageList(t, ctx).RunTestListRecently(t)
}

func TestImageListRepositoryListFirstLove(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewImageList(t, ctx).RunTestListFirstLove(t)
}

func TestImageListRepositoryListHallOfFame(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewImageList(t, ctx).RunTestListHallOfFame(t)
}

func TestImageListRepositoryListEngaged(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewImageList(t, ctx).RunTestListEngaged(t)
}
