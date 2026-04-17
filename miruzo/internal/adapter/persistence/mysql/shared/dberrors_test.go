package shared_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type timeoutError struct{}

func (timeoutError) Error() string   { return "i/o timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

var _ net.Error = timeoutError{}

func TestMapMySQLErrorMapsContextErrors(t *testing.T) {
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
			err := shared.MapMySQLError("ListLatest", tt.build())
			assert.ErrorIs(
				t,
				"MapMySQLError("+tt.name+")",
				err,
				tt.wantErr,
			)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestMapMySQLErrorMapsNumbers(t *testing.T) {
	tests := []struct {
		name    string
		number  uint16
		wantErr error
	}{
		// Canceled errors
		{"query_interrupted", 1317, persist.ErrQueryCanceled},

		// Connection errors
		{"access_denied_error", 1045, persist.ErrAuthorizationFailed},
		{"host_not_privileged", 1130, persist.ErrConnectionRefused},
		{"connection_error", 2002, persist.ErrConnectionInit},
		{"connection_host_error", 2003, persist.ErrConnectionInit},
		{"server_gone_error", 2006, persist.ErrConnectionLost},
		{"server_lost", 2013, persist.ErrConnectionLost},
		{"server_lost_extended", 2055, persist.ErrConnectionLost},

		// Constraint violations
		{"check_constraint_violated", 3819, persist.ErrCheckViolation},
		{"data_out_of_range", 1690, persist.ErrCheckViolation},
		{"data_too_long", 1406, persist.ErrCheckViolation},
		{"warn_data_out_of_range", 1264, persist.ErrCheckViolation},
		{"row_is_referenced_2", 1451, persist.ErrForeignKeyReferenced},
		{"no_referenced_row_2", 1452, persist.ErrForeignKeyReferenceNotFound},
		{"bad_null_error", 1048, persist.ErrNotNullViolation},
		{"no_default_for_field", 1364, persist.ErrNotNullViolation},
		{"duplicate_entry", 1062, persist.ErrUniqueViolation},
		{"connection_unknown_protocol", 2005, persist.ErrInvalidParam},
		{"truncated_wrong_value", 1292, persist.ErrInvalidParam},
		{"truncated_wrong_value_for_field", 1366, persist.ErrInvalidParam},
		{"warn_data_truncated", 1265, persist.ErrInvalidParam},

		// Contention errors
		{"lock_deadlock", 1213, persist.ErrDeadlockDetected},
		{"lock_wait_timeout", 1205, persist.ErrLockTimeout},
		{"cant_lock", 1015, persist.ErrLockUnavailable},

		// Resource exhaustion
		{"out_of_memory", 1037, persist.ErrOutOfMemory},
		{"disk_full", 1021, persist.ErrStorageFull},
		{"record_file_full", 1114, persist.ErrStorageFull},
		{"conn_count_error", 1040, persist.ErrTooManyConnections},
		{"too_many_user_connections", 1203, persist.ErrTooManyConnections},

		// Storage errors
		{"file_corrupt", 3000, persist.ErrStorageCorrupted},
		{"not_form_file", 1033, persist.ErrStorageCorrupted},
		{"not_key_file", 1034, persist.ErrStorageCorrupted},
		{"old_key_file", 1035, persist.ErrStorageCorrupted},
		{"open_as_read_only", 1036, persist.ErrStorageReadonly},
		{"read_only_mode", 1836, persist.ErrStorageReadonly},
		{"cant_open_file", 1016, persist.ErrStorageUnavailable},
		{"cant_read_dir", 1018, persist.ErrStorageUnavailable},
		{"error_on_read", 1024, persist.ErrStorageUnavailable},
		{"error_on_write", 1026, persist.ErrStorageUnavailable},
		{"file_not_found", 1017, persist.ErrStorageUnavailable},

		// Syntax errors
		{"parse_error", 1064, persist.ErrSyntax},
		{"bad_field_error", 1054, persist.ErrInvalidStatement},
		{"wrong_field_with_group", 1055, persist.ErrInvalidStatement},
		{"invalid_group_func_use", 1111, persist.ErrInvalidStatement},
		{"no_such_table", 1146, persist.ErrInvalidStatement},
		{"too_many_fields", 1117, persist.ErrInvalidStatement},
		{"too_big_row_size", 1118, persist.ErrInvalidStatement},
		{"net_packet_too_large", 1153, persist.ErrInvalidStatement},
		{"wrong_number_of_columns_in_select", 1222, persist.ErrInvalidStatement},
		{"derived_must_have_alias", 1248, persist.ErrInvalidStatement},
		{"stored_procedure_does_not_exist", 1305, persist.ErrInvalidStatement},
		{"stored_procedure_wrong_no_of_args", 1318, persist.ErrInvalidStatement},
		{"placeholder_many_param", 1390, persist.ErrInvalidStatement},
		{"wrong_param_count_to_native_func", 1582, persist.ErrInvalidStatement},

		// Transaction errors
		{"cant_execute_in_read_only_tx", 1792, persist.ErrTxReadonly},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapMySQLError(
				"ListLatest",
				fmt.Errorf(
					"query failed: %w",
					&mysql.MySQLError{
						Number:   tt.number,
						SQLState: [5]byte{'H', 'Y', '0', '0', '0'},
						Message:  "query failed",
					},
				),
			)
			assert.ErrorIs(
				t,
				fmt.Sprintf("MapMySQLError(%d)", tt.number),
				err,
				tt.wantErr,
			)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
			if !strings.Contains(err.Error(), fmt.Sprintf("number=%d", tt.number)) {
				t.Fatalf("expected number detail, got %v", err)
			}
			if !strings.Contains(err.Error(), "sqlstate=") {
				t.Fatalf("expected sqlstate detail, got %v", err)
			}
		})
	}
}

