import { Construct } from 'constructs';
import { KubePersistentVolumeClaim, Quantity } from '../../imports/k8s';

/**
 * Creates additional storage resources (PVCs) for the ElfHosted platform
 */
export function createAdditionalStorageResources(scope: Construct): void {
  // Transcode PVC
  new KubePersistentVolumeClaim(scope, 'transcode-1g-pvc', {
    metadata: {
      name: 'transcode-1g',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'transcode',
        'helm.sh/chart': 'transcode-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RD Backup Cache PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rd-backup-cache-pvc', {
    metadata: {
      name: 'volsync-rd-backup-cache',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RD Backup Dest PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rd-backup-dest-pvc', {
    metadata: {
      name: 'volsync-rd-backup-dest',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RD Config Cache PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rd-config-cache-pvc', {
    metadata: {
      name: 'volsync-rd-config-cache',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RD Config Dest PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rd-config-dest-pvc', {
    metadata: {
      name: 'volsync-rd-config-dest',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RD Symlinks Cache PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rd-symlinks-cache-pvc', {
    metadata: {
      name: 'volsync-rd-symlinks-cache',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RD Symlinks Dest PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rd-symlinks-dest-pvc', {
    metadata: {
      name: 'volsync-rd-symlinks-dest',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RS Backup Cache PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rs-backup-cache-pvc', {
    metadata: {
      name: 'volsync-rs-backup-cache',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RS Config Cache PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rs-config-cache-pvc', {
    metadata: {
      name: 'volsync-rs-config-cache',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });

  // VolSync RS Symlinks Cache PVC
  new KubePersistentVolumeClaim(scope, 'volsync-rs-symlinks-cache-pvc', {
    metadata: {
      name: 'volsync-rs-symlinks-cache',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'volsync',
        'helm.sh/chart': 'volsync-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'ceph-block',
      volumeMode: 'Filesystem'
    }
  });
} 