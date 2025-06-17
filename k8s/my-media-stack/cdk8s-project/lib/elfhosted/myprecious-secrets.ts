import { Construct } from 'constructs';
import { KubeSecret, KubeEndpoints, KubeIngress } from '../../imports/k8s';

/**
 * Creates all Secrets and additional resources for the myprecious chart
 * This includes API keys, database credentials, and other sensitive data
 */

interface SecretConfig {
  name: string;
  type?: string;
  data?: Record<string, string>;
  stringData?: Record<string, string>;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
}

// Helper function to base64 encode strings
function b64Encode(str: string): string {
  return Buffer.from(str, 'utf8').toString('base64');
}

const SECRET_CONFIGS: SecretConfig[] = [
  // API Keys and tokens
  {
    name: 'api-keys',
    type: 'Opaque',
    data: {
      'anthropic-api-key': b64Encode('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'),
      'gemini-api-key': b64Encode('AIzaSyB-1aUePBQbb7_rpIDz5kY7AMWxguMq7Y8'),
      'groq-api-key': b64Encode('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'),
      'mistral-api-key': b64Encode('rNZOe7bauJwpknHvDjOtnXV6iaKHzEt1'),
      'openai-api-key': b64Encode('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'),
      'openrouter-api-key': b64Encode('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'),
      'perplexity-api-key': b64Encode('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'),
      'firecrawl-api-key': b64Encode('fc-e0e316bead3d497fb24d42ea21d5daf6'),
      'meili-master-key': b64Encode('w9y6ZXi812GH1HEGMLMcrcS6AxuHiQ2jktizLd4ezcgoMW3S2Vyegm4l2DpRB'),
      'autobrr-api-key': b64Encode('74d290518e91075a37dca845c588b56a'),
      'bazarr-api-key': b64Encode('f1ededf2d38845756de95c135c917369'),
      'crowdsec-bouncer-api-key': b64Encode('usvi2urbIG5obProVWwcf2ImTrtB93x668jWDvrH'),
      'fanart-api-key': b64Encode('7a38941365b707a56080cf130f3fcec7'),
      'github-token': b64Encode('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'),
      'gitlab-token': b64Encode('glpat-grYxUyvaSakXzznsnzN2'),
      'jackett-api-key': b64Encode('nnx0n84pcj7umynyd2pbid2nl1zzuemz'),
      'jellyfin-api-key': b64Encode('48963bb342c14be4abccfdced3f8c927'),
      'jellyseerr-api-key': b64Encode('MTc0MDcxNDk2MTY5MDRmMDcxYmM1LTI5ODktNGJiMy05NzExLTdhYjYyMWQxMDAyOQ=='),
      'kavita-api-key': b64Encode('3e263495-10c1-4963-b003-9beb4a3dc553'),
      'kavita-token': b64Encode('7bAgcPhmVK7Es39XGd6g1pLBJbUPNVe4Fu5+sqnuMxX+UXdAgXKuWYEAegWPJbOdNpfamem0wh70tVHyk9XfW35BU3GGquAdfgQZkokLKDD1sXX6sFnd9S4Zx0ITmVpoH4Gx7yPY9kV0yLgU7gbIr8BAhlGd2sNvpsWJppjqnk1vKXB+GEBhgsaH95pgnoE5LShyuD2kVofXvy8NdDKouTlRp9uOdTp0CIyI5263Lx6wtW5ts5m3J0ZTynacw0Dlq80IJe0qjooH8JxhFNbVJDjGRiFcykyUyhrUaQgIL4xn4DjjDuDqg9KWnCMp0tJs8oryqqK3Ix8+xb1dBmQ=='),
      'lidarr-api-key': b64Encode('5cbcba61d75444d6a97c07df768a26e4'),
      'listrr-api-key': b64Encode('d962d6e9fd4e4db790711128ed6013b236e86b1bea534c7f9f0c17a3c4c72982'),
      'mdblist-api-key': b64Encode('so5wnihfjtqhacaa310gzmrq6'),
      'notifiarr-api-key': b64Encode('your-notifiarr-api-key-here'),
      'opensubtitles-api-key': b64Encode('XzkTEC85VNh8YYNMFezChYr52hnrZNW5'),
      'overseerr-api-key': b64Encode('MTc0MDcxNDk2MTY5MDRmMDcxYmM1LTI5ODktNGJiMy05NzExLTdhYjYyMWQxMDAyOQ=='),
      'plex-token': b64Encode('18rS6gM-Wec22uq5Gw-U'),
      'prowlarr-api-key': b64Encode('29440a82740d475cacb35327c62c87a1'),
      'radarr-api-key': b64Encode('fc44a10e3c7c4c0582cf99f1d7ff0f18fc44a10e3c7c4c0582cf99f1d7ff0f18'),
      'readarr-api-key': b64Encode('e5ae49d7512a41548588c9b5456d7010'),
      'replicate-api-key': b64Encode('r8_8ktlYwwthNrp8pPSLZF1DD7507tku082y7ptd'),
      'resend-api-key': b64Encode('re_Ns9xDLUJ_JxnRW8m9hpUgyjecDgNVmSzS'),
      'riven-api-key': b64Encode('3e0be750b16ff363727d05d18ba364fd'),
      'sabnzbd-api-key': b64Encode('f639219a19284f64926aba59d2fbd14d'),
      'sonarr-api-key': b64Encode('5fa91e23f1ce4bfda1159fb0d717b70e'),
      'speedtest-tracker-api-token': b64Encode('r2SHU2FL1jIQ5aQ5xVsk6T19HjRAGQQsUl0dB7ji140db905'),
      'tautulli-api-key': b64Encode('f86827c56ae740e6bff77c386e726ab8'),
      'tavily-api-key': b64Encode('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'),
      'tmdb-access-token': b64Encode('cec876f852b9c15d2c1b436b1117dff7'),
      'tmdb-api-key': b64Encode('cec876f852b9c15d2c1b436b1117dff7'),
      'tmdb-read-api-key': b64Encode('eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJjZWM4NzZmODUyYjljMTVkMmMxYjQzNmIxMTE3ZGZmNyIsIm5iZiI6MTczNjQzOTU1NC43NzUsInN1YiI6IjY3N2ZmNzAyMjE4ZmQ1N2FjZjRlOGQ2ZCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.yK5N4wwBXNKwauenlWmuEv_5jpoTiK8Bs2rit6HO1E'),
      'todoist-api-key': b64Encode('d9f0ebcdff604bd11c0944f52360758f7587f907'),
      'umami-app-secret': b64Encode('Cg2OAmj7fK1jBCmqT4kHx1gpIWp7tluJEWN5CDkjzBPQJ5LdepoYkSKMLNv6jquS'),
      'whisparr-api-key': b64Encode('fb6df473cfe14fda8498b8b881902985'),
    },
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/component': 'api-keys',
    },
  },

  // Database credentials
  {
    name: 'database-credentials',
    type: 'Opaque',
    data: {
      'postgres-username': b64Encode('postgres'),
      'postgres-password': b64Encode('postgres'),
      'postgres-database': b64Encode('myprecious'),
      'redis-password': b64Encode('c3ll0h3r0123'),
      'mongodb-username': b64Encode('mongodb'),
      'mongodb-password': b64Encode('c3ll0h3r0123'),
      'mongodb-database': b64Encode('myprecious'),
    },
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/component': 'database',
    },
  },

  // Debrid service credentials
  {
    name: 'debrid-credentials',
    type: 'Opaque',
    data: {
      'alldebrid-token': b64Encode('your-alldebrid-token-here'),
      'debridlink-token': b64Encode('your-debridlink-token-here'),
      'premiumize-token': b64Encode('6i6zgswm35baj3ur'),
      'realdebrid-token': b64Encode('FA5NDX4HOAQROKATB6WD7JTZLBZPWYLQSQE6UNBZ7XTUPAKG3YPQ'),
      'torbox-token': b64Encode('d490b7b1-1f6d-44e8-9cc8-3e735c87874b'),
    },
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/component': 'debrid',
    },
  },

  // VPN credentials
  {
    name: 'vpn-credentials',
    type: 'Opaque',
    data: {
      'vpn-username': b64Encode('117274388'),
      'vpn-password': b64Encode('6i6zgswm35baj3ur'),
      'wireguard-private-key': b64Encode('your-wireguard-private-key'),
      'openvpn-config': b64Encode(`remote vpn-us.premiumize.me
verify-x509-name CN=vpn-us.premiumize.me
auth-user-pass
client
dev tun
proto udp
cipher AES-256-CBC
resolv-retry infinite
nobind
persist-key
persist-tun
mute-replay-warnings
verb 3
reneg-sec 0
setenv CLIENT_CERT 0
ignore-unknown-option block-outside-dns
<ca>
-----BEGIN CERTIFICATE-----
MIIFJTCCAw2gAwIBAgIRAPAmbQRNE+PBqvFyFG8GOSIwDQYJKoZIhvcNAQELBQAw
LDEYMBYGA1UECgwPU2VjdXJlIFNlcnZpY2VzMRAwDgYDVQQDDAdSb290IFgxMB4X
DTIxMDEwNDE1MjEwM1oXDTQwMTIzMDE1MjEwM1owLDEYMBYGA1UECgwPU2VjdXJl
IFNlcnZpY2VzMRAwDgYDVQQDDAdSb290IFgxMIICIjANBgkqhkiG9w0BAQEFAAOC
Ag8AMIICCgKCAgEAvuhFcbO0Y5SXBr+h/XU1sPXo/OSjN4W32jzVZ3jmkqA2nH5D
dI5XYHYB9JkW23K37zOlvOWj9J3HiV6WYk0uqQ3cpqDMnIpi1MJCtSRxiaD7LTNO
XrvLpsREq9Vf+lN+zxFxhINv3W4jUyfx5zIjPyjY+vFgH5G56b12EeLKFLYUWwgF
9vTicXJvlV1h+TJVF8wZ5DugrgNZhSZ6QZzda0Zdu1dCwZOafd1GGDWPubgZ8enF
b9gWPGZmS59ZvlLsLQmtLvCEegELzRR5G/OI1CIntcImcKO9ZDdWP+HmqJ4Ss7Ng
0g5xQVx9VgMxwpHw4h/tMV8q1vQQZJXJHdoDRQr5RzVnMjq90FFMd3nfX1B6cheH
cBq31J84ifxcO2vIz/eK/6zNyynuj2+99Yz7MsYI+sfGznngglpYo9u3xNkRfRPP
TA/tXAzrGRqlHmJmMStfutIfHxea4jJCGOeb5eUETQ1kXf5sMSdf0+Ya64h9Kcxw
afMzU8qUAMk4ei+jTEAuUQ25Jq0orFc5D2Pv7RbU9CCzIML2KqqWTThKSSVQbW0x
ufQDxcFZKSKL03484zVPBDMQ65enBUC+9Or9YKT0S0F1IdOX4EUGpqzQAocHFF3u
G1zYz5cjYq2xHRBZAYs6SQH6f16bPrk66uA2B/vZJm5L9Kzni5q46IGsCUsCAwEA
AaNCMEAwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYE
FETGtuYkU3IlxCvmtkAzpP5Xuf+DMA0GCSqGSIb3DQEBCwUAA4ICAQBdTlQtxa6a
PXaZpzsi5gLGrDUhnoTF5cVvt9vYa6c7IaAuYIUMGN1WiwHUL0g1udxi7iC3oBzM
Mc6shYMijEMG9SvC0RJOky1BCPbKdlsRDQbYLgjN5X1fEDDu1cR/vInt/K6vTcLf
1ud6VFgDv8uYG1X9GVs6foAqjYTLt1A99cwHtyCZRxrBLyjgE8S2U07FDpakqlMo
j4S1lgvw4EuEN54Yo55wLv9J0pZ1w9LEpZ4VdEVzLmwYMPG25Ow7HFZDK7j0TZkA
ByvxUvNGWDk5MS71BtWzhrmmPh9kN3HobkfOllqUsE03QaLggg+pe+OYNdIlCBRz
9QJeMtGP7afdBZ5UUL8df0fWL0adUzk1N/1kvAmbvccNu7tej/4LLVv9j1aQtI76
iEXh0WIWTAOz9IsQyP21PhxG01Y3ScIrbOsP+IX+RL/VK4yfTZ2jICj7js8vm6ut
qw3FPyZgeZuBe0gk0cF4vPuFacS23PkEFXbVnt4h40xKwR5ezGshS5I36nXowTdh
uJPLnZaLGjsiFpprztwwDZkDLu5uT6YhT/M6Kb0V3s3a2W1N90eNjo5WEe/poCgr
1+NVJNJQ71s3kodfghA/t1iivJEl58scVLxxJXyhCK+k5IPAB0syq2cR4AKnr/fe
9WFpWZc7lKeI4rwzGej+kMCQ6ujinVPyRA==
-----END CERTIFICATE-----
</ca>
`),
    },
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/component': 'vpn',
    },
  },

  // Authentication secrets
  {
    name: 'auth-secrets',
    type: 'Opaque',
    data: {
      'jwt-secret': b64Encode('lhTCXQ3bH7hGoAuRBqP6ODk8ErJjdbbr15VBUC3dYEQ'),
      'session-secret': b64Encode('your-session-secret-key-here'),
      'oauth-client-id': b64Encode('your-oauth-client-id'),
      'oauth-client-secret': b64Encode('your-oauth-client-secret'),
      'oidc-client-secret': b64Encode('your-oidc-client-secret'),
    },
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/component': 'auth',
    },
  },

  // Notification service credentials
  {
    name: 'notification-credentials',
    type: 'Opaque',
    data: {
      'discord-webhook': b64Encode('your-discord-webhook-url'),
      'telegram-bot-token': b64Encode('your-telegram-bot-token'),
      'telegram-chat-id': b64Encode('your-telegram-chat-id'),
      'slack-webhook': b64Encode('your-slack-webhook-url'),
      'pushover-token': b64Encode('your-pushover-token'),
      'pushover-user': b64Encode('your-pushover-user'),
      'email-smtp-host': b64Encode('smtp.gmail.com'),
      'email-smtp-port': b64Encode('587'),
      'email-username': b64Encode('boden.crouch@gmail.com'),
      'email-password': b64Encode('your-email-password'),
    },
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/component': 'notifications',
    },
  },

  // Backup service credentials
  {
    name: 'backup-credentials',
    type: 'Opaque',
    data: {
      's3-access-key': b64Encode('your-s3-access-key'),
      's3-secret-key': b64Encode('your-s3-secret-key'),
      's3-bucket': b64Encode('your-backup-bucket'),
      's3-region': b64Encode('us-central-1'),
      'rclone-config': b64Encode('your-rclone-config-here'),
      'restic-password': b64Encode('your-restic-password'),
    },
    labels: {
      'app.kubernetes.io/instance': 'myprecious',
      'app.kubernetes.io/component': 'backup',
    },
  },
];

