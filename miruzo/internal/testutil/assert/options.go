package assert

import "testing"

type optional interface {
	IsPresent() bool
	IsAbsent() bool
}

func IsPresent(t *testing.T, name string, got any) {
	t.Helper()
	switch x := got.(type) {
	case optional:
		if x.IsAbsent() {
			t.Fatalf("%s is absent, want non-empty", name)
		}
	default:
		if x == nil {
			t.Fatalf("%s is nil, want non-empty", name)
		}
	}
}

func IsAbsent(t *testing.T, name string, got any) {
	t.Helper()
	switch x := got.(type) {
	case optional:
		if x.IsPresent() {
			t.Fatalf("%s is present, want empty", name)
		}
	default:
		if x != nil {
			t.Fatalf("%s is not nil, want nil", name)
		}
	}
}
