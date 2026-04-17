package mysql

import (
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestMySQLHarnessBackend(t *testing.T) {
	h := mysqlHarness{}
	assert.Equal(t, "Backend()", h.Backend(), backend.MySQL)
}

func TestMySQLHarnessSupportsCapabilityMatrix(t *testing.T) {
	h := mysqlHarness{}
	expected := map[contract.Capability]bool{
		contract.SupportsInfinityTimestamp:     false,
		contract.SupportsLastInsertID:          true,
		contract.SupportsNumberedPlaceholder:   false,
		contract.SupportsReturningClause:       false,
		contract.SupportsUnnumberedPlaceholder: true,
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
