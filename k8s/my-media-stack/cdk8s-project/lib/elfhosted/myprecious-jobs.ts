import { Construct } from 'constructs';
import { KubeJob, KubeCronJob, Quantity } from '../../imports/k8s';

/**
 * Creates all Jobs and CronJobs for the myprecious chart
 * This includes backup jobs, maintenance jobs, and other scheduled tasks
 */

interface JobConfig {
  name: string;
  image: string;
  command: string[];
  args?: string[];
  schedule?: string; // For CronJobs
  restartPolicy?: string;
  serviceAccount?: string;
  annotations?: Record<string, string>;
  resources?: {
    requests: { cpu: string; memory: string };
    limits: { cpu: string; memory: string };
  };
}

const JOB_CONFIGS: JobConfig[] = [
  // Backup jobs
  {
    name: 'backup-config',
    image: 'ghcr.io/elfhosted/tooling:focal-20240530',
    command: ['/bin/bash'],
    args: ['-c', 'echo "Backing up config..." && tar -czf /backup/config-$(date +%Y%m%d-%H%M%S).tar.gz -C /config .'],
    schedule: '0 2 * * *', // Daily at 2 AM
    resources: {
      requests: { cpu: '100m', memory: '256Mi' },
      limits: { cpu: '500m', memory: '1Gi' },
    },
  },
  {
    name: 'backup-database',
    image: 'postgres:15-alpine',
    command: ['/bin/bash'],
    args: ['-c', 'pg_dumpall -h postgresql -U postgres > /backup/database-$(date +%Y%m%d-%H%M%S).sql'],
    schedule: '0 3 * * *', // Daily at 3 AM
    resources: {
      requests: { cpu: '100m', memory: '256Mi' },
      limits: { cpu: '500m', memory: '1Gi' },
    },
  },

  // Maintenance jobs
  {
    name: 'cleanup-logs',
    image: 'ghcr.io/elfhosted/tooling:focal-20240530',
    command: ['/bin/bash'],
    args: ['-c', 'find /logs -name "*.log" -mtime +7 -delete && echo "Log cleanup completed"'],
    schedule: '0 4 * * 0', // Weekly on Sunday at 4 AM
    resources: {
      requests: { cpu: '50m', memory: '128Mi' },
      limits: { cpu: '200m', memory: '512Mi' },
    },
  },
  {
    name: 'cleanup-transcode',
    image: 'ghcr.io/elfhosted/tooling:focal-20240530',
    command: ['/bin/bash'],
    args: ['-c', 'find /transcode -type f -mtime +1 -delete && echo "Transcode cleanup completed"'],
    schedule: '*/30 * * * *', // Every 30 minutes
    resources: {
      requests: { cpu: '50m', memory: '128Mi' },
      limits: { cpu: '200m', memory: '512Mi' },
    },
  },

  // Health check jobs
  {
    name: 'health-check-services',
    image: 'curlimages/curl:latest',
    command: ['/bin/sh'],
    args: ['-c', `
      services="plex:32400 radarr:7878 sonarr:8989 prowlarr:9696"
      for service in $services; do
        host=$(echo $service | cut -d: -f1)
        port=$(echo $service | cut -d: -f2)
        if curl -f -s http://$host:$port/ping > /dev/null; then
          echo "$service is healthy"
        else
          echo "$service is unhealthy"
        fi
      done
    `],
    schedule: '*/5 * * * *', // Every 5 minutes
    resources: {
      requests: { cpu: '10m', memory: '64Mi' },
      limits: { cpu: '100m', memory: '128Mi' },
    },
  },

  // Recyclarr sync job
  {
    name: 'recyclarr-sync',
    image: 'ghcr.io/recyclarr/recyclarr:latest',
    command: ['recyclarr'],
    args: ['sync'],
    schedule: '0 6 * * *', // Daily at 6 AM
    resources: {
      requests: { cpu: '100m', memory: '256Mi' },
      limits: { cpu: '500m', memory: '1Gi' },
    },
  },

  // Notifiarr sync job
  {
    name: 'notifiarr-sync',
    image: 'ghcr.io/elfhosted/notifiarr:latest',
    command: ['/app/notifiarr'],
    args: ['--sync'],
    schedule: '0 */6 * * *', // Every 6 hours
    resources: {
      requests: { cpu: '50m', memory: '128Mi' },
      limits: { cpu: '200m', memory: '512Mi' },
    },
  },

  // Media library scan jobs
  {
    name: 'plex-library-scan',
    image: 'ghcr.io/elfhosted/tooling:focal-20240530',
    command: ['/bin/bash'],
    args: ['-c', 'curl -X POST "http://plex:32400/library/sections/all/refresh?X-Plex-Token=$PLEX_TOKEN"'],
    schedule: '0 5 * * *', // Daily at 5 AM
    resources: {
      requests: { cpu: '10m', memory: '64Mi' },
      limits: { cpu: '100m', memory: '128Mi' },
    },
  },

  // Database maintenance
  {
    name: 'postgres-vacuum',
    image: 'postgres:15-alpine',
    command: ['/bin/bash'],
    args: ['-c', 'psql -h postgresql -U postgres -c "VACUUM ANALYZE;"'],
    schedule: '0 1 * * 0', // Weekly on Sunday at 1 AM
    resources: {
      requests: { cpu: '100m', memory: '256Mi' },
      limits: { cpu: '500m', memory: '1Gi' },
    },
  },
];

