package assert_test

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestContains(t *testing.T) {
	requirePass(t, "contains", func(t *testing.T) {
		assert.Contains(t, "value", "operation=UpdateArticle affected_rows=0", "affected_rows=0")
	})
}
