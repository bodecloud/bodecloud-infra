import { Construct } from 'constructs';
import { KubeStatefulSet, KubeDaemonSet, Quantity, IntOrString } from '../../imports/k8s';

/**
 * Creates all StatefulSets and DaemonSets for the myprecious chart
 * This includes databases, caches, and other stateful applications
 */

interface StatefulSetConfig {
  name: string;
  image: string;
  port: number;
  replicas?: number;
  env?: Record<string, string>;
  volumes?: Array<{
    name: string;
    mountPath: string;
    size?: string;
    storageClass?: string;
    accessMode?: string;
  }>;
  resources?: {
    requests: { cpu: string; memory: string };
    limits: { cpu: string; memory: string };
  };
  serviceName?: string;
}

interface DaemonSetConfig {
  name: string;
  image: string;
  env?: Record<string, string>;
  volumes?: Array<{
    name: string;
    mountPath: string;
    hostPath?: string;
    configMapName?: string;
  }>;
  resources?: {
    requests: { cpu: string; memory: string };
    limits: { cpu: string; memory: string };
  };
  hostNetwork?: boolean;
  privileged?: boolean;
}

const STATEFULSET_CONFIGS: StatefulSetConfig[] = [
  // PostgreSQL database
  {
    name: 'postgresql',
    image: 'postgres:15-alpine',
    port: 5432,
    replicas: 1,
    serviceName: 'postgresql',
    env: {
      'POSTGRES_DB': 'myprecious',
      'POSTGRES_USER': 'postgres',
      'POSTGRES_PASSWORD': 'your-postgres-password',
      'PGDATA': '/var/lib/postgresql/data/pgdata',
    },
    volumes: [
      {
        name: 'postgresql-data',
        mountPath: '/var/lib/postgresql/data',
        size: '20Gi',
        storageClass: 'ceph-block-ssd',
        accessMode: 'ReadWriteOnce',
      },
    ],
    resources: {
      requests: { cpu: '100m', memory: '256Mi' },
      limits: { cpu: '1', memory: '2Gi' },
    },
  },

  // Redis cache
  {
    name: 'redis',
    image: 'redis:7-alpine',
    port: 6379,
    replicas: 1,
    serviceName: 'redis',
    env: {
      'REDIS_PASSWORD': 'your-redis-password',
    },
    volumes: [
      {
        name: 'redis-data',
        mountPath: '/data',
        size: '5Gi',
        storageClass: 'ceph-block-ssd',
        accessMode: 'ReadWriteOnce',
      },
    ],
    resources: {
      requests: { cpu: '50m', memory: '128Mi' },
      limits: { cpu: '500m', memory: '1Gi' },
    },
  },

  // MongoDB database
  {
    name: 'mongodb',
    image: 'mongo:7',
    port: 27017,
    replicas: 1,
    serviceName: 'mongodb',
    env: {
      'MONGO_INITDB_ROOT_USERNAME': 'mongodb',
      'MONGO_INITDB_ROOT_PASSWORD': 'your-mongodb-password',
      'MONGO_INITDB_DATABASE': 'myprecious',
    },
    volumes: [
      {
        name: 'mongodb-data',
        mountPath: '/data/db',
        size: '20Gi',
        storageClass: 'ceph-block-ssd',
        accessMode: 'ReadWriteOnce',
      },
    ],
    resources: {
      requests: { cpu: '100m', memory: '256Mi' },
      limits: { cpu: '1', memory: '2Gi' },
    },
  },

  // Elasticsearch (for logging/search)
  {
    name: 'elasticsearch',
    image: 'elasticsearch:8.11.0',
    port: 9200,
    replicas: 1,
    serviceName: 'elasticsearch',
    env: {
      'discovery.type': 'single-node',
      'ES_JAVA_OPTS': '-Xms512m -Xmx512m',
      'xpack.security.enabled': 'false',
    },
    volumes: [
      {
        name: 'elasticsearch-data',
        mountPath: '/usr/share/elasticsearch/data',
        size: '10Gi',
        storageClass: 'ceph-block-ssd',
        accessMode: 'ReadWriteOnce',
      },
    ],
    resources: {
      requests: { cpu: '200m', memory: '1Gi' },
      limits: { cpu: '1', memory: '2Gi' },
    },
  },

  // InfluxDB (for metrics)
  {
    name: 'influxdb',
    image: 'influxdb:2.7-alpine',
    port: 8086,
    replicas: 1,
    serviceName: 'influxdb',
    env: {
      'DOCKER_INFLUXDB_INIT_MODE': 'setup',
      'DOCKER_INFLUXDB_INIT_USERNAME': 'admin',
      'DOCKER_INFLUXDB_INIT_PASSWORD': 'your-influxdb-password',
      'DOCKER_INFLUXDB_INIT_ORG': 'elfhosted',
      'DOCKER_INFLUXDB_INIT_BUCKET': 'metrics',
    },
    volumes: [
      {
        name: 'influxdb-data',
        mountPath: '/var/lib/influxdb2',
        size: '10Gi',
        storageClass: 'ceph-block-ssd',
        accessMode: 'ReadWriteOnce',
      },
    ],
    resources: {
      requests: { cpu: '100m', memory: '256Mi' },
      limits: { cpu: '500m', memory: '1Gi' },
    },
  },
];

