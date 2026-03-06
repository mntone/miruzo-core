package assert_test

import "testing"

type testOptional struct {
	present bool
}

func (o testOptional) IsPresent() bool {
	return o.present
}

func (o testOptional) IsAbsent() bool {
	return !o.present
}

func requirePass(t *testing.T, name string, fn func(t *testing.T)) {
	t.Helper()
	if ok := t.Run(name, fn); !ok {
		t.Fatalf("%s: subtest failed, want pass", name)
	}
}
