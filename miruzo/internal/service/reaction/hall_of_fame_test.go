package reaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/stub"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	testutilDomain "github.com/mntone/miruzo-core/miruzo/internal/testutil/domain"
	"github.com/samber/mo"
)

// --- grant ---

func TestGrantHallOfFameUpdates(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID: ingestID,
		Score:    180,
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	response, err := service.GrantHallOfFame(context.Background(), ingestID)
	assert.NilError(t, "GrantHallOfFame() error", err)
	assert.EqualFn(t, "GrantHallOfFame().Stats.HallOfFameAt", response.Stats.HallOfFameAt, mo.Some(current))

	statsArgs := provider.StatsStub.ApplyHallOfFameGrantedArgs[0]
	assert.Equal(t, "statsArgs.IngestID", statsArgs.IngestID, ingestID)
	assert.EqualFn(t, "statsArgs.HallOfFameAt", statsArgs.HallOfFameAt, current)
	assert.Equal(t, "statsArgs.HallOfFameScoreThreshold", statsArgs.HallOfFameScoreThreshold, 180)

	actionArgs := provider.ActionStub.CreateArgs[0]
	assert.Equal(t, "actionArgs.IngestID", actionArgs.IngestID, ingestID)
	assert.Equal(t, "actionArgs.Type", actionArgs.Type, model.ActionTypeHallOfFameGranted)
	assert.EqualFn(t, "actionArgs.OccurredAt", actionArgs.OccurredAt, current)
}

func TestGrantHallOfFameReturnsConflict(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID: ingestID,
		Score:    170,
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.GrantHallOfFame(context.Background(), ingestID)
	assert.ErrorIs(t, "GrantHallOfFame() error", err, serviceerror.ErrConflict)

	args := provider.StatsStub.ApplyHallOfFameGrantedArgs[0]
	assert.Equal(t, "args.IngestID", args.IngestID, ingestID)
	assert.EqualFn(t, "args.HallOfFameAt", args.HallOfFameAt, current)
	assert.Equal(t, "args.HallOfFameScoreThreshold", args.HallOfFameScoreThreshold, 180)

	assert.Empty(t, "manager.ActionStub.CreateArgs", provider.ActionStub.CreateArgs)
	assert.Empty(t, "action count", provider.ActionStub.Store)
}

func TestGrantHallOfFameReturnsServiceUnavailableWhenStatsUpdateFails(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID: ingestID,
		Score:    180,
	})
	provider.StatsStub.ApplyHallOfFameGrantedError = persist.ErrUnavailable
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.GrantHallOfFame(context.Background(), ingestID)
	assert.ErrorIs(t, "GrantHallOfFame() error", err, serviceerror.ErrServiceUnavailable)
	assert.Empty(t, "manager.ActionStub.CreateArgs", provider.ActionStub.CreateArgs)
}

func TestGrantHallOfFameRollsBackWhenActionCreateFails(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID: ingestID,
		Score:    180,
	})
	provider.ActionStub.CreateError = persist.ErrUnavailable
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.GrantHallOfFame(context.Background(), ingestID)
	assert.ErrorIs(t, "GrantHallOfFame() error", err, serviceerror.ErrServiceUnavailable)

	stats := provider.StatsStub.Store[ingestID]
	assert.IsAbsent(t, "stats.HallOfFameAt", stats.HallOfFameAt)

	args := provider.StatsStub.ApplyHallOfFameGrantedArgs[0]
	assert.Equal(t, "args.IngestID", args.IngestID, ingestID)
	assert.EqualFn(t, "args.HallOfFameAt", args.HallOfFameAt, current)
	assert.Equal(t, "args.HallOfFameScoreThreshold", args.HallOfFameScoreThreshold, 180)

	assert.LenIs(t, "manager.ActionStub.CreateArgs", provider.ActionStub.CreateArgs, 1)
	assert.Empty(t, "action count", provider.ActionStub.Store)
}

// --- revoke ---

func TestRevokeHallOfFameUpdates(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID:     ingestID,
		Score:        180,
		HallOfFameAt: mo.Some(current.Add(-2 * time.Hour)),
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	response, err := service.RevokeHallOfFame(context.Background(), ingestID)
	assert.NilError(t, "RevokeHallOfFame() error", err)
	assert.IsAbsent(t, "RevokeHallOfFame().Stats.HallOfFameAt", response.Stats.HallOfFameAt)

	statsArgs := provider.StatsStub.ApplyHallOfFameRevokedArgs[0]
	assert.Equal(t, "statsArgs.IngestID", statsArgs, ingestID)

	actionArgs := provider.ActionStub.CreateArgs[0]
	assert.Equal(t, "actionArgs.IngestID", actionArgs.IngestID, ingestID)
	assert.Equal(t, "actionArgs.Type", actionArgs.Type, model.ActionTypeHallOfFameRevoked)
	assert.EqualFn(t, "actionArgs.OccurredAt", actionArgs.OccurredAt, current)
}

func TestRevokeHallOfFameReturnsConflict(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID: ingestID,
		Score:    170,
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.RevokeHallOfFame(context.Background(), ingestID)
	assert.ErrorIs(t, "RevokeHallOfFame() error", err, serviceerror.ErrConflict)

	args := provider.StatsStub.ApplyHallOfFameRevokedArgs[0]
	assert.Equal(t, "statsArgs.IngestID", args, ingestID)

	assert.Empty(t, "manager.ActionStub.CreateArgs", provider.ActionStub.CreateArgs)
	assert.Empty(t, "action count", provider.ActionStub.Store)
}

func TestRevokeHallOfFameReturnsServiceUnavailableWhenStatsUpdateFails(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID:     ingestID,
		Score:        180,
		HallOfFameAt: mo.Some(current.Add(-2 * time.Hour)),
	})
	provider.StatsStub.ApplyHallOfFameRevokedError = persist.ErrUnavailable
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.RevokeHallOfFame(context.Background(), ingestID)
	assert.ErrorIs(t, "RevokeHallOfFame() error", err, serviceerror.ErrServiceUnavailable)
	assert.Empty(t, "manager.ActionStub.CreateArgs", provider.ActionStub.CreateArgs)
}

func TestRevokeHallOfFameRollsBackWhenActionCreateFails(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	provider := stub.NewStubPersistenceProvider(0, model.Stats{
		IngestID:     ingestID,
		Score:        180,
		HallOfFameAt: mo.Some(current.Add(-2 * time.Hour)),
	})
	provider.ActionStub.CreateError = persist.ErrUnavailable
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		provider,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.RevokeHallOfFame(context.Background(), ingestID)
	assert.ErrorIs(t, "RevokeHallOfFame() error", err, serviceerror.ErrServiceUnavailable)

	stats := provider.StatsStub.Store[ingestID]
	assert.EqualFn(t, "stats.HallOfFameAt", stats.HallOfFameAt, mo.Some(current.Add(-2*time.Hour)))

	args := provider.StatsStub.ApplyHallOfFameRevokedArgs[0]
	assert.Equal(t, "args.IngestID", args, ingestID)

	assert.LenIs(t, "manager.ActionStub.CreateArgs", provider.ActionStub.CreateArgs, 1)
	assert.Empty(t, "action count", provider.ActionStub.Store)
}
