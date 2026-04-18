package contract_test

import (
	"os"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	testutilMySQL "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/mysql"
	testutilPostgres "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/postgres"
	testutilSQLite "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/sqlite"
)

var cleanupRegistry = &testutil.CleanupRegistry{}

func toHarnesses(t *testing.T, backends []backend.Backend) []contract.Harness {
	var harnesses []contract.Harness = make([]contract.Harness, 0, len(backends))
	for _, b := range backends {
		switch b {
		case backend.MySQL:
			harnesses = append(harnesses, testutilMySQL.NewHarness(t, cleanupRegistry))
		case backend.PostgreSQL:
			harnesses = append(harnesses, testutilPostgres.NewHarness(t, cleanupRegistry))
		case backend.SQLite:
			harnesses = append(harnesses, testutilSQLite.NewHarness(t, cleanupRegistry))
		}
	}
	return harnesses
}

func runHarnesses(t *testing.T, callback contract.HarnessCallback) {
	t.Helper()
	contract.RunDefaultHarnesses(t, toHarnesses, callback)
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	if err := cleanupRegistry.CloseAll(); err != nil && exitCode == 0 {
		exitCode = 1
	}
	os.Exit(exitCode)
}
