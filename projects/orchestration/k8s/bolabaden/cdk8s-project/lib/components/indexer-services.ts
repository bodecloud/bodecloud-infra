import { Construct } from 'constructs';
import { KubeDeployment, Quantity } from '../../imports/k8s';

/**
 * Creates indexer service deployments (Prowlarr, Jackett)
 */
export function createIndexerServices(scope: Construct): void {
  // Prowlarr deployment
  new KubeDeployment(scope, 'prowlarr-deployment', {
    metadata: {
      name: 'prowlarr',
      namespace: 'my-media-stack',
    },
    spec: {
      replicas: 1,
      selector: {},
      strategy: {
        type: 'Recreate',
      },
      template: {
        spec: {
          restartPolicy: 'Always',
          containers: [
            {
              name: 'prowlarr',
              image: 'linuxserver/prowlarr',
              imagePullPolicy: 'IfNotPresent',
              env: [
                { name: 'PGID', value: '988' },
                { name: 'PUID', value: '1001' },
                { name: 'TZ', value: 'America/Chicago' },
                { name: 'UMASK', value: '002' },
              ],
              livenessProbe: {
                exec: {
                  command: [
                    '/bin/sh',
                    '-c',
                    'curl -fs -H "Authorization: Bearer 29440a82740d475cacb35327c62c87a1" http://127.0.0.1:9696/api/v1/system/status || exit 1',
                  ],
                },
                periodSeconds: 30,
                initialDelaySeconds: 120,
                timeoutSeconds: 10,
                failureThreshold: 3,
              },
              ports: [
                {
                  name: 'prowlarr-9696',
                  containerPort: 9696,
                },
              ],
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/config',
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
                path: '/home/ubuntu/my-media-stack/configs/prowlarr/config',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });

  // Jackett deployment
  new KubeDeployment(scope, 'jackett-deployment', {
    metadata: {
      name: 'jackett',
      namespace: 'my-media-stack',
    },
    spec: {
      replicas: 1,
      selector: {},
      strategy: {
        type: 'Recreate',
      },
      template: {
        spec: {
          restartPolicy: 'Always',
          containers: [
            {
              name: 'jackett',
              image: 'linuxserver/jackett',
              imagePullPolicy: 'IfNotPresent',
              env: [
                { name: 'PGID', value: '988' },
                { name: 'PUID', value: '1002' },
                { name: 'TZ', value: 'America/Chicago' },
                { name: 'UMASK', value: '022' },
              ],
              securityContext: {
                runAsUser: 1002,
                runAsGroup: 988,
              },
              livenessProbe: {
                exec: {
                  command: [
                    '/bin/sh',
                    '-c',
                    'curl -fs http://127.0.0.1:9117/api/v2.0/indexers/all/results/torznab?t=indexers&apikey=nnx0n84pcj7umynyd2pbid2nl1zzuemz || exit 1',
                  ],
                },
                periodSeconds: 30,
                initialDelaySeconds: 60,
                timeoutSeconds: 10,
              },
              ports: [
                {
                  name: 'jackett-9117',
                  containerPort: 9117,
                },
              ],
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/config',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('256Mi'),
                  cpu: Quantity.fromString('100m'),
                },
                limits: {
                  memory: Quantity.fromString('512Mi'),
                  cpu: Quantity.fromString('300m'),
                },
              },
            },
          ],
          volumes: [
            {
              name: 'config',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/jackett/config',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });
} 