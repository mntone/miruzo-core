package httperror

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

func TestWriteServiceErrorWritesNothingForCanceled(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	WriteServiceError(responseRecorder, context.Canceled)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected default status 200, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		t.Fatalf("unexpected body: %q", responseRecorder.Body.String())
	}
}

func TestWriteServiceErrorReturnsConflictForConflict(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	WriteServiceError(
		responseRecorder,
		fmt.Errorf("duplicate key: %w", serviceerror.ErrConflict),
	)

	if responseRecorder.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "{\"type\":\"conflict\"}" {
		t.Fatalf("unexpected body: %q", responseRecorder.Body.String())
	}
}

func TestWriteServiceErrorReturnsServiceUnavailableForUnavailable(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	WriteServiceError(
		responseRecorder,
		fmt.Errorf("query failed: %w", serviceerror.ErrServiceUnavailable),
	)

	if responseRecorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "{\"type\":\"service_unavailable\"}" {
		t.Fatalf("unexpected body: %q", responseRecorder.Body.String())
	}
}

func TestWriteServiceErrorReturnsUnprocessableContent(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	WriteServiceError(
		responseRecorder,
		fmt.Errorf("invalid state: %w", serviceerror.ErrUnprocessableContent),
	)

	if responseRecorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "{\"type\":\"unprocessable_content\"}" {
		t.Fatalf("unexpected body: %q", responseRecorder.Body.String())
	}
}

func TestWriteServiceErrorReturnsGatewayTimeoutForTimeout(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	WriteServiceError(
		responseRecorder,
		fmt.Errorf("query timed out: %w", serviceerror.ErrGatewayTimeout),
	)

	if responseRecorder.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected status 504, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "{\"type\":\"gateway_timeout\"}" {
		t.Fatalf("unexpected body: %q", responseRecorder.Body.String())
	}
}

func TestWriteServiceErrorReturnsGatewayTimeoutForDeadlineExceeded(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	WriteServiceError(responseRecorder, context.DeadlineExceeded)

	if responseRecorder.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected status 504, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "{\"type\":\"gateway_timeout\"}" {
		t.Fatalf("unexpected body: %q", responseRecorder.Body.String())
	}
}

func TestWriteServiceErrorReturnsInternalServerErrorForUnknownError(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	WriteServiceError(responseRecorder, fmt.Errorf("unknown failure"))

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "{\"type\":\"internal_server_error\"}" {
		t.Fatalf("unexpected body: %q", responseRecorder.Body.String())
	}
}
