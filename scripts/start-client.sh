#!/bin/bash

PHANTOM_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG_DIR="$HOME/.phantom-client/configs"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ ! -f "$PHANTOM_DIR/phantom-client" ]; then
    echo -e "${RED}[!] phantom-client not found at: $PHANTOM_DIR${NC}"
    echo -e "${RED}    Run 'make' to build first.${NC}"
    exit 1
fi

# Check for config files
if [ ! -d "$CONFIG_DIR" ] || [ -z "$(ls -A "$CONFIG_DIR" 2>/dev/null)" ]; then
    echo -e "${RED}[!] No config files found at: $CONFIG_DIR${NC}"
    echo ""
    echo -e "${YELLOW}  Run setup first:${NC}"
    echo "    ./scripts/setup.sh"
    echo ""
    echo -e "${YELLOW}  Or manually:${NC}"
    echo "    1. Start server: ./scripts/start-server.sh"
    echo "    2. In server console: new-operator --name lab --lhost <kali-ip> --permissions all"
    echo "    3. mkdir -p $CONFIG_DIR"
    echo "    4. cp ~/phantom/lab_*.cfg $CONFIG_DIR/"
    echo ""
    exit 1
fi

# Show available configs
echo -e "${CYAN}[*] Available operator configs:${NC}"
ls "$CONFIG_DIR"/*.cfg 2>/dev/null | while read f; do
    echo "    - $(basename $f)"
done
echo ""

echo -e "${GREEN}[*] Starting Phantom C2 Client...${NC}"
echo ""

cd "$PHANTOM_DIR"
exec ./phantom-client console
