import { Construct } from 'constructs';
import { KubeNamespace } from '../../imports/k8s';

/**
 * Creates the namespace for ElfHosted resources
 */
export function createNamespace(scope: Construct): KubeNamespace {
  return new KubeNamespace(scope, 'bolabaden-namespace', {
    metadata: {
      name: 'bolabaden',
      labels: {
        name: 'bolabaden'
      }
    }
  });
} 