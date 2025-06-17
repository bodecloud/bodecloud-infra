import { Construct } from 'constructs';
import { KubeDeployment, Quantity, IntOrString } from '../../imports/k8s';

/**
 * Creates the Comet deployment
 */
export function createCometDeployment(scope: Construct): void {
  new KubeDeployment(scope, 'comet-deployment', {
    metadata: {
      name: 'comet',
      namespace: 'my-media-stack',
    },
    spec: {
      replicas: 1,
      selector: {},
      template: {
        spec: {
          containers: [
            {
              name: 'comet',
              image: 'g0ldyy/comet:latest',
              ports: [
                {
                  name: 'comet-2020',
                  containerPort: 2020,
                },
              ],
              env: [
                { name: 'TZ', value: 'America/Chicago' },
                { name: 'ADDON_ID', value: 'comet.bolabaden.org' },
                { name: 'ADDON_NAME', value: 'Comet' },
                { name: 'FASTAPI_HOST', value: '0.0.0.0' },
                { name: 'FASTAPI_PORT', value: '2020' },
                { name: 'FASTAPI_WORKERS', value: '4' },
                { name: 'USE_GUNICORN', value: 'True' },
                { name: 'DASHBOARD_ADMIN_PASSWORD', value: 'h4L0m4St3R327' },
                { name: 'DATABASE_TYPE', value: 'sqlite' },
                { name: 'DATABASE_PATH', value: '/data/comet.db' },
                { name: 'DATABASE_URL', value: 'postgres:postgres@comet-postgres:5432/comet' },
                { name: 'METADATA_CACHE_TTL', value: '2592000' },
                { name: 'TORRENT_CACHE_TTL', value: '1296000' },
                { name: 'DEBRID_CACHE_TTL', value: '86400' },
                { name: 'INDEXER_MANAGER_TYPE', value: 'prowlarr' },
                { name: 'INDEXER_MANAGER_URL', value: 'http://prowlarr:9696' },
                { name: 'INDEXER_MANAGER_API_KEY', value: '29440a82740d475cacb35327c62c87a1' },
                { name: 'INDEXER_MANAGER_TIMEOUT', value: '60' },
                { name: 'INDEXER_MANAGER_INDEXERS', value: '["1337x", "animetosho", "anirena", "limetorrents", "nyaasi", "solidtorrents", "thepiratebay", "torlock", "yts"]' },
                { name: 'GET_TORRENT_TIMEOUT', value: '2' },
                { name: 'DOWNLOAD_TORRENT_FILES', value: 'False' },
                { name: 'SCRAPE_COMET', value: 'false' },
                { name: 'COMET_URL', value: 'https://comet.elfhosted.com' },
                { name: 'SCRAPE_ZILEAN', value: 'true' },
                { name: 'ZILEAN_URL', value: 'https://zilean.elfhosted.com' },
                { name: 'SCRAPE_TORRENTIO', value: 'false' },
                { name: 'TORRENTIO_URL', value: 'https://torrentio.strem.fun' },
                { name: 'SCRAPE_MEDIAFUSION', value: 'false' },
                { name: 'MEDIAFUSION_URL', value: 'https://mediafusion.elfhosted.com' },
                { name: 'STREMTHRU_URL', value: 'https://stremthru.bolabaden.org' },
                { name: 'PROXY_DEBRID_STREAM', value: 'True' },
                { name: 'PROXY_DEBRID_STREAM_PASSWORD', value: 'h4L0m4St3R327' },
                { name: 'PROXY_DEBRID_STREAM_MAX_CONNECTIONS', value: '-1' },
                { name: 'PROXY_DEBRID_STREAM_DEBRID_DEFAULT_SERVICE', value: 'premiumize' },
                { name: 'PROXY_DEBRID_STREAM_DEBRID_DEFAULT_APIKEY', value: '6i6zgswm35baj3ur' },
                { name: 'REMOVE_ADULT_CONTENT', value: 'True' },
              ],
              volumeMounts: [
                {
                  name: 'comet-data',
                  mountPath: '/data',
                },
              ],
              resources: {
                requests: {
                  memory: Quantity.fromString('512Mi'),
                  cpu: Quantity.fromString('200m'),
                },
                limits: {
                  memory: Quantity.fromString('1Gi'),
                  cpu: Quantity.fromString('1000m'),
                },
              },
              livenessProbe: {
                httpGet: {
                  path: '/health',
                  port: IntOrString.fromNumber(2020),
                },
                initialDelaySeconds: 60,
                periodSeconds: 30,
                timeoutSeconds: 10,
                failureThreshold: 5,
              },
              readinessProbe: {
                httpGet: {
                  path: '/health',
                  port: IntOrString.fromNumber(2020),
                },
                initialDelaySeconds: 10,
                periodSeconds: 10,
                timeoutSeconds: 5,
                failureThreshold: 3,
              },
            },
          ],
          volumes: [
            {
              name: 'comet-data',
              hostPath: {
                path: '/home/ubuntu/my-media-stack/configs/stremio/addons/comet/data',
                type: 'DirectoryOrCreate',
              },
            },
          ],
        },
      },
    },
  });
} 