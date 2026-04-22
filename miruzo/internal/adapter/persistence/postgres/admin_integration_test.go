//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	testutilpersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
	testutilpostgres "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestPostgresAdminHandleLifecycle(t *testing.T) {
	registry := &testutil.CleanupRegistry{}
	t.Cleanup(func() {
		assert.NilError(t, "cleanupRegistry.CloseAll()", registry.CloseAll())
	})

	dbConfig := config.DefaultDatabaseConfig()
	dbConfig.Backend = backend.PostgreSQL
	dbConfig.DSN = testutilpostgres.GetPostgresTestDSN(t, registry)

	appDatabaseName := testutilpersistence.NewAdminTestDatabaseName("adm_postgres")
	appConfigDSN, err := withPostgresDatabaseName(dbConfig.DSN, appDatabaseName)
	assert.NilError(t, "withPostgresDatabaseName()", err)
	dbConfig.DSN = appConfigDSN

	testutilpersistence.RunAdminHandleLifecycle(
		t,
		func(
			ctx context.Context,
		) (persistence.DatabaseAdminHandle, error) {
			return postgres.OpenAdminHandle(
				ctx,
				dbConfig,
				shared.DatabaseAdminOptions{},
			)
		},
	)
}
