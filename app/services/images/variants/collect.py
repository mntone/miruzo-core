from collections.abc import Iterable, Iterator
from pathlib import Path

from app.services.images.variants.types import VariantFile
from app.services.images.variants.utils import load_image_info

_NON_VARIANT_SUFFIXES = ('orig',)


def collect_variant_directories(
	*,
	media_root: Path,
) -> Iterator[str]:
	"""
	Collect existing variant directories under media_root.

	This scans media_root and returns directories that look like
	variant directories (e.g. l0w320, l1w640), excluding symlinks
	and non-directory entries.

	This function does NOT validate naming rules strictly.
	"""

	try:
		entries = media_root.iterdir()
	except FileNotFoundError:
		return

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


def _collect_variant_file(
	file_path: Path,
	*,
	variant_dir: str,
	relative_path: Path,
) -> VariantFile | None:
	image_info = load_image_info(file_path)
	if image_info is None:
		return None

	file = VariantFile(
		variant_dir=variant_dir,
		relative_path=relative_path,
		file_info=image_info,
	)
	return file


def collect_variant_files(
	*,
	media_root: Path,
	variant_dirs: Iterable[str],
	relative_path_noext: Path,
) -> Iterator[VariantFile]:
	"""Yield VariantFile objects for files under the provided variant dirs."""

	for variant_dir in variant_dirs:
		output_path = media_root / variant_dir / relative_path_noext.parent
		if not output_path.is_dir():
			continue

		output_name = relative_path_noext.name

		for file in output_path.glob(f'{output_name}.*'):
			if (
				variant := _collect_variant_file(
					file,
					variant_dir=variant_dir,
					relative_path=relative_path_noext,
				)
			) is not None:
				yield variant
