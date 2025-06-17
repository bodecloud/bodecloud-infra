# Gatus Helm Chart

A Helm chart for deploying Gatus - A developer-oriented health dashboard that gives you the ability to monitor your services using HTTP, ICMP, TCP, and even custom conditions.

## Overview

Gatus is a powerful health monitoring dashboard that allows you to monitor the uptime and performance of your applications and services. This chart provides:

- **Health Monitoring** - HTTP, ICMP, TCP, and custom conditions
- **Alerting** - Email, Discord, Slack, Telegram notifications
- **Dashboard** - Clean web interface for monitoring status
- **Persistent Storage** - Optional data persistence
- **SMTP Integration** - Email notifications for alerts
- **Customizable Configuration** - Flexible endpoint monitoring

## Features

- 🎯 **Multi-Protocol Monitoring** - HTTP, ICMP, TCP support
- 📧 **Multiple Alert Channels** - Email, Discord, Slack, Telegram
- 📊 **Metrics** - Built-in metrics endpoint
- 🔐 **Security** - Non-root containers, security contexts
- 💾 **Persistence** - Optional SQLite data persistence
- ⚙️ **Configurable** - Extensive configuration options
- 🔍 **Health Checks** - Comprehensive probes

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support (if persistence enabled)

## Installation

### Quick Start

```bash
# Add the helm repository
helm repo add elfhosted https://charts.elfhosted.com
helm repo update

# Install with default values
helm install gatus elfhosted/gatus
```

### Custom Installation

```bash
# Install with custom values
helm install gatus elfhosted/gatus \
  --set persistence.config.size=5Gi \
  --set smtp.enabled=true \
  --set smtp.host=smtp.gmail.com
```

### Using a values file

```bash
# Install with custom values file
helm install gatus elfhosted/gatus -f my-values.yaml
```

## Configuration

The following table lists the configurable parameters of the Gatus chart and their default values.

### Global Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.imageRegistry` | Global Docker image registry | `ghcr.io/elfhosted` |
| `global.imagePullSecrets` | Global Docker registry secret names | `[]` |
| `global.storageClass` | Global StorageClass for Persistent Volume(s) | `local-path` |

### Image Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Gatus image repository | `gatus` |
| `image.tag` | Gatus image tag | `5.17.0` |
| `image.digest` | Gatus image digest | `sha256:e2752cf7e...` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |

### Application Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of Gatus replicas | `1` |
| `env.GATUS_CONFIG_PATH` | Path to Gatus configuration file | `/config/config.yaml` |
| `env.SMTP_FROM` | Default SMTP from address | `health@elfhosted.com` |
| `env.SMTP_PORT` | Default SMTP port | `587` |

### Service Account

| Parameter | Description | Default |
|-----------|-------------|---------|
| `serviceAccount.create` | Specifies whether a service account should be created | `false` |
| `serviceAccount.name` | The name of the service account | `default` |
| `serviceAccount.automount` | Automount service account token | `false` |
| `serviceAccount.annotations` | Annotations to add to the service account | `{}` |

### Security Context

| Parameter | Description | Default |
|-----------|-------------|---------|
| `securityContext.fsGroup` | Group ID for the container | `568` |
| `containerSecurityContext.readOnlyRootFilesystem` | Mount root filesystem as read-only | `true` |

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
| `ingress.hosts` | An array with hosts and paths | `[{host: "gatus.local", paths: [{path: "/", pathType: "Prefix"}]}]` |
| `ingress.tls` | An array with the tls configuration | `[]` |

### Persistence Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `persistence.config.enabled` | Enable persistence for configuration | `true` |
| `persistence.config.size` | Size of persistent volume claim | `1Gi` |
| `persistence.config.storageClass` | Type of persistent volume claim | `""` |
| `persistence.config.accessModes` | Persistent volume access modes | `["ReadWriteOnce"]` |

### ConfigMap Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `configMap.enabled` | Enable ConfigMap creation | `true` |
| `configMap.data` | Configuration data | See values.yaml |

### Resource Management

| Parameter | Description | Default |
|-----------|-------------|---------|
| `resources.limits.cpu` | CPU limit | `1` |
| `resources.limits.memory` | Memory limit | `128Mi` |
| `resources.requests.cpu` | CPU request | `0` |
| `resources.requests.memory` | Memory request | `20Mi` |

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

### SMTP Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `smtp.enabled` | Enable SMTP configuration via values | `false` |
| `smtp.host` | SMTP server hostname | `""` |
| `smtp.port` | SMTP server port | `587` |
| `smtp.username` | SMTP username | `""` |
| `smtp.password` | SMTP password | `""` |
| `smtp.from` | SMTP from address | `health@elfhosted.com` |
| `smtp.to` | SMTP to address | `""` |

