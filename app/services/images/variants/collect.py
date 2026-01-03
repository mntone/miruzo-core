import logging
from collections.abc import Iterable, Iterator
from pathlib import Path

from PIL import Image as PILImage
from PIL import UnidentifiedImageError as PILUnidentifiedImageError

from app.services.images.variants.path import (
	VariantBasePath,
	build_absolute_path,
	build_variant_relative_path,
)
from app.services.images.variants.types import FileInfo, VariantFile, VariantRelativePath
from app.services.images.variants.utils import get_image_info, parse_variant_slot

log = logging.getLogger(__name__)

_NON_VARIANT_SUFFIXES = ('orig',)


def collect_variant_directories(media_root: Path) -> Iterator[str]:
	"""
	Collect existing variant directories under media_root.

	This scans media_root and returns directories that look like
	variant directories (e.g. l0w320, l1w640), excluding symlinks
	and non-directory entries.

	This function does NOT validate naming rules strictly.
	"""

	entries = media_root.iterdir()
	for entry in entries:
		# skip non-directories
		if not entry.is_dir():
			continue

		# skip symlinks (e.g. gataku original link)
		if entry.is_symlink():
			continue

		# skip obvious non-variant directories
		# (e.g. "l0orig", "l9orig", etc.)
		if entry.name.endswith(_NON_VARIANT_SUFFIXES):
			continue

		relative_path = entry.relative_to(media_root)
		variant_dirname = relative_path.name
		yield variant_dirname


def _load_variant_file(
	absolute_path: Path,
	relative_path: VariantRelativePath,
	variant_dirname: str,
) -> VariantFile | None:
	try:
		stat = absolute_path.stat()
	except FileNotFoundError:
		log.debug('image not found: %s', absolute_path)
		return None

	try:
		with PILImage.open(absolute_path) as image:
			info = get_image_info(image)
	except FileNotFoundError:
		log.debug('image not found: %s', absolute_path)
		return None
	except PermissionError:
		log.warning('permission denied: %s', absolute_path)
		return None
	except PILUnidentifiedImageError:
		log.warning('unknown image format: %s', absolute_path)
		return None
	except OSError:
		log.warning('invalid image: %s', absolute_path)
		return None

	file_info = FileInfo(
		absolute_path=absolute_path,
		relative_path=relative_path,
		bytes=stat.st_size,
	)

	file = VariantFile(
		file_info=file_info,
		image_info=info,
		variant_dir=variant_dirname,
	)

	return file


def collect_variant_files(
	media_relpaths: Iterable[VariantRelativePath],
	*,
	under: Path,
) -> Iterator[VariantFile]:
	"""Yield VariantFile objects for files under the provided variant dirs."""

	# Normalize argument name for internal use
	media_root = under

	for relative_path in media_relpaths:
		target_base = build_absolute_path(relative_path, under=media_root)
		output_path = target_base.parent
		if not output_path.is_dir():
			continue

		variant_dirname = relative_path.parts[0]
		output_name = target_base.name

		for absolute_path in output_path.glob(f'{output_name}.*'):
			variant_file = _load_variant_file(absolute_path, relative_path, variant_dirname)
			if variant_file is not None:
				yield variant_file


def normalize_media_relative_paths(
	relative_path: VariantBasePath,
	*,
	under: Iterable[str],
) -> Iterator[VariantRelativePath]:
	# Normalize argument name for internal use
	variant_dirnames = under

	for variant_dirname in variant_dirnames:
		try:
			_ = parse_variant_slot(variant_dirname)
		except ValueError:
			continue

		variant_path = build_variant_relative_path(relative_path, under=variant_dirname)

		yield variant_path
