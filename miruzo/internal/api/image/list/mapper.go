package list

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	imageListService "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

func mapImage(
	entry model.Image,
	mediaURLBuilder variant.MediaURLBuilder,
) ImageListModel {
	return ImageListModel{
		IngestID:           entry.IngestID,
		VariantLayersModel: variant.MapVariantLayers(entry, mediaURLBuilder),
	}
}

func mapImageList(
	entries []model.Image,
	mediaURLBuilder variant.MediaURLBuilder,
) []ImageListModel {
	models := make([]ImageListModel, len(entries))

	for i, entry := range entries {
		models[i] = mapImage(entry, mediaURLBuilder)
	}

	return models
}

func mapImageListResponse[T persist.ImageListCursor](
	result imageListService.Result[T],
	mediaURLBuilder variant.MediaURLBuilder,
) ImageListResponse[T] {
	return ImageListResponse[T]{
		Items:  mapImageList(result.Items, mediaURLBuilder),
		Cursor: result.Cursor.ToPointer(),
	}
}
