package persistence

import (
	"fmt"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

var ingestSuiteBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

type IngestSuite SuiteBase[bool]

func (ste IngestSuite) RunTestIngestSchemaRejectsInvalidRelativePath(t *testing.T) {
	t.Helper()

	tests := []struct {
		name         string
		relativePath string
		wantErr      error
	}{
		{
			name:         "relative_path=.bin",
			relativePath: ".bin",
			wantErr:      persist.ErrCheckViolation,
		},
		{
			name:         "relative_path=../orig/test.png",
			relativePath: "../orig/test.png",
			wantErr:      persist.ErrCheckViolation,
		},
		{
			name:         "relative_path=/orig/test.png",
			relativePath: "/orig/test.png",
			wantErr:      persist.ErrCheckViolation,
		},
	}

	timestamp := ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := fmt.Sprintf(
				"INSERT INTO ingests(relative_path, fingerprint, ingested_at, captured_at, updated_at) VALUES('%s', '%s', '%s', '%s', '%s')",
				tt.relativePath,
				fmt.Sprintf("%064d", 1),
				timestamp,
				timestamp,
				timestamp,
			)
			err := ste.Operations.ExecuteStatement(stmt)
			assert.ErrorIs(t, "insert error", err, tt.wantErr)
		})
	}
}

// PostgreSQL only
func (ste IngestSuite) RunTestIngestSchemaRejectsInvalidOccurredAt(t *testing.T) {
	t.Helper()

	tests := []struct {
		name       string
		ingestedAt string
		capturedAt string
		updatedAt  string
		wantErr    error
	}{
		{
			name:       "ingested_at=infinity",
			ingestedAt: "infinity",
			capturedAt: ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			updatedAt:  ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			wantErr:    persist.ErrCheckViolation,
		},
		{
			name:       "ingested_at=-infinity",
			ingestedAt: "-infinity",
			capturedAt: ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			updatedAt:  ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			wantErr:    persist.ErrCheckViolation,
		},
		{
			name:       "captured_at=-infinity",
			ingestedAt: ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			capturedAt: "-infinity",
			updatedAt:  ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			wantErr:    persist.ErrCheckViolation,
		},
		{
			name:       "updated_at=infinity",
			ingestedAt: ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			capturedAt: ingestSuiteBaseTimeUTC.Format(time.RFC3339Nano),
			updatedAt:  "infinity",
			wantErr:    persist.ErrCheckViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := fmt.Sprintf(
				"INSERT INTO ingests(relative_path, fingerprint, ingested_at, captured_at, updated_at) VALUES('%s', '%s', '%s', '%s', '%s')",
				"orig/test.png",
				fmt.Sprintf("%064d", 1),
				tt.ingestedAt, tt.capturedAt, tt.updatedAt,
			)
			err := ste.Operations.ExecuteStatement(stmt)
			assert.ErrorIs(t, "insert error", err, tt.wantErr)
		})
	}
}
