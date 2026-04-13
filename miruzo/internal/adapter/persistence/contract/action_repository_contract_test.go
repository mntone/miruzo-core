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

func TestActionSchemaRejectsDuplicateDecayPerPeriod(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	periodStartAt := actionResolver.PeriodStart(baseTime)
	stmt := "INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(%s, %s, %s, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			dialectStmt := fmt.Sprintf(stmt, ops.ParamRange(1, 4)...)
			ops.MustExec(
				t,
				dialectStmt,
				ingest.ID,
				model.ActionTypeDecay,
				baseTime,
				periodStartAt,
			)
			ops.AssertExecErrorIs(
				t,
				c.DBErrorMappingDefault,
				persist.ErrUniqueViolation,
				dialectStmt,
				ingest.ID,
				model.ActionTypeDecay,
				baseTime.Add(time.Second),
				periodStartAt,
			)
		})
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

func TestActionRepositoryCreateReturnsConflictOnDuplicateDecayPerPeriod(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	periodStartAt := actionResolver.PeriodStart(baseTime)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())

			_, err := ops.Action().Create(
				t.Context(),
				ingest.ID,
				model.ActionTypeDecay,
				baseTime,
				periodStartAt,
			)
			assert.NilError(t, "Create() first error", err)

			_, err = ops.Action().Create(
				t.Context(),
				ingest.ID,
				model.ActionTypeDecay,
				baseTime.Add(time.Second),
				periodStartAt,
			)
			assert.ErrorIs(t, "Create() second error", err, persist.ErrUniqueViolation)
		})
	})
}

func TestActionRepositoryCreateAllowsDuplicatePeriodForNonDecay(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	periodStartAt := actionResolver.PeriodStart(baseTime)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())

			_, err := ops.Action().Create(
				t.Context(),
				ingest.ID,
				model.ActionTypeView,
				baseTime,
				periodStartAt,
			)
			assert.NilError(t, "Create() first error", err)

			_, err = ops.Action().Create(
				t.Context(),
				ingest.ID,
				model.ActionTypeView,
				baseTime.Add(time.Second),
				periodStartAt,
			)
			assert.NilError(t, "Create() second error", err)
		})
	})
}

// --- CreateDailyDecayIfAbsent ---

func TestActionRepositoryCreateDailyDecayIfAbsentReturnsConflictOnDuplicatePeriod(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	periodStartAt := actionResolver.PeriodStart(baseTime)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())

			err := ops.Action().CreateDailyDecayIfAbsent(
				t.Context(),
				ingest.ID,
				baseTime,
				periodStartAt,
			)
			assert.NilError(t, "CreateDailyDecayIfAbsent() first error", err)

			err = ops.Action().CreateDailyDecayIfAbsent(
				t.Context(),
				ingest.ID,
				baseTime.Add(time.Second),
				periodStartAt,
			)
			assert.ErrorIs(t, "CreateDailyDecayIfAbsent() second error", err, persist.ErrConflict)
		})
	})
}

func TestActionRepositoryCreateDailyDecayIfAbsentAllowsDuplicatePeriodForNonDecay(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	periodStartAt := actionResolver.PeriodStart(baseTime)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			_, err := ops.Action().Create(
				t.Context(),
				ingest.ID,
				model.ActionTypeView,
				baseTime,
				periodStartAt,
			)
			assert.NilError(t, "Create() error", err)

			err = ops.Action().CreateDailyDecayIfAbsent(
				t.Context(),
				ingest.ID,
				baseTime.Add(time.Second),
				periodStartAt,
			)
			assert.NilError(t, "CreateDailyDecayIfAbsent() error", err)
		})
	})
}

func TestActionRepositoryCreateDailyDecayIfAbsentAllowsDifferentPeriod(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	periodStartAt := actionResolver.PeriodStart(baseTime)
	nextPeriodStartAt := periodStartAt.Add(24 * time.Hour)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())

			err := ops.Action().CreateDailyDecayIfAbsent(
				t.Context(),
				ingest.ID,
				baseTime,
				periodStartAt,
			)
			assert.NilError(t, "CreateDailyDecayIfAbsent() first error", err)

			err = ops.Action().CreateDailyDecayIfAbsent(
				t.Context(),
				ingest.ID,
				baseTime.Add(24*time.Hour),
				nextPeriodStartAt,
			)
			assert.NilError(t, "CreateDailyDecayIfAbsent() second error", err)
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
			name: "MultiplePeriods",
			offsets: []time.Duration{
				0,
				24 * time.Hour,
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
