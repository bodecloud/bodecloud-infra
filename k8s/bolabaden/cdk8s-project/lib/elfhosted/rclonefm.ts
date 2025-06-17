import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates RcloneFM resources
 */
export function createRcloneFMResources(scope: Construct): void {
  // RcloneFM Deployment
  new KubeDeployment(scope, 'rclonefm-deployment', {
    metadata: {
      name: 'brunner56-rclonefm',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rclonefm',
        'app.kubernetes.io/version': '1.0.0',
        'helm.sh/chart': 'rclonefm-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'rclonefm-config',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'rclonefm',
        },
      },
      template: {
        metadata: {
          labels: {
            'app.elfhosted.com/name': 'rclonefm',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'rclonefm',
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
              name: 'rclonefm-config',
              configMap: {
                name: 'rclonefm-config',
                defaultMode: 420,
              },
            },
          ],
          containers: [
            {
              name: 'rclonefm',
              image: 'ghcr.io/elfhosted/rclone:1.68.1@sha256:b5c7f7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8',
              command: [
                '/bin/bash',
                '-c',
                'rclone rcd --rc-web-gui --rc-addr=0.0.0.0:5572 --rc-user=admin --rc-pass=admin --config=/config/rclone/rclone.conf --log-file=/logs/rclonefm/rclonefm.log --log-level=INFO',
              ],
              ports: [
                {
                  name: 'http',
                  containerPort: 5572,
                  protocol: 'TCP',
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
                  mountPath: '/config/',
                },
                {
                  name: 'logs',
                  mountPath: '/logs/',
                },
                {
                  name: 'rclone',
                  mountPath: '/storage/rclone/',
                },
                {
                  name: 'rclonemountrealdebridzurg',
                  mountPath: '/storage/realdebrid-zurg/',
                },
                {
                  name: 'symlinks',
                  mountPath: '/storage/symlinks/',
                },
                {
                  name: 'tmp',
                  mountPath: '/tmp',
                },
                {
                  name: 'rclonefm-config',
                  mountPath: '/config/rclone/rclone.conf',
                  subPath: 'rclone.conf',
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(5572),
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
                  port: IntOrString.fromNumber(5572),
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
                  port: IntOrString.fromNumber(5572),
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

  // RcloneFM Service
  new KubeService(scope, 'rclonefm-service', {
    metadata: {
      name: 'brunner56-rclonefm',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rclonefm',
        'app.kubernetes.io/service': 'brunner56-rclonefm',
        'helm.sh/chart': 'rclonefm-1.0.0',
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
          port: 5572,
          targetPort: IntOrString.fromString('http'),
        },
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'rclonefm',
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster',
    },
  });
} 