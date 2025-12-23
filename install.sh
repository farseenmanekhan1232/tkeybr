#!/bin/bash
# tkeybr installer - downloads and installs the latest release
# Usage: curl -fsSL https://raw.githubusercontent.com/farseenmanekhan1232/tkeybr/main/install.sh | bash

set -e

REPO="farseenmanekhan1232/tkeybr"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="tkeybr"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) error "Unsupported OS: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest release version
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | 
        grep '"tag_name"' | 
        sed -E 's/.*"([^"]+)".*/\1/'
}

# Main installation
main() {
    echo ""
    echo "  ╔═══════════════════════════════════════╗"
    echo "  ║        tkeybr Installer               ║"
    echo "  ║   Terminal-based typing practice      ║"
    echo "  ╚═══════════════════════════════════════╝"
    echo ""

    OS=$(detect_os)
    ARCH=$(detect_arch)
    
    info "Detected: ${OS}/${ARCH}"

    # Get latest version
    info "Fetching latest release..."
    VERSION=$(get_latest_version)
    
    if [ -z "$VERSION" ]; then
        error "Failed to get latest version"
    fi
    
    info "Latest version: ${VERSION}"

    # Construct download URL
    # GoReleaser format: tkeybr_0.4.0_darwin_amd64.tar.gz
    # Uses lowercase OS, version number, and amd64/arm64
    OS_NAME="$OS"
    ARCH_NAME="$ARCH"
    
    # Architecture naming matches GoReleaser (amd64, arm64)
    # No conversion needed for arch

    EXT="tar.gz"
    if [ "$OS" = "windows" ]; then
        EXT="zip"
    fi

    # Version without 'v' prefix for filename
    VERSION_NUM="${VERSION#v}"
    
    FILENAME="${BINARY_NAME}_${VERSION_NUM}_${OS_NAME}_${ARCH_NAME}.${EXT}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

    info "Downloading ${FILENAME}..."
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT
    
    # Download
    if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${FILENAME}"; then
        error "Download failed. URL: ${DOWNLOAD_URL}"
    fi

    # Extract
    info "Extracting..."
    cd "$TMP_DIR"
    
    if [ "$EXT" = "tar.gz" ]; then
        tar -xzf "$FILENAME"
    else
        unzip -q "$FILENAME"
    fi

    # Install
    info "Installing to ${INSTALL_DIR}..."
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        warn "Need sudo to install to ${INSTALL_DIR}"
        sudo mv "$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Verify installation
    if command -v "$BINARY_NAME" &> /dev/null; then
        echo ""
        info "✅ Installation complete!"
        echo ""
        echo "  Run 'tkeybr' to start practicing!"
        echo ""
        "${BINARY_NAME}" --version 2>/dev/null || true
    else
        warn "Installed but '${BINARY_NAME}' not found in PATH."
        warn "You may need to add ${INSTALL_DIR} to your PATH."
        echo ""
        echo "  Add this to your shell config (~/.bashrc or ~/.zshrc):"
        echo "    export PATH=\"\$PATH:${INSTALL_DIR}\""
        echo ""
    fi
}

main "$@"
