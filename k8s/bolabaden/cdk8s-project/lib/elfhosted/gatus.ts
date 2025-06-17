import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Gatus deployment and service
 */
export function createGatusResources(scope: Construct) {
  // Create Gatus Deployment
  new KubeDeployment(scope, 'gatus-deployment', {
    metadata: {
      name: 'brunner56-gatus',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'gatus',
        'helm.sh/chart': 'gatus-0.2.1',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'gatus-config',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'gatus'
        }
      },
      template: {
        metadata: {
          labels: {
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'gatus'
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
              name: 'gatus-config',
              configMap: {
                name: 'gatus-config',
                defaultMode: 420
              }
            }
          ],
          containers: [
            {
              name: 'brunner56-gatus',
              image: 'ghcr.io/elfhosted/gatus:5.17.0@sha256:e2752cf7e1781478b12e0ca48159a40549905195afe39fe0a5e39ac2c7256dac',
              ports: [
                {
                  name: 'http',
                  containerPort: 8080,
                  protocol: 'TCP'
                }
              ],
              envFrom: [
                {
                  secretRef: {
                    name: 'gatus-smtp-config'
                  }
                }
              ],
              env: [
                {
                  name: 'GATUS_CONFIG_PATH',
                  value: '/config/config.yaml'
                },
                {
                  name: 'SMTP_FROM',
                  value: 'health@elfhosted.com'
                },
                {
                  name: 'SMTP_PORT',
                  value: '587'
                }
              ],
              resources: {
                limits: {
                  cpu: Quantity.fromString('1'),
                  memory: Quantity.fromString('128Mi')
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('20Mi')
                }
              },
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/data/',
                  subPath: 'gatus'
                },
                {
                  name: 'gatus-config',
                  mountPath: '/config'
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
                readOnlyRootFilesystem: true,
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
          affinity: {
            podAffinity: {
              requiredDuringSchedulingIgnoredDuringExecution: [
                {
                  labelSelector: {
                    matchExpressions: [
                      {
                        key: 'app.elfhosted.com/role',
                        operator: 'In',
                        values: ['nodefinder']
                      }
                    ]
                  },
                  topologyKey: 'kubernetes.io/hostname'
                }
              ]
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

  // Create Gatus Service
  new KubeService(scope, 'gatus-service', {
    metadata: {
      name: 'brunner56-gatus',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'gatus',
        'app.kubernetes.io/service': 'brunner56-gatus',
        'helm.sh/chart': 'gatus-0.2.1',
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
        'app.kubernetes.io/name': 'gatus'
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster'
    }
  });
} 