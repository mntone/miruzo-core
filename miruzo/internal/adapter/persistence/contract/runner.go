package contract

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type HarnessCallback func(t *testing.T, h Harness)

func RunHarnesses(
	t *testing.T,
	harnesses []Harness,
	callback HarnessCallback,
) {
	t.Helper()
	for _, h := range harnesses {
		t.Run(h.Backend().String(), func(t *testing.T) {
			callback(t, h)
		})
	}
}

func RunDefaultHarnesses(
	t *testing.T,
	toHarnesses func(t *testing.T, backends []backend.Backend) []Harness,
	callback HarnessCallback,
) {
	t.Helper()
	backends := []backend.Backend{
		backend.PostgreSQL,
		backend.SQLite,
	}
	harnesses := toHarnesses(t, backends)
	RunHarnesses(t, harnesses, callback)
}
