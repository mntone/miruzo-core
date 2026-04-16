package persist

import (
	"errors"
	"testing"
)

func TestIsRecoverableReturnsTrueForRecoverableErrors(t *testing.T) {
	tests := []error{
		ErrConnectionTimeout,
		ErrConnectionLost,
		ErrDeadlockDetected,
		ErrLockTimeout,
		ErrLockUnavailable,
		ErrResourceBusy,
		ErrTxSerialization,
	}

	for _, tt := range tests {
		if !IsRecoverable(tt) {
			t.Fatalf("expect recoverable: %v", tt)
		}
	}
}

func TestIsRecoverableReturnsFalseForNonRecoverableErrors(t *testing.T) {
	tests := []error{
		nil,
		ErrTooManyConnections,
		ErrStorageCorrupted,
		ErrStorageUnavailable,
		ErrSyntax,
		ErrInvalidStatement,
		errors.New("unexpected"),
	}

	for _, tt := range tests {
		if IsRecoverable(tt) {
			t.Fatalf("did not expect recoverable: %v", tt)
		}
	}
}
