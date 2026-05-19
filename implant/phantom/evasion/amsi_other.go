//go:build !windows

package evasion

// PatchAMSI is a no-op on non-Windows platforms.
func PatchAMSI() error { return nil }

// PatchETW is a no-op on non-Windows platforms.
func PatchETW() error { return nil }

// PatchAll is a no-op on non-Windows platforms.
func PatchAll() {}
