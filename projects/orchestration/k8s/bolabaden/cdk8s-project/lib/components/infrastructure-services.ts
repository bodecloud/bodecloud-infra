import { Construct } from 'constructs';
import { KubeDeployment, Quantity } from '../../imports/k8s';

/**
 * Creates infrastructure service deployments (MongoDB, Redis, etc.)
 */
export function createInfrastructureServices(scope: Construct): void {
  // MongoDB deployment
  new KubeDeployment(scope, 'mongodb-deployment', {
    metadata: {
      name: 'mongodb',
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
              name: 'mongodb',
              image: 'mongo:latest',
              imagePullPolicy: 'IfNotPresent',
              command: ['mongod'],
              args: ['--bind_ip_all'],
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
                    'mongosh 127.0.0.1:27017/test --quiet --eval \'db.runCommand("ping").ok\' > /dev/null 2>&1 || exit 1',
                  ],
                },
                periodSeconds: 10,
                initialDelaySeconds: 40,
                timeoutSeconds: 10,
                failureThreshold: 5,
              },
              ports: [
                {
                  name: 'mongodb-27017',
                  containerPort: 27017,
                },
              ],
              volumeMounts: [
                {
                  name: 'data-db',
                  mountPath: '/data/db',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('256Mi'),
                  cpu: Quantity.fromString('200m'),
                },
                limits: {
                  memory: Quantity.fromString('1Gi'),
                  cpu: Quantity.fromString('1000m'),
                },
              },
            },
          ],
          volumes: [
            {
              name: 'data-db',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/mongodb/data',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });

  // Redis deployment
  new KubeDeployment(scope, 'redis-deployment', {
    metadata: {
      name: 'redis',
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
              name: 'redis',
              image: 'redis:latest',
              imagePullPolicy: 'IfNotPresent',
              command: ['redis-server', '--appendonly', 'yes', '--save', '60', '1'],
              livenessProbe: {
                exec: {
                  command: [
                    '/bin/sh',
                    '-c',
                    'redis-cli ping > /dev/null 2>&1 || exit 1',
                  ],
                },
                periodSeconds: 10,
                timeoutSeconds: 5,
                failureThreshold: 5,
              },
              ports: [
                {
                  name: 'redis-6379',
                  containerPort: 6379,
                },
              ],
              volumeMounts: [
                {
                  name: 'data',
                  mountPath: '/data',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('256Mi'),
                  cpu: Quantity.fromString('100m'),
                },
                limits: {
                  memory: Quantity.fromString('1Gi'),
                  cpu: Quantity.fromString('500m'),
                },
              },
            },
          ],
          volumes: [
            {
              name: 'data',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/redis',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });
} 