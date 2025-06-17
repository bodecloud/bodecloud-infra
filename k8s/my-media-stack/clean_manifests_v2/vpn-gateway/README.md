# Enhanced VPN Gateway for Kubernetes

This directory contains the configuration for an enhanced VPN gateway based on the [pod-gateway](https://github.com/angelnu/pod-gateway) chart with additional custom components for improved reliability and failover.

## Overview

This implementation preserves the core functionality of the original k8s-gitops vpn-gateway while adding improvements for reliability, monitoring, and automatic failover. It uses the Gluetun VPN client container to establish the VPN connection and route traffic through a VXLAN tunnel for improved isolation and security.

## Components

1. **VPN Gateway (pod-gateway)**: The main gateway that routes traffic through the VPN using the Helm chart (preserved from original implementation).
2. **VPN Failover Controller**: Automatically switches between VPN providers if one fails, with real-time status monitoring.
3. **Pod Gateway Network**: A DaemonSet that sets up the VXLAN tunnel and routing on each node.
4. **Network Policy**: Ensures that pods in the VPN namespace can only access the internet through the VPN.
5. **Status Monitoring**: Tracks VPN connection status and provides real-time information.

## VPN Providers

The following VPN providers are configured in order of priority:

1. Premiumize US (default)
2. Premiumize NL (fallback)
3. Premiumize DE (fallback)

## Architecture

This implementation combines several approaches:

1. **Original Pod Gateway Helm Chart**: Preserved from the k8s-gitops implementation.
2. **Custom Failover Controller**: Added to monitor VPN health and automatically switch between providers.
3. **Network DaemonSet**: Added to set up the necessary networking on each node for VXLAN tunneling.
4. **Status Tracking**: Added to maintain real-time status information in a ConfigMap.

## Usage

To route a namespace through the VPN, add the `vpn-routed: "true"` label to the namespace.

Example:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
  labels:
    vpn-routed: "true"
```

## Configuration

- VPN credentials are stored in the `vpn-credentials` secret.
- VPN configurations are stored in the `vpn-configs` ConfigMap.
- The failover controller automatically switches between VPN providers if one fails.
- VPN status is tracked in the `vpn-status` ConfigMap.

## Troubleshooting

If you encounter issues with the VPN connection, check the logs of the VPN gateway pod:

```bash
kubectl logs -n vpn-gateway -l app=vpn-gateway
```

To check the failover controller logs:

```bash
kubectl logs -n vpn-gateway -l app=vpn-failover-controller
```

To check the network setup on nodes:

```bash
kubectl logs -n vpn-gateway -l app=pod-gateway-network
```

To check the current VPN status:

```bash
kubectl get configmap vpn-status -n vpn-gateway -o yaml
``` 