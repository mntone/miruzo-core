package persistence

import (
	"fmt"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

var actionSuiteBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

type ActionSuite SuiteBase[persist.ActionRepository]

func (ste ActionSuite) RunTestActionSchemaRejectsInvalidKind(t *testing.T) {
	t.Helper()

	tests := []struct {
		name    string
		kind    int32
		wantErr error
	}{
		{
			name:    "kind=2",
			kind:    2,
			wantErr: persist.ErrCheckViolation,
		},
		{
			name:    "kind=10",
			kind:    10,
			wantErr: persist.ErrCheckViolation,
		},
		{
			name:    "kind=17",
			kind:    17,
			wantErr: persist.ErrCheckViolation,
		},
		{
			name:    "kind=999",
			kind:    999,
			wantErr: persist.ErrCheckViolation,
		},
	}

	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, actionSuiteBaseTimeUTC))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := fmt.Sprintf(
				"INSERT INTO actions(ingest_id, kind, occurred_at) VALUES(%d, %d, '%s')",
				ingest.ID,
				tt.kind,
				actionSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			)
			err := ste.Operations.ExecuteStatement(stmt)
			assert.ErrorIs(t, "insert error", err, tt.wantErr)
		})
	}
}

// PostgreSQL only
func (ste ActionSuite) RunTestActionSchemaRejectsInvalidOccurredAt(t *testing.T) {
	t.Helper()

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

	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, actionSuiteBaseTimeUTC))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := fmt.Sprintf(
				"INSERT INTO actions(ingest_id, occurred_at) VALUES(%d, '%s')",
				ingest.ID, tt.occurredAt,
			)
			err := ste.Operations.ExecuteStatement(stmt)
			assert.ErrorIs(t, "insert error", err, tt.wantErr)
		})
	}
}

func (ste ActionSuite) RunTestCreateAction(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, actionSuiteBaseTimeUTC))

	actionID, err := ste.Repository.Create(
		ste.Context,
		ingest.ID,
		model.ActionTypeView,
		suiteBaseTimeUTC.Add(2*time.Hour),
	)
	assert.NilError(t, "CreateAction() error", err)
	assert.Equal(t, "actionID", actionID, 1)
}

type testExistsSinceReturnsFalseActions struct {
	overrideIngestID mo.Option[model.IngestIDType]
	overrideType     mo.Option[model.ActionType]
	offset           time.Duration
}

func (ste ActionSuite) RunTestExistsSinceReturnsFalse(t *testing.T) {
	t.Helper()

	overrideIngestID := model.IngestIDType(2)
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
					overrideType: mo.Some(model.ActionTypeView),
					offset:       time.Hour,
				},
			},
		},
		{
			name: "IngestDoesNotMatch",
			actions: []testExistsSinceReturnsFalseActions{
				{
					overrideIngestID: mo.Some(overrideIngestID),
					offset:           time.Hour,
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

	baseTime := mb.GetDefaultStatsBaseTime()
	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, baseTime))
	ste.Operations.MustAddIngest(t, NewIngestFixture(overrideIngestID, baseTime))

	for _, tt := range tests {
		ste.Operations.MustTruncateActions(t)

		t.Run(tt.name, func(t *testing.T) {
			for _, action := range tt.actions {
				ste.Operations.MustAddAction(
					t,
					action.overrideIngestID.OrElse(ingest.ID),
					action.overrideType.OrElse(model.ActionTypeDecay),
					baseTime.Add(24*time.Hour+action.offset),
				)
			}

			exists, err := ste.Repository.ExistsSince(
				ste.Context,
				ingest.ID,
				model.ActionTypeDecay,
				baseTime.Add(24*time.Hour),
			)
			assert.NilError(t, "ExistsSince() error", err)
			assert.Equal(t, "exists", exists, false)
		})
	}
}

func (ste ActionSuite) RunTestExistsSinceReturnsTrue(t *testing.T) {
	t.Helper()

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

	baseTime := mb.GetDefaultStatsBaseTime()
	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, baseTime))

	for _, tt := range tests {
		ste.Operations.MustTruncateActions(t)

		t.Run(tt.name, func(t *testing.T) {
			for _, offset := range tt.offsets {
				ste.Operations.MustAddAction(
					t,
					ingest.ID,
					model.ActionTypeDecay,
					baseTime.Add(24*time.Hour+offset),
				)
			}

			exists, err := ste.Repository.ExistsSince(
				ste.Context,
				ingest.ID,
				model.ActionTypeDecay,
				baseTime.Add(24*time.Hour),
			)
			assert.NilError(t, "ExistsSince() error", err)
			assert.Equal(t, "exists", exists, true)
		})
	}
}
