package persist

import (
	"errors"
	"fmt"
)

var (
	ErrConflict = errors.New("conflict") // Kept for backward compatibility.
	ErrNoRows   = errors.New("no rows")

	ErrQuota          = errors.New("quota")
	ErrQuotaExceeded  = fmt.Errorf("%w exceeded", ErrQuota)
	ErrQuotaUnderflow = fmt.Errorf("%w underflow", ErrQuota)

	// Canceled errors
	ErrCanceled         = errors.New("canceled")
	ErrContextCanceled  = fmt.Errorf("%w: context canceled", ErrCanceled)
	ErrDeadlineExceeded = fmt.Errorf("%w: deadline exceeded", ErrCanceled)
	ErrQueryCanceled    = fmt.Errorf("%w: query canceled", ErrCanceled)

	// Connection errors (MySQL and PostgreSQL only)
	ErrConnection            = errors.New("connection")
	ErrConnectionInit        = fmt.Errorf("%w init failed", ErrConnection)
	ErrConnectionLost        = fmt.Errorf("%w lost", ErrConnection)
	ErrConnectionRefused     = fmt.Errorf("%w refused", ErrConnection)
	ErrConnectionTimeout     = fmt.Errorf("%w timeout", ErrConnection)
	ErrConnectionUnavailable = fmt.Errorf("%w unavailable", ErrConnection)

	// Constraint violations
	ErrConstraintViolation         = errors.New("constraint")
	ErrCheckViolation              = fmt.Errorf("%w: check violation", ErrConstraintViolation)
	ErrExclusionViolation          = fmt.Errorf("%w: exclusion violation", ErrConstraintViolation)
	ErrForeignKeyReferenceNotFound = fmt.Errorf("%w: foreign key reference not found", ErrConstraintViolation)
	ErrForeignKeyReferenced        = fmt.Errorf("%w: foreign key referenced", ErrConstraintViolation)
	ErrNotNullViolation            = fmt.Errorf("%w: not null violation", ErrConstraintViolation)
	ErrUniqueViolation             = fmt.Errorf("%w: unique violation", ErrConstraintViolation)

	ErrInvalidParam = errors.New("invalid parameter")

	// Contention errors
	ErrContention       = errors.New("contention")
	ErrDeadlockDetected = fmt.Errorf("%w: deadlock detected", ErrContention)
	ErrLockTimeout      = fmt.Errorf("%w: lock timeout", ErrContention)
	ErrLockUnavailable  = fmt.Errorf("%w: lock unavailable", ErrContention)
	ErrResourceBusy     = fmt.Errorf("%w: resource busy", ErrContention) // Used for SQLite ErrBusy/ErrLocked.
	ErrTxSerialization  = fmt.Errorf("%w: transaction serialization", ErrContention)

	// Resource exhaustion
	ErrResourceExhausted  = errors.New("resource exhausted")
	ErrOutOfMemory        = fmt.Errorf("%w: out of memory", ErrResourceExhausted)
	ErrStorageFull        = fmt.Errorf("%w: storage full", ErrResourceExhausted)
	ErrTooManyConnections = fmt.Errorf("%w: too many connections", ErrResourceExhausted)

	// Storage errors
	ErrStorage            = errors.New("storage")
	ErrStorageCorrupted   = fmt.Errorf("%w corrupted", ErrStorage)
	ErrStorageUnavailable = fmt.Errorf("%w unavailable", ErrStorage)

	// Syntax errors
	ErrSyntax           = errors.New("syntax error")
	ErrInvalidStatement = errors.New("invalid statement")
)

func IsRecoverable(err error) bool {
	return errors.Is(err, ErrConnectionTimeout) ||
		errors.Is(err, ErrConnectionLost) ||
		errors.Is(err, ErrContention)
}
