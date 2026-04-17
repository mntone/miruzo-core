package dberrors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

const (
	errCantLock                     uint16 = 1015 // ER_CANT_LOCK
	errCantOpenFile                 uint16 = 1016 // ER_CANT_OPEN_FILE
	errFileNotFound                 uint16 = 1017 // ER_FILE_NOT_FOUND
	errCantReadDir                  uint16 = 1018 // ER_CANT_READ_DIR
	errDiskFull                     uint16 = 1021 // ER_DISK_FULL
	errErrorOnRead                  uint16 = 1024 // ER_ERROR_ON_READ
	errErrorOnWrite                 uint16 = 1026 // ER_ERROR_ON_WRITE
	errGetErrno                     uint16 = 1030 // ER_GET_ERRNO
	errNotFormFile                  uint16 = 1033 // ER_NOT_FORM_FILE
	errNotKeyFile                   uint16 = 1034 // ER_NOT_KEYFILE
	errOldKeyFile                   uint16 = 1035 // ER_OLD_KEYFILE
	errOpenAsReadOnly               uint16 = 1036 // ER_OPEN_AS_READONLY
	errOutOfMemory                  uint16 = 1037 // ER_OUTOFMEMORY
	errConnCountError               uint16 = 1040 // ER_CON_COUNT_ERROR
	errAccessDeniedError            uint16 = 1045 // ER_ACCESS_DENIED_ERROR
	errBadNullError                 uint16 = 1048 // ER_BAD_NULL_ERROR
	errBadFieldError                uint16 = 1054 // ER_BAD_FIELD_ERROR
	errWrongFieldWithGroup          uint16 = 1055 // ER_WRONG_FIELD_WITH_GROUP
	errDuplicateEntry               uint16 = 1062 // ER_DUP_ENTRY
	errParseError                   uint16 = 1064 // ER_PARSE_ERROR
	errInvalidGroupFuncUse          uint16 = 1111 // ER_INVALID_GROUP_FUNC_USE
	errRecordFileFull               uint16 = 1114 // ER_RECORD_FILE_FULL
	errTooManyFields                uint16 = 1117 // ER_TOO_MANY_FIELDS
	errTooBigRowSize                uint16 = 1118 // ER_TOO_BIG_ROWSIZE
	errHostNotPrivileged            uint16 = 1130 // ER_HOST_NOT_PRIVILEGED
	errNoSuchTable                  uint16 = 1146 // ER_NO_SUCH_TABLE
	errNetPacketTooLarge            uint16 = 1153 // ER_NET_PACKET_TOO_LARGE
	errTooManyUserConnections       uint16 = 1203 // ER_TOO_MANY_USER_CONNECTIONS
	errLockWaitTimeout              uint16 = 1205 // ER_LOCK_WAIT_TIMEOUT
	errLockDeadlock                 uint16 = 1213 // ER_LOCK_DEADLOCK
	errWrongNumberOfColumnsInSelect uint16 = 1222 // ER_WRONG_NUMBER_OF_COLUMNS_IN_SELECT
	errDerivedMustHaveAlias         uint16 = 1248 // ER_DERIVED_MUST_HAVE_ALIAS
	errWarnDataOutOfRange           uint16 = 1264 // ER_WARN_DATA_OUT_OF_RANGE
	errWarnDataTruncated            uint16 = 1265 // ER_WARN_DATA_TRUNCATED
	errOptionPreventsStatement      uint16 = 1290 // ER_OPTION_PREVENTS_STATEMENT
	errTruncatedWrongValue          uint16 = 1292 // ER_TRUNCATED_WRONG_VALUE
	errStoredProcedureDoesNotExist  uint16 = 1305 // ER_SP_DOES_NOT_EXIST
	errQueryInterrupted             uint16 = 1317 // ER_QUERY_INTERRUPTED
	errStoredProcedureWrongNoOfArgs uint16 = 1318 // ER_SP_WRONG_NO_OF_ARGS
	errNoDefaultForField            uint16 = 1364 // ER_NO_DEFAULT_FOR_FIELD
	errTruncatedWrongValueForField  uint16 = 1366 // ER_TRUNCATED_WRONG_VALUE_FOR_FIELD
	errPlaceholderManyParam         uint16 = 1390 // ER_PS_MANY_PARAM
	errDataTooLong                  uint16 = 1406 // ER_DATA_TOO_LONG
	errRowIsReferenced2             uint16 = 1451 // ER_ROW_IS_REFERENCED_2
	errNoReferencedRow2             uint16 = 1452 // ER_NO_REFERENCED_ROW_2
	errWrongParamCountToNativeFunc  uint16 = 1582 // ER_WRONG_PARAMCOUNT_TO_NATIVE_FCT
	errDataOutOfRange               uint16 = 1690 // ER_DATA_OUT_OF_RANGE
	errCantExecuteInReadOnlyTx      uint16 = 1792 // ER_CANT_EXECUTE_IN_READ_ONLY_TRANSACTION
	errReadOnlyMode                 uint16 = 1836 // ER_READ_ONLY_MODE
	errConnectionError              uint16 = 2002 // CR_CONNECTION_ERROR
	errConnectionHostError          uint16 = 2003 // CR_CONN_HOST_ERROR
	errConnectionUnknownProtocol    uint16 = 2005 // CR_CONN_UNKNOWN_PROTOCOL
	errServerGoneError              uint16 = 2006 // CR_SERVER_GONE_ERROR
	errServerLost                   uint16 = 2013 // CR_SERVER_LOST
	errServerLostExtended           uint16 = 2055 // CR_SERVER_LOST_EXTENDED
	errFileCorrupt                  uint16 = 3000 // ER_FILE_CORRUPT
	errCheckConstraintViolated      uint16 = 3819 // ER_CHECK_CONSTRAINT_VIOLATED
)

