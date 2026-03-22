package context

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

// Lightweight metadata returned in context responses.
type imageSummaryModel struct {
	// IngestID is numeric primary key assigned in the database.
	IngestID model.IngestIDType `json:"id"`
	// IngestedAt is the time when the image was ingested.
	IngestedAt time.Time `json:"ingested_at"`
	// Type is categorization of the image content.
	Type model.ImageType `json:"type,omitempty"`
}

// Rich metadata returned in context responses.
type imageRichModel struct {
	// Level is response detail level for this context payload.
	Level string `json:"level"`
	imageSummaryModel
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

// Envelope returned by the context API for a single image.
type contextResponse[ImageModel imageSummaryModel | imageRichModel] struct {
	// Image is metadata for the requested image.
	Image ImageModel `json:"image"`
	// Stats is the latest statistics for the image.
	Stats statModel `json:"stats"`
}
