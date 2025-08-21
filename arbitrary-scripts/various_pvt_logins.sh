#!/bin/bash

# Exit on any error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper function for logging
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_exec() {
    echo -e "${CYAN}[EXEC]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Check if script is run as root
if [ "$EUID" -ne 0 ]; then
    echo ""
    error "Please run as root (use \`sudo -E $0\`)"
    exit 1
fi

echo ""
log_exec "apt update"
apt update

echo ""
log_exec "apt upgrade -y"
apt upgrade -y

echo ""
log_exec "apt install gh -y"
apt install gh -y

echo ""
log_exec "gh auth login --hostname github.com --git-protocol https --with-token <<< *********"
gh auth login \
  --hostname github.com \
  --git-protocol https \
  --with-token <<< "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"

echo ""
log_exec "docker login ghcr.io -u th3w1zard1 -p *******************"
docker login ghcr.io -u th3w1zard1 -p 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

echo ""
log_exec "docker login docker.io -u th3w1zard1 -p *********"
docker login docker.io -u th3w1zard1 -p h4L0m4St3R327