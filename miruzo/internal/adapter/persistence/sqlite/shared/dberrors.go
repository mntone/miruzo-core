package shared

import (
	"context"
	"errors"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapSQLiteBaseError(operation string, err error, foreignKeyError error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf(
			"%w: operation=%s: %v",
			persist.ErrTimeout,
			operation,
			err,
		)
	}

	if sqliteErr, ok := errors.AsType[sqlite3.Error](err); ok {
		switch sqliteErr.Code {
		case sqlite3.ErrInterrupt:
			return fmt.Errorf(
				"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
				persist.ErrTimeout,
				operation,
				sqliteErr.Code,
				sqliteErr.ExtendedCode,
				err,
			)

		case sqlite3.ErrBusy:
			return fmt.Errorf(
				"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
				persist.ErrRecoverableConflict,
				operation,
				sqliteErr.Code,
				sqliteErr.ExtendedCode,
				err,
			)

		case sqlite3.ErrLocked:
			if sqliteErr.ExtendedCode == sqlite3.ErrLockedSharedCache {
				return fmt.Errorf(
					"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
					persist.ErrRecoverableConflict,
					operation,
					sqliteErr.Code,
					sqliteErr.ExtendedCode,
					err,
				)
			}

		case sqlite3.ErrNomem,
			sqlite3.ErrIoErr,
			sqlite3.ErrCorrupt,
			sqlite3.ErrFull,
			sqlite3.ErrCantOpen:
			return fmt.Errorf(
				"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
				persist.ErrUnavailable,
				operation,
				sqliteErr.Code,
				sqliteErr.ExtendedCode,
				err,
			)

		case sqlite3.ErrConstraint:
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintCheck:
				return fmt.Errorf(
					"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
					persist.ErrCheckViolation,
					operation,
					sqliteErr.Code,
					sqliteErr.ExtendedCode,
					err,
				)

			case sqlite3.ErrConstraintForeignKey:
				return fmt.Errorf(
					"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
					foreignKeyError,
					operation,
					sqliteErr.Code,
					sqliteErr.ExtendedCode,
					err,
				)

			case sqlite3.ErrConstraintNotNull:
				return fmt.Errorf(
					"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
					persist.ErrNotNullViolation,
					operation,
					sqliteErr.Code,
					sqliteErr.ExtendedCode,
					err,
				)

			case sqlite3.ErrConstraintPrimaryKey,
				sqlite3.ErrConstraintUnique,
				sqlite3.ErrConstraintRowID:
				return fmt.Errorf(
					"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
					persist.ErrUniqueViolation,
					operation,
					sqliteErr.Code,
					sqliteErr.ExtendedCode,
					err,
				)
			}
		}
	}
	return err
}

func MapSQLiteError(operation string, err error) error {
	// MapSQLiteError is for non-DELETE operations (read/insert/update).
	return mapSQLiteBaseError(operation, err, persist.ErrForeignKeyReferenceNotFound)
}

func MapSQLiteDeleteError(operation string, err error) error {
	return mapSQLiteBaseError(operation, err, persist.ErrForeignKeyReferenced)
}
