package timezone

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}

	return e
}

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetDynamicTimeZoneInformation = modkernel32.NewProc("GetDynamicTimeZoneInformation")
)

type DynamicTimeZoneInformation struct {
	Bias                        int32
	StandardName                [32]uint16
	StandardDate                windows.Systemtime
	StandardBias                int32
	DaylightName                [32]uint16
	DaylightDate                windows.Systemtime
	DaylightBias                int32
	TimeZoneKeyName             [128]uint16
	DynamicDaylightTimeDisabled byte
}

func GetDynamicTimeZoneInformation(tzi *DynamicTimeZoneInformation) (rc uint32, err error) {
	r0, _, e1 := syscall.SyscallN(procGetDynamicTimeZoneInformation.Addr(), uintptr(unsafe.Pointer(tzi)))
	rc = uint32(r0)
	if rc == 0xffffffff {
		err = errnoErr(e1)
	}
	return
}
