import { Construct } from 'constructs';
import { KubeDeployment, Quantity } from '../../imports/k8s';

/**
 * Creates utility service deployments (Homepage, etc.)
 */
export function createUtilityServices(scope: Construct): void {
  // Homepage deployment
  new KubeDeployment(scope, 'homepage-deployment', {
    metadata: {
      name: 'homepage',
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
              name: 'homepage',
              image: 'ghcr.io/gethomepage/homepage',
              imagePullPolicy: 'IfNotPresent',
              env: [
                { name: 'HOMEPAGE_ALLOWED_HOSTS', value: '*' },
                { name: 'HOMEPAGE_VAR_HEADER_STYLE', value: '' },
                { name: 'HOMEPAGE_VAR_SEARCH_PROVIDER', value: 'google' },
                { name: 'HOMEPAGE_VAR_TITLE', value: 'my-media-stack' },
                { name: 'HOMEPAGE_VAR_WEATHER_CITY', value: 'Chicago' },
                { name: 'HOMEPAGE_VAR_WEATHER_LAT', value: '41.8781' },
                { name: 'HOMEPAGE_VAR_WEATHER_LONG', value: '-87.6298' },
                { name: 'HOMEPAGE_VAR_WEATHER_UNIT', value: 'fahrenheit' },
                { name: 'TZ', value: 'America/Chicago' },
              ],
              livenessProbe: {
                exec: {
                  command: [
                    '/bin/sh',
                    '-c',
                    'wget -qO- http://127.0.0.1:3000 > /dev/null 2>&1 || exit 1',
                  ],
                },
                periodSeconds: 30,
                timeoutSeconds: 15,
                failureThreshold: 3,
              },
              ports: [
                {
                  name: 'homepage-3000',
                  containerPort: 3000,
                },
              ],
              volumeMounts: [
                {
                  name: 'var-run-docker-sock',
                  mountPath: '/var/run/docker.sock',
                },
                {
                  name: 'app-config',
                  mountPath: '/app/config',
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
            {
              name: 'app-config',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/homepage',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });
} 