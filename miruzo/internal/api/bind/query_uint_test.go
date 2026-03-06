package bind

import (
	"net/url"
	"testing"
)

func TestParseUintQueryWithDefaultReturnsDefaultWhenMissing(t *testing.T) {
	query := url.Values{}
	defaultValue := uint16(12)

	parsedValue, errs := ParseUintQueryWithDefault(query, "limit", defaultValue)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if parsedValue != defaultValue {
		t.Fatalf("expected %d, got %d", defaultValue, parsedValue)
	}
}

func TestParseUintQueryWithDefaultParsesValidValue(t *testing.T) {
	query := url.Values{
		"limit": []string{"42"},
	}

	parsedValue, errs := ParseUintQueryWithDefault[uint16](query, "limit", 0)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if parsedValue != 42 {
		t.Fatalf("expected 42, got %d", parsedValue)
	}
}

func TestParseUintQueryWithDefaultReturnsValidationErrorWhenInvalid(t *testing.T) {
	query := url.Values{
		"limit": []string{"-1"},
	}

	_, errs := ParseUintQueryWithDefault[uint16](query, "limit", 0)
	if len(errs) != 1 {
		t.Fatalf("expected a single error, got %v", errs)
	}
	if errs[0].Path != "query.limit" {
		t.Fatalf("unexpected path: %s", errs[0].Path)
	}
	if errs[0].Type != "invalid_type" {
		t.Fatalf("unexpected type: %s", errs[0].Type)
	}
	if errs[0].Message != "limit must be an integer" {
		t.Fatalf("unexpected message: %s", errs[0].Message)
	}
}

func TestParseUintQueryWithDefaultReturnsValidationErrorWhenOutOfRange(t *testing.T) {
	query := url.Values{
		"limit": []string{"300"},
	}

	_, errs := ParseUintQueryWithDefault[uint8](query, "limit", 0)
	if len(errs) != 1 {
		t.Fatalf("expected a single error, got %v", errs)
	}
	if errs[0].Path != "query.limit" {
		t.Fatalf("unexpected path: %s", errs[0].Path)
	}
	if errs[0].Type != "invalid_type" {
		t.Fatalf("unexpected type: %s", errs[0].Type)
	}
	if errs[0].Message != "limit must be an integer" {
		t.Fatalf("unexpected message: %s", errs[0].Message)
	}
}
