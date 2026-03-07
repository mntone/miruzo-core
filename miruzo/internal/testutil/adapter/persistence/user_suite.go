package persistence

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type UserSuite SuiteBase[persist.UserRepository]

func (ste UserSuite) RunTestGetSingletonUser(t *testing.T) {
	t.Helper()

	user, err := ste.Repository.GetSingletonUser(ste.Context)
	assert.NilError(t, "GetSingletonUser() error", err)
	assert.Equal(t, "user.ID", user.ID, 1)
	assert.Equal(t, "user.DailyLoveUsed", user.DailyLoveUsed, 0)
}
