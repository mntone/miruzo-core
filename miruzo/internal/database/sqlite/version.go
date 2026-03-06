package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

type sqliteVersion struct {
	major, minor, patch int
}

func (v sqliteVersion) LessThan(min sqliteVersion) bool {
	if v.major != min.major {
		return v.major < min.major
	}
	if v.minor != min.minor {
		return v.minor < min.minor
	}
	return v.patch < min.patch
}

func (v sqliteVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func parseSQLiteVersion(s string) (sqliteVersion, error) {
	var v sqliteVersion
	if _, err := fmt.Sscanf(s, "%d.%d.%d", &v.major, &v.minor, &v.patch); err != nil {
		return sqliteVersion{}, fmt.Errorf("invalid sqlite_version %q: %w", s, err)
	}
	return v, nil
}

func verifySQLiteVersion(ctx context.Context, db *sql.DB, min sqliteVersion) error {
	var raw string
	if err := db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&raw); err != nil {
		return fmt.Errorf("read sqlite_version: %w", err)
	}

	version, err := parseSQLiteVersion(raw)
	if err != nil {
		return err
	}
	if version.LessThan(min) {
		return fmt.Errorf("unsupported sqlite version: got %s, require >= %s", version, min)
	}
	return nil
}

var minSQLiteForReturning = sqliteVersion{
	major: 3,
	minor: 35,
	patch: 0,
}

func supportsSQLiteReturningVersion(version sqliteVersion) bool {
	return !version.LessThan(minSQLiteForReturning)
}

func verifySQLiteSupportsReturning(ctx context.Context, db *sql.DB) error {
	var raw string
	if err := db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&raw); err != nil {
		return fmt.Errorf("read sqlite_version: %w", err)
	}

	version, err := parseSQLiteVersion(raw)
	if err != nil {
		return err
	}
	if !supportsSQLiteReturningVersion(version) {
		return fmt.Errorf("sqlite RETURNING requires version >= %s (detected %s)", minSQLiteForReturning, version)
	}

	return nil
}
