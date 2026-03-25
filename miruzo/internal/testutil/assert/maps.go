package assert

import (
	"maps"
	"testing"
)

func EqualMap[M ~map[K]V, K, V comparable](t *testing.T, name string, got, want M) {
	t.Helper()
	if !maps.Equal(got, want) {
		t.Fatalf("%s = \"%v\", want \"%v\"", name, got, want)
	}
}
