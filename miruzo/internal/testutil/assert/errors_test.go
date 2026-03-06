package assert_test

import (
	"errors"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNilError(t *testing.T) {
	requirePass(t, "nil", func(t *testing.T) {
		assert.NilError(t, "f()", nil)
	})
}

func TestErrorIs(t *testing.T) {
	root := errors.New("root")
	wrapped := errors.Join(errors.New("other"), root)

	requirePass(t, "match", func(t *testing.T) {
		assert.ErrorIs(t, "f()", wrapped, root)
	})
}
