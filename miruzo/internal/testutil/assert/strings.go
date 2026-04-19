package assert

import (
	"strings"
	"testing"
)

func Contains(t *testing.T, name string, gotVal, wantSubstr string) {
	t.Helper()
	if !strings.Contains(gotVal, wantSubstr) {
		t.Fatalf("%s = %v, want to contain %v", name, gotVal, wantSubstr)
	}
}
