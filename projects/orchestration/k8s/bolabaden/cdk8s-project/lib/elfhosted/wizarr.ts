import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Wizarr resources
 */
export function createWizarrResources(scope: Construct): void {
  // Wizarr Deployment
  new KubeDeployment(scope, 'wizarr-deployment', {
    metadata: {
      name: 'brunner56-wizarr',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'wizarr',
        'app.kubernetes.io/version': '4.1.1',
        'helm.sh/chart': 'wizarr-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'wizarr-env,wizarr-steps-plex,wizarr-steps-jellyfin',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'wizarr',
        },
      },
      template: {
        metadata: {
          labels: {
            'app.elfhosted.com/name': 'wizarr',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'wizarr',
          },
        },
        spec: {
          volumes: [
            {
              name: 'config',
              persistentVolumeClaim: {
                claimName: 'config',
              },
            },
            {
              name: 'logs',
              persistentVolumeClaim: {
                claimName: 'logs',
              },
            },
            {
              name: 'tmp',
              emptyDir: {},
            },
            {
              name: 'wizarr-steps-plex',
              configMap: {
                name: 'wizarr-steps-plex',
                defaultMode: 420,
              },
            },
            {
              name: 'wizarr-steps-jellyfin',
              configMap: {
                name: 'wizarr-steps-jellyfin',
                defaultMode: 420,
              },
            },
          ],
          containers: [
            {
              name: 'wizarr',
              image: 'ghcr.io/wizarrrr/wizarr:4.1.1',
              ports: [
                {
                  name: 'http',
                  containerPort: 5690,
                  protocol: 'TCP',
                },
              ],
              envFrom: [
                {
                  configMapRef: {
                    name: 'wizarr-env',
                  },
                },
              ],
              resources: {
                limits: {
                  cpu: Quantity.fromString('1'),
                  memory: Quantity.fromString('1Gi'),
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('1Mi'),
                },
              },
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/data/database/',
                  subPath: 'wizarr',
                },
                {
                  name: 'logs',
                  mountPath: '/data/logs/',
                  subPath: 'wizarr',
                },
                {
                  name: 'tmp',
                  mountPath: '/tmp',
                },
                {
                  name: 'wizarr-steps-plex',
                  mountPath: '/data/custom_libs/plex.json',
                  subPath: 'plex.json',
                },
                {
                  name: 'wizarr-steps-jellyfin',
                  mountPath: '/data/custom_libs/jellyfin.json',
                  subPath: 'jellyfin.json',
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(5690),
                  scheme: 'HTTP',
                },
                initialDelaySeconds: 30,
                timeoutSeconds: 30,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3,
              },
              readinessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(5690),
                  scheme: 'HTTP',
                },
                initialDelaySeconds: 30,
                timeoutSeconds: 30,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3,
              },
              startupProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(5690),
                  scheme: 'HTTP',
                },
                initialDelaySeconds: 30,
                timeoutSeconds: 30,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 30,
              },
              terminationMessagePath: '/dev/termination-log',
              terminationMessagePolicy: 'File',
              imagePullPolicy: 'IfNotPresent',
              securityContext: {
                capabilities: {
                  drop: ['ALL'],
                },
                runAsUser: 568,
                runAsGroup: 568,
                readOnlyRootFilesystem: true,
                allowPrivilegeEscalation: false,
                seccompProfile: {
                  type: 'RuntimeDefault',
                },
              },
            },
          ],
          restartPolicy: 'Always',
          terminationGracePeriodSeconds: 30,
          dnsPolicy: 'ClusterFirst',
          serviceAccountName: 'default',
          serviceAccount: 'default',
          securityContext: {
            fsGroup: 568,
            seccompProfile: {
              type: 'RuntimeDefault',
            },
          },
          schedulerName: 'default-scheduler',
          priorityClassName: 'tenant-normal',
        },
      },
      strategy: {
        type: 'RollingUpdate',
        rollingUpdate: {
          maxUnavailable: IntOrString.fromNumber(1),
          maxSurge: IntOrString.fromString('25%'),
        },
      },
      revisionHistoryLimit: 10,
      progressDeadlineSeconds: 600,
    },
  });

  // Wizarr Service
  new KubeService(scope, 'wizarr-service', {
    metadata: {
      name: 'brunner56-wizarr',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'wizarr',
        'app.kubernetes.io/service': 'brunner56-wizarr',
        'helm.sh/chart': 'wizarr-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      ports: [
        {
          name: 'http',
          protocol: 'TCP',
          port: 5690,
          targetPort: IntOrString.fromString('http'),
        },
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'wizarr',
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster',
    },
  });
} 