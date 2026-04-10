package stats_test

import (
	"context"
	"testing"
)

func TestStatsRepositoryApplyDecayUpdates(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyDecayUpdates(t)
}

func TestStatsRepositoryApplyDecayReturnsConflictWithoutStats(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyDecayReturnsConflictWithoutStats(t)
}

func TestStatsRepositoryApplyHallOfFameGrantedUpdates(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyHallOfFameGrantedUpdates(t)
}

func TestStatsRepositoryApplyHallOfFameGrantedReturnsConflict(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyHallOfFameGrantedReturnsConflict(t)
}

func TestStatsRepositoryApplyHallOfFameGrantedReturnsConflictWithoutStats(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyHallOfFameGrantedReturnsConflictWithoutStats(t)
}

func TestStatsRepositoryApplyHallOfFameRevokedUpdates(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyHallOfFameRevokedUpdates(t)
}

func TestStatsRepositoryApplyHallOfFameRevokedReturnsConflict(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyHallOfFameRevokedReturnsConflict(t)
}

func TestStatsRepositoryApplyHallOfFameRevokedReturnsConflictWithoutStats(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStats(t, ctx).RunTestApplyHallOfFameRevokedReturnsConflictWithoutStats(t)
}

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
