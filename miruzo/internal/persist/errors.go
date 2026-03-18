package persist

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrConflict               = errors.New("conflict")
	ErrRecoverableConflict    = errors.New("recoverable conflict")
	ErrUnavailable            = errors.New("unavailable")
	ErrRecoverableUnavailable = errors.New("recoverable unavailable")
	ErrTimeout                = errors.New("timeout")
	ErrQuotaExceeded          = errors.New("quota exceeded")
	ErrQuotaUnderflow         = errors.New("quota underflow")

	ErrCheckViolation              = errors.New("check violation")
	ErrForeignKeyReferenceNotFound = errors.New("foreign key reference not found")
	ErrForeignKeyReferenced        = errors.New("foreign key referenced")
	ErrExclusionViolation          = errors.New("exclusion violation")
	ErrNotNullViolation            = errors.New("not null violation")
	ErrUniqueViolation             = errors.New("unique violation")
)

func IsRecoverable(err error) bool {
	return errors.Is(err, ErrRecoverableUnavailable) || errors.Is(err, ErrRecoverableConflict)
}

func ToTerminalError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrRecoverableUnavailable):
		return fmt.Errorf("%w: %s", ErrUnavailable, err)

	case errors.Is(err, ErrRecoverableConflict):
		return fmt.Errorf("%w: %s", ErrConflict, err)

	default:
		return err
	}
}
