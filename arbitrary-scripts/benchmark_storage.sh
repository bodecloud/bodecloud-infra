#!/bin/bash

# benchmark_storage.sh
# A comprehensive storage benchmark script comparing rclone mounts with local storage
# Usage: ./benchmark_storage.sh <rclone_mount_path> <local_path>

# Set strict bash options
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print section headers
print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}\n"
}

# Function to print errors
print_error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
}

# Function to print success messages
print_success() {
    echo -e "${GREEN}$1${NC}"
}

# Check if required commands are installed
check_dependencies() {
    local missing_deps=()
    
    for cmd in dd fio rsync iozone; do
        if ! command -v "$cmd" &> /dev/null; then
            missing_deps+=("$cmd")
        fi
    done
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing dependencies: ${missing_deps[*]}"
        echo "Please install missing dependencies:"
        echo "sudo apt install fio rsync iozone3  # for Ubuntu/Debian"
        echo "sudo yum install fio rsync iozone   # for CentOS/RHEL"
        exit 1
    fi
}

# Function to run dd tests
run_dd_tests() {
    local path=$1
    local filename="$path/dd_testfile"
    
    print_header "Running DD Tests on $path"
    
    echo "Write speed test (1GB):"
    dd if=/dev/zero of="$filename" bs=1M count=1024 status=progress
    sync
    echo
    
    echo "Read speed test (1GB):"
    dd if="$filename" of=/dev/null bs=1M status=progress
    echo
    
    rm -f "$filename"
}

# Function to run fio tests
run_fio_tests() {
    local path=$1
    local filename="$path/fio_testfile"
    
    print_header "Running FIO Tests on $path"
    
    fio --name=test --filename="$filename" \
        --rw=readwrite --bs=4k --size=100M --direct=1 \
        --numjobs=1 --runtime=60 --group_reporting
    
    rm -f "$filename"
}

# Function to run rsync tests
run_rsync_tests() {
    local path=$1
    local source_dir="$path/source_test_dir"
    local dest_dir="$path/dest_test_dir"
    
    print_header "Running Rsync Tests on $path"
    
    # Create test directory with some files
    mkdir -p "$source_dir"
    for i in {1..10}; do
        dd if=/dev/urandom of="$source_dir/file$i" bs=1M count=10 2>/dev/null
    done
    
    echo "Copying 100MB of random files:"
    time rsync -av "$source_dir/" "$dest_dir/"
    
    # Cleanup
    rm -rf "$source_dir" "$dest_dir"
}

# Function to run iozone tests
run_iozone_tests() {
    local path=$1
    local filename="$path/iozone_testfile"
    
    print_header "Running IOzone Tests on $path"
    
    iozone -a -n 512M -g 1G -i 0 -i 1 -f "$filename"
    
    rm -f "$filename"
}

# Function to run small file operations test
run_small_files_test() {
    local path=$1
    local test_dir="$path/small_files_test"
    
    print_header "Running Small Files Test on $path"
    
    mkdir -p "$test_dir"
    
    echo "Creating 1000 small files:"
    time for i in {1..1000}; do
        touch "$test_dir/test_$i"
    done
    
    echo -e "\nDeleting 1000 small files:"
    time for i in {1..1000}; do
        rm "$test_dir/test_$i"
    done
    
    rmdir "$test_dir"
}

# Main function
main() {
    # Check arguments
    if [ $# -ne 2 ]; then
        print_error "Usage: $0 <rclone_mount_path> <local_path>"
        exit 1
    fi
    
    local rclone_path="$1"
    local local_path="$2"
    
    # Check if paths exist
    for path in "$rclone_path" "$local_path"; do
        if [ ! -d "$path" ]; then
            print_error "Directory does not exist: $path"
            exit 1
        fi
        if [ ! -w "$path" ]; then
            print_error "Directory is not writable: $path"
            exit 1
        fi
    done
    
    # Check dependencies
    check_dependencies
    
    # Start benchmarking
    print_header "Starting Benchmark Suite"
    echo "Rclone Mount: $rclone_path"
    echo "Local Path: $local_path"
    
    # Run tests for both paths
    for path in "$rclone_path" "$local_path"; do
        print_header "Testing $path"
        
        run_dd_tests "$path"
        run_fio_tests "$path"
        run_rsync_tests "$path"
        run_iozone_tests "$path"
        run_small_files_test "$path"
    done
    
    print_success "Benchmark completed successfully!"
}

# Run main function with all arguments
main "$@" 