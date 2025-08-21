#!/usr/bin/env python3
"""
VPN Gateway Deployment Generator
Generates Kubernetes manifests for all VPN gateways in the pod-gateway system
"""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml

# VPN Gateway configurations
VPN_GATEWAYS: list[dict[str, Any]] = [
    {"name": "gluetun-premiumize-at", "priority": 19, "type": "premiumize", "country": "at"},
    {"name": "gluetun-premiumize-au", "priority": 18, "type": "premiumize", "country": "au"},
    {"name": "gluetun-premiumize-be", "priority": 17, "type": "premiumize", "country": "be"},
    {"name": "gluetun-premiumize-ca", "priority": 16, "type": "premiumize", "country": "ca"},
    {"name": "gluetun-premiumize-ch", "priority": 15, "type": "premiumize", "country": "ch"},
    {"name": "gluetun-premiumize-cz", "priority": 14, "type": "premiumize", "country": "cz"},
    {"name": "gluetun-premiumize-de", "priority": 13, "type": "premiumize", "country": "de"},
    {"name": "gluetun-premiumize-es", "priority": 12, "type": "premiumize", "country": "es"},
    {"name": "gluetun-premiumize-fi", "priority": 11, "type": "premiumize", "country": "fi"},
    {"name": "gluetun-premiumize-fr", "priority": 10, "type": "premiumize", "country": "fr"},
    {"name": "gluetun-premiumize-gb", "priority": 9, "type": "premiumize", "country": "gb"},
    {"name": "gluetun-premiumize-gr", "priority": 8, "type": "premiumize", "country": "gr"},
    {"name": "gluetun-premiumize-it", "priority": 7, "type": "premiumize", "country": "it"},
    {"name": "gluetun-premiumize-jp", "priority": 6, "type": "premiumize", "country": "jp"},
    {"name": "gluetun-premiumize-nl", "priority": 3, "type": "premiumize", "country": "nl"},
    {"name": "gluetun-premiumize-pl", "priority": 5, "type": "premiumize", "country": "pl"},
    {"name": "gluetun-premiumize-sg", "priority": 4, "type": "premiumize", "country": "sg"},
    {"name": "gluetun-premiumize-us", "priority": 1, "type": "premiumize", "country": "us"},
    {"name": "gluetun-airvpn", "priority": 2, "type": "airvpn", "country": "us"},
    {"name": "warp-gateway", "priority": 20, "type": "warp", "country": "us"},
]


