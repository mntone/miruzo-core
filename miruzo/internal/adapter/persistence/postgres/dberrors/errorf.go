package dberrors

import (
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
)

func WrapKV(base error, operation string, keyValues ...any) error {
	return persistshared.WrapKV(base, operation, keyValues...)
}
