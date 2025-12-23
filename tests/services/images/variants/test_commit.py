from pathlib import Path

from tests.services.images.variants.utils import build_webp_info

from app.config.variant import WEBP_FORMAT, VariantSlotkey, VariantSpec
from app.services.images.variants.commit import _delete_variant_file
from app.services.images.variants.path import VariantRelativePath
from app.services.images.variants.types import FileInfo, VariantFile


def _build_variant_file(file_path: Path, file_name: Path) -> VariantFile:
	info = build_webp_info(width=100, height=80)
	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=1, width=200),
		layer_id=1,
		width=200,
		format=WEBP_FORMAT,
	)
	file_info = FileInfo(
		absolute_path=file_path,
		relative_path=VariantRelativePath(file_name),
		bytes=0,
	)
	return VariantFile(
		file_info=file_info,
		image_info=info,
		variant_dir=spec.slotkey.label,
	)


def test_delete_variant_file_succeeds_when_file_exists(tmp_path: Path) -> None:
	file_name = Path('foo.webp')
	file_path = tmp_path / file_name
	file_path.write_bytes(b'data')
	variant_file = _build_variant_file(file_path, file_name)

	result = _delete_variant_file(variant_file)

	assert result.action == 'delete'
	assert result.result == 'success'
	assert not file_path.exists()


def test_delete_variant_file_returns_missing_when_file_absent(tmp_path: Path) -> None:
	file_name = Path('missing.webp')
	file_path = tmp_path / file_name
	variant_file = _build_variant_file(file_path, file_name)

	result = _delete_variant_file(variant_file)

	assert result.result == 'failure'
	assert result.reason == 'file_already_missing'
