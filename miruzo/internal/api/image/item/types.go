package item

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type imageSummaryModel struct {
	IngestID   model.IngestIDType `json:"id"`
	IngestedAt time.Time          `json:"ingested_at"`
	Type       model.ImageType    `json:"type,omitempty"`
}

type imageRichModel struct {
	Level      string             `json:"level"`
	IngestID   model.IngestIDType `json:"id"`
	IngestedAt time.Time          `json:"ingested_at"`
	Type       model.ImageType    `json:"type,omitempty"`
	variant.VariantLayersModel
}

type statModel struct {
	Score                   model.ScoreType `json:"score"`
	FirstLovedAt            *time.Time      `json:"first_loved_at,omitempty"`
	LastLovedAt             *time.Time      `json:"last_loved_at,omitempty"`
	HallOfFameAt            *time.Time      `json:"hall_of_fame_at,omitempty"`
	LastViewedAt            *time.Time      `json:"last_viewed_at,omitempty"`
	ViewCount               int64           `json:"view_count"`
	ViewMilestoneCount      int64           `json:"view_milestone_count,omitzero"`
	ViewMilestoneArchivedAt *time.Time      `json:"view_milestone_archived_at,omitempty"`
}

type contextResponse[ImageModel imageSummaryModel | imageRichModel] struct {
	Image ImageModel `json:"image"`
	Stats statModel  `json:"stats"`
}
