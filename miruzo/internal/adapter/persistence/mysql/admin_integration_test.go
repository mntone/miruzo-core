//go:build integration

package mysql_test

import (
	"context"
	"testing"

	adapterpersistence "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	testutilpersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
	testutilmysql "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestMySQLAdminHandleLifecycle(t *testing.T) {
	registry := &testutil.CleanupRegistry{}
	t.Cleanup(func() {
		assert.NilError(t, "cleanupRegistry.CloseAll()", registry.CloseAll())
	})

	dbConfig := config.DefaultDatabaseConfig()
	dbConfig.Backend = backend.MySQL
	dbConfig.DSN = testutilmysql.GetMySQLTestDSN(t, registry)

	appDatabaseName := testutilpersistence.NewAdminTestDatabaseName("adm_mysql")
	appConfigDSN, err := withMySQLDatabaseName(dbConfig.DSN, appDatabaseName)
	assert.NilError(t, "withMySQLDatabaseName()", err)
	dbConfig.DSN = appConfigDSN

	testutilpersistence.RunAdminHandleLifecycle(
		t,
		func(
			ctx context.Context,
		) (adapterpersistence.DatabaseAdminHandle, error) {
			return mysql.OpenAdminHandle(ctx, dbConfig, "")
		},
	)
}
