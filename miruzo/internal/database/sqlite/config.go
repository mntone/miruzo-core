package sqlite

import "github.com/mntone/miruzo-core/miruzo/internal/database/shared"

type ConnectConfig struct {
	DSN string
	shared.ConnectionTuning
}
