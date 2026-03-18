package user_test

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

func TestUserRepositoryUserSchemaRejectsInvalidDailyLoveUsed(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestUserSchemaRejectsInvalidDailyLoveUsed(t)
}

func TestUserRepositoryGetSingletonUser(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestGetSingletonUser(t)
}

func TestUserRepositoryIncrementDailyLoveUsedIncrements(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestIncrementDailyLoveUsedIncrements(t)
}

func TestUserRepositoryIncrementDailyLoveUsedReturnsQuotaExceededWhenMissing(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestIncrementDailyLoveUsedReturnsQuotaExceededWhenMissing(t)
}

func TestUserRepositoryIncrementDailyLoveUsedReturnsQuotaExceededWhenLimitReached(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestIncrementDailyLoveUsedReturnsQuotaExceededWhenLimitReached(t)
}

func TestUserRepositoryDecrementDailyLoveUsedDecrements(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestDecrementDailyLoveUsedDecrements(t)
}

func TestUserRepositoryDecrementDailyLoveUsedReturnsNotFound(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestDecrementDailyLoveUsedReturnsNotFound(t)
}

func TestUserRepositoryDecrementDailyLoveUsedReturnsQuotaUnderflow(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestDecrementDailyLoveUsedReturnsQuotaUnderflow(t)
}
