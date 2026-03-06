package list

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	imageListService "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

func mapImage(
	entry persist.Image,
	cfg []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) ImageListModel {
	return ImageListModel{
		IngestID:           entry.IngestID,
		VariantLayersModel: variant.MapVariantLayers(entry, cfg, mediaURLBuilder),
	}
}

func mapImageList(
	entries []persist.Image,
	cfg []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) []ImageListModel {
	models := make([]ImageListModel, len(entries))

	for i, entry := range entries {
		models[i] = mapImage(entry, cfg, mediaURLBuilder)
	}

	return models
}

func mapImageListResponse[T persist.ImageListCursor](
	result imageListService.Result[T],
	cfg []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) ImageListResponse[T] {
	return ImageListResponse[T]{
		Items:  mapImageList(result.Items, cfg, mediaURLBuilder),
		Cursor: result.Cursor.ToPointer(),
	}
}
