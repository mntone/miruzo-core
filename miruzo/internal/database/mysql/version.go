package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

var minMySQLForCheck = shared.Version{
	Major: 8,
	Minor: 0,
	Patch: 16,
}

func supportsMySQLCheckVersion(version shared.Version) bool {
	return !version.LessThan(minMySQLForCheck)
}

// verifyMySQLSupportsCheck rejects MariaDB first, then validates MySQL version.
// ParseVersion accepts suffixes, so ordering here is intentional.
func verifyMySQLSupportsCheck(ctx context.Context, db *sql.DB) error {
	var raw string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&raw); err != nil {
		return fmt.Errorf("read mysql version: %w", err)
	}

	if strings.Contains(raw, "MariaDB") {
		return fmt.Errorf("MariaDB is not supported (detected version: %s)", raw)
	}

	version, err := shared.ParseVersion(raw)
	if err != nil {
		return err
	}
	if !supportsMySQLCheckVersion(version) {
		return fmt.Errorf("mysql CHECK requires version >= %s (detected %s)", minMySQLForCheck, version)
	}

	return nil
}
