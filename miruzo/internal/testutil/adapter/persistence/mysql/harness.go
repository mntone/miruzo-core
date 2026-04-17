package mysql

import (
	"database/sql"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type mysqlHarness struct {
	mysqlDialect
	db *sql.DB
}

func NewHarness(t testing.TB, reg *testutil.CleanupRegistry) mysqlHarness {
	return mysqlHarness{
		db: GetMySQLTestDB(t, reg),
	}
}

func (mysqlHarness) Backend() backend.Backend {
	return backend.MySQL
}

func (mysqlHarness) Supports(capability contract.Capability) bool {
	switch capability {
	case contract.SupportsLastInsertID,
		contract.SupportsUnnumberedPlaceholder:
		return true
	}
	return false
}

func (h mysqlHarness) RequireCapability(t testing.TB, capability contract.Capability) {
	t.Helper()
	contract.RequireCapability(t, h, capability)
}

func (h mysqlHarness) BeginTx(t testing.TB) contract.TxSession {
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

func (h mysqlHarness) RunInTx(t *testing.T, callback contract.TxCallback) {
	t.Helper()

	session := h.BeginTx(t)
	defer session.Rollback(t)
	callback(t, session)
}
