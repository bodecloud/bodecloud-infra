import { Construct } from 'constructs';
import { KubeService, IntOrString } from '../../imports/k8s';

/**
 * Creates ElfTerm resources for the ElfHosted platform
 */
export function createElfTermResources(scope: Construct): void {
  // Create ElfTerm Service
  new KubeService(scope, 'elfterm-service', {
    metadata: {
      name: 'elfterm',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'elfterm',
        'app.kubernetes.io/service': 'elfterm',
        'helm.sh/chart': 'elfterm-1.0.0',
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
          port: 7681,
          targetPort: IntOrString.fromString('http')
        }
      ],
      selector: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/name': 'elfterm'
      },
      type: 'ClusterIP',
      sessionAffinity: 'None',
      ipFamilies: ['IPv4'],
      ipFamilyPolicy: 'SingleStack',
      internalTrafficPolicy: 'Cluster'
    }
  });
} 