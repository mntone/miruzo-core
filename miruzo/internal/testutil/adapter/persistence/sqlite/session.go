package sqlite

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	database "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	dbshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/user"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

type sqliteTxSession struct {
	tx      *sql.Tx
	queries *gen.Queries
	nextID  model.IngestIDType
}

func newTxSession(tx *sql.Tx) *sqliteTxSession {
	return &sqliteTxSession{
		tx:      tx,
		queries: gen.New(tx),
		nextID:  modelbuilder.GetNextID(),
	}
}

// --- Dialect ---

func (s sqliteTxSession) Backend() backend.Backend {
	return backend.SQLite
}

func (s sqliteTxSession) MapError(
	operation string,
	err error,
	mapping contract.DBErrorMapping,
) error {
	switch mapping {
	case contract.DBErrorMappingDefault:
		err = dbshared.MapSQLiteError(operation, err)
	case contract.DBErrorMappingDelete:
		err = dbshared.MapSQLiteDeleteError(operation, err)
	}
	return err
}

func (s sqliteTxSession) BindVarStyle() contract.BindVarStyle {
	return contract.BindVarStyleQuestion
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
		return dbshared.MapSQLiteError("InsertImage", err)
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

	e.ScoreEvaluatedAt = dbshared.OptionTimeFromSql(scoreEvaluatedAt)
	e.FirstLovedAt = dbshared.OptionTimeFromSql(firstLovedAt)
	e.LastLovedAt = dbshared.OptionTimeFromSql(lastLovedAt)
	e.HallOfFameAt = dbshared.OptionTimeFromSql(hallOfFameAt)
	e.LastViewedAt = dbshared.OptionTimeFromSql(lastViewedAt)
	e.ViewMilestoneArchivedAt = dbshared.OptionTimeFromSql(viewMilestoneArchivedAt)
	return e, nil
}

// --- RepositoryProvider ---

func (s sqliteTxSession) Action() persist.ActionRepository {
	return action.NewRepository(s.queries)
}

func (s sqliteTxSession) ImageList() persist.ImageListRepository {
	return imagelist.NewRepository(s.queries)
}

func (s sqliteTxSession) Job() persist.JobRepository {
	return database.NewJobRepository(s.queries)
}

func (s sqliteTxSession) Settings() persist.SettingsRepository {
	return database.NewSettingsRepository(s.queries)
}

func (s sqliteTxSession) Stats() persist.StatsRepository {
	return stats.NewRepository(s.queries)
}

func (s sqliteTxSession) User() persist.SessionUserRepository {
	return user.NewRepository(s.queries)
}

func (s sqliteTxSession) View() persist.ViewRepository {
	return database.NewViewRepository(s.queries)
}
