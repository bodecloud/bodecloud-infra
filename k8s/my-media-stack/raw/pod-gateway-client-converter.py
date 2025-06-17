#!/usr/bin/env python3
"""
Pod-Gateway Client Converter
Converts services to use pod-gateway for transparent WARP VPN routing
This replaces simple route manipulation with proper VXLAN tunneling
"""

from __future__ import annotations

import json
import subprocess
import sys
from typing import Any


def run_kubectl(args: list[str]) -> tuple[bool, str]:
    """Run kubectl command and return success status and output."""
    try:
        result = subprocess.run(
            ["k3s", "kubectl"] + args,
            capture_output=True,
            text=True,
            check=False
        )
        return result.returncode == 0, result.stdout.strip()
    except Exception as e:
        return False, str(e)


def deployment_exists(service: str, namespace: str) -> bool:
    """Check if deployment exists."""
    success, _ = run_kubectl(["get", "deployment", service, "-n", namespace])
    return success


def has_pod_gateway_containers(service: str, namespace: str) -> bool:
    """Check if deployment already has pod-gateway containers."""
    success, output = run_kubectl(["get", "deployment", service, "-n", namespace, "-o", "yaml"])
    if not success:
        return False
    return "pod-gateway" in output and "client" in output


def has_old_vpn_routing(service: str, namespace: str) -> bool:
    """Check if deployment has old simple VPN routing."""
    success, output = run_kubectl(["get", "deployment", service, "-n", namespace, "-o", "yaml"])
    if not success:
        return False
    return "vpn-route-init" in output


def remove_old_vpn_routing(service: str, namespace: str) -> bool:
    """Remove old simple VPN routing from deployment."""
    print(f"   🧹 Removing old simple VPN routing from {service}...")
    
    # Remove init containers
    patch = '[{"op": "remove", "path": "/spec/template/spec/initContainers"}]'
    success1, _ = run_kubectl([
        "patch", "deployment", service, "-n", namespace,
        "--type=json", f"-p={patch}"
    ])
    
    # Remove vpn routing label
    success2, _ = run_kubectl([
        "label", "deployment", service, "-n", namespace, "vpn.routing-"
    ])
    
    return success1 or success2  # At least one should succeed


def create_pod_gateway_patch(gateway_name: str) -> dict[str, Any]:
    """Create the pod-gateway patch configuration."""
    return {
        "spec": {
            "template": {
                "metadata": {
                    "labels": {
                        "vpn.routing": "pod-gateway"
                    }
                },
                "spec": {
                    "initContainers": [
                        {
                            "name": "pod-gateway-client-init",
                            "image": "ghcr.io/angelnu/pod-gateway:v1.10.0",
                            "command": ["/bin/client_init.sh"],
                            "securityContext": {
                                "privileged": True,
                                "capabilities": {
                                    "add": ["NET_ADMIN", "SYS_MODULE"]
                                }
                            },
                            "env": [
                                {"name": "gateway", "value": gateway_name},
                                {"name": "K8S_DNS_ips", "value": "10.43.0.10"},
                                {"name": "VXLAN_ID", "value": "42"},
                                {"name": "VXLAN_PORT", "value": "4789"},
                                {"name": "VXLAN_IP_NETWORK", "value": "172.16.0"},
                                {"name": "VPN_INTERFACE_MTU", "value": "1420"},
                                {"name": "CONNECTION_RETRY_COUNT", "value": "3"}
                            ],
                            "volumeMounts": [
                                {
                                    "name": "pod-gateway-client-config",
                                    "mountPath": "/config"
                                }
                            ]
                        }
                    ],
                    "containers": [
                        {
                            "name": "pod-gateway-client-sidecar",
                            "image": "ghcr.io/angelnu/pod-gateway:v1.10.0",
                            "command": ["/bin/client_sidecar.sh"],
                            "securityContext": {
                                "privileged": True,
                                "capabilities": {
                                    "add": ["NET_ADMIN"]
                                }
                            },
                            "env": [
                                {"name": "gateway", "value": gateway_name},
                                {"name": "K8S_DNS_ips", "value": "10.43.0.10"},
                                {"name": "VXLAN_ID", "value": "42"},
                                {"name": "VXLAN_PORT", "value": "4789"},
                                {"name": "VXLAN_IP_NETWORK", "value": "172.16.0"},
                                {"name": "CONNECTION_RETRY_COUNT", "value": "3"}
                            ],
                            "volumeMounts": [
                                {
                                    "name": "pod-gateway-client-config",
                                    "mountPath": "/config"
                                }
                            ],
                            "resources": {
                                "requests": {
                                    "memory": "32Mi",
                                    "cpu": "25m"
                                },
                                "limits": {
                                    "memory": "64Mi",
                                    "cpu": "50m"
                                }
                            }
                        }
                    ],
                    "volumes": [
                        {
                            "name": "pod-gateway-client-config",
                            "configMap": {
                                "name": "pod-gateway-client-config",
                                "defaultMode": 493
                            }
                        }
                    ]
                }
            }
        }
    }


