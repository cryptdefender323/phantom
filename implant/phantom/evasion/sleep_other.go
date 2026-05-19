//go:build !windows

package evasion

import "time"

// EnableSleepObfuscation is a no-op on non-Windows platforms.
func EnableSleepObfuscation() {}

// ObfuscatedSleep falls back to normal sleep on non-Windows platforms.
func ObfuscatedSleep(ms uint32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
