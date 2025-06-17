#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration variables
NAMESPACE="media-stack"
DOMAIN=""
VALUES_DIR="./helm-values"
ENABLED_SERVICES=""
CHART_REPO="oci://ghcr.io/elfhosted/charts"
MYPRECIOUS_CHART="myprecious"
CHART_VERSION="latest"

# Function to display help message
show_help() {
    echo -e "${BLUE}Elfhosted Media Stack Deployment Script${NC}"
    echo "This script deploys media applications from elfhosted's charts."
    echo
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  -d, --domain DOMAIN          Domain name for services (required)"
    echo "  -n, --namespace NAMESPACE    Kubernetes namespace (default: media-stack)"
    echo "  -s, --services SERVICES      Comma-separated list of services to enable (default: all shown in UI)"
    echo "  -v, --values-dir DIR         Directory for Helm values files (default: ./helm-values)"
    echo "  -c, --chart-version VERSION  Chart version to use (default: latest)"
    echo "  -r, --repo REPO              Helm chart repository (default: oci://ghcr.io/elfhosted/charts)"
    echo "  -h, --help                   Show this help message"
    echo
    echo "Examples:"
    echo "  # Deploy all core services:"
    echo "  $0 --domain example.com"
    echo
    echo "  # Deploy specific services:"
    echo "  $0 --domain example.com --services plex,wizarr,riven,zurg,homer"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--domain)
            DOMAIN="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -s|--services)
            ENABLED_SERVICES="$2"
            shift 2
            ;;
        -v|--values-dir)
            VALUES_DIR="$2"
            shift 2
            ;;
        -c|--chart-version)
            CHART_VERSION="$2"
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

# Validation checks
if [[ -z "$DOMAIN" ]]; then
    echo -e "${RED}Error: Domain name is required (--domain)${NC}"
    show_help
    exit 1
fi

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}Error: kubectl is not installed${NC}"
    exit 1
fi

# Check if helm is available
if ! command -v helm &> /dev/null; then
    echo -e "${RED}Error: helm is not installed${NC}"
    exit 1
fi

# Create values directory if it doesn't exist
mkdir -p "$VALUES_DIR"

# Create namespace if it doesn't exist
echo -e "${BLUE}Creating namespace ${NAMESPACE}...${NC}"
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# If no services specified, use defaults shown in UI
if [[ -z "$ENABLED_SERVICES" ]]; then
    ENABLED_SERVICES="plex,wizarr,riven,zurg,filebrowser,rclonefm,rcloneui,homer,kubernetes-dashboard"
fi

# Create persistent volume claims for storage
echo -e "${BLUE}Creating persistent volume claims...${NC}"
cat > pvc.yaml << EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: config
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 10Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: media
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 100Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: downloads
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 50Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: transcode
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 20Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: backup
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 5Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: symlinks
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: logs
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: rclone
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: realdebrid-zurg
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 1Gi
EOF

kubectl apply -f pvc.yaml
rm pvc.yaml

# Add Helm repository
echo -e "${BLUE}Adding Helm repository...${NC}"
helm repo add bjw-s https://bjw-s.github.io/helm-charts/
helm repo update

# Configure values for myprecious chart
echo -e "${BLUE}Creating values file for myprecious chart...${NC}"
cat > "$VALUES_DIR/myprecious-values.yaml" << EOF
# Global settings
global:
  namespace: ${NAMESPACE}
  domain: ${DOMAIN}
  storageClass: longhorn
  uid: 1000
  gid: 1000
  umask: "002"
  
# Service enablement
$(for service in $(echo $ENABLED_SERVICES | tr ',' ' '); do
  echo "${service}:"
  echo "  enabled: true"
  echo "  ingress:"
  echo "    hosts:"
  echo "      - host: ${service}.${DOMAIN}"
  echo "        paths:"
  echo "          - path: /"
  echo "            pathType: Prefix"
  echo "    tls:"
  echo "      - secretName: ${service}-tls"
  echo "        hosts:"
  echo "          - ${service}.${DOMAIN}"
  echo "  persistence:"
  echo "    config:"
  echo "      enabled: true"
  echo "      existingClaim: config"
  echo "    media:"
  echo "      enabled: true"
  echo "      existingClaim: media"
  echo "    downloads:"
  echo "      enabled: true"
  echo "      existingClaim: downloads"
  echo "    transcode:"
  echo "      enabled: true"
  echo "      existingClaim: transcode"
  echo "    backup:"
  echo "      enabled: true"
  echo "      existingClaim: backup"
  echo "    symlinks:"
  echo "      enabled: true"
  echo "      existingClaim: symlinks"
  echo "    logs:"
  echo "      enabled: true"
  echo "      existingClaim: logs"
  echo "    rclone:"
  echo "      enabled: true"
  echo "      existingClaim: rclone"
  echo "    realdebrid-zurg:"
  echo "      enabled: true"
  echo "      existingClaim: realdebrid-zurg"
done)

