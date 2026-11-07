#!/bin/sh
set -eu

export TZ="${TZ:-UTC}"
export CONFIG_DIR="${CONFIG_DIR:-/config}"
export DATA_DIR="${DATA_DIR:-/data}"
export LOGS_DIR="${LOGS_DIR:-/logs}"

exec /usr/local/bin/quotes --address "${ADDRESS:-0.0.0.0}" --port "${PORT:-80}" "$@"
