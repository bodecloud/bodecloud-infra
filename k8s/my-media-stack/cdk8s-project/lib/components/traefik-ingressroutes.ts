import { Construct } from 'constructs';
import { ApiObject } from 'cdk8s';

/**
 * Creates Traefik IngressRoutes and middleware
 */
export function createTraefikIngressRoutes(scope: Construct): void {
  // Traefik dashboard IngressRoute (HTTPS)
  new ApiObject(scope, 'traefik-dashboard-ingressroute', {
    apiVersion: 'traefik.io/v1alpha1',
    kind: 'IngressRoute',
    metadata: {
      name: 'traefik-dashboard',
      namespace: 'my-media-stack',
    },
    spec: {
      entryPoints: ['websecure'],
      routes: [
        {
          match: 'Host(`traefik.beatapostapita.duckdns.org`)',
          kind: 'Rule',
          services: [
            {
              name: 'traefik-dashboard',
              port: 8080,
            },
          ],
        },
      ],
      tls: {
        certResolver: 'beatapostapita_duckdns_letsencrypt',
      },
    },
  });

  // Traefik dashboard redirect IngressRoute (HTTP to HTTPS)
  new ApiObject(scope, 'traefik-dashboard-redirect', {
    apiVersion: 'traefik.io/v1alpha1',
    kind: 'IngressRoute',
    metadata: {
      name: 'traefik-dashboard-redirect',
      namespace: 'my-media-stack',
    },
    spec: {
      entryPoints: ['web'],
      routes: [
        {
          match: 'Host(`traefik.beatapostapita.duckdns.org`)',
          kind: 'Rule',
          services: [
            {
              name: 'traefik-dashboard',
              port: 8080,
            },
          ],
          middlewares: [
            {
              name: 'redirect-to-https',
              namespace: 'my-media-stack',
            },
          ],
        },
      ],
    },
  });

  // HTTPS redirect middleware
  new ApiObject(scope, 'redirect-to-https-middleware', {
    apiVersion: 'traefik.io/v1alpha1',
    kind: 'Middleware',
    metadata: {
      name: 'redirect-to-https',
      namespace: 'my-media-stack',
    },
    spec: {
      redirectScheme: {
        scheme: 'https',
        permanent: true,
      },
    },
  });
} 