export function createAllMyPreciousSecrets(scope: Construct) {
  SECRET_CONFIGS.forEach(config => {
    new KubeSecret(scope, `secret-${config.name}`, {
      metadata: {
        name: config.name,
        ...(config.labels && { labels: config.labels }),
        ...(config.annotations && { annotations: config.annotations }),
      },
      type: config.type || 'Opaque',
      ...(config.data && { data: config.data }),
    });
  });
}

/**
 * Creates external service endpoints for services not running in the cluster
 */
export function createExternalEndpoints(scope: Construct) {
  // External database endpoint (if using external database)
  new KubeEndpoints(scope, 'external-postgres', {
    metadata: {
      name: 'external-postgres',
      labels: {
        'app.kubernetes.io/instance': 'myprecious',
        'app.kubernetes.io/component': 'external-database',
      },
    },
    subsets: [{
      addresses: [{
        ip: '10.0.0.100', // Replace with actual external database IP
      }],
      ports: [{
        name: 'postgres',
        port: 5432,
        protocol: 'TCP',
      }],
    }],
  });

  // External Redis endpoint (if using external Redis)
  new KubeEndpoints(scope, 'external-redis', {
    metadata: {
      name: 'external-redis',
      labels: {
        'app.kubernetes.io/instance': 'myprecious',
        'app.kubernetes.io/component': 'external-cache',
      },
    },
    subsets: [{
      addresses: [{
        ip: '10.0.0.101', // Replace with actual external Redis IP
      }],
      ports: [{
        name: 'redis',
        port: 6379,
        protocol: 'TCP',
      }],
    }],
  });
}