func TestMapMySQLErrorMapsGetErrnoByMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr error
	}{
		{"storage_full", "Got errno 28 from storage engine", persist.ErrStorageFull},
		{"storage_unavailable", "Got error 5 from storage engine", persist.ErrStorageUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapMySQLError(
				"ListLatest",
				&mysql.MySQLError{
					Number:   1030,
					SQLState: [5]byte{'H', 'Y', '0', '0', '0'},
					Message:  tt.message,
				},
			)
			assert.ErrorIs(
				t,
				"MapMySQLError(1030)",
				err,
				tt.wantErr,
			)
		})
	}
}

func TestMapMySQLErrorMapsOptionPreventsStatementReadOnly(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr error
	}{
		{"read_only", "The MySQL server is running in read only mode", persist.ErrStorageReadonly},
		{"super_read_only", "The MySQL server is running with super_read_only enabled", persist.ErrStorageReadonly},
		{"other", "This command is not allowed", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := &mysql.MySQLError{
				Number:   1290,
				SQLState: [5]byte{'H', 'Y', '0', '0', '0'},
				Message:  tt.message,
			}
			err := shared.MapMySQLError("ListLatest", source)
			if tt.wantErr == nil {
				if !errors.Is(err, source) {
					t.Fatalf("expected pass-through error, got %v", err)
				}
				return
			}
			assert.ErrorIs(t, "MapMySQLError(1290)", err, tt.wantErr)
		})
	}
}

func TestMapMySQLErrorMapsMySQLAndNetworkErrors(t *testing.T) {
	tests := []struct {
		name    string
		inErr   error
		wantErr error
	}{
		{"invalid_conn", mysql.ErrInvalidConn, persist.ErrConnectionLost},
		{"network_timeout", timeoutError{}, persist.ErrConnectionTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapMySQLError("ListLatest", tt.inErr)
			assert.ErrorIs(t, "MapMySQLError("+tt.name+")", err, tt.wantErr)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestMapMySQLErrorPassesThroughErrors(t *testing.T) {
	tests := []struct {
		name  string
		inErr error
	}{
		{"nil_error", nil},
		{"unknown_mysql_error", &mysql.MySQLError{Number: 9999}},
		{"unknown_error", errors.New("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapMySQLError("ListLatest", tt.inErr)
			if !errors.Is(err, tt.inErr) {
				t.Fatalf("err = got %v, want original error", err)
			}
		})
	}
}
