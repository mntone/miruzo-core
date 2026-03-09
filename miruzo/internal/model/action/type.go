package action

type ActionID = int64

type ActionType uint32

const (
	ActionTypeUnspecified       ActionType = 0
	ActionTypeDecay             ActionType = 1
	ActionTypeView              ActionType = 11
	ActionTypeMemo              ActionType = 12
	ActionTypeLove              ActionType = 13
	ActionTypeLoveCanceled      ActionType = 14
	ActionTypeHallOfFameAdded   ActionType = 15
	ActionTypeHallOfFameRemoved ActionType = 16
)
