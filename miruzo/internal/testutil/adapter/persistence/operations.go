package persistence

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type Operations struct {
	ctx    context.Context
	Action persist.ActionRepository
	test   TestRepository
}

func NewOperations(
	ctx context.Context,
	action persist.ActionRepository,
	test TestRepository,
) Operations {
	return Operations{
		ctx:    ctx,
		Action: action,
		test:   test,
	}
}

func (ops Operations) MustAddAction(
	t testing.TB,
	ingestID model.IngestIDType,
	kind model.ActionType,
	occurredAt time.Time,
) model.ActionIDType {
	t.Helper()

	actionID, err := ops.Action.Create(ops.ctx, ingestID, kind, occurredAt)
	if err != nil {
		t.Fatalf("add action: %v", err)
	}
	return actionID
}

func (ops Operations) AddIngest(entry persist.Ingest) error {
	return ops.test.CreateIngest(
		ops.ctx,
		entry.ID,
		entry.RelativePath,
		entry.Fingerprint,
		entry.IngestedAt,
		entry.CapturedAt,
	)
}

func (ops Operations) MustAddIngest(t testing.TB, entry persist.Ingest) persist.Ingest {
	t.Helper()

	err := ops.AddIngest(entry)
	if err != nil {
		t.Fatalf("add ingest: %v", err)
	}
	return entry
}

func createVariant(
	id model.IngestIDType,
	layerID media.LayerIDType,
	format string,
	width uint16,
) persist.Variant {
	var codecs string
	if format == "webp" {
		codecs = "vp8"
	}

	return persist.Variant{
		RelativePath: fmt.Sprintf("l%dw%d/%d.%s", layerID, width, id, format),
		LayerID:      layerID,
		Format:       format,
		Codecs:       codecs,
		Bytes:        234,
		Width:        width,
		Height:       uint16(math.Round(0.75 * float64(width))),
		Quality:      nil,
	}
}

func (ops Operations) AddIngestAndImage(entry persist.Ingest) error {
	err := ops.AddIngest(entry)
	if err != nil {
		return err
	}

	return ops.test.CreateImage(
		ops.ctx,
		entry.ID,
		entry.IngestedAt,
		createVariant(entry.ID, 1, "webp", 768),
		mo.None[persist.Variant](),
		[]persist.Variant{
			createVariant(entry.ID, 1, "webp", 320),
			createVariant(entry.ID, 1, "webp", 480),
			createVariant(entry.ID, 1, "webp", 640),
			createVariant(entry.ID, media.FallbackLayerID, "jpeg", 320),
		},
	)
}

func (ops Operations) MustAddIngestAndImage(t testing.TB, entry persist.Ingest) persist.Ingest {
	t.Helper()

	err := ops.AddIngestAndImage(entry)
	if err != nil {
		t.Fatalf("add ingest and image: %v", err)
	}

	return entry
}

func (ops Operations) MustAddStat(t testing.TB, entry model.Stats) model.Stats {
	t.Helper()

	err := ops.test.CreateStat(
		ops.ctx,
		entry.IngestID,
		entry.Score,
		entry.ScoreEvaluated,
		entry.LastViewedAt,
		entry.FirstLovedAt,
		entry.LastLovedAt,
		entry.HallOfFameAt,
		entry.ViewCount,
	)
	if err != nil {
		t.Fatalf("add stat: %v", err)
	}

	return entry
}

func (ops Operations) ExecuteStatement(stmt string) error {
	return ops.test.ExecuteStatement(ops.ctx, stmt, false)
}

func (ops Operations) MustRemoveUser(t testing.TB) {
	t.Helper()

	rowCount, err := ops.test.ExecuteStatementAndReturnRowCount(ops.ctx, "DELETE FROM users WHERE id=1", true)
	if err != nil {
		t.Fatalf("remove user: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("remove user: row_count=%d", rowCount)
	}
}

func (ops Operations) MustSetDailyLoveUsed(t testing.TB, dailyLoveUsed model.QuotaInt) {
	t.Helper()

	rowCount, err := ops.test.ExecuteStatementAndReturnRowCount(
		ops.ctx,
		fmt.Sprintf("UPDATE users SET daily_love_used=%d WHERE id=1", dailyLoveUsed),
		false,
	)
	if err != nil {
		t.Fatalf("set daily_love_used to user: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("set daily_love_used to user: row_count=%d", rowCount)
	}
}

func (ops Operations) MustTruncateActions(t testing.TB) {
	t.Helper()

	err := ops.test.TruncateActions(ops.ctx)
	if err != nil {
		t.Fatalf("truncate actions: %v", err)
	}
}

func (ops Operations) MustTruncateStats(t testing.TB) {
	t.Helper()

	err := ops.test.TruncateStats(ops.ctx)
	if err != nil {
		t.Fatalf("truncate stats: %v", err)
	}
}
