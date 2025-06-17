#!/bin/bash

# RClone Zurg Mount Service Management Script
# This script provides easy management of the rclone mount service

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SERVICE_NAME="rclone-zurg-mount"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[RCLONE SERVICE]${NC} $1"
}

# Function to check if service exists
check_service_exists() {
    if ! systemctl list-unit-files | grep -q "$SERVICE_NAME.service"; then
        print_error "Service $SERVICE_NAME is not installed."
        print_status "Run './install-rclone-service.sh' to install it first."
        exit 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 {start|stop|restart|status|logs|enable|disable|install}"
    echo ""
    echo "Commands:"
    echo "  start    - Start the rclone mount service"
    echo "  stop     - Stop the rclone mount service"
    echo "  restart  - Restart the rclone mount service"
    echo "  status   - Show service status"
    echo "  logs     - Show service logs (follow mode)"
    echo "  enable   - Enable service to start on boot"
    echo "  disable  - Disable service from starting on boot"
    echo "  install  - Install the service (requires sudo)"
    echo ""
}

# Check if running as root for certain operations
check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This operation requires root privileges. Please run with sudo."
        exit 1
    fi
}

case "$1" in
    start)
        check_service_exists
        check_root
        print_header "Starting $SERVICE_NAME service..."
        systemctl start "$SERVICE_NAME"
        print_status "Service started successfully!"
        ;;
    stop)
        check_service_exists
        check_root
        print_header "Stopping $SERVICE_NAME service..."
        systemctl stop "$SERVICE_NAME"
        print_status "Service stopped successfully!"
        ;;
    restart)
        check_service_exists
        check_root
        print_header "Restarting $SERVICE_NAME service..."
        systemctl restart "$SERVICE_NAME"
        print_status "Service restarted successfully!"
        ;;
    status)
        check_service_exists
        print_header "Service status:"
        systemctl status "$SERVICE_NAME" --no-pager
        ;;
    logs)
        check_service_exists
        print_header "Following service logs (Ctrl+C to exit):"
        journalctl -u "$SERVICE_NAME" -f
        ;;
    enable)
        check_service_exists
        check_root
        print_header "Enabling $SERVICE_NAME service..."
        systemctl enable "$SERVICE_NAME"
        print_status "Service enabled to start on boot!"
        ;;
    disable)
        check_service_exists
        check_root
        print_header "Disabling $SERVICE_NAME service..."
        systemctl disable "$SERVICE_NAME"
        print_status "Service disabled from starting on boot!"
        ;;
    install)
        check_root
        print_header "Installing $SERVICE_NAME service..."
        ./install-rclone-service.sh
        ;;
    *)
        show_usage
        exit 1
        ;;
esac 