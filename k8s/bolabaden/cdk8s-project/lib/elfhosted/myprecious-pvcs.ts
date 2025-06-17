import { Construct } from 'constructs';
import { KubePersistentVolumeClaim, Quantity } from '../../imports/k8s';

/**
 * Creates all Persistent Volume Claims for the myprecious chart
 * This includes all storage requirements from the original Helm chart
 */

interface PvcConfig {
  name: string;
  size: string;
  accessMode?: string;
  storageClass?: string;
}

const PVC_CONFIGS: PvcConfig[] = [
  // Core storage
  { name: 'config', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'logs', size: '5Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'backup', size: '50Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },

  // Rclone mounts
  { name: 'rclone', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'symlinks', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },

  // Debrid services storage
  { name: 'realdebrid-zurg', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'alldebrid', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'torbox', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'debridav', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'debridlink', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'premiumize', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },

  // Media covers and metadata
  { name: 'mediacovers', size: '5Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },

  // Transcoding storage
  { name: 'transcode-1g', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'transcode-50g', size: '50Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },

  // VolSync backup destinations (if volsync is enabled)
  { name: 'volsync-rd-backup-cache', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'volsync-rd-backup-dest', size: '50Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'volsync-rd-config-cache', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'volsync-rd-config-dest', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'volsync-rd-symlinks-cache', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'volsync-rd-symlinks-dest', size: '10Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },

  // VolSync replication sources (if volsync is enabled)
  { name: 'volsync-rs-backup-cache', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'volsync-rs-config-cache', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
  { name: 'volsync-rs-symlinks-cache', size: '1Gi', accessMode: 'ReadWriteOnce', storageClass: 'local-path' },
];

export function createAllMyPreciousPVCs(scope: Construct) {
  PVC_CONFIGS.forEach(config => {
    new KubePersistentVolumeClaim(scope, `pvc-${config.name}`, {
      metadata: {
        name: config.name,
        labels: {
          'app.kubernetes.io/instance': 'myprecious',
          'app.kubernetes.io/managed-by': 'cdk8s',
        },
      },
      spec: {
        accessModes: [config.accessMode || 'ReadWriteOnce'],
        storageClassName: config.storageClass || 'local-path',
        resources: {
          requests: {
            storage: Quantity.fromString(config.size),
          },
        },
      },
    });
  });
}

/**
 * Creates additional PVCs for specific services that need dedicated storage
 */
export function createServiceSpecificPVCs(scope: Construct) {
  // Application-specific PVCs that might be needed
  const serviceSpecificPVCs = [
    // Database storage for services that need it
    { name: 'postgresql-data', size: '20Gi', storageClass: 'local-path' },
    { name: 'redis-data', size: '5Gi', storageClass: 'local-path' },
    { name: 'mongodb-data', size: '20Gi', storageClass: 'local-path' },

    // Large storage for specific applications
    { name: 'jellyfin-cache', size: '10Gi', storageClass: 'local-path' },
    { name: 'plex-cache', size: '10Gi', storageClass: 'local-path' },
    { name: 'emby-cache', size: '10Gi', storageClass: 'local-path' },

    // Download client storage
    { name: 'qbittorrent-downloads', size: '100Gi', storageClass: 'local-path' },
    { name: 'deluge-downloads', size: '100Gi', storageClass: 'local-path' },
    { name: 'rutorrent-downloads', size: '100Gi', storageClass: 'local-path' },
    { name: 'sabnzbd-downloads', size: '100Gi', storageClass: 'local-path' },
    { name: 'nzbget-downloads', size: '100Gi', storageClass: 'local-path' },
  ];

  serviceSpecificPVCs.forEach(config => {
    new KubePersistentVolumeClaim(scope, `pvc-${config.name}`, {
      metadata: {
        name: config.name,
        labels: {
          'app.kubernetes.io/instance': 'myprecious',
          'app.kubernetes.io/managed-by': 'cdk8s',
          'storage.elfhosted.com/type': 'service-specific',
        },
      },
      spec: {
        accessModes: ['ReadWriteOnce'],
        storageClassName: config.storageClass,
        resources: {
          requests: {
            storage: Quantity.fromString(config.size),
          },
        },
      },
    });
  });
}
