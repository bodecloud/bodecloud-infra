import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Riven backend and frontend resources
 */
export function createRivenResources(scope: Construct): void {
  // Riven Backend Deployment
  new KubeDeployment(scope, 'riven-deployment', {
    metadata: {
      name: 'brunner56-riven',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'riven',
        'app.kubernetes.io/version': '1.0.0',
        'helm.sh/chart': 'riven-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'riven-env,riven-setup',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'riven',
        },
      },
      template: {
        metadata: {
          labels: {
            'app.elfhosted.com/name': 'riven',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'riven',
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
              name: 'rclone',
              persistentVolumeClaim: {
                claimName: 'rclone',
              },
            },
            {
              name: 'rclonemountrealdebridzurg',
              persistentVolumeClaim: {
                claimName: 'realdebrid-zurg',
              },
            },
            {
              name: 'symlinks',
              persistentVolumeClaim: {
                claimName: 'symlinks',
              },
            },
            {
              name: 'tmp',
              emptyDir: {},
            },
            {
              name: 'riven-env',
              configMap: {
                name: 'riven-env',
                defaultMode: 420,
              },
            },
            {
              name: 'riven-setup',
              configMap: {
                name: 'riven-setup',
                defaultMode: 493,
              },
            },
          ],
          initContainers: [
            {
              name: 'setup',
              image: 'ghcr.io/elfhosted/tooling:focal-20240530@sha256:458d1f3b54e9455b5cdad3c341d6853a6fdd75ac3f1120931ca3c09ac4b588de',
              command: ['/bin/bash', '/setup/setup.sh'],
              resources: {},
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/config/',
                  subPath: 'riven',
                },
                {
                  name: 'riven-setup',
                  mountPath: '/setup/',
                },
              ],
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
          containers: [
            {
              name: 'riven',
              image: 'spoked/riven:latest',
              ports: [
                {
                  name: 'http',
                  containerPort: 8080,
                  protocol: 'TCP',
                },
              ],
              envFrom: [
                {
                  configMapRef: {
                    name: 'riven-env',
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
                  mountPath: '/riven/data/',
                  subPath: 'riven',
                },
                {
                  name: 'logs',
                  mountPath: '/riven/data/logs/',
                  subPath: 'riven',
                },
                {
                  name: 'rclone',
                  mountPath: '/riven/data/mounts/remote/',
                },
                {
                  name: 'rclonemountrealdebridzurg',
                  mountPath: '/riven/data/mounts/realdebrid/',
                },
                {
                  name: 'symlinks',
                  mountPath: '/riven/data/symlinks/',
                },
                {
                  name: 'tmp',
                  mountPath: '/tmp',
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(8080),
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
                  port: IntOrString.fromNumber(8080),
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
                  port: IntOrString.fromNumber(8080),
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
              imagePullPolicy: 'Always',
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

  // Riven Frontend Deployment
  new KubeDeployment(scope, 'riven-frontend-deployment', {
    metadata: {
      name: 'brunner56-rivenfrontend',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rivenfrontend',
        'app.kubernetes.io/version': '1.0.0',
        'helm.sh/chart': 'rivenfrontend-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'riven-frontend-config,riven-frontend-env',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'rivenfrontend',
        },
      },
      template: {
        metadata: {
          labels: {
            'app.elfhosted.com/name': 'rivenfrontend',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'rivenfrontend',
          },
        },
        spec: {
          volumes: [
            {
              name: 'riven-frontend-config',
              configMap: {
                name: 'riven-frontend-config',
                defaultMode: 420,
              },
            },
          ],
          containers: [
            {
              name: 'rivenfrontend',
              image: 'spoked/riven-frontend:latest',
              ports: [
                {
                  name: 'http',
                  containerPort: 3000,
                  protocol: 'TCP',
                },
              ],
              envFrom: [
                {
                  configMapRef: {
                    name: 'riven-frontend-env',
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
                  name: 'riven-frontend-config',
                  mountPath: '/app/src/lib/config.ts',
                  subPath: 'config.ts',
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(3000),
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
                  port: IntOrString.fromNumber(3000),
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
                  port: IntOrString.fromNumber(3000),
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
              imagePullPolicy: 'Always',
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

  // Riven Backend Service
  new KubeService(scope, 'riven-service', {
    metadata: {
      name: 'brunner56-riven',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'riven',
        'app.kubernetes.io/service': 'brunner56-riven',
        'helm.sh/chart': 'riven-1.0.0',
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
          port: 8080,
          targetPort: IntOrString.fromString('http'),
        },
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'riven',
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster',
    },
  });

  // Riven Frontend Service
  new KubeService(scope, 'riven-frontend-service', {
    metadata: {
      name: 'brunner56-rivenfrontend',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rivenfrontend',
        'app.kubernetes.io/service': 'brunner56-rivenfrontend',
        'helm.sh/chart': 'rivenfrontend-1.0.0',
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
          port: 3000,
          targetPort: IntOrString.fromString('http'),
        },
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'rivenfrontend',
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster',
    },
  });
} 