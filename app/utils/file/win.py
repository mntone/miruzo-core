import ctypes
from ctypes import wintypes
from enum import IntEnum, IntFlag, auto
from pathlib import Path

advapi32 = ctypes.windll.advapi32


class _SecurityInformation(IntFlag):
	OWNER_SECURITY_INFORMATION = 0x00000001
	GROUP_SECURITY_INFORMATION = 0x00000002
	DACL_SECURITY_INFORMATION = 0x00000004
	SACL_SECURITY_INFORMATION = 0x00000008
	LABEL_SECURITY_INFORMATION = 0x00000010
	ATTRIBUTE_SECURITY_INFORMATION = 0x00000020
	SCOPE_SECURITY_INFORMATION = 0x00000040
	PROCESS_TRUST_LABEL_SECURITY_INFORMATION = 0x00000080
	ACCESS_FILTER_SECURITY_INFORMATION = 0x00000100
	BACKUP_SECURITY_INFORMATION = 0x00010000
	PROTECTED_DACL_SECURITY_INFORMATION = 0x80000000
	PROTECTED_SACL_SECURITY_INFORMATION = 0x40000000
	UNPROTECTED_DACL_SECURITY_INFORMATION = 0x20000000
	UNPROTECTED_SACL_SECURITY_INFORMATION = 0x10000000


class _SystemError(IntEnum):
	SUCCESS = 0x0
	FILE_NOT_FOUND = 0x2
	PATH_NOT_FOUND = 0x3
	TOO_MANY_OPEN_FILES = 0x4
	ACCESS_DENIED = 0x5
	NOT_ENOUGH_MEMORY = 0x8
	OUTOFMEMORY = 0xE
	INVALID_DRIVE = 0xF
	INVALID_PARAMETER = 0x57
	INVALID_NAME = 0x7B
	BAD_PATHNAME = 0xA1
	PRIVILEGE_NOT_HELD = 0x522


class _SecurityObjectType(IntEnum):
	UNKNOWN_OBJECT = 0
	FILE_OBJECT = auto()
	SERVICE = auto()
	PRINTER = auto()
	REGISTRY_KEY = auto()
	LMSHARE = auto()
	KERNEL_OBJECT = auto()
	WINDOW_OBJECT = auto()
	DS_OBJECT = auto()
	DS_OBJECT_ALL = auto()
	PROVIDER_DIFINED_OBJECT = auto()
	WMIGUID_OBJECT = auto()
	REGISTRY_WOW64_32KEY = auto()
	REGISTRY_WOW64_64KEY = auto()


class _FileAccess(IntFlag):
	READ_DATA = 0x0001  # file & pipe
	LIST_DIRECTORY = 0x0001  # directory
	WRITE_DATA = 0x0002  # file & pipe
	ADD_FILE = 0x0002  # directory
	APPEND_DATA = 0x0004  # file
	ADD_SUBDIRECTORY = 0x0004  # directory
	CREATE_PIPE_INSTANCE = 0x0004  # named pipe
	READ_EA = 0x0008  # file & directory
	WRITE_EA = 0x0010  # file & directory
	EXECUTE = 0x0020  # file
	TRAVERSE = 0x0020  # directory
	DELETE_CHILD = 0x0040  # directory
	READ_ATTRIBUTES = 0x0080  # all
	WRITE_ATTRIBUTES = 0x0100  # all


class _GENERIC_MAPPING(ctypes.Structure):
	_fields_ = [
		('GenericRead', wintypes.DWORD),
		('GenericWrite', wintypes.DWORD),
		('GenericExecute', wintypes.DWORD),
		('GenericAll', wintypes.DWORD),
	]


_FILE_GENERIC_MAPPING = _GENERIC_MAPPING(
	GenericRead=0x120089,  # FILE_GENERIC_READ
	GenericWrite=0x120116,  # FILE_GENERIC_WRITE
	GenericExecute=0x1200A0,  # FILE_GENERIC_EXECUTE
	GenericAll=0x1F01FF,  # FILE_ALL_ACCESS
)

_PSECURITY_DESCRIPTOR = wintypes.LPVOID

advapi32.GetNamedSecurityInfoW.argtypes = [
	wintypes.LPWSTR,  # path
	wintypes.DWORD,  # SE_OBJECT_TYPE
	wintypes.DWORD,  # SECURITY_INFORMATION
	ctypes.POINTER(wintypes.LPVOID),  # owner
	ctypes.POINTER(wintypes.LPVOID),  # group
	ctypes.POINTER(wintypes.LPVOID),  # dacl
	ctypes.POINTER(wintypes.LPVOID),  # sacl
	ctypes.POINTER(_PSECURITY_DESCRIPTOR),
]
advapi32.GetNamedSecurityInfoW.restype = wintypes.DWORD


