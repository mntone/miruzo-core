package model

import "time"

type ActionIDType = int64

type ActionType uint32

const (
	ActionTypeUnspecified       ActionType = 0
	ActionTypeDecay             ActionType = 1
	ActionTypeView              ActionType = 11
	ActionTypeMemo              ActionType = 12
	ActionTypeLove              ActionType = 13
	ActionTypeLoveCanceled      ActionType = 14
	ActionTypeHallOfFameGranted ActionType = 15
	ActionTypeHallOfFameRevoked ActionType = 16
)

type Action struct {
	ID         ActionIDType
	IngestID   IngestIDType
	Type       ActionType
	OccurredAt time.Time
}
