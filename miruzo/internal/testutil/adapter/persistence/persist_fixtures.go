package testutil

import (
	"fmt"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model/ingest"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

func NewIngestFixtureWithCapturedAt(
	id persist.IngestID,
	ingestedAt time.Time,
	capturedAt time.Time,
) persist.Ingest {
	return persist.Ingest{
		ID:           id,
		Process:      ingest.ProcessStatusProcessing,
		Visibility:   ingest.VisibilityStatusPrivate,
		RelativePath: fmt.Sprintf("orig/%d.png", id),
		Fingerprint:  fmt.Sprintf("%064d", id),
		IngestedAt:   ingestedAt,
		CapturedAt:   capturedAt,
		UpdatedAt:    ingestedAt,
		Executions:   []ingest.Execution{},
	}
}

func NewIngestFixture(id persist.IngestID, ingestedAt time.Time) persist.Ingest {
	return NewIngestFixtureWithCapturedAt(
		id,
		ingestedAt,
		ingestedAt.Add(-30*time.Minute),
	)
}

func DefaultIngestFixture(ingestedAt time.Time) persist.Ingest {
	return NewIngestFixture(0, ingestedAt)
}

func NewStatFixture(id persist.IngestID) persist.Stat {
	return persist.Stat{
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
	id persist.IngestID,
	score int16,
	scoreEvaluatedAt time.Time,
) persist.Stat {
	entry := NewStatFixture(id)
	entry.Score = score
	entry.ScoreEvaluated = score
	entry.ScoreEvaluatedAt = mo.Some(scoreEvaluatedAt)
	return entry
}

func NewStatFixtureWithLastViewedAt(
	id persist.IngestID,
	viewCount int64,
	lastViewedAt time.Time,
) persist.Stat {
	entry := NewStatFixture(id)
	entry.ViewCount = viewCount
	entry.LastViewedAt = mo.Some(lastViewedAt)
	return entry
}

func NewStatFixtureWithLastLovedAt(
	id persist.IngestID,
	lastLovedAt time.Time,
) persist.Stat {
	entry := NewStatFixture(id)
	entry.FirstLovedAt = mo.Some(lastLovedAt.Add(-24 * time.Hour))
	entry.LastLovedAt = mo.Some(lastLovedAt)
	return entry
}

func NewStatFixtureWithHallOfFameAt(
	id persist.IngestID,
	hallOfFameAt time.Time,
) persist.Stat {
	entry := NewStatFixture(id)
	entry.HallOfFameAt = mo.Some(hallOfFameAt)
	return entry
}