/**
 * Creates ingress resources for external access
 */
export function createIngressResources(scope: Construct) {
  const ingressConfigs = [
    { name: 'filebrowser', host: 'filebrowser.elfhosted.com', service: 'filebrowser', port: 80 },
    { name: 'gatus', host: 'gatus.elfhosted.com', service: 'gatus', port: 8080 },
    { name: 'homer', host: 'homer.elfhosted.com', service: 'homer', port: 8080 },
    { name: 'jellyfin', host: 'jellyfin.elfhosted.com', service: 'jellyfin', port: 8096 },
    { name: 'jellyseerr', host: 'jellyseerr.elfhosted.com', service: 'jellyseerr', port: 5055 },
    { name: 'lidarr', host: 'lidarr.elfhosted.com', service: 'lidarr', port: 8686 },
    { name: 'overseerr', host: 'overseerr.elfhosted.com', service: 'overseerr', port: 5055 },
    { name: 'plex', host: 'plex.elfhosted.com', service: 'plex', port: 32400 },
    { name: 'prowlarr', host: 'prowlarr.elfhosted.com', service: 'prowlarr', port: 9696 },
    { name: 'qbittorrent', host: 'qbittorrent.elfhosted.com', service: 'qbittorrent', port: 8080 },
    { name: 'radarr', host: 'radarr.elfhosted.com', service: 'radarr', port: 7878 },
    { name: 'readarr', host: 'readarr.elfhosted.com', service: 'readarr', port: 8787 },
    { name: 'riven-backend', host: 'riven-backend.elfhosted.com', service: 'riven', port: 8080 },
    { name: 'riven', host: 'riven.elfhosted.com', service: 'riven-frontend', port: 3000 },
    { name: 'sonarr', host: 'sonarr.elfhosted.com', service: 'sonarr', port: 8989 },
    { name: 'tautulli', host: 'tautulli.elfhosted.com', service: 'tautulli', port: 8181 },
    { name: 'wizarr', host: 'wizarr.elfhosted.com', service: 'wizarr', port: 5690 },
    { name: 'zurg', host: 'zurg.elfhosted.com', service: 'zurg', port: 9999 },
  ];

  ingressConfigs.forEach(config => {
    new KubeIngress(scope, `ingress-${config.name}`, {
      metadata: {
        name: `${config.name}-ingress`,
        labels: {
          'app.kubernetes.io/instance': 'myprecious',
          'app.kubernetes.io/component': 'ingress',
        },
        annotations: {
          'kubernetes.io/ingress.class': 'traefik',
          'traefik.ingress.kubernetes.io/router.entrypoints': 'websecure',
          'traefik.ingress.kubernetes.io/router.tls': 'true',
          'cert-manager.io/cluster-issuer': 'letsencrypt-prod',
        },
      },
      spec: {
        tls: [{
          hosts: [config.host],
          secretName: `${config.name}-tls`,
        }],
        rules: [{
          host: config.host,
          http: {
            paths: [{
              path: '/',
              pathType: 'Prefix',
              backend: {
                service: {
                  name: config.service,
                  port: {
                    number: config.port,
                  },
                },
              },
            }],
          },
        }],
      },
    });
  });
}

