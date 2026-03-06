package bind

import (
	"net/url"
	"testing"
	"time"
)

func TestParseTimeQueryWithDefaultReturnsDefaultWhenMissing(t *testing.T) {
	query := url.Values{}
	defaultValue := time.Date(2026, 3, 2, 1, 2, 3, 400000000, time.UTC)

	parsedValue, errs := ParseTimeQueryWithDefault(query, "cursor", defaultValue)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if !parsedValue.Equal(defaultValue) {
		t.Fatalf("expected %v, got %v", defaultValue, parsedValue)
	}
}

func TestParseTimeQueryWithDefaultParsesValidValue(t *testing.T) {
	query := url.Values{
		"cursor": []string{"2026-03-02T10:20:30.123456Z"},
	}
	expected := time.Date(2026, 3, 2, 10, 20, 30, 123456000, time.UTC)

	parsedValue, errs := ParseTimeQueryWithDefault(query, "cursor", time.Time{})
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if !parsedValue.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, parsedValue)
	}
}

func TestParseTimeQueryWithDefaultReturnsValidationErrorWhenInvalid(t *testing.T) {
	query := url.Values{
		"cursor": []string{"not-a-time"},
	}

	_, errs := ParseTimeQueryWithDefault(query, "cursor", time.Time{})
	if len(errs) != 1 {
		t.Fatalf("expected a single error, got %v", errs)
	}
	if errs[0].Path != "query.cursor" {
		t.Fatalf("unexpected path: %s", errs[0].Path)
	}
	if errs[0].Type != "invalid_type" {
		t.Fatalf("unexpected type: %s", errs[0].Type)
	}
	expectedMessage := "cursor must be a UTC timestamp in the format 2006-01-02T15:04:05.999999Z"
	if errs[0].Message != expectedMessage {
		t.Fatalf("unexpected message: %s", errs[0].Message)
	}
}
