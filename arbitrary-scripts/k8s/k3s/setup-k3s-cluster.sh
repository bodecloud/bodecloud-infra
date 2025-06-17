#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration variables
K3S_VERSION="v1.28.6+k3s2"
CLUSTER_NAME="bolabaden"
NODE_TOKEN=""
MASTER_IP=""
MASTER_NODE=false
WORKER_NODE=false
INSTALL_LONGHORN=true
INSTALL_METALLB=true
INSTALL_FLUX=true
INSTALL_TRAEFIK=true
SETUP_DNS=true
DOMAIN=""
ENABLE_METRICS=true
DISABLE_TRAEFIK=true
DISABLE_SERVICELB=true
CILIUM_ENABLED=true

# Function to display help message
show_help() {
    echo -e "${BLUE}K3s Bolabaden Cluster Setup Script${NC}"
    echo "This script sets up a K3s cluster for running bolabaden."
    echo
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  -m, --master                 Configure this node as a master node"
    echo "  -w, --worker                 Configure this node as a worker node"
    echo "  -i, --ip IP                  Master node IP address (required for worker nodes)"
    echo "  -t, --token TOKEN            Node token (required for worker nodes)"
    echo "  -n, --name NAME              Cluster name (default: bolabaden)"
    echo "  -d, --domain DOMAIN          Domain name for the cluster"
    echo "  --no-longhorn                Skip Longhorn installation"
    echo "  --no-metallb                 Skip MetalLB installation"
    echo "  --no-flux                    Skip Flux installation"
    echo "  --no-traefik                 Skip Traefik installation"
    echo "  --no-dns                     Skip DNS setup"
    echo "  --no-metrics                 Disable metrics server"
    echo "  --enable-k3s-traefik         Keep K3s built-in Traefik"
    echo "  --enable-k3s-servicelb       Keep K3s built-in ServiceLB"
    echo "  --disable-cilium             Use Flannel instead of Cilium"
    echo "  -h, --help                   Show this help message"
    echo
    echo "Examples:"
    echo "  # Set up a master node:"
    echo "  $0 --master --domain example.com"
    echo
    echo "  # Set up a worker node:"
    echo "  $0 --worker --ip 192.168.1.10 --token TOKEN"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--master)
            MASTER_NODE=true
            shift
            ;;
        -w|--worker)
            WORKER_NODE=true
            shift
            ;;
        -i|--ip)
            MASTER_IP="$2"
            shift 2
            ;;
        -t|--token)
            NODE_TOKEN="$2"
            shift 2
            ;;
        -n|--name)
            CLUSTER_NAME="$2"
            shift 2
            ;;
        -d|--domain)
            DOMAIN="$2"
            shift 2
            ;;
        --no-longhorn)
            INSTALL_LONGHORN=false
            shift
            ;;
        --no-metallb)
            INSTALL_METALLB=false
            shift
            ;;
        --no-flux)
            INSTALL_FLUX=false
            shift
            ;;
        --no-traefik)
            INSTALL_TRAEFIK=false
            shift
            ;;
        --no-dns)
            SETUP_DNS=false
            shift
            ;;
        --no-metrics)
            ENABLE_METRICS=false
            shift
            ;;
        --enable-k3s-traefik)
            DISABLE_TRAEFIK=false
            shift
            ;;
        --enable-k3s-servicelb)
            DISABLE_SERVICELB=false
            shift
            ;;
        --disable-cilium)
            CILIUM_ENABLED=false
            shift
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
if [[ "$MASTER_NODE" == "false" && "$WORKER_NODE" == "false" ]]; then
    echo -e "${RED}Error: You must specify either --master or --worker${NC}"
    show_help
    exit 1
fi

if [[ "$MASTER_NODE" == "true" && "$WORKER_NODE" == "true" ]]; then
    echo -e "${RED}Error: You cannot specify both --master and --worker${NC}"
    show_help
    exit 1
fi

if [[ "$WORKER_NODE" == "true" && -z "$MASTER_IP" ]]; then
    echo -e "${RED}Error: Worker node requires master IP (--ip)${NC}"
    show_help
    exit 1
fi

