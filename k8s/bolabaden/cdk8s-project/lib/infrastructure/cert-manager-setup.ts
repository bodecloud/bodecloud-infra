import { Construct } from 'constructs';
import { KubeJob, Quantity } from '../../imports/k8s';

/**
 * Creates the cert-manager setup job that waits for cert-manager to be ready
 * and then creates the ClusterIssuer and Certificate
 */
export function createCertManagerSetup(scope: Construct): void {
  new KubeJob(scope, 'cert-manager-setup', {
    metadata: {
      name: 'cert-manager-setup',
      namespace: 'cert-manager',
      labels: {
        app: 'cert-manager-setup',
      },
    },
    spec: {
      template: {
        spec: {
          restartPolicy: 'OnFailure',
          serviceAccountName: 'cert-manager',
          containers: [
            {
              name: 'setup',
              image: 'bitnami/kubectl:latest',
              command: ['/bin/bash', '-c'],
              args: [
                `set -e
echo "Waiting for cert-manager to be ready..."
kubectl wait --for=condition=available deployment/cert-manager -n cert-manager --timeout=300s

echo "Waiting for cert-manager CRDs to be established..."
kubectl wait --for condition=established --timeout=60s crd/clusterissuers.cert-manager.io
kubectl wait --for condition=established --timeout=60s crd/certificates.cert-manager.io

echo "Creating vpn-gateway namespace if it doesn't exist..."
kubectl create namespace vpn-gateway --dry-run=client -o yaml | kubectl apply -f -

echo "Creating self-signed ClusterIssuer..."
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
EOF

echo "Waiting for ClusterIssuer to be ready..."
sleep 5

echo "Creating proxy-injector certificate..."
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: proxy-injector-cert
  namespace: vpn-gateway
spec:
  secretName: proxy-injector-certs
  issuerRef:
    name: selfsigned-issuer
    kind: ClusterIssuer
  dnsNames:
  - proxy-injector.vpn-gateway.svc
  - proxy-injector.vpn-gateway.svc.cluster.local
EOF

echo "Waiting for certificate to be ready..."
kubectl wait --for=condition=ready certificate/proxy-injector-cert -n vpn-gateway --timeout=120s

echo "cert-manager setup completed successfully!"`,
              ],
              resources: {
                requests: {
                  cpu: Quantity.fromString('10m'),
                  memory: Quantity.fromString('32Mi'),
                },
                limits: {
                  cpu: Quantity.fromString('100m'),
                  memory: Quantity.fromString('128Mi'),
                },
              },
            },
          ],
        },
      },
    },
  });
} 