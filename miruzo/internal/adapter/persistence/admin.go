package persistence

import "context"

type DatabaseAdmin interface {
	Create(ctx context.Context) error
	Drop(ctx context.Context) error
	Exists(ctx context.Context) (bool, error)
}
