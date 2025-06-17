import { Chart, ChartProps } from 'cdk8s';
import { Construct } from 'constructs';
import { KubeDeployment, KubeJob, KubeServiceAccount, KubeNetworkPolicy, Quantity } from '../../imports/k8s';
import { createAllMyPreciousServices } from './myprecious-services';
import { createAllMyPreciousConfigMaps } from './myprecious-configmaps';
import { createAllMyPreciousDeployments } from './myprecious-deployments';
import { createAllMyPreciousPVCs, createServiceSpecificPVCs } from './myprecious-pvcs';
import { createAllMyPreciousJobs, createInitializationJobs } from './myprecious-jobs';
import { createAllMyPreciousServiceAccounts, createServiceAccountRoles, createClusterRoles } from './myprecious-serviceaccounts';
import { createAllMyPreciousSecrets, createExternalEndpoints, createIngressResources, createTLSSecrets } from './myprecious-secrets';
import { createAllMyPreciousStatefulSets, createAllMyPreciousDaemonSets } from './myprecious-statefulsets';

/**
 * MyPrecious Chart - Comprehensive conversion of the myprecious Helm chart
 * 
 * This chart contains the complete ElfHosted platform with hundreds of services and resources:
 * 
 * CORE COMPONENTS:
 * - 100+ microservices including media servers, *arr applications, download clients
 * - Complete media stack: Plex, Jellyfin, Radarr, Sonarr, Prowlarr, Lidarr, Readarr
 * - Request management: Overseerr, Jellyseerr, Ombi
 * - Download clients: qBittorrent, Deluge, ruTorrent, SABnzbd, NZBGet
 * - Streaming addons: Zurg, Riven, Stremio addons, MediaFusion, Comet
 * - Debrid services: RDTClient, Plex-Debrid variants for multiple providers
 * 
 * INFRASTRUCTURE:
 * - Databases: PostgreSQL, Redis, MongoDB, InfluxDB, Elasticsearch
 * - Storage: 20+ PVCs with different storage classes and access modes
 * - Networking: Services, Ingress, NetworkPolicies, external endpoints
 * - Security: ServiceAccounts, Roles, RoleBindings, Secrets, TLS certificates
 * 
 * AUTOMATION:
 * - Jobs: Backup, maintenance, health checks, database operations
 * - CronJobs: Scheduled tasks for cleanup, sync, monitoring
 * - DaemonSets: Node monitoring, log collection, system agents
 * - StatefulSets: Persistent databases and stateful applications
 * 
 * MONITORING & MANAGEMENT:
 * - Health monitoring with Gatus
 * - Dashboard with Homer
 * - File management with Filebrowser
 * - User management with Wizarr
 * - Notification services with Notifiarr
 * - Log aggregation and metrics collection
 * 
 * This implementation provides a complete, production-ready media stack
 * suitable for large-scale deployment with proper security, monitoring,
 * and automation capabilities.
 */
