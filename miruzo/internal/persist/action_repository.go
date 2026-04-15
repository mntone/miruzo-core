package persist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type LoveActionType uint32

const (
	LoveActionTypeLove         LoveActionType = LoveActionType(model.ActionTypeLove)
	LoveActionTypeLoveCanceled LoveActionType = LoveActionType(model.ActionTypeLoveCanceled)
)

type ActionRepository interface {
	Create(
		ctx context.Context,
		ingestID model.IngestIDType,
		kind model.ActionType,
		occurredAt time.Time,
		periodStartAt time.Time,
	) (model.ActionIDType, error)

	CreateDailyDecayIfAbsent(
		ctx context.Context,
		ingestID model.IngestIDType,
		occurredAt time.Time,
		periodStartAt time.Time,
	) error

	CreateLoveIfAbsent(
		ctx context.Context,
		ingestID model.IngestIDType,
		loveType LoveActionType,
		occurredAt time.Time,
		periodStartAt time.Time,
	) error

	ExistsSince(
		ctx context.Context,
		ingestID model.IngestIDType,
		kind model.ActionType,
		sinceOccurredAt time.Time,
	) (bool, error)
}
