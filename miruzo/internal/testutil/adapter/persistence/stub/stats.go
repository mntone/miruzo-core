package stub

import (
	"context"
	"maps"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type statsRepositoryApplyHallOfFameGrantedArgs struct {
	IngestID                 model.IngestIDType
	HallOfFameAt             time.Time
	HallOfFameScoreThreshold model.ScoreType
}

type statsRepositoryApplyLoveArgs struct {
	IngestID           model.IngestIDType
	ScoreDelta         model.ScoreType
	LovedAt            time.Time
	LoveScoreThreshold model.ScoreType
	PeriodStartAt      time.Time
}

type statsRepositoryApplyLoveCanceledArgs struct {
	IngestID       model.IngestIDType
	ScoreDelta     model.ScoreType
	PeriodStartAt  time.Time
	DayStartOffset time.Duration
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

	ApplyHallOfFameGrantedError error
	ApplyHallOfFameGrantedArgs  []statsRepositoryApplyHallOfFameGrantedArgs
	ApplyHallOfFameRevokedError error
	ApplyHallOfFameRevokedArgs  []model.IngestIDType
	ApplyLoveError              error
	ApplyLoveArgs               []statsRepositoryApplyLoveArgs
	ApplyLoveCanceledError      error
	ApplyLoveCanceledArgs       []statsRepositoryApplyLoveCanceledArgs
	ApplyViewError              error
	ApplyViewArgs               []statsRepositoryApplyViewArgs
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

func (repo *statsRepository) ApplyHallOfFameGranted(
	ctx context.Context,
	ingestID model.IngestIDType,
	hallOfFameAt time.Time,
	hallOfFameScoreThreshold model.ScoreType,
) error {
	repo.ApplyHallOfFameGrantedArgs = append(repo.ApplyHallOfFameGrantedArgs, statsRepositoryApplyHallOfFameGrantedArgs{
		IngestID:                 ingestID,
		HallOfFameAt:             hallOfFameAt,
		HallOfFameScoreThreshold: hallOfFameScoreThreshold,
	})

	if repo.ApplyHallOfFameGrantedError != nil {
		return repo.ApplyHallOfFameGrantedError
	}

	stats, ok := repo.Store[ingestID]
	if !ok {
		return persist.ErrConflict
	}

	if stats.HallOfFameAt.IsPresent() || stats.Score < hallOfFameScoreThreshold {
		return persist.ErrConflict
	}

	stats.HallOfFameAt = mo.Some(hallOfFameAt)

	repo.Store[ingestID] = stats
	return nil
}

func (repo *statsRepository) ApplyHallOfFameRevoked(
	ctx context.Context,
	ingestID model.IngestIDType,
) error {
	repo.ApplyHallOfFameRevokedArgs = append(repo.ApplyHallOfFameRevokedArgs, ingestID)

	if repo.ApplyHallOfFameRevokedError != nil {
		return repo.ApplyHallOfFameRevokedError
	}

	stats, ok := repo.Store[ingestID]
	if !ok {
		return persist.ErrConflict
	}

	if stats.HallOfFameAt.IsAbsent() {
		return persist.ErrConflict
	}

	stats.HallOfFameAt = mo.None[time.Time]()

	repo.Store[ingestID] = stats
	return nil
}

func (repo *statsRepository) ApplyLove(
	_ context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	lovedAt time.Time,
	loveScoreThreshold model.ScoreType,
	periodStartAt time.Time,
) (persist.LoveStats, error) {
	repo.ApplyLoveArgs = append(repo.ApplyLoveArgs, statsRepositoryApplyLoveArgs{
		IngestID:           ingestID,
		ScoreDelta:         scoreDelta,
		LovedAt:            lovedAt,
		LoveScoreThreshold: loveScoreThreshold,
		PeriodStartAt:      periodStartAt,
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
	if stats.Score >= loveScoreThreshold {
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

func (repo *statsRepository) ApplyLoveCanceled(
	_ context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	periodStartAt time.Time,
	dayStartOffset time.Duration,
) (persist.LoveStats, error) {
	repo.ApplyLoveCanceledArgs = append(repo.ApplyLoveCanceledArgs, statsRepositoryApplyLoveCanceledArgs{
		IngestID:       ingestID,
		ScoreDelta:     scoreDelta,
		PeriodStartAt:  periodStartAt,
		DayStartOffset: dayStartOffset,
	})

	if repo.ApplyLoveCanceledError != nil {
		return persist.LoveStats{}, repo.ApplyLoveCanceledError
	}

	stats, ok := repo.Store[ingestID]
	if !ok {
		return persist.LoveStats{}, persist.ErrConflict
	}

	lastLovedAt, present := stats.LastLovedAt.Get()
	if !present || lastLovedAt.Compare(periodStartAt) < 0 {
		return persist.LoveStats{}, persist.ErrConflict
	}
	stats.Score += scoreDelta

	latest := mo.None[time.Time]()
	if firstLovedAt, firstPresent := stats.FirstLovedAt.Get(); firstPresent && firstLovedAt.Compare(periodStartAt) < 0 {
		latest = mo.Some(firstLovedAt)
	}
	stats.FirstLovedAt = latest
	stats.LastLovedAt = latest

	repo.Store[ingestID] = stats
	return persist.LoveStats{
		Score:        stats.Score,
		FirstLovedAt: stats.FirstLovedAt,
		LastLovedAt:  stats.LastLovedAt,
	}, nil
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
