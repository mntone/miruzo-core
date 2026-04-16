package shared

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapPersistError(
	operation string,
	persistError error,
	sqliteError *sqlite3.Error,
	err error,
) error {
	return fmt.Errorf(
		"%w: operation=%s sqlite_code=%d sqlite_extended_code=%d: %v",
		persistError,
		operation,
		sqliteError.Code,
		sqliteError.ExtendedCode,
		err,
	)
}

func mapSQLiteBaseError(operation string, err error, foreignKeyError error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return fmt.Errorf(
			"%w: operation=%s: %v",
			persist.ErrContextCanceled,
			operation,
			err,
		)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf(
			"%w: operation=%s: %v",
			persist.ErrDeadlineExceeded,
			operation,
			err,
		)
	}

	if sqliteError, ok := errors.AsType[sqlite3.Error](err); ok {
		switch sqliteError.Code {
		// Canceled errors
		case sqlite3.ErrInterrupt:
			return mapPersistError(operation, persist.ErrQueryCanceled, &sqliteError, err)

		// Constraint violations
		case sqlite3.ErrConstraint:
			switch sqliteError.ExtendedCode {
			case sqlite3.ErrConstraintCheck:
				return mapPersistError(operation, persist.ErrCheckViolation, &sqliteError, err)
			case sqlite3.ErrConstraintForeignKey:
				return mapPersistError(operation, foreignKeyError, &sqliteError, err)
			case sqlite3.ErrConstraintNotNull:
				return mapPersistError(operation, persist.ErrNotNullViolation, &sqliteError, err)
			case sqlite3.ErrConstraintPrimaryKey,
				sqlite3.ErrConstraintUnique,
				sqlite3.ErrConstraintRowID:
				return mapPersistError(operation, persist.ErrUniqueViolation, &sqliteError, err)
			}
		case sqlite3.ErrTooBig:
			return mapPersistError(operation, persist.ErrCheckViolation, &sqliteError, err)
		case sqlite3.ErrError:
			if strings.Contains(sqliteError.Error(), "integer overflow") {
				return mapPersistError(operation, persist.ErrCheckViolation, &sqliteError, err)
			}
			if strings.Contains(sqliteError.Error(), "syntax error") {
				return mapPersistError(operation, persist.ErrSyntax, &sqliteError, err)
			}
			if strings.Contains(sqliteError.Error(), "too many columns") ||
				strings.Contains(sqliteError.Error(), "too many SQL variables") {
				return mapPersistError(operation, persist.ErrInvalidStatement, &sqliteError, err)
			}
		case sqlite3.ErrMismatch:
			return mapPersistError(operation, persist.ErrInvalidParam, &sqliteError, err)

		// Contention errors
		case sqlite3.ErrBusy, sqlite3.ErrLocked:
			return mapPersistError(operation, persist.ErrResourceBusy, &sqliteError, err)

		// Resource exhaustion
		case sqlite3.ErrNomem:
			return mapPersistError(operation, persist.ErrOutOfMemory, &sqliteError, err)
		case sqlite3.ErrFull:
			return mapPersistError(operation, persist.ErrStorageFull, &sqliteError, err)

		// Storage errors
		case sqlite3.ErrIoErr, sqlite3.ErrCantOpen:
			return mapPersistError(operation, persist.ErrStorageUnavailable, &sqliteError, err)
		case sqlite3.ErrCorrupt:
			return mapPersistError(operation, persist.ErrStorageCorrupted, &sqliteError, err)
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