def generate_premiumize_deployment(gateway_config: dict[str, Any]) -> dict[str, Any]:
    """Generate deployment for Premiumize VPN gateway"""
    name: str = gateway_config["name"]
    country: str = gateway_config["country"]
    priority: int = gateway_config["priority"]

    return {
        "apiVersion": "apps/v1",
        "kind": "Deployment",
        "metadata": {
            "name": name,
            "namespace": "vpn-gateway",
            "labels": {
                "app": name,
                "component": "vpn-gateway",
                "priority": str(priority),
            },
        },
        "spec": {
            "replicas": 0,  # Managed by controller
            "selector": {"matchLabels": {"app": name}},
            "template": {
                "metadata": {"labels": {"app": name, "component": "vpn-gateway"}},
                "spec": {
                    "securityContext": {
                        "sysctls": [
                            {"name": "net.ipv6.conf.all.disable_ipv6", "value": "1"}
                        ]
                    },
                    "initContainers": [
                        {
                            "name": "hostname-resolver",
                            "image": "alpine:latest",
                            "command": [
                                "/bin/sh",
                                "-c",
                                f"""
                                apk add --no-cache bind-tools
                                
                                COUNTRY="{country}"
                                HOSTNAME="vpn-${{COUNTRY}}.premiumize.me"
                                CONFIG_FILE="/gluetun/premiumize-${{COUNTRY}}.ovpn"
                                AUTH_FILE="/gluetun/auth.conf"
                                
                                echo "=== Resolving $HOSTNAME ==="
                                
                                # Resolve hostname with retries
                                IP=""
                                for i in $(seq 1 10); do
                                    if IP=$(dig +short $HOSTNAME | head -1); then
                                        if echo "$IP" | grep -E '^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$'; then
                                            echo "✓ Resolved $HOSTNAME -> $IP (attempt $i)"
                                            break
                                        fi
                                    fi
                                    echo "Resolution attempt $i failed, retrying in 2s..."
                                    sleep 2
                                done
                                
                                if [ -z "$IP" ]; then
                                    echo "❌ Failed to resolve $HOSTNAME after 10 attempts"
                                    exit 1
                                fi
                                
                                # Create OpenVPN config
                                cat > "$CONFIG_FILE" << 'EOF'
                                remote $IP 1194
                                verify-x509-name CN=$HOSTNAME
                                client
                                dev tun
                                proto udp
                                cipher AES-256-CBC
                                resolv-retry infinite
                                nobind
                                persist-key
                                persist-tun
                                auth-user-pass /gluetun/auth.conf
                                verb 3
                                auth SHA256
                                key-direction 1
                                <ca>
                                -----BEGIN CERTIFICATE-----
                                MIIFqzCCA5OgAwIBAgIJAKZ7D5Kv2XSUMA0GCSqGSIb3DQEBCwUAMGwxCzAJBgNV
                                BAYTAkRFMRAwDgYDVQQIDAdCYXZhcmlhMQ8wDQYDVQQHDAZNdW5pY2gxEzARBgNV
                                BAoMClByZW1pdW1pemUxJTAjBgNVBAMMHFByZW1pdW1pemUgQ2VydGlmaWNhdGUg
                                QXV0aG9yaXR5MB4XDTE5MDkxOTE0NDUwNFoXDTI5MDkxNjE0NDUwNFowbDELMAkG
                                A1UEBhMCREUxEDAOBgNVBAgMB0JhdmFyaWExDzANBgNVBAcMBk11bmljaDETMBEG
                                A1UECgwKUHJlbWl1bWl6ZTElMCMGA1UEAwwcUHJlbWl1bWl6ZSBDZXJ0aWZpY2F0
                                ZSBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDMsJwI
                                -----END CERTIFICATE-----
                                </ca>
                                <tls-auth>
                                -----BEGIN OpenVPN Static key V1-----
                                # Add your actual tls-auth key here
                                -----END OpenVPN Static key V1-----
                                </tls-auth>
                                EOF
                                
                                # Replace placeholders
                                sed -i "s/\\$IP/$IP/g" "$CONFIG_FILE"
                                sed -i "s/\\$HOSTNAME/$HOSTNAME/g" "$CONFIG_FILE"
                                
                                # Create auth file
                                echo "117274388" > "$AUTH_FILE"
                                echo "6i6zgswm35baj3ur" >> "$AUTH_FILE"
                                
                                echo "✅ Configuration complete for {country.upper()}"
                                """,
                            ],
                            "volumeMounts": [
                                {"name": "gluetun-config", "mountPath": "/gluetun"}
                            ],
                        },
                        {
                            "name": "gateway-init",
                            "image": "ghcr.io/k8s-at-home/pod-gateway:latest",
                            "securityContext": {
                                "capabilities": {"add": ["NET_ADMIN"]},
                                "privileged": True,
                            },
                            "command": [
                                "/bin/bash",
                                "-c",
                                """
                                source /config/settings.sh
                                
                                # Enable IP forwarding
                                echo 1 > /proc/sys/net/ipv4/ip_forward
                                
                                # Create VXLAN interface
                                VXLAN_GATEWAY_IP="${VXLAN_IP_NETWORK}.1"
                                ip link add vxlan0 type vxlan id $VXLAN_ID dev eth0 dstport $VXLAN_PORT || true
                                ip addr add ${VXLAN_GATEWAY_IP}/24 dev vxlan0 || true
                                ip link set up dev vxlan0
                                ip link set mtu $VPN_INTERFACE_MTU dev vxlan0 || true
                                
                                # Set routing rules
                                ip rule add from all lookup main suppress_prefixlength 0 preference 50 || true
                                
                                # Configure iptables for NAT and forwarding
                                iptables -t nat -A POSTROUTING -j MASQUERADE || true
                                iptables -A INPUT -i eth0 -p udp --dport=$VXLAN_PORT -j ACCEPT || true
                                iptables -A INPUT -i vxlan0 -p udp --sport=68 --dport=67 -j ACCEPT || true
                                
                                echo "Gateway init complete"
                                """,
                            ],
                            "volumeMounts": [
                                {"name": "config", "mountPath": "/config"}
                            ],
                        },
                    ],
                    "containers": [
                        {
                            "name": "gluetun",
                            "image": "ghcr.io/qdm12/gluetun:latest",
                            "securityContext": {
                                "capabilities": {"add": ["NET_ADMIN"]},
                                "privileged": True,
                            },
                            "env": [
                                {"name": "VPN_SERVICE_PROVIDER", "value": "custom"},
                                {"name": "VPN_TYPE", "value": "openvpn"},
                                {
                                    "name": "OPENVPN_CUSTOM_CONFIG",
                                    "value": f"/gluetun/premiumize-{country}.ovpn",
                                },
                                {"name": "HTTP_PROXY", "value": "on"},
                                {"name": "HTTP_PROXY_ADDRESS", "value": "0.0.0.0:8888"},
                                {"name": "SHADOWSOCKS", "value": "on"},
                                {
                                    "name": "SHADOWSOCKS_ADDRESS",
                                    "value": "0.0.0.0:8388",
                                },
                                {
                                    "name": "HEALTH_SERVER_ADDRESS",
                                    "value": "0.0.0.0:9999",
                                },
                            ],
                            "ports": [
                                {"containerPort": 8888, "name": "http-proxy"},
                                {"containerPort": 8388, "name": "shadowsocks"},
                                {"containerPort": 9999, "name": "health"},
                            ],
                            "volumeMounts": [
                                {"name": "gluetun-config", "mountPath": "/gluetun"}
                            ],
                            "livenessProbe": {
                                "httpGet": {"path": "/", "port": 9999},
                                "initialDelaySeconds": 60,
                                "periodSeconds": 30,
                            },
                            "readinessProbe": {
                                "httpGet": {"path": "/", "port": 9999},
                                "initialDelaySeconds": 30,
                                "periodSeconds": 10,
                            },
                            "resources": {
                                "requests": {"memory": "128Mi", "cpu": "100m"},
                                "limits": {"memory": "512Mi", "cpu": "300m"},
                            },
                        },
                        {
                            "name": "gateway-sidecar",
                            "image": "ghcr.io/k8s-at-home/pod-gateway:latest",
                            "command": [
                                "/bin/bash",
                                "-c",
                                """
                                source /config/settings.sh
                                
                                # Start DHCP server for VXLAN clients
                                cat > /tmp/dnsmasq.conf << EOF
                                interface=vxlan0
                                bind-interfaces
                                dhcp-range=${VXLAN_IP_NETWORK}.${VXLAN_GATEWAY_FIRST_DYNAMIC_IP},${VXLAN_IP_NETWORK}.254,12h
                                dhcp-option=3,${VXLAN_IP_NETWORK}.1
                                dhcp-option=6,${VXLAN_IP_NETWORK}.1
                                server=10.96.0.10
                                EOF
                                
                                # Start dnsmasq for DHCP and DNS
                                dnsmasq --conf-file=/tmp/dnsmasq.conf --no-daemon &
                                
                                # Health monitoring
                                while true; do
                                  if ip link show tun0 >/dev/null 2>&1; then
                                    if timeout 5 curl -s --interface tun0 http://httpbin.org/ip >/dev/null 2>&1; then
                                      echo "$(date): VPN Gateway healthy - tun0 connected"
                                    else
                                      echo "$(date): VPN Gateway unhealthy - no external connectivity"
                                    fi
                                  else
                                    echo "$(date): VPN Gateway unhealthy - tun0 not found"
                                  fi
                                  sleep 30
                                done
                                """,
                            ],
                            "volumeMounts": [
                                {"name": "config", "mountPath": "/config"}
                            ],
                            "resources": {
                                "requests": {"memory": "64Mi", "cpu": "50m"},
                                "limits": {"memory": "128Mi", "cpu": "100m"},
                            },
                        },
                    ],
                    "volumes": [
                        {
                            "name": "config",
                            "configMap": {
                                "name": "pod-gateway-config",
                                "defaultMode": 0o755,
                            },
                        },
                        {
                            "name": "gluetun-config",
                            "hostPath": {
                                "path": f"/home/ubuntu/my-media-stack/configs/gluetun/premiumize/{country}",
                                "type": "DirectoryOrCreate",
                            },
                        },
                    ],
                },
            },
        },
    }


