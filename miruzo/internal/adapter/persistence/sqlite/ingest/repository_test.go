package ingest_test

import (
	"testing"

	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

func TestIngestSchemaRejectsInvalidRelativePath(t *testing.T) {
	testutilSQLite.NewIngestSuite(t).RunTestIngestSchemaRejectsInvalidRelativePath(t)
}
