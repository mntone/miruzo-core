class DomainError(RuntimeError):
	"""Domain base error"""


class InvalidStateError(DomainError):
	"""Raised when state is invalid."""


class InvariantViolationError(DomainError):
	"""Raised when a domain invariant is violated."""


class QuotaExceededError(DomainError):
	"""Raised when quota exceeded"""


class SingletonUserMissingError(InvariantViolationError):
	"""Raised when singleton user row is missing."""
