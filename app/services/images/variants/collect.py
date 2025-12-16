from collections.abc import Iterable, Iterator
from pathlib import Path

from app.services.images.variants.path import NormalizedRelativePath, VariantDirectoryPath
from app.services.images.variants.types import VariantFile
from app.services.images.variants.utils import load_image_info, parse_variant_slotkey

_NON_VARIANT_SUFFIXES = ('orig',)


def collect_variant_directories(media_root: Path) -> Iterator[str]:
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


def _collect_variant_file(path: Path, variant_dirname: str) -> VariantFile | None:
	info = load_image_info(path)
	if info is None:
		return None

	file = VariantFile(variant_dirname, info)
	return file


def collect_variant_files(
	variant_dirpaths: Iterable[VariantDirectoryPath],
	*,
	rel_to: NormalizedRelativePath,
) -> Iterator[VariantFile]:
	"""Yield VariantFile objects for files under the provided variant dirs."""

	# Normalize argument name for internal use
	relative_path = rel_to

	for variant_dirpath in variant_dirpaths:
		output_path = variant_dirpath / relative_path.parent
		if not output_path.is_dir():
			continue

		variant_dirname = variant_dirpath.name
		output_name = relative_path.name

		for file_path in output_path.glob(f'{output_name}.*'):
			if (variant := _collect_variant_file(file_path, variant_dirname)) is not None:
				yield variant


def normalize_variant_directories(
	variant_dirnames: Iterable[str],
	*,
	under: Path,
) -> Iterator[VariantDirectoryPath]:
	# Normalize argument name for internal use
	media_root = under

	for variant_dirname in variant_dirnames:
		try:
			_ = parse_variant_slotkey(variant_dirname)
		except ValueError:
			continue  # or raise, depending on policy

		# NOTE:
		# The result of `collect_variant_directories` is treated as a trust boundary.
		# All filesystem-level validation and filtering is intentionally performed
		# during collection, so this stage assumes the input to be structurally sound.
		#
		# The responsibility here is limited to interpreting the directory name as a
		# variant slot key and materializing the corresponding path, without re-checking
		# filesystem properties.
		variant_dirpath = media_root / variant_dirname

		yield VariantDirectoryPath(variant_dirpath)
