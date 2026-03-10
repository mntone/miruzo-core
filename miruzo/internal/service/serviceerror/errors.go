package serviceerror

import (
	"context"
	"errors"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

var (
	ErrNotFound             = errors.New("not found")             // 404 Not Found
	ErrConflict             = errors.New("conflict")              // 409 Conflict
	ErrUnprocessableContent = errors.New("unprocessable content") // 422 Unprocessable Content

	ErrServiceUnavailable = errors.New("service unavailable") // 503 Service Unavailable
	ErrGatewayTimeout     = errors.New("gateway timeout")     // 504 Gateway Timeout
)

func MapPersistError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, context.Canceled):
		return err

	case errors.Is(err, persist.ErrNotFound):
		return fmt.Errorf("%w: %v", ErrNotFound, err)

	case errors.Is(err, persist.ErrConflict),
		errors.Is(err, persist.ErrRecoverableConflict),
		errors.Is(err, persist.ErrUniqueViolation),
		errors.Is(err, persist.ErrExclusionViolation),
		errors.Is(err, persist.ErrForeignKeyReferenced):
		return fmt.Errorf("%w: %v", ErrConflict, err)

	case errors.Is(err, persist.ErrNotNullViolation),
		errors.Is(err, persist.ErrCheckViolation),
		errors.Is(err, persist.ErrForeignKeyReferenceNotFound):
		return fmt.Errorf("%w: %v", ErrUnprocessableContent, err)

	case errors.Is(err, persist.ErrTimeout):
		return fmt.Errorf("%w: %v", ErrGatewayTimeout, err)

	case errors.Is(err, persist.ErrUnavailable),
		errors.Is(err, persist.ErrRecoverableUnavailable):
		return fmt.Errorf("%w: %v", ErrServiceUnavailable, err)

	default:
		return err
	}
}
