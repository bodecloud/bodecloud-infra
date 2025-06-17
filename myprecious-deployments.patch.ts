/**
 * This file contains patches to fix the CDK8s project
 * 
 * Key improvements:
 * 1. Use local-path storage class instead of ceph-block-ssd
 * 2. Use alternative image sources instead of Docker Hub to avoid rate limits
 * 3. Reduce excessive resource requests and limits
 * 4. Set reasonable replica counts to avoid pod sprawl
 * 
 * Apply these changes to your CDK8s project and redeploy for a stable cluster
 */

// In MyPreciousChart class, update the storage class for all volumeClaimTemplates
// Change 'ceph-block-ssd' to 'local-path' in all PVC definitions

// Update image sources to avoid Docker Hub rate limits
// Example fix for prom/node-exporter:
// new DaemonSet(this, 'node-exporter', {
//   containers: [{
//     name: 'node-exporter',
//     image: 'quay.io/prometheus/node-exporter:latest', // Using Quay.io instead of Docker Hub
//     resources: {
//       limits: {
//         cpu: Quantity.fromString('100m'),
//         memory: Quantity.fromString('128Mi')
//       }
//     }
//   }]
// });

// Lower resource requests/limits for deployments
// For most services, use something like:
// resources: {
//   requests: {
//     cpu: Quantity.fromString('10m'),
//     memory: Quantity.fromString('64Mi')
//   },
//   limits: {
//     cpu: Quantity.fromString('200m'),
//     memory: Quantity.fromString('256Mi')
//   }
// }

// Ensure replica count is always set to 1 unless explicitly needed:
// replicas: 1,
