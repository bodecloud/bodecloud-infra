#!/bin/bash

# This script finds all timestamped files in the cdk8s-project directory,
# keeps only the most recent version of each file, and renames it to remove the timestamp.

# Set the directory to process
DIR="cdk8s-project"

# Check if dry run mode is enabled
DRY_RUN=false
if [ "$1" = "--dry-run" ]; then
  DRY_RUN=true
  echo "Running in dry run mode - no changes will be made"
fi

# Function to process files
process_files() {
  local file=$1
  local dry_run=$2

  # Extract the directory, base name without timestamp, and extension
  dir=$(dirname "$file")
  filename=$(basename "$file")
  base=$(echo "$filename" | sed -E 's/(.+)_[0-9]+\.(.*)/\1/')
  ext=$(echo "$filename" | sed -E 's/.*\.([^.]+)$/\1/')

  # Create a temporary file to store the list of files with this base name
  tmp_file=$(mktemp)

  # Find all files with the same base name in the same directory
  find "$dir" -type f -name "${base}_[0-9]*.$ext" | sort >"$tmp_file"

  # Get the most recent file (last in sorted list)
  most_recent=$(tail -n 1 "$tmp_file")

  # If this is the most recent file, keep it and rename
  if [ "$file" = "$most_recent" ]; then
    new_name="${dir}/${base}.${ext}"
    echo "Keeping $file as the most recent version"
    echo "Renaming to $new_name"

    # Delete any existing non-timestamped file
    if [ -f "$new_name" ]; then
      echo "Removing existing file $new_name"
      if [ "$dry_run" = false ]; then
        rm "$new_name"
      fi
    fi

    # Rename the most recent file
    if [ "$dry_run" = false ]; then
      mv "$file" "$new_name"
    fi
  else
    # Delete older versions
    echo "Removing older version: $file"
    if [ "$dry_run" = false ]; then
      rm "$file"
    fi
  fi

  # Clean up temp file
  rm "$tmp_file"
}

# Create a list of processed base filenames to avoid duplicate processing
processed_files=$(mktemp)

# Find all timestamped files
echo "Finding all timestamped files..."
find "$DIR" -type f -name "*_[0-9]*.*" | while read -r file; do
  # Extract the directory, base name without timestamp, and extension
  dir=$(dirname "$file")
  filename=$(basename "$file")
  base=$(echo "$filename" | sed -E 's/(.+)_[0-9]+\.(.*)/\1/')
  ext=$(echo "$filename" | sed -E 's/.*\.([^.]+)$/\1/')

  # Check if this base filename has been processed already
  if ! grep -q "^${dir}/${base}.${ext}$" "$processed_files"; then
    # Process all files with this base name
    echo "${dir}/${base}.${ext}" >>"$processed_files"

    # Find all files with the same base name in the same directory
    find "$dir" -type f -name "${base}_[0-9]*.$ext" | sort | while read -r similar_file; do
      process_files "$similar_file" "$DRY_RUN"
    done
  fi
done

# Clean up temp file
rm "$processed_files"

echo "Cleanup completed!"
>