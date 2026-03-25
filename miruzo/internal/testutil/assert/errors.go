package assert

import (
	"errors"
	"testing"
)

func Error(t *testing.T, name string, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s = nil, want error", name)
	}
}

func NilError(t *testing.T, name string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s = \"%v\", want nil", name, err)
	}
}

func ErrorIs(t *testing.T, name string, gotErr error, wantErr error) {
	t.Helper()
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("%s = \"%v\", want \"%v\"", name, gotErr, wantErr)
	}
}