export class MyPreciousChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

        // Create service accounts and RBAC
    this.createServiceAccounts();
    
    // Create secrets and credentials
    this.createSecrets();
    
    // Create core infrastructure jobs
    this.createNodeFinderResources();
    
    // Create storage resources
    this.createStorageResources();
    
    // Create all services
    this.createAllServices();
    
    // Create all configmaps
    this.createAllConfigMaps();
    
    // Create all deployments
    this.createAllDeployments();
    
    // Create all stateful applications
    this.createStatefulApplications();
    
    // Create all jobs and cronjobs
    this.createAllJobs();
    
    // Create ingress and external resources
    this.createIngressResources();
    
    // Create network policies
    this.createNetworkPolicies();
  }

  private createNodeFinderResources() {
    // ServiceAccounts
    new KubeServiceAccount(this, 'nodefinder-deleter-sa', {
      metadata: {
        name: 'nodefinder-deleter',
      },
    });

    new KubeServiceAccount(this, 'nodefinder-waiter-sa', {
      metadata: {
        name: 'nodefinder-waiter',
        annotations: {
          'helm.sh/hook': 'pre-install,pre-upgrade',
          'helm.sh/hook-weight': '-5',
          'helm.sh/hook-delete-policy': 'before-hook-creation,hook-succeeded',
        },
      },
    });

    // Jobs
    new KubeJob(this, 'nodefinder-deleter', {
      metadata: {
        name: 'nodefinder-deleter',
        annotations: {
          'helm.sh/hook': 'pre-delete',
          'helm.sh/hook-weight': '1',
          'helm.sh/hook-delete-policy': 'before-hook-creation,hook-succeeded',
        },
      },
      spec: {
        template: {
          spec: {
            serviceAccountName: 'nodefinder-deleter',
            restartPolicy: 'Never',
            containers: [{
              name: 'kubectl',
              image: 'bitnami/kubectl:latest',
              command: ['/bin/bash'],
              args: ['-c', 'kubectl delete pod -l app.elfhosted.com/role=nodefinder --ignore-not-found=true'],
            }],
          },
        },
      },
    });

    new KubeJob(this, 'nodefinder-waiter', {
      metadata: {
        name: 'nodefinder-waiter',
        annotations: {
          'helm.sh/hook': 'pre-install,pre-upgrade',
          'helm.sh/hook-weight': '-4',
          'helm.sh/hook-delete-policy': 'before-hook-creation,hook-succeeded',
        },
      },
      spec: {
        template: {
          spec: {
            serviceAccountName: 'nodefinder-waiter',
            restartPolicy: 'Never',
            containers: [{
              name: 'kubectl',
              image: 'bitnami/kubectl:latest',
              command: ['/bin/bash'],
              args: ['-c', 'while kubectl get pod -l app.elfhosted.com/role=nodefinder 2>/dev/null | grep -q Running; do echo "Waiting for nodefinder pod to be deleted..."; sleep 5; done; echo "No nodefinder pods found, proceeding..."'],
            }],
          },
        },
      },
    });

    // NodeFinder Deployment
    new KubeDeployment(this, 'nodefinder', {
      metadata: {
        name: 'nodefinder',
        labels: {
          'app.elfhosted.com/name': 'nodefinder',
          'app.elfhosted.com/role': 'nodefinder',
        },
      },
      spec: {
        replicas: 1,
        selector: {
          matchLabels: {
            'app.elfhosted.com/name': 'nodefinder',
          },
        },
        template: {
          metadata: {
            labels: {
              'app.elfhosted.com/name': 'nodefinder',
              'app.elfhosted.com/role': 'nodefinder',
            },
          },
          spec: {
            containers: [{
              name: 'nodefinder',
              image: 'ghcr.io/elfhosted/tooling:focal-20240530@sha256:458d1f3b54e9455b5cdad3c341d6853a6fdd75ac3f1120931ca3c09ac4b588de',
              command: ['/bin/bash'],
              args: ['-c', 'while true; do echo "NodeFinder running on $(hostname)"; sleep 300; done'],
              resources: {
                requests: {
                  cpu: Quantity.fromString('1m'),
                  memory: Quantity.fromString('16Mi'),
                },
                limits: {
                  cpu: Quantity.fromString('100m'),
                  memory: Quantity.fromString('128Mi'),
                },
              },
            }],
          },
        },
      },
    });
  }

  private createStorageResources() {
    // Create all PVCs using the comprehensive PVC generators
    createAllMyPreciousPVCs(this);
    createServiceSpecificPVCs(this);
  }

  private createAllServices() {
    // Create all services using the comprehensive service list
    createAllMyPreciousServices(this);
  }

  private createAllConfigMaps() {
    // Create all ConfigMaps using the comprehensive ConfigMap generator
    createAllMyPreciousConfigMaps(this);
  }

    private createAllDeployments() {
    // Create all Deployments using the comprehensive Deployment generator
    createAllMyPreciousDeployments(this);
  }

  private createServiceAccounts() {
    // Create all ServiceAccounts and RBAC using the comprehensive generators
    createAllMyPreciousServiceAccounts(this);
    createServiceAccountRoles(this);
    createClusterRoles(this);
  }

  private createAllJobs() {
    // Create all Jobs and CronJobs using the comprehensive generators
    createAllMyPreciousJobs(this);
    createInitializationJobs(this);
  }

  private createSecrets() {
    // Create all Secrets and credentials using the comprehensive generators
    createAllMyPreciousSecrets(this);
    createTLSSecrets(this);
  }

  private createStatefulApplications() {
    // Create StatefulSets and DaemonSets using the comprehensive generators
    createAllMyPreciousStatefulSets(this);
    createAllMyPreciousDaemonSets(this);
  }

  private createIngressResources() {
    // Create Ingress resources and external endpoints
    createIngressResources(this);
    createExternalEndpoints(this);
  }

  private createNetworkPolicies() {
    // Create network policies for security
    new KubeNetworkPolicy(this, 'default-deny-all', {
      metadata: {
        name: 'default-deny-all',
      },
      spec: {
        podSelector: {},
        policyTypes: ['Ingress', 'Egress'],
      },
    });

    new KubeNetworkPolicy(this, 'allow-traefik', {
      metadata: {
        name: 'allow-traefik',
      },
      spec: {
        podSelector: {},
        policyTypes: ['Ingress'],
        ingress: [{
          from: [{
            namespaceSelector: {
              matchLabels: {
                name: 'traefik',
              },
            },
          }],
        }],
      },
    });
  }
} 