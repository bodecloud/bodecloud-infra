# Filebrowser Helm Chart

A Helm chart for deploying Filebrowser - a web-based file manager with terminal capabilities, optimized for the ElfHosted platform.

## Overview

Filebrowser provides a web interface for managing files and directories on a server. This chart includes:

- **Main Filebrowser container** - Web-based file manager
- **Elfbot sidecar** - Automation and orchestration helper
- **Init containers** - For setup and configuration
- **Persistent storage** - For configs, backups, logs, and file storage
- **Service account** - With appropriate permissions

## Features

- 🗂️ **File Management** - Upload, download, move, copy, delete files
- 🖥️ **Terminal Access** - Built-in terminal with customizable commands  
- 🔐 **Authentication** - Integration with proxy authentication
- 📁 **Multiple Storage** - Config, backup, logs, rclone, symlinks
- 🤖 **Elfbot Integration** - Automated management and restarts
- 📊 **Monitoring** - Health checks and probes
- 🔒 **Security** - Non-root containers, security contexts

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure

## Installation

### Quick Start

```bash
# Add the helm repository
helm repo add elfhosted https://charts.elfhosted.com
helm repo update

# Install with default values
helm install filebrowser elfhosted/filebrowser
```

### Custom Installation

```bash
# Install with custom values
helm install filebrowser elfhosted/filebrowser \
  --set persistence.config.size=20Gi \
  --set persistence.backup.size=50Gi \
  --set branding.name="My FileBrowser"
```

### Using a values file

```bash
# Install with custom values file
helm install filebrowser elfhosted/filebrowser -f my-values.yaml
```

## Configuration

The following table lists the configurable parameters of the Filebrowser chart and their default values.

### Global Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.imageRegistry` | Global Docker image registry | `ghcr.io/elfhosted` |
| `global.imagePullSecrets` | Global Docker registry secret names | `[]` |
| `global.storageClass` | Global StorageClass for Persistent Volume(s) | `local-path` |

### Image Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Filebrowser image repository | `filebrowser` |
| `image.tag` | Filebrowser image tag | `2.23.0` |
| `image.digest` | Filebrowser image digest | `sha256:296e3a3d08...` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |

### Application Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of Filebrowser replicas | `1` |
| `filebrowser.allowedCommands` | Comma-separated list of allowed terminal commands | `zip,unzip,rar,unrar,ls,pwd,cd,mv,cp,ln,find,echo,grep,cat,touch,tar,gzip,rm,tree,du,mlocate,updatedb,locate,elfbot,Elfbot` |
| `branding.name` | Custom branding name | `ElfHosted FileBrowser 🧝` |
| `env.TZ` | Timezone | `UTC` |

### Service Account

| Parameter | Description | Default |
|-----------|-------------|---------|
| `serviceAccount.create` | Specifies whether a service account should be created | `true` |
| `serviceAccount.name` | The name of the service account | `filebrowser` |
| `serviceAccount.automount` | Automount service account token | `true` |
| `serviceAccount.annotations` | Annotations to add to the service account | `{}` |

### Security Context

| Parameter | Description | Default |
|-----------|-------------|---------|
| `securityContext.fsGroup` | Group ID for the container | `568` |
| `securityContext.fsGroupChangePolicy` | Policy for changing ownership and permissions | `OnRootMismatch` |
| `containerSecurityContext.runAsUser` | User ID for the container | `568` |
| `containerSecurityContext.runAsGroup` | Group ID for the container | `568` |
| `containerSecurityContext.readOnlyRootFilesystem` | Mount root filesystem as read-only | `false` |
| `containerSecurityContext.allowPrivilegeEscalation` | Allow privilege escalation | `false` |

### Service Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.port` | Service port | `8080` |
| `service.targetPort` | Target port | `http` |
| `service.annotations` | Annotations to add to service | `{}` |

### Ingress Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable ingress controller resource | `false` |
| `ingress.className` | IngressClass that will be used | `""` |
| `ingress.annotations` | Additional annotations for the Ingress resource | `{}` |
| `ingress.hosts` | An array with hosts and paths | `[{host: "filebrowser.local", paths: [{path: "/", pathType: "Prefix"}]}]` |
| `ingress.tls` | An array with the tls configuration | `[]` |

