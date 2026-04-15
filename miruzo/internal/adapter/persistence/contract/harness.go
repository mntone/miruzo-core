package contract

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type TxCallback func(t *testing.T, ops TxSession)

type Harness interface {
	Dialect

	Backend() backend.Backend
	Supports(capability Capability) bool
	RequireCapability(t testing.TB, capability Capability)

	BeginTx(t testing.TB) TxSession
	RunInTx(t *testing.T, callback TxCallback)
}

func RequireCapability(
	t testing.TB, h interface {
		Backend() backend.Backend
		Supports(cap Capability) bool
	}, capability Capability,
) {
	t.Helper()
	if !h.Supports(capability) {
		t.Skipf("%s does not support %s", h.Backend(), capability)
	}
}
