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

# Function to display help message
show_help() {
    echo -e "${BLUE}Media Stack Deployment Script${NC}"
    echo "This script deploys media applications to your K3s cluster."
    echo
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  -d, --domain DOMAIN          Domain name for services (required)"
    echo "  -n, --namespace NAMESPACE    Kubernetes namespace (default: media-stack)"
    echo "  -h, --help                   Show this help message"
    echo
    echo "Example:"
    echo "  $0 --domain example.com"
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

# Create namespace if it doesn't exist
echo -e "${BLUE}Creating namespace ${NAMESPACE}...${NC}"
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Create storage directories
echo -e "${BLUE}Setting up storage directories...${NC}"
mkdir -p ./storage/{config,media,downloads,transcode,backups}

# Create persistent volume claims for storage
echo -e "${BLUE}Creating persistent volume claims...${NC}"
cat > pvc-config.yaml << EOF
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
  name: backups
  namespace: ${NAMESPACE}
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: 5Gi
EOF

kubectl apply -f pvc-config.yaml
rm pvc-config.yaml

# Add Helm repositories
echo -e "${BLUE}Adding Helm repositories...${NC}"
helm repo add k8s-at-home https://k8s-at-home.com/charts/
helm repo add bjw-s https://bjw-s.github.io/helm-charts/
helm repo add truecharts https://charts.truecharts.org
helm repo update

# Deploy Homer (Dashboard)
echo -e "${BLUE}Deploying Homer dashboard...${NC}"
cat > homer-values.yaml << EOF
image:
  repository: b4bz/homer
  tag: latest
  pullPolicy: IfNotPresent

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "dashboard.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: homer-tls
        hosts:
          - "dashboard.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /www/assets

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000

# Create default dashboard configuration
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
          url: "https://files.${DOMAIN}"
          subtitle: "Manage Files"
          tag: "tools"
        - name: Rclone UI
          logo: "assets/tools/rclone.png"
          url: "https://rclone.${DOMAIN}"
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
EOF

helm install homer bjw-s/app-template -n $NAMESPACE -f homer-values.yaml
rm homer-values.yaml

# Deploy Plex
echo -e "${BLUE}Deploying Plex...${NC}"
cat > plex-values.yaml << EOF
image:
  repository: linuxserver/plex
  tag: latest
  pullPolicy: IfNotPresent

env:
  TZ: "UTC"
  PLEX_CLAIM: ""
  ADVERTISE_IP: "https://plex.${DOMAIN}"
  ALLOWED_NETWORKS: "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

service:
  main:
    ports:
      http:
        port: 32400

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "plex.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: plex-tls
        hosts:
          - "plex.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /config
  media:
    enabled: true
    existingClaim: media
    mountPath: /media
  transcode:
    enabled: true
    existingClaim: transcode
    mountPath: /transcode

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000

resources:
  requests:
    cpu: 200m
    memory: 500Mi
  limits:
    memory: 4Gi
EOF

helm install plex bjw-s/app-template -n $NAMESPACE -f plex-values.yaml
rm plex-values.yaml

# Deploy Wizarr
echo -e "${BLUE}Deploying Wizarr...${NC}"
cat > wizarr-values.yaml << EOF
image:
  repository: wizarr/wizarr
  tag: latest
  pullPolicy: IfNotPresent

env:
  TZ: "UTC"
  PGID: "1000"
  PUID: "1000"

service:
  main:
    ports:
      http:
        port: 5690

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "wizarr.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: wizarr-tls
        hosts:
          - "wizarr.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /config

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
EOF

helm install wizarr bjw-s/app-template -n $NAMESPACE -f wizarr-values.yaml
rm wizarr-values.yaml

# Deploy Riven (Debrid media management)
echo -e "${BLUE}Deploying Riven...${NC}"
cat > riven-values.yaml << EOF
image:
  repository: ghcr.io/elfhosted/riven
  tag: latest
  pullPolicy: IfNotPresent

env:
  TZ: "UTC"
  PGID: "1000"
  PUID: "1000"

service:
  main:
    ports:
      http:
        port: 7531

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "riven.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: riven-tls
        hosts:
          - "riven.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /config
  media:
    enabled: true
    existingClaim: media
    mountPath: /media
  downloads:
    enabled: true
    existingClaim: downloads
    mountPath: /downloads

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
EOF

helm install riven bjw-s/app-template -n $NAMESPACE -f riven-values.yaml
rm riven-values.yaml

# Deploy Zurg (RealDebrid mount)
echo -e "${BLUE}Deploying Zurg...${NC}"
cat > zurg-values.yaml << EOF
image:
  repository: ghcr.io/debridmediamanager/zurg
  tag: latest
  pullPolicy: IfNotPresent

env:
  TZ: "UTC"
  PGID: "1000"
  PUID: "1000"

service:
  main:
    ports:
      http:
        port: 9999

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "zurg.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: zurg-tls
        hosts:
          - "zurg.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /config
  media:
    enabled: true
    existingClaim: media
    mountPath: /mnt/unionfs/zurg

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
EOF

helm install zurg bjw-s/app-template -n $NAMESPACE -f zurg-values.yaml
rm zurg-values.yaml

# Deploy FileBrowser
echo -e "${BLUE}Deploying FileBrowser...${NC}"
cat > filebrowser-values.yaml << EOF
image:
  repository: filebrowser/filebrowser
  tag: latest
  pullPolicy: IfNotPresent

env:
  TZ: "UTC"
  FB_PORT: "80"

service:
  main:
    ports:
      http:
        port: 80

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "files.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: filebrowser-tls
        hosts:
          - "files.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /config
  media:
    enabled: true
    existingClaim: media
    mountPath: /media
  downloads:
    enabled: true
    existingClaim: downloads
    mountPath: /downloads
  database:
    enabled: true
    type: emptyDir
    mountPath: /database

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
EOF

helm install filebrowser bjw-s/app-template -n $NAMESPACE -f filebrowser-values.yaml
rm filebrowser-values.yaml

# Deploy Rclone UI
echo -e "${BLUE}Deploying Rclone UI...${NC}"
cat > rcloneui-values.yaml << EOF
image:
  repository: rclone/rclone
  tag: latest
  pullPolicy: IfNotPresent

command:
  - rclone
  - rcd
  - --rc-web-gui
  - --rc-addr=0.0.0.0:5572
  - --rc-user=admin
  - --rc-pass=admin
  - --rc-serve

service:
  main:
    ports:
      http:
        port: 5572

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "rclone.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: rcloneui-tls
        hosts:
          - "rclone.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /config/rclone
  media:
    enabled: true
    existingClaim: media
    mountPath: /media

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
EOF

helm install rcloneui bjw-s/app-template -n $NAMESPACE -f rcloneui-values.yaml
rm rcloneui-values.yaml

# Deploy Rclone FM
echo -e "${BLUE}Deploying Rclone FM...${NC}"
cat > rclonefm-values.yaml << EOF
image:
  repository: romancin/rclonebrowser
  tag: latest
  pullPolicy: IfNotPresent

env:
  TZ: "UTC"
  PUID: "1000"
  PGID: "1000"

service:
  main:
    ports:
      http:
        port: 5800

ingress:
  main:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: "traefik"
      traefik.ingress.kubernetes.io/router.entrypoints: "websecure"
    hosts:
      - host: "rclonefm.${DOMAIN}"
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: rclonefm-tls
        hosts:
          - "rclonefm.${DOMAIN}"

persistence:
  config:
    enabled: true
    existingClaim: config
    mountPath: /config
  media:
    enabled: true
    existingClaim: media
    mountPath: /media

podSecurityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
EOF

helm install rclonefm bjw-s/app-template -n $NAMESPACE -f rclonefm-values.yaml
rm rclonefm-values.yaml

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

echo -e "${GREEN}Media applications deployment completed!${NC}"
echo "Access your applications at the following URLs:"
echo "- Dashboard: https://dashboard.${DOMAIN}"
echo "- Plex: https://plex.${DOMAIN}"
echo "- Wizarr: https://wizarr.${DOMAIN}"
echo "- Riven: https://riven.${DOMAIN}"
echo "- Zurg: https://zurg.${DOMAIN}"
echo "- FileBrowser: https://files.${DOMAIN}"
echo "- Rclone UI: https://rclone.${DOMAIN}"
echo "- Rclone FM: https://rclonefm.${DOMAIN}"
echo "- Kubernetes Dashboard: https://k8s.${DOMAIN}" 