def generate_airvpn_deployment(gateway_config: dict[str, Any]) -> dict[str, Any]:
    """Generate deployment for AirVPN gateway"""
    name: str = gateway_config["name"]
    priority: int = gateway_config["priority"]

    return {
        "apiVersion": "apps/v1",
        "kind": "Deployment",
        "metadata": {
            "name": name,
            "namespace": "vpn-gateway",
            "labels": {
                "app": name,
                "component": "vpn-gateway",
                "priority": str(priority),
            },
        },
        "spec": {
            "replicas": 0,
            "selector": {"matchLabels": {"app": name}},
            "template": {
                "metadata": {"labels": {"app": name, "component": "vpn-gateway"}},
                "spec": {
                    "securityContext": {
                        "sysctls": [
                            {"name": "net.ipv6.conf.all.disable_ipv6", "value": "1"}
                        ]
                    },
                    "initContainers": [
                        {
                            "name": "gateway-init",
                            "image": "ghcr.io/k8s-at-home/pod-gateway:latest",
                            "securityContext": {
                                "capabilities": {"add": ["NET_ADMIN"]},
                                "privileged": True,
                            },
                            "command": [
                                "/bin/bash",
                                "-c",
                                """
                                source /config/settings.sh
                                
                                # Enable IP forwarding
                                echo 1 > /proc/sys/net/ipv4/ip_forward
                                
                                # Create VXLAN interface
                                VXLAN_GATEWAY_IP="${VXLAN_IP_NETWORK}.1"
                                ip link add vxlan0 type vxlan id $VXLAN_ID dev eth0 dstport $VXLAN_PORT || true
                                ip addr add ${VXLAN_GATEWAY_IP}/24 dev vxlan0 || true
                                ip link set up dev vxlan0
                                ip link set mtu $VPN_INTERFACE_MTU dev vxlan0 || true
                                
                                # Set routing rules
                                ip rule add from all lookup main suppress_prefixlength 0 preference 50 || true
                                
                                # Configure iptables for NAT and forwarding
                                iptables -t nat -A POSTROUTING -j MASQUERADE || true
                                iptables -A INPUT -i eth0 -p udp --dport=$VXLAN_PORT -j ACCEPT || true
                                iptables -A INPUT -i vxlan0 -p udp --sport=68 --dport=67 -j ACCEPT || true
                                
                                echo "Gateway init complete"
                                """,
                            ],
                            "volumeMounts": [
                                {"name": "config", "mountPath": "/config"}
                            ],
                        }
                    ],
                    "containers": [
                        {
                            "name": "gluetun",
                            "image": "ghcr.io/qdm12/gluetun:latest",
                            "securityContext": {
                                "capabilities": {"add": ["NET_ADMIN"]},
                                "privileged": True,
                            },
                            "env": [
                                {"name": "VPN_SERVICE_PROVIDER", "value": "airvpn"},
                                {"name": "OPENVPN_CIPHERS", "value": "aes-256-gcm"},
                                {"name": "SERVER_COUNTRIES", "value": "United States"},
                                {"name": "HTTP_PROXY", "value": "on"},
                                {"name": "HTTP_PROXY_ADDRESS", "value": "0.0.0.0:8888"},
                                {"name": "SHADOWSOCKS", "value": "on"},
                                {
                                    "name": "SHADOWSOCKS_ADDRESS",
                                    "value": "0.0.0.0:8388",
                                },
                                {
                                    "name": "HEALTH_SERVER_ADDRESS",
                                    "value": "0.0.0.0:9999",
                                },
                            ],
                            "ports": [
                                {"containerPort": 8888, "name": "http-proxy"},
                                {"containerPort": 8388, "name": "shadowsocks"},
                                {"containerPort": 9999, "name": "health"},
                            ],
                            "volumeMounts": [
                                {"name": "gluetun-config", "mountPath": "/gluetun"}
                            ],
                            "livenessProbe": {
                                "httpGet": {"path": "/", "port": 9999},
                                "initialDelaySeconds": 60,
                                "periodSeconds": 30,
                            },
                            "readinessProbe": {
                                "httpGet": {"path": "/", "port": 9999},
                                "initialDelaySeconds": 30,
                                "periodSeconds": 10,
                            },
                            "resources": {
                                "requests": {"memory": "128Mi", "cpu": "100m"},
                                "limits": {"memory": "512Mi", "cpu": "300m"},
                            },
                        },
                        {
                            "name": "gateway-sidecar",
                            "image": "ghcr.io/k8s-at-home/pod-gateway:latest",
                            "command": [
                                "/bin/bash",
                                "-c",
                                """
                                source /config/settings.sh
                                
                                # Start DHCP server for VXLAN clients
                                cat > /tmp/dnsmasq.conf << EOF
                                interface=vxlan0
                                bind-interfaces
                                dhcp-range=${VXLAN_IP_NETWORK}.${VXLAN_GATEWAY_FIRST_DYNAMIC_IP},${VXLAN_IP_NETWORK}.254,12h
                                dhcp-option=3,${VXLAN_IP_NETWORK}.1
                                dhcp-option=6,${VXLAN_IP_NETWORK}.1
                                server=10.96.0.10
                                EOF
                                
                                # Start dnsmasq for DHCP and DNS
                                dnsmasq --conf-file=/tmp/dnsmasq.conf --no-daemon &
                                
                                # Health monitoring
                                while true; do
                                  if ip link show tun0 >/dev/null 2>&1; then
                                    if timeout 5 curl -s --interface tun0 http://httpbin.org/ip >/dev/null 2>&1; then
                                      echo "$(date): VPN Gateway healthy - tun0 connected"
                                    else
                                      echo "$(date): VPN Gateway unhealthy - no external connectivity"
                                    fi
                                  else
                                    echo "$(date): VPN Gateway unhealthy - tun0 not found"
                                  fi
                                  sleep 30
                                done
                                """,
                            ],
                            "volumeMounts": [
                                {"name": "config", "mountPath": "/config"}
                            ],
                            "resources": {
                                "requests": {"memory": "64Mi", "cpu": "50m"},
                                "limits": {"memory": "128Mi", "cpu": "100m"},
                            },
                        },
                    ],
                    "volumes": [
                        {
                            "name": "config",
                            "configMap": {
                                "name": "pod-gateway-config",
                                "defaultMode": 0o755,
                            },
                        },
                        {
                            "name": "gluetun-config",
                            "hostPath": {
                                "path": "/home/ubuntu/my-media-stack/configs/gluetun/airvpn",
                                "type": "DirectoryOrCreate",
                            },
                        },
                    ],
                },
            },
        },
    }


