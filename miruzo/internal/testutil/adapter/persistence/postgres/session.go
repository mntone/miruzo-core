package postgres

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres"
	dbshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

type postgresTxSession struct {
	postgresDialect
	tx     pgx.Tx
	nextID model.IngestIDType
	persist.SessionRepositories
}

func newTxSession(tx pgx.Tx) *postgresTxSession {
	return &postgresTxSession{
		tx:                  tx,
		nextID:              modelbuilder.GetNextID(),
		SessionRepositories: postgres.NewSessionRepositories(gen.New(tx)),
	}
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

func (s postgresTxSession) ExecAndGetRowCount(t testing.TB, stmt string, args ...any) (int64, error) {
	result, err := s.tx.Exec(t.Context(), stmt, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (s postgresTxSession) Rollback(t testing.TB) {
	t.Helper()
	err := s.tx.Rollback(t.Context())
	assert.NilError(t, "Rollback()", err)

	modelbuilder.SetNextID(s.nextID)
}

func (s postgresTxSession) InsertImage(t testing.TB, e persist.Image) error {
	_, err := s.tx.Exec(
		t.Context(),
		"INSERT INTO images(ingest_id,ingested_at,kind,original,fallback,variants)VALUES($1,$2,$3,$4,$5,$6)",
		e.IngestID, e.IngestedAt, e.Type,
		e.Original, e.Fallback.ToPointer(), e.Layers,
	)
	if err != nil {
		return dbshared.MapPostgreError("InsertImage", err)
	}

	return nil
}

func (s postgresTxSession) SelectStats(t testing.TB, id model.IngestIDType) (model.Stats, error) {
	row := s.tx.QueryRow(t.Context(), "SELECT * FROM stats WHERE ingest_id=$1", id)

	var e model.Stats
	var scoreEvaluatedAt, firstLovedAt, lastLovedAt, hallOfFameAt, lastViewedAt, viewMilestoneArchivedAt *time.Time
	err := row.Scan(
		&e.IngestID,
		&e.Score,
		&e.ScoreEvaluated,
		&scoreEvaluatedAt,
		&firstLovedAt,
		&lastLovedAt,
		&hallOfFameAt,
		&lastViewedAt,
		&e.ViewCount,
		&e.ViewMilestoneCount,
		&viewMilestoneArchivedAt,
	)
	if err != nil {
		return e, err
	}

	e.ScoreEvaluatedAt = mo.PointerToOption(scoreEvaluatedAt)
	e.FirstLovedAt = mo.PointerToOption(firstLovedAt)
	e.LastLovedAt = mo.PointerToOption(lastLovedAt)
	e.HallOfFameAt = mo.PointerToOption(hallOfFameAt)
	e.LastViewedAt = mo.PointerToOption(lastViewedAt)
	e.ViewMilestoneArchivedAt = mo.PointerToOption(viewMilestoneArchivedAt)
	return e, nil
}
