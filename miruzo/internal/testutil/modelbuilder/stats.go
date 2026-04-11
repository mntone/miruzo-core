package modelbuilder

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

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
	return defaultBaseTime
}

func Stats(id model.IngestIDType) *statsBuilder {
	if id <= 0 {
		panic("invalid ingest id")
	}

	return &statsBuilder{
		BaseTime: defaultBaseTime,

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

func (b *statsBuilder) ScoreOption(value mo.Option[model.ScoreType]) *statsBuilder {
	if score, present := value.Get(); present {
		return b.Score(score)
	}
	return b
}

func (b *statsBuilder) EvaluateScore(at time.Time) *statsBuilder {
	b.scoreEvaluated = b.score
	b.scoreEvaluatedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) EvaluateScoreOffset(v any) *statsBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.EvaluateScore(at)
	}
	return b
}

func (b *statsBuilder) Loved(at time.Time) *statsBuilder {
	lovedAt := mo.Some(at)
	b.firstLovedAt = lovedAt
	b.lastLovedAt = lovedAt
	return b
}

func (b *statsBuilder) LovedOffset(v any) *statsBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.Loved(at)
	}
	return b
}

func (b *statsBuilder) FirstLoved(at time.Time) *statsBuilder {
	b.firstLovedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) FirstLovedOffset(v any) *statsBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.FirstLoved(at)
	}
	return b
}

func (b *statsBuilder) LastLoved(at time.Time) *statsBuilder {
	b.lastLovedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) LastLovedOffset(v any) *statsBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.LastLoved(at)
	}
	return b
}

func (b *statsBuilder) HallOfFame(at time.Time) *statsBuilder {
	b.hallOfFameAt = mo.Some(at)
	return b
}

func (b *statsBuilder) HallOfFameOffset(v any) *statsBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.HallOfFame(at)
	}
	return b
}

func (b *statsBuilder) Viewed(count int64, at time.Time) *statsBuilder {
	b.viewCount = count
	b.lastViewedAt = mo.Some(at)
	return b
}

func (b *statsBuilder) ViewedOffset(count int64, v any) *statsBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.Viewed(count, at)
	}
	return b
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
