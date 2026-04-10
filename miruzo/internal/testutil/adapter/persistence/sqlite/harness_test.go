package sqlite

import (
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestSQLiteHarnessBackend(t *testing.T) {
	h := sqliteHarness{}
	assert.Equal(t, "Backend()", h.Backend(), backend.SQLite)
}

func TestSQLiteHarnessSupportsCapabilityMatrix(t *testing.T) {
	h := sqliteHarness{}
	expected := map[contract.Capability]bool{
		contract.SupportsLastInsertID:      true,
		contract.SupportsReturningClause:   true,
		contract.SupportsInfinityTimestamp: false,
	}

	for _, capability := range contract.AllCapabilities() {
		want, ok := expected[capability]
		if !ok {
			t.Fatalf("missing expected capability: %s", capability)
		}

		got := h.Supports(capability)
		assert.Equal(t, fmt.Sprintf("Supports(%s)", capability), got, want)
	}
}
