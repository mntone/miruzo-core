package list

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
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

func mapCursor[ScalarType model.ImageListCursorScalar](
	result imageListService.Result[ScalarType],
	encodeCursor func(cursor model.ImageListCursorKey[ScalarType]) (string, error),
) (string, error) {
	cursor, present := result.Cursor.Get()
	if !present {
		return "", nil
	}

	encoded, err := encodeCursor(cursor)
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func mapDatetimeImageListResponse(
	result imageListService.Result[time.Time],
	mediaURLBuilder variant.MediaURLBuilder,
	mode imageListCursorMode,
) (ImageListResponse, error) {
	cursor, err := mapCursor(
		result,
		func(cursor model.ImageListCursorKey[time.Time]) (string, error) {
			return encodeTimeImageListCursor(mode, cursor)
		},
	)
	if err != nil {
		return ImageListResponse{}, err
	}

	return ImageListResponse{
		Items:  mapImageList(result.Items, mediaURLBuilder),
		Cursor: cursor,
	}, nil
}

func mapEngagedImageListResponse(
	result imageListService.Result[model.ScoreType],
	mediaURLBuilder variant.MediaURLBuilder,
) (ImageListResponse, error) {
	cursor, err := mapCursor(
		result,
		func(cursor model.ImageListCursorKey[model.ScoreType]) (string, error) {
			return encodeUint8ImageListCursor(imageListCursorModeEngaged, cursor)
		},
	)
	if err != nil {
		return ImageListResponse{}, err
	}

	return ImageListResponse{
		Items:  mapImageList(result.Items, mediaURLBuilder),
		Cursor: cursor,
	}, nil
}