### Persistence Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `persistence.config.enabled` | Enable persistence for configuration | `true` |
| `persistence.config.size` | Size of persistent volume claim | `10Gi` |
| `persistence.config.storageClass` | Type of persistent volume claim | `""` |
| `persistence.config.accessModes` | Persistent volume access modes | `["ReadWriteOnce"]` |
| `persistence.backup.enabled` | Enable persistence for backups | `true` |
| `persistence.backup.size` | Size of backup persistent volume claim | `20Gi` |
| `persistence.backup.annotations` | Annotations for backup PVC | `{helm.sh/resource-policy: keep}` |
| `persistence.logs.enabled` | Enable persistence for logs | `true` |
| `persistence.logs.size` | Size of logs persistent volume claim | `5Gi` |
| `persistence.rclone.enabled` | Enable persistence for rclone | `true` |
| `persistence.rclone.size` | Size of rclone persistent volume claim | `1Gi` |
| `persistence.symlinks.enabled` | Enable persistence for symlinks | `true` |
| `persistence.symlinks.size` | Size of symlinks persistent volume claim | `1Gi` |
| `persistence.realdebridZurg.enabled` | Enable RealDebrid Zurg storage | `true` |
| `persistence.realdebridZurg.existingClaim` | Name of existing claim for RealDebrid | `realdebrid-zurg` |

### Resource Management

| Parameter | Description | Default |
|-----------|-------------|---------|
| `resources.limits.cpu` | CPU limit | `1` |
| `resources.limits.memory` | Memory limit | `1Gi` |
| `resources.requests.cpu` | CPU request | `0` |
| `resources.requests.memory` | Memory request | `6Mi` |

### Elfbot Sidecar

| Parameter | Description | Default |
|-----------|-------------|---------|
| `elfbot.enabled` | Enable Elfbot sidecar container | `true` |
| `elfbot.image.repository` | Elfbot image repository | `tooling` |
| `elfbot.image.tag` | Elfbot image tag | `focal-20240530` |
| `elfbot.resources.limits.cpu` | Elfbot CPU limit | `1` |
| `elfbot.resources.limits.memory` | Elfbot memory limit | `4Gi` |

### Health Checks

| Parameter | Description | Default |
|-----------|-------------|---------|
| `livenessProbe.tcpSocket.port` | Liveness probe port | `8080` |
| `livenessProbe.timeoutSeconds` | Liveness probe timeout | `1` |
| `livenessProbe.periodSeconds` | Liveness probe period | `10` |
| `readinessProbe.tcpSocket.port` | Readiness probe port | `8080` |
| `readinessProbe.timeoutSeconds` | Readiness probe timeout | `1` |
| `readinessProbe.periodSeconds` | Readiness probe period | `10` |
| `startupProbe.tcpSocket.port` | Startup probe port | `8080` |
| `startupProbe.failureThreshold` | Startup probe failure threshold | `30` |

## Examples

### Minimal Installation

```yaml
# values-minimal.yaml
persistence:
  config:
    size: 5Gi
  backup:
    enabled: false
  logs:
    size: 1Gi
elfbot:
  enabled: false
```

```bash
helm install filebrowser elfhosted/filebrowser -f values-minimal.yaml
```

### Production Installation

```yaml
# values-production.yaml
persistence:
  config:
    size: 50Gi
    storageClass: "fast-ssd"
  backup:
    size: 100Gi
    storageClass: "backup-storage"
  logs:
    size: 10Gi

resources:
  limits:
    cpu: "2"
    memory: 2Gi
  requests:
    cpu: "100m"
    memory: 256Mi

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/proxy-body-size: "100m"
  hosts:
    - host: files.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: filebrowser-tls
      hosts:
        - files.yourdomain.com

branding:
  name: "Company FileBrowser"
```

```bash
helm install filebrowser elfhosted/filebrowser -f values-production.yaml
```

### Custom Commands

```yaml
# values-custom-commands.yaml
filebrowser:
  allowedCommands: "ls,pwd,cd,mv,cp,find,grep,cat,tree,du,custom-script"

configMaps:
  filebrowserElfbotScript:
    data:
      elfbot: |
        #!/bin/bash
        # Custom elfbot script
        echo "Custom automation script"
```

## Troubleshooting

### Common Issues

1. **Pod stuck in Pending state**
   - Check if PersistentVolumeClaims are bound
   - Verify storage class exists
   - Check node resources

2. **Permission denied errors**
   - Verify security context settings
   - Check PVC permissions
   - Ensure fsGroup is set correctly

3. **Elfbot not working**
   - Check service account permissions
   - Verify RBAC settings
   - Check elfbot container logs

### Useful Commands

```bash
# Check pod status
kubectl get pods -l app.kubernetes.io/name=filebrowser

# Check logs
kubectl logs -l app.kubernetes.io/name=filebrowser -c filebrowser
kubectl logs -l app.kubernetes.io/name=filebrowser -c elfbot

# Check PVC status
kubectl get pvc

# Check service
kubectl get svc -l app.kubernetes.io/name=filebrowser

# Port forward for local access
kubectl port-forward svc/filebrowser 8080:8080
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This Helm chart is licensed under the Apache 2.0 License. See the LICENSE file for details.

## Support

- **Documentation**: [ElfHosted Docs](https://elfhosted.com/docs)
- **Issues**: [GitHub Issues](https://github.com/elfhosted/charts/issues)
- **Community**: [Discord](https://discord.gg/elfhosted)
