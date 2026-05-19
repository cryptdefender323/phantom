//go:build !windows

package evasion

// SpawnWithSpoofedParent is a no-op on non-Windows platforms.
func SpawnWithSpoofedParent(targetExe string, parentName string) (*SpoofedProcess, error) {
	return nil, nil
}

// SpoofedProcess stub for non-Windows builds.
type SpoofedProcess struct {
	PID    uint32
	Handle uintptr
}
