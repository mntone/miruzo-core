package assert

import "testing"

func Nil[T any](t *testing.T, name string, val *T) {
	t.Helper()
	if val != nil {
		t.Fatalf("%s = %v, want nil", name, val)
	}
}

func NotNil[T any](t *testing.T, name string, val *T) {
	t.Helper()
	if val == nil {
		t.Fatalf("%s = nil, want non-nil", name)
	}
}
