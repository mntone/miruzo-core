package postgres

import (
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestPostgreHarnessBackend(t *testing.T) {
	h := postgresHarness{}
	assert.Equal(t, "Backend()", h.Backend(), backend.PostgreSQL)
}

func TestPostgreHarnessSupportsCapabilityMatrix(t *testing.T) {
	h := postgresHarness{}
	expected := map[contract.Capability]bool{
		contract.SupportsInfinityTimestamp:     true,
		contract.SupportsLastInsertID:          false,
		contract.SupportsNumberedPlaceholder:   true,
		contract.SupportsReturningClause:       true,
		contract.SupportsUnnumberedPlaceholder: false,
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
