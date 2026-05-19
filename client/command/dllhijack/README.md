# client/command/dllhijack

## Overview

Implements the 'dllhijack' command group for the Phantom client console.

## Go Files

- `commands.go` – Registers the dllhijack command, its flags, and completions for selecting remote DLL targets.
- `dllhijack.go` – Performs the RPC call that plants a Phantom DLL alongside a target executable for hijacking.
