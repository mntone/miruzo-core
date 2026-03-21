package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpmedia "github.com/mntone/miruzo-core/miruzo/internal/api/http/media"
	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func TestRequireAcceptAnyOf(t *testing.T) {
	tests := []struct {
		name          string
		accept        string
		allowed       []string
		wantStatus    int
		wantNextCalls int
	}{
		{
			name:          "empty accept allows",
			accept:        "",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNoContent,
			wantNextCalls: 1,
		},
		{
			name:          "exact match allows",
			accept:        "application/json",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNoContent,
			wantNextCalls: 1,
		},
		{
			name:          "json with q allows",
			accept:        "application/json;q=0.8",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNoContent,
			wantNextCalls: 1,
		},
		{
			name:          "q zero rejects",
			accept:        "application/json;q=0",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNotAcceptable,
			wantNextCalls: 0,
		},
		{
			name:          "specific q zero wins over wildcard",
			accept:        "application/json;q=0,*/*;q=1",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNotAcceptable,
			wantNextCalls: 0,
		},
		{
			name:          "type wildcard allows",
			accept:        "application/*;q=0.5",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNoContent,
			wantNextCalls: 1,
		},
		{
			name:          "all wildcard allows",
			accept:        "*/*;q=0.1",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNoContent,
			wantNextCalls: 1,
		},
		{
			name:          "unsupported media type rejects",
			accept:        "text/plain",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNotAcceptable,
			wantNextCalls: 0,
		},
		{
			name:          "invalid q segment ignored",
			accept:        "application/json;q=abc,text/plain",
			allowed:       []string{"application/json"},
			wantStatus:    http.StatusNotAcceptable,
			wantNextCalls: 0,
		},
		{
			name:          "multiple allowed types picks acceptable one",
			accept:        "application/xml;q=0.2",
			allowed:       []string{"application/json", "application/xml"},
			wantStatus:    http.StatusNoContent,
			wantNextCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalls := 0
			parsed := httpmedia.ParseMediaTypes(tt.allowed)
			handler := m.RequireAcceptAnyOf(parsed, func(responseWriter http.ResponseWriter, request *http.Request) {
				nextCalls++
				responseWriter.WriteHeader(http.StatusNoContent)
			})

			request := httptest.NewRequest(http.MethodGet, "http://example.com/api/test", nil)
			if tt.accept != "" {
				request.Header.Set("Accept", tt.accept)
			}
			responseRecorder := httptest.NewRecorder()

			handler(responseRecorder, request)

			if nextCalls != tt.wantNextCalls {
				t.Fatalf("next calls = %d, want %d", nextCalls, tt.wantNextCalls)
			}
			if responseRecorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", responseRecorder.Code, tt.wantStatus)
			}
		})
	}
}

func TestRequireAcceptJson(t *testing.T) {
	nextCalls := 0
	handler := m.RequireAcceptJson(func(responseWriter http.ResponseWriter, request *http.Request) {
		nextCalls++
		responseWriter.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "http://example.com/api/test", nil)
	request.Header.Set("Accept", "application/json;q=0.7")
	responseRecorder := httptest.NewRecorder()

	handler(responseRecorder, request)

	if nextCalls != 1 {
		t.Fatalf("next calls = %d, want %d", nextCalls, 1)
	}
	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", responseRecorder.Code, http.StatusNoContent)
	}
}