export function createAllMyPreciousJobs(scope: Construct) {
  JOB_CONFIGS.forEach(config => {
    if (config.schedule) {
      // Create CronJob
      new KubeCronJob(scope, `cronjob-${config.name}`, {
        metadata: {
          name: config.name,
          labels: {
            'app.kubernetes.io/name': config.name,
            'app.kubernetes.io/instance': 'myprecious',
            'app.kubernetes.io/component': 'cronjob',
          },
          ...(config.annotations && { annotations: config.annotations }),
        },
        spec: {
          schedule: config.schedule,
          jobTemplate: {
            spec: {
              template: {
                spec: {
                  restartPolicy: config.restartPolicy || 'OnFailure',
                  ...(config.serviceAccount && { serviceAccountName: config.serviceAccount }),
                  containers: [{
                    name: config.name,
                    image: config.image,
                    command: config.command,
                    ...(config.args && { args: config.args }),
                    ...(config.resources && {
                      resources: {
                        requests: {
                          cpu: Quantity.fromString(config.resources.requests.cpu),
                          memory: Quantity.fromString(config.resources.requests.memory),
                        },
                        limits: {
                          cpu: Quantity.fromString(config.resources.limits.cpu),
                          memory: Quantity.fromString(config.resources.limits.memory),
                        },
                      },
                    }),
                  }],
                },
              },
            },
          },
        },
      });
    } else {
      // Create regular Job
      new KubeJob(scope, `job-${config.name}`, {
        metadata: {
          name: config.name,
          labels: {
            'app.kubernetes.io/name': config.name,
            'app.kubernetes.io/instance': 'myprecious',
            'app.kubernetes.io/component': 'job',
          },
          ...(config.annotations && { annotations: config.annotations }),
        },
        spec: {
          template: {
            spec: {
              restartPolicy: config.restartPolicy || 'OnFailure',
              ...(config.serviceAccount && { serviceAccountName: config.serviceAccount }),
              containers: [{
                name: config.name,
                image: config.image,
                command: config.command,
                ...(config.args && { args: config.args }),
                ...(config.resources && {
                  resources: {
                    requests: {
                      cpu: Quantity.fromString(config.resources.requests.cpu),
                      memory: Quantity.fromString(config.resources.requests.memory),
                    },
                    limits: {
                      cpu: Quantity.fromString(config.resources.limits.cpu),
                      memory: Quantity.fromString(config.resources.limits.memory),
                    },
                  },
                }),
              }],
            },
          },
        },
      });
    }
  });
}

/**
 * Creates one-time initialization jobs
 */
export function createInitializationJobs(scope: Construct) {
  // Database initialization
  new KubeJob(scope, 'init-database', {
    metadata: {
      name: 'init-database',
      labels: {
        'app.kubernetes.io/name': 'init-database',
        'app.kubernetes.io/instance': 'myprecious',
        'app.kubernetes.io/component': 'init-job',
      },
      annotations: {
        'helm.sh/hook': 'post-install',
        'helm.sh/hook-weight': '1',
        'helm.sh/hook-delete-policy': 'hook-succeeded',
      },
    },
    spec: {
      template: {
        spec: {
          restartPolicy: 'OnFailure',
          containers: [{
            name: 'init-database',
            image: 'postgres:15-alpine',
            command: ['/bin/bash'],
            args: ['-c', `
              echo "Initializing databases..."
              createdb -h postgresql -U postgres riven || echo "Database riven already exists"
              createdb -h postgresql -U postgres notifiarr || echo "Database notifiarr already exists"
              echo "Database initialization completed"
            `],
            resources: {
              requests: {
                cpu: Quantity.fromString('100m'),
                memory: Quantity.fromString('256Mi'),
              },
              limits: {
                cpu: Quantity.fromString('500m'),
                memory: Quantity.fromString('1Gi'),
              },
            },
          }],
        },
      },
    },
  });

  // Configuration setup
  new KubeJob(scope, 'setup-config', {
    metadata: {
      name: 'setup-config',
      labels: {
        'app.kubernetes.io/name': 'setup-config',
        'app.kubernetes.io/instance': 'myprecious',
        'app.kubernetes.io/component': 'init-job',
      },
      annotations: {
        'helm.sh/hook': 'post-install',
        'helm.sh/hook-weight': '2',
        'helm.sh/hook-delete-policy': 'hook-succeeded',
      },
    },
    spec: {
      template: {
        spec: {
          restartPolicy: 'OnFailure',
          containers: [{
            name: 'setup-config',
            image: 'ghcr.io/elfhosted/tooling:focal-20240530',
            command: ['/bin/bash'],
            args: ['-c', `
              echo "Setting up initial configuration..."
              mkdir -p /config/{plex,radarr,sonarr,prowlarr,lidarr,readarr,jellyfin,overseerr,jellyseerr,tautulli,qbittorrent,filebrowser,wizarr,zurg,riven}
              chown -R 568:568 /config
              echo "Configuration setup completed"
            `],
            volumeMounts: [{
              name: 'config',
              mountPath: '/config',
            }],
            resources: {
              requests: {
                cpu: Quantity.fromString('50m'),
                memory: Quantity.fromString('128Mi'),
              },
              limits: {
                cpu: Quantity.fromString('200m'),
                memory: Quantity.fromString('512Mi'),
              },
            },
          }],
          volumes: [{
            name: 'config',
            persistentVolumeClaim: {
              claimName: 'config',
            },
          }],
        },
      },
    },
  });
} 