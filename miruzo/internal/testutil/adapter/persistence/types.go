package persistence

import "context"

type SuiteBase[R any] struct {
	Context    context.Context
	Operations Operations
	Repository R
}
