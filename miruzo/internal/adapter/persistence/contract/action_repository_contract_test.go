package contract_test

import (
	"fmt"
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

const actionDayStartOffset = 5 * time.Hour

var actionResolver = period.NewDailyResolver(actionDayStartOffset)

// --- Schema ---

func TestActionSchemaRejectsInvalidKind(t *testing.T) {
	tests := []int32{-1, 2, 10, 17, 999}

	baseTime := mb.GetDefaultBaseTime()
	stmt := "INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(%s, %s, %s, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			ops := h.BeginTx(t)
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			t.Run(fmt.Sprintf("kind=%d", tt), func(t *testing.T) {
				ops.AssertExecErrorIs(
					t,
					c.DBErrorMappingDefault,
					persist.ErrCheckViolation,
					fmt.Sprintf(stmt, ops.ParamRange(1, 4)...),
					ingest.ID,
					tt,
					baseTime,
					baseTime,
				)
			})
			ops.Rollback(t)
		}
	})
}

func TestActionSchemaRejectsInvalidOccurredAt(t *testing.T) {
	tests := []struct {
		name       string
		occurredAt string
		wantErr    error
	}{
		{
			name:       "occurred_at=infinity",
			occurredAt: "infinity",
			wantErr:    persist.ErrCheckViolation,
		},
		{
			name:       "occurred_at=-infinity",
			occurredAt: "-infinity",
			wantErr:    persist.ErrCheckViolation,
		},
	}

	baseTime := mb.GetDefaultBaseTime()
	stmt := "INSERT INTO actions(ingest_id, occurred_at, period_start_at) VALUES(%s, %s, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RequireCapability(t, c.SupportsInfinityTimestamp)

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						tt.wantErr,
						fmt.Sprintf(stmt, ops.ParamRange(1, 3)...),
						ingest.ID,
						tt.occurredAt,
						baseTime,
					)
				})
			})
		}
	})
}

func TestActionSchemaRejectsInvalidPeriodStartAt(t *testing.T) {
	tests := []struct {
		name          string
		periodStartAt string
		wantErr       error
	}{
		{
			name:          "period_start_at=infinity",
			periodStartAt: "infinity",
			wantErr:       persist.ErrCheckViolation,
		},
		{
			name:          "period_start_at=-infinity",
			periodStartAt: "-infinity",
			wantErr:       persist.ErrCheckViolation,
		},
	}

	baseTime := mb.GetDefaultBaseTime()
	stmt := "INSERT INTO actions(ingest_id, occurred_at, period_start_at) VALUES(%s, %s, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RequireCapability(t, c.SupportsInfinityTimestamp)

		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						tt.wantErr,
						fmt.Sprintf(stmt, ops.ParamRange(1, 3)...),
						ingest.ID,
						baseTime,
						tt.periodStartAt,
					)
				})
			})
		}
	})
}

// --- Create ---

func TestActionRepositoryCreates(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			actionID, err := ops.Action().Create(
				t.Context(),
				ingest.ID,
				model.ActionTypeView,
				baseTime,
				baseTime,
			)
			assert.NilError(t, "Create() error", err)
			assert.GreaterThan(t, "actionID", actionID, 0)
		})
	})
}

// --- ExistsSince ---

func TestActionRepositoryExistsSinceReturnsFalse(t *testing.T) {
	type testExistsSinceReturnsFalseActions struct {
		offset       time.Duration
		nextIngestID bool
		overrideType mo.Option[model.ActionType]
	}

	tests := []struct {
		name    string
		actions []testExistsSinceReturnsFalseActions
	}{
		{
			name:    "Missing",
			actions: []testExistsSinceReturnsFalseActions{},
		},
		{
			name: "KindDoesNotMatch",
			actions: []testExistsSinceReturnsFalseActions{
				{
					offset:       time.Hour,
					overrideType: mo.Some(model.ActionTypeView),
				},
			},
		},
		// verify filtering by ingest_id, not missing-ingest behavior
		{
			name: "IngestDoesNotMatch",
			actions: []testExistsSinceReturnsFalseActions{
				{
					offset:       time.Hour,
					nextIngestID: true,
				},
			},
		},
		{
			name: "OccurredAtBeforeSince",
			actions: []testExistsSinceReturnsFalseActions{
				{
					offset: -time.Second,
				},
			},
		},
	}

	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
				nextIngest := ops.MustAddIngest(t, mb.Ingest().Build())
				t.Run(tt.name, func(t *testing.T) {
					for _, action := range tt.actions {
						var ingestID model.IngestIDType
						if action.nextIngestID {
							ingestID = nextIngest.ID
						} else {
							ingestID = ingest.ID
						}
						at := baseTime.Add(24*time.Hour + action.offset)
						ops.MustAddAction(
							t, ingestID,
							action.overrideType.OrElse(model.ActionTypeDecay),
							at,
							actionResolver.PeriodStart(at),
						)
					}

					exists, err := ops.Action().ExistsSince(
						t.Context(),
						ingest.ID,
						model.ActionTypeDecay,
						baseTime.Add(24*time.Hour),
					)
					assert.NilError(t, "ExistsSince() error", err)
					assert.Equal(t, "exists", exists, false)
				})
			})
		}
	})
}

func TestActionRepositoryExistsSinceReturnsTrue(t *testing.T) {
	tests := []struct {
		name    string
		offsets []time.Duration
	}{
		{
			name: "Unique",
			offsets: []time.Duration{
				0,
			},
		},
		{
			name: "Multiple",
			offsets: []time.Duration{
				0,
				time.Hour,
			},
		},
		{
			name: "MultipleOnTie",
			offsets: []time.Duration{
				0,
				0,
			},
		},
	}

	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
				t.Run(tt.name, func(t *testing.T) {
					for _, offset := range tt.offsets {
						at := baseTime.Add(24*time.Hour + offset)
						ops.MustAddAction(
							t,
							ingest.ID,
							model.ActionTypeDecay,
							at,
							actionResolver.PeriodStart(at),
						)
					}

					exists, err := ops.Action().ExistsSince(
						t.Context(),
						ingest.ID,
						model.ActionTypeDecay,
						baseTime.Add(24*time.Hour),
					)
					assert.NilError(t, "ExistsSince() error", err)
					assert.Equal(t, "exists", exists, true)
				})
			})
		}
	})
}
