package imagelist

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	testutilDomain "github.com/mntone/miruzo-core/miruzo/internal/testutil/domain"
)

func TestListBaseReturnsLimitedItemsAndNextCursor(t *testing.T) {
	base := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	spec := 42

	gotSpec := 0
	result, err := listBase(
		context.Background(),
		2,
		[]media.ImageFormat{},
		func(requestContext context.Context, loadSpec int) ([]persist.ImageWithCursor[time.Time], error) {
			gotSpec = loadSpec
			return []persist.ImageWithCursor[time.Time]{
				{
					Image:  persist.Image{IngestID: 1001},
					Cursor: base.Add(3 * time.Hour),
				},
				{
					Image:  persist.Image{IngestID: 1002},
					Cursor: base.Add(2 * time.Hour),
				},
				{
					Image:  persist.Image{IngestID: 1003},
					Cursor: base.Add(1 * time.Hour),
				},
			}, nil
		},
		spec,
		testutilDomain.NewTestVariantLayersBuilder(),
		backoff.NoRetryPolicy{},
	)
	assert.NilError(t, "listBase() error", err)
	if gotSpec != spec {
		t.Fatalf("load spec = %d, want %d", gotSpec, spec)
	}

	assert.LenIs(t, "result.Items", result.Items, 2)
	if got := result.Items[0].IngestID; got != 1001 {
		t.Fatalf("result.Items[0].IngestID = %d, want 1001", got)
	}
	if got := result.Items[1].IngestID; got != 1002 {
		t.Fatalf("result.Items[1].IngestID = %d, want 1002", got)
	}

	assert.IsPresent(t, "result.Cursor", result.Cursor)
	assert.EqualFn(t, "result.Cursor", result.Cursor.MustGet(), base.Add(2*time.Hour))
}

func TestListBaseReturnsNoNextCursorWhenNoMoreItems(t *testing.T) {
	base := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	result, err := listBase(
		context.Background(),
		2,
		[]media.ImageFormat{},
		func(requestContext context.Context, loadSpec int) ([]persist.ImageWithCursor[time.Time], error) {
			return []persist.ImageWithCursor[time.Time]{
				{
					Image:  persist.Image{IngestID: 2001},
					Cursor: base.Add(2 * time.Hour),
				},
				{
					Image:  persist.Image{IngestID: 2002},
					Cursor: base.Add(1 * time.Hour),
				},
			}, nil
		},
		0,
		testutilDomain.NewTestVariantLayersBuilder(),
		backoff.NoRetryPolicy{},
	)
	assert.NilError(t, "listBase() error", err)
	assert.LenIs(t, "result.Items", result.Items, 2)
	assert.IsAbsent(t, "result.Cursor", result.Cursor)
}

func TestListBaseMapsPersistErrorToServiceError(t *testing.T) {

	_, err := listBase(
		context.Background(),
		2,
		[]media.ImageFormat{},
		func(requestContext context.Context, loadSpec int) ([]persist.ImageWithCursor[time.Time], error) {
			return nil, fmt.Errorf("database unavailable: %w", persist.ErrUnavailable)
		},
		0,
		testutilDomain.NewTestVariantLayersBuilder(),
		backoff.NoRetryPolicy{},
	)
	assert.ErrorIs(t, "listBase() error", err, serviceerror.ErrServiceUnavailable)
}

func TestListBaseReturnsContextCanceledWithoutCallingLoad(t *testing.T) {
	requestContext, cancel := context.WithCancel(context.Background())
	cancel()

	loadCalled := false
	_, err := listBase(
		requestContext,
		2,
		[]media.ImageFormat{},
		func(innerContext context.Context, loadSpec int) ([]persist.ImageWithCursor[time.Time], error) {
			loadCalled = true
			return nil, nil
		},
		0,
		testutilDomain.NewTestVariantLayersBuilder(),
		backoff.NoRetryPolicy{},
	)
	assert.ErrorIs(t, "listBase() error", err, context.Canceled)
	if loadCalled {
		t.Fatal("load was called, want not called")
	}
}
