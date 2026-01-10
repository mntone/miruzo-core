from pathlib import Path

from importers.common.models import GatakuImageRow
from importers.common.origin import OriginResolver


def test_origin_resolver_returns_resolution_for_existing_asset(tmp_path: Path) -> None:
	gataku_root = tmp_path / 'gataku'
	assets_root = gataku_root / 'out' / 'downloads'
	assets_root.mkdir(parents=True)

	relative_path = Path('subdir/sample.jpg')
	asset_path = assets_root / relative_path
	asset_path.parent.mkdir(parents=True, exist_ok=True)
	asset_path.write_text('x', encoding='utf-8')

	row = GatakuImageRow(
		filepath=Path('out/downloads') / relative_path,
		sha256='x' * 64,
		created_at=None,
	)

	resolver = OriginResolver(gataku_root=gataku_root, gataku_assets_root=assets_root)
	resolution = resolver.resolve(row)

	assert resolution is not None
	assert resolution.src_path == asset_path
	assert resolution.origin_relative_path == relative_path


def test_origin_resolver_accepts_absolute_path(tmp_path: Path) -> None:
	assets_root = tmp_path / 'gataku' / 'out' / 'downloads'
	assets_root.mkdir(parents=True)

	asset_path = assets_root / 'absolute.jpg'
	asset_path.write_text('x', encoding='utf-8')

	row = GatakuImageRow(
		filepath=asset_path,
		sha256='x' * 64,
		created_at=None,
	)

	resolver = OriginResolver(gataku_root=tmp_path, gataku_assets_root=assets_root)
	resolution = resolver.resolve(row)

	assert resolution is not None
	assert resolution.src_path == asset_path
	assert resolution.origin_relative_path == Path('absolute.jpg')


def test_origin_resolver_counts_missing_file(tmp_path: Path) -> None:
	assets_root = tmp_path / 'gataku' / 'out' / 'downloads'
	assets_root.mkdir(parents=True)

	row = GatakuImageRow(
		filepath=Path('missing.jpg'),
		sha256='x' * 64,
		created_at=None,
	)

	resolver = OriginResolver(gataku_root=tmp_path, gataku_assets_root=assets_root)
	resolution = resolver.resolve(row)

	assert resolution is None


def test_origin_resolver_rejects_asset_outside_root(tmp_path: Path) -> None:
	gataku_root = tmp_path / 'gataku'
	gataku_root.mkdir()

	assets_root = tmp_path / 'gataku' / 'out' / 'downloads'
	assets_root.mkdir(parents=True)

	other_root = tmp_path / 'other'
	other_root.mkdir()
	asset_path = other_root / 'outside.jpg'
	asset_path.write_text('x', encoding='utf-8')

	row = GatakuImageRow(
		filepath=asset_path,
		sha256='x' * 64,
		created_at=None,
	)

	resolver = OriginResolver(gataku_root=gataku_root, gataku_assets_root=assets_root)
	resolution = resolver.resolve(row)

	assert resolution is None
