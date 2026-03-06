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