if [[ "$WORKER_NODE" == "true" && -z "$NODE_TOKEN" ]]; then
    echo -e "${RED}Error: Worker node requires node token (--token)${NC}"
    show_help
    exit 1
fi

# Install required packages
echo -e "${BLUE}Installing required packages...${NC}"
apt-get update
apt-get install -y curl wget jq apt-transport-https ca-certificates gnupg lsb-release nfs-common open-iscsi

# Configure system settings
echo -e "${BLUE}Configuring system settings...${NC}"
cat > /etc/sysctl.d/99-kubernetes.conf << EOF
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
vm.max_map_count = 262144
EOF
sysctl --system

# Enable and start iscsid for Longhorn
systemctl enable iscsid
systemctl start iscsid

# Create K3s configuration
if [[ "$MASTER_NODE" == "true" ]]; then
    echo -e "${BLUE}Creating K3s configuration for master node...${NC}"
    mkdir -p /etc/rancher/k3s
    
    # Create K3s config file
    cat > /etc/rancher/k3s/config.yaml << EOF
# K3s Configuration
cluster-name: ${CLUSTER_NAME}
EOF

    # Configure disabled components
    DISABLED_COMPONENTS=""
    
    if [[ "$DISABLE_TRAEFIK" == "true" ]]; then
        DISABLED_COMPONENTS="${DISABLED_COMPONENTS}traefik,"
    fi
    
    if [[ "$DISABLE_SERVICELB" == "true" ]]; then
        DISABLED_COMPONENTS="${DISABLED_COMPONENTS}servicelb,"
    fi
    
    if [[ "$CILIUM_ENABLED" == "true" ]]; then
        DISABLED_COMPONENTS="${DISABLED_COMPONENTS}flannel,"
    fi
    
    # Add additional settings
    cat >> /etc/rancher/k3s/config.yaml << EOF
