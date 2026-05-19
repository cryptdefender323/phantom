#!/bin/bash

set -e

PHANTOM_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG_DIR="$HOME/.phantom-client/configs"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${CYAN}"
echo "  ██████╗ ██╗  ██╗ █████╗ ███╗   ██╗████████╗ ██████╗ ███╗   ███╗"
echo "  ██╔══██╗██║  ██║██╔══██╗████╗  ██║╚══██╔══╝██╔═══██╗████╗ ████║"
echo "  ██████╔╝███████║███████║██╔██╗ ██║   ██║   ██║   ██║██╔████╔██║"
echo "  ██╔═══╝ ██╔══██║██╔══██║██║╚██╗██║   ██║   ██║   ██║██║╚██╔╝██║"
echo "  ██║     ██║  ██║██║  ██║██║ ╚████║   ██║   ╚██████╔╝██║ ╚═╝ ██║"
echo "  ╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝ ╚═╝     ╚═╝"
echo -e "${NC}"
echo -e "${YELLOW}  C2 Framework - First Time Setup${NC}"
echo ""

# Check binary exists
if [ ! -f "$PHANTOM_DIR/phantom-server" ]; then
    echo -e "${RED}[!] phantom-server not found. Run 'make' first.${NC}"
    exit 1
fi

# Detect local IP
LHOST=$(ip route get 1 2>/dev/null | awk '{print $7; exit}' || hostname -I | awk '{print $1}')
echo -e "${CYAN}[*] Detected local IP: ${LHOST}${NC}"
read -p "    Use this IP? (Enter to confirm, or type new IP): " INPUT_IP
if [ -n "$INPUT_IP" ]; then
    LHOST="$INPUT_IP"
fi

# Operator name
read -p "[*] Operator name [default: operator]: " OPERATOR_NAME
OPERATOR_NAME="${OPERATOR_NAME:-operator}"

echo ""
echo -e "${CYAN}[*] Starting server temporarily to generate config...${NC}"

# Create config dir
mkdir -p "$CONFIG_DIR"

# Start server in background, wait for it to be ready, then generate operator
"$PHANTOM_DIR/phantom-server" &
SERVER_PID=$!

echo -e "${CYAN}[*] Waiting for server to start...${NC}"
sleep 3

# Check server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo -e "${RED}[!] Server failed to start.${NC}"
    exit 1
fi

# Use expect-style approach: pipe commands to server stdin
CONFIG_OUTPUT=$(mktemp)
echo "new-operator --name $OPERATOR_NAME --lhost $LHOST --permissions all" | \
    timeout 10 "$PHANTOM_DIR/phantom-server" 2>&1 | tee "$CONFIG_OUTPUT" || true

# Kill background server
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

# Find generated config
CFG_FILE=$(find "$PHANTOM_DIR" -maxdepth 2 -name "${OPERATOR_NAME}_*.cfg" 2>/dev/null | head -1)
if [ -z "$CFG_FILE" ]; then
    CFG_FILE=$(find ~ -maxdepth 3 -name "${OPERATOR_NAME}_*.cfg" 2>/dev/null | head -1)
fi

if [ -n "$CFG_FILE" ]; then
    cp "$CFG_FILE" "$CONFIG_DIR/"
    echo -e "${GREEN}[+] Config copied to: $CONFIG_DIR/$(basename $CFG_FILE)${NC}"
else
    echo -e "${YELLOW}[!] Config not auto-detected. Manual steps:${NC}"
    echo "    1. Run: ./phantom-server"
    echo "    2. In server console: new-operator --name $OPERATOR_NAME --lhost $LHOST --permissions all"
    echo "    3. Copy the generated .cfg to: $CONFIG_DIR/"
fi

rm -f "$CONFIG_OUTPUT"

echo ""
echo -e "${GREEN}[+] Setup complete!${NC}"
echo ""
echo "  Next steps:"
echo "    1. Start server:  ./scripts/start-server.sh"
echo "    2. Start client:  ./scripts/start-client.sh"
echo ""
