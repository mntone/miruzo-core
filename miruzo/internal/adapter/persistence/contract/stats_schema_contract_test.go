package contract_test

import (
	"fmt"
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

func TestStatsSchemaRejectsInvalidScore(t *testing.T) {
	tests := []struct {
		name  string
		score int32
	}{
		{
			name:  "score=-32769",
			score: -32769,
		},
		{
			name:  "score=32768",
			score: 32768,
		},
	}

	ingest := mb.Ingest().Build()
	stmt := "INSERT INTO stats(ingest_id, score, score_evaluated) VALUES(%s, %%d, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		dialectStmtFmt := fmt.Sprintf(stmt, h.ParamRange(1, 2)...)

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ops.MustAddIngest(t, ingest)
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(dialectStmtFmt, tt.score),
						ingest.ID,
						100,
					)
				})
			})
		}
	})
}

func TestStatsSchemaRejectsInvalidScoreEvaluated(t *testing.T) {
	tests := []struct {
		name             string
		scoreEvaluated   int32
		scoreEvaluatedAt time.Time
	}{
		{
			name:             "score_evaluated=-32769",
			scoreEvaluated:   -32769,
			scoreEvaluatedAt: mb.GetDefaultBaseTime().Add(20 * time.Minute),
		},
		{
			name:             "score_evaluated=32768",
			scoreEvaluated:   32768,
			scoreEvaluatedAt: mb.GetDefaultBaseTime().Add(40 * time.Minute),
		},
	}

	ingest := mb.Ingest().Build()
	stmt := "INSERT INTO stats(ingest_id, score, score_evaluated, score_evaluated_at) VALUES(%s, %s, %%d, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		dialectStmtFmt := fmt.Sprintf(stmt, h.ParamRange(1, 3)...)

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ops.MustAddIngest(t, ingest)
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(dialectStmtFmt, tt.scoreEvaluated),
						ingest.ID,
						100,
						tt.scoreEvaluatedAt,
					)
				})
			})
		}
	})
}

func TestStatsSchemaRejectsInvalidViewCount(t *testing.T) {
	tests := []struct {
		name      string
		viewCount string
	}{
		{
			name:      "view_count=-1",
			viewCount: "-1",
		},
		{
			name:      "view_count=9223372036854775808(math.MaxInt64+1)",
			viewCount: "9223372036854775808",
		},
	}

	ingest := mb.Ingest().Build()
	stmt := "INSERT INTO stats(ingest_id, score, score_evaluated, view_count) VALUES(%s, 100, 100, %%s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		dialectStmtFmt := fmt.Sprintf(stmt, h.Param(1))

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ops.MustAddIngest(t, ingest)
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(dialectStmtFmt, tt.viewCount),
						ingest.ID,
					)
				})
			})
		}
	})
}

func TestStatsSchemaRejectsInvalidViewMilestoneCount(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	tests := []struct {
		name                    string
		viewCount               int64
		viewMilestoneCount      int64
		viewMilestoneArchivedAt time.Time
	}{
		{
			name:                    "view_milestone_count=-100(<1)",
			viewCount:               1,
			viewMilestoneCount:      -100,
			viewMilestoneArchivedAt: baseTime.Add(35 * time.Minute),
		},
		{
			name:                    "view_milestone_count=200(>199)",
			viewCount:               199,
			viewMilestoneCount:      200,
			viewMilestoneArchivedAt: baseTime.Add(70 * time.Minute),
		},
	}

	ingest := mb.Ingest().Build()
	stmt := "INSERT INTO stats(ingest_id, score, score_evaluated, view_count, view_milestone_count, view_milestone_archived_at) VALUES(%s, 100, 100, %s, %%d, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		dialectStmtFmt := fmt.Sprintf(stmt, h.ParamRange(1, 3)...)

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ops.MustAddIngest(t, ingest)
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(dialectStmtFmt, tt.viewMilestoneCount),
						ingest.ID,
						tt.viewCount,
						tt.viewMilestoneArchivedAt,
					)
				})
			})
		}
	})
}

func TestStatsSchemaRejectsInvalidOccurredAt(t *testing.T) {
	tests := []struct {
		name string
		stmt string
	}{
		{
			name: "score_evaluated_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, score_evaluated_at) VALUES(%s, %s, %s, %s)",
		},
		{
			name: "first_loved_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, first_loved_at, last_loved_at) VALUES(%s, %s, %s, %s, %s)",
		},
		{
			name: "last_loved_at=-infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, first_loved_at, last_loved_at) VALUES(%s, %s, %s, %s, %s)",
		},
		{
			name: "hall_of_fame_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, hall_of_fame_at) VALUES(%s, %s, %s, %s)",
		},
		{
			name: "last_viewed_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, last_viewed_at, view_count) VALUES(%s, %s, %s, %s, %s)",
		},
	}

	ingest := mb.Ingest().Build()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RequireCapability(t, c.SupportsInfinityTimestamp)

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ops.MustAddIngest(t, ingest)
				t.Run(tt.name, func(t *testing.T) {
					args := []any{ingest.ID, 100, 100}
					switch tt.name {
					case "score_evaluated_at=infinity":
						args = append(args, "infinity")
					case "first_loved_at=infinity":
						args = append(args, "infinity", "infinity")
					case "last_loved_at=-infinity":
						args = append(args, "-infinity", "-infinity")
					case "hall_of_fame_at=infinity":
						args = append(args, "infinity")
					case "last_viewed_at=infinity":
						args = append(args, "infinity", 1)
					}

					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(tt.stmt, ops.ParamRange(1, int32(len(args)))...),
						args...,
					)
				})
			})
		}
	})
}