class _LUID(ctypes.Structure):
	_fields_ = [('LowPart', wintypes.DWORD), ('HighPart', wintypes.LONG)]


class _LUID_AND_ATTRIBUTES(ctypes.Structure):
	_fields_ = [('Luid', _LUID), ('Attributes', wintypes.DWORD)]


class _PRIVILEGE_SET(ctypes.Structure):
	_fields_ = [
		('PrivilegeCount', wintypes.DWORD),
		('Control', wintypes.DWORD),
		('Privileges', _LUID_AND_ATTRIBUTES * 1),
	]


advapi32.AccessCheck.argtypes = [
	_PSECURITY_DESCRIPTOR,
	wintypes.HANDLE,
	wintypes.DWORD,
	ctypes.POINTER(_GENERIC_MAPPING),
	ctypes.POINTER(_PRIVILEGE_SET),
	wintypes.LPDWORD,
	wintypes.LPDWORD,
	wintypes.LPBOOL,
]
advapi32.AccessCheck.restype = wintypes.BOOL


def _get_named_security_info(
	object_name: Path,
	object_type: _SecurityObjectType,
	security_information: _SecurityInformation,
	security_descriptor: _PSECURITY_DESCRIPTOR,
) -> None:
	result = advapi32.GetNamedSecurityInfoW(
		object_name.__str__(),
		wintypes.DWORD(object_type.value),
		wintypes.DWORD(security_information.value),
		None,
		None,
		None,
		None,
		ctypes.byref(security_descriptor),
	)
	if result != _SystemError.SUCCESS:
		base = OSError(result, f'GetNamedSecurityInfoW failed: err={result} path="{object_name!s}"')
		match result:
			case _SystemError.FILE_NOT_FOUND, _SystemError.PATH_NOT_FOUND:
				raise FileNotFoundError() from base
			case _SystemError.TOO_MANY_OPEN_FILES:
				raise IOError() from base
			case (
				_SystemError.INVALID_NAME,
				_SystemError.BAD_PATHNAME,
				_SystemError.INVALID_PARAMETER,
				_SystemError.INVALID_DRIVE,
			):
				raise ValueError(f'Invalid path: {object_name!s}') from base
			case _SystemError.NOT_ENOUGH_MEMORY, _SystemError.OUTOFMEMORY:
				raise MemoryError() from base
			case _SystemError.ACCESS_DENIED, _SystemError.PRIVILEGE_NOT_HELD:
				raise PermissionError() from base
			case _:
				raise base


def _get_current_thread_effective_token() -> wintypes.HANDLE:
	return wintypes.HANDLE(-6)


def _access_check(
	security_descriptor: _PSECURITY_DESCRIPTOR,
	client_token: wintypes.HANDLE,
	desired_access: _FileAccess,
) -> tuple[bool, int]:
	privilege_set = _PRIVILEGE_SET()
	privilege_length = wintypes.DWORD(ctypes.sizeof(privilege_set))
	dword_granted_access = wintypes.DWORD()
	bool_status = wintypes.BOOL()

	res = advapi32.AccessCheck(
		security_descriptor,
		client_token,
		wintypes.DWORD(desired_access.value),
		ctypes.byref(_FILE_GENERIC_MAPPING),
		ctypes.byref(privilege_set),
		ctypes.byref(privilege_length),
		ctypes.byref(dword_granted_access),
		ctypes.byref(bool_status),
	)
	if not res:
		error = ctypes.get_last_error()
		raise OSError(error, 'AccessCheck failed')

	status = bool(bool_status.value)
	granted_access = int(dword_granted_access.value)

	return status, granted_access


def ensure_directory_access(path: Path, name: str | None = None) -> None:
	name = name or 'path'

	security_desc = _PSECURITY_DESCRIPTOR()
	try:
		_get_named_security_info(
			path,
			_SecurityObjectType.FILE_OBJECT,
			_SecurityInformation.DACL_SECURITY_INFORMATION,
			security_desc,
		)

		allowed, _ = _access_check(
			security_desc,
			_get_current_thread_effective_token(),
			_FileAccess.TRAVERSE | _FileAccess.LIST_DIRECTORY,
		)
		if not allowed:
			raise RuntimeError(name + ' is not accessible as a directory: ' + path.__str__())

	finally:
		ctypes.windll.kernel32.LocalFree(security_desc)
