package imagelist_test

import (
	"strings"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNewReturnsErrorWhenVariantLayersBuilderIsNil(t *testing.T) {
	_, err := imagelist.New(nil, nil, 0, nil)
	assert.Error(t, "New() error", err)
	if !strings.Contains(err.Error(), "variantLayersBuilder must not be nil") {
		t.Fatalf("New() error = %q, want to include %q", err.Error(), "variantLayersBuilder must not be nil")
	}
}

func TestNewReturnsServiceWhenVariantLayersBuilderIsPresent(t *testing.T) {
	builder := media.NewVariantLayerBuilder(media.VariantLayersSpec{})

	_, err := imagelist.New(nil, nil, 0, builder)
	assert.NilError(t, "New() error", err)
}
