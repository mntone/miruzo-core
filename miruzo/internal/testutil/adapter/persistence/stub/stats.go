package stub

import (
	"context"
	"maps"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type statsRepositoryApplyLoveArgs struct {
	IngestID      model.IngestIDType
	ScoreDelta    model.ScoreType
	LovedAt       time.Time
	PeriodStartAt time.Time
}

type statsRepositoryApplyViewArgs struct {
	IngestID   model.IngestIDType
	ScoreDelta model.ScoreType
	ViewedAt   time.Time
}

type statsStorage struct {
	Store map[model.IngestIDType]persist.Stats
}

type statsRepository struct {
	statsStorage

	ApplyLoveError error
	ApplyLoveArgs  []statsRepositoryApplyLoveArgs
	ApplyViewError error
	ApplyViewArgs  []statsRepositoryApplyViewArgs
}

func NewStubStatsRepository(stats ...persist.Stats) *statsRepository {
	var store map[model.IngestIDType]persist.Stats
	if len(stats) != 0 {
		store = make(map[model.IngestIDType]persist.Stats)
		for _, s := range stats {
			store[s.IngestID] = s
		}
	}
	return &statsRepository{
		statsStorage: statsStorage{
			Store: store,
		},
	}
}

func (repo statsRepository) snapshot() statsStorage {
	var store map[model.IngestIDType]persist.Stats
	if repo.Store != nil {
		store = make(map[model.IngestIDType]persist.Stats, len(repo.Store))
		maps.Copy(store, repo.Store)
	}
	return statsStorage{
		Store: store,
	}
}

func (repo *statsRepository) ApplyLove(
	_ context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	lovedAt time.Time,
	periodStartAt time.Time,
) (persist.LoveStats, error) {
	repo.ApplyLoveArgs = append(repo.ApplyLoveArgs, statsRepositoryApplyLoveArgs{
		IngestID:      ingestID,
		ScoreDelta:    scoreDelta,
		LovedAt:       lovedAt,
		PeriodStartAt: periodStartAt,
	})

	if repo.ApplyLoveError != nil {
		return persist.LoveStats{}, repo.ApplyLoveError
	}

	stats, ok := repo.Store[ingestID]
	if !ok {
		return persist.LoveStats{}, persist.ErrNotFound
	}

	lastLovedAt, present := stats.LastLovedAt.Get()
	if present && lastLovedAt.Compare(periodStartAt) >= 0 {
		return persist.LoveStats{}, persist.ErrConflict
	}

	stats.Score += scoreDelta
	stats.LastLovedAt = mo.Some(lovedAt)
	if stats.FirstLovedAt.IsAbsent() {
		stats.FirstLovedAt = mo.Some(lovedAt)
	}

	repo.Store[ingestID] = stats
	return persist.LoveStats{
		Score:        stats.Score,
		FirstLovedAt: stats.FirstLovedAt,
		LastLovedAt:  stats.LastLovedAt,
	}, nil
}

func (repo statsRepository) ApplyLoveCanceled(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	periodStartAt time.Time,
	dayStartOffset time.Duration,
) (persist.LoveStats, error) {
	return persist.LoveStats{}, persist.ErrUnavailable
}

func (repo *statsRepository) ApplyView(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	viewedAt time.Time,
) error {
	repo.ApplyViewArgs = append(repo.ApplyViewArgs, statsRepositoryApplyViewArgs{
		IngestID:   ingestID,
		ScoreDelta: scoreDelta,
		ViewedAt:   viewedAt,
	})

	if repo.ApplyViewError != nil {
		return repo.ApplyViewError
	}

	stats, ok := repo.Store[ingestID]
	if !ok {
		return persist.ErrNotFound
	}

	stats.Score += scoreDelta
	stats.LastViewedAt = mo.Some(viewedAt)
	stats.ViewCount += 1

	repo.Store[ingestID] = stats
	return nil
}

func (repo *statsRepository) ApplyViewWithMilestone(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	viewedAt time.Time,
) error {
	repo.ApplyViewArgs = append(repo.ApplyViewArgs, statsRepositoryApplyViewArgs{
		IngestID:   ingestID,
		ScoreDelta: scoreDelta,
		ViewedAt:   viewedAt,
	})

	if repo.ApplyViewError != nil {
		return repo.ApplyViewError
	}

	stats, ok := repo.Store[ingestID]
	if !ok {
		return persist.ErrNotFound
	}

	stats.Score += scoreDelta
	stats.LastViewedAt = mo.Some(viewedAt)
	stats.ViewCount += 1
	stats.ViewMilestoneCount += 1
	stats.ViewMilestoneArchivedAt = mo.Some(viewedAt)

	repo.Store[ingestID] = stats
	return nil
}
