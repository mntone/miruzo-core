from sqlalchemy.ext.compiler import compiles
from sqlalchemy.sql.functions import FunctionElement


class _EmptyJsonArrayDefault(FunctionElement):
	inherit_cache = True


@compiles(_EmptyJsonArrayDefault, 'mysql')
def _compile_empty_json_array_default_mysql(
	_element: _EmptyJsonArrayDefault,
	_compiler: object,
	**_kwargs: object,
) -> str:
	return 'JSON_ARRAY()'


@compiles(_EmptyJsonArrayDefault, 'postgresql')
def _compile_empty_json_array_default_postgresql(
	_element: _EmptyJsonArrayDefault,
	_compiler: object,
	**_kwargs: object,
) -> str:
	return "'[]'::jsonb"


@compiles(_EmptyJsonArrayDefault)
def _compile_empty_json_array_default_default(
	_element: _EmptyJsonArrayDefault,
	_compiler: object,
	**_kwargs: object,
) -> str:
	return "'[]'"


def empty_json_array_default() -> FunctionElement:
	return _EmptyJsonArrayDefault()
