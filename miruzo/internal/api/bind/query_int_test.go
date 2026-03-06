package bind

import (
	"net/url"
	"testing"
)

func TestParseIntQueryWithDefaultReturnsDefaultWhenMissing(t *testing.T) {
	query := url.Values{}
	defaultValue := int16(-12)

	parsedValue, errs := ParseIntQueryWithDefault(query, "offset", defaultValue)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if parsedValue != defaultValue {
		t.Fatalf("expected %d, got %d", defaultValue, parsedValue)
	}
}

func TestParseIntQueryWithDefaultParsesValidValue(t *testing.T) {
	query := url.Values{
		"offset": []string{"-42"},
	}

	parsedValue, errs := ParseIntQueryWithDefault[int16](query, "offset", 0)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if parsedValue != -42 {
		t.Fatalf("expected -42, got %d", parsedValue)
	}
}

func TestParseIntQueryWithDefaultReturnsValidationErrorWhenInvalid(t *testing.T) {
	query := url.Values{
		"offset": []string{"12.3"},
	}

	_, errs := ParseIntQueryWithDefault[int16](query, "offset", 0)
	if len(errs) != 1 {
		t.Fatalf("expected a single error, got %v", errs)
	}
	if errs[0].Path != "query.offset" {
		t.Fatalf("unexpected path: %s", errs[0].Path)
	}
	if errs[0].Type != "invalid_type" {
		t.Fatalf("unexpected type: %s", errs[0].Type)
	}
	if errs[0].Message != "offset must be an integer" {
		t.Fatalf("unexpected message: %s", errs[0].Message)
	}
}

func TestParseIntQueryWithDefaultReturnsValidationErrorWhenOutOfRange(t *testing.T) {
	query := url.Values{
		"offset": []string{"200"},
	}

	_, errs := ParseIntQueryWithDefault[int8](query, "offset", 0)
	if len(errs) != 1 {
		t.Fatalf("expected a single error, got %v", errs)
	}
	if errs[0].Path != "query.offset" {
		t.Fatalf("unexpected path: %s", errs[0].Path)
	}
	if errs[0].Type != "invalid_type" {
		t.Fatalf("unexpected type: %s", errs[0].Type)
	}
	if errs[0].Message != "offset must be an integer" {
		t.Fatalf("unexpected message: %s", errs[0].Message)
	}
}
