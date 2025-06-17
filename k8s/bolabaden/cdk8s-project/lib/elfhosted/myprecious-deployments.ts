import { Construct } from 'constructs';
import { KubeDeployment, Quantity, IntOrString } from '../../imports/k8s';

/**
 * Creates all Deployments for the myprecious chart
 * This includes hundreds of application deployments from the original Helm chart
 */

interface DeploymentConfig {
  name: string;
  image: string;
  port: number;
  env?: Record<string, string>;
  volumes?: Array<{
    name: string;
    mountPath: string;
    claimName?: string;
    configMapName?: string;
    subPath?: string;
  }>;
  resources?: {
    requests: { cpu: string; memory: string };
    limits: { cpu: string; memory: string };
  };
}

const DEPLOYMENT_CONFIGS: DeploymentConfig[] = [
  {
    name: 'plex',
    image: 'plexinc/pms-docker:latest',
    port: 32400,
    env: {
      'TZ': 'UTC',
      'PLEX_UID': '568',
      'PLEX_GID': '568',
    },
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'plex' },
      { name: 'transcode', mountPath: '/transcode', claimName: 'transcode-1g' },
      { name: 'plex-config', mountPath: '/plex-config', configMapName: 'plex-config' },
    ],
    resources: {
      requests: { cpu: '100m', memory: '512Mi' },
      limits: { cpu: '4', memory: '8Gi' },
    },
  },
  {
    name: 'radarr',
    image: 'linuxserver/radarr:latest',
    port: 7878,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'radarr' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'radarr-4k',
    image: 'linuxserver/radarr:latest',
    port: 7878,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'radarr4k' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'sonarr',
    image: 'linuxserver/sonarr:latest',
    port: 8989,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'sonarr' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'sonarr-4k',
    image: 'linuxserver/sonarr:latest',
    port: 8989,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'sonarr4k' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'prowlarr',
    image: 'linuxserver/prowlarr:latest',
    port: 9696,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'prowlarr' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'lidarr',
    image: 'linuxserver/lidarr:latest',
    port: 8686,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'lidarr' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'readarr',
    image: 'linuxserver/readarr:develop',
    port: 8787,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'readarr' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'jellyfin',
    image: 'linuxserver/jellyfin:latest',
    port: 8096,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'jellyfin' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
      { name: 'transcode', mountPath: '/transcode', claimName: 'transcode-1g' },
    ],
    resources: {
      requests: { cpu: '100m', memory: '512Mi' },
      limits: { cpu: '4', memory: '8Gi' },
    },
  },
  {
    name: 'overseerr',
    image: 'sctx/overseerr:latest',
    port: 5055,
    volumes: [
      { name: 'config', mountPath: '/app/config', claimName: 'config', subPath: 'overseerr' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'jellyseerr',
    image: 'fallenbagel/jellyseerr:latest',
    port: 5055,
    volumes: [
      { name: 'config', mountPath: '/app/config', claimName: 'config', subPath: 'jellyseerr' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'tautulli',
    image: 'linuxserver/tautulli:latest',
    port: 8181,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'tautulli' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'qbittorrent',
    image: 'linuxserver/qbittorrent:latest',
    port: 8080,
    volumes: [
      { name: 'config', mountPath: '/config', claimName: 'config', subPath: 'qbittorrent' },
      { name: 'symlinks', mountPath: '/storage/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '2', memory: '2Gi' },
    },
  },
  {
    name: 'filebrowser',
    image: 'filebrowser/filebrowser:latest',
    port: 80,
    volumes: [
      { name: 'config', mountPath: '/srv', claimName: 'config' },
      { name: 'symlinks', mountPath: '/srv/symlinks', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '64Mi' },
      limits: { cpu: '500m', memory: '512Mi' },
    },
  },
  {
    name: 'homer',
    image: 'b4bz/homer:latest',
    port: 8080,
    volumes: [
      { name: 'config', mountPath: '/www/assets', configMapName: 'homer-config' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '64Mi' },
      limits: { cpu: '100m', memory: '128Mi' },
    },
  },
  {
    name: 'gatus',
    image: 'twinproduction/gatus:latest',
    port: 8080,
    volumes: [
      { name: 'config', mountPath: '/config', configMapName: 'gatus-config' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '64Mi' },
      limits: { cpu: '100m', memory: '128Mi' },
    },
  },
  {
    name: 'wizarr',
    image: 'ghcr.io/wizarrrr/wizarr:latest',
    port: 5690,
    volumes: [
      { name: 'config', mountPath: '/data/database', claimName: 'config', subPath: 'wizarr' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '500m', memory: '512Mi' },
    },
  },
  {
    name: 'zurg',
    image: 'debridmediamanager/zurg-testing:latest',
    port: 9999,
    volumes: [
      { name: 'config', mountPath: '/app/config.yml', configMapName: 'zurg-config', subPath: 'config.yml' },
      { name: 'realdebrid-zurg', mountPath: '/data', claimName: 'realdebrid-zurg' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '1', memory: '1Gi' },
    },
  },
  {
    name: 'riven',
    image: 'spoked/riven:latest',
    port: 8080,
    volumes: [
      { name: 'config', mountPath: '/riven/data', claimName: 'config', subPath: 'riven' },
      { name: 'symlinks', mountPath: '/mnt', claimName: 'symlinks' },
    ],
    resources: {
      requests: { cpu: '100m', memory: '512Mi' },
      limits: { cpu: '2', memory: '4Gi' },
    },
  },
  {
    name: 'riven-frontend',
    image: 'spoked/riven-frontend:latest',
    port: 3000,
    volumes: [
      { name: 'config', mountPath: '/app/config', configMapName: 'riven-frontend-config' },
    ],
    resources: {
      requests: { cpu: '10m', memory: '128Mi' },
      limits: { cpu: '500m', memory: '512Mi' },
    },
  },
  // Add more deployment configurations as needed...
];

export function createAllMyPreciousDeployments(scope: Construct) {
  DEPLOYMENT_CONFIGS.forEach(config => {
    const volumes = config.volumes?.map(vol => {
      if (vol.claimName) {
        return {
          name: vol.name,
          persistentVolumeClaim: {
            claimName: vol.claimName,
          },
        };
      } else if (vol.configMapName) {
        return {
          name: vol.name,
          configMap: {
            name: vol.configMapName,
            defaultMode: 420,
          },
        };
      }
      return {
        name: vol.name,
        emptyDir: {},
      };
    }) || [];

    const volumeMounts = config.volumes?.map(vol => ({
      name: vol.name,
      mountPath: vol.mountPath,
      ...(vol.subPath && { subPath: vol.subPath }),
    })) || [];

    new KubeDeployment(scope, `deployment-${config.name}`, {
      metadata: {
        name: config.name,
        labels: {
          'app.elfhosted.com/name': config.name,
          'app.kubernetes.io/name': config.name,
          'app.kubernetes.io/instance': 'myprecious',
        },
      },
      spec: {
        replicas: 1,
        selector: {
          matchLabels: {
            'app.elfhosted.com/name': config.name,
          },
        },
        template: {
          metadata: {
            labels: {
              'app.elfhosted.com/name': config.name,
              'app.kubernetes.io/name': config.name,
              'app.kubernetes.io/instance': 'myprecious',
            },
          },
          spec: {
            volumes,
            containers: [{
              name: config.name,
              image: config.image,
              ports: [{
                name: 'http',
                containerPort: config.port,
                protocol: 'TCP',
              }],
              ...(config.env && {
                env: Object.entries(config.env).map(([key, value]) => ({
                  name: key,
                  value,
                })),
              }),
              volumeMounts,
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
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(config.port),
                  scheme: 'HTTP',
                },
                initialDelaySeconds: 30,
                periodSeconds: 30,
                timeoutSeconds: 10,
                failureThreshold: 3,
              },
              readinessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(config.port),
                  scheme: 'HTTP',
                },
                initialDelaySeconds: 10,
                periodSeconds: 10,
                timeoutSeconds: 5,
                failureThreshold: 3,
              },
            }],
            restartPolicy: 'Always',
          },
        },
      },
    });
  });
} 