package stats_test

import (
	"context"
	"testing"
)

func TestStatsRepositoryApplyLoveUpdatesWhenEmpty(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyLoveUpdatesWhenEmpty(t)
}

func TestStatsRepositoryApplyLoveReturnsConflict(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyLoveReturnsConflict(t)
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
