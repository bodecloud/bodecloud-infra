import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Homer deployment and service
 */
export function createHomerResources(scope: Construct) {
  // Create Homer Deployment
  new KubeDeployment(scope, 'homer-deployment', {
    metadata: {
      name: 'brunner56-homer',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'homer',
        'app.kubernetes.io/version': 'v22.07.2',
        'helm.sh/chart': 'homer-8.0.2',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'homer-config, elfbot-homer',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'homer'
        }
      },
      template: {
        metadata: {
          labels: {
            'app.elfhosted.com/role': 'nodefinder',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'homer'
          }
        },
        spec: {
          volumes: [
            {
              name: 'backup',
              persistentVolumeClaim: {
                claimName: 'backup'
              }
            },
            {
              name: 'config',
              persistentVolumeClaim: {
                claimName: 'config'
              }
            },
            {
              name: 'config-yml',
              configMap: {
                name: 'homer-config',
                defaultMode: 420
              }
            },
            {
              name: 'custom-css',
              configMap: {
                name: 'homer-config',
                defaultMode: 420
              }
            },
            {
              name: 'disk-usage',
              configMap: {
                name: 'homer-config',
                defaultMode: 420
              }
            },
            {
              name: 'gatus-config',
              configMap: {
                name: 'gatus-config',
                defaultMode: 420
              }
            },
            {
              name: 'logs',
              persistentVolumeClaim: {
                claimName: 'logs'
              }
            },
            {
              name: 'message',
              emptyDir: {}
            },
            {
              name: 'rclone',
              persistentVolumeClaim: {
                claimName: 'rclone'
              }
            },
            {
              name: 'rclonemountrealdebridzurg',
              persistentVolumeClaim: {
                claimName: 'realdebrid-zurg'
              }
            },
            {
              name: 'symlinks',
              persistentVolumeClaim: {
                claimName: 'symlinks'
              }
            },
            {
              name: 'tmp',
              emptyDir: {}
            }
          ],
          containers: [
            {
              name: 'brunner56-homer',
              image: 'ghcr.io/elfhosted/tooling:focal-20230605@sha256:6088a1e9fc0ce83aec9910af0899661c23b5f2025428d7da631b9b9390241b6c',
              command: [
                '/bin/bash',
                '/usr/local/bin/disk_usage.sh'
              ],
              ports: [
                {
                  name: 'http',
                  containerPort: 8080,
                  protocol: 'TCP'
                }
              ],
              env: [
                {
                  name: 'TZ',
                  value: 'UTC'
                }
              ],
              resources: {
                limits: {
                  cpu: Quantity.fromString('200m'),
                  memory: Quantity.fromString('1Gi')
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('1Mi')
                }
              },
              volumeMounts: [
                {
                  name: 'backup',
                  mountPath: '/backup'
                },
                {
                  name: 'config',
                  mountPath: '/config'
                },
                {
                  name: 'config-yml',
                  mountPath: '/config-yml',
                  subPath: 'config.yml'
                },
                {
                  name: 'custom-css',
                  mountPath: '/custom-css',
                  subPath: 'custom-css'
                },
                {
                  name: 'disk-usage',
                  mountPath: '/usr/local/bin/disk_usage.sh',
                  subPath: 'disk_usage.sh'
                },
                {
                  name: 'gatus-config',
                  mountPath: '/gatus-config'
                },
                {
                  name: 'logs',
                  mountPath: '/logs'
                },
                {
                  name: 'message',
                  mountPath: '/www/assets/message'
                },
                {
                  name: 'rclone',
                  mountPath: '/storage/rclone'
                },
                {
                  name: 'rclonemountrealdebridzurg',
                  mountPath: '/storage/realdebrid-zurg'
                },
                {
                  name: 'symlinks',
                  mountPath: '/storage/symlinks'
                },
                {
                  name: 'tmp',
                  mountPath: '/tmp'
                }
              ],
              livenessProbe: {
                tcpSocket: {
                  port: IntOrString.fromNumber(8080)
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              readinessProbe: {
                tcpSocket: {
                  port: IntOrString.fromNumber(8080)
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              startupProbe: {
                tcpSocket: {
                  port: IntOrString.fromNumber(8080)
                },
                timeoutSeconds: 1,
                periodSeconds: 5,
                successThreshold: 1,
                failureThreshold: 30
              },
              terminationMessagePath: '/dev/termination-log',
              terminationMessagePolicy: 'File',
              imagePullPolicy: 'IfNotPresent',
              securityContext: {
                runAsUser: 568,
                runAsGroup: 568,
                runAsNonRoot: false,
                readOnlyRootFilesystem: true,
                seccompProfile: {
                  type: 'RuntimeDefault'
                }
              }
            },
            {
              name: 'ui',
              image: 'ghcr.io/elfhosted/homer:v25.05.2@sha256:60772bd0292281282161f10aa2bea105bba344158937f1c0022c6986ef5aedc3',
              resources: {
                limits: {
                  cpu: Quantity.fromString('1'),
                  memory: Quantity.fromString('4Gi')
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('1Mi')
                }
              },
              volumeMounts: [
                {
                  name: 'config-yml',
                  mountPath: '/www/assets/config.yml',
                  subPath: 'config.yml'
                },
                {
                  name: 'custom-css',
                  mountPath: '/www/assets/custom.css',
                  subPath: 'custom.css'
                },
                {
                  name: 'message',
                  mountPath: '/www/assets/message'
                },
                {
                  name: 'config',
                  readOnly: true,
                  mountPath: '/www/assets/backgrounds',
                  subPath: 'homer/backgrounds'
                }
              ],
              terminationMessagePath: '/dev/termination-log',
              terminationMessagePolicy: 'File',
              imagePullPolicy: 'IfNotPresent',
              securityContext: {
                capabilities: {
                  drop: ['ALL']
                },
                runAsUser: 568,
                runAsGroup: 568,
                readOnlyRootFilesystem: true,
                allowPrivilegeEscalation: false,
                seccompProfile: {
                  type: 'RuntimeDefault'
                }
              }
            }
          ],
          restartPolicy: 'Always',
          terminationGracePeriodSeconds: 30,
          dnsPolicy: 'ClusterFirst',
          serviceAccountName: 'default',
          serviceAccount: 'default',
          automountServiceAccountToken: false,
          securityContext: {
            runAsNonRoot: false,
            fsGroup: 568,
            fsGroupChangePolicy: 'Always',
            seccompProfile: {
              type: 'RuntimeDefault'
            }
          },
          schedulerName: 'default-scheduler',
          priorityClassName: 'tenant-normal',
          enableServiceLinks: false
        }
      },
      strategy: {
        type: 'RollingUpdate',
        rollingUpdate: {
          maxUnavailable: IntOrString.fromNumber(1),
          maxSurge: IntOrString.fromString('25%')
        }
      },
      revisionHistoryLimit: 3,
      progressDeadlineSeconds: 600
    }
  });

  // Create Homer Service
  new KubeService(scope, 'homer-service', {
    metadata: {
      name: 'brunner56-homer',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'homer',
        'app.kubernetes.io/service': 'brunner56-homer',
        'app.kubernetes.io/version': 'v22.07.2',
        'helm.sh/chart': 'homer-8.0.2',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      ports: [
        {
          name: 'http',
          protocol: 'TCP',
          port: 8080,
          targetPort: IntOrString.fromString('http')
        }
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'homer'
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster'
    }
  });
}
