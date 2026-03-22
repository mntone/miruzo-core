package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/samber/mo"
)

func getNormalizedOrigin(headers http.Header) (string, bool) {
	rawOrigins := headers.Values("Origin")
	if len(rawOrigins) != 1 {
		return "", false
	}

	origin := strings.TrimSpace(rawOrigins[0])
	return origin, true
}

type requestMethodState int32

const (
	requestMethodStateEmpty requestMethodState = iota
	requestMethodStateValid
	requestMethodStateMultiple
)

func getNormalizedRequestMethod(headers http.Header) (string, requestMethodState) {
	rawMethods := headers.Values("Access-Control-Request-Method")
	if len(rawMethods) == 0 {
		return "", requestMethodStateEmpty
	}
	if len(rawMethods) >= 2 {
		return "", requestMethodStateMultiple
	}

	method := strings.ToUpper(strings.TrimSpace(rawMethods[0]))
	return method, requestMethodStateValid
}

func isPreflightRequest(request *http.Request) bool {
	if request.Method != http.MethodOptions {
		return false
	}

	_, hasRequestMethod := getNormalizedRequestMethod(request.Header)
	return hasRequestMethod != requestMethodStateEmpty
}

func splitRequestHeaders(headers string) []string {
	var result []string
	for rawHeader := range strings.SplitSeq(headers, ",") {
		header := strings.TrimSpace(rawHeader)
		if header == "" {
			continue
		}

		result = append(result, header)
	}
	return result
}

type corsPolicyBase struct {
	allowOrigins map[string]struct{}
	maxAge       mo.Option[time.Duration]
}

func (p corsPolicyBase) resolveAllowOrigin(origin string) (string, bool) {
	if origin == "" {
		return "", false
	}

	if p.allowOrigins == nil {
		return "*", true
	}

	if _, ok := p.allowOrigins[origin]; !ok {
		return "", false
	}

	return origin, true
}

type CORSAllowHeaders map[string]struct{}

func (allow CORSAllowHeaders) validateRequestHeaders(requestHeaders []string) bool {
	if allow == nil {
		if len(requestHeaders) == 0 {
			return true
		}

		return false
	}

	for _, header := range requestHeaders {
		lowerHeader := strings.ToLower(header)
		if _, ok := allow[lowerHeader]; !ok {
			return false
		}
	}

	return true
}

type CORSAllowMethods map[string]CORSAllowHeaders

type corsPolicy struct {
	*corsPolicyBase
	allowMethods CORSAllowMethods
}

func (p corsPolicy) resolveAllowMethod(method string) (CORSAllowHeaders, bool) {
	if method == "" {
		return nil, false
	}

	if p.allowMethods == nil {
		return nil, true
	}

	headers, ok := p.allowMethods[method]
	return headers, ok
}

func cors(policy corsPolicy, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Resolve the request origin.
		origin, hasOrigin := getNormalizedOrigin(request.Header)
		if !hasOrigin {
			if isPreflightRequest(request) {
				responseWriter.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(responseWriter, request)
			return
		}

		// Validate the origin against the configured policy.
		allowOrigin, isOriginAllowed := policy.resolveAllowOrigin(origin)
		if !isOriginAllowed {
			if isPreflightRequest(request) {
				responseWriter.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(responseWriter, request)
			return
		}

		// For non-preflight requests, add CORS headers and continue.
		if request.Method != http.MethodOptions {
			h := responseWriter.Header()
			h.Set("Access-Control-Allow-Origin", allowOrigin)
			h.Set("Vary", "Origin")
			next.ServeHTTP(responseWriter, request)
			return
		}

		// Resolve the requested method from preflight headers.
		requestMethod, requestMethodState := getNormalizedRequestMethod(request.Header)
		switch requestMethodState {
		case requestMethodStateEmpty:
			h := responseWriter.Header()
			h.Set("Access-Control-Allow-Origin", allowOrigin)
			h.Set("Vary", "Origin")
			next.ServeHTTP(responseWriter, request)
			return
		case requestMethodStateMultiple:
			responseWriter.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Validate the requested method and get the allowed headers.
		allowHeaders, isMethodAllowed := policy.resolveAllowMethod(requestMethod)
		if !isMethodAllowed {
			responseWriter.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Validate requested headers from Access-Control-Request-Headers.
		requestHeaders := splitRequestHeaders(request.Header.Get("Access-Control-Request-Headers"))
		if !allowHeaders.validateRequestHeaders(requestHeaders) {
			responseWriter.WriteHeader(http.StatusForbidden)
			return
		}

		h := responseWriter.Header()
		h.Set("Access-Control-Allow-Origin", allowOrigin)
		h.Set("Access-Control-Allow-Methods", requestMethod)
		h.Set("Access-Control-Allow-Headers", strings.Join(requestHeaders, ","))
		if maxAge, present := policy.maxAge.Get(); present && maxAge > 0 {
			seconds := int64(maxAge / time.Second)
			if seconds > 0 {
				h.Set("Access-Control-Max-Age", strconv.FormatInt(seconds, 10))
			}
		}
		h.Set("Vary", "Origin,Access-Control-Request-Method,Access-Control-Request-Headers")
		responseWriter.WriteHeader(http.StatusNoContent)
	})
}

type CORSFactory struct {
	*corsPolicyBase
	getPolicy  corsPolicy
	postPolicy corsPolicy
}

func NewCORSFactory(origins []string, maxAge *time.Duration) *CORSFactory {
	var allowOrigins map[string]struct{}
	if len(origins) != 0 {
		allowOrigins = make(map[string]struct{}, len(origins))
		for _, origin := range origins {
			normalizedOrigin := strings.TrimSpace(origin)
			if normalizedOrigin == "" {
				continue
			}
			if normalizedOrigin == "*" {
				allowOrigins = nil
				break
			}
			allowOrigins[normalizedOrigin] = struct{}{}
		}
	}

	base := &corsPolicyBase{
		allowOrigins: allowOrigins,
		maxAge:       mo.PointerToOption(maxAge),
	}
	return &CORSFactory{
		corsPolicyBase: base,
		getPolicy: corsPolicy{
			corsPolicyBase: base,
			allowMethods: CORSAllowMethods{
				"GET": CORSAllowHeaders{},
			},
		},
		postPolicy: corsPolicy{
			corsPolicyBase: base,
			allowMethods: CORSAllowMethods{
				"POST": CORSAllowHeaders{},
			},
		},
	}
}

func (fty *CORSFactory) New(allowMethods CORSAllowMethods, next http.HandlerFunc) http.HandlerFunc {
	policy := corsPolicy{
		corsPolicyBase: fty.corsPolicyBase,
		allowMethods:   allowMethods,
	}
	return cors(policy, next)
}

func (fty *CORSFactory) GET(next http.HandlerFunc) http.HandlerFunc {
	return cors(fty.getPolicy, next)
}

func (fty *CORSFactory) POST(next http.HandlerFunc) http.HandlerFunc {
	return cors(fty.postPolicy, next)
}
