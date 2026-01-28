#!/bin/bash

# Termius From Walmart Build Script
# This script builds the termius-from-walmart application

set -e

echo "================================"
echo "Termius From Walmart Build Script"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed$(go version)${NC}"

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then 
    echo -e "${RED}Error: Go version $REQUIRED_VERSION or higher is required${NC}"
    echo "Current version: $GO_VERSION"
    exit 1
fi

echo -e "${GREEN}✓ Go version is compatible${NC}"
echo ""

# Download dependencies
echo "Downloading dependencies..."
go mod download
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Dependencies downloaded${NC}"
else
    echo -e "${RED}✗ Failed to download dependencies${NC}"
    exit 1
fi

# Tidy up go.mod
echo "Tidying dependencies..."
go mod tidy
echo ""

# Build the application
echo "Building termius-from-walmart..."
go build -o termius-from-walmart

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Build successful!${NC}"
    echo ""
    echo "The binary 'termius-from-walmart' has been created in the current directory."
    echo ""
    echo "To run: ./termius-from-walmart"
    echo "To install: sudo mv termius-from-walmart /usr/local/bin/"
    echo ""
    
    # Check if sshpass is installed
    if ! command -v sshpass &> /dev/null; then
        echo -e "${YELLOW}Note: 'sshpass' is not installed.${NC}"
        echo "For password-based SSH connections, install sshpass:"
        echo "  Ubuntu/Debian: sudo apt-get install sshpass"
        echo "  macOS: brew install hudochenkov/sshpass/sshpass"
        echo "  Fedora: sudo dnf install sshpass"
        echo ""
        echo "Alternatively, use SSH keys (recommended)."
    else
        echo -e "${GREEN}✓ sshpass is installed${NC}"
    fi
    
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi
