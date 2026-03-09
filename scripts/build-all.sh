#!/usr/bin/env bash

# build-all.sh - Convenience script to build everything
# This script builds all tools for multiple platforms

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Building all tools for all platforms...${NC}"
echo ""

# Build for current platform
echo -e "${GREEN}Building for current platform...${NC}"
./scripts/build.sh all

echo ""
echo -e "${GREEN}Building release binaries...${NC}"
./scripts/build.sh -r all

# Cross-compile for other platforms (optional)
if [[ "$1" == "--all-platforms" ]]; then
    echo ""
    echo -e "${GREEN}Cross-compiling for Linux...${NC}"
    GOOS=linux GOARCH=amd64 ./scripts/build.sh -r all -o bin/linux
    
    echo ""
    echo -e "${GREEN}Cross-compiling for Windows...${NC}"
    GOOS=windows GOARCH=amd64 ./scripts/build.sh -r all -o bin/windows
    
    echo ""
    echo -e "${GREEN}Cross-compiling for macOS...${NC}"
    GOOS=darwin GOARCH=amd64 ./scripts/build.sh -r all -o bin/darwin
    GOOS=darwin GOARCH=arm64 ./scripts/build.sh -r all -o bin/darwin-arm64
fi

echo ""
echo -e "${GREEN}✨ All builds complete!${NC}"