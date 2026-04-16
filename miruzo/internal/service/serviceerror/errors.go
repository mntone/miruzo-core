package serviceerror

import (
	"errors"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

var (
	ErrNotFound             = errors.New("not found")             // 404 Not Found
	ErrConflict             = errors.New("conflict")              // 409 Conflict
	ErrUnprocessableContent = errors.New("unprocessable content") // 422 Unprocessable Content
	ErrTooManyRequests      = errors.New("too many requests")     // 429 Too Many Request
	ErrClientClosedRequest  = errors.New("client closed request") // 499 Client Closed Request

	ErrServiceUnavailable = errors.New("service unavailable") // 503 Service Unavailable
	ErrGatewayTimeout     = errors.New("gateway timeout")     // 504 Gateway Timeout
)

func MapPersistError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	// App errors
	case errors.Is(err, persist.ErrNoRows):
		return fmt.Errorf("%w: %v", ErrNotFound, err)
	case errors.Is(err, persist.ErrQuotaExceeded):
		return fmt.Errorf("%w: %v", ErrTooManyRequests, err)
	case errors.Is(err, persist.ErrQuotaUnderflow):
		return fmt.Errorf("%w: %v", ErrConflict, err)

	// Canceled and connection errors
	case errors.Is(err, persist.ErrContextCanceled):
		return fmt.Errorf("%w: %v", ErrClientClosedRequest, err)
	case errors.Is(err, persist.ErrDeadlineExceeded),
		errors.Is(err, persist.ErrQueryCanceled),
		errors.Is(err, persist.ErrConnectionTimeout):
		return fmt.Errorf("%w: %v", ErrGatewayTimeout, err)
	case errors.Is(err, persist.ErrConnection):
		return fmt.Errorf("%w: %v", ErrServiceUnavailable, err)

	// Constraint violations
	case errors.Is(err, persist.ErrConflict),
		errors.Is(err, persist.ErrExclusionViolation),
		errors.Is(err, persist.ErrForeignKeyReferenced),
		errors.Is(err, persist.ErrUniqueViolation):
		return fmt.Errorf("%w: %v", ErrConflict, err)
	case errors.Is(err, persist.ErrCheckViolation),
		errors.Is(err, persist.ErrForeignKeyReferenceNotFound),
		errors.Is(err, persist.ErrNotNullViolation),
		errors.Is(err, persist.ErrInvalidParam):
		return fmt.Errorf("%w: %v", ErrUnprocessableContent, err)

	// Too many requests
	case errors.Is(err, persist.ErrTooManyConnections):
		return fmt.Errorf("%w: %v", ErrTooManyRequests, err)

	// Contention, resource exhaustion and storage errors
	case errors.Is(err, persist.ErrContention),
		errors.Is(err, persist.ErrResourceExhausted),
		errors.Is(err, persist.ErrStorage):
		return fmt.Errorf("%w: %v", ErrServiceUnavailable, err)

	default:
		return err
	}
}
