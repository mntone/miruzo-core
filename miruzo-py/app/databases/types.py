from sqlalchemy import JSON, Integer, SmallInteger
from sqlalchemy.dialects import mysql, postgresql

LEAST8_INT = SmallInteger().with_variant(mysql.TINYINT, 'mysql').with_variant(Integer, 'sqlite')

JSON_VALUE = JSON(none_as_null=True).with_variant(postgresql.JSONB(none_as_null=True), 'postgresql')
