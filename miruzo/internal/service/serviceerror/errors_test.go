package serviceerror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

func TestMapPersistErrorMapsNotFound(t *testing.T) {
	err := serviceerror.MapPersistError(fmt.Errorf("not found: %w", persist.ErrNotFound))
	if !errors.Is(err, serviceerror.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
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
	conflictErr := serviceerror.MapPersistError(fmt.Errorf("conflict: %w", persist.ErrConflict))
	if !errors.Is(conflictErr, serviceerror.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", conflictErr)
	}

	uniqueErr := serviceerror.MapPersistError(fmt.Errorf("conflict: %w", persist.ErrUniqueViolation))
	if !errors.Is(uniqueErr, serviceerror.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", uniqueErr)
	}

	recoverableErr := serviceerror.MapPersistError(
		fmt.Errorf("recoverable conflict: %w", persist.ErrRecoverableConflict),
	)
	if !errors.Is(recoverableErr, serviceerror.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", recoverableErr)
	}

	foreignKeyReferencedErr := serviceerror.MapPersistError(
		fmt.Errorf("conflict: %w", persist.ErrForeignKeyReferenced),
	)
	if !errors.Is(foreignKeyReferencedErr, serviceerror.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", foreignKeyReferencedErr)
	}

	exclusionErr := serviceerror.MapPersistError(fmt.Errorf("conflict: %w", persist.ErrExclusionViolation))
	if !errors.Is(exclusionErr, serviceerror.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", exclusionErr)
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
