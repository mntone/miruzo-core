package assert_test

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNil(t *testing.T) {
	requirePass(t, "nil", func(t *testing.T) {
		var p *int = nil
		assert.Nil(t, "f()", p)
	})
}

func TestNotNil(t *testing.T) {
	requirePass(t, "i", func(t *testing.T) {
		i := 0
		assert.NotNil(t, "f()", &i)
	})
}
