package persistence

import (
	"fmt"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
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
