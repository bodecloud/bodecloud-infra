import { Construct } from 'constructs';
import { KubePersistentVolumeClaim, Quantity } from '../../imports/k8s';

/**
 * Creates Persistent Volume Claims for ElfHosted applications
 */
export function createStorageResources(scope: Construct) {
  // Create Backup PVC
  new KubePersistentVolumeClaim(scope, 'backup-pvc', {
    metadata: {
      name: 'backup',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'helm.sh/resource-policy': 'keep',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('20Gi')
        }
      },
      storageClassName: 'local-path',
      volumeMode: 'Filesystem'
    }
  });

  // Create Config PVC
  new KubePersistentVolumeClaim(scope, 'config-pvc', {
    metadata: {
      name: 'config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'helm.sh/resource-policy': 'keep',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'local-path',
      volumeMode: 'Filesystem'
    }
  });

  // Create Logs PVC
  new KubePersistentVolumeClaim(scope, 'logs-pvc', {
    metadata: {
      name: 'logs',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'helm.sh/resource-policy': 'keep',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'local-path',
      volumeMode: 'Filesystem'
    }
  });

  // Create Rclone PVC
  new KubePersistentVolumeClaim(scope, 'rclone-pvc', {
    metadata: {
      name: 'rclone',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'helm.sh/resource-policy': 'keep',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'local-path',
      volumeMode: 'Filesystem'
    }
  });

  // Create Symlinks PVC
  new KubePersistentVolumeClaim(scope, 'symlinks-pvc', {
    metadata: {
      name: 'symlinks',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'helm.sh/resource-policy': 'keep',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'local-path',
      volumeMode: 'Filesystem'
    }
  });

  // Create RealDebrid-Zurg PVC
  new KubePersistentVolumeClaim(scope, 'realdebrid-zurg-pvc', {
    metadata: {
      name: 'realdebrid-zurg',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'helm.sh/resource-policy': 'keep',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'local-path',
      volumeMode: 'Filesystem'
    }
  });

  // Create Transcode PVC
  new KubePersistentVolumeClaim(scope, 'transcode-1g-pvc', {
    metadata: {
      name: 'transcode-1g',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'helm.sh/resource-policy': 'keep',
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    spec: {
      accessModes: ['ReadWriteOnce'],
      resources: {
        requests: {
          storage: Quantity.fromString('1Gi')
        }
      },
      storageClassName: 'local-path',
      volumeMode: 'Filesystem'
    }
  });
} 