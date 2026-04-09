package persist

import "context"

type MigrationRunner interface {
	Migrate(ctx context.Context, version int) error
	Step(ctx context.Context, steps int) error
	Down(ctx context.Context) error
	Up(ctx context.Context) error

	Version(ctx context.Context) (int, error)
	SetVersion(ctx context.Context, version int) error
}
