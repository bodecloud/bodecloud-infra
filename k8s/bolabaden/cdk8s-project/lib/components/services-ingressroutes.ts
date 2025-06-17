import { Construct } from 'constructs';
import { ApiObject } from 'cdk8s';

/**
 * Creates IngressRoutes for all media stack services
 */
export function createServicesIngressRoutes(scope: Construct): void {
  const services = [
    { name: 'dozzle', port: 8080, subdomain: 'dozzle' },
    { name: 'gptr', port: 8000, subdomain: 'gptr' },
    { name: 'homepage', port: 3000, subdomain: 'homepage' },
    { name: 'jackett', port: 9117, subdomain: 'jackett' },
    { name: 'mediafusion', port: 8000, subdomain: 'mediafusion' },
    { name: 'prowlarr', port: 9696, subdomain: 'prowlarr' },
    { name: 'searxng', port: 8080, subdomain: 'searxng' },
    { name: 'speedtest', port: 80, subdomain: 'speedtest' },
    { name: 'stremio', port: 11470, subdomain: 'stremio' },
    { name: 'whoami', port: 80, subdomain: 'whoami' },
  ];

  services.forEach((service) => {
    new ApiObject(scope, `${service.name}-ingressroute`, {
      apiVersion: 'traefik.io/v1alpha1',
      kind: 'IngressRoute',
      metadata: {
        name: `${service.name}-ingressroute`,
        namespace: 'my-media-stack',
      },
      spec: {
        entryPoints: ['websecure'],
        routes: [
          {
            match: `Host(\`${service.subdomain}.beatapostapita.duckdns.org\`)`,
            kind: 'Rule',
            services: [
              {
                name: service.name,
                port: service.port,
              },
            ],
          },
        ],
        tls: {
          certResolver: 'beatapostapita_duckdns_letsencrypt',
        },
      },
    });
  });
} 