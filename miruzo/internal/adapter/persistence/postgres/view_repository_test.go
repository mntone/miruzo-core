package postgres_test

import (
	"context"
	"testing"
)

func TestViewRepositoryGetImageWithStats(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewView(t, ctx).RunTestGetImageWithStats(t)
}
