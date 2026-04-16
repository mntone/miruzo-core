package shared

import (
	"errors"
	"strings"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type timeoutMarkerError struct{}

func (timeoutMarkerError) Error() string {
	return "timeout marker"
}

func TestMapPostgreErrorMapsPgconnTimeoutByHook(t *testing.T) {
	original := isPgconnTimeoutError
	isPgconnTimeoutError = func(err error) bool {
		_, ok := err.(timeoutMarkerError)
		return ok
	}
	t.Cleanup(func() {
		isPgconnTimeoutError = original
	})

	err := MapPostgreError("ListLatest", timeoutMarkerError{})
	assert.ErrorIs(
		t,
		"MapPostgreError(pgconn.Timeout)",
		err,
		persist.ErrConnectionTimeout,
	)
	if !strings.Contains(err.Error(), "operation=ListLatest") {
		t.Fatalf("expected operation detail, got %v", err)
	}
}

func TestMapPostgreErrorPassesThroughWhenNotTimeoutByHook(t *testing.T) {
	original := isPgconnTimeoutError
	isPgconnTimeoutError = func(error) bool { return false }
	t.Cleanup(func() {
		isPgconnTimeoutError = original
	})

	source := errors.New("unknown")
	err := MapPostgreError("ListLatest", source)
	if !errors.Is(err, source) {
		t.Fatalf("expected pass-through error, got %v", err)
	}
}
