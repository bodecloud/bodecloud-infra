import { Construct } from 'constructs';
import { KubeDeployment, KubeService, IntOrString, Quantity } from '../../imports/k8s';

/**
 * Creates Kubernetes Dashboard resources
 */
export function createKubernetesDashboardResources(scope: Construct): void {
  // Kubernetes Dashboard Deployment
  new KubeDeployment(scope, 'kubernetesdashboard-deployment', {
    metadata: {
      name: 'brunner56-kubernetesdashboard',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/component': 'kubernetes-dashboard',
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'kubernetesdashboard',
        'app.kubernetes.io/version': '2.6.1',
        'helm.sh/chart': 'kubernetesdashboard-5.10.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          'app.kubernetes.io/component': 'kubernetes-dashboard',
          'app.kubernetes.io/instance': 'brunner56',
          'app.kubernetes.io/name': 'kubernetesdashboard',
        },
      },
      template: {
        metadata: {
          labels: {
            'app.kubernetes.io/component': 'kubernetes-dashboard',
            'app.kubernetes.io/instance': 'brunner56',
            'app.kubernetes.io/managed-by': 'Helm',
            'app.kubernetes.io/name': 'kubernetesdashboard',
            'app.kubernetes.io/version': '2.6.1',
            'helm.sh/chart': 'kubernetesdashboard-5.10.0',
          },
        },
        spec: {
          volumes: [
            {
              name: 'kubernetes-dashboard-certs',
              secret: {
                secretName: 'brunner56-kubernetesdashboard-certs',
                defaultMode: 420,
              },
            },
            {
              name: 'tmp-volume',
              emptyDir: {},
            },
          ],
          containers: [
            {
              name: 'kubernetesdashboard',
              image: 'kubernetesui/dashboard:v2.6.1',
              args: [
                '--namespace=bolabaden',
                '--sidecar-host=http://127.0.0.1:8000',
                '--enable-skip-login',
                '--enable-insecure-login',
                '--system-banner=Built with ❤️ by @funkypenguin and friends (join us!)',
              ],
              ports: [
                {
                  name: 'http',
                  containerPort: 9090,
                  protocol: 'TCP',
                },
              ],
              resources: {
                limits: {
                  cpu: Quantity.fromString('1'),
                  memory: Quantity.fromString('256Mi'),
                },
                requests: {
                  cpu: Quantity.fromString('0'),
                  memory: Quantity.fromString('64Mi'),
                },
              },
              volumeMounts: [
                {
                  name: 'kubernetes-dashboard-certs',
                  mountPath: '/certs',
                },
                {
                  name: 'tmp-volume',
                  mountPath: '/tmp',
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(9090),
                  scheme: 'HTTP',
                },
                initialDelaySeconds: 30,
                timeoutSeconds: 30,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3,
              },
              terminationMessagePath: '/dev/termination-log',
              terminationMessagePolicy: 'File',
              imagePullPolicy: 'IfNotPresent',
              securityContext: {
                runAsUser: 1001,
                runAsGroup: 2001,
                readOnlyRootFilesystem: true,
                allowPrivilegeEscalation: false,
              },
            },
            {
              name: 'dashboard-metrics-scraper',
              image: 'kubernetesui/metrics-scraper:v1.0.9',
              ports: [
                {
                  containerPort: 8000,
                  protocol: 'TCP',
                },
              ],
              resources: {},
              volumeMounts: [
                {
                  name: 'tmp-volume',
                  mountPath: '/tmp',
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/',
                  port: IntOrString.fromNumber(8000),
                  scheme: 'HTTP',
                },
                initialDelaySeconds: 30,
                timeoutSeconds: 30,
                periodSeconds: 10,
                successThreshold: 1,
                failureThreshold: 3,
              },
              terminationMessagePath: '/dev/termination-log',
              terminationMessagePolicy: 'File',
              imagePullPolicy: 'IfNotPresent',
              securityContext: {
                runAsUser: 1001,
                runAsGroup: 2001,
                readOnlyRootFilesystem: true,
                allowPrivilegeEscalation: false,
              },
            },
          ],
          restartPolicy: 'Always',
          terminationGracePeriodSeconds: 30,
          dnsPolicy: 'ClusterFirst',
          serviceAccountName: 'kubernetes-dashboard',
          serviceAccount: 'kubernetes-dashboard',
          securityContext: {
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
          maxSurge: IntOrString.fromNumber(0),
        },
      },
      revisionHistoryLimit: 10,
      progressDeadlineSeconds: 600,
    },
  });

  // Kubernetes Dashboard Service
  new KubeService(scope, 'kubernetesdashboard-service', {
    metadata: {
      name: 'brunner56-kubernetesdashboard',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/component': 'kubernetes-dashboard',
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'kubernetesdashboard',
        'app.kubernetes.io/service': 'brunner56-kubernetesdashboard',
        'helm.sh/chart': 'kubernetesdashboard-5.10.0',
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
          port: 9090,
          targetPort: IntOrString.fromNumber(9090),
        },
      ],
      selector: {
        'app.kubernetes.io/component': 'kubernetes-dashboard',
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'kubernetesdashboard',
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster',
    },
  });
} 