### Gatus Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `gatus.metrics.enabled` | Enable metrics endpoint | `true` |
| `gatus.storage.type` | Storage type (memory, sqlite, postgres) | `memory` |
| `gatus.ui.title` | Dashboard title | `ElfHosted Health Dashboard` |
| `gatus.ui.description` | Dashboard description | `Monitor your ElfHosted services` |
| `gatus.endpoints.defaultInterval` | Default monitoring interval | `30s` |
| `gatus.endpoints.defaultTimeout` | Default request timeout | `10s` |

## Examples

### Minimal Installation

```yaml
# values-minimal.yaml
persistence:
  config:
    enabled: false

configMap:
  data:
    config.yaml: |
      metrics: true
      storage:
        type: memory
      endpoints:
        - name: "Google"
          url: "https://google.com"
          interval: 30s
          conditions:
            - "[STATUS] == 200"
```

```bash
helm install gatus elfhosted/gatus -f values-minimal.yaml
```

### Production with SMTP

```yaml
# values-production.yaml
persistence:
  config:
    enabled: true
    size: 5Gi
    storageClass: "fast-ssd"

smtp:
  enabled: true
  host: "smtp.gmail.com"
  port: 587
  username: "monitoring@company.com"
  password: "app-password"
  from: "monitoring@company.com"
  to: "alerts@company.com"

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: gatus.company.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: gatus-tls
      hosts:
        - gatus.company.com

gatus:
  storage:
    type: "sqlite"
  ui:
    title: "Company Health Dashboard"
    description: "Monitor all company services"

configMap:
  data:
    config.yaml: |
      metrics: true
      storage:
        type: sqlite
        path: /data/data.db
      
      alerting:
        email:
          from: "${SMTP_FROM}"
          host: "${SMTP_HOST}"
          port: ${SMTP_PORT}
          username: "${SMTP_USERNAME}"
          password: "${SMTP_PASSWORD}"
          to: "${SMTP_TO}"
          default-alert:
            failure-threshold: 3
            success-threshold: 2
            send-on-resolved: true
      
      endpoints:
        - name: "Main Website"
          url: "https://company.com"
          interval: 30s
          conditions:
            - "[STATUS] == 200"
            - "[RESPONSE_TIME] < 1000"
          alerts:
            - type: email
        
        - name: "API Service"
          url: "https://api.company.com/health"
          interval: 15s
          conditions:
            - "[STATUS] == 200"
            - "[BODY].status == UP"
          alerts:
            - type: email
```

```bash
helm install gatus elfhosted/gatus -f values-production.yaml
```

### Multi-Service Monitoring

```yaml
# values-multi-service.yaml
configMap:
  data:
    config.yaml: |
      metrics: true
      storage:
        type: memory
      
      endpoints:
        # Web Services
        - name: "Frontend"
          url: "https://app.company.com"
          interval: 30s
          conditions:
            - "[STATUS] == 200"
            - "[RESPONSE_TIME] < 2000"
        
        - name: "API Gateway"
          url: "https://api.company.com/health"
          interval: 15s
          conditions:
            - "[STATUS] == 200"
            - "[BODY].status == healthy"
        
        # Database Health
        - name: "Database"
          url: "tcp://db.company.com:5432"
          interval: 60s
          conditions:
            - "[CONNECTED] == true"
        
        # Internal Services
        - name: "Redis Cache"
          url: "tcp://redis.company.com:6379"
          interval: 30s
          conditions:
            - "[CONNECTED] == true"
        
        # External Dependencies
        - name: "Payment Gateway"
          url: "https://api.stripe.com/v1"
          interval: 120s
          conditions:
            - "[STATUS] == 200"
```

## Troubleshooting

### Common Issues

1. **Pod stuck in Pending state**
   - Check if PersistentVolumeClaim is bound (if persistence enabled)
   - Verify storage class exists
   - Check node resources

2. **Configuration not loading**
   - Verify ConfigMap is created correctly
   - Check configuration syntax in YAML
   - Review container logs

3. **SMTP alerts not working**
   - Verify SMTP credentials are correct
   - Check network connectivity to SMTP server
   - Review Gatus logs for SMTP errors

### Useful Commands

```bash
# Check pod status
kubectl get pods -l app.kubernetes.io/name=gatus

# Check logs
kubectl logs -l app.kubernetes.io/name=gatus

# Check ConfigMap
kubectl get configmap -l app.kubernetes.io/name=gatus

# Check service
kubectl get svc -l app.kubernetes.io/name=gatus

# Port forward for local access
kubectl port-forward svc/gatus 8080:8080
```

### Debug Configuration

```bash
# View current configuration
kubectl get configmap <release-name>-gatus-config -o yaml

# Test configuration locally
docker run --rm -v $(pwd)/config.yaml:/config/config.yaml:ro \
  twinproduction/gatus:latest
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

- **Documentation**: [Gatus Documentation](https://gatus.io)
- **ElfHosted Docs**: [ElfHosted Documentation](https://elfhosted.com/docs)
- **Issues**: [GitHub Issues](https://github.com/elfhosted/charts/issues)
- **Community**: [Discord](https://discord.gg/elfhosted) 