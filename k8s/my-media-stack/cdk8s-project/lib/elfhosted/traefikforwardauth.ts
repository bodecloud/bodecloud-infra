import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Traefik Forward Auth resources for the ElfHosted platform
 */
export function createTraefikForwardAuthResources(scope: Construct): void {
  // Create Traefik Forward Auth Deployment
  new KubeDeployment(scope, 'traefikforwardauth-deployment', {
    metadata: {
      name: 'brunner56-traefikforwardauth',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'traefikforwardauth',
        'app.kubernetes.io/version': '2.2.0',
        'helm.sh/chart': 'traefikforwardauth-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'configmap.reloader.stakater.com/reload': 'traefik-forward-auth-config',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'traefikforwardauth'
        }
      },
      template: {
        metadata: {
          labels: {
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/name': 'traefikforwardauth'
          }
        },
        spec: {
          volumes: [
            {
              name: 'config',
              configMap: {
                name: 'traefik-forward-auth-config',
                defaultMode: 420
              }
            }
          ],
          containers: [
            {
              name: 'brunner56-traefikforwardauth',
              image: 'ghcr.io/elfhosted/traefik-forward-auth:2.2.0@sha256:a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2',
              ports: [
                {
                  name: 'http',
                  containerPort: 4181,
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
                  cpu: Quantity.fromString('100m'),
                  memory: Quantity.fromString('128Mi')
                },
                requests: {
                  cpu: Quantity.fromString('10m'),
                  memory: Quantity.fromString('64Mi')
                }
              },
              volumeMounts: [
                {
                  name: 'config',
                  mountPath: '/config'
                }
              ],
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(4181),
                  scheme: 'HTTP'
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              readinessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(4181),
                  scheme: 'HTTP'
                },
                timeoutSeconds: 1,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3
              },
              startupProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(4181),
                  scheme: 'HTTP'
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
                readOnlyRootFilesystem: true,
                allowPrivilegeEscalation: false,
                capabilities: {
                  drop: ['ALL']
                },
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
        type: 'RollingUpdate',
        rollingUpdate: {
          maxUnavailable: IntOrString.fromString('25%'),
          maxSurge: IntOrString.fromString('25%')
        }
      },
      revisionHistoryLimit: 3,
      progressDeadlineSeconds: 600
    }
  });

  // Create Traefik Forward Auth Service
  new KubeService(scope, 'traefikforwardauth-service', {
    metadata: {
      name: 'brunner56-traefikforwardauth',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'traefikforwardauth',
        'app.kubernetes.io/service': 'brunner56-traefikforwardauth',
        'helm.sh/chart': 'traefikforwardauth-1.0.0',
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
          port: 4181,
          targetPort: IntOrString.fromString('http')
        }
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'traefikforwardauth'
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster'
    }
  });
} 