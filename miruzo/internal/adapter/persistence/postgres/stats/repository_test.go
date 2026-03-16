package stats_test

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

func TestStatsRepositoryApplyLoveUpdatesWhenEmpty(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyLoveUpdatesWhenEmpty(t)
}

func TestStatsRepositoryApplyLoveRejectsCurrentPeriod(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyLoveRejectsCurrentPeriod(t)
}

func TestStatsRepositoryApplyLoveCanceledUpdatesTimestamps(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyLoveCanceledUpdatesTimestamps(t)
}

func TestStatsRepositoryApplyLoveCanceledReturnsConflict(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyLoveCanceledReturnsConflict(t)
}

func TestStatsRepositoryApplyLoveCanceledReturnsConflictWithoutStats(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyLoveCanceledReturnsConflictWithoutStats(t)
}

func TestStatsRepositoryApplyView(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyView(t)
}

func TestStatsRepositoryApplyViewNotFound(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyViewNotFound(t)
}

func TestStatsRepositoryApplyViewWithMilestone(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyViewWithMilestone(t)
}

func TestStatsRepositoryApplyViewWithMilestoneNotFound(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyViewWithMilestoneNotFound(t)
}
