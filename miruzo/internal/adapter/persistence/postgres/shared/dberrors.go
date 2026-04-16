package shared

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

var pgToPersistError = map[string]error{
	// Canceled errors
	pgerrcode.QueryCanceled: persist.ErrQueryCanceled,

	// Connection errors
	pgerrcode.SQLClientUnableToEstablishSQLConnection:       persist.ErrConnectionInit,
	pgerrcode.ConnectionDoesNotExist:                        persist.ErrConnectionLost,
	pgerrcode.ConnectionFailure:                             persist.ErrConnectionLost,
	pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection: persist.ErrConnectionRefused,
	pgerrcode.CannotConnectNow:                              persist.ErrConnectionUnavailable,
	pgerrcode.ConnectionException:                           persist.ErrConnectionUnavailable,

	// Constraint violations
	pgerrcode.CheckViolation:     persist.ErrCheckViolation,
	pgerrcode.ExclusionViolation: persist.ErrExclusionViolation,
	pgerrcode.NotNullViolation:   persist.ErrNotNullViolation,
	pgerrcode.UniqueViolation:    persist.ErrUniqueViolation,

	pgerrcode.StringDataRightTruncationDataException: persist.ErrCheckViolation,
	pgerrcode.NumericValueOutOfRange:                 persist.ErrCheckViolation,
	pgerrcode.InvalidDatetimeFormat:                  persist.ErrCheckViolation,
	pgerrcode.DatetimeFieldOverflow:                  persist.ErrCheckViolation,

	pgerrcode.InvalidCharacterValueForCast: persist.ErrInvalidParam,
	pgerrcode.InvalidTextRepresentation:    persist.ErrInvalidParam,

	// Contention errors
	pgerrcode.DeadlockDetected:     persist.ErrDeadlockDetected,
	pgerrcode.SerializationFailure: persist.ErrTxSerialization,

	// Resource exhaustion
	pgerrcode.OutOfMemory:        persist.ErrOutOfMemory,
	pgerrcode.DiskFull:           persist.ErrStorageFull,
	pgerrcode.TooManyConnections: persist.ErrTooManyConnections,

	// Storage errors
	pgerrcode.DataCorrupted:  persist.ErrStorageCorrupted,
	pgerrcode.IndexCorrupted: persist.ErrStorageCorrupted,
	pgerrcode.IOError:        persist.ErrStorageUnavailable,

	// Syntax errors
	pgerrcode.SyntaxError:           persist.ErrSyntax,
	pgerrcode.DatatypeMismatch:      persist.ErrInvalidStatement,
	pgerrcode.FeatureNotSupported:   persist.ErrInvalidStatement,
	pgerrcode.IndeterminateDatatype: persist.ErrInvalidStatement,
	pgerrcode.ProgramLimitExceeded:  persist.ErrInvalidStatement,
	pgerrcode.StatementTooComplex:   persist.ErrInvalidStatement,
	pgerrcode.TooManyColumns:        persist.ErrInvalidStatement,
	pgerrcode.TooManyArguments:      persist.ErrInvalidStatement,
	pgerrcode.UndefinedColumn:       persist.ErrInvalidStatement,
	pgerrcode.UndefinedFunction:     persist.ErrInvalidStatement,
	pgerrcode.UndefinedParameter:    persist.ErrInvalidStatement,
	pgerrcode.UndefinedTable:        persist.ErrInvalidStatement,
}

var isPgconnTimeoutError = pgconn.Timeout

func toPersistError(operation string, persistError error, pgError *pgconn.PgError) error {
	return fmt.Errorf(
		"%w: operation=%s sqlstate=%s: %s",
		persistError,
		operation,
		pgError.Code,
		pgError.Message,
	)
}

func mapPostgreBaseError(operation string, err error, foreignKeyError error) error {
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

	if isPgconnTimeoutError(err) {
		return fmt.Errorf(
			"%w: operation=%s: %v",
			persist.ErrConnectionTimeout,
			operation,
			err,
		)
	}

	if pgError, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgError.Code == pgerrcode.ForeignKeyViolation {
			return toPersistError(operation, foreignKeyError, pgError)
		}

		if pgError.Code == pgerrcode.LockNotAvailable {
			if strings.Contains(pgError.Message, "lock timeout") {
				return toPersistError(operation, persist.ErrLockTimeout, pgError)
			}
			return toPersistError(operation, persist.ErrLockUnavailable, pgError)
		}

		if persistError, ok := pgToPersistError[pgError.Code]; ok {
			return toPersistError(operation, persistError, pgError)
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
