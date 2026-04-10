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

	stmt := "INSERT INTO stats(ingest_id, score, score_evaluated) VALUES(%s, %d, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(stmt, ops.Param(1), tt.score, ops.Param(2)),
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

	stmt := "INSERT INTO stats(ingest_id, score, score_evaluated, score_evaluated_at) VALUES(%s, %s, %d, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(stmt, ops.Param(1), ops.Param(2), tt.scoreEvaluated, ops.Param(3)),
						ingest.ID,
						100,
						tt.scoreEvaluatedAt,
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

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RequireCapability(t, c.SupportsInfinityTimestamp)

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
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

					params := make([]any, len(args))
					for i := range args {
						params[i] = ops.Param(int32(i + 1))
					}
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(tt.stmt, params...),
						args...,
					)
				})
			})
		}
	})
}
