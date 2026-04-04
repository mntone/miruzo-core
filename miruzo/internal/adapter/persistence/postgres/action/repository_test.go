package action_test

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

func TestActionRepositoryActionSchemaRejectsInvalidKind(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewAction(t, ctx).RunTestActionSchemaRejectsInvalidKind(t)
}

func TestActionRepositoryActionSchemaRejectsInvalidOccurredAt(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewAction(t, ctx).RunTestActionSchemaRejectsInvalidOccurredAt(t)
}

func TestActionRepositoryCreate(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewAction(t, ctx).RunTestCreate(t)
}

func TestActionRepositoryExistsSinceReturnsFalse(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewAction(t, ctx).RunTestExistsSinceReturnsFalse(t)
}

func TestActionRepositoryExistsSinceReturnsTrue(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewAction(t, ctx).RunTestExistsSinceReturnsTrue(t)
}
