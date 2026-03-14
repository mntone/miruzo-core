package shared_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestMapPostgreErrorReturnsNilForNilInput(t *testing.T) {
	err := shared.MapPostgreError("ListLatest", nil)
	assert.NilError(t, "MapPostgreError(nil)", err)
}

func TestMapPostgreDeleteErrorReturnsNilForNilInput(t *testing.T) {
	err := shared.MapPostgreDeleteError("DeleteImage", nil)
	assert.NilError(t, "MapPostgreDeleteError(nil)", err)
}

func TestMapPostgreErrorMapsQueryCanceledToTimeout(t *testing.T) {
	err := shared.MapPostgreError(
		"ListLatest",
		fmt.Errorf("query failed: %w", &pgconn.PgError{Code: "57014"}),
	)
	assert.ErrorIs(
		t,
		"MapPostgreError(57014)",
		err,
		persist.ErrTimeout,
	)
	if !strings.Contains(err.Error(), "operation=ListLatest") {
		t.Fatalf("expected operation detail, got %v", err)
	}
	if !strings.Contains(err.Error(), "sqlstate=57014") {
		t.Fatalf("expected sqlstate detail, got %v", err)
	}
}

func TestMapPostgreErrorMapsConnectionAndCapacitySQLStates(t *testing.T) {
	tests := []struct {
		name     string
		sqlState string
		wantErr  error
	}{
		{
			name:     "sqlclient_unable_to_establish_sqlconnection",
			sqlState: "08001",
			wantErr:  persist.ErrUnavailable,
		},
		{
			name:     "connection_does_not_exist",
			sqlState: "08003",
			wantErr:  persist.ErrRecoverableUnavailable,
		},
		{
			name:     "sqlserver_rejected_establishment_of_sqlconnection",
			sqlState: "08004",
			wantErr:  persist.ErrUnavailable,
		},
		{
			name:     "connection_failure",
			sqlState: "08006",
			wantErr:  persist.ErrUnavailable,
		},
		{
			name:     "out_of_memory",
			sqlState: "53200",
			wantErr:  persist.ErrUnavailable,
		},
		{
			name:     "too_many_connections",
			sqlState: "53300",
			wantErr:  persist.ErrRecoverableUnavailable,
		},
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
				err,
				tt.wantErr,
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

func TestMapPostgreErrorMapsRecoverableConflictSQLStates(t *testing.T) {
	tests := []struct {
		name     string
		sqlState string
	}{
		{
			name:     "serialization_failure",
			sqlState: "40001",
		},
		{
			name:     "deadlock_detected",
			sqlState: "40P01",
		},
		{
			name:     "lock_not_available",
			sqlState: "55P03",
		},
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
				err,
				persist.ErrRecoverableConflict,
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

func TestMapPostgreErrorMapsViolations(t *testing.T) {
	tests := []struct {
		name     string
		sqlState string
		wantErr  error
	}{
		{
			name:     "not_null_violation",
			sqlState: "23502",
			wantErr:  persist.ErrNotNullViolation,
		},
		{
			name:     "foreign_key_violation",
			sqlState: "23503",
			wantErr:  persist.ErrForeignKeyReferenceNotFound,
		},
		{
			name:     "unique_violation",
			sqlState: "23505",
			wantErr:  persist.ErrUniqueViolation,
		},
		{
			name:     "check_violation",
			sqlState: "23514",
			wantErr:  persist.ErrCheckViolation,
		},
		{
			name:     "exclusion_violation",
			sqlState: "23P01",
			wantErr:  persist.ErrExclusionViolation,
		},
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
				err,
				tt.wantErr,
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

func TestMapPostgreErrorOtherErrorIsNotMapped(t *testing.T) {
	err := shared.MapPostgreError(
		"ListLatest",
		&pgconn.PgError{Code: "22001"},
	)

	if errors.Is(err, persist.ErrRecoverableConflict) {
		t.Fatalf("did not expect ErrRecoverableConflict, got %v", err)
	}
	if errors.Is(err, persist.ErrUnavailable) {
		t.Fatalf("did not expect ErrUnavailable, got %v", err)
	}
	if errors.Is(err, persist.ErrRecoverableUnavailable) {
		t.Fatalf("did not expect ErrRecoverableUnavailable, got %v", err)
	}
	if errors.Is(err, persist.ErrTimeout) {
		t.Fatalf("did not expect ErrTimeout, got %v", err)
	}
	if errors.Is(err, persist.ErrCheckViolation) {
		t.Fatalf("did not expect ErrCheckViolation, got %v", err)
	}
	if errors.Is(err, persist.ErrExclusionViolation) {
		t.Fatalf("did not expect ErrExclusionViolation, got %v", err)
	}
	if errors.Is(err, persist.ErrForeignKeyReferenceNotFound) {
		t.Fatalf("did not expect ErrForeignKeyReferenceNotFound, got %v", err)
	}
	if errors.Is(err, persist.ErrForeignKeyReferenced) {
		t.Fatalf("did not expect ErrForeignKeyReferenced, got %v", err)
	}
	if errors.Is(err, persist.ErrNotNullViolation) {
		t.Fatalf("did not expect ErrNotNullViolation, got %v", err)
	}
	if errors.Is(err, persist.ErrUniqueViolation) {
		t.Fatalf("did not expect ErrUniqueViolation, got %v", err)
	}

	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) {
		t.Fatalf("expected *pgconn.PgError, got %T", err)
	}
	if pgError.Code != "22001" {
		t.Fatalf("unexpected sqlstate: %s", pgError.Code)
	}
}
