import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Zurg resources
 */
export function createZurgResources(scope: Construct): void {
  // Zurg Deployment
  new KubeDeployment(scope, 'zurg-deployment', {
    metadata: {
      name: 'brunner56-zurg',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'zurg',
        'app.kubernetes.io/version': '0.9.3-hotfix.6',
        'helm.sh/chart': 'zurg-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'zurg-config,zurg-env',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'zurg',
        },
      },
      template: {
        metadata: {
          labels: {
            'app.elfhosted.com/name': 'zurg',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'zurg',
          },
          annotations: {
            'kubernetes.io/egress-bandwidth': '100M',
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
              name: 'rclonemountrealdebridzurg',
              persistentVolumeClaim: {
                claimName: 'realdebrid-zurg',
              },
            },
            {
              name: 'tmp',
              emptyDir: {},
            },
            {
              name: 'zurg-config',
              configMap: {
                name: 'zurg-config',
                defaultMode: 420,
              },
            },
          ],
          containers: [
            {
              name: 'zurg',
              image: 'ghcr.io/debridmediamanager/zurg-testing:v0.9.3-hotfix.6',
              ports: [
                {
                  name: 'http',
                  containerPort: 9999,
                  protocol: 'TCP',
                },
              ],
              envFrom: [
                {
                  configMapRef: {
                    name: 'zurg-env',
                  },
                },
              ],
              resources: {
                limits: {
                  cpu: Quantity.fromString('4'),
                  memory: Quantity.fromString('8Gi'),
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('1Mi'),
                },
              },
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/app/config/',
                  subPath: 'zurg',
                },
                {
                  name: 'logs',
                  mountPath: '/app/logs/',
                  subPath: 'zurg',
                },
                {
                  name: 'rclonemountrealdebridzurg',
                  mountPath: '/app/RD/',
                },
                {
                  name: 'tmp',
                  mountPath: '/tmp',
                },
                {
                  name: 'zurg-config',
                  mountPath: '/app/config.yml',
                  subPath: 'config.yml',
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/dav/',
                  port: IntOrString.fromNumber(9999),
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
                  path: '/dav/',
                  port: IntOrString.fromNumber(9999),
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
                  path: '/dav/',
                  port: IntOrString.fromNumber(9999),
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
        type: 'Recreate',
      },
      revisionHistoryLimit: 10,
      progressDeadlineSeconds: 600,
    },
  });

  // Zurg Service
  new KubeService(scope, 'zurg-service', {
    metadata: {
      name: 'brunner56-zurg',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'zurg',
        'app.kubernetes.io/service': 'brunner56-zurg',
        'helm.sh/chart': 'zurg-1.0.0',
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
          port: 9999,
          targetPort: IntOrString.fromString('http'),
        },
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'zurg',
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster',
    },
  });
} 