package httperror

import (
	"context"
	"errors"
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

func WriteNotFound(responseWriter http.ResponseWriter) {
	_ = response.WriteJSONText(
		responseWriter,
		http.StatusNotFound,
		"{\"type\":\"not_found\"}",
	)
}

func WriteInternalServerError(responseWriter http.ResponseWriter) {
	_ = response.WriteJSONText(
		responseWriter,
		http.StatusInternalServerError,
		"{\"type\":\"internal_server_error\"}",
	)
}

func WriteServiceError(responseWriter http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		return

	case errors.Is(err, serviceerror.ErrNotFound):
		_ = response.WriteJSONText(
			responseWriter,
			http.StatusNotFound,
			"{\"type\":\"not_found\"}",
		)

	case errors.Is(err, serviceerror.ErrServiceUnavailable):
		_ = response.WriteJSONText(
			responseWriter,
			http.StatusServiceUnavailable,
			"{\"type\":\"service_unavailable\"}",
		)

	case errors.Is(err, serviceerror.ErrConflict):
		_ = response.WriteJSONText(
			responseWriter,
			http.StatusConflict,
			"{\"type\":\"conflict\"}",
		)

	case errors.Is(err, serviceerror.ErrUnprocessableContent):
		_ = response.WriteJSONText(
			responseWriter,
			http.StatusUnprocessableEntity,
			"{\"type\":\"unprocessable_content\"}",
		)

	case errors.Is(err, serviceerror.ErrGatewayTimeout), errors.Is(err, context.DeadlineExceeded):
		_ = response.WriteJSONText(
			responseWriter,
			http.StatusGatewayTimeout,
			"{\"type\":\"gateway_timeout\"}",
		)

	default:
		WriteInternalServerError(responseWriter)
	}
}
