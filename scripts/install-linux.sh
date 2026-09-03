#!/bin/bash
# install-linux.sh - Distro-agnostic installer for quotes
# Supports: systemd, OpenRC, init.d, runit
# Auto-detects: architecture, init system, package manager

set -e

PROJECTNAME="quotes"
GITHUB_REPO="apimgr/quotes"
VERSION="latest"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== ${PROJECTNAME} Installer ===${NC}"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "Architecture: $ARCH"

# Detect if running as root
if [ "$EUID" -eq 0 ]; then
    IS_ROOT=true
    BIN_DIR="/usr/local/bin"
    CONFIG_DIR="/etc/${PROJECTNAME}"
    DATA_DIR="/var/lib/${PROJECTNAME}"
    LOG_DIR="/var/log/${PROJECTNAME}"
else
    IS_ROOT=false
    BIN_DIR="$HOME/.local/bin"
    CONFIG_DIR="$HOME/.config/${PROJECTNAME}"
    DATA_DIR="$HOME/.local/share/${PROJECTNAME}"
    LOG_DIR="$HOME/.local/state/${PROJECTNAME}"
fi

echo "Install mode: $([ "$IS_ROOT" = true ] && echo "System (root)" || echo "User")"

# Create directories
mkdir -p "$BIN_DIR" "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"
mkdir -p "$DATA_DIR/db"

# Download binary
echo "Downloading ${PROJECTNAME}-linux-${ARCH}..."
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/${VERSION}/download/${PROJECTNAME}-linux-${ARCH}"

if command -v curl &> /dev/null; then
    curl -L -o "${BIN_DIR}/${PROJECTNAME}" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -O "${BIN_DIR}/${PROJECTNAME}" "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: curl or wget required${NC}"
    exit 1
fi

chmod +x "${BIN_DIR}/${PROJECTNAME}"
echo -e "${GREEN}✓ Binary installed to ${BIN_DIR}/${PROJECTNAME}${NC}"

# Detect init system
detect_init() {
    if [ -d /run/systemd/system ] || command -v systemctl &> /dev/null; then
        echo "systemd"
    elif [ -f /sbin/openrc-run ] || [ -d /etc/init.d ] && grep -q "openrc" /sbin/init 2>/dev/null; then
        echo "openrc"
    elif [ -d /etc/init.d ] && [ ! -d /run/systemd/system ]; then
        echo "sysvinit"
    elif command -v sv &> /dev/null; then
        echo "runit"
    else
        echo "unknown"
    fi
}

INIT_SYSTEM=$(detect_init)
echo "Init system: $INIT_SYSTEM"

# Install service based on init system
case $INIT_SYSTEM in
    systemd)
        if [ "$IS_ROOT" = true ]; then
            cat > /etc/systemd/system/${PROJECTNAME}.service << EOF
[Unit]
Description=Quotes API Server
After=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
ExecStart=${BIN_DIR}/${PROJECTNAME}
Restart=always
RestartSec=5
Environment="CONFIG_DIR=${CONFIG_DIR}"
Environment="DATA_DIR=${DATA_DIR}"
Environment="LOGS_DIR=${LOG_DIR}"
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
            systemctl daemon-reload
            systemctl enable ${PROJECTNAME}
            systemctl start ${PROJECTNAME}
            echo -e "${GREEN}✓ Service installed and started${NC}"
            echo "  systemctl status ${PROJECTNAME}"
            echo "  journalctl -u ${PROJECTNAME} -f"
        else
            echo -e "${YELLOW}User install: Run manually with: ${PROJECTNAME}${NC}"
        fi
        ;;
    *)
        echo -e "${YELLOW}Init system ${INIT_SYSTEM} - manual start required${NC}"
        echo "Run: ${BIN_DIR}/${PROJECTNAME}"
        ;;
esac

echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo "Admin credentials: ${CONFIG_DIR}/admin-credentials.txt"
echo "Configuration: ${CONFIG_DIR}/"
echo "Data: ${DATA_DIR}/"
echo "Logs: ${LOG_DIR}/"
