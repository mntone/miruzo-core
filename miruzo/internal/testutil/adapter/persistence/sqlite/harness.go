package sqlite

import (
	"database/sql"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type sqliteHarness struct {
	sqliteDialect
	db *sql.DB
}

func NewHarness(t testing.TB, reg *testutil.CleanupRegistry) sqliteHarness {
	return sqliteHarness{
		db: GetSQLiteTestDB(t, reg),
	}
}

func (h sqliteHarness) Backend() backend.Backend {
	return backend.SQLite
}

func (sqliteHarness) Supports(capability contract.Capability) bool {
	switch capability {
	case contract.SupportsLastInsertID,
		contract.SupportsNumberedPlaceholder,
		contract.SupportsReturningClause,
		contract.SupportsUnnumberedPlaceholder:
		return true
	}
	return false
}

func (h sqliteHarness) RequireCapability(t testing.TB, capability contract.Capability) {
	t.Helper()
	contract.RequireCapability(t, h, capability)
}

func (h sqliteHarness) BeginTx(t testing.TB) contract.TxSession {
	t.Helper()

	tx, err := h.db.BeginTx(t.Context(), nil)
	assert.NilError(t, "BeginTx()", err)

	txSession := newTxSession(tx)
	return contract.TxSession{
		Dialect:               txSession,
		TransactionOperations: txSession,
		RepositoryProvider:    txSession,
	}
}

func (h sqliteHarness) RunInTx(t *testing.T, callback contract.TxCallback) {
	t.Helper()

	session := h.BeginTx(t)
	defer session.Rollback(t)
	callback(t, session)
}
