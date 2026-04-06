from collections.abc import Iterable
from datetime import datetime, timezone

from tests.services.images.utils import build_variant

from app.config.variant import _FALLBACK_LAYER_ID
from app.models.enums import ImageKind
from app.models.image import Image


def make_image_fixture(
	ingest_id: int,
	*,
	ingested_at: datetime | None = None,
	kind: ImageKind = ImageKind.UNSPECIFIED,
	widths: Iterable[int] = [320, 480, 640, 960],
) -> Image:
	return Image(
		ingest_id=ingest_id,
		ingested_at=ingested_at or datetime.now(timezone.utc),
		kind=kind,
		original=build_variant('png', 1024, layer_id=0, label='original'),
		fallback=None,
		variants=[
			*[build_variant('webp', w) for w in widths],
			build_variant('jpeg', 320, layer_id=_FALLBACK_LAYER_ID, label='fallback'),
		],
	)
