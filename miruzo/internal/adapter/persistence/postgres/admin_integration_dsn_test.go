package postgres_test

import (
	"fmt"
	"net/url"
	"path"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func withPostgresDatabaseName(dsn string, databaseName string) (string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}
	if parsedURL.Scheme != "postgres" && parsedURL.Scheme != "postgresql" {
		return "", fmt.Errorf("unsupported postgres dsn scheme: %s", parsedURL.Scheme)
	}

	parsedURL.Path = path.Join("/", databaseName)
	return parsedURL.String(), nil
}

func TestWithPostgresDatabaseName(t *testing.T) {
	t.Parallel()

	const srcDSN = "postgresql://app:secret@localhost:5432/miruzo_test?sslmode=disable"
	gotDSN, err := withPostgresDatabaseName(srcDSN, "miruzo_tmp")
	assert.NilError(t, "withPostgresDatabaseName() error", err)

	gotCfg, err := pgxpool.ParseConfig(gotDSN)
	assert.NilError(t, "ParseConfig(gotDSN) error", err)
	assert.Equal(t, "database", gotCfg.ConnConfig.Database, "miruzo_tmp")
	assert.Equal(t, "user", gotCfg.ConnConfig.User, "app")
	assert.Equal(t, "host", gotCfg.ConnConfig.Host, "localhost")

	gotURL, err := url.Parse(gotDSN)
	assert.NilError(t, "url.Parse(gotDSN) error", err)
	assert.Equal(t, "sslmode", gotURL.Query().Get("sslmode"), "disable")
}

func TestWithPostgresDatabaseNameSupportsPostgresScheme(t *testing.T) {
	t.Parallel()

	const srcDSN = "postgres://app:secret@localhost:5432/miruzo_test?sslmode=disable"
	gotDSN, err := withPostgresDatabaseName(srcDSN, "miruzo_tmp")
	assert.NilError(t, "withPostgresDatabaseName() error", err)

	gotCfg, err := pgxpool.ParseConfig(gotDSN)
	assert.NilError(t, "ParseConfig(gotDSN) error", err)
	assert.Equal(t, "database", gotCfg.ConnConfig.Database, "miruzo_tmp")

	gotURL, err := url.Parse(gotDSN)
	assert.NilError(t, "url.Parse(gotDSN) error", err)
	assert.Equal(t, "sslmode", gotURL.Query().Get("sslmode"), "disable")
}

func TestWithPostgresDatabaseNameReturnsErrorOnUnsupportedScheme(t *testing.T) {
	t.Parallel()

	_, err := withPostgresDatabaseName("mysql://root@localhost/miruzo", "miruzo_tmp")
	assert.Error(t, "withPostgresDatabaseName() error", err)
}
