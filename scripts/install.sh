#!/usr/bin/env bash
set -euo pipefail

# Wand installer script
# Usage: curl -sSL https://raw.githubusercontent.com/ochairo/wand/main/scripts/install.sh | bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="ochairo/wand"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ${NC} $*"
}

log_success() {
    echo -e "${GREEN}✓${NC} $*"
}

log_error() {
    echo -e "${RED}✗${NC} $*" >&2
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $*"
}

# Check requirements
check_requirements() {
    local missing=()

    for cmd in curl tar grep sed; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing+=("$cmd")
        fi
    done

    if [ ${#missing[@]} -gt 0 ]; then
        log_error "Missing required commands: ${missing[*]}"
        exit 1
    fi
}

# Detect OS and architecture
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Darwin) os="darwin" ;;
        Linux) os="linux" ;;
        *)
            log_error "Unsupported OS: $(uname -s)"
            log_info "Supported: Darwin (macOS), Linux"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            log_info "Supported: x86_64 (amd64), aarch64 (arm64)"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local version

    version=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

    if [ -z "$version" ]; then
        log_error "Failed to fetch latest version from GitHub"
        log_info "Check your internet connection or try again later"
        exit 1
    fi

    echo "$version"
}

# Verify checksum
verify_checksum() {
    local file="$1"
    local checksum_url="$2"
    local checksum_file="$file.sha256"

    log_info "Verifying checksum..."

    if ! curl -fsSL "$checksum_url" -o "$checksum_file"; then
        log_warn "Checksum file not available, skipping verification"
        return 0
    fi

    if command -v shasum >/dev/null 2>&1; then
        if shasum -a 256 -c "$checksum_file" >/dev/null 2>&1; then
            log_success "Checksum verified"
            return 0
        fi
    elif command -v sha256sum >/dev/null 2>&1; then
        if sha256sum -c "$checksum_file" >/dev/null 2>&1; then
            log_success "Checksum verified"
            return 0
        fi
    fi

    log_error "Checksum verification failed"
    exit 1
}

# Download and install
install_wand() {
    local platform version download_url checksum_url tmp_dir

    platform=$(detect_platform)
    log_info "Installing Wand for ${platform}..."

    version=$(get_latest_version)
    log_success "Latest version: ${version}"

    download_url="https://github.com/${REPO}/releases/download/${version}/wand-${platform}.tar.gz"
    checksum_url="${download_url}.sha256"

    # Create temporary directory
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT

    # Download
    log_info "Downloading from ${download_url}..."
    if ! curl -fsSL "$download_url" -o "$tmp_dir/wand.tar.gz"; then
        log_error "Failed to download wand"
        log_info "URL: ${download_url}"
        exit 1
    fi

    # Verify checksum
    (cd "$tmp_dir" && verify_checksum "wand.tar.gz" "$checksum_url")

    # Extract
    log_info "Extracting..."
    if ! tar xzf "$tmp_dir/wand.tar.gz" -C "$tmp_dir"; then
        log_error "Failed to extract archive"
        exit 1
    fi

    # Verify binary exists
    if [ ! -f "$tmp_dir/wand" ]; then
        log_error "Binary not found in archive"
        exit 1
    fi

    # Install
    log_info "Installing to ${INSTALL_DIR}..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_dir/wand" "$INSTALL_DIR/wand"
        chmod +x "$INSTALL_DIR/wand"
    else
        if ! command -v sudo >/dev/null 2>&1; then
            log_error "sudo is required but not available"
            log_info "Try setting INSTALL_DIR to a writable location:"
            log_info "  INSTALL_DIR=\$HOME/.local/bin curl -sSL ... | bash"
            exit 1
        fi
        sudo mv "$tmp_dir/wand" "$INSTALL_DIR/wand"
        sudo chmod +x "$INSTALL_DIR/wand"
    fi

    # Verify installation
    if command -v wand >/dev/null 2>&1; then
        local installed_version
        installed_version=$(wand version 2>/dev/null | head -n1 || echo "unknown")
        log_success "Successfully installed wand ${installed_version}"
    else
        log_success "Successfully installed wand to ${INSTALL_DIR}/wand"
        log_warn "${INSTALL_DIR} might not be in your PATH"
        log_info "Add it to your PATH with:"
        log_info "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

# Main
main() {
    echo -e "${GREEN}"
    cat << 'EOF'
╔══════════════════════════════════════╗
║         🪄 Wand Installer            ║
╚══════════════════════════════════════╝
EOF
    echo -e "${NC}"

    check_requirements
    install_wand

    echo ""
    log_info "Get started with:"
    echo "  wand --help"
    echo "  wand install kubectl"
    echo ""
}

main "$@"
