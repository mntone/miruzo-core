package middleware

import (
	"net/http"
	"slices"
	"strings"

	httperror "github.com/mntone/miruzo-core/miruzo/internal/api/http/error"
)

// RequireMethodOf ensures that the request method matches the expected method.
//
// If the method does not match, it returns 405 Method Not Allowed
// and sets the Allow header to the expected method.
func RequireMethodOf(
	method string,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != method {
			responseWriter.Header().Set("Allow", method)
			httperror.WriteMethodNotAllowed(responseWriter)
		} else {
			next(responseWriter, request)
		}
	}
}

// RequireMethodsOf ensures that the request methods match the expected method.
//
// If the method does not match, it returns 405 Method Not Allowed
// and sets the Allow header to the expected method.
func RequireMethodsOf(
	methods []string,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if !slices.Contains(methods, request.Method) {
			responseWriter.Header().Set("Allow", strings.Join(methods, ","))
			httperror.WriteMethodNotAllowed(responseWriter)
		} else {
			next(responseWriter, request)
		}
	}
}

// RequireMethodGet is a shorthand for RequireMethodOf(http.MethodGet, ...).
func RequireMethodGet(next http.HandlerFunc) http.HandlerFunc {
	return RequireMethodOf(http.MethodGet, next)
}

// RequireMethodPost is a shorthand for RequireMethodOf(http.MethodPost, ...).
func RequireMethodPost(next http.HandlerFunc) http.HandlerFunc {
	return RequireMethodOf(http.MethodPost, next)
}
