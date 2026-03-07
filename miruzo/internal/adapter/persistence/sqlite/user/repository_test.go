package user_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestUserRepositoryGetSingletonUser(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestGetSingletonUser(t)
}
