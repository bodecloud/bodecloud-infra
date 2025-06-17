#!/bin/bash

# Script to run kubectl-neat on all YAML files recursively
# Outputs cleaned files to 'neat' subfolder with same name but .yaml extension

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🧹 Running kubectl-neat on all YAML files recursively...${NC}"

# Check if kubectl-neat is installed
if ! command -v kubectl-neat &> /dev/null; then
    echo -e "${RED}❌ kubectl-neat is not installed or not in PATH${NC}"
    echo "Install it with: kubectl krew install neat"
    exit 1
fi

# Counter for processed files
processed=0
errors=0

# Function to process a single file
process_file() {
    local file="$1"
    
    # Skip files that are already in 'neat' folders
    if [[ "$file" == */neat/* ]]; then
        echo -e "${YELLOW}⏭️  Skipping file already in neat folder: $file${NC}"
        return 0
    fi
    
    # Generate output filename
    local dir=$(dirname "$file")
    local basename=$(basename "$file")
    local filename="${basename%.*}"
    
    # Create neat subfolder if it doesn't exist
    local neat_dir="$dir/neat"
    mkdir -p "$neat_dir"
    
    # Always use .yaml extension in the neat folder
    local output_file="$neat_dir/$filename.yaml"
    
    # Run kubectl-neat
    if kubectl-neat -f "$file" > "$output_file" 2>/dev/null; then
        echo -e "${GREEN}✅ Successfully cleaned: $file${NC}"
        processed=$((processed + 1))
    else
        echo -e "${RED}❌ Failed to process: $file${NC}"
        # Remove empty output file if created
        [[ -f "$output_file" ]] && rm "$output_file"
        errors=$((errors + 1))
    fi
}

# Find all YAML files recursively, excluding hidden directories and common non-k8s directories
echo -e "${YELLOW}🔍 Finding YAML files recursively...${NC}"

# Use find to get all YAML files, excluding hidden dirs, node_modules, .git, etc.
while IFS= read -r -d '' file; do
    process_file "$file"
done < <(find . -type f \( -name "*.yaml" -o -name "*.yml" \) \
    -not -path "./.git/*" \
    -not -path "./node_modules/*" \
    -not -path "./.next/*" \
    -not -path "./dist/*" \
    -not -path "./build/*" \
    -not -path "*/neat/*" \
    -print0)

echo ""
echo -e "${GREEN}🎉 kubectl-neat processing complete!${NC}"
echo -e "${GREEN}📊 Processed: $processed files${NC}"
if [[ $errors -gt 0 ]]; then
    echo -e "${RED}⚠️  Errors: $errors files${NC}"
else
    echo -e "${GREEN}🎯 All files processed successfully!${NC}"
fi