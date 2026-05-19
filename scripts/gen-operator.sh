#!/bin/bash

CONFIG_DIR="$HOME/.phantom-client/configs"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

LHOST=$(ip route get 1 2>/dev/null | awk '{print $7; exit}' || hostname -I | awk '{print $1}')

echo -e "${CYAN}[*] Generate Operator Config${NC}"
echo ""

read -p "Operator name [default: lab]: " OPERATOR_NAME
OPERATOR_NAME="${OPERATOR_NAME:-lab}"

echo -e "${CYAN}[*] Detected IP: ${LHOST}${NC}"
read -p "Server IP (lhost) [Enter to use $LHOST]: " INPUT_IP
if [ -n "$INPUT_IP" ]; then
    LHOST="$INPUT_IP"
fi

read -p "Server port [default: 31337]: " LPORT
LPORT="${LPORT:-31337}"

echo ""
echo -e "${YELLOW}[*] Run this command in your phantom-server console:${NC}"
echo ""
echo -e "${GREEN}  new-operator --name $OPERATOR_NAME --lhost $LHOST --permissions all${NC}"
echo ""
echo -e "${YELLOW}[*] Then copy the generated .cfg file:${NC}"
echo ""
echo "  mkdir -p $CONFIG_DIR"
echo "  cp ~/${OPERATOR_NAME}_${LHOST}.cfg $CONFIG_DIR/"
echo "  find ~ -name '${OPERATOR_NAME}_*.cfg' 2>/dev/null"
echo ""
echo -e "${CYAN}[*] After copying, start client with:${NC}"
echo "  ./scripts/start-client.sh"
echo ""
