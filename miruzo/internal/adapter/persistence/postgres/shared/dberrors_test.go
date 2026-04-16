package shared_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestMapPostgreErrorMapsContextErrors(t *testing.T) {
	tests := []struct {
		name    string
		build   func() error
		wantErr error
	}{
		{
			name: "context_canceled",
			build: func() error {
				return context.Canceled
			},
			wantErr: persist.ErrContextCanceled,
		},
		{
			name: "context_deadline_exceeded",
			build: func() error {
				return context.DeadlineExceeded
			},
			wantErr: persist.ErrDeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapPostgreError("ListLatest", tt.build())
			assert.ErrorIs(
				t,
				"MapPostgreError("+tt.name+")",
				err,
				tt.wantErr,
			)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestMapPostgreErrorMapsSQLStates(t *testing.T) {
	tests := []struct {
		name     string
		sqlState string
		wantErr  error
	}{
		// Canceled errors
		{"query_canceled", pgerrcode.QueryCanceled, persist.ErrQueryCanceled},

		// Connection errors
		{"sqlclient_unable_to_establish_sqlconnection", pgerrcode.SQLClientUnableToEstablishSQLConnection, persist.ErrConnectionInit},
		{"connection_does_not_exist", pgerrcode.ConnectionDoesNotExist, persist.ErrConnectionLost},
		{"connection_failure", pgerrcode.ConnectionFailure, persist.ErrConnectionLost},
		{"sqlserver_rejected_establishment_of_sqlconnection", pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection, persist.ErrConnectionRefused},
		{"cannot_connect_now", pgerrcode.CannotConnectNow, persist.ErrConnectionUnavailable},
		{"connection_exception", pgerrcode.ConnectionException, persist.ErrConnectionUnavailable},

		// Constraint violations
		{"check_violation", pgerrcode.CheckViolation, persist.ErrCheckViolation},
		{"exclusion_violation", pgerrcode.ExclusionViolation, persist.ErrExclusionViolation},
		{"foreign_key_violation", pgerrcode.ForeignKeyViolation, persist.ErrForeignKeyReferenceNotFound},
		{"not_null_violation", pgerrcode.NotNullViolation, persist.ErrNotNullViolation},
		{"unique_violation", pgerrcode.UniqueViolation, persist.ErrUniqueViolation},

		{"string_data_right_truncation", pgerrcode.StringDataRightTruncationDataException, persist.ErrCheckViolation},
		{"numeric_value_out_of_range", pgerrcode.NumericValueOutOfRange, persist.ErrCheckViolation},
		{"invalid_datetime_format", pgerrcode.InvalidDatetimeFormat, persist.ErrCheckViolation},
		{"datetime_field_overflow", pgerrcode.DatetimeFieldOverflow, persist.ErrCheckViolation},

		{"invalid_character_value_for_cast", pgerrcode.InvalidCharacterValueForCast, persist.ErrInvalidParam},
		{"invalid_text_representation", pgerrcode.InvalidTextRepresentation, persist.ErrInvalidParam},

		// Contention errors
		{"deadlock_detected", pgerrcode.DeadlockDetected, persist.ErrDeadlockDetected},
		{"lock_not_available", pgerrcode.LockNotAvailable, persist.ErrLockUnavailable},
		{"serialization_failure", pgerrcode.SerializationFailure, persist.ErrTxSerialization},

		// Resource exhaustion
		{"out_of_memory", pgerrcode.OutOfMemory, persist.ErrOutOfMemory},
		{"disk_full", pgerrcode.DiskFull, persist.ErrStorageFull},
		{"too_many_connections", pgerrcode.TooManyConnections, persist.ErrTooManyConnections},

		// Storage errors
		{"data_corrupted", pgerrcode.DataCorrupted, persist.ErrStorageCorrupted},
		{"index_corrupted", pgerrcode.IndexCorrupted, persist.ErrStorageCorrupted},
		{"io_error", pgerrcode.IOError, persist.ErrStorageUnavailable},

		// Syntax errors
		{"syntax_error", pgerrcode.SyntaxError, persist.ErrSyntax},
		{"datatype_mismatch", pgerrcode.DatatypeMismatch, persist.ErrInvalidStatement},
		{"feature_not_supported", pgerrcode.FeatureNotSupported, persist.ErrInvalidStatement},
		{"indeterminate_datatype", pgerrcode.IndeterminateDatatype, persist.ErrInvalidStatement},
		{"program_limit_exceeded", pgerrcode.ProgramLimitExceeded, persist.ErrInvalidStatement},
		{"statement_too_complex", pgerrcode.StatementTooComplex, persist.ErrInvalidStatement},
		{"too_many_columns", pgerrcode.TooManyColumns, persist.ErrInvalidStatement},
		{"too_many_arguments", pgerrcode.TooManyArguments, persist.ErrInvalidStatement},
		{"undefined_column", pgerrcode.UndefinedColumn, persist.ErrInvalidStatement},
		{"undefined_function", pgerrcode.UndefinedFunction, persist.ErrInvalidStatement},
		{"undefined_parameter", pgerrcode.UndefinedParameter, persist.ErrInvalidStatement},
		{"undefined_table", pgerrcode.UndefinedTable, persist.ErrInvalidStatement},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapPostgreError(
				"ListLatest",
				fmt.Errorf("query failed: %w", &pgconn.PgError{Code: tt.sqlState}),
			)
			assert.ErrorIs(
				t,
				fmt.Sprintf("MapPostgreError(%s)", tt.sqlState),
				err, tt.wantErr,
			)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
			if !strings.Contains(err.Error(), "sqlstate="+tt.sqlState) {
				t.Fatalf("expected sqlstate detail, got %v", err)
			}
		})
	}
}

func TestMapPostgreErrorMapsLockNotAvailableWithTimeoutMessage(t *testing.T) {
	err := shared.MapPostgreError(
		"ListLatest",
		fmt.Errorf(
			"query failed: %w",
			&pgconn.PgError{
				Code:    pgerrcode.LockNotAvailable,
				Message: "canceling statement due to lock timeout",
			},
		),
	)
	assert.ErrorIs(
		t,
		"MapPostgreError(55P03 lock timeout)",
		err,
		persist.ErrLockTimeout,
	)
	if !strings.Contains(err.Error(), "operation=ListLatest") {
		t.Fatalf("expected operation detail, got %v", err)
	}
	if !strings.Contains(err.Error(), "sqlstate="+pgerrcode.LockNotAvailable) {
		t.Fatalf("expected sqlstate detail, got %v", err)
	}
}

func TestMapPostgreDeleteErrorMapsForeignKeyViolationToReferenced(t *testing.T) {
	err := shared.MapPostgreDeleteError(
		"DeleteImage",
		fmt.Errorf("query failed: %w", &pgconn.PgError{Code: "23503"}),
	)
	assert.ErrorIs(
		t,
		"MapPostgreDeleteError(23503)",
		err,
		persist.ErrForeignKeyReferenced,
	)
	if !strings.Contains(err.Error(), "operation=DeleteImage") {
		t.Fatalf("expected operation detail, got %v", err)
	}
	if !strings.Contains(err.Error(), "sqlstate=23503") {
		t.Fatalf("expected sqlstate detail, got %v", err)
	}
}

func TestMapPostgreErrorPassesThroughErrors(t *testing.T) {
	tests := []struct {
		name  string
		inErr error
	}{
		// Nil error
		{"nil error", nil},

		// Unknown error
		{"unknown error", &pgconn.PgError{Code: pgerrcode.StringDataLengthMismatch}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapPostgreError("ListLatest", tt.inErr)
			if !errors.Is(err, tt.inErr) {
				t.Fatalf("err = got %v, want original error", err)
			}

			inPgErr, ok := errors.AsType[*pgconn.PgError](tt.inErr)
			if !ok {
				return
			}

			var pgError *pgconn.PgError
			if !errors.As(err, &pgError) {
				t.Fatalf("expected *pgconn.PgError, got %T", err)
			}
			if pgError.Code != inPgErr.Code {
				t.Fatalf("unexpected sqlstate: %s", pgError.Code)
			}
		})
	}
}
