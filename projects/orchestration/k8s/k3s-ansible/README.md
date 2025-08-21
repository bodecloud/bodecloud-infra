# Media Stack K3s Cluster

This repository contains Ansible playbooks and configurations for deploying a standardized K3s Kubernetes cluster optimized for media streaming applications. It is based on the official [k3s-ansible](https://github.com/k3s-io/k3s-ansible) approach but customized for a media stack deployment.

## Features

- **Standardized Deployment**: Uses the official k3s-ansible approach for deploying K3s
- **Tailscale Integration**: Designed to work with Tailscale for secure inter-node communication
- **Easy Operation**: Simple commands for installation, verification, and reset
- **Customizable**: Easily adaptable to different environments and requirements
- **Idempotent**: Can be run multiple times without breaking things

## Prerequisites

- **Ansible 2.15+** on the control node
- **Ubuntu** or **Debian** on target nodes
- **SSH key-based access** to all nodes
- **Tailscale** configured on all nodes
- **Sudo privileges** on all nodes

## Directory Structure

```shell
k3s-ansible/
├── collections/           # Ansible collections requirements
├── inventory/             # Inventory files
│   └── hosts.yml          # Main inventory file
├── playbooks/             # Ansible playbooks
│   ├── reset.yml          # Playbook to uninstall k3s
│   └── site.yml           # Main installation playbook
├── roles/                 # Custom roles (if needed)
├── templates/             # Templates for configuration files
│   ├── k3s-agent.service.j2   # Agent systemd service template
│   └── k3s-server.service.j2  # Server systemd service template
├── ansible.cfg            # Ansible configuration
├── deploy.sh              # Main deployment script
├── Makefile               # Make commands for easier operation
└── README.md              # This file
```

## Quick Start

1. Edit the `inventory/hosts.yml` file to match your environment.

2. Run the deployment script:

    ```bash
    ./deploy.sh
    ```

3. After successful deployment, access your Kubernetes cluster:

    ```bash
    export KUBECONFIG=./kubeconfig
    kubectl get nodes
    ```

## Usage

### Installation

```bash
# Full installation
./deploy.sh

# Or using make
make install
```

### Reset/Uninstall

```bash
# Reset the cluster
make reset
```

### Check Node Connectivity

```bash
# Ping all nodes
make ping
```

### Run a Specific Playbook

```bash
# Apply a specific playbook
make apply playbook=<playbook-name>.yml
```

## Configuration

### Inventory Configuration

The `inventory/hosts.yml` file contains all node definitions and configuration:

```yaml
k3s_cluster:
  children:
    server:
      hosts:
        # Server nodes go here
    agent:
      hosts:
        # Agent nodes go here
  vars:
    # Global variables go here
```

### Important Configuration Options

- **k3s_version**: Version of K3s to install (default: `v1.33.0+k3s1`)
- **k3s_token**: Cluster token (auto-generated if not specified)
- **k3s_server_args**: Arguments for the K3s server

## Next Steps After Deployment

After successful deployment, consider the following next steps:

1. **Install a Storage Provider**:

   ```bash
   kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/master/deploy/longhorn.yaml
   ```

2. **Install MetalLB for Load Balancing**:

   ```bash
   kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml
   ```

3. **Install Traefik Ingress Controller**:

   ```bash
   helm repo add traefik https://helm.traefik.io/traefik
   helm repo update
   helm install traefik traefik/traefik -n kube-system
   ```

4. **Install Cert-Manager for TLS certificates**:

   ```bash
   kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml
   ```

## Using Helmfile for Application Deployment

After setting up the K3s cluster, you can use Helmfile for declarative application deployment:

1. Install Helmfile:

   ```bash
   wget https://github.com/helmfile/helmfile/releases/download/v0.157.0/helmfile_0.157.0_linux_amd64.tar.gz
   tar -xzf helmfile_0.157.0_linux_amd64.tar.gz
   sudo mv helmfile /usr/local/bin/
   ```

2. Create a `helmfile.yaml` file:

   ```yaml
   repositories:
   - name: prometheus-community
     url: https://prometheus-community.github.io/helm-charts
   
   releases:
   - name: prometheus
     namespace: monitoring
     chart: prometheus-community/prometheus
     createNamespace: true
     values:
     - values/prometheus.yaml
   ```

3. Apply the Helmfile:

   ```bash
   helmfile apply
   ```

## Troubleshooting

### Common Issues

1. **Nodes Not Joining the Cluster**:
   - Verify Tailscale is working properly
   - Check firewall settings
   - Verify the token is correctly distributed

2. **Failed Installation**:
   - Run `make reset` to clean up
   - Check the logs at `/var/log/syslog` on the nodes
   - Try with `./deploy.sh --check-only` to verify prerequisites

### Reset and Retry

If the installation fails, you can reset the cluster and try again:

```bash
make reset
./deploy.sh
```

## License

This project is licensed under the MIT License.
