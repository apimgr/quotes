#!/bin/sh
# install-bsd.sh - BSD installer with rc.d service
# Supports: FreeBSD, OpenBSD, NetBSD

set -e

PROJECTNAME="quotes"
GITHUB_REPO="apimgr/quotes"
VERSION="latest"

echo "=== ${PROJECTNAME} Installer for BSD ==="

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Architecture: $ARCH"

# Directories
BIN_DIR="/usr/local/bin"
CONFIG_DIR="/usr/local/etc/${PROJECTNAME}"
DATA_DIR="/var/db/${PROJECTNAME}"
LOG_DIR="/var/log/${PROJECTNAME}"
RC_DIR="/usr/local/etc/rc.d"

# Create directories
mkdir -p "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR" "$DATA_DIR/db"

# Download binary
echo "Downloading ${PROJECTNAME}-bsd-${ARCH}..."
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/${VERSION}/download/${PROJECTNAME}-bsd-${ARCH}"

# Try fetch first (BSD native), fallback to curl
if command -v fetch >/dev/null 2>&1; then
    fetch -o "${BIN_DIR}/${PROJECTNAME}" "$DOWNLOAD_URL"
else
    curl -L -o "${BIN_DIR}/${PROJECTNAME}" "$DOWNLOAD_URL"
fi

chmod +x "${BIN_DIR}/${PROJECTNAME}"
echo "✓ Binary installed to ${BIN_DIR}/${PROJECTNAME}"

# Create rc.d script
cat > "${RC_DIR}/${PROJECTNAME}" << 'RCEOF'
#!/bin/sh
#
# PROVIDE: quotes
# REQUIRE: DAEMON
# KEYWORD: shutdown

. /etc/rc.subr

name="quotes"
rcvar="quotes_enable"
procname="/usr/local/bin/quotes"
command="/usr/local/bin/quotes"
pidfile="/var/run/${name}.pid"
quotes_config="/usr/local/etc/quotes"
quotes_data="/var/db/quotes"
quotes_logs="/var/log/quotes"

command_args=""
quotes_env="CONFIG_DIR=${quotes_config} DATA_DIR=${quotes_data} LOGS_DIR=${quotes_logs}"

load_rc_config $name
: ${quotes_enable:=NO}

run_rc_command "$1"
RCEOF

chmod +x "${RC_DIR}/${PROJECTNAME}"
echo "✓ RC script created"

# Enable and start service
sysrc quotes_enable="YES"
service ${PROJECTNAME} start

echo "✓ Installation complete!"
echo ""
echo "Admin credentials: ${CONFIG_DIR}/admin-credentials.txt"
echo "Service management:"
echo "  service ${PROJECTNAME} status"
echo "  service ${PROJECTNAME} restart"
echo "  service ${PROJECTNAME} stop"
