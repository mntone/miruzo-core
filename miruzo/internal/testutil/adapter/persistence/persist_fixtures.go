package persistence

import (
	"fmt"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
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

func NewStatFixture(id model.IngestIDType) persist.Stats {
	return persist.Stats{
		IngestID:         id,
		Score:            100,
		ScoreEvaluated:   100,
		ScoreEvaluatedAt: mo.None[time.Time](),

		FirstLovedAt: mo.None[time.Time](),
		LastLovedAt:  mo.None[time.Time](),
		HallOfFameAt: mo.None[time.Time](),
		LastViewedAt: mo.None[time.Time](),

		ViewCount:               0,
		ViewMilestoneCount:      0,
		ViewMilestoneArchivedAt: mo.None[time.Time](),
	}
}

func NewStatFixtureWithScore(
	id model.IngestIDType,
	score int16,
	scoreEvaluatedAt time.Time,
) persist.Stats {
	entry := NewStatFixture(id)
	entry.Score = score
	entry.ScoreEvaluated = score
	entry.ScoreEvaluatedAt = mo.Some(scoreEvaluatedAt)
	return entry
}

func NewStatFixtureWithLastViewedAt(
	id model.IngestIDType,
	viewCount int64,
	lastViewedAt time.Time,
) persist.Stats {
	entry := NewStatFixture(id)
	entry.ViewCount = viewCount
	entry.LastViewedAt = mo.Some(lastViewedAt)
	return entry
}

func NewStatFixtureWithLastLovedAt(
	id model.IngestIDType,
	lastLovedAt time.Time,
) persist.Stats {
	entry := NewStatFixture(id)
	entry.FirstLovedAt = mo.Some(lastLovedAt.Add(-24 * time.Hour))
	entry.LastLovedAt = mo.Some(lastLovedAt)
	return entry
}

func NewStatFixtureWithHallOfFameAt(
	id model.IngestIDType,
	hallOfFameAt time.Time,
) persist.Stats {
	entry := NewStatFixture(id)
	entry.HallOfFameAt = mo.Some(hallOfFameAt)
	return entry
}
