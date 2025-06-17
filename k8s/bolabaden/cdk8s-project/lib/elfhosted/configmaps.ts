import { Construct } from 'constructs';
import { KubeConfigMap } from '../../imports/k8s';

/**
 * Creates ConfigMaps for ElfHosted applications
 */
export function createConfigMaps(scope: Construct) {
  // Create B2-S3cmd-Restore ConfigMap
  new KubeConfigMap(scope, 'b2-s3cmd-restore-config', {
    metadata: {
      name: 'b2-s3cmd-restore-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    data: {
      's3cfg': `[default]

bucket_location = us-west-002

host_base = s3.us-west-002.backblazeb2.com

host_bucket = %(bucket)s.s3.us-west-002.backblazeb2.com`
    }
  });

  // Create Filebrowser Environment ConfigMap
  new KubeConfigMap(scope, 'filebrowser-env', {
    metadata: {
      name: 'filebrowser-env',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    data: {
      'FB_BASEURL': '',
      'FB_DATABASE': '/tmp/filebrowser.db',
      'FB_ROOT': '/storage'
    }
  });

  // Create Filebrowser Script ConfigMap
  new KubeConfigMap(scope, 'filebrowser-elfbot-script', {
    metadata: {
      name: 'filebrowser-elfbot-script',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    data: {
      'elfbot': `#!/bin/bash
# elfbot command dispatcher
# This script is run by the user via filebrowser commands
# It will create a file in the /elfbot directory to trigger a command

# Check if we have at least one parameter
if [ -z "$1" ]; then
  echo "elfbot: missing parameter"
  exit 1
fi

# Get the first parameter
CMD="$1"

# If there's a second parameter, add it to the command
if [ -n "$2" ]; then
  CMD="$CMD=$2"
fi

# Write the command to a file in /elfbot with the name of the service
echo "$CMD" > /elfbot/$1

echo "elfbot: command sent: $CMD"
`
    }
  });

  // Create Gatus Config ConfigMap
  new KubeConfigMap(scope, 'gatus-config', {
    metadata: {
      name: 'gatus-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    data: {
      'config.yaml': `endpoints:
  - name: Plex 
    group: Services
    url: https://brunner56-plex.bolabaden.svc.cluster.local:32400/web/index.html
    interval: 1m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: custom
        enabled: true
        description: "Plex is healthy"
        send-on-resolved: true
        failure-threshold: 3
        success-threshold: 1
        custom: {}
  - name: Filebrowser
    group: Services
    url: http://brunner56-filebrowser.bolabaden.svc.cluster.local:8080/
    interval: 1m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: custom
        enabled: true
        description: "Filebrowser is healthy"
        send-on-resolved: true
        failure-threshold: 3
        success-threshold: 1
        custom: {}
`
    }
  });

  // Create Homer Config ConfigMap
  new KubeConfigMap(scope, 'homer-config', {
    metadata: {
      name: 'homer-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    data: {
      'config.yml': `---
title: "ElfHosted Dashboard"
subtitle: "brunner56"
logo: "assets/logo.png"
icon: "fas fa-elf"
header: true
footer: '<p>Created with <span class="has-text-danger">❤️</span> by <a href="https://elfhosted.com">ElfHosted</a></p>'

settings:
  showGlobalHealthCheck: true
  useExternalJs: false

services:
  - name: "ElfHosted Services"
    icon: "fas fa-home"
    items:
      - name: "Filebrowser"
        subtitle: "File manager to access your folders"
        logo: "assets/tools/filebrowser.png"
        url: "#"
      - name: "Plex"
        subtitle: "Media server"
        logo: "assets/tools/plex.png"
        url: "#"

  - name: "System Status"
    icon: "fas fa-heartbeat"
    items:
      - name: "Health Checks"
        subtitle: "Health monitoring for your services"
        logo: "assets/tools/gatus.png"
        url: "#"
`,
      'custom.css': `/* Hide things we don't want to show */
.dashboard-item[href=""], .dashboard-item:not([href]) {
  display: none;
}`,
      'disk_usage.sh': `#!/bin/bash
# Serve a tiny web server using netcat to provide an HTTP endpoint for Homer
# This returns the available storage space for all the mounted volumes

PORT=\${PORT:-8080}
echo "Starting simple HTTP server on port \$PORT"

while true; do
  echo -ne "HTTP/1.0 200 OK\r\nContent-Type: text/html\r\n\r\n" | nc -l -p \$PORT -q 1
done
`
    }
  });

  // Create Plex Config ConfigMap
  new KubeConfigMap(scope, 'plex-config', {
    metadata: {
      name: 'plex-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    data: {
      'PLEX_CLAIM': '',
      'ADVERTISE_IP': 'https://brunner56-plex.elfhosted.com:443',
      'ALLOWED_NETWORKS': '10.0.0.0/8,172.16.0.0/12,192.168.0.0/16',
      'PLEX_PREFERENCE_1': 'FriendlyName=brunner56 by ElfHosted 🧝',
      'PLEX_PREFERENCE_2': 'EnableIPv6=0',
      'PLEX_PREFERENCE_3': 'logDebug=0',
      'PLEX_PREFERENCE_4': 'DisableTLSv1_0=1',
      'PLEX_PREFERENCE_5': 'secureConnections=1',
      'PLEX_PREFERENCE_6': 'ManualPortMappingMode=1',
      'PLEX_PREFERENCE_7': 'ManualPortMappingPort=443'
    }
  });

  // Create Recyclarr Config ConfigMap
  new KubeConfigMap(scope, 'recyclarr-config', {
    metadata: {
      name: 'recyclarr-config',
      namespace: 'bolabaden',
      labels: {
        'app.kubernetes.io/managed-by': 'Helm',
        'helm.toolkit.fluxcd.io/name': 'brunner56',
        'helm.toolkit.fluxcd.io/namespace': 'bolabaden'
      },
      annotations: {
        'meta.helm.sh/release-name': 'brunner56',
        'meta.helm.sh/release-namespace': 'bolabaden'
      }
    },
    data: {
      'recyclarr.yaml': `# Configuration specific to Sonarr
sonarr:
  series:
    base_url: http://sonarr:8989
    api_key: !env_var SONARR_API_KEY

    # Quality definitions from the guide to sync to Sonarr. Choices: series, anime
    quality_definition:
      type: series

    custom_formats:
      # A list of custom formats to sync to Sonarr.
      # Use 'recyclarr list custom-formats sonarr' for values you can put here.
      # https://recyclarr.dev/wiki/yaml/config-reference/custom-formats/
      - trash_ids:
          # Audio
          - 496f355514737f7d83bf7aa4d24f8169 # TrueHD Atmos
          - 2f22d89048b01681dde8afe73aabb0b9 # DTS-X
          - 417804f7f2c4308c1f4c5d380d4c4475 # ATMOS (undefined)
          - 1af239278386be2919e1bcee0bde047e # DD+ ATMOS
          - 3cafb66171b47f226146a0770576870f # TrueHD
          - dcf3ec6938fa32445f590a4da84256cd # DTS-HD MA
          - a570d4a0e56a2874b64e5bfa55202a1b # FLAC
          - e7c2fcae07cbada050a0af3357491d7b # PCM
          - 8e109e50e0a0b83b5098b056e13bf6db # DTS-HD HRA
          - 185f1dd7264c4562b9022d963ac37424 # DD+
          - f9f847ac70a0af62ea4a08280b859636 # DTS-ES
          - 1c1a4c5e823891c75bc50380a6866f73 # DTS
          - 240770601cc226190c367ef59aba7463 # AAC
          - c2998bd0d90ed5621d8df281e839436e # DD
        quality_profiles:
          - name: WEB-1080p
            score: 0 # Adjust scoring as desired

# Configuration specific to Radarr
radarr:
  movies:
    base_url: http://radarr:7878
    api_key: !env_var RADARR_API_KEY

    # Quality definitions from the guide to sync to Radarr. Choices: movie, sqp-1-1080p, sqp-1-2160p
    quality_definition:
      type: movie

    custom_formats:
      # A list of custom formats to sync to Radarr.
      # Use 'recyclarr list custom-formats radarr' for values you can put here.
      # https://recyclarr.dev/wiki/yaml/config-reference/custom-formats/
      - trash_ids:
          # Audio
          - 496f355514737f7d83bf7aa4d24f8169 # TrueHD Atmos
          - 2f22d89048b01681dde8afe73aabb0b9 # DTS-X
          - 417804f7f2c4308c1f4c5d380d4c4475 # ATMOS (undefined)
          - 1af239278386be2919e1bcee0bde047e # DD+ ATMOS
          - 3cafb66171b47f226146a0770576870f # TrueHD
          - dcf3ec6938fa32445f590a4da84256cd # DTS-HD MA
          - a570d4a0e56a2874b64e5bfa55202a1b # FLAC
          - e7c2fcae07cbada050a0af3357491d7b # PCM
          - 8e109e50e0a0b83b5098b056e13bf6db # DTS-HD HRA
          - 185f1dd7264c4562b9022d963ac37424 # DD+
          - f9f847ac70a0af62ea4a08280b859636 # DTS-ES
          - 1c1a4c5e823891c75bc50380a6866f73 # DTS
          - 240770601cc226190c367ef59aba7463 # AAC
          - c2998bd0d90ed5621d8df281e839436e # DD
        quality_profiles:
          - name: HD-1080p
            score: 0 # Adjust scoring as desired
`
    }
  });
} 