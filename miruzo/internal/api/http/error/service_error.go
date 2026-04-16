package error

import (
	"errors"
	"net/http"
	"unsafe"

	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

const (
	jsonNotFound             = "{\"type\":\"not_found\"}"
	jsonMethodNotAllowed     = "{\"type\":\"method_not_allowed\"}"
	jsonNotAcceptable        = "{\"type\":\"not_acceptable\"}"
	jsonConflict             = "{\"type\":\"conflict\"}"
	jsonUnprocessableContent = "{\"type\":\"unprocessable_content\"}"
	jsonTooManyRequests      = "{\"type\":\"too_many_requests\"}"

	jsonInternalServerError = "{\"type\":\"internal_server_error\"}"
	jsonServiceUnavailable  = "{\"type\":\"service_unavailable\"}"
	jsonGatewayTimeout      = "{\"type\":\"gateway_timeout\"}"
)

func getUnsafeByteSlice(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func WriteNotFound(responseWriter http.ResponseWriter) {
	_ = response.WriteJSONBytes(
		responseWriter,
		http.StatusNotFound,
		getUnsafeByteSlice(jsonNotFound),
	)
}

func WriteMethodNotAllowed(responseWriter http.ResponseWriter) {
	_ = response.WriteJSONBytes(
		responseWriter,
		http.StatusMethodNotAllowed,
		getUnsafeByteSlice(jsonMethodNotAllowed),
	)
}

func WriteNotAcceptable(responseWriter http.ResponseWriter) {
	_ = response.WriteJSONBytes(
		responseWriter,
		http.StatusNotAcceptable,
		getUnsafeByteSlice(jsonNotAcceptable),
	)
}

func WriteInternalServerError(responseWriter http.ResponseWriter) {
	_ = response.WriteJSONBytes(
		responseWriter,
		http.StatusInternalServerError,
		getUnsafeByteSlice(jsonInternalServerError),
	)
}

func WriteServiceError(responseWriter http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, serviceerror.ErrClientClosedRequest):
		return

	case errors.Is(err, serviceerror.ErrNotFound):
		_ = response.WriteJSONBytes(
			responseWriter,
			http.StatusNotFound,
			getUnsafeByteSlice(jsonNotFound),
		)

	case errors.Is(err, serviceerror.ErrServiceUnavailable):
		_ = response.WriteJSONBytes(
			responseWriter,
			http.StatusServiceUnavailable,
			getUnsafeByteSlice(jsonServiceUnavailable),
		)

	case errors.Is(err, serviceerror.ErrConflict):
		_ = response.WriteJSONBytes(
			responseWriter,
			http.StatusConflict,
			getUnsafeByteSlice(jsonConflict),
		)

	case errors.Is(err, serviceerror.ErrUnprocessableContent):
		_ = response.WriteJSONBytes(
			responseWriter,
			http.StatusUnprocessableEntity,
			getUnsafeByteSlice(jsonUnprocessableContent),
		)

	case errors.Is(err, serviceerror.ErrTooManyRequests):
		_ = response.WriteJSONBytes(
			responseWriter,
			http.StatusTooManyRequests,
			getUnsafeByteSlice(jsonTooManyRequests),
		)

	case errors.Is(err, serviceerror.ErrGatewayTimeout):
		_ = response.WriteJSONBytes(
			responseWriter,
			http.StatusGatewayTimeout,
			getUnsafeByteSlice(jsonGatewayTimeout),
		)

	default:
		WriteInternalServerError(responseWriter)
	}
}
