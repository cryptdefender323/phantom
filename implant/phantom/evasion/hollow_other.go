//go:build !windows

package evasion

// InjectShellcode is a no-op on non-Windows platforms.
func InjectShellcode(shellcode []byte, hostProcess string, parentName string) error {
	return nil
}
