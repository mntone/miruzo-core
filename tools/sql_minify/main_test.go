package main

import (
	"strings"
	"testing"
)

func TestMinifySQLKeepsSpaceBeforeDollarQuote(t *testing.T) {
	in := "CREATE FUNCTION f() RETURNS trigger AS $$\nBEGIN\n\tRETURN NEW;\nEND;\n$$ LANGUAGE plpgsql;"
	out := string(minifySQL([]byte(in), dialectPostgres))

	if strings.Contains(out, "AS$$") {
		t.Fatalf("unexpected minified SQL: %s", out)
	}
	if !strings.Contains(out, " AS $$") {
		t.Fatalf("expected space before $$ delimiter: %s", out)
	}
}

func TestMinifySQLKeepsSpaceBeforeTaggedDollarQuote(t *testing.T) {
	in := "CREATE FUNCTION f() RETURNS trigger AS $fn$\nBEGIN\n\tRETURN NEW;\nEND;\n$fn$ LANGUAGE plpgsql;"
	out := string(minifySQL([]byte(in), dialectPostgres))

	if strings.Contains(out, "AS$fn$") {
		t.Fatalf("unexpected minified SQL: %s", out)
	}
	if !strings.Contains(out, " AS $fn$") {
		t.Fatalf("expected space before tagged delimiter: %s", out)
	}
}

func TestMinifySQLKeepsSpaceBeforeDollarQuoteAfterComment(t *testing.T) {
	in := "CREATE FUNCTION f() RETURNS trigger AS -- comment\n$$\nBEGIN\n\tRETURN NEW;\nEND;\n$$ LANGUAGE plpgsql;"
	out := string(minifySQL([]byte(in), dialectPostgres))

	if strings.Contains(out, "AS$$") {
		t.Fatalf("unexpected minified SQL: %s", out)
	}
	if !strings.Contains(out, " AS $$") {
		t.Fatalf("expected space before $$ delimiter: %s", out)
	}
}

func TestMinifySQLKeepsSingleSpaceBeforeDollarQuoteFromManySpaces(t *testing.T) {
	in := "CREATE FUNCTION f() RETURNS trigger AS    $$\nBEGIN\n\tRETURN NEW;\nEND;\n$$ LANGUAGE plpgsql;"
	out := string(minifySQL([]byte(in), dialectPostgres))

	if strings.Contains(out, "AS$$") {
		t.Fatalf("unexpected minified SQL: %s", out)
	}
	if !strings.Contains(out, " AS $$") {
		t.Fatalf("expected normalized single space before $$ delimiter: %s", out)
	}
}
