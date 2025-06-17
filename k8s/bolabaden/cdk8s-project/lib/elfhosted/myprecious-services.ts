import { Construct } from 'constructs';
import { KubeService, IntOrString } from '../../imports/k8s';

/**
 * Service definitions for all myprecious services
 * This contains the comprehensive list of services from the myprecious Helm chart
 */

interface ServiceConfig {
  name: string;
  ports: Array<{
    port: number;
    name: string;
    targetPort?: number;
  }>;
}

export const MYPRECIOUS_SERVICES: ServiceConfig[] = [
  // Media Services
  { name: 'plex', ports: [{ port: 32400, name: 'http' }, { port: 3002, name: 'speed' }, { port: 8888, name: 'tinyproxy' }] },
  { name: 'jellyfin', ports: [{ port: 8096, name: 'http' }] },
  { name: 'emby', ports: [{ port: 8096, name: 'http' }, { port: 8920, name: 'https' }] },

  // *arr Services
  { name: 'radarr', ports: [{ port: 7878, name: 'http' }] },
  { name: 'radarr-4k', ports: [{ port: 7878, name: 'http' }] },
  { name: 'sonarr', ports: [{ port: 8989, name: 'http' }] },
  { name: 'sonarr-4k', ports: [{ port: 8989, name: 'http' }] },
  { name: 'lidarr', ports: [{ port: 8686, name: 'http' }] },
  { name: 'readarr', ports: [{ port: 8787, name: 'http' }] },
  { name: 'readarraudio', ports: [{ port: 8787, name: 'http' }] },
  { name: 'prowlarr', ports: [{ port: 9696, name: 'http' }] },
  { name: 'bazarr', ports: [{ port: 6767, name: 'http' }] },
  { name: 'bazarr4k', ports: [{ port: 6767, name: 'http' }] },

  // Download Clients
  { name: 'qbittorrent', ports: [{ port: 8080, name: 'http' }, { port: 8999, name: 'gluetun' }] },
  { name: 'deluge', ports: [{ port: 8112, name: 'http' }, { port: 58846, name: 'daemon' }, { port: 8999, name: 'gluetun' }] },
  { name: 'rutorrent', ports: [{ port: 80, name: 'http' }, { port: 5000, name: 'scgi' }, { port: 6881, name: 'bt-tcp' }, { port: 6882, name: 'bt-udp' }, { port: 8999, name: 'gluetun' }] },
  { name: 'sabnzbd', ports: [{ port: 8080, name: 'http' }] },
  { name: 'nzbget', ports: [{ port: 6789, name: 'http' }] },

  // Request Management
  { name: 'overseerr', ports: [{ port: 5055, name: 'http' }] },
  { name: 'jellyseerr', ports: [{ port: 5055, name: 'http' }] },
  { name: 'ombi', ports: [{ port: 3579, name: 'http' }] },
  { name: 'requestrr', ports: [{ port: 4545, name: 'http' }] },

  // Monitoring & Management
  { name: 'tautulli', ports: [{ port: 8181, name: 'http' }] },
  { name: 'gatus', ports: [{ port: 8080, name: 'http' }] },
  { name: 'homer', ports: [{ port: 8080, name: 'http' }] },
  { name: 'homepage', ports: [{ port: 3000, name: 'http' }] },
  { name: 'uptime-kuma', ports: [{ port: 3001, name: 'http' }] },
  { name: 'gotify', ports: [{ port: 80, name: 'http' }] },

  // File Management
  { name: 'filebrowser', ports: [{ port: 80, name: 'http' }] },
  { name: 'rclonefm', ports: [{ port: 5572, name: 'http' }] },
  { name: 'rcloneui', ports: [{ port: 5572, name: 'http' }] },

  // Streaming & Addons
  { name: 'zurg', ports: [{ port: 9999, name: 'http' }] },
  { name: 'riven', ports: [{ port: 8080, name: 'http' }, { port: 8001, name: 'backend' }] },
  { name: 'riven-frontend', ports: [{ port: 3000, name: 'http' }] },
  { name: 'stremio-jackett', ports: [{ port: 7000, name: 'http' }] },
  { name: 'stremthru', ports: [{ port: 8080, name: 'http' }] },
  { name: 'stremify', ports: [{ port: 8080, name: 'http' }] },
  { name: 'mediafusion', ports: [{ port: 8000, name: 'http' }] },
  { name: 'comet', ports: [{ port: 8000, name: 'http' }] },
  { name: 'aiostreams', ports: [{ port: 8080, name: 'http' }] },
  { name: 'knightcrawler', ports: [{ port: 8080, name: 'http' }] },
  { name: 'jackettio', ports: [{ port: 3000, name: 'http' }] },
  { name: 'torrentio', ports: [{ port: 7000, name: 'http' }] },

  // Debrid Services
  { name: 'rdtclient', ports: [{ port: 6500, name: 'http' }] },
  { name: 'rdtclient-alldebrid', ports: [{ port: 6500, name: 'http' }] },
  { name: 'rdtclient-premiumize', ports: [{ port: 6500, name: 'http' }] },
  { name: 'rdtclient-torbox', ports: [{ port: 6500, name: 'http' }] },
  { name: 'plex-debrid', ports: [{ port: 3000, name: 'http' }] },
  { name: 'plex-debrid-alldebrid', ports: [{ port: 3000, name: 'http' }] },
  { name: 'plex-debrid-premiumize', ports: [{ port: 3000, name: 'http' }] },
  { name: 'plex-debrid-torbox', ports: [{ port: 3000, name: 'http' }] },
  { name: 'rdebrid-ui', ports: [{ port: 3000, name: 'http' }] },

  // Proxy Services
  { name: 'plex-proxy', ports: [{ port: 32400, name: 'http' }, { port: 8888, name: 'tinyproxy' }] },
  { name: 'premiumize', ports: [{ port: 8080, name: 'http' }] },
  { name: 'alldebrid', ports: [{ port: 9999, name: 'http' }] },
  { name: 'torbox', ports: [{ port: 9999, name: 'http' }] },
  { name: 'debridlink', ports: [{ port: 8080, name: 'http' }] },

  // Utility Services
  { name: 'wizarr', ports: [{ port: 5690, name: 'http' }] },
  { name: 'recyclarr', ports: [{ port: 80, name: 'http' }] },
  { name: 'notifiarr', ports: [{ port: 5454, name: 'http' }] },
  { name: 'flaresolverr', ports: [{ port: 8191, name: 'http' }] },
  //{ name: 'prowlarr', ports: [{ port: 9696, name: 'http' }] },

  // Books & Reading
  { name: 'calibre-web', ports: [{ port: 8083, name: 'http' }] },
  { name: 'calibre', ports: [{ port: 8080, name: 'http' }, { port: 8081, name: 'webserver' }] },
  { name: 'komga', ports: [{ port: 8080, name: 'http' }] },
  { name: 'kavita', ports: [{ port: 5000, name: 'http' }] },
  { name: 'mylar', ports: [{ port: 8090, name: 'http' }] },
  { name: 'lazylibrarian', ports: [{ port: 5299, name: 'http' }] },
  { name: 'openbooks', ports: [{ port: 80, name: 'http' }] },

  // Music
  { name: 'navidrome', ports: [{ port: 4533, name: 'http' }] },

  // TV & Live Streaming
  { name: 'threadfin', ports: [{ port: 34400, name: 'http' }] },
  { name: 'tunarr', ports: [{ port: 8000, name: 'http' }] },
  { name: 'ersatztv', ports: [{ port: 8409, name: 'http' }] },
  { name: 'youriptv', ports: [{ port: 8080, name: 'http' }] },

  // Communication
  { name: 'thelounge', ports: [{ port: 9000, name: 'http' }] },
  { name: 'gotosocial', ports: [{ port: 8080, name: 'http' }] },

  // Storage & Sync
  { name: 'seafile', ports: [{ port: 8000, name: 'http' }] },
  { name: 'seafile-memcached', ports: [{ port: 11211, name: 'memcached' }] },
  { name: 'resilio-sync', ports: [{ port: 8888, name: 'http' }] },
  { name: 'syncthing', ports: [{ port: 8384, name: 'http' }] },

  // Development & Tools
  { name: 'vaultwarden', ports: [{ port: 80, name: 'http' }] },
  { name: 'privatebin', ports: [{ port: 8080, name: 'http' }] },
  { name: 'pairdrop', ports: [{ port: 3000, name: 'http' }] },
  { name: 'wallabag', ports: [{ port: 80, name: 'http' }] },

  // Database & Admin
  { name: 'pgadmin', ports: [{ port: 80, name: 'http' }] },
  { name: 'mongoexpress', ports: [{ port: 8081, name: 'http' }] },
  { name: 'redisinsight', ports: [{ port: 8001, name: 'http' }] },

  // Authentication & Proxy
  { name: 'traefik-forward-auth', ports: [{ port: 4181, name: 'http' }] },
  { name: 'elfterm', ports: [{ port: 7681, name: 'http' }] },
  { name: 'storagehub', ports: [{ port: 8080, name: 'http' }, { port: 9090, name: 'metrics' }] },

  // WebDAV
  { name: 'webdav', ports: [{ port: 8080, name: 'http' }] },
  { name: 'webdav-plus', ports: [{ port: 8080, name: 'http' }] },

  // Additional Services
  { name: 'symlink-cleaner', ports: [{ port: 8080, name: 'http' }] },
  { name: 'scannarr', ports: [{ port: 8080, name: 'http' }] },
  { name: 'scannarr4k', ports: [{ port: 8080, name: 'http' }] },
  { name: 'storyteller', ports: [{ port: 8080, name: 'http' }] },
  { name: 'webstreamr', ports: [{ port: 8080, name: 'http' }] },
  { name: 'nuviostreams', ports: [{ port: 8080, name: 'http' }] },
  { name: 'suggestarr', ports: [{ port: 8080, name: 'http' }] },
  { name: 'profilarr', ports: [{ port: 8080, name: 'http' }] },
  { name: 'pulsarr', ports: [{ port: 8080, name: 'http' }] },
  { name: 'dispatcharr', ports: [{ port: 8080, name: 'http' }] },
  { name: 'decypharr', ports: [{ port: 8282, name: 'http' }] },
];

/**
 * Creates all services for the myprecious chart
 */
export function createAllMyPreciousServices(scope: Construct) {
  MYPRECIOUS_SERVICES.forEach(serviceConfig => {
    new KubeService(scope, `service-${serviceConfig.name}`, {
      metadata: {
        name: serviceConfig.name,
        labels: {
          'app.kubernetes.io/name': serviceConfig.name,
          'app.kubernetes.io/instance': 'myprecious',
        },
      },
      spec: {
        type: 'ClusterIP',
        ports: serviceConfig.ports.map(port => ({
          port: port.port,
          targetPort: IntOrString.fromNumber(port.targetPort || port.port),
          protocol: 'TCP',
          name: port.name,
        })),
        selector: {
          'app.elfhosted.com/name': serviceConfig.name,
        },
      },
    });
  });
} 