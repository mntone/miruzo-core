from dataclasses import dataclass
from pathlib import Path
from typing import final

from PIL import Image as PILImage


@final
@dataclass(frozen=True, slots=True)
class ImagePathes:
	relpath_str: str
	relpath: Path
	path: Path


def new_image_file_fixture(
	tmp_path: Path,
	*,
	relative_path: str = 'l0orig/sample.png',
	image_color: float | str = 'blue',
	image_size: tuple[int, int] | list[int] = (10, 8),
) -> ImagePathes:
	origin_relpath = Path(relative_path)
	origin_path = tmp_path / origin_relpath
	origin_path.parent.mkdir(parents=True, exist_ok=True)
	PILImage.new('RGB', image_size, image_color).save(origin_path)
	return ImagePathes(
		relpath_str=relative_path,
		path=origin_path,
		relpath=origin_relpath,
	)
