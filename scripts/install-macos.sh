#!/bin/bash
# install-macos.sh - macOS installer with launchd service
# Supports: Intel (amd64) and Apple Silicon (arm64)

set -e

PROJECTNAME="quotes"
PROJECT_DISPLAY="Quotes"
GITHUB_REPO="apimgr/quotes"
VERSION="latest"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== ${PROJECT_DISPLAY} Installer for macOS ===${NC}"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)  ARCH="amd64" ;;
    arm64)   ARCH="arm64" ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "Architecture: $ARCH ($(uname -m))"

# Detect if running with sudo
if [ "$EUID" -eq 0 ]; then
    IS_ROOT=true
    BIN_DIR="/usr/local/bin"
    CONFIG_DIR="/Library/Application Support/${PROJECT_DISPLAY}"
    DATA_DIR="/Library/Application Support/${PROJECT_DISPLAY}/data"
    LOG_DIR="/Library/Logs/${PROJECT_DISPLAY}"
    LAUNCHD_DIR="/Library/LaunchDaemons"
else
    IS_ROOT=false
    BIN_DIR="$HOME/.local/bin"
    CONFIG_DIR="$HOME/Library/Application Support/${PROJECT_DISPLAY}"
    DATA_DIR="$HOME/Library/Application Support/${PROJECT_DISPLAY}/data"
    LOG_DIR="$HOME/Library/Logs/${PROJECT_DISPLAY}"
    LAUNCHD_DIR="$HOME/Library/LaunchAgents"
fi

echo "Install mode: $([ "$IS_ROOT" = true ] && echo "System (sudo)" || echo "User")"

# Create directories
mkdir -p "$BIN_DIR" "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR" "$LAUNCHD_DIR"
mkdir -p "$DATA_DIR/db"

# Download binary
echo "Downloading ${PROJECTNAME}-macos-${ARCH}..."
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/${VERSION}/download/${PROJECTNAME}-macos-${ARCH}"

curl -L -o "${BIN_DIR}/${PROJECTNAME}" "$DOWNLOAD_URL"
chmod +x "${BIN_DIR}/${PROJECTNAME}"
echo -e "${GREEN}✓ Binary installed to ${BIN_DIR}/${PROJECTNAME}${NC}"

# Create launchd plist
PLIST_FILE="${LAUNCHD_DIR}/com.${GITHUB_REPO//\//.}.${PROJECTNAME}.plist"

cat > "$PLIST_FILE" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.apimgr.quotes</string>
    <key>ProgramArguments</key>
    <array>
        <string>${BIN_DIR}/${PROJECTNAME}</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>CONFIG_DIR</key>
        <string>${CONFIG_DIR}</string>
        <key>DATA_DIR</key>
        <string>${DATA_DIR}</string>
        <key>LOGS_DIR</key>
        <string>${LOG_DIR}</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${LOG_DIR}/stdout.log</string>
    <key>StandardErrorPath</key>
    <string>${LOG_DIR}/stderr.log</string>
</dict>
</plist>
EOF

# Load service
if [ "$IS_ROOT" = true ]; then
    launchctl load "$PLIST_FILE"
    echo -e "${GREEN}✓ Service installed and started${NC}"
    echo "  launchctl list | grep ${PROJECTNAME}"
else
    launchctl load "$PLIST_FILE"
    echo -e "${GREEN}✓ User service installed and started${NC}"
fi

echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo "Admin credentials: ${CONFIG_DIR}/admin-credentials.txt"
echo "Configuration: ${CONFIG_DIR}/"
echo "Data: ${DATA_DIR}/"
echo "Logs: ${LOG_DIR}/"
echo ""
echo "Service management:"
echo "  launchctl list | grep ${PROJECTNAME}"
echo "  launchctl stop com.apimgr.quotes"
echo "  launchctl start com.apimgr.quotes"
