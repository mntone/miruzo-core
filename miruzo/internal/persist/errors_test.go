package persist

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsRecoverableReturnsTrueForRecoverableErrors(t *testing.T) {
	recoverableUnavailable := fmt.Errorf("wrapped: %w", ErrRecoverableUnavailable)
	if !IsRecoverable(recoverableUnavailable) {
		t.Fatalf("expected recoverable for ErrRecoverableUnavailable")
	}

	recoverableConflict := fmt.Errorf("wrapped: %w", ErrRecoverableConflict)
	if !IsRecoverable(recoverableConflict) {
		t.Fatalf("expected recoverable for ErrRecoverableConflict")
	}
}

func TestIsRecoverableReturnsFalseForNonRecoverableErrors(t *testing.T) {
	tests := []error{
		nil,
		ErrUnavailable,
		ErrConflict,
		errors.New("unexpected"),
	}

	for _, tt := range tests {
		if IsRecoverable(tt) {
			t.Fatalf("did not expect recoverable: %v", tt)
		}
	}
}

func TestToTerminalErrorNilReturnsNil(t *testing.T) {
	if err := ToTerminalError(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestToTerminalErrorConvertsRecoverableUnavailable(t *testing.T) {
	source := fmt.Errorf("temporary unavailable: %w", ErrRecoverableUnavailable)
	err := ToTerminalError(source)

	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
	if errors.Is(err, ErrRecoverableUnavailable) {
		t.Fatalf("did not expect ErrRecoverableUnavailable, got %v", err)
	}
}

func TestToTerminalErrorConvertsRecoverableConflict(t *testing.T) {
	source := fmt.Errorf("serialization retry: %w", ErrRecoverableConflict)
	err := ToTerminalError(source)

	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if errors.Is(err, ErrRecoverableConflict) {
		t.Fatalf("did not expect ErrRecoverableConflict, got %v", err)
	}
}

func TestToTerminalErrorLeavesOtherErrorsUntouched(t *testing.T) {
	source := errors.New("unexpected")
	err := ToTerminalError(source)

	if !errors.Is(err, source) {
		t.Fatalf("expected original error, got %v", err)
	}
}
