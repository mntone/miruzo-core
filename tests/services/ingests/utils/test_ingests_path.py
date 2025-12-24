from pathlib import Path

import pytest

from app.config.environments import env
from app.services.ingests.utils.path import map_relative_to_output_path, resolve_origin_absolute_path


def _setup_roots(tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> Path:
	assets_root = tmp_path / 'assets'
	media_root = tmp_path / 'media'
	assets_root.mkdir()
	media_root.mkdir()

	monkeypatch.setattr(env, 'gataku_assets_root', assets_root)
	monkeypatch.setattr(env, 'media_root', media_root)
	monkeypatch.setattr(env, 'gataku_symlink_dirname', 'gataku')

	return assets_root


def test_resolve_origin_absolute_path_accepts_files_under_assets(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	relative_path = Path('foo') / 'bar.png'
	target = assets_root / relative_path
	target.parent.mkdir(parents=True)
	target.write_bytes(b'data')

	resolved = resolve_origin_absolute_path(relative_path)

	assert resolved == target


def test_resolve_origin_absolute_path_rejects_directories(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	relative_path = Path('foo')
	(assets_root / relative_path).mkdir()

	with pytest.raises(ValueError, match='must be a file'):
		resolve_origin_absolute_path(relative_path)


@pytest.mark.parametrize('input_path', [Path('/etc/passwd'), Path('/foo/bar')])
def test_resolve_origin_absolute_path_rejects_absolute(
	input_path: Path,
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	_setup_roots(tmp_path, monkeypatch)
	with pytest.raises(ValueError):
		resolve_origin_absolute_path(input_path)


@pytest.mark.parametrize('input_path', [Path('../foo'), Path('foo/../bar'), Path('..')])
def test_resolve_origin_absolute_path_rejects_path_traversal(
	input_path: Path,
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	_setup_roots(tmp_path, monkeypatch)
	with pytest.raises(ValueError):
		resolve_origin_absolute_path(input_path)


@pytest.mark.parametrize(
	'input_path',
	[
		Path(''),
		Path('.'),
		Path('foo\x00bar'),
		Path('foo|bar'),
		Path('foo '),
		Path('foo\u3000'),
		Path('foo.'),
	],
)
def test_resolve_origin_absolute_path_rejects_empty_dot_and_invalid_chars(
	input_path: Path,
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	_setup_roots(tmp_path, monkeypatch)
	with pytest.raises(ValueError):
		resolve_origin_absolute_path(input_path)


def test_map_relative_to_output_path_builds_copy_paths(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	_setup_roots(tmp_path, monkeypatch)
	target = Path('foo') / 'bar.png'

	output_path = map_relative_to_output_path(target)

	assert output_path == tmp_path / 'media' / 'l0orig' / 'foo' / 'bar.png'
