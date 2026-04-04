package shared

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapPostgreBaseError(operation string, err error, foreignKeyError error) error {
	if err == nil {
		return nil
	}

	if pgconn.Timeout(err) {
		return fmt.Errorf(
			"%w: operation=%s: %v",
			persist.ErrTimeout,
			operation,
			err,
		)
	}

	if pgError, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgError.Code {
		case pgerrcode.QueryCanceled:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrTimeout,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.ConnectionDoesNotExist,
			pgerrcode.TooManyConnections:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrRecoverableUnavailable,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.SerializationFailure,
			pgerrcode.DeadlockDetected,
			pgerrcode.LockNotAvailable:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrRecoverableConflict,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.ConnectionFailure,
			pgerrcode.SQLClientUnableToEstablishSQLConnection,
			pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
			pgerrcode.DiskFull,
			pgerrcode.OutOfMemory,
			pgerrcode.IOError,
			pgerrcode.DataCorrupted,
			pgerrcode.IndexCorrupted:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrUnavailable,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.NotNullViolation:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrNotNullViolation,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.ForeignKeyViolation:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				foreignKeyError,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.UniqueViolation:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrUniqueViolation,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.CheckViolation, pgerrcode.StringDataRightTruncationDataException, pgerrcode.NumericValueOutOfRange:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrCheckViolation,
				operation,
				pgError.Code,
				err,
			)

		case pgerrcode.ExclusionViolation:
			return fmt.Errorf(
				"%w: operation=%s sqlstate=%s: %v",
				persist.ErrExclusionViolation,
				operation,
				pgError.Code,
				err,
			)
		}
	}
	return err
}

func MapPostgreError(operation string, err error) error {
	// MapPostgreError is for non-DELETE operations (read/insert/update).
	return mapPostgreBaseError(operation, err, persist.ErrForeignKeyReferenceNotFound)
}

func MapPostgreDeleteError(operation string, err error) error {
	return mapPostgreBaseError(operation, err, persist.ErrForeignKeyReferenced)
}
