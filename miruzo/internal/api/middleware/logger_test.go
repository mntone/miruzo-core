package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestLogWithLogger(t *testing.T) {
	var buffer bytes.Buffer
	testLogger := log.New(&buffer, "", 0)

	handler := RequestLogWithLogger(
		http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			responseWriter.WriteHeader(http.StatusNotAcceptable)
		}),
		testLogger,
	)

	request := httptest.NewRequest(http.MethodGet, "http://example.com/api/i/latest?limit=10", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	got := buffer.String()
	want := "[miruzo-api] 192.0.2.1:1234 - \"GET /api/i/latest\" 406 Not Acceptable\n"
	if got != want {
		t.Fatalf("log = %q, want %q", got, want)
	}
}

func TestRequestLogWithLogger_DefaultStatusOK(t *testing.T) {
	var buffer bytes.Buffer
	testLogger := log.New(&buffer, "", 0)

	handler := RequestLogWithLogger(
		http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {}),
		testLogger,
	)

	request := httptest.NewRequest(http.MethodPost, "http://example.com/", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	got := buffer.String()
	want := "[miruzo-api] 192.0.2.1:1234 - \"POST /\" 200 OK\n"
	if got != want {
		t.Fatalf("log = %q, want %q", got, want)
	}
}