/**
 * Creates TLS certificates for services
 */
export function createTLSSecrets(scope: Construct) {
  const domains = [
    'plex.elfhosted.com',
    'radarr.elfhosted.com',
    'sonarr.elfhosted.com',
    'prowlarr.elfhosted.com',
    'lidarr.elfhosted.com',
    'readarr.elfhosted.com',
    'jellyfin.elfhosted.com',
    'overseerr.elfhosted.com',
    'jellyseerr.elfhosted.com',
    'tautulli.elfhosted.com',
    'qbittorrent.elfhosted.com',
    'filebrowser.elfhosted.com',
    'homer.elfhosted.com',
    'gatus.elfhosted.com',
    'wizarr.elfhosted.com',
    'zurg.elfhosted.com',
    'riven.elfhosted.com',
    'riven-backend.elfhosted.com',
  ];

  domains.forEach(domain => {
    const serviceName = domain.split('.')[0];
    new KubeSecret(scope, `tls-${serviceName}`, {
      metadata: {
        name: `${serviceName}-tls`,
        labels: {
          'app.kubernetes.io/instance': 'myprecious',
          'app.kubernetes.io/component': 'tls',
        },
      },
      type: 'kubernetes.io/tls',
      data: {
        'tls.crt': b64Encode(`-----BEGIN CERTIFICATE-----
MIIDGzCCAgOgAwIBAgIUecJULBeRe6 +/ralwgf8GFYSi6FAwDQYJKoZIhvcNAQEL
BQAwHTEbMBkGA1UEAwwScGxleC5lbGZob3N0ZWQuY29tMB4XDTI1MDYxMjAxMTAw
MFoXDTI2MDYxMjAxMTAwMFowHTEbMBkGA1UEAwwScGxleC5lbGZob3N0ZWQuY29t
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4Si7k5BOn7RerIVjI2Ap
AGCoTvcMqLrwfmmvsJMg + CyVYQ2t / ic8Ze0N5MVRqtycQvkL7 / jM7OXPQlAy49SH
JDZmclrNC5gpeyYli7RWQnlI4rDyj / ldq5HX5SED6KocX1E6h3ziISJYzauZWQcU
LdG2qvYtgVUQYYcYvEfRA6r8rJk78PhCaNlWQ844G68 / Exaxf1Ez / Tdb + 7WCH2 + X
vAJk4Lehk + s5CcptjNBeLi6CSakCmfnAJToCfXAaHVLtkFGUy2TnvTEcBzncwPAb
KnQUvV3cflaGQf3n4UZmOUx5qB5dwBM6Z1SDFkZROO8IqmnjxNqins3eOSWtIhXN
WQIDAQABo1MwUTAdBgNVHQ4EFgQUBHV2hFx1vyu6VfXR1TASYexLz4EwHwYDVR0j
BBgwFoAUBHV2hFx1vyu6VfXR1TASYexLz4EwDwYDVR0TAQH / BAUwAwEB / zANBgkq
hkiG9w0BAQsFAAOCAQEAohDWKRXpA7JGoyF / cp / nRNygXi7ShoymlBIi7UjOilFz
DVYPXUB + LcnhVWcMQi + VtGDJZmPzd9l0YQ4Ax20B4pRktuULUD5d4EGJTVix6AuL
NxovjkaO9ciBhuEwOc1VrIPySR0ocmYLoavNr / 6 + 6XCSmUPDT5GBUBnCcjo6eK4m
0IvhbrJJ9590uBLD0dWOobq9Vxf4XvTFRS5Ms32zpoMYTNcfSKohCVrWrjZX4myg
GNdOQyZ / TmnQhYpwb95UOom8GzgoD5 + qbTZwpO7XDhEWov6KwXDgAbyWXj70t + fL
9S + JMitrcqeIys95Eax8EphYEUio3o + BxPcRjOagAQ ==
        ----- END CERTIFICATE----- `), // Placeholder certificate
        'tls.key': b64Encode(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDhKLuTkE6ftF6s
hWMjYCkAYKhO9wyouvB + aa + wkyD4LJVhDa3 + Jzxl7Q3kxVGq3JxC + Qvv + Mzs5c9C
UDLj1IckNmZyWs0LmCl7JiWLtFZCeUjisPKP + V2rkdflIQPoqhxfUTqHfOIhIljN
q5lZBxQt0baq9i2BVRBhhxi8R9EDqvysmTvw + EJo2VZDzjgbrz8TFrF / UTP9N1v7
tYIfb5e8AmTgt6GT6zkJym2M0F4uLoJJqQKZ + cAlOgJ9cBodUu2QUZTLZOe9MRwH
OdzA8BsqdBS9Xdx + VoZB / efhRmY5THmoHl3AEzpnVIMWRlE47wiqaePE2qKezd45
Ja0iFc1ZAgMBAAECggEAWyJwDhHeaREKMv1Ie9Sbs + 3rog6n / pGv7xLbDHb45Mqy
6dYuI02BSmYxdfQSErz5uLfyR37mf0qjYd1xQ7tNufAU9ltfXD6KJ7DwiIQFWCGc
STlC7NiLu / qrrq / 205ErK7 + Gl7mKE8xBsnmu95SAa + 1S6Q78qzkhiMA0WF0fMkKc
Y7cM59x1CIWPBVrMapnkdR5DSMb0d6d9BOurtacZeQh8GtRPPZeApjCnBw + FC88v
0UQpRoNzPQzNxyCQa / ERY + 5018CvHjiJyItucq / 4aCIUnH4Dyo / gDRFguIJrtByd
KP9mGBbu7uycFw / ccUF0hC0uMML85i4pc8bL1j + HswKBgQDzxmlJZiQAPMwURzLj
qwAMYaMPI46qqvT6KSqGZ0sIl / w0Eqs1XnB1 / yHug2Q8tUUswBvKddGgwWwU6 / Ly
FY5Xtg4Q4w4NVFUqAiWEPJVb3geWhwEHJB9NMKiyf5v4af8UH08zQhEGpIy1 / h3l
hYZrKQYANJp2kbLPXtwJgkM4zwKBgQDsc1RqlqrFYjEEo041UUraoVP514X5yjLv
OI8UKVaIwqOcqcbv0RrU0jy0hh50sm6vyB4UxQb / t1c3UKRxqlMpeFlT / Vz1X + lT
ehZvasL0dCftVIKh7 / qZzWR1bf0EjwLNtAoCqe3FeUZ5 + JOMbhzTUkm1MTpXa4Dx
vK6ArHxRVwKBgQCt9 + DXgs8aZEj4B7 + nfjdwnpUxjpyX650ckhhJBpojreNMfi40
zgrQCp16i8YTFQIi545ttBs / 8Alj / ObKINwOeFwdbQxwMsj8S7 / eWSX2A8PChuIS
6JJ2Ec2yZSM36t0gzR9GY1WnOfM5Rfqr + 9hrzUD9EI1TJLNJDldVaeLzPQKBgAUT
iMltCKeKLyE5XFF6uE + vTP09KkwtkiBep3u4U3pGK3sOjg3SAHB3PwRlKLw6pHOz
qSmq / TZ6Oi4e1hj2nihyxAAwnVFLSNgY8 + hac2sKH11SBifx3gB1T2XSAa + aXmYK
KnjaKxelPeUaeBh4uLe0uY5hSy5bSX5nHZv3mAerAoGAT2unIbNZK15eyhxx6pO2
BlzXk2DzdM / wWvOD6RR52jLVzqmVa6M4oAbt8WWbfrnXXodQyOM//W5rnSewJuDD
mg + qzYKROXZDwhyuwYnGflvPWb1V5xuf39XeSdQ8XCtJeGOCSDUMG85JZAnQL4K4
5E9Qf / C6AMQwAHC / 7H34S74 =
        ----- END PRIVATE KEY----- `), // Placeholder private key
      },
    });
  });
} 