//go:build !windows

package evasion

// InitIndirectSyscalls is a no-op on non-Windows platforms.
func InitIndirectSyscalls() error { return nil }

// GetSyscallNumber always returns not-found on non-Windows.
func GetSyscallNumber(name string) (uint16, bool) { return 0xFFFF, false }
