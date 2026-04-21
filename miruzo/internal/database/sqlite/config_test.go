package sqlite

import (
	"testing"
	"time"

	sharedDB "github.com/mntone/miruzo-core/miruzo/internal/database/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNewConnectConfigFromDSNAppliesOptionsAndPragmas(t *testing.T) {
	t.Parallel()

	options := ConnectOptions{
		ConnectionTuning: sharedDB.ConnectionTuning{
			ConnectionTimeout:     7 * time.Second,
			PoolWarmConnections:   2,
			MaxOpenConnections:    4,
			MaxConnectionIdleTime: 11 * time.Second,
			MaxConnectionLifeTime: 13 * time.Second,
		},
	}

	cfg, err := NewConnectConfigFromDSN("file:app.sqlite", options)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	assert.Equal(t, "dsn.Scheme", cfg.dsn.Scheme, "file")
	assert.Equal(t, "dsn.Opaque", cfg.dsn.Opaque, "app.sqlite")
	assert.Equal(t, "dsn.Query().Get(\"_txlock\")", cfg.dsn.Query().Get("_txlock"), "immediate")
	assert.Equal(t, "dsn.Query().Get(\"_foreign_keys\")", cfg.dsn.Query().Get("_foreign_keys"), "1")
	assert.Equal(t, "dsn.Query().Get(\"_journal_mode\")", cfg.dsn.Query().Get("_journal_mode"), "WAL")
	if cfg.connTune != options.ConnectionTuning {
		t.Fatalf("connTune = %#v, want %#v", cfg.connTune, options.ConnectionTuning)
	}
}

func TestNewConnectConfigFromDSNRespectsExistingPragmaSettings(t *testing.T) {
	t.Parallel()

	cfg, err := NewConnectConfigFromDSN(
		"file:app.sqlite?_txlock=exclusive&_foreign_keys=0&_synchronous=OFF",
		ConnectOptions{},
	)
	assert.NilError(t, "NewConnectConfigFromDSN() error", err)

	assert.Equal(t, "dsn.Query().Get(\"_txlock\")", cfg.dsn.Query().Get("_txlock"), "exclusive")
	assert.Equal(t, "dsn.Query().Get(\"_foreign_keys\")", cfg.dsn.Query().Get("_foreign_keys"), "0")
	assert.Equal(t, "dsn.Query().Get(\"_synchronous\")", cfg.dsn.Query().Get("_synchronous"), "OFF")
	assert.Equal(t, "dsn.Query().Get(\"_journal_mode\")", cfg.dsn.Query().Get("_journal_mode"), "")
}

func TestNewConnectConfigFromDSNReturnsErrorOnInvalidDSN(t *testing.T) {
	t.Parallel()

	_, err := NewConnectConfigFromDSN("file:app.sqlite\n", ConnectOptions{})
	assert.Error(t, "NewConnectConfigFromDSN() error", err)
}