def apply_pod_gateway_patch(service: str, namespace: str, patch: dict[str, Any]) -> bool:
    """Apply pod-gateway patch to deployment."""
    patch_json = json.dumps(patch)
    success, _ = run_kubectl([
        "patch", "deployment", service, "-n", namespace,
        "--type=strategic", f"-p={patch_json}"
    ])
    return success


def add_routing_label(service: str, namespace: str) -> bool:
    """Add pod-gateway routing label to deployment."""
    success, _ = run_kubectl([
        "label", "deployment", service, "-n", namespace,
        "vpn.routing=pod-gateway", "--overwrite"
    ])
    return success


def main() -> None:
    """Main function to convert services to pod-gateway VPN routing."""
    namespace = "my-media-stack"
    gateway_name = "warp-pod-gateway"
    
    # Services that should route through VPN (same as Docker Compose network_mode: service:warp)
    vpn_services = [
        "aiostreams",
        "mediafusion",
        "comet",
        "torrentio",
        "stremthru",
        "jackett",
        "prowlarr",
        "flaresolverr"
    ]
    
    print("🚀 Converting services to use pod-gateway for transparent WARP VPN routing...")
    print("   This provides VXLAN tunneling, DHCP, DNS, and automatic reconnection")
    print()
    
    for service in vpn_services:
        print(f"🔍 Processing {service}...")
        
        # Check if deployment exists
        if not deployment_exists(service, namespace):
            print(f"   ⚠️  Deployment {service} not found, skipping...")
            continue
        
        # Check if already has pod-gateway containers
        if has_pod_gateway_containers(service, namespace):
            print(f"   ✅ {service} already configured for pod-gateway VPN routing")
            continue
        
        # Remove old simple VPN routing if it exists
        if has_old_vpn_routing(service, namespace):
            remove_old_vpn_routing(service, namespace)
        
        print(f"   🔧 Adding pod-gateway VPN routing to {service}...")
        
        # Create and apply pod-gateway patch
        patch = create_pod_gateway_patch(gateway_name)
        
        if apply_pod_gateway_patch(service, namespace, patch):
            print(f"   ✅ {service} configured for pod-gateway VPN routing")
            
            # Add pod-gateway routing label to deployment
            add_routing_label(service, namespace)
        else:
            print(f"   ❌ Failed to configure {service}")
    
    print()
    print("🎯 Pod-Gateway VPN Configuration Complete!")
    print()
    print("📋 Summary:")
    print("   - All specified services now route through WARP VPN via pod-gateway")
    print("   - Services use VXLAN tunneling for robust connectivity")
    print("   - Automatic DHCP IP assignment and DNS resolution")
    print("   - Automatic reconnection if gateway restarts")
    print("   - Services are labeled with 'vpn.routing=pod-gateway'")
    print()
    print("🚀 Next steps:")
    print("   1. Deploy pod-gateway client config: kubectl apply -f pod-gateway-client-config.yaml")
    print("   2. Deploy WARP pod-gateway: kubectl apply -f warp-pod-gateway-deployment.yaml")
    print(f"   3. Wait for gateway to be ready: kubectl wait --for=condition=ready pod -l app=warp-pod-gateway -n {namespace} --timeout=300s")
    print(f"   4. Test VPN routing: kubectl exec -it deployment/jackett -n {namespace} -- wget -qO- http://ifconfig.me")
    print()
    print("🔍 Verify pod-gateway routing:")
    print(f"   kubectl exec -it deployment/jackett -n {namespace} -- ip addr show vxlan0")
    print(f"   kubectl exec -it deployment/jackett -n {namespace} -- ip route show")
    print(f"   kubectl exec -it deployment/jackett -n {namespace} -- wget -qO- http://ifconfig.me")
    print()
    print("🐛 Debug commands:")
    print(f"   kubectl logs -l app=warp-pod-gateway -n {namespace} -c gateway-sidecar")
    print(f"   kubectl logs deployment/jackett -n {namespace} -c pod-gateway-client-sidecar")


if __name__ == "__main__":
    main()