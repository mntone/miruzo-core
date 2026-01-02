from dataclasses import dataclass, field
from datetime import datetime
from typing import ClassVar


@dataclass(slots=True)
class StubJob:
	_NAME: ClassVar[str] = 'stub_job'

	run_called_with: list[datetime] = field(default_factory=list[datetime])

	@property
	def name(self) -> str:
		return StubJob._NAME

	def run(self, *, evaluated_at: datetime) -> None:
		self.run_called_with.append(evaluated_at)
