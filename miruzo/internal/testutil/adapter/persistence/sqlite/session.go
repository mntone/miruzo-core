package sqlite

import (
	"database/sql"
	"encoding/json"
	"testing"

	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	sqliteshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

type sqliteTxSession struct {
	sqliteDialect
	tx     *sql.Tx
	nextID model.IngestIDType
	persist.SessionRepositories
}

func newTxSession(tx *sql.Tx) *sqliteTxSession {
	return &sqliteTxSession{
		tx:                  tx,
		nextID:              modelbuilder.GetNextID(),
		SessionRepositories: sqlite.NewSessionRepositories(gen.New(tx)),
	}
}

// --- TransactionOperations ---

func (s sqliteTxSession) Exec(t testing.TB, stmt string, args ...any) error {
	_, err := s.tx.ExecContext(t.Context(), stmt, args...)
	return err
}

func (s sqliteTxSession) ExecInsertAndGetID(t testing.TB, stmt string, args ...any) (int64, error) {
	row, err := s.tx.ExecContext(t.Context(), stmt, args...)
	if err != nil {
		return -1, err
	}

	return row.LastInsertId()
}

func (s sqliteTxSession) ExecReturningInt64(t testing.TB, stmt string, args ...any) (int64, error) {
	row := s.tx.QueryRowContext(t.Context(), stmt, args...)
	var ret int64 = -1
	err := row.Scan(&ret)
	return ret, err
}

func (s sqliteTxSession) ExecAndGetRowCount(t testing.TB, stmt string, args ...any) (int64, error) {
	result, err := s.tx.ExecContext(t.Context(), stmt, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s sqliteTxSession) Rollback(t testing.TB) {
	t.Helper()
	err := s.tx.Rollback()
	assert.NilError(t, "Rollback()", err)

	modelbuilder.SetNextID(s.nextID)
}

func (s sqliteTxSession) InsertImage(t testing.TB, e persist.Image) error {
	originalBytes, err := json.Marshal(e.Original)
	if err != nil {
		return err
	}

	var fallbackBytes *[]byte
	if fallbackValue, present := e.Fallback.Get(); present {
		bytes, err := json.Marshal(fallbackValue)
		if err != nil {
			return err
		}

		fallbackBytes = &bytes
	}

	layersBytes, err := json.Marshal(e.Layers)
	if err != nil {
		return err
	}

	_, err = s.tx.ExecContext(
		t.Context(),
		"INSERT INTO images(ingest_id,ingested_at,kind,original,fallback,variants)VALUES(?,?,?,?,?,?)",
		e.IngestID, e.IngestedAt, e.Type,
		originalBytes, fallbackBytes, layersBytes,
	)
	if err != nil {
		return sqliteshared.MapSQLiteError("InsertImage", err)
	}

	return nil
}

func (s sqliteTxSession) SelectStats(t testing.TB, id model.IngestIDType) (model.Stats, error) {
	row := s.tx.QueryRowContext(t.Context(), "SELECT * FROM stats WHERE ingest_id=?", id)

	var e model.Stats
	var scoreEvaluatedAt, firstLovedAt, lastLovedAt, hallOfFameAt, lastViewedAt, viewMilestoneArchivedAt sql.NullTime
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

	e.ScoreEvaluatedAt = persistshared.OptionTimeFromSql(scoreEvaluatedAt)
	e.FirstLovedAt = persistshared.OptionTimeFromSql(firstLovedAt)
	e.LastLovedAt = persistshared.OptionTimeFromSql(lastLovedAt)
	e.HallOfFameAt = persistshared.OptionTimeFromSql(hallOfFameAt)
	e.LastViewedAt = persistshared.OptionTimeFromSql(lastViewedAt)
	e.ViewMilestoneArchivedAt = persistshared.OptionTimeFromSql(viewMilestoneArchivedAt)
	return e, nil
}
