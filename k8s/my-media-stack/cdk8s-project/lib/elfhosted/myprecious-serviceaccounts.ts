import { Construct } from 'constructs';
import { KubeServiceAccount, KubeRole, KubeRoleBinding, KubeClusterRole, KubeClusterRoleBinding } from '../../imports/k8s';

/**
 * Creates all ServiceAccounts, Roles, and RoleBindings for the myprecious chart
 * This includes service accounts for various applications and their required permissions
 */

interface ServiceAccountConfig {
  name: string;
  annotations?: Record<string, string>;
  labels?: Record<string, string>;
  automountServiceAccountToken?: boolean;
  imagePullSecrets?: string[];
  secrets?: string[];
}

const SERVICE_ACCOUNT_CONFIGS: ServiceAccountConfig[] = [
  // Core service accounts
  {
    name: 'default',
    automountServiceAccountToken: false,
  },
  {
    name: 'myprecious-default',
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/managed-by': 'cdk8s',
    },
  },

  // Application-specific service accounts
  {
    name: 'plex',
    labels: {
      'app.elfhosted.com/name': 'plex',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'radarr',
    labels: {
      'app.elfhosted.com/name': 'radarr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'sonarr',
    labels: {
      'app.elfhosted.com/name': 'sonarr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'prowlarr',
    labels: {
      'app.elfhosted.com/name': 'prowlarr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'lidarr',
    labels: {
      'app.elfhosted.com/name': 'lidarr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'readarr',
    labels: {
      'app.elfhosted.com/name': 'readarr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'jellyfin',
    labels: {
      'app.elfhosted.com/name': 'jellyfin',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'overseerr',
    labels: {
      'app.elfhosted.com/name': 'overseerr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'jellyseerr',
    labels: {
      'app.elfhosted.com/name': 'jellyseerr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'tautulli',
    labels: {
      'app.elfhosted.com/name': 'tautulli',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'qbittorrent',
    labels: {
      'app.elfhosted.com/name': 'qbittorrent',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'filebrowser',
    labels: {
      'app.elfhosted.com/name': 'filebrowser',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'homer',
    labels: {
      'app.elfhosted.com/name': 'homer',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'gatus',
    labels: {
      'app.elfhosted.com/name': 'gatus',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'wizarr',
    labels: {
      'app.elfhosted.com/name': 'wizarr',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'zurg',
    labels: {
      'app.elfhosted.com/name': 'zurg',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'riven',
    labels: {
      'app.elfhosted.com/name': 'riven',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'riven-frontend',
    labels: {
      'app.elfhosted.com/name': 'riven-frontend',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },

  // Backup and maintenance service accounts
  {
    name: 'backup-operator',
    labels: {
      'app.kubernetes.io/component': 'backup',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'maintenance-operator',
    labels: {
      'app.kubernetes.io/component': 'maintenance',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },

  // Monitoring service accounts
  {
    name: 'monitoring-operator',
    labels: {
      'app.kubernetes.io/component': 'monitoring',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },

  // Database service accounts
  {
    name: 'postgresql',
    labels: {
      'app.elfhosted.com/name': 'postgresql',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'redis',
    labels: {
      'app.elfhosted.com/name': 'redis',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
  {
    name: 'mongodb',
    labels: {
      'app.elfhosted.com/name': 'mongodb',
      'app.kubernetes.io/instance': 'myprecious',
    },
  },
];

export function createAllMyPreciousServiceAccounts(scope: Construct) {
  SERVICE_ACCOUNT_CONFIGS.forEach(config => {
    new KubeServiceAccount(scope, `sa-${config.name}`, {
      metadata: {
        name: config.name,
        ...(config.labels && { labels: config.labels }),
        ...(config.annotations && { annotations: config.annotations }),
      },
      ...(config.automountServiceAccountToken !== undefined && {
        automountServiceAccountToken: config.automountServiceAccountToken,
      }),
      ...(config.imagePullSecrets && {
        imagePullSecrets: config.imagePullSecrets.map(secret => ({ name: secret })),
      }),
      ...(config.secrets && {
        secrets: config.secrets.map(secret => ({ name: secret })),
      }),
    });
  });
}

/**
 * Creates roles and role bindings for service accounts that need special permissions
 */
export function createServiceAccountRoles(scope: Construct) {
  // Backup operator role
  new KubeRole(scope, 'backup-operator-role', {
    metadata: {
      name: 'backup-operator',
    },
    rules: [
      {
        apiGroups: [''],
        resources: ['persistentvolumeclaims', 'persistentvolumes'],
        verbs: ['get', 'list', 'watch'],
      },
      {
        apiGroups: [''],
        resources: ['pods'],
        verbs: ['get', 'list', 'watch', 'create', 'delete'],
      },
      {
        apiGroups: ['batch'],
        resources: ['jobs'],
        verbs: ['get', 'list', 'watch', 'create', 'delete'],
      },
    ],
  });

  new KubeRoleBinding(scope, 'backup-operator-binding', {
    metadata: {
      name: 'backup-operator',
    },
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      kind: 'Role',
      name: 'backup-operator',
    },
    subjects: [{
      kind: 'ServiceAccount',
      name: 'backup-operator',
    }],
  });

  // Maintenance operator role
  new KubeRole(scope, 'maintenance-operator-role', {
    metadata: {
      name: 'maintenance-operator',
    },
    rules: [
      {
        apiGroups: [''],
        resources: ['pods'],
        verbs: ['get', 'list', 'watch', 'delete'],
      },
      {
        apiGroups: ['apps'],
        resources: ['deployments', 'replicasets'],
        verbs: ['get', 'list', 'watch', 'patch'],
      },
      {
        apiGroups: ['batch'],
        resources: ['jobs', 'cronjobs'],
        verbs: ['get', 'list', 'watch', 'create', 'delete'],
      },
    ],
  });

  new KubeRoleBinding(scope, 'maintenance-operator-binding', {
    metadata: {
      name: 'maintenance-operator',
    },
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      kind: 'Role',
      name: 'maintenance-operator',
    },
    subjects: [{
      kind: 'ServiceAccount',
      name: 'maintenance-operator',
    }],
  });

  // Monitoring operator role
  new KubeRole(scope, 'monitoring-operator-role', {
    metadata: {
      name: 'monitoring-operator',
    },
    rules: [
      {
        apiGroups: [''],
        resources: ['pods', 'services', 'endpoints'],
        verbs: ['get', 'list', 'watch'],
      },
      {
        apiGroups: ['apps'],
        resources: ['deployments', 'replicasets', 'daemonsets', 'statefulsets'],
        verbs: ['get', 'list', 'watch'],
      },
      {
        apiGroups: ['metrics.k8s.io'],
        resources: ['pods', 'nodes'],
        verbs: ['get', 'list'],
      },
    ],
  });

  new KubeRoleBinding(scope, 'monitoring-operator-binding', {
    metadata: {
      name: 'monitoring-operator',
    },
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      kind: 'Role',
      name: 'monitoring-operator',
    },
    subjects: [{
      kind: 'ServiceAccount',
      name: 'monitoring-operator',
    }],
  });
}

/**
 * Creates cluster-level roles for service accounts that need cluster-wide permissions
 */
export function createClusterRoles(scope: Construct) {
  // Node finder cluster role (for the nodefinder service)
  new KubeClusterRole(scope, 'nodefinder-cluster-role', {
    metadata: {
      name: 'nodefinder',
    },
    rules: [
      {
        apiGroups: [''],
        resources: ['nodes'],
        verbs: ['get', 'list', 'watch'],
      },
      {
        apiGroups: [''],
        resources: ['pods'],
        verbs: ['get', 'list', 'watch', 'delete'],
      },
    ],
  });

  new KubeClusterRoleBinding(scope, 'nodefinder-cluster-binding', {
    metadata: {
      name: 'nodefinder',
    },
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      kind: 'ClusterRole',
      name: 'nodefinder',
    },
    subjects: [
      {
        kind: 'ServiceAccount',
        name: 'nodefinder-deleter',
        namespace: 'default', // This should be the actual namespace
      },
      {
        kind: 'ServiceAccount',
        name: 'nodefinder-waiter',
        namespace: 'default', // This should be the actual namespace
      },
    ],
  });
} 