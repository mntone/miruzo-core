package contract

import (
	"fmt"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type TransactionOperations interface {
	Exec(t testing.TB, stmt string, args ...any) error
	ExecInsertAndGetID(t testing.TB, stmt string, args ...any) (int64, error)
	ExecReturningInt64(t testing.TB, stmt string, args ...any) (int64, error)
	ExecAndGetRowCount(t testing.TB, stmt string, args ...any) (int64, error)
	Rollback(t testing.TB)

	InsertImage(t testing.TB, e persist.Image) error
	SelectStats(t testing.TB, id model.IngestIDType) (model.Stats, error)
}

type RepositoryProvider interface {
	Action() persist.ActionRepository
	ImageList() persist.ImageListRepository
	Job() persist.JobRepository
	Settings() persist.SettingsRepository
	Stats() persist.StatsRepository
	User() persist.SessionUserRepository
	View() persist.ViewRepository
}

type TxSession struct {
	Dialect
	TransactionOperations
	RepositoryProvider
}

// --- operation helpers ---

func (s TxSession) MustExec(t testing.TB, stmt string, args ...any) {
	t.Helper()
	err := s.Exec(t, stmt, args...)
	assert.NilError(t, "Exec() error", err)
}

func (s TxSession) MustExecInsertAndGetID(t testing.TB, stmt string, args ...any) int64 {
	t.Helper()
	lastInsertID, err := s.ExecInsertAndGetID(t, stmt, args...)
	assert.NilError(t, "ExecInsertAndGetID() error", err)
	return lastInsertID
}

func (s TxSession) MustExecReturningInt64(t testing.TB, stmt string, args ...any) int64 {
	t.Helper()
	retID, err := s.ExecReturningInt64(t, stmt, args...)
	assert.NilError(t, "ExecReturningInt64() error", err)
	return retID
}

func (s TxSession) MustExecAndGetRowCount(t testing.TB, stmt string, args ...any) int64 {
	t.Helper()
	rowCount, err := s.ExecAndGetRowCount(t, stmt, args...)
	assert.NilError(t, "ExecAndGetRowCount() error", err)
	return rowCount
}

func (s TxSession) AssertExecErrorIs(
	t *testing.T,
	errorMapping DBErrorMapping,
	wantError error,
	stmt string,
	args ...any,
) {
	t.Helper()
	err := s.ToPersistError("Exec()", s.Exec(t, stmt, args...), errorMapping)
	assert.ErrorIs(t, "Exec() error", err, wantError)
}

// --- add fixture ---

func (s TxSession) MustAddIngest(t testing.TB, e model.Ingest) model.Ingest {
	t.Helper()

	stmt := fmt.Sprintf(
		"INSERT INTO ingests(id, process, visibility, relative_path, fingerprint, ingested_at, captured_at, updated_at) VALUES(%s, %s, %s, %s, %s, %s, %s, %s)",
		s.ParamRange(1, 8)...,
	)
	s.MustExec(
		t, stmt,
		e.ID,
		e.Process, e.Visibility,
		e.RelativePath, e.Fingerprint,
		e.IngestedAt, e.CapturedAt, e.UpdatedAt,
	)
	return e
}

func (s TxSession) MustAddImage(t testing.TB, e persist.Image) persist.Image {
	t.Helper()

	err := s.InsertImage(t, e)
	assert.NilError(t, "InsertImage() error", err)
	return e
}

func (s TxSession) MustAddStats(t testing.TB, e model.Stats) model.Stats {
	t.Helper()

	stmt := fmt.Sprintf(
		"INSERT INTO stats(ingest_id, score, score_evaluated, first_loved_at, last_loved_at, hall_of_fame_at, last_viewed_at, view_count) VALUES(%s, %s, %s, %s, %s, %s, %s, %s)",
		s.ParamRange(1, 8)...,
	)
	s.MustExec(
		t, stmt,
		e.IngestID,
		e.Score, e.ScoreEvaluated,
		e.FirstLovedAt.ToPointer(),
		e.LastLovedAt.ToPointer(),
		e.HallOfFameAt.ToPointer(),
		e.LastViewedAt.ToPointer(),
		e.ViewCount,
	)
	return e
}

func (s TxSession) MustAddAction(
	t testing.TB,
	ingestID model.IngestIDType,
	kind model.ActionType,
	at time.Time,
	periodStartAt time.Time,
) model.ActionIDType {
	t.Helper()

	actionID, err := s.Action().Create(t.Context(), ingestID, kind, at, periodStartAt)
	assert.NilError(t, "MustAddAction() error", err)
	return actionID
}

func (s TxSession) MustRemoveUser(t *testing.T) {
	t.Helper()
	rowCount := s.MustExecAndGetRowCount(t, "DELETE FROM users WHERE id=1")
	assert.Equal(t, "MustRemoveUser() row_count", rowCount, 1)
}

func (s TxSession) MustSetDailyLoveUsed(t *testing.T, dailyLoveUsed model.QuotaInt) {
	t.Helper()
	rowCount := s.MustExecAndGetRowCount(
		t,
		fmt.Sprintf("UPDATE users SET daily_love_used=%s WHERE id=1", s.Param(1)),
		dailyLoveUsed,
	)
	assert.Equal(t, "MustSetDailyLoveUsed() row_count", rowCount, 1)
}

// --- get each row ---

func (s TxSession) MustGetStats(t testing.TB, id model.IngestIDType) model.Stats {
	t.Helper()
	stats, err := s.SelectStats(t, id)
	assert.NilError(t, "MustGetStats() error", err)
	return stats
}
