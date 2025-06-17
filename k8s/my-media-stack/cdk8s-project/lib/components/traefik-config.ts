import { Construct } from 'constructs';
import { ApiObject } from 'cdk8s';

/**
 * Creates the Traefik HelmChartConfig for K3s
 */
export function createTraefikHelmChartConfig(scope: Construct): void {
  new ApiObject(scope, 'traefik-helmchartconfig', {
    apiVersion: 'helm.cattle.io/v1',
    kind: 'HelmChartConfig',
    metadata: {
      name: 'traefik',
      namespace: 'kube-system',
    },
    spec: {
      valuesContent: `additionalArguments:
  - "--global.checknewversion=true"
  - "--global.sendanonymoususage=false"
  - "--log.level=DEBUG"
  - "--api.dashboard=true"
  - "--api.debug=false"
  - "--api.disableDashboardAd=true"
  - "--api.insecure=true"
  - "--ping.terminatingStatusCode=503"
  - "--serverstransport.insecureskipverify=true"
  - "--accesslog.format=json"
  - "--providers.kubernetesingress=true"
  - "--providers.kubernetesingress.allowexternalnameservices=true"
  - "--providers.kubernetescrd=true"
  - "--providers.kubernetescrd.allowexternalnameservices=true"
  - "--certificatesresolvers.letsencrypt.acme.email=boden.crouch@gmail.com"
  - "--certificatesresolvers.letsencrypt.acme.storage=/data/acme.json"
  - "--certificatesresolvers.letsencrypt.acme.dnschallenge.provider=cloudflare"
  - "--certificatesresolvers.letsencrypt.acme.dnschallenge.resolvers=dane.ns.cloudflare.com,tori.ns.cloudflare.com,1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4"
  - "--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web"
  - "--certificatesresolvers.letsencrypt.acme.tlschallenge=true"
  - "--certificatesresolvers.letsencrypt.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
  - "--certificatesresolvers.bolabaden_duckdns_letsencrypt.acme.email=boden.crouch@gmail.com"
  - "--certificatesresolvers.bolabaden_duckdns_letsencrypt.acme.storage=/data/acme_duckdns.json"
  - "--certificatesresolvers.bolabaden_duckdns_letsencrypt.acme.dnschallenge.provider=duckdns"
  - "--certificatesresolvers.bolabaden_duckdns_letsencrypt.acme.dnschallenge.resolvers=1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4"
  - "--certificatesresolvers.bolabaden_duckdns_letsencrypt.acme.httpchallenge.entrypoint=web"
  - "--certificatesresolvers.bolabaden_duckdns_letsencrypt.acme.tlschallenge=true"
  - "--certificatesresolvers.bolabaden_duckdns_letsencrypt.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
  - "--certificatesresolvers.beatapostapita_duckdns_letsencrypt.acme.email=boden.crouch@gmail.com"
  - "--certificatesresolvers.beatapostapita_duckdns_letsencrypt.acme.storage=/data/acme_beatapostapita_duckdns.json"
  - "--certificatesresolvers.beatapostapita_duckdns_letsencrypt.acme.dnschallenge.provider=duckdns"
  - "--certificatesresolvers.beatapostapita_duckdns_letsencrypt.acme.dnschallenge.resolvers=1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4"
  - "--certificatesresolvers.beatapostapita_duckdns_letsencrypt.acme.httpchallenge.entrypoint=web"
  - "--certificatesresolvers.beatapostapita_duckdns_letsencrypt.acme.tlschallenge=true"
  - "--certificatesresolvers.beatapostapita_duckdns_letsencrypt.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"

persistence:
  enabled: true
  size: 1Gi
  path: /data

env:
  - name: DUCKDNS_TOKEN
    value: "a7c4d0ad-114e-40ef-ba1d-d217904a50f2"
  - name: CLOUDFLARE_EMAIL
    value: "boden.crouch@gmail.com"
  - name: CLOUDFLARE_API_KEY
    value: "34a4d9943975b89dd2430870ce146f26aa5bf"

ports:
  web:
    redirections:
      entryPoint:
        to: websecure
        scheme: https
        permanent: true
  websecure:
    tls:
      enabled: true
      certResolver: "beatapostapita_duckdns_letsencrypt"
      domains:
        - main: "beatapostapita.duckdns.org"
          sans:
            - "*.beatapostapita.duckdns.org"
  traefik:
    port: 8080
    expose:
      default: true

ingressRoute:
  dashboard:
    enabled: false

service:
  type: LoadBalancer`,
    },
  });
} 