# client/command/backdoor

## Overview

Implements the 'backdoor' command group for the Phantom client console.

## Go Files

- `backdoor.go` – Implements the RPC workflow that backdoors a remote binary with Phantom shellcode for the active session.
- `commands.go` – Registers the Cobra command, binds flags/completions, and restricts usage to applicable targets.
