# Self-Hosted Media Stack Setup Guide

This guide will help you set up a self-hosted version of the elfhosted media stack on your own infrastructure using Kubernetes (K3s).

## Overview

This setup leverages:

- **K3s**: A lightweight Kubernetes distribution perfect for edge and IoT deployments
- **Longhorn**: For distributed storage
- **MetalLB**: For load balancing
- **Traefik**: For ingress routing
- **Elfhosted's myprecious chart**: For deploying media applications

## Prerequisites

- One or more VPS instances with at least 4GB RAM and 2 CPU cores
- A domain name with DNS pointing to your server IP
- Basic familiarity with Linux and command line

## Quick Start

### 1. Prepare the environment

Clone this repository to your server:

```shell
git clone https://github.com/yourusername/my-media-stack.git
cd my-media-stack
```

Make the scripts executable:

```shell
chmod +x scripts/*.sh
```

### 2. Set up the K3s cluster

For a single-node setup:

```shell
./scripts/setup-k3s-cluster.sh --master --domain yourdomain.com
```

For a multi-node setup:

1. On the master node:

   ```shell
   ./scripts/setup-k3s-cluster.sh --master --domain yourdomain.com
   ```

   Note the node token and IP address displayed at the end of the setup.

2. On worker nodes:

   ```shell
   ./scripts/setup-k3s-cluster.sh --worker --ip <MASTER_IP> --token <NODE_TOKEN>
   ```

### 3. Set up access to elfhosted's myprecious chart

```shell
./scripts/setup-myprecious.sh
```

This will:

- Pull the elfhosted myprecious chart
- Extract information about available services
- Create a template values file
- Provide documentation

### 4. Configure your services

Create a values file for the myprecious chart:

```shell
cp helm-charts/myprecious/values-template.yaml helm-values/myprecious-values.yaml
```

Edit the values file to enable and configure the services you want:

```shell
nano helm-values/myprecious-values.yaml
```

Enable the services you want by setting `enabled: true` under each service. For example:

```yaml
plex:
  enabled: true
  env:
    TZ: "America/New_York"
    PLEX_CLAIM: "your-claim-token"
```

### 5. Deploy the media stack

Deploy the media stack using the elfhosted charts:

```shell
./scripts/deploy-elfhosted-charts.sh --domain yourdomain.com
```

This will:

- Create necessary persistent volumes
- Deploy the Kubernetes Dashboard
- Deploy all enabled services using elfhosted's myprecious chart

## Available Applications

The media stack includes many applications, including:

- **Media Servers**: Plex, Jellyfin
- **Media Management**: Wizarr, Riven
- **Download Clients**: Zurg (RealDebrid)
- **File Management**: FileBrowser, RcloneUI, RcloneFM
- **Dashboards**: Homer, Kubernetes Dashboard
- ***Arr Suite**: Sonarr, Radarr, Prowlarr, etc.
- **And many more**: Check `helm-charts/myprecious/available-services.txt` for the full list

## Accessing Your Applications

After deployment, your applications will be available at:

- Plex: <https://plex.yourdomain.com>
- Wizarr: <https://wizarr.yourdomain.com>
- Riven: <https://riven.yourdomain.com>
- Homer Dashboard: <https://homer.yourdomain.com>
- Kubernetes Dashboard: <https://k8s.yourdomain.com>

Use the token provided at the end of the deployment script to log in to the Kubernetes Dashboard.

## Customizing Your Setup

### Adding More Services

To enable additional services:

1. Edit your `helm-values/myprecious-values.yaml` file to enable more services
2. Run the deployment script again:

   ```shell
   ./scripts/deploy-elfhosted-charts.sh --domain yourdomain.com
   ```

### Scaling Your Cluster

To add more worker nodes, run the worker setup script on new VPS instances:

```shell
./scripts/setup-k3s-cluster.sh --worker --ip <MASTER_IP> --token <TOKEN>
```

## Maintenance Tasks

### Checking Application Status

```shell
kubectl get pods -n media-stack
```

### Viewing Application Logs

```shell
kubectl logs -n media-stack deployment/app-name
```

### Restarting an Application

```shell
kubectl rollout restart -n media-stack deployment/app-name
```

### Updating Applications

The elfhosted charts are regularly updated. To update to the latest versions:

```shell
./scripts/deploy-elfhosted-charts.sh --domain yourdomain.com
```

## Troubleshooting

### Common Issues

1. **Services not accessible via domain name**
   - Check that your DNS is configured correctly
   - Verify that Traefik ingress is working: `kubectl get ingress -A`

2. **Persistent volume claims pending**
   - Ensure Longhorn is installed and running: `kubectl get pods -n longhorn-system`

3. **Application pods crashing**
   - Check pod logs: `kubectl logs -n media-stack pod/<pod-name>`
   - Check pod events: `kubectl describe pod -n media-stack <pod-name>`

### Getting Help

If you encounter issues:

1. Check the Kubernetes Dashboard for pod status and events
2. Review logs using `kubectl logs`
3. Check the elfhosted repository and documentation for updates or known issues

## Advanced Configuration

### Custom Storage Classes

If you want to use a different storage class than Longhorn:

1. Edit your `helm-values/myprecious-values.yaml`:

   ```yaml
   global:
     storageClass: your-storage-class
   ```

2. Run the deployment script again

### Using Specific Chart Versions

To use a specific version of the myprecious chart:

```shell
./scripts/deploy-elfhosted-charts.sh --domain yourdomain.com --chart-version 1.2.3
```

## Security Considerations

By default, this setup:

- Uses HTTPS with Let's Encrypt certificates
- Configures applications with basic security settings

For production use, consider:

- Setting up proper authentication for all services
- Implementing network policies
- Using secrets management for sensitive configuration
- Regular security updates

## Backup and Recovery

### Creating Backups

Longhorn provides volume snapshots for backup:

```shell
kubectl apply -f longhorn-snapshot.yaml
```

### Restoring Backups

To restore from a Longhorn snapshot:

1. Stop the affected deployments
2. Restore the volume from the snapshot in Longhorn UI
3. Restart the deployments

## Conclusion

You now have a self-hosted media stack running on your own infrastructure. This setup leverages the same components used by elfhosted, allowing you to benefit from their updates and improvements while maintaining control of your own infrastructure.

For more information, refer to:

- [K3s Documentation](https://docs.k3s.io/)
- [Longhorn Documentation](https://longhorn.io/docs/)
- [Elfhosted's GitHub Repository](https://github.com/elfhosted)
