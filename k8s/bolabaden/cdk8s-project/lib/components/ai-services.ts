import { Construct } from 'constructs';
import { KubeDeployment, Quantity } from '../../imports/k8s';

/**
 * Creates AI service deployments (LobeChat, etc.)
 */
export function createAiServices(scope: Construct): void {
  // LobeChat deployment
  new KubeDeployment(scope, 'lobechat-deployment', {
    metadata: {
      name: 'lobechat',
      namespace: 'ai-services',
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          app: 'lobechat',
        },
      },
      template: {
        metadata: {
          labels: {
            app: 'lobechat',
          },
        },
        spec: {
          containers: [
            {
              name: 'lobechat',
              image: 'lobehub/lobe-chat:latest',
              ports: [
                {
                  name: 'lobechat-3210',
                  containerPort: 3210,
                },
              ],
              env: [
                { name: 'ACCESS_CODE', value: 'brunner56' },
                { name: 'PGID', value: '988' },
                { name: 'PUID', value: '1002' },
                { name: 'TZ', value: 'America/Chicago' },
                { name: 'UMASK', value: '002' },
              ],
              volumeMounts: [
                {
                  name: 'host-config',
                  mountPath: '/configs/lobechat/host_config',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('256Mi'),
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
              name: 'host-config',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/lobechat/host_config',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });
} 