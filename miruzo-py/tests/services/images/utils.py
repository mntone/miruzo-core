from pathlib import Path
from tempfile import gettempdir

from app.config.variant import VariantFormat, VariantSlot, VariantSpec
from app.models.types import VariantEntry

TEST_VARIANT_ROOT = Path(gettempdir()) / 'miruzo-test-variants'


def build_variant(fmt: str, width: int, *, layer_id: int = 1, label: str = 'primary') -> VariantEntry:
	filepath = TEST_VARIANT_ROOT / f'l{layer_id}w{width}-{label}.{fmt}'
	filepath.parent.mkdir(parents=True, exist_ok=True)
	payload = f'{label}-{width}-{fmt}'.encode('utf-8')
	filepath.write_bytes(payload)
	return {
		'rel': filepath.as_posix(),
		'layer_id': layer_id,
		'format': fmt,
		'codecs': None,
		'bytes': filepath.stat().st_size,
		'width': width,
		'height': round(width * 0.75),
		'quality': None,
	}


def build_variant_spec(
	layer_id: int,
	width: int,
	*,
	container: str = 'jpeg',
	codecs: str | None = None,
	quality: int | None = None,
	required: bool = False,
) -> VariantSpec:
	return VariantSpec(
		slot=VariantSlot(layer_id, width),
		layer_id=layer_id,
		width=width,
		format=VariantFormat(container=container, codecs=codecs, file_extension=f'.{container}'),  # pyright: ignore[reportArgumentType]
		quality=quality,
		required=required,
	)
