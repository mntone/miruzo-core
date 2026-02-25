from fastapi import FastAPI
from starlette.requests import Request
from starlette.responses import JSONResponse, Response

from app.errors import InvalidStateError, QuotaExceededError


def _handle_quota_exceeded(_req: Request, _exc: Exception) -> Response:
	return JSONResponse(
		status_code=409,
		content={
			'error': 'quota_exceeded',
		},
	)


def _handle_invalid_state(_req: Request, _exc: Exception) -> Response:
	return JSONResponse(
		status_code=409,
		content={
			'error': 'invalid_state',
		},
	)


def register_exception_handlers(app: FastAPI) -> None:
	app.add_exception_handler(QuotaExceededError, _handle_quota_exceeded)
	app.add_exception_handler(InvalidStateError, _handle_invalid_state)
