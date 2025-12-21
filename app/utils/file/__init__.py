import os

if os.name == 'nt':
	from .win import ensure_directory_access
else:
	from .nix import ensure_directory_access

__all__ = ['ensure_directory_access']
