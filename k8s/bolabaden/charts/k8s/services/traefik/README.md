# Traefik Kubernetes Setup

This directory contains the necessary Kubernetes manifests to deploy Traefik as an ingress controller in a Kubernetes cluster.

## Components

1. **RBAC Resources**:
   - `00-role.yml`: ClusterRole with permissions for Traefik to access Kubernetes resources
   - `00-account.yml`: ServiceAccount for Traefik
   - `01-role-binding.yml`: ClusterRoleBinding to bind the role to the service account

2. **Traefik Configuration**:
   - `traefik-env-configmap.yml`: ConfigMap with environment variables for Traefik

3. **Traefik Resources**:
   - `02-traefik.yml`: Deployment for Traefik
   - `02-traefik-services.yml`: Services for Traefik dashboard, web, and websecure endpoints

4. **Demo Application**:
   - `03-whoami.yml`: Deployment for the whoami application
   - `03-whoami-services.yml`: Service for the whoami application
   - `04-whoami-ingress.yml`: Ingress for the whoami application

## Deployment

To deploy all resources, run:

```bash
./deploy.sh
```

## Accessing Services

- **Traefik Dashboard**: http://localhost:8080
- **Whoami Application**: http://whoami.bolabaden.org (or the domain configured in DOMAIN environment variable)

## Configuration

The setup closely mirrors the docker-compose configuration:

1. **Volumes**:
   - Mounts the same directories from the host system as the docker-compose version
   - Uses hostPath volumes to directly access files on the host

2. **Environment Variables**:
   - All environment variables from docker-compose are provided via ConfigMap
   - Variables are referenced in the Deployment via valueFrom/configMapKeyRef

3. **Command**:
   - Uses the same command as the docker-compose setup: `--configFile=/config/traefik.yaml`

## Notes

- The configuration uses direct host paths for volume mounts, matching the docker-compose configuration
- Environment variables are stored in a ConfigMap for maintainability
- The deployment checks for and starts the K3s service if needed 