package assert_test

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNotEmpty(t *testing.T) {
	requirePass(t, "slice", func(t *testing.T) {
		assert.NotEmpty(t, "slice", []int{1})
	})
}

func TestEmpty(t *testing.T) {
	requirePass(t, "slice", func(t *testing.T) {
		assert.Empty(t, "slice", []int{})
	})
}

func TestLenIs(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	requirePass(t, "slice", func(t *testing.T) {
		assert.LenIs(t, "slice", []int{1, 2}, 2)
	})
	requirePass(t, "array", func(t *testing.T) {
		assert.LenIs(t, "array", [2]int{1, 2}, 2)
	})
	requirePass(t, "map", func(t *testing.T) {
		assert.LenIs(t, "map", m, 2)
	})
	requirePass(t, "string", func(t *testing.T) {
		assert.LenIs(t, "str", "ab", 2)
	})
}

func TestLenIsPanicsForUnsupportedType(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("LenIs() panic = nil, want non-nil")
		}
	}()

	assert.LenIs(t, "value", 10, 1)
}
