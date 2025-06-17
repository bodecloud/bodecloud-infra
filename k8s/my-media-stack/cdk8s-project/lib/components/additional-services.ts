import { Construct } from 'constructs';
import { KubeDeployment, Quantity } from '../../imports/k8s';

/**
 * Creates additional service deployments (SearXNG, Dozzle, Speedtest, etc.)
 */
export function createAdditionalServices(scope: Construct): void {
  // SearXNG deployment
  new KubeDeployment(scope, 'searxng-deployment', {
    metadata: {
      name: 'searxng',
      namespace: 'my-media-stack',
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          app: 'searxng',
        },
      },
      strategy: {
        type: 'Recreate',
      },
      template: {
        metadata: {
          labels: {
            app: 'searxng',
          },
        },
        spec: {
          restartPolicy: 'Always',
          containers: [
            {
              name: 'searxng',
              image: 'searxng/searxng:latest',
              imagePullPolicy: 'IfNotPresent',
              env: [
                { name: 'SEARXNG_BASE_URL', value: 'http://searxng:8080' },
                { name: 'TZ', value: 'America/Chicago' },
              ],
              livenessProbe: {
                exec: {
                  command: [
                    '/bin/sh',
                    '-c',
                    'wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/ || exit 1',
                  ],
                },
                periodSeconds: 30,
                timeoutSeconds: 10,
                failureThreshold: 3,
              },
              ports: [
                {
                  name: 'searxng-8080',
                  containerPort: 8080,
                },
              ],
              volumeMounts: [
                {
                  name: 'etc-searxng',
                  mountPath: '/etc/searxng',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('256Mi'),
                  cpu: Quantity.fromString('50m'),
                },
                limits: {
                  memory: Quantity.fromString('512Mi'),
                  cpu: Quantity.fromString('200m'),
                },
              },
            },
          ],
          volumes: [
            {
              name: 'etc-searxng',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/searxng',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });

  // Dozzle deployment
  new KubeDeployment(scope, 'dozzle-deployment', {
    metadata: {
      name: 'dozzle',
      namespace: 'my-media-stack',
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          app: 'dozzle',
        },
      },
      strategy: {
        type: 'Recreate',
      },
      template: {
        metadata: {
          labels: {
            app: 'dozzle',
          },
        },
        spec: {
          restartPolicy: 'Always',
          containers: [
            {
              name: 'dozzle',
              image: 'amir20/dozzle',
              imagePullPolicy: 'IfNotPresent',
              ports: [
                {
                  name: 'dozzle-8080',
                  containerPort: 8080,
                },
              ],
              volumeMounts: [
                {
                  name: 'var-run-docker-sock',
                  mountPath: '/var/run/docker.sock',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('128Mi'),
                  cpu: Quantity.fromString('100m'),
                },
                limits: {
                  memory: Quantity.fromString('512Mi'),
                  cpu: Quantity.fromString('500m'),
                },
              },
            },
          ],
          volumes: [
            {
              name: 'var-run-docker-sock',
              hostPath: {
                path: '/var/run/docker.sock',
                type: 'Socket',
              },
            },
          ],
        },
      },
    },
  });

  // Speedtest deployment
  new KubeDeployment(scope, 'speedtest-deployment', {
    metadata: {
      name: 'speedtest',
      namespace: 'my-media-stack',
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          app: 'speedtest',
        },
      },
      strategy: {
        type: 'Recreate',
      },
      template: {
        metadata: {
          labels: {
            app: 'speedtest',
          },
        },
        spec: {
          restartPolicy: 'Always',
          containers: [
            {
              name: 'speedtest',
              image: 'linuxserver/speedtest-tracker',
              imagePullPolicy: 'IfNotPresent',
              env: [
                { name: 'ADMIN_EMAIL', value: 'boden.crouch@gmail.com' },
                { name: 'ADMIN_NAME', value: 'ubuntu' },
                { name: 'ADMIN_PASSWORD', value: 'h4L0m4St3R327' },
                { name: 'API_RATE_LIMIT', value: '60' },
                { name: 'APP_KEY', value: 'base64:OZVnVd3/7lJTtOTtVZ0YqxXJKtle2R48NIrfFv11CkM=' },
                { name: 'APP_NAME', value: 'Speedtest Tracker' },
                { name: 'APP_TIMEZONE', value: 'America/Chicago' },
                { name: 'APP_URL', value: '' },
                { name: 'ASSET_URL', value: '' },
                { name: 'CHART_BEGIN_AT_ZERO', value: 'true' },
                { name: 'CHART_DATETIME_FORMAT', value: 'j/m G:i' },
                { name: 'CONTENT_WIDTH', value: '7xl' },
                { name: 'DATETIME_FORMAT', value: 'j M Y, G:i:s' },
                { name: 'DB_CONNECTION', value: 'sqlite' },
                { name: 'DISPLAY_TIMEZONE', value: 'America/Chicago' },
                { name: 'PGID', value: '988' },
                { name: 'PRUNE_RESULTS_OLDER_THAN', value: '0' },
                { name: 'PUBLIC_DASHBOARD', value: 'true' },
                { name: 'PUID', value: '1001' },
                { name: 'SPEEDTEST_BLOCKED_SERVERS', value: '' },
                { name: 'SPEEDTEST_INTERFACE', value: '' },
                { name: 'SPEEDTEST_SCHEDULE', value: '0 * * * *' },
                { name: 'SPEEDTEST_SERVERS', value: '' },
                { name: 'SPEEDTEST_SKIP_IPS', value: '' },
                { name: 'THRESHOLD_DOWNLOAD', value: '900' },
                { name: 'THRESHOLD_ENABLED', value: 'true' },
                { name: 'THRESHOLD_PING', value: '25' },
                { name: 'THRESHOLD_UPLOAD', value: '900' },
                { name: 'TZ', value: 'America/Chicago' },
                { name: 'UMASK', value: '002' },
              ],
              livenessProbe: {
                exec: {
                  command: [
                    '/bin/sh',
                    '-c',
                    'curl -fs http://127.0.0.1:80 > /dev/null || exit 1',
                  ],
                },
                periodSeconds: 30,
                initialDelaySeconds: 20,
                timeoutSeconds: 15,
              },
              ports: [
                {
                  name: 'speedtest-443',
                  containerPort: 443,
                },
                {
                  name: 'speedtest-80',
                  containerPort: 80,
                },
              ],
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/config',
                },
                {
                  name: 'config-keys',
                  mountPath: '/config/keys',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('128Mi'),
                  cpu: Quantity.fromString('100m'),
                },
                limits: {
                  memory: Quantity.fromString('512Mi'),
                  cpu: Quantity.fromString('500m'),
                },
              },
            },
          ],
          volumes: [
            {
              name: 'config',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/speedtest-tracker/config',
                type: 'DirectoryOrCreate',
              },
            },
            {
              name: 'config-keys',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/certs/speedtest-tracker/keys',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });
} 