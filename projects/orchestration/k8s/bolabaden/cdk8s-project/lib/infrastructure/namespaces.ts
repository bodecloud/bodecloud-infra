import { Construct } from 'constructs';
import { KubeNamespace } from '../../imports/k8s';

/**
 * Creates namespaces for the media stack
 */
export function createNamespaces(scope: Construct): void {
  // Main media stack namespace
  new KubeNamespace(scope, 'my-media-stack-namespace', {
    metadata: {
      name: 'my-media-stack',
      labels: {
        name: 'my-media-stack',
        'app.kubernetes.io/managed-by': 'cdk8s',
        'app.kubernetes.io/part-of': 'media-stack',
      },
    },
  });

  // VPN Gateway namespace
  new KubeNamespace(scope, 'vpn-gateway-namespace', {
    metadata: {
      name: 'vpn-gateway',
      labels: {
        name: 'vpn-gateway',
        'app.kubernetes.io/managed-by': 'cdk8s',
        'app.kubernetes.io/part-of': 'media-stack',
      },
    },
  });

  // Monitoring namespace
  new KubeNamespace(scope, 'monitoring-namespace', {
    metadata: {
      name: 'monitoring',
      labels: {
        'app.kubernetes.io/managed-by': 'cdk8s',
        'app.kubernetes.io/part-of': 'media-stack',
      },
    },
  });

  // Development namespace
  new KubeNamespace(scope, 'development-namespace', {
    metadata: {
      name: 'development',
      labels: {
        'app.kubernetes.io/managed-by': 'cdk8s',
        'app.kubernetes.io/part-of': 'media-stack',
      },
    },
  });

  // AI services namespace
  new KubeNamespace(scope, 'ai-services-namespace', {
    metadata: {
      name: 'ai-services',
      labels: {
        'app.kubernetes.io/managed-by': 'cdk8s',
        'app.kubernetes.io/part-of': 'media-stack',
      },
    },
  });
} 