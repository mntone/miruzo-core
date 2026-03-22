package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/samber/mo"
)

var defaultCORSPolicy = corsPolicy{
	&corsPolicyBase{
		allowOrigins: map[string]struct{}{
			"https://example.com": {},
		},
	},
	CORSAllowMethods{
		http.MethodGet: CORSAllowHeaders{
			"content-type": {},
		},
	},
}

func TestCORSPassesThroughWhenNoOrigin(t *testing.T) {
	called := false
	handler := cors(defaultCORSPolicy, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/api/i/latest", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "called", called, true)
	assert.Equal(t, "status code", responseRecorder.Code, http.StatusOK)
	assert.Equal(
		t,
		"Access-Control-Allow-Origin",
		responseRecorder.Header().Get("Access-Control-Allow-Origin"),
		"",
	)
}

func TestCORSAddsHeadersForActualRequest(t *testing.T) {
	called := false
	handler := cors(defaultCORSPolicy, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/api/i/latest", nil)
	request.Header.Set("Origin", "https://example.com")
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "called", called, true)
	assert.Equal(t, "status code", responseRecorder.Code, http.StatusOK)
	assert.Equal(
		t,
		"Access-Control-Allow-Origin",
		responseRecorder.Header().Get("Access-Control-Allow-Origin"),
		"https://example.com",
	)
}

func TestCORSPassesThroughNonPreflightOptions(t *testing.T) {
	called := false
	handler := cors(defaultCORSPolicy, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodOptions, "/api/i/latest", nil)
	request.Header.Set("Origin", "https://example.com")
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "called", called, true)
	assert.Equal(t, "status code", responseRecorder.Code, http.StatusOK)
	assert.Equal(
		t,
		"Access-Control-Allow-Origin",
		responseRecorder.Header().Get("Access-Control-Allow-Origin"),
		"https://example.com",
	)
}

func TestCORSHandlesAllowedPreflightRequest(t *testing.T) {
	policy := corsPolicy{
		&corsPolicyBase{
			allowOrigins: map[string]struct{}{
				"https://example.com": {},
			},
			maxAge: mo.Some(10 * time.Minute),
		},
		CORSAllowMethods{
			http.MethodGet: CORSAllowHeaders{
				"content-type": {},
			},
			http.MethodPost: CORSAllowHeaders{
				"authorization": {},
				"content-type":  {},
			},
		},
	}

	called := false
	handler := cors(policy, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodOptions, "/api/i/latest", nil)
	request.Header.Set("Origin", "https://example.com")
	request.Header.Set("Access-Control-Request-Method", "POST")
	request.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "called", called, false)
	assert.Equal(t, "status code", responseRecorder.Code, http.StatusNoContent)
	assert.Equal(
		t,
		"Access-Control-Allow-Origin",
		responseRecorder.Header().Get("Access-Control-Allow-Origin"),
		"https://example.com",
	)
	assert.Equal(
		t,
		"Access-Control-Allow-Methods",
		responseRecorder.Header().Get("Access-Control-Allow-Methods"),
		"POST",
	)
	assert.Equal(
		t,
		"Access-Control-Allow-Headers",
		responseRecorder.Header().Get("Access-Control-Allow-Headers"),
		"Content-Type,Authorization",
	)
	assert.Equal(
		t,
		"Access-Control-Max-Age",
		responseRecorder.Header().Get("Access-Control-Max-Age"),
		"600",
	)
	assert.Equal(
		t,
		"Vary",
		responseRecorder.Header().Get("Vary"),
		"Origin,Access-Control-Request-Method,Access-Control-Request-Headers",
	)
}

func TestCORSRejectsDisallowedPreflightOrigin(t *testing.T) {
	called := false
	handler := cors(defaultCORSPolicy, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodOptions, "/api/i/latest", nil)
	request.Header.Set("Origin", "https://forbidden.example")
	request.Header.Set("Access-Control-Request-Method", "GET")
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "called", called, false)
	assert.Equal(t, "status code", responseRecorder.Code, http.StatusForbidden)
}

func TestCORSRejectsDisallowedPreflightMethod(t *testing.T) {
	tests := []struct {
		name    string
		methods []string
	}{
		{
			name:    "DisallowedMethod",
			methods: []string{"POST"},
		},
		{
			name:    "MultipleMethodsInSingleHeaderValue",
			methods: []string{"GET, POST"},
		},
		{
			name:    "MultipleMethodsInMultipleHeaderValues",
			methods: []string{"GET", "POST"},
		},
	}

	called := false
	handler := cors(defaultCORSPolicy, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	for _, tt := range tests {
		called = false
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodOptions, "/api/i/latest", nil)
			request.Header.Set("Origin", "https://example.com")
			for _, method := range tt.methods {
				request.Header.Add("Access-Control-Request-Method", method)
			}
			responseRecorder := httptest.NewRecorder()
			handler.ServeHTTP(responseRecorder, request)

			assert.Equal(t, "called", called, false)
			assert.Equal(t, "status code", responseRecorder.Code, http.StatusMethodNotAllowed)
		})
	}
}

func TestCORSRejectsDisallowedPreflightHeader(t *testing.T) {
	called := false
	handler := cors(defaultCORSPolicy, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodOptions, "/api/i/latest", nil)
	request.Header.Set("Origin", "https://example.com")
	request.Header.Set("Access-Control-Request-Method", "GET")
	request.Header.Set("Access-Control-Request-Headers", "X-Requested-With")
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "called", called, false)
	assert.Equal(t, "status code", responseRecorder.Code, http.StatusForbidden)
}
