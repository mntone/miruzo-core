package user_test

import (
	"context"
	"log"
	"os"
	"testing"

	testutilPostgre "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/postgre"
)

var factory *testutilPostgre.SuiteFactory

func TestMain(m *testing.M) {
	ctx := context.Background()

	// setup
	testFactory, err := testutilPostgre.NewSuiteFactory(ctx)
	if err != nil {
		log.Printf("setup postgre test suite: %v", err)
		os.Exit(1)
	}
	factory = testFactory

	exitCode := m.Run()

	// teardown
	if factory != nil {
		if err := factory.Close(); err != nil {
			log.Printf("teardown postgre test suite: %v", err)
			if exitCode == 0 {
				exitCode = 1
			}
		}
	}

	os.Exit(exitCode)
}

func TestUserRepositoryGetSingletonUser(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewUser(t, ctx).RunTestGetSingletonUser(t)
}
