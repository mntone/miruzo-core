package postgres_test

import (
	"context"
	"testing"
)

func TestJobRepositoryMarks(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewJob(t, ctx).RunTestMarks(t)
}

func TestJobRepositoryMarkStartedReturnsConflict(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewJob(t, ctx).RunTestMarkStartedReturnsConflict(t)
}

func TestJobRepositoryMarkFinishedReturnsConflict(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewJob(t, ctx).RunTestMarkFinishedReturnsConflict(t)
}
