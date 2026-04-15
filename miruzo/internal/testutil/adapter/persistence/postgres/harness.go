package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type postgresHarness struct {
	postgresDialect
	pool *pgxpool.Pool
}

func NewHarness(t testing.TB, reg *testutil.CleanupRegistry) postgresHarness {
	return postgresHarness{
		pool: GetPostgresTestPool(t, reg),
	}
}

func (h postgresHarness) Backend() backend.Backend {
	return backend.PostgreSQL
}

func (h postgresHarness) Supports(capability contract.Capability) bool {
	switch capability {
	case contract.SupportsReturningClause,
		contract.SupportsInfinityTimestamp:
		return true
	}
	return false
}

func (h postgresHarness) RequireCapability(t testing.TB, capability contract.Capability) {
	t.Helper()
	contract.RequireCapability(t, h, capability)
}

func (h postgresHarness) BeginTx(t testing.TB) contract.TxSession {
	t.Helper()

	tx, err := h.pool.BeginTx(t.Context(), pgx.TxOptions{})
	assert.NilError(t, "BeginTx()", err)

	txSession := newTxSession(tx)
	return contract.TxSession{
		Dialect:               txSession,
		TransactionOperations: txSession,
		RepositoryProvider:    txSession,
	}
}

func (h postgresHarness) RunInTx(t *testing.T, callback contract.TxCallback) {
	t.Helper()

	session := h.BeginTx(t)
	defer session.Rollback(t)
	callback(t, session)
}
