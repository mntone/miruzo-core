package persist

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type Ingest struct {
	ID           model.IngestIDType
	Process      model.ProcessStatus
	Visibility   model.VisibilityStatus
	RelativePath string
	Fingerprint  string
	IngestedAt   time.Time
	CapturedAt   time.Time
	UpdatedAt    time.Time
	Executions   []model.Execution
}
