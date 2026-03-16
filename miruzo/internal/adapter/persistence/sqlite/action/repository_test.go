package action_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestActionRepositoryActionSchemaRejectsInvalidKind(t *testing.T) {
	testutilSQLite.NewActionSuite(t).RunTestActionSchemaRejectsInvalidKind(t)
}

func TestActionRepositoryCreateAction(t *testing.T) {
	testutilSQLite.NewActionSuite(t).RunTestCreateAction(t)
}
