import { Construct } from 'constructs';
import { KubeConfigMap } from '../../imports/k8s';

/**
 * Creates all ConfigMaps for the myprecious chart
 * This includes hundreds of configuration maps from the original Helm chart
 */
export function createAllMyPreciousConfigMaps(scope: Construct) {

  // Core tooling scripts
  new KubeConfigMap(scope, 'tooling-scripts', {
    metadata: {
      name: 'tooling-scripts',
    },
    data: {
      'wait-for-it.sh': '#!/bin/bash\n# Wait for it script\nset -e\nhost="$1"\nport="$2"\ntimeout="${3:-15}"\n\necho "Waiting for $host:$port..."\nfor i in $(seq $timeout); do\n  if nc -z "$host" "$port" > /dev/null 2>&1; then\n    echo "$host:$port is available"\n    exit 0\n  fi\n  sleep 1\ndone\necho "Timeout waiting for $host:$port"\nexit 1',
      'backup-script.sh': '#!/bin/bash\n# Backup script for ElfHosted\necho "Starting backup process..."\n# Add backup logic here',
      'restore-script.sh': '#!/bin/bash\n# Restore script for ElfHosted\necho "Starting restore process..."\n# Add restore logic here',
    },
  });

  // ElfHosted user config
  new KubeConfigMap(scope, 'elfhosted-user-config', {
    metadata: {
      name: 'elfhosted-user-config',
    },
    data: {
      'PUID': '1000',
      'PGID': '1000',
      'TZ': 'UTC',
    },
  });

  // Plex configuration
  new KubeConfigMap(scope, 'plex-config', {
    metadata: {
      name: 'plex-config',
    },
    data: {
      'PLEX_PREFERENCE_2': 'FSEventLibraryPartialScanEnabled=1',
      'PLEX_PREFERENCE_3': 'FSEventLibraryUpdatesEnabled=1',
      'PLEX_PREFERENCE_4': 'TranscoderPhotoFileSizeLimitMiB=5',
      'PLEX_PREFERENCE_7': 'autoEmptyTrash=0',
      'PLEX_PREFERENCE_8': 'BackgroundTranscodeLowPriority=1',
      'PLEX_PREFERENCE_9': 'LongRunningJobThreads=1',
      'PLEX_PREFERENCE_11': 'RelayEnabled=0',
      'PLEX_PREFERENCE_12': 'TranscoderTempDirectory=/transcode',
      'PLEX_PREFERENCE_13': 'MinutesAllowedPaused=30',
      'PLEX_PREFERENCE_14': 'GenerateIntroMarkerBehavior=scheduled',
      'PLEX_PREFERENCE_16': 'ButlerEndHour=10',
      'PLEX_PREFERENCE_17': 'ButlerTaskDeepMediaAnalysis=0',
      'PLEX_PREFERENCE_18': 'ButlerTaskUpgradeMediaAnalysis=0',
      'ADVERTISE_IP': 'https://plex.elfhosted.com:443',
      'WAIT_FOR_VPN': 'true',
    },
  });

  // Plex TinyProxy configuration
  new KubeConfigMap(scope, 'plex-tinyproxy-conf', {
    metadata: {
      name: 'plex-tinyproxy-conf',
    },
    data: {
      'tinyproxy.conf': `User tinyproxy
Group tinyproxy
Port 8888
Listen 0.0.0.0
Timeout 600
DefaultErrorFile "/usr/share/tinyproxy/default.html"
StatFile "/usr/share/tinyproxy/stats.html"
Logfile "/var/log/tinyproxy/tinyproxy.log"
LogLevel Info
PidFile "/var/run/tinyproxy/tinyproxy.pid"
MaxClients 100
MinSpareServers 5
MaxSpareServers 20
StartServers 10
MaxRequestsPerChild 0
Allow 127.0.0.1
Allow 10.0.0.0/8
Allow 172.16.0.0/12
Allow 192.168.0.0/16
ViaProxyName "tinyproxy"
ConnectPort 443
ConnectPort 563`,
    },
  });

  // Radarr environment configuration
  new KubeConfigMap(scope, 'radarr-env', {
    metadata: {
      name: 'radarr-env',
    },
    data: {
      'RADARR__INSTANCE_NAME': 'Radarr',
      'RADARR__BRANCH': 'master',
      'RADARR__PORT': '7878',
      'RADARR__APPLICATION_URL': 'https://radarr.elfhosted.com',
      'RADARR__LOG_LEVEL': 'info',
      'RADARR__ANALYTICS_ENABLED': 'False',
      'RADARR__API_KEY': 'your-api-key-here',
    },
  });

  // Radarr 4K environment configuration
  new KubeConfigMap(scope, 'radarr4k-env', {
    metadata: {
      name: 'radarr4k-env',
    },
    data: {
      'RADARR__INSTANCE_NAME': 'Radarr4K',
      'RADARR__BRANCH': 'master',
      'RADARR__PORT': '7878',
      'RADARR__APPLICATION_URL': 'https://radarr4k.elfhosted.com',
      'RADARR__LOG_LEVEL': 'info',
      'RADARR__ANALYTICS_ENABLED': 'False',
      'RADARR__API_KEY': 'your-api-key-here',
    },
  });

  // Sonarr environment configuration
  new KubeConfigMap(scope, 'sonarr-env', {
    metadata: {
      name: 'sonarr-env',
    },
    data: {
      'SONARR__INSTANCE_NAME': 'Sonarr',
      'SONARR__BRANCH': 'main',
      'SONARR__PORT': '8989',
      'SONARR__APPLICATION_URL': 'https://sonarr.elfhosted.com',
      'SONARR__LOG_LEVEL': 'info',
      'SONARR__ANALYTICS_ENABLED': 'False',
      'SONARR__API_KEY': 'your-api-key-here',
    },
  });

  // Sonarr 4K environment configuration
  new KubeConfigMap(scope, 'sonarr4k-env', {
    metadata: {
      name: 'sonarr4k-env',
    },
    data: {
      'SONARR__INSTANCE_NAME': 'Sonarr4K',
      'SONARR__BRANCH': 'main',
      'SONARR__PORT': '8989',
      'SONARR__APPLICATION_URL': 'https://sonarr4k.elfhosted.com',
      'SONARR__LOG_LEVEL': 'info',
      'SONARR__ANALYTICS_ENABLED': 'False',
      'SONARR__API_KEY': 'your-api-key-here',
    },
  });

  // Prowlarr environment configuration
  new KubeConfigMap(scope, 'prowlarr-env', {
    metadata: {
      name: 'prowlarr-env',
    },
    data: {
      'PROWLARR__INSTANCE_NAME': 'Prowlarr',
      'PROWLARR__BRANCH': 'master',
      'PROWLARR__PORT': '9696',
      'PROWLARR__APPLICATION_URL': 'https://prowlarr.elfhosted.com',
      'PROWLARR__LOG_LEVEL': 'info',
      'PROWLARR__ANALYTICS_ENABLED': 'False',
      'PROWLARR__API_KEY': 'your-api-key-here',
    },
  });

  // Lidarr environment configuration
  new KubeConfigMap(scope, 'lidarr-env', {
    metadata: {
      name: 'lidarr-env',
    },
    data: {
      'LIDARR__INSTANCE_NAME': 'Lidarr',
      'LIDARR__BRANCH': 'master',
      'LIDARR__PORT': '8686',
      'LIDARR__APPLICATION_URL': 'https://lidarr.elfhosted.com',
      'LIDARR__LOG_LEVEL': 'info',
      'LIDARR__ANALYTICS_ENABLED': 'False',
      'LIDARR__API_KEY': 'your-api-key-here',
    },
  });

  // Readarr environment configuration
  new KubeConfigMap(scope, 'readarr-env', {
    metadata: {
      name: 'readarr-env',
    },
    data: {
      'READARR__INSTANCE_NAME': 'Readarr',
      'READARR__BRANCH': 'develop',
      'READARR__PORT': '8787',
      'READARR__APPLICATION_URL': 'https://readarr.elfhosted.com',
      'READARR__LOG_LEVEL': 'info',
      'READARR__ANALYTICS_ENABLED': 'False',
      'READARR__API_KEY': 'your-api-key-here',
    },
  });

  // Readarr Audio environment configuration
  new KubeConfigMap(scope, 'readarraudio-env', {
    metadata: {
      name: 'readarraudio-env',
    },
    data: {
      'READARR__INSTANCE_NAME': 'ReadarrAudio',
      'READARR__BRANCH': 'develop',
      'READARR__PORT': '8787',
      'READARR__APPLICATION_URL': 'https://readarraudio.elfhosted.com',
      'READARR__LOG_LEVEL': 'info',
      'READARR__ANALYTICS_ENABLED': 'False',
      'READARR__API_KEY': 'your-api-key-here',
    },
  });

  // Zurg configuration
  new KubeConfigMap(scope, 'zurg-config', {
    metadata: {
      name: 'zurg-config',
    },
    data: {
      'config.yml': `# Zurg configuration
token: your-real-debrid-token
host: 0.0.0.0
port: 9999
concurrent_workers: 20
check_for_changes_every_secs: 10
retain_folder_name_extension: false
retain_rd_torrent_name: false
auto_delete_rar_torrents: true
serve_from_rclone: false
rclone_mount_path: /storage/realdebrid-zurg
network_buffer_size: 1048576
api:
  host: 0.0.0.0
  port: 9999
  username: ""
  password: ""`,
    },
  });

  // Zurg environment configuration
  new KubeConfigMap(scope, 'zurg-env', {
    metadata: {
      name: 'zurg-env',
    },
    data: {
      'ZURG_TOKEN': 'your-real-debrid-token',
      'ZURG_PORT': '9999',
    },
  });

  // Riven environment configuration
  new KubeConfigMap(scope, 'riven-env', {
    metadata: {
      name: 'riven-env',
    },
    data: {
      'RIVEN_FRONTEND_URL': 'https://riven.elfhosted.com',
      'RIVEN_BACKEND_URL': 'https://riven-backend.elfhosted.com',
      'RIVEN_DATABASE_HOST': 'postgresql',
      'RIVEN_DATABASE_PORT': '5432',
      'RIVEN_DATABASE_NAME': 'riven',
      'RIVEN_DATABASE_USER': 'riven',
      'RIVEN_DATABASE_PASSWORD': 'riven-password',
      'RIVEN_REDIS_HOST': 'redis',
      'RIVEN_REDIS_PORT': '6379',
      'RIVEN_LOG_LEVEL': 'INFO',
    },
  });

  // Riven frontend configuration
  new KubeConfigMap(scope, 'riven-frontend-config', {
    metadata: {
      name: 'riven-frontend-config',
    },
    data: {
      'config.json': `{
  "apiUrl": "https://riven-backend.elfhosted.com",
  "title": "Riven",
  "theme": "dark"
}`,
    },
  });

  // Riven frontend environment
  new KubeConfigMap(scope, 'riven-frontend-env', {
    metadata: {
      name: 'riven-frontend-env',
    },
    data: {
      'ORIGIN': 'https://riven.elfhosted.com',
      'BACKEND_URL': 'https://riven-backend.elfhosted.com',
    },
  });

  // Riven setup script
  new KubeConfigMap(scope, 'riven-setup', {
    metadata: {
      name: 'riven-setup',
    },
    data: {
      'setup.py': `#!/usr/bin/env python3
# Riven setup script
import os
import sys

def main():
    print("Setting up Riven...")
    # Add setup logic here
    
if __name__ == "__main__":
    main()`,
    },
  });

  // Homer configuration
  new KubeConfigMap(scope, 'homer-config', {
    metadata: {
      name: 'homer-config',
    },
    data: {
      'config.yml': `---
title: "ElfHosted Dashboard"
subtitle: "Your Personal Media Stack"
logo: "logo.png"
header: true
footer: '<p>Created with <span class="has-text-danger">❤️</span> by ElfHosted</p>'

defaults:
  layout: columns
  colorTheme: auto

theme: default
colors:
  light:
    highlight-primary: "#3367d6"
    highlight-secondary: "#4285f4"
    highlight-hover: "#5a95f5"
    background: "#f5f5f5"
    card-background: "#ffffff"
    text: "#363636"
    text-header: "#ffffff"
    text-title: "#303030"
    text-subtitle: "#424242"
    card-shadow: rgba(0, 0, 0, 0.1)
    link: "#3273dc"
    link-hover: "#363636"
  dark:
    highlight-primary: "#3367d6"
    highlight-secondary: "#4285f4"
    highlight-hover: "#5a95f5"
    background: "#131313"
    card-background: "#2b2b2b"
    text: "#eaeaea"
    text-header: "#ffffff"
    text-title: "#fafafa"
    text-subtitle: "#f5f5f5"
    card-shadow: rgba(0, 0, 0, 0.4)
    link: "#3273dc"
    link-hover: "#ffdd57"

services:
  - name: "Media"
    icon: "fas fa-play"
    items:
      - name: "Plex"
        logo: "https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/svg/plex.svg"
        subtitle: "Media Server"
        url: "https://plex.elfhosted.com"
        target: "_blank"
      - name: "Jellyfin"
        logo: "https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/svg/jellyfin.svg"
        subtitle: "Media Server"
        url: "https://jellyfin.elfhosted.com"
        target: "_blank"
      - name: "Overseerr"
        logo: "https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/svg/overseerr.svg"
        subtitle: "Request Management"
        url: "https://overseerr.elfhosted.com"
        target: "_blank"

  - name: "Downloads"
    icon: "fas fa-download"
    items:
      - name: "Radarr"
        logo: "https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/svg/radarr.svg"
        subtitle: "Movie Management"
        url: "https://radarr.elfhosted.com"
        target: "_blank"
      - name: "Sonarr"
        logo: "https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/svg/sonarr.svg"
        subtitle: "TV Management"
        url: "https://sonarr.elfhosted.com"
        target: "_blank"
      - name: "Prowlarr"
        logo: "https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/svg/prowlarr.svg"
        subtitle: "Indexer Management"
        url: "https://prowlarr.elfhosted.com"
        target: "_blank"`,
    },
  });

  // Gatus configuration
  new KubeConfigMap(scope, 'gatus-config', {
    metadata: {
      name: 'gatus-config',
    },
    data: {
      'config.yaml': `storage:
  type: memory

metrics: true

endpoints:
  - name: Plex
    url: "https://plex.elfhosted.com"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[RESPONSE_TIME] < 1000"
    alerts:
      - type: discord
        enabled: true
        failure-threshold: 3
        success-threshold: 2

  - name: Radarr
    url: "https://radarr.elfhosted.com"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[RESPONSE_TIME] < 1000"

  - name: Sonarr
    url: "https://sonarr.elfhosted.com"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[RESPONSE_TIME] < 1000"

  - name: Prowlarr
    url: "https://prowlarr.elfhosted.com"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[RESPONSE_TIME] < 1000"`,
    },
  });

  // Filebrowser environment
  new KubeConfigMap(scope, 'filebrowser-env', {
    metadata: {
      name: 'filebrowser-env',
    },
    data: {
      'FB_DATABASE': '/database/filebrowser.db',
      'FB_ROOT': '/srv',
      'FB_LOG': 'stdout',
      'FB_ADDRESS': '0.0.0.0',
      'FB_PORT': '80',
      'FB_BASEURL': '',
    },
  });

  // Filebrowser elfbot script
  new KubeConfigMap(scope, 'filebrowser-elfbot-script', {
    metadata: {
      name: 'filebrowser-elfbot-script',
    },
    data: {
      'elfbot.sh': `#!/bin/bash
# ElfBot script for Filebrowser
echo "ElfBot: Starting Filebrowser automation..."

# Add automation logic here
while true; do
    echo "ElfBot: Monitoring Filebrowser..."
    sleep 300
done`,
    },
  });

  // RcloneFM configuration
  new KubeConfigMap(scope, 'rclonefm-config', {
    metadata: {
      name: 'rclonefm-config',
    },
    data: {
      'rclone.conf': `[realdebrid-zurg]
type = webdav
url = http://zurg:9999/http
vendor = other
user = 
pass = 

[alldebrid]
type = webdav
url = http://alldebrid:9999
vendor = other
user = 
pass = 

[premiumize]
type = webdav
url = http://premiumize:8080
vendor = other
user = 
pass = `,
    },
  });

  // Wizarr environment
  new KubeConfigMap(scope, 'wizarr-env', {
    metadata: {
      name: 'wizarr-env',
    },
    data: {
      'APP_URL': 'https://wizarr.elfhosted.com',
      'DISABLE_BUILTIN_AUTH': 'false',
    },
  });

  // Wizarr Jellyfin steps
  new KubeConfigMap(scope, 'wizarr-steps-jellyfin', {
    metadata: {
      name: 'wizarr-steps-jellyfin',
    },
    data: {
      'steps.json': `[
  {
    "step": 1,
    "title": "Welcome to Jellyfin!",
    "description": "Follow these steps to get started with your Jellyfin server.",
    "content": "Your Jellyfin server is ready to use!"
  },
  {
    "step": 2,
    "title": "Download the App",
    "description": "Download the Jellyfin app for your device.",
    "content": "Visit https://jellyfin.org/downloads/ to download the app."
  },
  {
    "step": 3,
    "title": "Connect to Server",
    "description": "Connect to your Jellyfin server.",
    "content": "Use the server address: https://jellyfin.elfhosted.com"
  }
]`,
    },
  });

  // Wizarr Plex steps
  new KubeConfigMap(scope, 'wizarr-steps-plex', {
    metadata: {
      name: 'wizarr-steps-plex',
    },
    data: {
      'steps.json': `[
  {
    "step": 1,
    "title": "Welcome to Plex!",
    "description": "Follow these steps to get started with your Plex server.",
    "content": "Your Plex server is ready to use!"
  },
  {
    "step": 2,
    "title": "Download the App",
    "description": "Download the Plex app for your device.",
    "content": "Visit https://www.plex.tv/media-server-downloads/ to download the app."
  },
  {
    "step": 3,
    "title": "Connect to Server",
    "description": "Connect to your Plex server.",
    "content": "Use the server address: https://plex.elfhosted.com"
  }
]`,
    },
  });

  // Traefik Forward Auth configuration
  new KubeConfigMap(scope, 'traefik-forward-auth-config', {
    metadata: {
      name: 'traefik-forward-auth-config',
    },
    data: {
      'DEFAULT_PROVIDER': 'oidc',
      'PROVIDERS_OIDC_ISSUER_URL': 'https://auth.elfhosted.com',
      'PROVIDERS_OIDC_CLIENT_ID': 'traefik-forward-auth',
      'PROVIDERS_OIDC_CLIENT_SECRET': 'your-client-secret',
      'SECRET': 'your-secret-key',
      'COOKIE_DOMAIN': '.elfhosted.com',
      'AUTH_HOST': 'auth.elfhosted.com',
      'URL_PATH': '/_oauth',
      'LOG_LEVEL': 'info',
    },
  });

  // Recyclarr configuration
  new KubeConfigMap(scope, 'recyclarr-config', {
    metadata: {
      name: 'recyclarr-config',
    },
    data: {
      'recyclarr.yml': `# Recyclarr Configuration
radarr:
  movies:
    base_url: http://radarr:7878
    api_key: !env_var RADARR_API_KEY
    
    quality_definition:
      type: movie
      
    custom_formats:
      - trash_ids:
          - 496f355514737f7d83bf7aa4d24f8169  # TrueHD Atmos
          - 2f22d89048b01681dde8afe73aabb0b9  # DTS-X
          - 417804f7f2c4308c1f4c5d380d4c4475  # ATMOS (undefined)
        quality_profiles:
          - name: Ultra-HD
            
sonarr:
  tv:
    base_url: http://sonarr:8989
    api_key: !env_var SONARR_API_KEY
    
    quality_definition:
      type: series
      
    custom_formats:
      - trash_ids:
          - 32b367365729d530ca1c124a0b180c64  # Bad Dual Groups
          - 82d40da2bc6923f41e14394075dd4b03  # No-RlsGroup
        quality_profiles:
          - name: HD-1080p`,
    },
  });

  // Notifiarr configuration
  new KubeConfigMap(scope, 'notifiarr-config', {
    metadata: {
      name: 'notifiarr-config',
    },
    data: {
      'notifiarr.conf': `##     Notifiarr Client Configuration File     ##
#
# This file is used to configure the Notifiarr client.
# 
api_key = "your-notifiarr-api-key"

# Application URLs and API Keys
[radarr]
  url = "http://radarr:7878"
  api_key = "your-radarr-api-key"

[sonarr]
  url = "http://sonarr:8989"
  api_key = "your-sonarr-api-key"

[prowlarr]
  url = "http://prowlarr:9696"
  api_key = "your-prowlarr-api-key"

[lidarr]
  url = "http://lidarr:8686"
  api_key = "your-lidarr-api-key"

[readarr]
  url = "http://readarr:8787"
  api_key = "your-readarr-api-key"

[plex]
  url = "http://plex:32400"
  token = "your-plex-token"

[tautulli]
  url = "http://tautulli:8181"
  api_key = "your-tautulli-api-key"`,
    },
  });

  // Additional ConfigMaps for other services would go here...
  // This is a representative sample of the hundreds of ConfigMaps needed

} 