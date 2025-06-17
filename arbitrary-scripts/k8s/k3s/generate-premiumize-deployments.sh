#!/bin/bash
set -euo pipefail

# List of Premiumize server countries
COUNTRIES=(
    "at" "au" "be" "ca" "ch" "cz" "de" "es" "fi" "fr" "gb" "gr" "it" "jp" "nl" "pl" "sg" "us"
)

# Base directory for configs
CONFIG_BASE="/home/ubuntu/my-media-stack/configs/gluetun/premiumize"

# Function to create deployment for a country
create_deployment() {
    local country="$1"
    local config_dir="${CONFIG_BASE}/config"
    local config_file="vpn-${country}.premiumize.me.ovpn"
    
    echo "Creating deployment for Premiumize ${country^^}..."
    
    # Create config directory if it doesn't exist
    mkdir -p "${CONFIG_BASE}-${country}"
    
    # Copy the config file if it exists
    if [[ -f "${config_dir}/${config_file}" ]]; then
        cp "${config_dir}/${config_file}" "${CONFIG_BASE}-${country}/premiumize-${country}.ovpn"
        echo "Copied config file for ${country}"
    else
        echo "Warning: Config file ${config_file} not found, skipping ${country}"
        return
    fi
    
    # Generate the deployment YAML
    cat > "k8s/my-media-stack/gluetun-premiumize-${country}-deployment.yaml" << EOF
---
# Gluetun Premiumize ${country^^} Gateway with Dynamic Hostname Resolution
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gluetun-premiumize-${country}
  namespace: vpn-gateway
  labels:
    app: gluetun-premiumize-${country}
    component: vpn-gateway
    country: ${country}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gluetun-premiumize-${country}
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: gluetun-premiumize-${country}
        component: vpn-gateway
        country: ${country}
    spec:
      securityContext:
        sysctls:
          - name: net.ipv6.conf.all.disable_ipv6
            value: "1"
      initContainers:
        - name: hostname-resolver
          image: busybox:latest
          command:
            - /bin/sh
            - -c
            - |
              echo "=== Dynamic Hostname Resolution for Premiumize ${country^^} ==="
              
              config_file="/gluetun/premiumize-${country}.ovpn"
              backup_file="/gluetun/premiumize-${country}.ovpn.backup"
              
              # Create backup if it doesn't exist
              if [[ ! -f "\$backup_file" ]]; then
                cp "\$config_file" "\$backup_file"
                echo "Created backup: \$backup_file"
              fi
              
              # Extract hostname from the backup (original) file
              hostname=\$(grep "^remote " "\$backup_file" | awk '{print \$2}' | head -1)
              
              if [[ -z "\$hostname" ]]; then
                echo "No hostname found in \$config_file"
                exit 1
              fi
              
              echo "Resolving hostname: \$hostname"
              
              # Resolve hostname to IP with retries
              ip=""
              retries=5
              for ((i=1; i<=retries; i++)); do
                if ip=\$(getent hosts "\$hostname" | awk '{print \$1}' | head -1); then
                  if [[ -n "\$ip" && "\$ip" =~ ^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+\$ ]]; then
                    echo "Resolved \$hostname -> \$ip (attempt \$i)"
                    break
                  fi
                fi
                echo "Resolution attempt \$i failed, retrying..."
                sleep 2
              done
              
              if [[ -z "\$ip" ]]; then
                echo "ERROR: Failed to resolve \$hostname after \$retries attempts"
                # Use original config as fallback
                cp "\$backup_file" "\$config_file"
                exit 1
              fi
              
              # Update config file with resolved IP
              cp "\$backup_file" "\$config_file"
              sed -i "s/^remote \$hostname/remote \$ip/" "\$config_file"
              
              # Keep verify-x509-name with original hostname for certificate validation
              if grep -q "^verify-x509-name" "\$config_file"; then
                sed -i "s/^verify-x509-name CN=\$ip/verify-x509-name CN=\$hostname/" "\$config_file"
              fi
              
              echo "Updated \$config_file: \$hostname -> \$ip"
              echo "=== Hostname resolution complete for ${country^^} ==="
          volumeMounts:
            - name: gluetun-config
              mountPath: /gluetun
      containers:
        - name: gluetun
          image: ghcr.io/qdm12/gluetun:latest
          securityContext:
            capabilities:
              add:
                - NET_ADMIN
            privileged: true
          env:
            - name: TZ
              value: "America/Chicago"
            - name: VPN_SERVICE_PROVIDER
              value: "custom"
            - name: VPN_TYPE
              value: "openvpn"
            - name: OPENVPN_CUSTOM_CONFIG
              value: "/gluetun/premiumize-${country}.ovpn"
            - name: OPENVPN_USER
              value: "117274388"
            - name: OPENVPN_PASSWORD
              value: "6i6zgswm35baj3ur"
            - name: BLOCK_ADS
              value: "on"
            - name: BLOCK_MALICIOUS
              value: "on"
            - name: BLOCK_SURVEILLANCE
              value: "on"
            - name: FIREWALL_OUTBOUND_SUBNETS
              value: "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
            - name: HEALTH_SERVER_ADDRESS
              value: "0.0.0.0:9999"
            - name: HEALTH_SUCCESS_WAIT_DURATION
              value: "5s"
            - name: HEALTH_TARGET_ADDRESS
              value: "cloudflare.com:443"
            - name: HEALTH_VPN_DURATION_ADDITION
              value: "5s"
            - name: HEALTH_VPN_DURATION_INITIAL
              value: "6s"
            - name: OPENVPN_IPV6
              value: "off"
            - name: VERSION_INFORMATION
              value: "on"
            - name: HTTP_PROXY
              value: "on"
            - name: HTTP_PROXY_ADDRESS
              value: "0.0.0.0:8888"
            - name: SHADOWSOCKS
              value: "on"
            - name: SHADOWSOCKS_ADDRESS
              value: "0.0.0.0:8388"
          ports:
            - containerPort: 8888
              name: http-proxy
            - containerPort: 8388
              name: shadowsocks
            - containerPort: 9999
              name: health
          volumeMounts:
            - name: gluetun-config
              mountPath: /gluetun
          livenessProbe:
            httpGet:
              path: /
              port: 9999
            initialDelaySeconds: 60
            periodSeconds: 30
            timeoutSeconds: 10
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /
              port: 9999
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          resources:
            requests:
              memory: "128Mi"
              cpu: "50m"
            limits:
              memory: "512Mi"
              cpu: "200m"
      volumes:
        - name: gluetun-config
          hostPath:
            path: /home/ubuntu/my-media-stack/configs/gluetun/premiumize-${country}
            type: DirectoryOrCreate

---
# Service for Gluetun Premiumize ${country^^}
apiVersion: v1
kind: Service
metadata:
  name: gluetun-premiumize-${country}
  namespace: vpn-gateway
  labels:
    app: gluetun-premiumize-${country}
    component: vpn-gateway
    country: ${country}
spec:
  selector:
    app: gluetun-premiumize-${country}
  ports:
    - name: http-proxy
      port: 8888
      targetPort: 8888
    - name: shadowsocks
      port: 8388
      targetPort: 8388
    - name: health
      port: 9999
      targetPort: 9999
EOF

    echo "Created deployment file: gluetun-premiumize-${country}-deployment.yaml"
}

# Main execution
echo "Generating Premiumize VPN deployments for all countries..."

for country in "${COUNTRIES[@]}"; do
    create_deployment "$country"
done

echo ""
echo "=== Generation Complete ==="
echo "Created deployments for: ${COUNTRIES[*]}"
echo ""
echo "To deploy all Premiumize VPN gateways, run:"
echo "kubectl apply -f k8s/my-media-stack/gluetun-premiumize-*-deployment.yaml" 