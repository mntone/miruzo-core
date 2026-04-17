package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

func verifySQLiteVersion(ctx context.Context, db *sql.DB, min shared.Version) error {
	var raw string
	if err := db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&raw); err != nil {
		return fmt.Errorf("read sqlite_version: %w", err)
	}

	version, err := shared.ParseVersion(raw)
	if err != nil {
		return err
	}
	if version.LessThan(min) {
		return fmt.Errorf("unsupported sqlite version: got %s, require >= %s", version, min)
	}
	return nil
}

var minSQLiteForReturningAndStrict = shared.Version{
	Major: 3,
	Minor: 37,
	Patch: 0,
}

func supportsSQLiteReturningAndStrictVersion(version shared.Version) bool {
	return !version.LessThan(minSQLiteForReturningAndStrict)
}

func verifySQLiteSupportsReturningAndStrict(ctx context.Context, db *sql.DB) error {
	var raw string
	if err := db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&raw); err != nil {
		return fmt.Errorf("read sqlite_version: %w", err)
	}

	version, err := shared.ParseVersion(raw)
	if err != nil {
		return err
	}
	if !supportsSQLiteReturningAndStrictVersion(version) {
		return fmt.Errorf("sqlite RETURNING & STRICT requires version >= %s (detected %s)", minSQLiteForReturningAndStrict, version)
	}

	return nil
}
