package modelbuilder

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

var defaultStatsBaseTime = time.Date(2026, 1, 10, 5, 0, 0, 0, time.UTC)

type statsBuilder struct {
	BaseTime time.Time

	ingestID         model.IngestIDType
	score            model.ScoreType
	scoreEvaluated   model.ScoreType
	scoreEvaluatedAt mo.Option[time.Time]

	firstLovedAt mo.Option[time.Time]
	lastLovedAt  mo.Option[time.Time]
	hallOfFameAt mo.Option[time.Time]
	lastViewedAt mo.Option[time.Time]

	viewCount               int64
	viewMilestoneCount      int64
	viewMilestoneArchivedAt mo.Option[time.Time]
}

func GetDefaultStatsBaseTime() time.Time {
	return defaultStatsBaseTime
}

func Stats(id model.IngestIDType) *statsBuilder {
	if id <= 0 {
		panic("invalid ingest id")
	}

	return &statsBuilder{
		BaseTime: defaultStatsBaseTime,

		ingestID:         id,
		score:            100,
		scoreEvaluated:   100,
		scoreEvaluatedAt: mo.None[time.Time](),
	}
}

func (b *statsBuilder) ChangeBaseTime(value time.Time) *statsBuilder {
	b.BaseTime = value
	return b
}

func (b *statsBuilder) Score(value model.ScoreType) *statsBuilder {
	b.score = value
	return b
}

func (b *statsBuilder) EvaluateScore(at time.Time) *statsBuilder {
	b.scoreEvaluated = b.score
	b.scoreEvaluatedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) Loved(at time.Time) *statsBuilder {
	lovedAt := mo.Some(at)
	b.firstLovedAt = lovedAt
	b.lastLovedAt = lovedAt
	return b
}

func (b *statsBuilder) LovedOffset(v any) *statsBuilder {
	switch value := v.(type) {
	case time.Duration:
		return b.Loved(b.BaseTime.Add(value))
	case mo.Option[time.Duration]:
		if duration, present := value.Get(); present {
			return b.Loved(b.BaseTime.Add(duration))
		}
		return b
	case int:
		return b.Loved(b.BaseTime.Add(time.Duration(value) * time.Second))
	}
	panic("invalid offset")
}

func (b *statsBuilder) FirstLoved(at time.Time) *statsBuilder {
	b.firstLovedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) FirstLovedOffset(v any) *statsBuilder {
	switch value := v.(type) {
	case time.Duration:
		return b.FirstLoved(b.BaseTime.Add(value))
	case mo.Option[time.Duration]:
		if duration, present := value.Get(); present {
			return b.FirstLoved(b.BaseTime.Add(duration))
		}
		return b
	case int:
		return b.FirstLoved(b.BaseTime.Add(time.Duration(value) * time.Second))
	}
	panic("invalid offset")
}

func (b *statsBuilder) LastLoved(at time.Time) *statsBuilder {
	b.lastLovedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) LastLovedOffset(v any) *statsBuilder {
	switch value := v.(type) {
	case time.Duration:
		return b.LastLoved(b.BaseTime.Add(value))
	case mo.Option[time.Duration]:
		if duration, present := value.Get(); present {
			return b.LastLoved(b.BaseTime.Add(duration))
		}
		return b
	case int:
		return b.LastLoved(b.BaseTime.Add(time.Duration(value) * time.Second))
	}
	panic("invalid offset")
}

func (b *statsBuilder) HallOfFame(at time.Time) *statsBuilder {
	b.hallOfFameAt = mo.Some(at)
	return b
}

func (b *statsBuilder) HallOfFameOffset(v any) *statsBuilder {
	switch value := v.(type) {
	case time.Duration:
		return b.HallOfFame(b.BaseTime.Add(value))
	case mo.Option[time.Duration]:
		if duration, present := value.Get(); present {
			return b.HallOfFame(b.BaseTime.Add(duration))
		}
		return b
	case int:
		return b.HallOfFame(b.BaseTime.Add(time.Duration(value) * time.Second))
	}
	panic("invalid offset")
}

func (b *statsBuilder) Viewed(count int64, at time.Time) *statsBuilder {
	b.viewCount = count
	b.lastViewedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) ViewedOffset(count int64, v any) *statsBuilder {
	switch value := v.(type) {
	case time.Duration:
		return b.Viewed(count, b.BaseTime.Add(value))
	case mo.Option[time.Duration]:
		if duration, present := value.Get(); present {
			return b.Viewed(count, b.BaseTime.Add(duration))
		}
		return b
	case int:
		return b.Viewed(count, b.BaseTime.Add(time.Duration(value)*time.Second))
	}
	panic("invalid offset")
}

func (b *statsBuilder) Build() model.Stats {
	return model.Stats{
		IngestID:         b.ingestID,
		Score:            b.score,
		ScoreEvaluated:   b.scoreEvaluated,
		ScoreEvaluatedAt: b.scoreEvaluatedAt,

		FirstLovedAt: b.firstLovedAt,
		LastLovedAt:  b.lastLovedAt,
		HallOfFameAt: b.hallOfFameAt,
		LastViewedAt: b.lastViewedAt,

		ViewCount:               b.viewCount,
		ViewMilestoneCount:      b.viewMilestoneCount,
		ViewMilestoneArchivedAt: b.viewMilestoneArchivedAt,
	}
}
