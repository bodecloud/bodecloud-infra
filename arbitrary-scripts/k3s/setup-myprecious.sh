#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CHART_REPO="oci://ghcr.io/elfhosted/charts"
TARGET_DIR="./helm-charts/myprecious"
TEMP_DIR="./tmp-charts"

# Function to display help message
show_help() {
    echo -e "${BLUE}Elfhosted Myprecious Chart Setup Script${NC}"
    echo "This script sets up access to elfhosted's myprecious chart repository"
    echo "and downloads chart information for reference."
    echo
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  -t, --target-dir DIR    Target directory for chart info (default: ./helm-charts/myprecious)"
    echo "  -r, --repo REPO         Chart repository URL (default: oci://ghcr.io/elfhosted/charts)"
    echo "  -h, --help              Show this help message"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--target-dir)
            TARGET_DIR="$2"
            shift 2
            ;;
        -r|--repo)
            CHART_REPO="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Unknown option $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Check if helm is available
if ! command -v helm &> /dev/null; then
    echo -e "${RED}Error: helm is not installed${NC}"
    exit 1
fi

# Create target directory
mkdir -p "$TARGET_DIR"
mkdir -p "$TEMP_DIR"

# Enable Helm OCI support
export HELM_EXPERIMENTAL_OCI=1

echo -e "${BLUE}Fetching chart information from ${CHART_REPO}...${NC}"

# Pull the chart to examine available services and configuration options
helm pull "$CHART_REPO/myprecious" --untar --untardir="$TEMP_DIR"

# Check if the chart was successfully pulled
if [ ! -d "$TEMP_DIR/myprecious" ]; then
    echo -e "${RED}Error: Failed to pull the chart from repository${NC}"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Copy chart files to target directory for reference
cp -r "$TEMP_DIR/myprecious"/* "$TARGET_DIR/"

# Create a list of available services from the chart's values.yaml
echo -e "${BLUE}Analyzing chart to extract available services...${NC}"
SERVICES_FILE="$TARGET_DIR/available-services.txt"
VALUES_FILE="$TEMP_DIR/myprecious/values.yaml"

# Parse values.yaml to get list of services
if [ -f "$VALUES_FILE" ]; then
    # Extract top-level keys that are likely services (exclude known non-service sections)
    grep -E "^[a-z0-9-]+:" "$VALUES_FILE" | 
    grep -v "^global:" | 
    grep -v "^serviceAccount:" | 
    grep -v "^persistence:" | 
    grep -v "^resources:" | 
    grep -v "^nodeSelector:" | 
    grep -v "^tolerations:" | 
    cut -d ':' -f 1 > "$SERVICES_FILE"
    
    echo -e "${GREEN}Found $(wc -l < "$SERVICES_FILE") services in the chart${NC}"
    echo -e "Services list saved to ${SERVICES_FILE}"
else
    echo -e "${RED}Warning: Could not find values.yaml in the chart${NC}"
fi

# Create a template for values.yaml
echo -e "${BLUE}Creating template values.yaml for myprecious chart...${NC}"
cat > "$TARGET_DIR/values-template.yaml" << EOF
# Global settings
global:
  namespace: media-stack
  domain: example.com
  storageClass: longhorn
  uid: 1000
  gid: 1000
  umask: "002"

# Available services to enable
# Add 'enabled: true' under each service you want to enable
# For example:
# plex:
#   enabled: true
#   ingress:
#     hosts:
#       - host: plex.example.com
#         paths:
#           - path: /
#             pathType: Prefix
#     tls:
#       - secretName: plex-tls
#         hosts:
#           - plex.example.com

# List of commonly used services:
EOF

# Add common services to the template
if [ -f "$SERVICES_FILE" ]; then
    COMMON_SERVICES=("plex" "wizarr" "riven" "zurg" "filebrowser" "homer" "rclonefm" "rcloneui" "sonarr" "radarr" "jellyfin" "jackett")
    
    for service in "${COMMON_SERVICES[@]}"; do
        if grep -q "^$service$" "$SERVICES_FILE"; then
            cat >> "$TARGET_DIR/values-template.yaml" << EOF
${service}:
  enabled: false
  # Add your custom configuration here
EOF
        fi
    done
    
    # Add comment about other available services
    cat >> "$TARGET_DIR/values-template.yaml" << EOF

# Other available services (uncomment and set enabled: true to use):
EOF
    
    # Add other services to template as comments
    for service in $(cat "$SERVICES_FILE"); do
        if ! echo "${COMMON_SERVICES[@]}" | grep -q -w "$service"; then
            cat >> "$TARGET_DIR/values-template.yaml" << EOF
# ${service}:
#   enabled: false
EOF
        fi
    done
fi

# Create documentation
echo -e "${BLUE}Creating documentation for using the chart...${NC}"
cat > "$TARGET_DIR/README.md" << EOF
# Elfhosted Myprecious Chart

This directory contains reference information for the elfhosted/myprecious Helm chart.

## Available Services

The \`available-services.txt\` file contains a list of all services that can be enabled in the chart.

## Using the Chart

1. Create a values file based on the \`values-template.yaml\` template:

   \`\`\`bash
   cp helm-charts/myprecious/values-template.yaml helm-values/myprecious-values.yaml
   \`\`\`

2. Edit the values file to enable and configure desired services.

3. Deploy using the \`deploy-elfhosted-charts.sh\` script:

   \`\`\`bash
   ./scripts/deploy-elfhosted-charts.sh --domain yourdomain.com
   \`\`\`

## Adding Custom Configuration

For each service, you can add custom configuration under its section. For example:

\`\`\`yaml
plex:
  enabled: true
  env:
    TZ: "America/New_York"
    PLEX_CLAIM: "your-claim-token"
  persistence:
    config:
      enabled: true
      existingClaim: plex-config
    media:
      enabled: true
      existingClaim: media
\`\`\`

## Updating the Chart

To update to the latest version of the chart:

\`\`\`bash
./scripts/deploy-elfhosted-charts.sh --domain yourdomain.com
\`\`\`

This will automatically use the latest version of the chart.
EOF

# Clean up temporary directory
rm -rf "$TEMP_DIR"

echo -e "${GREEN}Setup completed successfully!${NC}"
echo "- Chart reference files are in: $TARGET_DIR"
echo "- Template values file: $TARGET_DIR/values-template.yaml"
echo "- Documentation: $TARGET_DIR/README.md"
echo
echo "Next steps:"
echo "1. Review available services in $SERVICES_FILE"
echo "2. Create your values file based on the template"
echo "3. Run the deploy-elfhosted-charts.sh script to deploy the chart" 