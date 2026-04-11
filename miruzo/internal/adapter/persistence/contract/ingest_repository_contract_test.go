package contract_test

import (
	"fmt"
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

// --- Schema ---

func TestIngestSchemaRejectsInvalidRelativePath(t *testing.T) {
	tests := []struct {
		name         string
		relativePath string
	}{
		{
			name:         "relative_path=.bin",
			relativePath: ".bin",
		},
		{
			name:         "relative_path=../orig/test.png",
			relativePath: "../orig/test.png",
		},
		{
			name:         "relative_path=/orig/test.png",
			relativePath: "/orig/test.png",
		},
	}

	baseTime := time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)
	stmt := "INSERT INTO ingests(relative_path, fingerprint, ingested_at, captured_at, updated_at) VALUES(%s, %s, %s, %s, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for i, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(stmt, ops.ParamRange(1, 5)...),
						tt.relativePath,
						fmt.Sprintf("%064d", i+1),
						baseTime,
						baseTime,
						baseTime,
					)
				})
			})
		}
	})
}

func TestIngestSchemaRejectsInvalidOccurredAt(t *testing.T) {
	baseTime := time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		ingestedAt any
		capturedAt any
		updatedAt  any
	}{
		{
			name:       "ingested_at=infinity",
			ingestedAt: "infinity",
			capturedAt: baseTime,
			updatedAt:  baseTime,
		},
		{
			name:       "ingested_at=-infinity",
			ingestedAt: "-infinity",
			capturedAt: baseTime,
			updatedAt:  baseTime,
		},
		{
			name:       "captured_at=-infinity",
			ingestedAt: baseTime,
			capturedAt: "-infinity",
			updatedAt:  baseTime,
		},
		{
			name:       "updated_at=infinity",
			ingestedAt: baseTime,
			capturedAt: baseTime,
			updatedAt:  "infinity",
		},
	}

	stmt := "INSERT INTO ingests(relative_path, fingerprint, ingested_at, captured_at, updated_at) VALUES(%s, %s, %s, %s, %s)"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RequireCapability(t, c.SupportsInfinityTimestamp)

		for i, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				t.Run(tt.name, func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(stmt, ops.ParamRange(1, 5)...),
						"orig/test.png",
						fmt.Sprintf("%064d", i+1),
						tt.ingestedAt,
						tt.capturedAt,
						tt.updatedAt,
					)
				})
			})
		}
	})
}
