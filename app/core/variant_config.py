from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class VariantFormat:
	"""Image encoding details used when emitting a variant."""

	name: str
	codecs: str | None
	extension: str
	lossless: bool = False
	default_quality: int | None = None


@dataclass(frozen=True, slots=True)
class VariantSpec:
	"""Concrete thumbnail definition (size + format + output directory)."""

	label: str
	width: int
	format: VariantFormat
	quality: int | None = None
	required: bool = False


@dataclass(frozen=True, slots=True)
class VariantLayer:
	"""Logical layer (e.g. primary, fallback) composed of several specs."""

	name: str
	layer_id: int
	specs: tuple[VariantSpec, ...]


WEBP_FORMAT = VariantFormat(
	name='webp',
	codecs='vp8',
	extension='.webp',
	lossless=False,
	default_quality=80,
)

LOSSLESS_WEBP_FORMAT = VariantFormat(
	name='webp',
	codecs='vp8l',
	extension='.webp',
	lossless=True,
	default_quality=80,
)

JPEG_FORMAT = VariantFormat(
	name='jpeg',
	codecs=None,
	extension='.jpg',
	lossless=False,
	default_quality=85,
)


def _spec(
	fmt: VariantFormat,
	*,
	width: int,
	quality: int | None = None,
	label: str | None = None,
	required: bool = False,
) -> VariantSpec:
	if label is None:
		label = f'w{width}'

	return VariantSpec(
		label=label,
		width=width,
		format=fmt,
		quality=quality if quality is not None else fmt.default_quality,
		required=required,
	)


DEFAULT_VARIANT_LAYERS: tuple[VariantLayer, ...] = (
	VariantLayer(
		name='primary',
		layer_id=1,
		specs=(
			_spec(WEBP_FORMAT, width=320, quality=80, required=True),
			_spec(WEBP_FORMAT, width=480, quality=70),
			_spec(WEBP_FORMAT, width=640, quality=60),
			_spec(WEBP_FORMAT, width=960, quality=50),
			_spec(WEBP_FORMAT, width=1120, quality=40),
		),
	),
	VariantLayer(
		name='fallback',
		layer_id=9,
		specs=(_spec(JPEG_FORMAT, width=320, quality=85, required=True),),
	),
)
