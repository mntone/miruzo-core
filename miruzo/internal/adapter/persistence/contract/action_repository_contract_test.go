package contract_test

import (
	"fmt"
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

// --- Schema ---

func TestActionRepositorySchemaRejectsInvalidKind(t *testing.T) {
	tests := []int32{-1, 2, 10, 17, 999}

	stmt := "INSERT INTO actions(ingest_id, kind, occurred_at) VALUES(%s, %s, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			ops := h.BeginTx(t)
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			t.Run(fmt.Sprintf("kind=%d", tt), func(t *testing.T) {
				ops.AssertExecErrorIs(
					t,
					c.DBErrorMappingDefault,
					persist.ErrCheckViolation,
					fmt.Sprintf(stmt, ops.ParamRange(1, 3)...),
					ingest.ID,
					tt,
					mb.GetDefaultBaseTime(),
				)
			})
			ops.Rollback(t)
		}
	})
}

func TestActionRepositorySchemaRejectsInvalidOccurredAt(t *testing.T) {
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

	stmt := "INSERT INTO actions(ingest_id, occurred_at) VALUES(%s, %s)"

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
						fmt.Sprintf(stmt, ops.ParamRange(1, 2)...),
						ingest.ID,
						tt.occurredAt,
					)
				})
			})
		}
	})
}

// --- Create ---

func TestActionRepositoryCreates(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			actionID, err := ops.Action().Create(
				t.Context(),
				ingest.ID,
				model.ActionTypeView,
				mb.GetDefaultBaseTime(),
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
						ops.MustAddAction(
							t, ingestID,
							action.overrideType.OrElse(model.ActionTypeDecay),
							mb.GetDefaultBaseTime().Add(24*time.Hour+action.offset),
						)
					}

					exists, err := ops.Action().ExistsSince(
						t.Context(),
						ingest.ID,
						model.ActionTypeDecay,
						mb.GetDefaultBaseTime().Add(24*time.Hour),
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

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				ingest := ops.MustAddIngest(t, mb.Ingest().Build())
				t.Run(tt.name, func(t *testing.T) {
					for _, offset := range tt.offsets {
						ops.MustAddAction(
							t,
							ingest.ID,
							model.ActionTypeDecay,
							mb.GetDefaultBaseTime().Add(24*time.Hour+offset),
						)
					}

					exists, err := ops.Action().ExistsSince(
						t.Context(),
						ingest.ID,
						model.ActionTypeDecay,
						mb.GetDefaultBaseTime().Add(24*time.Hour),
					)
					assert.NilError(t, "ExistsSince() error", err)
					assert.Equal(t, "exists", exists, true)
				})
			})
		}
	})
}
