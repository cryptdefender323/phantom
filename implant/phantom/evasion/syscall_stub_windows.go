package evasion

// directSyscall invokes a Windows NT syscall directly using the syscall
// number (SSN), bypassing any userland hooks in ntdll.dll.

import (
	"syscall"
)

// directSyscall executes a Windows NT syscall with up to 15 arguments.
// Syscall15 signature: (trap, nargs, a1..a15) → (r1, r2, errno)
func directSyscall(ssn uintptr, args ...uintptr) (uintptr, uintptr, error) {
	var a [15]uintptr
	copy(a[:], args)

	r1, r2, errno := syscall.Syscall15(
		ssn,
		uintptr(len(args)),
		a[0], a[1], a[2], a[3], a[4],
		a[5], a[6], a[7], a[8], a[9],
		a[10], a[11], a[12], a[13], a[14],
	)

	var err error
	if errno != 0 {
		err = errno
	}
	return r1, r2, err
}
