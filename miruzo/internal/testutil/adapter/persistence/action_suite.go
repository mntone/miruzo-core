package persistence

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

var actionSuiteBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

type ActionSuite SuiteBase[persist.ActionRepository]

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
