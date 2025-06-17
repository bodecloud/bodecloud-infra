import { Construct } from 'constructs';
import { KubeService, IntOrString } from '../../imports/k8s';

/**
 * Creates StorageHub resources for the ElfHosted platform
 */
export function createStorageHubResources(scope: Construct): void {
  // Create StorageHub Service
  new KubeService(scope, 'storagehub-service', {
    metadata: {
      name: 'storagehub',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'storagehub',
        'app.kubernetes.io/service': 'storagehub',
        'helm.sh/chart': 'storagehub-1.0.0',
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
          port: 80,
          targetPort: IntOrString.fromString('http')
        }
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'storagehub'
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster'
    }
  });
} 