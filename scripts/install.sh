#!/usr/bin/env bash

# install.sh - Install spoke-tool binaries to GOPATH/bin

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Get GOPATH
GOPATH=${GOPATH:-$(go env GOPATH)}
INSTALL_DIR="$GOPATH/bin"

echo -e "${YELLOW}Installing spoke-tool to $INSTALL_DIR${NC}"
echo ""

# Build first
echo -e "${GREEN}Building binaries...${NC}"
./scripts/build.sh -r all

echo ""
echo -e "${GREEN}Installing...${NC}"

# Copy binaries
cp bin/readmegen "$INSTALL_DIR/"
cp bin/testgen "$INSTALL_DIR/"

# Make executable
chmod +x "$INSTALL_DIR/readmegen"
chmod +x "$INSTALL_DIR/testgen"

echo -e "${GREEN}✅ Installed to $INSTALL_DIR${NC}"
echo ""
echo "You can now run:"
echo "  readmegen -h"
echo "  testgen -h"