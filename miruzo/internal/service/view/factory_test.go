package view_test

import (
	"strings"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/service/view"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestNewReturnsErrorWhenVariantLayersBuilderIsNil(t *testing.T) {
	_, err := view.New(nil, nil, nil, score.Calculator{}, nil, nil)
	assert.Error(t, "New() error", err)
	if !strings.Contains(err.Error(), "variantLayersBuilder must not be nil") {
		t.Fatalf("New() error = %q, want to include %q", err.Error(), "variantLayersBuilder must not be nil")
	}
}

func TestNewReturnsServiceWhenVariantLayersBuilderIsPresent(t *testing.T) {
	builder := media.NewVariantLayerBuilder(media.VariantLayersSpec{})

	_, err := view.New(nil, nil, nil, score.Calculator{}, builder, nil)
	assert.NilError(t, "New() error", err)
}
