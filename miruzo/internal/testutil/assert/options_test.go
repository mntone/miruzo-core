package assert_test

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestIsPresent(t *testing.T) {
	requirePass(t, "optional-present", func(t *testing.T) {
		assert.IsPresent(t, "option", testOptional{present: true})
	})
	requirePass(t, "non-optional-non-nil", func(t *testing.T) {
		assert.IsPresent(t, "value", struct{}{})
	})
}

func TestIsAbsent(t *testing.T) {
	requirePass(t, "optional-absent", func(t *testing.T) {
		assert.IsAbsent(t, "option", testOptional{present: false})
	})
	requirePass(t, "non-optional-nil", func(t *testing.T) {
		assert.IsAbsent(t, "value", nil)
	})
}
