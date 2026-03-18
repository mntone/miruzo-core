package serviceerror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestMapPersistErrorMapsNotFound(t *testing.T) {
	err := serviceerror.MapPersistError(fmt.Errorf("not found: %w", persist.ErrNotFound))
	if !errors.Is(err, serviceerror.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMapPersistErrorMapsQuotaExceeded(t *testing.T) {
	err := serviceerror.MapPersistError(fmt.Errorf("not found: %w", persist.ErrQuotaExceeded))
	if !errors.Is(err, serviceerror.ErrTooManyRequests) {
		t.Fatalf("expected ErrTooManyRequests, got %v", err)
	}
}

func TestMapPersistErrorMapsTimeout(t *testing.T) {
	err := serviceerror.MapPersistError(fmt.Errorf("query timeout: %w", persist.ErrTimeout))
	if !errors.Is(err, serviceerror.ErrGatewayTimeout) {
		t.Fatalf("expected ErrGatewayTimeout, got %v", err)
	}
}

func TestMapPersistErrorMapsUnavailable(t *testing.T) {
	err := serviceerror.MapPersistError(fmt.Errorf("unavailable: %w", persist.ErrUnavailable))
	if !errors.Is(err, serviceerror.ErrServiceUnavailable) {
		t.Fatalf("expected ErrServiceUnavailable, got %v", err)
	}

	recoverableErr := serviceerror.MapPersistError(
		fmt.Errorf("recoverable unavailable: %w", persist.ErrRecoverableUnavailable),
	)
	if !errors.Is(recoverableErr, serviceerror.ErrServiceUnavailable) {
		t.Fatalf("expected ErrServiceUnavailable, got %v", recoverableErr)
	}
}

func TestMapPersistErrorMapsConflict(t *testing.T) {
	tests := []error{
		fmt.Errorf("conflict: %w", persist.ErrConflict),
		fmt.Errorf("recoverable conflict: %w", persist.ErrRecoverableConflict),
		fmt.Errorf("unique violation: %w", persist.ErrUniqueViolation),
		fmt.Errorf("exclusion violation: %w", persist.ErrExclusionViolation),
		fmt.Errorf("foreign key referenced: %w", persist.ErrForeignKeyReferenced),
		fmt.Errorf("quota underflow: %w", persist.ErrQuotaUnderflow),
	}

	for _, tt := range tests {
		gotErr := serviceerror.MapPersistError(tt)
		assert.ErrorIs(t, gotErr.Error(), gotErr, serviceerror.ErrConflict)
	}
}

func TestMapPersistErrorMapsUnprocessableContent(t *testing.T) {
	checkErr := serviceerror.MapPersistError(fmt.Errorf("check: %w", persist.ErrCheckViolation))
	if !errors.Is(checkErr, serviceerror.ErrUnprocessableContent) {
		t.Fatalf("expected ErrUnprocessableContent for check, got %v", checkErr)
	}

	notNullErr := serviceerror.MapPersistError(fmt.Errorf("not null: %w", persist.ErrNotNullViolation))
	if !errors.Is(notNullErr, serviceerror.ErrUnprocessableContent) {
		t.Fatalf("expected ErrUnprocessableContent for not null, got %v", notNullErr)
	}

	foreignKeyReferenceNotFoundErr := serviceerror.MapPersistError(
		fmt.Errorf(
			"foreign key reference not found: %w",
			persist.ErrForeignKeyReferenceNotFound,
		),
	)
	if !errors.Is(foreignKeyReferenceNotFoundErr, serviceerror.ErrUnprocessableContent) {
		t.Fatalf(
			"expected ErrUnprocessableContent for foreign key reference not found, got %v",
			foreignKeyReferenceNotFoundErr,
		)
	}
}

func TestMapPersistErrorPassesThroughUnknownError(t *testing.T) {
	source := errors.New("unknown error")
	err := serviceerror.MapPersistError(source)
	if !errors.Is(err, source) {
		t.Fatalf("expected original error, got %v", err)
	}
}
