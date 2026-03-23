package persistence

import (
	"fmt"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func NewIngestFixtureWithCapturedAt(
	id model.IngestIDType,
	ingestedAt time.Time,
	capturedAt time.Time,
) persist.Ingest {
	return persist.Ingest{
		ID:           id,
		Process:      model.ProcessStatusProcessing,
		Visibility:   model.VisibilityStatusPrivate,
		RelativePath: fmt.Sprintf("orig/%d.png", id),
		Fingerprint:  fmt.Sprintf("%064d", id),
		IngestedAt:   ingestedAt,
		CapturedAt:   capturedAt,
		UpdatedAt:    ingestedAt,
		Executions:   []model.Execution{},
	}
}

func NewIngestFixture(id model.IngestIDType, ingestedAt time.Time) persist.Ingest {
	return NewIngestFixtureWithCapturedAt(
		id,
		ingestedAt,
		ingestedAt.Add(-30*time.Minute),
	)
}

func DefaultIngestFixture(ingestedAt time.Time) persist.Ingest {
	return NewIngestFixture(0, ingestedAt)
}
