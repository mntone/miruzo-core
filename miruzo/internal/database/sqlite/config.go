package sqlite

import (
	"net/url"

	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

type ConnectOptions struct {
	shared.ConnectionTuning
}

type ConnectConfig struct {
	dsn      *url.URL
	connTune shared.ConnectionTuning
}

func NewConnectConfigFromDSN(
	dsn string,
	options ConnectOptions,
) (*ConnectConfig, error) {
	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	queries := parsedDSN.Query()
	if !queries.Has("_txlock") {
		queries.Set("_txlock", "immediate")
	}
	if !queries.Has("_foreign_keys") {
		queries.Set("_foreign_keys", "1") // 1=ON
	}
	if !queries.Has("_journal_mode") && !queries.Has("_synchronous") {
		queries.Set("_journal_mode", "WAL")
	}
	parsedDSN.RawQuery = queries.Encode()

	return &ConnectConfig{
		dsn:      parsedDSN,
		connTune: options.ConnectionTuning,
	}, nil
}
