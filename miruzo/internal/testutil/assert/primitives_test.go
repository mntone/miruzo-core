package assert_test

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestEqual(t *testing.T) {
	requirePass(t, "equal", func(t *testing.T) {
		assert.Equal(t, "value", 10, 10)
	})
}

func TestEqualFn(t *testing.T) {
	now := time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC)

	requirePass(t, "equal", func(t *testing.T) {
		assert.EqualFn(t, "time", now, now)
	})
}

func TestGreaterThan(t *testing.T) {
	requirePass(t, "int", func(t *testing.T) {
		assert.GreaterThan(t, "value", 11, 10)
	})
	requirePass(t, "string", func(t *testing.T) {
		assert.GreaterThan(t, "value", "b", "a")
	})
}

func TestLessThan(t *testing.T) {
	requirePass(t, "int", func(t *testing.T) {
		assert.LessThan(t, "value", 9, 10)
	})
	requirePass(t, "string", func(t *testing.T) {
		assert.LessThan(t, "value", "a", "b")
	})
}
