#!/bin/bash

# Script to restart Docker Compose containers one by one with cleanup
# This helps manage disk space by pruning unused resources between restarts

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.yml"
SLEEP_BETWEEN_CONTAINERS=5
DRY_RUN=false

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -f, --file FILE     Docker compose file (default: docker-compose.yml)"
    echo "  -s, --sleep SECONDS Sleep between containers (default: 5)"
    echo "  -d, --dry-run       Show what would be done without executing"
    echo "  -h, --help          Show this help message"
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--file)
            COMPOSE_FILE="$2"
            shift 2
            ;;
        -s|--sleep)
            SLEEP_BETWEEN_CONTAINERS="$2"
            shift 2
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            print_error "Unknown option: $1"
            usage
            ;;
    esac
done

# Check if docker-compose.yml exists
if [[ ! -f "$COMPOSE_FILE" ]]; then
    print_error "Docker compose file '$COMPOSE_FILE' not found!"
    exit 1
fi

# Check if docker and docker compose are available
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed or not in PATH"
    exit 1
fi

if ! docker compose version &> /dev/null; then
    print_error "Docker Compose is not available"
    exit 1
fi

print_status "Parsing $COMPOSE_FILE to extract service names..."

# Extract service names from docker-compose.yml
# This regex looks for lines that start with service names (no leading spaces, followed by colon)
# and excludes YAML anchors (starting with x-) and the services: keyword itself
services=($(awk '
/^services:/ { in_services=1; next }
/^[a-zA-Z]/ && in_services { in_services=0 }
in_services && /^  [a-zA-Z0-9_-]+:/ && !/^  x-/ {
    gsub(/:.*$/, ""); 
    gsub(/^  /, ""); 
    print
}' "$COMPOSE_FILE" | sort))

if [[ ${#services[@]} -eq 0 ]]; then
    print_error "No services found in $COMPOSE_FILE"
    exit 1
fi

print_success "Found ${#services[@]} services:"
for service in "${services[@]}"; do
    echo "  - $service"
done

if [[ "$DRY_RUN" == "true" ]]; then
    print_warning "DRY RUN MODE - No actual changes will be made"
    echo
    print_status "Would execute the following sequence for each service:"
    print_status "  1. Check if service is already running (skip if yes)"
    print_status "  2. docker compose -f $COMPOSE_FILE up -d <service>"
    print_status "  3. docker system prune -af"
    print_status "  4. Sleep for $SLEEP_BETWEEN_CONTAINERS seconds"
    echo
    print_status "Services that are already running (due to dependencies) will be skipped"
    print_status "This prevents unnecessary restarts and reduces total execution time"
    exit 0
fi

# Confirm before proceeding
echo
read -p "Do you want to proceed with restarting all containers? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warning "Operation cancelled by user"
    exit 0
fi

# Get initial disk usage
print_status "Checking initial disk usage..."
initial_usage=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
print_status "Initial disk usage: ${initial_usage}%"

echo
print_status "Starting container restart sequence..."
echo "========================================"

# Counter for progress
total_services=${#services[@]}
current_service=0
processed_services=()

# Function to check if a service is running
is_service_running() {
    local service_name="$1"
    # Try using jq first (more reliable)
    if command -v jq &> /dev/null; then
        docker compose -f "$COMPOSE_FILE" ps "$service_name" --format json 2>/dev/null | jq -r '.State' 2>/dev/null | grep -q "running"
    else
        # Fallback method without jq
        docker compose -f "$COMPOSE_FILE" ps "$service_name" 2>/dev/null | grep -q "Up"
    fi
}

# Restart each service individually
for service in "${services[@]}"; do
    current_service=$((current_service + 1))
    
    echo
    print_status "[$current_service/$total_services] Processing service: $service"
    
    # Check if the service is already running
    if is_service_running "$service"; then
        print_warning "Service $service is already running (likely started as dependency), skipping..."
        processed_services+=("$service")
        continue
    fi
    
    # Start the service
    print_status "Starting $service..."
    if docker compose -f "$COMPOSE_FILE" up -d "$service"; then
        print_success "Successfully started $service"
        processed_services+=("$service")
    else
        print_error "Failed to start $service"
        continue
    fi
    
    # Wait a moment for the container to stabilize
    sleep 2
    
    # Run docker system prune
    print_status "Running docker system prune..."
    if docker system prune -af; then
        print_success "Docker system prune completed"
    else
        print_warning "Docker system prune had issues"
    fi
    
    # Check current disk usage
    current_usage=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
    print_status "Current disk usage: ${current_usage}%"
    
    # Sleep between containers (except for the last one)
    if [[ $current_service -lt $total_services ]]; then
        print_status "Waiting $SLEEP_BETWEEN_CONTAINERS seconds before next container..."
        sleep "$SLEEP_BETWEEN_CONTAINERS"
    fi
done

# Final status
echo
echo "========================================"
print_success "All containers have been processed!"

# Calculate statistics
total_processed=${#processed_services[@]}
total_skipped=$((total_services - total_processed))

print_status "Processing summary:"
print_status "  Total services: $total_services"
print_status "  Services processed: $total_processed"
print_status "  Services skipped: $total_skipped"

if [[ $total_skipped -gt 0 ]]; then
    print_status "Skipped services were already running (likely started as dependencies)"
fi

# Get final disk usage
final_usage=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
disk_saved=$((initial_usage - final_usage))

print_status "Final disk usage: ${final_usage}%"
if [[ $disk_saved -gt 0 ]]; then
    print_success "Disk space saved: ${disk_saved}%"
elif [[ $disk_saved -lt 0 ]]; then
    print_warning "Disk usage increased by: $((disk_saved * -1))%"
else
    print_status "No change in disk usage"
fi

# Show final container status
echo
print_status "Final container status:"
docker compose -f "$COMPOSE_FILE" ps

print_success "Script completed successfully!" 