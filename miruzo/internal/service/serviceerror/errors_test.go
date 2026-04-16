package serviceerror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestMapPersistErrorMapsErrors(t *testing.T) {
	tests := []struct {
		name    string
		inErr   error
		wantErr error
	}{
		// App errors
		{"conflict", fmt.Errorf("x: %w", persist.ErrConflict), serviceerror.ErrConflict},
		{"no_rows", fmt.Errorf("x: %w", persist.ErrNoRows), serviceerror.ErrNotFound},
		{"quota_exceeded", fmt.Errorf("x: %w", persist.ErrQuotaExceeded), serviceerror.ErrTooManyRequests},
		{"quota_underflow", fmt.Errorf("x: %w", persist.ErrQuotaUnderflow), serviceerror.ErrConflict},

		// Canceled errors
		{"context_canceled", fmt.Errorf("x: %w", persist.ErrContextCanceled), serviceerror.ErrClientClosedRequest},
		{"deadline_exceeded", fmt.Errorf("x: %w", persist.ErrDeadlineExceeded), serviceerror.ErrGatewayTimeout},
		{"query_canceled", fmt.Errorf("x: %w", persist.ErrQueryCanceled), serviceerror.ErrGatewayTimeout},

		// Connection errors
		{"conn_init", fmt.Errorf("x: %w", persist.ErrConnectionInit), serviceerror.ErrServiceUnavailable},
		{"conn_lost", fmt.Errorf("x: %w", persist.ErrConnectionLost), serviceerror.ErrServiceUnavailable},
		{"conn_refused", fmt.Errorf("x: %w", persist.ErrConnectionRefused), serviceerror.ErrServiceUnavailable},
		{"conn_timeout", fmt.Errorf("x: %w", persist.ErrConnectionTimeout), serviceerror.ErrGatewayTimeout},
		{"conn_unavailable", fmt.Errorf("x: %w", persist.ErrConnectionUnavailable), serviceerror.ErrServiceUnavailable},

		// Constraint violations
		{"check", fmt.Errorf("x: %w", persist.ErrCheckViolation), serviceerror.ErrUnprocessableContent},
		{"exclusion", fmt.Errorf("x: %w", persist.ErrExclusionViolation), serviceerror.ErrConflict},
		{"fk_not_found", fmt.Errorf("x: %w", persist.ErrForeignKeyReferenceNotFound), serviceerror.ErrUnprocessableContent},
		{"fk_referenced", fmt.Errorf("x: %w", persist.ErrForeignKeyReferenced), serviceerror.ErrConflict},
		{"not_null", fmt.Errorf("x: %w", persist.ErrNotNullViolation), serviceerror.ErrUnprocessableContent},
		{"unique", fmt.Errorf("x: %w", persist.ErrUniqueViolation), serviceerror.ErrConflict},
		{"invalid_param", fmt.Errorf("x: %w", persist.ErrInvalidParam), serviceerror.ErrUnprocessableContent},

		// Contention errors
		{"deadlock", fmt.Errorf("x: %w", persist.ErrDeadlockDetected), serviceerror.ErrServiceUnavailable},
		{"lock_timeout", fmt.Errorf("x: %w", persist.ErrLockTimeout), serviceerror.ErrServiceUnavailable},
		{"lock_unavailable", fmt.Errorf("x: %w", persist.ErrLockUnavailable), serviceerror.ErrServiceUnavailable},
		{"resource_busy", fmt.Errorf("x: %w", persist.ErrResourceBusy), serviceerror.ErrServiceUnavailable},
		{"tx_serialization", fmt.Errorf("x: %w", persist.ErrTxSerialization), serviceerror.ErrServiceUnavailable},

		// Resource exhaustion
		{"out_of_memory", fmt.Errorf("x: %w", persist.ErrOutOfMemory), serviceerror.ErrServiceUnavailable},
		{"resource_exhausted", fmt.Errorf("x: %w", persist.ErrResourceExhausted), serviceerror.ErrServiceUnavailable},
		{"storage_full", fmt.Errorf("x: %w", persist.ErrStorageFull), serviceerror.ErrServiceUnavailable},
		{"too_many_connections", fmt.Errorf("x: %w", persist.ErrTooManyConnections), serviceerror.ErrTooManyRequests},

		// Storage errors
		{"storage_corrupted", fmt.Errorf("x: %w", persist.ErrStorageCorrupted), serviceerror.ErrServiceUnavailable},
		{"storage_unavailable", fmt.Errorf("x: %w", persist.ErrStorageUnavailable), serviceerror.ErrServiceUnavailable},

		// Syntax errors
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serviceerror.MapPersistError(tt.inErr)
			if tt.wantErr == nil {
				assert.ErrorIs(t, "err", got, tt.inErr)
			} else {
				assert.ErrorIs(t, "err", got, tt.wantErr)
			}
		})
	}
}

func TestMapPersistErrorPassesThroughErrors(t *testing.T) {
	tests := []struct {
		name  string
		inErr error
	}{
		// Nil error
		{"nil error", nil},

		// Syntax errors
		{"syntax_error", fmt.Errorf("x: %w", persist.ErrSyntax)},
		{"invalid_statement", fmt.Errorf("x: %w", persist.ErrInvalidStatement)},

		// Unknown error
		{"unknown error", errors.New("unknown error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := serviceerror.MapPersistError(tt.inErr)
			if !errors.Is(err, tt.inErr) {
				t.Fatalf("err = got %v, want original error", err)
			}
		})
	}
}
