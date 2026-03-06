package middleware

import (
	"net/http"
	"strings"

	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
)

// splitHeaderList splits a comma-separated HTTP header list, ignoring empty segments.
// It does not implement full RFC parsing for quoted strings, but is sufficient for typical Accept headers.
func splitHeaderList(headerValue string) []string {
	parts := strings.Split(headerValue, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

// RequireAcceptAnyOf ensures that the request Accept header allows
// at least one of the specified media types.
//
// Rules:
// - If Accept header is missing or empty: allow (treated as "*/*").
// - If Accept contains "*/*": allow.
// - If Accept contains one of allowedMediaTypes (ignoring parameters): allow.
// - Otherwise: return 406 Not Acceptable.
func RequireAcceptAnyOf(
	allowedMediaTypes []string,
	next http.HandlerFunc,
) http.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowedMediaTypes))
	for _, mediaType := range allowedMediaTypes {
		allowedSet[strings.ToLower(strings.TrimSpace(mediaType))] = struct{}{}
	}

	return func(responseWriter http.ResponseWriter, request *http.Request) {
		acceptHeaderValue := strings.TrimSpace(request.Header.Get("Accept"))
		if acceptHeaderValue == "" {
			next(responseWriter, request)
			return
		}

		for _, acceptPart := range splitHeaderList(acceptHeaderValue) {
			mediaType := strings.ToLower(strings.TrimSpace(acceptPart))
			if mediaType == "" {
				continue
			}

			// Strip parameters, e.g. "application/json; q=0.9"
			if semicolonIndex := strings.IndexByte(mediaType, ';'); semicolonIndex >= 0 {
				mediaType = strings.TrimSpace(mediaType[:semicolonIndex])
			}

			if mediaType == "*/*" {
				next(responseWriter, request)
				return
			}

			if _, ok := allowedSet[mediaType]; ok {
				next(responseWriter, request)
				return
			}
		}

		response.WriteJSONText(
			responseWriter,
			http.StatusNotAcceptable,
			"{\"type\":\"not_acceptable\"}",
		)
	}
}

// RequireAcceptJson is a shorthand for requiring "application/json".
func RequireAcceptJson(next http.HandlerFunc) http.HandlerFunc {
	return RequireAcceptAnyOf(
		[]string{"application/json"},
		next,
	)
}
