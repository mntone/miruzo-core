package dberrors_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestSQLiteToPersistMapsContextErrors(t *testing.T) {
	tests := []struct {
		name    string
		inErr   error
		wantErr error
	}{
		{"context_canceled", context.Canceled, persist.ErrContextCanceled},
		{"context_deadline_exceeded", context.DeadlineExceeded, persist.ErrDeadlineExceeded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrors.ToPersist("ListLatest", tt.inErr)
			assert.ErrorIs(t, "ToPersist("+tt.name+")", err, tt.wantErr)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestSQLiteToPersistMapsByCode(t *testing.T) {
	tests := []struct {
		name    string
		code    sqlite3.ErrNo
		wantErr error
	}{
		{"interrupt", sqlite3.ErrInterrupt, persist.ErrQueryCanceled},
		{"too_big", sqlite3.ErrTooBig, persist.ErrCheckViolation},
		{"mismatch", sqlite3.ErrMismatch, persist.ErrInvalidParam},
		{"busy", sqlite3.ErrBusy, persist.ErrResourceBusy},
		{"locked", sqlite3.ErrLocked, persist.ErrResourceBusy},
		{"no_memory", sqlite3.ErrNomem, persist.ErrOutOfMemory},
		{"disk_full", sqlite3.ErrFull, persist.ErrStorageFull},
		{"readonly", sqlite3.ErrReadonly, persist.ErrStorageReadonly},
		{"io_error", sqlite3.ErrIoErr, persist.ErrStorageUnavailable},
		{"cant_open", sqlite3.ErrCantOpen, persist.ErrStorageUnavailable},
		{"corrupt", sqlite3.ErrCorrupt, persist.ErrStorageCorrupted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrors.ToPersist(
				"ListLatest",
				sqlite3.Error{Code: tt.code},
			)
			assert.ErrorIs(t, "ToPersist("+tt.name+")", err, tt.wantErr)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestSQLiteToPersistMapsConstraintByExtendedCode(t *testing.T) {
	tests := []struct {
		name         string
		extendedCode sqlite3.ErrNoExtended
		wantErr      error
	}{
		{"check", sqlite3.ErrConstraintCheck, persist.ErrCheckViolation},
		{"foreign_key", sqlite3.ErrConstraintForeignKey, persist.ErrForeignKeyReferenceNotFound},
		{"not_null", sqlite3.ErrConstraintNotNull, persist.ErrNotNullViolation},
		{"primary_key", sqlite3.ErrConstraintPrimaryKey, persist.ErrUniqueViolation},
		{"unique", sqlite3.ErrConstraintUnique, persist.ErrUniqueViolation},
		{"rowid", sqlite3.ErrConstraintRowID, persist.ErrUniqueViolation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrors.ToPersist(
				"ListLatest",
				sqlite3.Error{
					Code:         sqlite3.ErrConstraint,
					ExtendedCode: tt.extendedCode,
				},
			)
			assert.ErrorIs(t, "ToPersist("+tt.name+")", err, tt.wantErr)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestSQLiteToPersistDeleteErrorMapsDifferences(t *testing.T) {
	err := dberrors.ToPersistDelete(
		"DeleteImage",
		sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintForeignKey},
	)
	assert.ErrorIs(
		t,
		"ToPersistDelete(foreign_key)",
		err,
		persist.ErrForeignKeyReferenced,
	)
	if !strings.Contains(err.Error(), "operation=DeleteImage") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestSQLiteToPersistMapsByMessageWithDB(t *testing.T) {
	tests := []struct {
		name    string
		build   func(*sql.DB) error
		wantErr error
	}{
		{
			name: "integer_overflow",
			build: func(db *sql.DB) error {
				_, err := db.Exec("SELECT abs(-9223372036854775808)")
				return err
			},
			wantErr: persist.ErrCheckViolation,
		},
		{
			name: "syntax_error",
			build: func(db *sql.DB) error {
				_, err := db.Exec("SELEC 1")
				return err
			},
			wantErr: persist.ErrSyntax,
		},
		{
			name: "too_many_columns",
			build: func(db *sql.DB) error {
				columns := make([]string, 2100)
				for i := range columns {
					columns[i] = strconv.Itoa(i)
				}
				_, err := db.Exec("SELECT " + strings.Join(columns, ","))
				return err
			},
			wantErr: persist.ErrInvalidStatement,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sql.Open("sqlite3", ":memory:")
			assert.NilError(t, "sql.Open()", err)
			defer func() {
				assert.NilError(t, "db.Close()", db.Close())
			}()

			rawErr := tt.build(db)
			assert.Error(t, "rawErr", rawErr)

			err = dberrors.ToPersist("ListLatest", rawErr)
			assert.ErrorIs(t, fmt.Sprintf("ToPersist(%s)", tt.name), err, tt.wantErr)
			if !strings.Contains(err.Error(), "operation=ListLatest") {
				t.Fatalf("expected operation detail, got %v", err)
			}
		})
	}
}

func TestSQLiteToPersistMapsAbortRollbackToTxAborted(t *testing.T) {
	err := dberrors.ToPersist(
		"ListLatest",
		sqlite3.Error{
			Code:         sqlite3.ErrAbort,
			ExtendedCode: sqlite3.ErrAbortRollback,
		},
	)
	assert.ErrorIs(
		t,
		"ToPersist(abort_rollback)",
		err,
		persist.ErrTxAborted,
	)
	if !strings.Contains(err.Error(), "operation=ListLatest") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestSQLiteToPersistPassesThrough(t *testing.T) {
	if err := dberrors.ToPersist("ListLatest", nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	source := sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintVTab}
	err := dberrors.ToPersist("ListLatest", source)
	if !errors.Is(err, source) {
		t.Fatalf("expected pass-through sqlite error, got %v", err)
	}

	generic := errors.New("unknown")
	err = dberrors.ToPersist("ListLatest", generic)
	if !errors.Is(err, generic) {
		t.Fatalf("expected pass-through generic error, got %v", err)
	}
}