def generate_warp_deployment(gateway_config: dict[str, Any]) -> dict[str, Any]:
    """Generate deployment for WARP gateway"""
    name = gateway_config["name"]
    priority = gateway_config["priority"]

    return {
        "apiVersion": "apps/v1",
        "kind": "Deployment",
        "metadata": {
            "name": name,
            "namespace": "vpn-gateway",
            "labels": {
                "app": name,
                "component": "vpn-gateway",
                "priority": str(priority),
            },
        },
        "spec": {
            "replicas": 0,
            "selector": {"matchLabels": {"app": name}},
            "template": {
                "metadata": {"labels": {"app": name, "component": "vpn-gateway"}},
                "spec": {
                    "securityContext": {
                        "sysctls": [
                            {"name": "net.ipv6.conf.all.disable_ipv6", "value": "1"},
                            {"name": "net.ipv4.conf.all.src_valid_mark", "value": "1"},
                        ]
                    },
                    "containers": [
                        {
                            "name": "warp",
                            "image": "caomingjun/warp:latest",
                            "securityContext": {
                                "capabilities": {
                                    "add": ["NET_ADMIN", "MKNOD", "AUDIT_WRITE"]
                                },
                                "privileged": True,
                            },
                            "env": [
                                {"name": "WARP_SLEEP", "value": "2"},
                                {
                                    "name": "WARP_LICENSE_KEY",
                                    "value": "eyJhIjoiZTRlYjNkYmViMTJiZWFhY2MxNzcwNDEyMzE3OTA0NTQiLCJ0IjoiM2ExODhhOTMtYzQwNC00Zjg5LTg4NzItOThlMDkxNjNiYzAzIiwicyI6Ill6SmpNVFl5WmpNdE9EYzROUzAwWlRrMUxUazFORFl0TWpnd1lXWXhZVEpsT1dNMiJ9",
                                },
                            ],
                            "ports": [
                                {"containerPort": 1080, "name": "socks5-proxy"},
                                {"containerPort": 3128, "name": "http-proxy"},
                            ],
                            "volumeMounts": [
                                {
                                    "name": "warp-data",
                                    "mountPath": "/var/lib/cloudflare-warp",
                                }
                            ],
                            "livenessProbe": {
                                "exec": {
                                    "command": [
                                        "/bin/sh",
                                        "-c",
                                        "curl -s --max-time 5 --proxy socks5://127.0.0.1:1080 http://httpbin.org/ip >/dev/null",
                                    ]
                                },
                                "initialDelaySeconds": 60,
                                "periodSeconds": 30,
                            },
                            "readinessProbe": {
                                "exec": {
                                    "command": [
                                        "/bin/sh",
                                        "-c",
                                        "curl -s --max-time 5 --proxy socks5://127.0.0.1:1080 http://httpbin.org/ip >/dev/null",
                                    ]
                                },
                                "initialDelaySeconds": 30,
                                "periodSeconds": 10,
                            },
                            "resources": {
                                "requests": {"memory": "128Mi", "cpu": "100m"},
                                "limits": {"memory": "512Mi", "cpu": "300m"},
                            },
                        }
                    ],
                    "volumes": [{"name": "warp-data", "emptyDir": {}}],
                },
            },
        },
    }


