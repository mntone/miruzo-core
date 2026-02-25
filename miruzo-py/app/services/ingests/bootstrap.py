from collections.abc import Iterable
from pathlib import Path

from app.config.environments import Settings
from app.config.variant import VariantLayerSpec


def _ensure_media_root(media_root: Path) -> None:
	if media_root.exists():
		if media_root.is_symlink():
			raise RuntimeError('media_root must not be a symlink')
		if media_root.is_file():
			raise RuntimeError('media_root must be a directory')
	else:
		media_root.mkdir(parents=True, exist_ok=True)


def _ensure_original_root(media_root: Path) -> None:
	original_root = media_root / 'l0orig'
	if original_root.exists():
		if original_root.is_symlink():
			raise RuntimeError('Original directory must not be a symlink: l0orig')
		if original_root.is_file():
			raise RuntimeError('Original directory must be a directory: l0orig')
	else:
		original_root.mkdir(parents=True)


def _ensure_variant_roots(
	media_root: Path,
	layers: Iterable[VariantLayerSpec],
) -> None:
	seen: set[str] = set()
	for layer in layers:
		for spec in layer.specs:
			variant_dirname = spec.slot.key
			if variant_dirname in seen:
				continue

			variant_relpath = media_root / variant_dirname
			if variant_relpath.exists():
				if variant_relpath.is_symlink():
					raise RuntimeError('Variant directory must not be a symlink: ' + variant_dirname)
				if variant_relpath.is_file():
					raise RuntimeError('Variant directory must be a directory: ' + variant_dirname)
			else:
				variant_relpath.mkdir(parents=True)

			seen.add(variant_dirname)


def ensure_ingest_layout(env: Settings) -> None:
	"""Ensure ingest-related directories exist and are not symlinked."""
	_ensure_media_root(env.media_root)
	_ensure_original_root(env.media_root)
	_ensure_variant_roots(env.media_root, env.variant_layers)
