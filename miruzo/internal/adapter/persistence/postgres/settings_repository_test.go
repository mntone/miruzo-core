package postgres_test

import (
	"context"
	"testing"
)

func TestSettingsRepositoryGetValueReturnsNotFound(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewSettings(t, ctx).RunTestGetValueReturnsNotFound(t)
}

func TestSettingsRepositoryUpdateAndGetValue(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewSettings(t, ctx).RunTestUpdateAndGetValue(t)
}

func TestSettingsRepositoryUpdateValueReturnsCheckViolation(t *testing.T) {
	if factory == nil {
		t.Fatal("suite is nil")
	}

	ctx := context.Background()
	factory.MustReset(t, ctx)
	factory.NewSettings(t, ctx).RunTestUpdateValueReturnsCheckViolation(t)
}
