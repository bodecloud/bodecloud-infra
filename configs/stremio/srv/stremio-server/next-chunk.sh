#!/bin/bash

# Script to extract the next chunk of server.js for refactoring
# Usage: ./next-chunk.sh start_line [lines_count]

START_LINE=$1
LINES_COUNT=${2:-1000}

if [ -z "$START_LINE" ]; then
  echo "Usage: ./next-chunk.sh start_line [lines_count]"
  echo "Example: ./next-chunk.sh 1000 500"
  exit 1
fi

END_LINE=$((START_LINE + LINES_COUNT - 1))

echo "Extracting lines $START_LINE-$END_LINE from server.js"
sed -n "${START_LINE},${END_LINE}p" server.js > next-chunk.txt

echo "Content saved to next-chunk.txt"
echo "Next extraction should start at line $((END_LINE + 1))" 