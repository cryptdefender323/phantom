// Package dylib converts macOS dylib/shared libraries into position-independent shellcode.
package dylib

import (
	"github.com/cryptdefender3232/beignet"
)

// Options controls shellcode generation from a dylib.
type Options struct {
	EntrySymbol string
	Compress    bool
}

// ToShellcode converts a dylib byte slice into shellcode using the provided options.
func ToShellcode(dylibData []byte, opts Options) ([]byte, error) {
	return beignet.DylibToShellcode(dylibData, beignet.Options{
		EntrySymbol: opts.EntrySymbol,
		Compress:    opts.Compress,
	})
}

// FileToShellcode reads a dylib from disk and converts it into shellcode.
func FileToShellcode(path string, opts Options) ([]byte, error) {
	return beignet.DylibFileToShellcode(path, beignet.Options{
		EntrySymbol: opts.EntrySymbol,
		Compress:    opts.Compress,
	})
}
