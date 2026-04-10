package assert

import (
	"cmp"
	"testing"
)

func Equal[T comparable](t *testing.T, name string, gotVal, wantVal T) {
	t.Helper()
	if gotVal != wantVal {
		t.Fatalf("%s = %v, want %v", name, gotVal, wantVal)
	}
}

type equatable[T any] interface {
	Equal(other T) bool
}

func EqualFn[T any, E equatable[T]](t *testing.T, name string, gotVal E, wantVal T) {
	t.Helper()
	if !gotVal.Equal(wantVal) {
		t.Fatalf("%s = %v, want %v", name, gotVal, wantVal)
	}
}

func GreaterThan[T cmp.Ordered](t *testing.T, name string, gotVal, wantVal T) {
	t.Helper()
	if gotVal <= wantVal {
		t.Fatalf("%s <= %v, want > %v", name, gotVal, wantVal)
	}
}

func LessThan[T cmp.Ordered](t *testing.T, name string, gotVal, wantVal T) {
	t.Helper()
	if gotVal >= wantVal {
		t.Fatalf("%s >= %v, want < %v", name, gotVal, wantVal)
	}
}
