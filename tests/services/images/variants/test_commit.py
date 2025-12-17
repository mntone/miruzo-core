from pathlib import Path

from tests.services.images.variants.utils import build_webp_info

from app.config.variant import WEBP_FORMAT, VariantSlotkey, VariantSpec
from app.services.images.variants.commit import _delete_variant_file
from app.services.images.variants.types import VariantFile


def _build_variant_file(file_path: Path) -> VariantFile:
	info = build_webp_info(width=100, height=80)
	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=1, width=200),
		layer_id=1,
		width=200,
		format=WEBP_FORMAT,
	)
	return VariantFile(
		bytes=0,
		info=info,
		path=file_path,
		variant_dir=spec.slotkey.label,
	)


def test_delete_variant_file_succeeds_when_file_exists(tmp_path: Path) -> None:
	file_path = tmp_path / 'foo.webp'
	file_path.write_bytes(b'data')
	variant_file = _build_variant_file(file_path)

	result = _delete_variant_file(variant_file)

	assert result.action == 'delete'
	assert result.result == 'success'
	assert not file_path.exists()


def test_delete_variant_file_returns_missing_when_file_absent(tmp_path: Path) -> None:
	file_path = tmp_path / 'missing.webp'
	variant_file = _build_variant_file(file_path)

	result = _delete_variant_file(variant_file)

	assert result.result == 'failure'
	assert result.reason == 'file_already_missing'
