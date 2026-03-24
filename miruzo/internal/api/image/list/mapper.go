package list

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	imageListService "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

func mapImage(
	entry persist.Image,
	spec media.VariantLayersSpec,
	mediaURLBuilder variant.MediaURLBuilder,
) ImageListModel {
	return ImageListModel{
		IngestID:           entry.IngestID,
		VariantLayersModel: variant.MapVariantLayers(entry, spec, mediaURLBuilder),
	}
}

func mapImageList(
	entries []persist.Image,
	spec media.VariantLayersSpec,
	mediaURLBuilder variant.MediaURLBuilder,
) []ImageListModel {
	models := make([]ImageListModel, len(entries))

	for i, entry := range entries {
		models[i] = mapImage(entry, spec, mediaURLBuilder)
	}

	return models
}

func mapImageListResponse[T persist.ImageListCursor](
	result imageListService.Result[T],
	spec media.VariantLayersSpec,
	mediaURLBuilder variant.MediaURLBuilder,
) ImageListResponse[T] {
	return ImageListResponse[T]{
		Items:  mapImageList(result.Items, spec, mediaURLBuilder),
		Cursor: result.Cursor.ToPointer(),
	}
}
