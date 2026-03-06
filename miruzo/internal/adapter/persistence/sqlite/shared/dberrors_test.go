package shared_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestMapSQLiteErrorReturnsNilForNilInput(t *testing.T) {
	err := shared.MapSQLiteError("ListLatest", nil)
	assert.NilError(t, "MapSQLiteError(nil)", err)
}

func TestMapSQLiteDeleteErrorReturnsNilForNilInput(t *testing.T) {
	err := shared.MapSQLiteDeleteError("DeleteImage", nil)
	assert.NilError(t, "MapSQLiteDeleteError(nil)", err)
}

func TestMapSQLiteErrorDeadlineExceededMapsToTimeout(t *testing.T) {
	err := shared.MapSQLiteError(
		"ListLatest",
		context.DeadlineExceeded,
	)
	assert.ErrorIs(
		t,
		"MapSQLiteError(context.DeadlineExceeded)",
		err,
		persist.ErrTimeout,
	)
	if !strings.Contains(err.Error(), "operation=ListLatest") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestMapSQLiteErrorInterruptMapsToTimeout(t *testing.T) {
	err := shared.MapSQLiteError(
		"ListRecently",
		sqlite3.Error{
			Code: sqlite3.ErrInterrupt,
		},
	)
	assert.ErrorIs(
		t,
		"MapSQLiteError(sqlite3.ErrInterrupt)",
		err,
		persist.ErrTimeout,
	)
	if !strings.Contains(err.Error(), "operation=ListRecently") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestMapSQLiteErrorBusyMapsToRecoverableConflict(t *testing.T) {
	err := shared.MapSQLiteError(
		"ListLatest",
		sqlite3.Error{
			Code:         sqlite3.ErrBusy,
			ExtendedCode: sqlite3.ErrBusySnapshot,
		},
	)
	assert.ErrorIs(
		t,
		"MapSQLiteError(sqlite3.ErrBusySnapshot)",
		err,
		persist.ErrRecoverableConflict,
	)
	if !strings.Contains(err.Error(), "operation=ListLatest") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestMapSQLiteErrorLockedMapsToRecoverableConflict(t *testing.T) {
	err := shared.MapSQLiteError(
		"ListRecently",
		sqlite3.Error{
			Code:         sqlite3.ErrLocked,
			ExtendedCode: sqlite3.ErrLockedSharedCache,
		},
	)
	assert.ErrorIs(
		t,
		"MapSQLiteError(sqlite3.ErrLockedSharedCache)",
		err,
		persist.ErrRecoverableConflict,
	)
	if !strings.Contains(err.Error(), "operation=ListRecently") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestMapSQLiteErrorConstraintsMapsViolations(t *testing.T) {
	tests := []struct {
		name         string
		extendedCode sqlite3.ErrNoExtended
		wantErr      error
	}{
		{
			name:         "check",
			extendedCode: sqlite3.ErrConstraintCheck,
			wantErr:      persist.ErrCheckViolation,
		},
		{
			name:         "foreign_key",
			extendedCode: sqlite3.ErrConstraintForeignKey,
			wantErr:      persist.ErrForeignKeyReferenceNotFound,
		},
		{
			name:         "not_null",
			extendedCode: sqlite3.ErrConstraintNotNull,
			wantErr:      persist.ErrNotNullViolation,
		},
		{
			name:         "primary_key",
			extendedCode: sqlite3.ErrConstraintPrimaryKey,
			wantErr:      persist.ErrUniqueViolation,
		},
		{
			name:         "unique",
			extendedCode: sqlite3.ErrConstraintUnique,
			wantErr:      persist.ErrUniqueViolation,
		},
		{
			name:         "rowid",
			extendedCode: sqlite3.ErrConstraintRowID,
			wantErr:      persist.ErrUniqueViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.MapSQLiteError(
				"ListLatest",
				sqlite3.Error{
					Code:         sqlite3.ErrConstraint,
					ExtendedCode: tt.extendedCode,
				},
			)
			assert.ErrorIs(
				t,
				fmt.Sprintf("MapSQLiteError(%v)", sqlite3.Error{
					Code:         sqlite3.ErrConstraint,
					ExtendedCode: tt.extendedCode,
				}),
				err,
				tt.wantErr,
			)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestMapSQLiteDeleteErrorMapsForeignKeyViolationToReferenced(t *testing.T) {
	err := shared.MapSQLiteDeleteError(
		"DeleteImage",
		sqlite3.Error{
			Code:         sqlite3.ErrConstraint,
			ExtendedCode: sqlite3.ErrConstraintForeignKey,
		},
	)
	assert.ErrorIs(
		t,
		"MapSQLiteDeleteError(sqlite3.ErrConstraintForeignKey)",
		err,
		persist.ErrForeignKeyReferenced,
	)
	if !strings.Contains(err.Error(), "operation=DeleteImage") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestMapSQLiteErrorOtherErrorIsNotMapped(t *testing.T) {
	err := shared.MapSQLiteError(
		"ListLatest",
		sqlite3.Error{
			Code:         sqlite3.ErrConstraint,
			ExtendedCode: sqlite3.ErrConstraintVTab,
		},
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

	var sqliteErr sqlite3.Error
	if !errors.As(err, &sqliteErr) {
		t.Fatalf("expected sqlite3.Error, got %T", err)
	}
	if sqliteErr.Code != sqlite3.ErrConstraint {
		t.Fatalf("unexpected code: %d", sqliteErr.Code)
	}
}
