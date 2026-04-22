package persistence

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	adapterpersistence "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type OpenAdminHandleFunc func(
	ctx context.Context,
) (adapterpersistence.DatabaseAdminHandle, error)

var adminTestSequence uint64

func NewAdminTestDatabaseName(prefix string) string {
	sequence := atomic.AddUint64(&adminTestSequence, 1)
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), sequence)
}

func RunAdminHandleLifecycle(
	t *testing.T,
	openHandle OpenAdminHandleFunc,
) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	adminHandle, err := openHandle(ctx)
	assert.NilError(t, "OpenAdminHandle()", err)
	t.Cleanup(func() {
		_ = adminHandle.Drop(context.Background())
		_ = adminHandle.Close()
	})

	exists, err := adminHandle.Exists(ctx)
	assert.NilError(t, "Exists() before Create error", err)
	if exists {
		err := adminHandle.Drop(ctx)
		assert.NilError(t, "Drop() pre-cleanup error", err)

		exists, err = adminHandle.Exists(ctx)
		assert.NilError(t, "Exists() after pre-cleanup Drop error", err)
		assert.Equal(t, "Exists() after pre-cleanup Drop", exists, false)
	}

	err = adminHandle.Create(ctx)
	assert.NilError(t, "Create() error", err)

	exists, err = adminHandle.Exists(ctx)
	assert.NilError(t, "Exists() after Create error", err)
	assert.Equal(t, "Exists() after Create", exists, true)

	err = adminHandle.Drop(ctx)
	assert.NilError(t, "Drop() error", err)

	exists, err = adminHandle.Exists(ctx)
	assert.NilError(t, "Exists() after Drop error", err)
	assert.Equal(t, "Exists() after Drop", exists, false)
}
