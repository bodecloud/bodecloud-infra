import { Construct } from 'constructs';
import { KubeConfigMap } from '../../imports/k8s';

/**
 * Creates additional ConfigMaps for the ElfHosted platform
 */
export function createAdditionalConfigMaps(scope: Construct): void {
  // RcloneUI ConfigMap
  new KubeConfigMap(scope, 'rcloneui-configmap', {
    metadata: {
      name: 'rcloneui-configmap',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rcloneui',
        'helm.sh/chart': 'rcloneui-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'rclone.conf': `[premiumize]
type = premiumizeme
user = your-username
pass = your-password

[realdebrid]
type = realdebrid
api_key = your-api-key

[local]
type = local
`
    }
  });

  // Traefik Forward Auth Config ConfigMap
  new KubeConfigMap(scope, 'traefik-forward-auth-config-configmap', {
    metadata: {
      name: 'traefik-forward-auth-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'traefikforwardauth',
        'helm.sh/chart': 'traefikforwardauth-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'config.ini': `[generic-oauth]
oauth-uri = https://auth.elfhosted.com
oauth-path = /oauth
oauth-logout-path = /logout
client-id = your-client-id
client-secret = your-client-secret
scope = openid profile email
token-style = header
login-url = https://auth.elfhosted.com/login
logout-url = https://auth.elfhosted.com/logout
`
    }
  });

  // Tooling Scripts ConfigMap
  new KubeConfigMap(scope, 'tooling-scripts-configmap', {
    metadata: {
      name: 'tooling-scripts',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'tooling',
        'helm.sh/chart': 'tooling-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'backup.sh': `#!/bin/bash
# Backup script for ElfHosted
echo "Starting backup..."
# Add backup logic here
echo "Backup completed."
`,
      'restore.sh': `#!/bin/bash
# Restore script for ElfHosted
echo "Starting restore..."
# Add restore logic here
echo "Restore completed."
`,
      'maintenance.sh': `#!/bin/bash
# Maintenance script for ElfHosted
echo "Starting maintenance..."
# Add maintenance logic here
echo "Maintenance completed."
`
    }
  });

  // Mattermost Backup ConfigMap
  new KubeConfigMap(scope, 'mattermost-backup-configmap', {
    metadata: {
      name: 'mattermost-backup',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'mattermost',
        'helm.sh/chart': 'mattermost-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'backup.sh': `#!/bin/bash
# Mattermost backup script
echo "Starting Mattermost backup..."
# Add Mattermost backup logic here
echo "Mattermost backup completed."
`
    }
  });

  // Notifiarr Config ConfigMap
  new KubeConfigMap(scope, 'notifiarr-config-configmap', {
    metadata: {
      name: 'notifiarr-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'notifiarr',
        'helm.sh/chart': 'notifiarr-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'notifiarr.conf': `# Notifiarr Configuration
api_key = your-api-key
bind_addr = 0.0.0.0:5454
log_level = info
timeout = 10s

# Service configurations
[services]
sonarr = true
radarr = true
lidarr = true
readarr = true
prowlarr = true
`
    }
  });

  // Rutorrent Config ConfigMap
  new KubeConfigMap(scope, 'rutorrent-config-configmap', {
    metadata: {
      name: 'rutorrent-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'rutorrent',
        'helm.sh/chart': 'rutorrent-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'rtorrent.rc': `# rtorrent configuration
directory.default.set = /downloads
session.path.set = /session
network.port_range.set = 6881-6999
network.port_random.set = yes
dht.mode.set = auto
protocol.pex.set = yes
trackers.use_udp.set = yes
encryption = allow_incoming,try_outgoing,enable_retry
network.xmlrpc.size_limit.set = 4M
`
    }
  });

  // Samba Config ConfigMap
  new KubeConfigMap(scope, 'samba-config-configmap', {
    metadata: {
      name: 'samba-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'samba',
        'helm.sh/chart': 'samba-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'smb.conf': `[global]
workgroup = WORKGROUP
server string = ElfHosted Samba Server
security = user
map to guest = bad user
dns proxy = no

[media]
path = /media
browseable = yes
writable = yes
guest ok = yes
read only = no
`
    }
  });

  // Seafile Config ConfigMap
  new KubeConfigMap(scope, 'seafile-config-configmap', {
    metadata: {
      name: 'seafile-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'seafile',
        'helm.sh/chart': 'seafile-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'seafile.conf': `[database]
type = mysql
host = mysql
port = 3306
user = seafile
password = seafile
db_name = seafile

[fileserver]
port = 8082
host = 0.0.0.0

[seahub]
port = 8000
host = 0.0.0.0
`
    }
  });

  // B2 S3cmd Restore Config ConfigMap
  new KubeConfigMap(scope, 'b2-s3cmd-restore-config-configmap', {
    metadata: {
      name: 'b2-s3cmd-restore-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'b2-restore',
        'helm.sh/chart': 'b2-restore-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      's3cfg': `[default]
access_key = your-access-key
secret_key = your-secret-key
host_base = s3.us-west-002.backblazeb2.com
host_bucket = %(bucket)s.s3.us-west-002.backblazeb2.com
use_https = True
`
    }
  });

  // Kubernetes Dashboard Settings ConfigMap
  new KubeConfigMap(scope, 'kubernetes-dashboard-settings-configmap', {
    metadata: {
      name: 'kubernetes-dashboard-settings',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/instance': 'brunner56',
        'app.kubernetes.io/managed-by': 'Helm',
        'app.kubernetes.io/name': 'kubernetes-dashboard',
        'helm.sh/chart': 'kubernetes-dashboard-1.0.0',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden',
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden',
      },
    },
    data: {
      'settings.json': `{
  "clusterName": "ElfHosted Cluster",
  "itemsPerPage": 10,
  "labelsLimit": 3,
  "logsAutoRefreshTimeInterval": 5,
  "resourceAutoRefreshTimeInterval": 5,
  "disableAccessDeniedNotifications": false
}`
    }
  });
} 