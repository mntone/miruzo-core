package persist

import "context"

type SettingsRepository interface {
	GetValue(
		ctx context.Context,
		key string,
	) (string, error)

	UpdateValue(
		ctx context.Context,
		key string,
		value string,
	) error
}
