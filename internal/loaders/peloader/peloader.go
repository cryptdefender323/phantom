// Package peloader converts Windows PE files into position-independent shellcode.
package peloader

import (
	"context"

	wasmdonut "github.com/cryptdefender3232/wasm-donut"
)

// Architecture constants.
const (
	ArchX86 = 1 // 32-bit x86
	ArchX64 = 2 // 64-bit x86
	ArchX84 = 3 // x86+x64 combined stub
)

// Entropy level constants.
const (
	EntropyNone    = 1
	EntropyRandom  = 2
	EntropyDefault = 3
)

// Compression constants.
const (
	CompressNone  = 1
	CompressAPLib = 2
)

// Exit behaviour constants.
const (
	ExitThread  = 1
	ExitProcess = 2
	ExitBlock   = 3
)

// AMSI/ETW bypass constants.
const (
	BypassNone     = 1
	BypassAbort    = 2
	BypassContinue = 3
)

// PE header handling constants.
const (
	HeadersOverwrite = 1
	HeadersKeep      = 2
)

// GenerateOptions controls shellcode generation from a PE binary.
type GenerateOptions struct {
	Ext      string
	Args     string
	Class    string
	Method   string
	Domain   string
	Runtime  string
	Arch     int
	Bypass   int
	Headers  int
	Entropy  int
	Compress int
	ExitOpt  int
	Thread   bool
	Unicode  bool
	OEP      uint32
}

// GenerateResult holds the output of a shellcode generation run.
type GenerateResult struct {
	Loader []byte
}

// Generate converts a PE binary (exe or dll) into shellcode.
// ext should be ".exe" or ".dll".
func Generate(ctx context.Context, input []byte, ext string, opts GenerateOptions) (GenerateResult, error) {
	upstream := wasmdonut.GenerateOptions{
		Ext:      opts.Ext,
		Args:     opts.Args,
		Class:    opts.Class,
		Method:   opts.Method,
		Domain:   opts.Domain,
		Runtime:  opts.Runtime,
		Arch:     opts.Arch,
		Bypass:   opts.Bypass,
		Headers:  opts.Headers,
		Entropy:  opts.Entropy,
		Compress: opts.Compress,
		ExitOpt:  opts.ExitOpt,
		Thread:   opts.Thread,
		Unicode:  opts.Unicode,
		OEP:      opts.OEP,
	}
	res, err := wasmdonut.Generate(ctx, input, ext, upstream)
	if err != nil {
		return GenerateResult{}, err
	}
	return GenerateResult{Loader: res.Loader}, nil
}
