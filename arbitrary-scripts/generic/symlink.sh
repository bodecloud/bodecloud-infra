#!/bin/bash

# Check if two arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <symlink_path> <destination_path>"
    exit 1
fi

SYMLINK_PATH="$1"
DEST_PATH="$2"

# Call rm twice, first time without -r.
# Otherwise the symlink will be followed and all the target files will be deleted. (if symlink_path was already symlinked somewhere)
sudo mkdir -p "$SYMLINK_PATH"
sudo rm -rf "$SYMLINK_PATH" >> /dev/null 2>&1
sudo rm -rvf "$SYMLINK_PATH"
sudo mkdir -p "$DEST_PATH"
sudo ln -s "$DEST_PATH" "$SYMLINK_PATH"
