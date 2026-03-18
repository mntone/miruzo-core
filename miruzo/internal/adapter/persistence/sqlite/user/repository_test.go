package user_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestUserRepositoryUserSchemaRejectsInvalidDailyLoveUsed(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestUserSchemaRejectsInvalidDailyLoveUsed(t)
}

func TestUserRepositoryGetSingletonUser(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestGetSingletonUser(t)
}

func TestUserRepositoryIncrementDailyLoveUsedIncrements(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestIncrementDailyLoveUsedIncrements(t)
}

func TestUserRepositoryIncrementDailyLoveUsedReturnsQuotaExceededWhenMissing(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestIncrementDailyLoveUsedReturnsQuotaExceededWhenMissing(t)
}

func TestUserRepositoryIncrementDailyLoveUsedReturnsQuotaExceededWhenLimitReached(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestIncrementDailyLoveUsedReturnsQuotaExceededWhenLimitReached(t)
}

func TestUserRepositoryDecrementDailyLoveUsedDecrements(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestDecrementDailyLoveUsedDecrements(t)
}

func TestUserRepositoryDecrementDailyLoveUsedReturnsNotFound(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestDecrementDailyLoveUsedReturnsNotFound(t)
}

func TestUserRepositoryDecrementDailyLoveUsedReturnsQuotaUnderflow(t *testing.T) {
	testutilSQLite.NewUserSuite(t).RunTestDecrementDailyLoveUsedReturnsQuotaUnderflow(t)
}
