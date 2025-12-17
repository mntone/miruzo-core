from pydantic import BaseModel
from starlette.responses import Response


def build_response(
	data: BaseModel,
	*,
	status_code: int = 200,
	exclude_defaults: bool = True,
) -> Response:
	payload = data.model_dump_json(
		exclude_defaults=exclude_defaults,
		exclude_none=True,
	)
	return Response(
		content=payload,
		status_code=status_code,
		media_type='application/json',
	)