def generate_service(gateway_config: dict[str, Any]) -> dict[str, Any]:
    """Generate service for VPN gateway"""
    name: str = gateway_config["name"]
    gateway_type: str = gateway_config["type"]

    if gateway_type == "warp":
        ports: list[dict[str, Any]] = [
            {"name": "socks5-proxy", "port": 1080, "targetPort": 1080},
            {"name": "http-proxy", "port": 3128, "targetPort": 3128},
        ]
    else:
        ports = [
            {"name": "http-proxy", "port": 8888, "targetPort": 8888},
            {"name": "shadowsocks", "port": 8388, "targetPort": 8388},
            {"name": "health", "port": 9999, "targetPort": 9999},
        ]

    return {
        "apiVersion": "v1",
        "kind": "Service",
        "metadata": {
            "name": name,
            "namespace": "vpn-gateway",
            "labels": {"app": name, "component": "vpn-gateway"},
        },
        "spec": {"selector": {"app": name}, "ports": ports},
    }


def main():
    """Generate all VPN gateway deployments"""
    output_file = Path.cwd().joinpath("vpn-gateways-all.yaml")

    manifests: list[dict[str, Any]] = []

    # Generate deployments and services for each gateway
    for gateway_config in VPN_GATEWAYS:
        gateway_type: str = gateway_config["type"]

        if gateway_type == "premiumize":
            deployment: dict[str, Any] = generate_premiumize_deployment(gateway_config)
        elif gateway_type == "airvpn":
            deployment = generate_airvpn_deployment(gateway_config)
        elif gateway_type == "warp":
            deployment = generate_warp_deployment(gateway_config)
        else:
            continue

        service: dict[str, Any] = generate_service(gateway_config)

        manifests.extend([deployment, service])

    # Write all manifests to file
    with open(output_file, "w") as f:
        f.write("# Generated VPN Gateway Deployments for Pod-Gateway System\n")
        f.write("# This file contains all VPN gateway deployments and services\n\n")

        for i, manifest in enumerate(manifests):
            if i > 0:
                f.write("\n---\n")
            yaml.dump(manifest, f, default_flow_style=False, sort_keys=False)

    print(f"Generated {len(manifests)} manifests in {output_file}")
    print(f"Generated {len(VPN_GATEWAYS)} VPN gateways:")
    for gateway in VPN_GATEWAYS:
        print(
            f"  - {gateway['name']} (priority {gateway['priority']}, {gateway['type']}, {gateway['country']})"
        )


if __name__ == "__main__":
    main()