const DAEMONSET_CONFIGS: DaemonSetConfig[] = [
  // Node exporter for monitoring
  {
    name: 'node-exporter',
    image: 'prom/node-exporter:latest',
    hostNetwork: true,
    volumes: [
      {
        name: 'proc',
        mountPath: '/host/proc',
        hostPath: '/proc',
      },
      {
        name: 'sys',
        mountPath: '/host/sys',
        hostPath: '/sys',
      },
      {
        name: 'root',
        mountPath: '/rootfs',
        hostPath: '/',
      },
    ],
    resources: {
      requests: { cpu: '10m', memory: '32Mi' },
      limits: { cpu: '100m', memory: '128Mi' },
    },
  },

  // Filebeat for log collection
  {
    name: 'filebeat',
    image: 'elastic/filebeat:8.11.0',
    volumes: [
      {
        name: 'varlog',
        mountPath: '/var/log',
        hostPath: '/var/log',
      },
      {
        name: 'varlibdockercontainers',
        mountPath: '/var/lib/docker/containers',
        hostPath: '/var/lib/docker/containers',
      },
      {
        name: 'filebeat-config',
        mountPath: '/usr/share/filebeat/filebeat.yml',
        configMapName: 'filebeat-config',
      },
    ],
    resources: {
      requests: { cpu: '50m', memory: '128Mi' },
      limits: { cpu: '200m', memory: '512Mi' },
    },
  },

  // Fluent Bit for log processing
  {
    name: 'fluent-bit',
    image: 'fluent/fluent-bit:latest',
    volumes: [
      {
        name: 'varlog',
        mountPath: '/var/log',
        hostPath: '/var/log',
      },
      {
        name: 'fluent-bit-config',
        mountPath: '/fluent-bit/etc/fluent-bit.conf',
        configMapName: 'fluent-bit-config',
      },
    ],
    resources: {
      requests: { cpu: '50m', memory: '128Mi' },
      limits: { cpu: '200m', memory: '512Mi' },
    },
  },
];

export function createAllMyPreciousStatefulSets(scope: Construct) {
  STATEFULSET_CONFIGS.forEach(config => {
    const volumeClaimTemplates = config.volumes?.filter(vol => vol.size).map(vol => ({
      metadata: {
        name: vol.name,
      },
      spec: {
        accessModes: [vol.accessMode || 'ReadWriteOnce'],
        storageClassName: vol.storageClass || 'ceph-block-ssd',
        resources: {
          requests: {
            storage: Quantity.fromString(vol.size!),
          },
        },
      },
    })) || [];

    const volumes = config.volumes?.filter(vol => !vol.size).map(vol => ({
      name: vol.name,
      emptyDir: {},
    })) || [];

    const volumeMounts = config.volumes?.map(vol => ({
      name: vol.name,
      mountPath: vol.mountPath,
    })) || [];

    new KubeStatefulSet(scope, `statefulset-${config.name}`, {
      metadata: {
        name: config.name,
        labels: {
          'app.elfhosted.com/name': config.name,
          'app.kubernetes.io/name': config.name,
          'app.kubernetes.io/instance': 'myprecious',
          'app.kubernetes.io/component': 'database',
        },
      },
      spec: {
        serviceName: config.serviceName || config.name,
        replicas: config.replicas || 1,
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
                name: 'main',
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
                tcpSocket: {
                  port: IntOrString.fromNumber(config.port),
                },
                initialDelaySeconds: 30,
                periodSeconds: 30,
                timeoutSeconds: 10,
                failureThreshold: 3,
              },
              readinessProbe: {
                tcpSocket: {
                  port: IntOrString.fromNumber(config.port),
                },
                initialDelaySeconds: 10,
                periodSeconds: 10,
                timeoutSeconds: 5,
                failureThreshold: 3,
              },
            }],
          },
        },
        volumeClaimTemplates,
      },
    });
  });
}

export function createAllMyPreciousDaemonSets(scope: Construct) {
  DAEMONSET_CONFIGS.forEach(config => {
    const volumes = config.volumes?.map(vol => {
      if (vol.hostPath) {
        return {
          name: vol.name,
          hostPath: {
            path: vol.hostPath,
          },
        };
      } else if (vol.configMapName) {
        return {
          name: vol.name,
          configMap: {
            name: vol.configMapName,
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
      ...(vol.hostPath && { readOnly: true }),
    })) || [];

    new KubeDaemonSet(scope, `daemonset-${config.name}`, {
      metadata: {
        name: config.name,
        labels: {
          'app.elfhosted.com/name': config.name,
          'app.kubernetes.io/name': config.name,
          'app.kubernetes.io/instance': 'myprecious',
          'app.kubernetes.io/component': 'monitoring',
        },
      },
      spec: {
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
            ...(config.hostNetwork && { hostNetwork: true }),
            volumes,
            containers: [{
              name: config.name,
              image: config.image,
              ...(config.env && {
                env: Object.entries(config.env).map(([key, value]) => ({
                  name: key,
                  value,
                })),
              }),
              volumeMounts,
              ...(config.privileged && {
                securityContext: {
                  privileged: true,
                },
              }),
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
            tolerations: [{
              operator: 'Exists',
              effect: 'NoSchedule',
            }],
          },
        },
      },
    });
  });
} 