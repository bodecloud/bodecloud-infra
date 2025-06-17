import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Plex deployment and service
 */
export function createPlexResources(scope: Construct) {
  // Create Plex Deployment
  new KubeDeployment(scope, 'plex-deployment', {
    metadata: {
      name: 'brunner56-plex',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'plex',
        'app.kubernetes.io/version': '1.32.5.7349-8f4248874',
        'helm.sh/chart': 'plex-6.4.3',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'plex-config,elfbot-plex',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'plex'
        }
      },
      template: {
        metadata: {
          labels: {
            'app.elfhosted.com/name': 'plex',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'plex'
          },
          annotations: {
            'kubernetes.io/egress-bandwidth': '1G'
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
              name: 'logs',
              persistentVolumeClaim: {
                claimName: 'logs'
              }
            },
            {
              name: 'plex-config',
              configMap: {
                name: 'plex-config',
                defaultMode: 420
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
            },
            {
              name: 'transcode',
              persistentVolumeClaim: {
                claimName: 'transcode-1g'
              }
            }
          ],
          containers: [
            {
              name: 'brunner56-plex',
              image: 'ghcr.io/elfhosted/plex:1.32.5.7349-8f4248874@sha256:0f62b6b5b8b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5',
              ports: [
                {
                  name: 'plex',
                  containerPort: 32400,
                  protocol: 'TCP'
                }
              ],
              envFrom: [
                {
                  configMapRef: {
                    name: 'plex-config'
                  }
                }
              ],
              env: [
                {
                  name: 'TZ',
                  value: 'UTC'
                },
                {
                  name: 'PLEX_UID',
                  value: '568'
                },
                {
                  name: 'PLEX_GID',
                  value: '568'
                }
              ],
              resources: {
                limits: {
                  cpu: Quantity.fromString('4'),
                  memory: Quantity.fromString('8Gi')
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('128Mi')
                }
              },
              volumeMounts: [
                {
                  name: 'backup',
                  mountPath: '/storage/backup'
                },
                {
                  name: 'config',
                  mountPath: '/config',
                  subPath: 'plex'
                },
                {
                  name: 'logs',
                  mountPath: '/storage/logs'
                },
                {
                  name: 'plex-config',
                  mountPath: '/plex-config'
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
                },
                {
                  name: 'transcode',
                  mountPath: '/transcode'
                }
              ],
              livenessProbe: {
                httpGet: {
                  path: '/identity',
                  port: IntOrString.fromNumber(32400),
                  scheme: 'HTTP'
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              readinessProbe: {
                httpGet: {
                  path: '/identity',
                  port: IntOrString.fromNumber(32400),
                  scheme: 'HTTP'
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              startupProbe: {
                httpGet: {
                  path: '/identity',
                  port: IntOrString.fromNumber(32400),
                  scheme: 'HTTP'
                },
                timeoutSeconds: 1,
                periodSeconds: 5,
                successThreshold: 1,
                failureThreshold: 60
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
            fsGroupChangePolicy: 'OnRootMismatch',
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

  // Create Plex Service
  new KubeService(scope, 'plex-service', {
    metadata: {
      name: 'brunner56-plex',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'plex',
        'app.kubernetes.io/service': 'brunner56-plex',
        'app.kubernetes.io/version': '1.32.5.7349-8f4248874',
        'helm.sh/chart': 'plex-6.4.3',
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
          name: 'plex',
          protocol: 'TCP',
          port: 32400,
          targetPort: IntOrString.fromString('plex')
        }
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'plex'
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster'
    }
  });
} 