var mysqlToPersistError = map[uint16]error{
	// Canceled errors
	errQueryInterrupted: persist.ErrQueryCanceled,

	// Connection errors
	errAccessDeniedError:   persist.ErrAuthorizationFailed,
	errConnectionError:     persist.ErrConnectionInit,
	errConnectionHostError: persist.ErrConnectionInit,
	errServerGoneError:     persist.ErrConnectionLost,
	errServerLost:          persist.ErrConnectionLost,
	errServerLostExtended:  persist.ErrConnectionLost,
	errHostNotPrivileged:   persist.ErrConnectionRefused,

	// Constraint violations
	errCheckConstraintViolated: persist.ErrCheckViolation,
	errDataOutOfRange:          persist.ErrCheckViolation,
	errDataTooLong:             persist.ErrCheckViolation,
	errWarnDataOutOfRange:      persist.ErrCheckViolation,
	errRowIsReferenced2:        persist.ErrForeignKeyReferenced,
	errNoReferencedRow2:        persist.ErrForeignKeyReferenceNotFound,
	errBadNullError:            persist.ErrNotNullViolation,
	errNoDefaultForField:       persist.ErrNotNullViolation,
	errDuplicateEntry:          persist.ErrUniqueViolation,

	errConnectionUnknownProtocol:   persist.ErrInvalidParam,
	errTruncatedWrongValue:         persist.ErrInvalidParam,
	errTruncatedWrongValueForField: persist.ErrInvalidParam,
	errWarnDataTruncated:           persist.ErrInvalidParam,

	// Contention errors
	errLockDeadlock:    persist.ErrDeadlockDetected,
	errLockWaitTimeout: persist.ErrLockTimeout,
	errCantLock:        persist.ErrLockUnavailable,

	// Resource exhaustion
	errOutOfMemory:            persist.ErrOutOfMemory,
	errDiskFull:               persist.ErrStorageFull,
	errRecordFileFull:         persist.ErrStorageFull,
	errConnCountError:         persist.ErrTooManyConnections,
	errTooManyUserConnections: persist.ErrTooManyConnections,

	// Storage errors
	errFileCorrupt:    persist.ErrStorageCorrupted,
	errNotFormFile:    persist.ErrStorageCorrupted,
	errNotKeyFile:     persist.ErrStorageCorrupted,
	errOldKeyFile:     persist.ErrStorageCorrupted,
	errOpenAsReadOnly: persist.ErrStorageReadonly,
	errReadOnlyMode:   persist.ErrStorageReadonly,
	errCantOpenFile:   persist.ErrStorageUnavailable,
	errCantReadDir:    persist.ErrStorageUnavailable,
	errErrorOnRead:    persist.ErrStorageUnavailable,
	errErrorOnWrite:   persist.ErrStorageUnavailable,
	errFileNotFound:   persist.ErrStorageUnavailable,

	// Syntax errors
	errParseError:                   persist.ErrSyntax,
	errBadFieldError:                persist.ErrInvalidStatement,
	errWrongFieldWithGroup:          persist.ErrInvalidStatement,
	errInvalidGroupFuncUse:          persist.ErrInvalidStatement,
	errNoSuchTable:                  persist.ErrInvalidStatement,
	errTooManyFields:                persist.ErrInvalidStatement,
	errTooBigRowSize:                persist.ErrInvalidStatement,
	errNetPacketTooLarge:            persist.ErrInvalidStatement,
	errWrongNumberOfColumnsInSelect: persist.ErrInvalidStatement,
	errDerivedMustHaveAlias:         persist.ErrInvalidStatement,
	errStoredProcedureDoesNotExist:  persist.ErrInvalidStatement,
	errStoredProcedureWrongNoOfArgs: persist.ErrInvalidStatement,
	errPlaceholderManyParam:         persist.ErrInvalidStatement,
	errWrongParamCountToNativeFunc:  persist.ErrInvalidStatement,

	// Transaction errors
	errCantExecuteInReadOnlyTx: persist.ErrTxReadonly,
}

func wrapMySQLError(operation string, persistError error, mySQLError *mysql.MySQLError) error {
	return fmt.Errorf(
		"%w: operation=%s number=%d sqlstate=%s: %s",
		persistError,
		operation,
		mySQLError.Number,
		mySQLError.SQLState,
		mySQLError.Message,
	)
}

func ToPersist(operation string, err error) error {
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

	if errors.Is(err, mysql.ErrInvalidConn) {
		return fmt.Errorf(
			"%w: operation=%s: %v",
			persist.ErrConnectionLost,
			operation,
			err,
		)
	}

	if netError, ok := errors.AsType[net.Error](err); ok && netError.Timeout() {
		return fmt.Errorf(
			"%w: operation=%s: %v",
			persist.ErrConnectionTimeout,
			operation,
			err,
		)
	}

	if mySQLError, ok := errors.AsType[*mysql.MySQLError](err); ok {
		switch mySQLError.Number {
		case errGetErrno:
			if strings.Contains(mySQLError.Message, "errno 28") {
				return wrapMySQLError(operation, persist.ErrStorageFull, mySQLError)
			}
			return wrapMySQLError(operation, persist.ErrStorageUnavailable, mySQLError)
		case errOptionPreventsStatement:
			messageLower := strings.ToLower(mySQLError.Message)
			if strings.Contains(messageLower, "read only") ||
				strings.Contains(messageLower, "super_read_only") {
				return wrapMySQLError(operation, persist.ErrStorageReadonly, mySQLError)
			}
		}

		if persistError, ok := mysqlToPersistError[mySQLError.Number]; ok {
			return wrapMySQLError(operation, persistError, mySQLError)
		}
	}
	return err
}
