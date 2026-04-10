package testutil

import (
	"errors"
	"slices"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestCleanupRegistryCloseAllRunsLIFO(t *testing.T) {
	registry := &CleanupRegistry{}
	calls := []string{}

	registry.Register(func() error {
		calls = append(calls, "first")
		return nil
	})
	registry.Register(func() error {
		calls = append(calls, "second")
		return nil
	})
	registry.Register(func() error {
		calls = append(calls, "third")
		return nil
	})

	err := registry.CloseAll()
	assert.NilError(t, "CloseAll() error", err)

	want := []string{"third", "second", "first"}
	if !slices.Equal(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestCleanupRegistryCloseAllJoinsErrors(t *testing.T) {
	registry := &CleanupRegistry{}
	firstErr := errors.New("first")
	secondErr := errors.New("second")

	registry.Register(func() error { return firstErr })
	registry.Register(func() error { return secondErr })

	err := registry.CloseAll()
	assert.Error(t, "CloseAll() error", err)
	if !errors.Is(err, firstErr) {
		t.Fatalf("CloseAll() error = %v, must include firstErr", err)
	}
	if !errors.Is(err, secondErr) {
		t.Fatalf("CloseAll() error = %v, must include secondErr", err)
	}
}

func TestCleanupRegistryCloseAllIsIdempotent(t *testing.T) {
	registry := &CleanupRegistry{}
	count := 0

	registry.Register(func() error {
		count++
		return nil
	})

	err := registry.CloseAll()
	assert.NilError(t, "[1st] CloseAll() error", err)
	assert.Equal(t, "[1st] count", count, 1)

	err = registry.CloseAll()
	assert.NilError(t, "[2nd] CloseAll() error", err)
	assert.Equal(t, "[2nd] count", count, 1)
}

func TestCleanupRegistryRegisterNilAndAfterCloseDoesNothing(t *testing.T) {
	registry := &CleanupRegistry{}
	count := 0

	registry.Register(nil)
	registry.Register(func() error {
		count++
		return nil
	})

	err := registry.CloseAll()
	assert.NilError(t, "[1st] CloseAll() error", err)
	assert.Equal(t, "[1st] count", count, 1)

	registry.Register(func() error {
		count++
		return nil
	})
	err = registry.CloseAll()
	assert.NilError(t, "[2nd] CloseAll() error", err)
	assert.Equal(t, "[2nd] count", count, 1)
}
