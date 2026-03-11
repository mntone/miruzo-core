package assert

import (
	"reflect"
	"testing"
)

func NilArray[T any](t *testing.T, name string, val []T) {
	t.Helper()
	if val != nil {
		t.Fatalf("%s = \"%v\", want nil", name, val)
	}
}

func rlen(val any) int {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		panic("unsupported type")
	}

	switch v.Kind() {
	case reflect.Slice,
		reflect.Map,
		reflect.Array,
		reflect.String,
		reflect.Chan:
		return v.Len()
	}
	panic("unsupported type")
}

func NotEmpty(t *testing.T, name string, got any) {
	t.Helper()
	if rlen(got) == 0 {
		t.Fatalf("len(%s) = 0, want > 0", name)
	}
}

func Empty(t *testing.T, name string, got any) {
	t.Helper()
	if gotLen := rlen(got); gotLen != 0 {
		t.Fatalf("len(%s) = %d, want 0", name, gotLen)
	}
}

func LenIs(t *testing.T, name string, got any, wantLen int) {
	t.Helper()
	if gotLen := rlen(got); gotLen != wantLen {
		t.Fatalf("len(%s) = %d, want %d", name, gotLen, wantLen)
	}
}
