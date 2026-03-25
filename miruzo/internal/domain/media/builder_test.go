package media

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestVariantLayersBuilderGroupVariantsByLayer(t *testing.T) {
	builder := NewVariantLayerBuilder(VariantLayersSpec{
		{
			LayerID: 2,
			Variants: []VariantSpec{
				{LayerID: 2, Width: 320},
			},
		},
		{
			LayerID: 1,
			Variants: []VariantSpec{
				{LayerID: 1, Width: 320},
			},
		},
		{
			LayerID: FallbackLayerID,
			Variants: []VariantSpec{
				{LayerID: FallbackLayerID, Width: 320},
			},
		},
	})

	got := builder.GroupVariantsByLayer([]Variant{
		{LayerID: 1, Width: 640, RelativePath: "l1w640"},
		{LayerID: 999, Width: 777, RelativePath: "unknown"},
		{LayerID: 2, Width: 960, RelativePath: "l2w960"},
		{LayerID: 1, Width: 320, RelativePath: "l1w320"},
		{LayerID: 2, Width: 480, RelativePath: "l2w480"},
	})

	assert.LenIs(t, "grouped layers", got, 2)

	assert.Equal(t, "grouped[0][0].LayerID", got[0][0].LayerID, LayerIDType(2))
	assert.Equal(t, "grouped[0][0].Width", got[0][0].Width, uint16(480))
	assert.Equal(t, "grouped[0][1].Width", got[0][1].Width, uint16(960))

	assert.Equal(t, "grouped[1][0].LayerID", got[1][0].LayerID, LayerIDType(1))
	assert.Equal(t, "grouped[1][0].Width", got[1][0].Width, uint16(320))
	assert.Equal(t, "grouped[1][1].Width", got[1][1].Width, uint16(640))
}

func TestVariantLayersBuilderGroupVariantsByLayerReturnsEmptyForNoKnownLayers(t *testing.T) {
	builder := NewVariantLayerBuilder(VariantLayersSpec{
		{
			LayerID: 1,
			Variants: []VariantSpec{
				{LayerID: 1, Width: 320},
			},
		},
	})

	got := builder.GroupVariantsByLayer([]Variant{
		{LayerID: 999, Width: 777},
	})

	assert.Empty(t, "grouped layers", got)
}

func TestVariantLayersBuilderGroupVariantsByLayerIncludesFallbackLayer(t *testing.T) {
	builder := NewVariantLayerBuilder(VariantLayersSpec{
		{
			LayerID: 1,
			Variants: []VariantSpec{
				{LayerID: 1, Width: 320},
			},
		},
		{
			LayerID: FallbackLayerID,
			Variants: []VariantSpec{
				{LayerID: FallbackLayerID, Width: 320},
			},
		},
	})

	got := builder.GroupVariantsByLayer([]Variant{
		{LayerID: FallbackLayerID, Width: 320, RelativePath: "fallback"},
		{LayerID: 1, Width: 640, RelativePath: "l1w640"},
	})

	assert.LenIs(t, "grouped layers", got, 2)
	assert.Equal(t, "grouped[0][0].LayerID", got[0][0].LayerID, LayerIDType(1))
	assert.Equal(t, "grouped[1][0].LayerID", got[1][0].LayerID, FallbackLayerID)
}

func TestVariantLayersBuilderGroupVariantsByLayerLogsUnknownLayer(t *testing.T) {
	builder := NewVariantLayerBuilder(VariantLayersSpec{
		{
			LayerID: 1,
			Variants: []VariantSpec{
				{LayerID: 1, Width: 320},
			},
		},
	})

	var buffer bytes.Buffer
	originalWriter := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&buffer)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(originalWriter)
		log.SetFlags(originalFlags)
	})

	_ = builder.GroupVariantsByLayer([]Variant{
		{LayerID: 999, Width: 640, RelativePath: "unknown-variant"},
		{LayerID: 888, Width: 480, RelativePath: "unknown-variant-2"},
		{LayerID: 999, Width: 320, RelativePath: "unknown-variant-3"},
	})

	logged := strings.TrimSpace(buffer.String())
	if !strings.Contains(logged, "unknown variant layers dropped: count=3 layer_ids=[888 999]") {
		t.Fatalf("unexpected log output: %q", logged)
	}

	if strings.Count(logged, "\n")+1 != 1 {
		t.Fatalf("log lines = %d, want 1: %q", strings.Count(logged, "\n")+1, logged)
	}
}
