import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates RcloneUI resources for the ElfHosted platform
 */
export function createRcloneUIResources(scope: Construct): void {
  // Create RcloneUI Deployment
  new KubeDeployment(scope, 'rcloneui-deployment', {
    metadata: {
      name: 'brunner56-rcloneui',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rcloneui',
        'app.kubernetes.io/version': '2.0.5',
        'helm.sh/chart': 'rcloneui-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'rcloneui-configmap',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'rcloneui'
        }
      },
      template: {
        metadata: {
          labels: {
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'rcloneui'
          }
        },
        spec: {
          volumes: [
            {
              name: 'config',
              persistentVolumeClaim: {
                claimName: 'config'
              }
            },
            {
              name: 'rclone',
              persistentVolumeClaim: {
                claimName: 'rclone'
              }
            },
            {
              name: 'rcloneui-config',
              configMap: {
                name: 'rcloneui-configmap',
                defaultMode: 420
              }
            }
          ],
          containers: [
            {
              name: 'brunner56-rcloneui',
              image: 'ghcr.io/elfhosted/rclone:1.68.2@sha256:b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0',
              ports: [
                {
                  name: 'http',
                  containerPort: 5572,
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
                  cpu: Quantity.fromString('1'),
                  memory: Quantity.fromString('1Gi')
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('1Mi')
                }
              },
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/storage/config',
                  subPath: 'rcloneui'
                },
                {
                  name: 'rclone',
                  mountPath: '/storage/rclone'
                },
                {
                  name: 'rcloneui-config',
                  mountPath: '/config/rclone',
                  subPath: 'rclone.conf'
                }
              ],
              livenessProbe: {
                tcpSocket: {
                  port: IntOrString.fromNumber(5572)
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              readinessProbe: {
                tcpSocket: {
                  port: IntOrString.fromNumber(5572)
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              startupProbe: {
                tcpSocket: {
                  port: IntOrString.fromNumber(5572)
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
                readOnlyRootFilesystem: false,
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
            fsGroup: 568,
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
        type: 'Recreate'
      },
      revisionHistoryLimit: 3,
      progressDeadlineSeconds: 600
    }
  });

  // Create RcloneUI Service
  new KubeService(scope, 'rcloneui-service', {
    metadata: {
      name: 'brunner56-rcloneui',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rcloneui',
        'app.kubernetes.io/service': 'brunner56-rcloneui',
        'helm.sh/chart': 'rcloneui-1.0.0',
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
          targetPort: IntOrString.fromString('http')
        }
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'rcloneui'
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster'
    }
  });
} 