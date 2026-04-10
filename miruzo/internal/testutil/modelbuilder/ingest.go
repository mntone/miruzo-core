package modelbuilder

import (
	"fmt"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

var defaultBaseTime = time.Date(2026, 1, 10, 5, 0, 0, 0, time.UTC)
var currentNextID model.IngestIDType = 0

func GetDefaultBaseTime() time.Time {
	return defaultBaseTime
}

func GetNextID() model.IngestIDType {
	return currentNextID
}

func SetNextID(nextID model.IngestIDType) {
	currentNextID = nextID
}

type ingestBuilder struct {
	BaseTime time.Time

	ID           model.IngestIDType
	Process      model.ProcessStatus
	Visibility   model.VisibilityStatus
	RelativePath string
	Fingerprint  string
	IngestedAt   time.Time
	CapturedAt   mo.Option[time.Time]
	UpdatedAt    mo.Option[time.Time]
	Executions   []model.Execution
}

func Ingest() *ingestBuilder {
	currentNextID += 1
	return &ingestBuilder{
		BaseTime: defaultBaseTime,

		ID:           currentNextID,
		Process:      model.ProcessStatusProcessing,
		Visibility:   model.VisibilityStatusPrivate,
		RelativePath: fmt.Sprintf("orig/%d.png", currentNextID),
		Fingerprint:  fmt.Sprintf("%064d", currentNextID),
		IngestedAt:   defaultBaseTime,
	}
}

func (b *ingestBuilder) ChangeBaseTime(value time.Time) *ingestBuilder {
	b.BaseTime = value
	return b
}

func (b *ingestBuilder) IngestID(id model.IngestIDType) *ingestBuilder {
	if id < currentNextID {
		panic("invalid ingest id")
	}
	currentNextID = id

	b.ID = id
	b.RelativePath = fmt.Sprintf("orig/%d.png", id)
	b.Fingerprint = fmt.Sprintf("%064d", id)
	return b
}

func (b *ingestBuilder) Processing() *ingestBuilder {
	b.Process = model.ProcessStatusProcessing
	return b
}

func (b *ingestBuilder) Finished() *ingestBuilder {
	b.Process = model.ProcessStatusFinished
	return b
}

func (b *ingestBuilder) Private() *ingestBuilder {
	b.Visibility = model.VisibilityStatusPrivate
	return b
}

func (b *ingestBuilder) Public() *ingestBuilder {
	b.Visibility = model.VisibilityStatusPublic
	return b
}

func (b *ingestBuilder) Ingested(at time.Time) *ingestBuilder {
	b.IngestedAt = at
	return b
}

func (b *ingestBuilder) IngestedOffset(v any) *ingestBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.Ingested(at)
	}
	return b
}

func (b *ingestBuilder) Captured(at time.Time) *ingestBuilder {
	b.CapturedAt = mo.Some(at)
	return b
}

func (b *ingestBuilder) CapturedOffset(v any) *ingestBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.Captured(at)
	}
	return b
}

func (b *ingestBuilder) Updated(at time.Time) *ingestBuilder {
	b.UpdatedAt = mo.Some(at)
	return b
}

func (b *ingestBuilder) UpdatedOffset(v any) *ingestBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.Updated(at)
	}
	return b
}

func (b *ingestBuilder) AppendExecution(execution model.Execution) *ingestBuilder {
	b.Executions = append(b.Executions, execution)
	return b
}

func (b *ingestBuilder) Build() model.Ingest {
	return model.Ingest{
		ID:           b.ID,
		Process:      b.Process,
		Visibility:   b.Visibility,
		RelativePath: b.RelativePath,
		Fingerprint:  b.Fingerprint,
		IngestedAt:   b.IngestedAt,
		CapturedAt:   b.CapturedAt.OrElse(b.IngestedAt),
		UpdatedAt:    b.UpdatedAt.OrElse(b.IngestedAt),
		Executions:   b.Executions,
	}
}
