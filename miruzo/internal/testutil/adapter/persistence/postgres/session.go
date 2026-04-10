package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/action"
	dbshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

type postgresTxSession struct {
	tx      pgx.Tx
	queries *gen.Queries
	nextID  model.IngestIDType
}

func newTxSession(tx pgx.Tx) *postgresTxSession {
	return &postgresTxSession{
		tx:      tx,
		queries: gen.New(tx),
		nextID:  modelbuilder.GetNextID(),
	}
}

// --- Dialect ---

func (s postgresTxSession) Backend() backend.Backend {
	return backend.PostgreSQL
}

func (s postgresTxSession) MapError(
	operation string,
	err error,
	mapping contract.DBErrorMapping,
) error {
	switch mapping {
	case contract.DBErrorMappingDefault:
		err = dbshared.MapPostgreError(operation, err)
	case contract.DBErrorMappingDelete:
		err = dbshared.MapPostgreDeleteError(operation, err)
	}
	return err
}

func (s postgresTxSession) BindVarStyle() contract.BindVarStyle {
	return contract.BindVarStyleDollar
}

// --- TransactionOperations ---

func (s postgresTxSession) Exec(t testing.TB, stmt string, args ...any) error {
	_, err := s.tx.Exec(t.Context(), stmt, args...)
	return err
}

func (s postgresTxSession) ExecInsertAndGetID(t testing.TB, stmt string, args ...any) (int64, error) {
	return -1, contract.ErrUnsupportedCapability
}

func (s postgresTxSession) ExecReturningInt64(t testing.TB, stmt string, args ...any) (int64, error) {
	row := s.tx.QueryRow(t.Context(), stmt, args...)
	var ret int64 = -1
	err := row.Scan(&ret)
	return ret, err
}

func (s postgresTxSession) Rollback(t testing.TB) {
	t.Helper()
	err := s.tx.Rollback(t.Context())
	assert.NilError(t, "Rollback()", err)

	modelbuilder.SetNextID(s.nextID)
}

// --- RepositoryProvider ---

func (s postgresTxSession) Action() persist.ActionRepository {
	return action.NewRepository(s.queries)
}
