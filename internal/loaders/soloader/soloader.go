// Package soloader converts Linux ELF shared objects into position-independent shellcode.
package soloader

import (
	"github.com/cryptdefender3232/malasada"
)

// ConvertSharedObject converts a Linux .so file into shellcode that calls the named export.
// Set compress to true to apply aplib compression to the output.
func ConvertSharedObject(soPath string, exportName string, compress bool) ([]byte, error) {
	return malasada.ConvertSharedObject(soPath, exportName, compress)
}