# Homer dashboard configuration
homer:
  enabled: true
  config:
    title: "Media Dashboard"
    subtitle: "Media Server Management"
    logo: "assets/logo.png"
    footer: "Made with ♥"
    theme: default
    columns: "auto"
    connectivityCheck: true
    services:
      - name: Consume Media
        icon: "fas fa-play-circle"
        items:
          - name: Plex
            logo: "assets/tools/plex.png"
            url: "https://plex.${DOMAIN}"
            subtitle: "Watch Movies/TV"
            tag: "media"
      - name: Manage Media
        icon: "fas fa-sliders-h"
        items:
          - name: Wizarr
            logo: "assets/tools/wizarr.png"
            url: "https://wizarr.${DOMAIN}"
            subtitle: "Manage Plex/Jellyfin Accounts"
            tag: "management"
          - name: Riven
            logo: "assets/tools/riven.png"
            url: "https://riven.${DOMAIN}"
            subtitle: "All-in-one Debrid Media Management"
            tag: "management"
      - name: Download Media
        icon: "fas fa-download"
        items:
          - name: Zurg
            logo: "assets/tools/zurg.png"
            url: "https://zurg.${DOMAIN}"
            subtitle: "Manage RealDebrid Mount"
            tag: "download"
      - name: Tools
        icon: "fas fa-tools"
        items:
          - name: FileBrowser
            logo: "assets/tools/filebrowser.png"
            url: "https://filebrowser.${DOMAIN}"
            subtitle: "Manage Files"
            tag: "tools"
          - name: Rclone UI
            logo: "assets/tools/rclone.png"
            url: "https://rcloneui.${DOMAIN}"
            subtitle: "Rclone UI"
            tag: "tools"
          - name: Rclone FM
            logo: "assets/tools/rclone.png"
            url: "https://rclonefm.${DOMAIN}"
            subtitle: "Rclone File Manager"
            tag: "tools"
          - name: Kubernetes Dashboard
            logo: "assets/tools/kubernetes.png"
            url: "https://k8s.${DOMAIN}"
            subtitle: "Kubernetes Dashboard"
            tag: "tools"

# Specific service configurations
plex:
  env:
    TZ: "UTC"
    PLEX_CLAIM: ""
    ADVERTISE_IP: "https://plex.${DOMAIN}"
    ALLOWED_NETWORKS: "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

wizarr:
  env:
    TZ: "UTC"
    PGID: "1000"
    PUID: "1000"

riven:
  env:
    TZ: "UTC"
    PGID: "1000"
    PUID: "1000"

zurg:
  env:
    TZ: "UTC"
    PGID: "1000"
    PUID: "1000"

filebrowser:
  env:
    TZ: "UTC"
    FB_PORT: "80"

rcloneui:
  command:
    - rclone
    - rcd
    - --rc-web-gui
    - --rc-addr=0.0.0.0:5572
    - --rc-user=admin
    - --rc-pass=admin
    - --rc-serve

rclonefm:
  env:
    TZ: "UTC"
    PUID: "1000"
    PGID: "1000"
EOF

# Deploy Kubernetes Dashboard
echo -e "${BLUE}Deploying Kubernetes Dashboard...${NC}"
DASHBOARD_VERSION=$(curl -s https://api.github.com/repos/kubernetes/dashboard/releases/latest | grep "tag_name" | cut -d '"' -f 4)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/${DASHBOARD_VERSION}/aio/deploy/recommended.yaml

# Create admin user for Kubernetes Dashboard
cat > dashboard-admin-user.yaml << EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kubernetes-dashboard
EOF

kubectl apply -f dashboard-admin-user.yaml
rm dashboard-admin-user.yaml

# Create Ingress for Kubernetes Dashboard
cat > dashboard-ingress.yaml << EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
  annotations:
    kubernetes.io/ingress.class: "traefik"
    traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
spec:
  tls:
  - hosts:
    - "k8s.${DOMAIN}"
    secretName: kubernetes-dashboard-tls
  rules:
  - host: "k8s.${DOMAIN}"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: kubernetes-dashboard
            port:
              number: 443
EOF

kubectl apply -f dashboard-ingress.yaml
rm dashboard-ingress.yaml

# Get token for Kubernetes Dashboard
echo -e "${BLUE}Getting token for Kubernetes Dashboard...${NC}"
TOKEN=$(kubectl -n kubernetes-dashboard create token admin-user)
echo -e "${GREEN}Kubernetes Dashboard Token:${NC}"
echo "$TOKEN"

# Clone the myprecious chart from the elfhosted repository if needed
if [[ "$CHART_VERSION" == "latest" ]]; then
    echo -e "${BLUE}Installing media applications using elfhosted's myprecious chart...${NC}"
    
    # Set up Helm OCI support
    export HELM_EXPERIMENTAL_OCI=1
    
    # Install or upgrade the chart
    helm upgrade --install myprecious \
      $CHART_REPO/$MYPRECIOUS_CHART \
      --namespace "$NAMESPACE" \
      --create-namespace \
      -f "$VALUES_DIR/myprecious-values.yaml"
else
    echo -e "${BLUE}Installing media applications using elfhosted's myprecious chart version ${CHART_VERSION}...${NC}"
    
    # Set up Helm OCI support
    export HELM_EXPERIMENTAL_OCI=1
    
    # Install or upgrade the chart
    helm upgrade --install myprecious \
      $CHART_REPO/$MYPRECIOUS_CHART \
      --version "$CHART_VERSION" \
      --namespace "$NAMESPACE" \
      --create-namespace \
      -f "$VALUES_DIR/myprecious-values.yaml"
fi

echo -e "${GREEN}Media applications deployment completed!${NC}"
echo "Access your applications at the following URLs:"

# List the enabled services
for service in $(echo $ENABLED_SERVICES | tr ',' ' '); do
  if [[ "$service" == "kubernetes-dashboard" ]]; then
    echo "- Kubernetes Dashboard: https://k8s.${DOMAIN}"
  else
    echo "- ${service^}: https://${service}.${DOMAIN}"
  fi
done 