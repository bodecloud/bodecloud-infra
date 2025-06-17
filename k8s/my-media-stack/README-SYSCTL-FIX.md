# Fixing SysctlForbidden Errors in Kubernetes

This document explains how to fix the `SysctlForbidden` error that occurs when pods try to use sysctls that are not allowed by default in Kubernetes.

## Problem

The `SysctlForbidden` error occurs when a pod tries to use certain system controls (sysctls) that are not allowed by default in Kubernetes for security reasons. In our case, the warp-gateway pods are trying to set sysctls like `net.ipv4.ip_forward` that are forbidden.

## Solution

The solution is to configure the kubelet to allow these unsafe sysctls. We've created a script that automates this process:

1. Run the fix-sysctls.sh script:

   ```bash
   sudo ./k8s/my-media-stack/fix-sysctls.sh
   ```

2. This script will:
   - Create an override configuration for the k3s service
   - Configure the kubelet to allow the required sysctls
   - Restart the k3s service
   - Apply the fixed warp-gateway deployment
   - Clean up any failed pods

## Manual Fix

If you prefer to fix the issue manually, follow these steps:

1. Create or modify the k3s service override configuration:

   ```bash
   sudo mkdir -p /etc/systemd/system/k3s.service.d/
   ```

2. Create the override.conf file:

   ```bash
   cat << EOF | sudo tee /etc/systemd/system/k3s.service.d/override.conf
   [Service]
   ExecStart=
   ExecStart=/usr/local/bin/k3s server --kubelet-arg=allowed-unsafe-sysctls=net.ipv4.ip_forward,net.ipv4.conf.all.forwarding,net.ipv6.conf.all.forwarding,net.ipv6.conf.all.disable_ipv6,net.ipv4.conf.all.src_valid_mark
   EOF
   ```

3. Reload systemd configuration:

   ```bash
   sudo systemctl daemon-reload
   ```

4. Restart k3s service:

   ```bash
   sudo systemctl restart k3s
   ```

5. Apply the fixed warp-gateway deployment:

   ```bash
   kubectl delete deployment warp-gateway -n vpn-gateway --ignore-not-found
   kubectl apply -f ./k8s/my-media-stack/raw_v2/fixed-warp-gateway-deployment.yaml
   ```

## Required Sysctls

The following sysctls are required for the warp-gateway to function properly:

- `net.ipv4.ip_forward=1`
- `net.ipv4.conf.all.forwarding=1`
- `net.ipv6.conf.all.forwarding=1`
- `net.ipv6.conf.all.disable_ipv6=0`
- `net.ipv4.conf.all.src_valid_mark=1`

## Verification

To verify that the fix has been applied successfully, check the status of the warp-gateway pods:

```bash
kubectl get pods -n vpn-gateway
```

The pods should now be in the Running state, rather than showing the SysctlForbidden error.
