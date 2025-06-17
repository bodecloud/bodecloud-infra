import { Construct } from 'constructs';
import { KubeResourceQuota, KubeLimitRange, Quantity } from '../../imports/k8s';

/**
 * Creates resource quotas and limit ranges for the media stack namespaces
 */
export function createResourceQuotas(scope: Construct): void {
  // Resource Quota for my-media-stack namespace
  new KubeResourceQuota(scope, 'my-media-stack-quota', {
    metadata: {
      name: 'my-media-stack-quota',
      namespace: 'my-media-stack',
    },
    spec: {
      hard: {
        // Memory limits
        'requests.memory': Quantity.fromString('16Gi'),
        'limits.memory': Quantity.fromString('20Gi'),

        // CPU limits  
        'requests.cpu': Quantity.fromString('8000m'),
        'limits.cpu': Quantity.fromString('12000m'),

        // Pod limits
        'pods': Quantity.fromString('50'),

        // Storage limits
        'persistentvolumeclaims': Quantity.fromString('20'),
        'requests.storage': Quantity.fromString('100Gi'),
      },
    },
  });

  // Resource Quota for vpn-gateway namespace
  new KubeResourceQuota(scope, 'vpn-gateway-quota', {
    metadata: {
      name: 'vpn-gateway-quota',
      namespace: 'vpn-gateway',
    },
    spec: {
      hard: {
        // Memory limits
        'requests.memory': Quantity.fromString('1.5Gi'),
        'limits.memory': Quantity.fromString('3Gi'),

        // CPU limits
        'requests.cpu': Quantity.fromString('800m'),
        'limits.cpu': Quantity.fromString('1600m'),

        // Pod limits
        'pods': Quantity.fromString('10'),
      },
    },
  });

  // Limit Range for my-media-stack namespace
  new KubeLimitRange(scope, 'my-media-stack-limits', {
    metadata: {
      name: 'my-media-stack-limits',
      namespace: 'my-media-stack',
    },
    spec: {
      limits: [
        {
          default: {
            memory: Quantity.fromString('512Mi'),
            cpu: Quantity.fromString('500m'),
          },
          defaultRequest: {
            memory: Quantity.fromString('128Mi'),
            cpu: Quantity.fromString('100m'),
          },
          max: {
            memory: Quantity.fromString('4Gi'),
            cpu: Quantity.fromString('2000m'),
          },
          min: {
            memory: Quantity.fromString('32Mi'),
            cpu: Quantity.fromString('10m'),
          },
          type: 'Container',
        },
        {
          default: {
            storage: Quantity.fromString('10Gi'),
          },
          max: {
            storage: Quantity.fromString('100Gi'),
          },
          min: {
            storage: Quantity.fromString('1Gi'),
          },
          type: 'PersistentVolumeClaim',
        },
      ],
    },
  });

  // Limit Range for vpn-gateway namespace
  new KubeLimitRange(scope, 'vpn-gateway-limits', {
    metadata: {
      name: 'vpn-gateway-limits',
      namespace: 'vpn-gateway',
    },
    spec: {
      limits: [
        {
          default: {
            memory: Quantity.fromString('512Mi'),
            cpu: Quantity.fromString('200m'),
          },
          defaultRequest: {
            memory: Quantity.fromString('128Mi'),
            cpu: Quantity.fromString('50m'),
          },
          max: {
            memory: Quantity.fromString('1Gi'),
            cpu: Quantity.fromString('500m'),
          },
          min: {
            memory: Quantity.fromString('64Mi'),
            cpu: Quantity.fromString('10m'),
          },
          type: 'Container',
        },
      ],
    },
  });
} 