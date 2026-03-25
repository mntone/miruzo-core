package assert_test

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestEqualMap(t *testing.T) {
	requirePass(t, "equal", func(t *testing.T) {
		assert.EqualMap(
			t,
			"map",
			map[string]int{"a": 1, "b": 2},
			map[string]int{"a": 1, "b": 2},
		)
	})
}
