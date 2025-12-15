from __future__ import annotations

from dataclasses import dataclass
from typing import Literal


@dataclass(frozen=True, slots=True)
class VariantFormat:
	"""Image encoding details used when emitting a variant."""

	container: Literal['jpeg', 'webp']
	codecs: str | None
	file_extension: str
	default_quality: int | None = None


@dataclass(frozen=True, slots=True)
class VariantSlotkey:
	layer_id: int
	width: int

	@property
	def label(self) -> str:
		return f'l{self.layer_id}w{self.width}'


@dataclass(frozen=True, slots=True)
class VariantSpec:
	"""Concrete thumbnail definition (size + format + output directory)."""

	slotkey: VariantSlotkey
	layer_id: int
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
	container='webp',
	codecs='vp8',
	file_extension='.webp',
	default_quality=80,
)

LOSSLESS_WEBP_FORMAT = VariantFormat(
	container='webp',
	codecs='vp8l',
	file_extension='.webp',
	default_quality=80,
)

JPEG_FORMAT = VariantFormat(
	container='jpeg',
	codecs=None,
	file_extension='.jpg',
	default_quality=85,
)


def _spec(
	fmt: VariantFormat,
	*,
	layer_id: int,
	width: int,
	quality: int | None = None,
	required: bool = False,
) -> VariantSpec:
	return VariantSpec(
		slotkey=VariantSlotkey(layer_id, width),
		layer_id=layer_id,
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
			_spec(WEBP_FORMAT, layer_id=1, width=320, quality=80, required=True),
			_spec(WEBP_FORMAT, layer_id=1, width=480, quality=70),
			_spec(WEBP_FORMAT, layer_id=1, width=640, quality=60),
			_spec(WEBP_FORMAT, layer_id=1, width=960, quality=50),
			_spec(WEBP_FORMAT, layer_id=1, width=1120, quality=40),
		),
	),
	VariantLayer(
		name='fallback',
		layer_id=9,
		specs=(_spec(JPEG_FORMAT, layer_id=9, width=320, quality=85, required=True),),
	),
)
