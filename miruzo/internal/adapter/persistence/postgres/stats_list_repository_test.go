package postgres_test

import (
	"context"
	"testing"
)

func TestStatsListRepositoryIterateStatsPaginatesWithoutDuplicates(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStatsList(t, ctx).RunTestIterateStatsPaginatesWithoutDuplicates(t)
}

func TestStatsListRepositoryIterateStatsReturnsEmpty(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewStatsList(t, ctx).RunTestIterateStatsReturnsEmpty(t)
}
