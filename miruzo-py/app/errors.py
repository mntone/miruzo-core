class DomainError(RuntimeError):
	"""Domain base error"""


class InvariantViolationError(DomainError):
	"""Raised when a domain invariant is violated."""


class SingletonUserMissingError(InvariantViolationError):
	"""Raised when singleton user row is missing."""