disable: 
$(echo $DISABLED_COMPONENTS | tr ',' '\n' | sed 's/^/- "/' | sed 's/$/"/')
$(if [[ "$CILIUM_ENABLED" == "true" ]]; then echo "flannel-backend: none"; fi)
$(if [[ "$ENABLE_METRICS" == "true" ]]; then echo "kube-controller-manager-arg:"; echo "- \"bind-address=0.0.0.0\""; echo "kube-scheduler-arg:"; echo "- \"bind-address=0.0.0.0\""; fi)
kubelet-arg:
- "max-pods=250"
node-ip: $(hostname -I | awk '{print $1}')
advertise-address: $(hostname -I | awk '{print $1}')
tls-san:
- $(hostname -I | awk '{print $1}')
etcd-expose-metrics: true
EOF

    # Install K3s as master
    echo -e "${BLUE}Installing K3s master node...${NC}"
    curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=${K3S_VERSION} sh -

    # Wait for K3s to be ready
    echo -e "${BLUE}Waiting for K3s to be ready...${NC}"
    until kubectl get nodes; do
        echo "Waiting for K3s to start..."
        sleep 5
    done

    # Get node token for future worker nodes
    NODE_TOKEN=$(cat /var/lib/rancher/k3s/server/node-token)
    echo -e "${GREEN}Node token for joining workers:${NC} ${NODE_TOKEN}"
    echo -e "${GREEN}Master node IP:${NC} $(hostname -I | awk '{print $1}')"

    # Configure kubectl
    echo -e "${BLUE}Configuring kubectl...${NC}"
    mkdir -p $HOME/.kube
    cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/config
    chmod 600 $HOME/.kube/config
    export KUBECONFIG=$HOME/.kube/config
    echo "export KUBECONFIG=$HOME/.kube/config" >> $HOME/.bashrc

    # Install Helm
    echo -e "${BLUE}Installing Helm...${NC}"
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

    # Setup Cilium if enabled
    if [[ "$CILIUM_ENABLED" == "true" ]]; then
        echo -e "${BLUE}Installing Cilium...${NC}"
        helm repo add cilium https://helm.cilium.io/
        helm repo update
        helm install cilium cilium/cilium --version 1.14.4 \
            --namespace kube-system \
            --set kubeProxyReplacement=strict \
            --set k8sServiceHost=$(hostname -I | awk '{print $1}') \
            --set k8sServicePort=6443
    fi

    # Install MetalLB if enabled
    if [[ "$INSTALL_METALLB" == "true" ]]; then
        echo -e "${BLUE}Installing MetalLB...${NC}"
        METALLB_VERSION=$(curl -s https://api.github.com/repos/metallb/metallb/releases/latest | grep "tag_name" | cut -d '"' -f 4)
        kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/${METALLB_VERSION}/config/manifests/metallb-native.yaml

        # Wait for MetalLB to be ready
        echo "Waiting for MetalLB to be ready..."
        kubectl wait --namespace metallb-system \
            --for=condition=ready pod \
            --selector=app=metallb \
            --timeout=90s
        
        # Get the node's IP range for MetalLB
        NODE_IP=$(hostname -I | awk '{print $1}')
        IP_PREFIX=$(echo $NODE_IP | cut -d. -f1-3)
        
        # Create MetalLB configuration
        cat > metallb-config.yaml << EOF
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: first-pool
  namespace: metallb-system
spec:
  addresses:
  - ${IP_PREFIX}.200-${IP_PREFIX}.250
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: l2-advert
  namespace: metallb-system
spec:
  ipAddressPools:
  - first-pool
EOF
        
        kubectl apply -f metallb-config.yaml
        rm metallb-config.yaml
    fi

    # Install Traefik if enabled
    if [[ "$INSTALL_TRAEFIK" == "true" ]]; then
        echo -e "${BLUE}Installing Traefik...${NC}"
        helm repo add traefik https://traefik.github.io/charts
        helm repo update
        
        # Create values file for Traefik
        cat > traefik-values.yaml << EOF
deployment:
  replicas: 1
ingressRoute:
  dashboard:
    enabled: true
ports:
  web:
    redirectTo: websecure
  websecure:
    tls:
      enabled: true
service:
  enabled: true
  type: LoadBalancer
persistence:
  enabled: true
  size: 128Mi
additionalArguments:
  - "--api.dashboard=true"
  - "--providers.kubernetesingress.ingressclass=traefik"
  - "--log.level=INFO"
EOF
        
        helm install traefik traefik/traefik -n kube-system -f traefik-values.yaml
        rm traefik-values.yaml
    fi

    # Install Longhorn if enabled
    if [[ "$INSTALL_LONGHORN" == "true" ]]; then
        echo -e "${BLUE}Installing Longhorn...${NC}"
        LONGHORN_VERSION=$(curl -s https://api.github.com/repos/longhorn/longhorn/releases/latest | grep "tag_name" | cut -d '"' -f 4)
        kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/${LONGHORN_VERSION}/deploy/longhorn.yaml
        
        # Wait for Longhorn to be ready
        echo "Waiting for Longhorn to be ready..."
        kubectl -n longhorn-system wait --for=condition=ready pod --all --timeout=300s
    fi

    # Install Flux if enabled
    if [[ "$INSTALL_FLUX" == "true" ]]; then
        echo -e "${BLUE}Installing Flux...${NC}"
        # Install Flux CLI
        curl -s https://fluxcd.io/install.sh | bash
        
        # Bootstrap Flux
        flux install --components=source-controller,kustomize-controller,helm-controller,notification-controller
    fi
    
    echo -e "${GREEN}Master node setup complete!${NC}"
    echo "You can now join worker nodes using:"
    echo "./setup-k3s-cluster.sh --worker --ip $(hostname -I | awk '{print $1}') --token ${NODE_TOKEN}"

elif [[ "$WORKER_NODE" == "true" ]]; then
    echo -e "${BLUE}Installing K3s worker node...${NC}"
    
    # Install K3s as worker
    curl -sfL https://get.k3s.io | K3S_URL=https://${MASTER_IP}:6443 K3S_TOKEN=${NODE_TOKEN} INSTALL_K3S_VERSION=${K3S_VERSION} sh -
    
    echo -e "${GREEN}Worker node joined the cluster successfully!${NC}"
fi

echo -e "${GREEN}K3s installation completed successfully!${NC}" 