package persistence

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/model/media"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type Operations struct {
	ctx  context.Context
	repo TestRepository
}

func NewOperations(
	ctx context.Context,
	repo TestRepository,
) Operations {
	return Operations{
		ctx:  ctx,
		repo: repo,
	}
}

func (ops Operations) AddIngest(entry persist.Ingest) error {
	return ops.repo.CreateIngest(
		ops.ctx,
		entry.ID,
		entry.RelativePath,
		entry.Fingerprint,
		entry.IngestedAt,
		entry.CapturedAt,
	)
}

func (ops Operations) MustAddIngest(t testing.TB, entry persist.Ingest) persist.Ingest {
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
) media.Variant {
	var codecs string
	if format == "webp" {
		codecs = "vp8"
	}

	return media.Variant{
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

	return ops.repo.CreateImage(
		ops.ctx,
		entry.ID,
		entry.IngestedAt,
		createVariant(entry.ID, 1, "webp", 768),
		mo.None[media.Variant](),
		[]media.Variant{
			createVariant(entry.ID, 1, "webp", 320),
			createVariant(entry.ID, 1, "webp", 480),
			createVariant(entry.ID, 1, "webp", 640),
			createVariant(entry.ID, media.FallbackLayerID, "jpeg", 320),
		},
	)
}

func (ops Operations) MustAddIngestAndImage(t testing.TB, entry persist.Ingest) persist.Ingest {
	err := ops.AddIngestAndImage(entry)
	if err != nil {
		t.Fatalf("add ingest and image: %v", err)
	}

	return entry
}

func (ops Operations) AddStat(entry persist.Stats) error {
	return ops.repo.CreateStat(
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
}

func (ops Operations) MustAddStat(t testing.TB, entry persist.Stats) persist.Stats {
	err := ops.AddStat(entry)
	if err != nil {
		t.Fatalf("add stat: %v", err)
	}

	return entry
}

func (ops Operations) RemoveUser() error {
	return ops.repo.DeleteUser(ops.ctx)
}

func (ops Operations) MustRemoveUser(t testing.TB) {
	err := ops.RemoveUser()
	if err != nil {
		t.Fatalf("remove user: %v", err)
	}
}
