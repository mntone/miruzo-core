package item

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapSummaryImage(e persist.Image) imageSummaryModel {
	return imageSummaryModel{
		IngestID:   e.IngestID,
		IngestedAt: e.IngestedAt,
		Type:       e.Type,
	}
}

func mapRichImage(
	e persist.Image,
	cfg []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) imageRichModel {
	return imageRichModel{
		Level:              "rich",
		IngestID:           e.IngestID,
		IngestedAt:         e.IngestedAt,
		Type:               e.Type,
		VariantLayersModel: variant.MapVariantLayers(e, cfg, mediaURLBuilder),
	}
}

func mapStats(e persist.Stats) statModel {
	return statModel{
		Score:                   e.Score,
		FirstLovedAt:            e.FirstLovedAt.ToPointer(),
		LastLovedAt:             e.LastLovedAt.ToPointer(),
		HallOfFameAt:            e.HallOfFameAt.ToPointer(),
		LastViewedAt:            e.LastViewedAt.ToPointer(),
		ViewCount:               e.ViewCount,
		ViewMilestoneCount:      e.ViewMilestoneCount,
		ViewMilestoneArchivedAt: e.ViewMilestoneArchivedAt.ToPointer(),
	}
}

func mapSummaryContextResponse(e persist.ImageWithStats) contextResponse[imageSummaryModel] {
	return contextResponse[imageSummaryModel]{
		Image: mapSummaryImage(e.Image),
		Stats: mapStats(e.Stats),
	}
}

func mapRichContextResponse(
	e persist.ImageWithStats,
	cfg []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) contextResponse[imageRichModel] {
	return contextResponse[imageRichModel]{
		Image: mapRichImage(e.Image, cfg, mediaURLBuilder),
		Stats: mapStats(e.Stats),
	}
}
