#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="zhaori96/luamerge"
BINARY_NAME="luamerge"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${VERSION:-latest}"

# Show usage information
usage() {
    cat << EOF
${BLUE}luamerge installer${NC}

Usage:
    curl -sSL https://raw.githubusercontent.com/$REPO/main/install.sh | bash
    curl -sSL https://raw.githubusercontent.com/$REPO/main/install.sh | VERSION=v1.0.0 bash

Environment variables:
    VERSION       Version to install (default: latest)
                  Examples: latest, v1.0.0, v1.2.3
    INSTALL_DIR   Installation directory (default: \$HOME/.local/bin)

Examples:
    # Install latest version
    curl -sSL https://raw.githubusercontent.com/$REPO/main/install.sh | bash

    # Install specific version
    curl -sSL https://raw.githubusercontent.com/$REPO/main/install.sh | VERSION=v1.0.0 bash

    # Install to custom directory
    curl -sSL https://raw.githubusercontent.com/$REPO/main/install.sh | INSTALL_DIR=/usr/local/bin bash

EOF
}

# Parse command line arguments
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    usage
    exit 0
fi

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case $OS in
    linux|darwin)
        ;;
    mingw*|msys*|cygwin*)
        OS="windows"
        BINARY_NAME="${BINARY_NAME}.exe"
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Installing luamerge...${NC}"
echo "OS: $OS"
echo "Architecture: $ARCH"
echo "Install directory: $INSTALL_DIR"
echo ""

# Determine version to install
if [ "$VERSION" = "latest" ]; then
    echo -e "${YELLOW}Fetching latest release...${NC}"
    RELEASE_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$RELEASE_VERSION" ]; then
        echo -e "${RED}Failed to fetch latest release${NC}"
        exit 1
    fi

    echo "Latest version: $RELEASE_VERSION"
else
    RELEASE_VERSION="$VERSION"
    echo "Installing version: $RELEASE_VERSION"

    # Validate that the version exists
    echo -e "${YELLOW}Checking if version exists...${NC}"
    STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://api.github.com/repos/$REPO/releases/tags/$RELEASE_VERSION")

    if [ "$STATUS_CODE" != "200" ]; then
        echo -e "${RED}Version $RELEASE_VERSION not found${NC}"
        echo "Please check available versions at: https://github.com/$REPO/releases"
        exit 1
    fi
fi

# Construct download URL for archive
ARCH_NAME="${ARCH}"
if [ "$ARCH" = "amd64" ]; then
    ARCH_NAME="x86_64"
elif [ "$ARCH" = "arm" ]; then
    ARCH_NAME="armv7"
fi

OS_NAME="$(echo $OS | sed 's/\b\(.\)/\u\1/')" # Capitalize first letter

if [ "$OS" = "windows" ]; then
    ARCHIVE_NAME="${BINARY_NAME}_${RELEASE_VERSION#v}_Windows_${ARCH_NAME}.zip"
else
    ARCHIVE_NAME="${BINARY_NAME}_${RELEASE_VERSION#v}_${OS_NAME}_${ARCH_NAME}.tar.gz"
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$RELEASE_VERSION/$ARCHIVE_NAME"

echo -e "${YELLOW}Downloading from: $DOWNLOAD_URL${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download archive
if ! curl -L -o "$TMP_DIR/archive" "$DOWNLOAD_URL"; then
    echo -e "${RED}Failed to download archive${NC}"
    echo "URL: $DOWNLOAD_URL"
    exit 1
fi

# Extract binary
echo -e "${YELLOW}Extracting binary...${NC}"
cd "$TMP_DIR"
if [ "$OS" = "windows" ]; then
    unzip -q archive
else
    tar -xzf archive
fi

# Make binary executable
chmod +x "$BINARY_NAME"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Install binary
echo -e "${YELLOW}Installing to $INSTALL_DIR/$BINARY_NAME${NC}"
mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

# Check if install directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo "Add the following line to your ~/.bashrc or ~/.zshrc:"
    echo ""
    echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
fi

echo -e "${GREEN}✓ luamerge $RELEASE_VERSION installed successfully!${NC}"
echo ""
echo "Run 'luamerge --help' to get started"
