import { Construct } from 'constructs';
import { KubeDeployment, Quantity } from '../../imports/k8s';

/**
 * Creates development service deployments (Code Server, etc.)
 */
export function createDevelopmentServices(scope: Construct): void {
  // Code Server development deployment
  new KubeDeployment(scope, 'code-dev-deployment', {
    metadata: {
      name: 'code-dev',
      namespace: 'development',
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          app: 'code-dev',
        },
      },
      strategy: {
        type: 'Recreate',
      },
      template: {
        metadata: {
          labels: {
            app: 'code-dev',
          },
        },
        spec: {
          restartPolicy: 'Always',
          containers: [
            {
              name: 'code-dev',
              image: 'linuxserver/code-server:latest',
              imagePullPolicy: 'IfNotPresent',
              env: [
                { name: 'CODESERVER_PASSWORD', value: 'ubuntu' },
                { name: 'CODESERVER_SUDO_PASSWORD', value: 'c3ll0h3r0123' },
                { name: 'DEFAULT_WORKSPACE', value: '/workspace' },
                { name: 'PGID', value: '988' },
                { name: 'PUID', value: '1001' },
                { name: 'TZ', value: 'America/Chicago' },
                { name: 'UMASK', value: '002' },
              ],
              ports: [
                {
                  name: 'code-dev-8443',
                  containerPort: 8443,
                },
              ],
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/config',
                },
                {
                  name: 'workspace',
                  mountPath: '/workspace',
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
                path: '/home/ubuntu/my-media-stack/configs/code-server/dev/config',
                type: 'DirectoryOrCreate',
              },
            },
            {
              name: 'workspace',
              hostPath: {
                path: '/home/ubuntu/my-media-stack',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });
} 