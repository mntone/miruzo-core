package persist

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model/ingest"
)

type IngestID = int64

type Ingest struct {
	ID           IngestID
	Process      ingest.ProcessStatus
	Visibility   ingest.VisibilityStatus
	RelativePath string
	Fingerprint  string
	IngestedAt   time.Time
	CapturedAt   time.Time
	UpdatedAt    time.Time
	Executions   []ingest.Execution
}
