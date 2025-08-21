#!/bin/bash

# Script to run helmify on YAML files individually
# Creates separate Helm charts for each valid manifest file

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}⚙️  Running helmify on YAML files individually...${NC}"

# Check if helmify is installed
if ! command -v helmify &> /dev/null; then
    echo -e "${RED}❌ helmify is not installed or not in PATH${NC}"
    echo "Install it from: https://github.com/Arttor/helmify"
    exit 1
fi

# Counters
processed=0
errors=0
skipped=0

# Create output directory for charts
charts_dir="./helmify-charts"
mkdir -p "$charts_dir"

# Create logs directory
logs_dir="./helmify-logs"
mkdir -p "$logs_dir"

# Function to process a single file
process_file() {
    local file="$1"
    local basename=$(basename "$file" .yaml)
    local chart_name="chart-$basename"
    local log_file="$logs_dir/$basename.log"
    
    echo -e "${BLUE}🔍 Processing: $file${NC}"
    
    # Check if file is empty or contains only '{}'
    if [[ ! -s "$file" ]] || [[ "$(cat "$file" | tr -d '[:space:]')" == "{}" ]]; then
        echo -e "${YELLOW}⏭️  Skipping empty file: $file${NC}"
        echo "SKIPPED: Empty file or contains only '{}'" > "$log_file"
        skipped=$((skipped + 1))
        return 0
    fi
    
    # Check if file contains 'kind:' or 'Kind:'
    if ! grep -q -i "kind:" "$file"; then
        echo -e "${YELLOW}⏭️  Skipping file without 'kind' field: $file${NC}"
        echo "SKIPPED: No 'kind' field found" > "$log_file"
        skipped=$((skipped + 1))
        return 0
    fi
    
    # Create chart directory
    local chart_dir="$charts_dir/$chart_name"
    
    echo -e "${GREEN}🔧 Creating chart: $chart_name${NC}"
    
    # Run helmify and capture output
    if cat "$file" | helmify "$chart_name" > "$log_file" 2>&1; then
        # Move the chart to our charts directory
        if [[ -d "$chart_name" ]]; then
            mv "$chart_name" "$chart_dir"
            echo -e "${GREEN}✅ Successfully created chart: $chart_name${NC}"
            echo "SUCCESS: Chart created at $chart_dir" >> "$log_file"
            processed=$((processed + 1))
        else
            echo -e "${YELLOW}⚠️  Chart directory not created for: $file${NC}"
            echo "WARNING: Chart directory not found after helmify" >> "$log_file"
            errors=$((errors + 1))
        fi
    else
        echo -e "${RED}❌ Failed to process: $file${NC}"
        echo "ERROR: Helmify failed - see log for details" >> "$log_file"
        errors=$((errors + 1))
    fi
}

# Function to process directory
process_directory() {
    local target_dir="$1"
    
    if [[ ! -d "$target_dir" ]]; then
        echo -e "${RED}❌ Directory not found: $target_dir${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}🔍 Finding YAML files in: $target_dir${NC}"
    
    # Find all YAML files in the directory
    while IFS= read -r -d '' file; do
        process_file "$file"
    done < <(find "$target_dir" -type f \( -name "*.yaml" -o -name "*.yml" \) -print0)
}

# Function to create combined chart
create_combined_chart() {
    local target_dir="$1"
    local combined_chart_name="elfhosted-combined"
    local combined_log="$logs_dir/combined.log"
    
    echo -e "${BLUE}🔗 Creating combined chart from all valid files...${NC}"
    
    # Create temporary directory with only valid files
    local temp_dir=$(mktemp -d)
    local valid_files=0
    
    while IFS= read -r -d '' file; do
        # Check if file is valid
        if [[ -s "$file" ]] && [[ "$(cat "$file" | tr -d '[:space:]')" != "{}" ]] && grep -q -i "kind:" "$file"; then
            cp "$file" "$temp_dir/"
            valid_files=$((valid_files + 1))
        fi
    done < <(find "$target_dir" -type f \( -name "*.yaml" -o -name "*.yml" \) -print0)
    
    if [[ $valid_files -gt 0 ]]; then
        echo -e "${GREEN}📦 Found $valid_files valid files for combined chart${NC}"
        
        if helmify -f "$temp_dir" -r "$combined_chart_name" > "$combined_log" 2>&1; then
            if [[ -d "$combined_chart_name" ]]; then
                mv "$combined_chart_name" "$charts_dir/"
                echo -e "${GREEN}✅ Successfully created combined chart: $combined_chart_name${NC}"
            else
                echo -e "${YELLOW}⚠️  Combined chart directory not created${NC}"
            fi
        else
            echo -e "${RED}❌ Failed to create combined chart${NC}"
            echo "Check $combined_log for details"
        fi
    else
        echo -e "${YELLOW}⚠️  No valid files found for combined chart${NC}"
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
}

# Main execution
target_directory="${1:-./neat}"

echo -e "${GREEN}🎯 Target directory: $target_directory${NC}"
echo -e "${GREEN}📁 Charts output: $charts_dir${NC}"
echo -e "${GREEN}📝 Logs output: $logs_dir${NC}"
echo ""

# Process individual files
process_directory "$target_directory"

echo ""
echo -e "${BLUE}🔗 Attempting to create combined chart...${NC}"
create_combined_chart "$target_directory"

echo ""
echo -e "${GREEN}🎉 Helmify processing complete!${NC}"
echo -e "${GREEN}📊 Processed: $processed files${NC}"
echo -e "${YELLOW}⏭️  Skipped: $skipped files${NC}"
if [[ $errors -gt 0 ]]; then
    echo -e "${RED}⚠️  Errors: $errors files${NC}"
    echo -e "${YELLOW}💡 Check logs in $logs_dir for details${NC}"
else
    echo -e "${GREEN}🎯 All valid files processed successfully!${NC}"
fi

echo ""
echo -e "${BLUE}📋 Summary of created charts:${NC}"
if [[ -d "$charts_dir" ]]; then
    ls -la "$charts_dir"
else
    echo -e "${YELLOW}No charts were created${NC}"
fi 