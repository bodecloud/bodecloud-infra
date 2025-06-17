#!/usr/bin/env python3
"""
Auto-generate IngressRoute resources for all services in my-media-stack namespace
This script provides a semi-dynamic approach to expose services
"""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path
from typing import List, Tuple


def run_kubectl_command(command: list[str]) -> str:
    """Run kubectl command and return output"""
    try:
        result: subprocess.CompletedProcess[str] = subprocess.run(command, capture_output=True, text=True, check=True)
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"Error running kubectl command: {e}")
        sys.exit(1)


def get_services(namespace: str) -> List[Tuple[str, str, str, str, str, str]]:
    """Get all services in the namespace"""
    command: list[str] = ["k3s", "kubectl", "get", "services", "-n", namespace, "--no-headers"]
    output: str = run_kubectl_command(command)

    services: list[tuple[str, str, str, str, str, str]] = []
    for line in output.split("\n"):
        if line.strip():
            parts: list[str] = line.split()
            if len(parts) >= 6:
                name: str = parts[0]
                service_type: str = parts[1]
                cluster_ip: str = parts[2]
                external_ip: str = parts[3]
                ports: str = parts[4]
                age: str = parts[5]
                services.append(
                    (name, service_type, cluster_ip, external_ip, ports, age)
                )

    return services


def extract_port(ports_str: str) -> str:
    """Extract the first port from the ports string"""
    if ports_str == "<none>":
        return ""

    # Extract the first port (assuming it's the main port)
    port: str = ports_str.split("/")[0].split(":")[0]
    return port


def generate_ingressroute(
    name: str,
    port: str,
    namespace: str,
    domain: str,
) -> str:
    """Generate IngressRoute YAML for a service"""
    return f"""---
# Auto-generated IngressRoute for {name}
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: {name}-auto-ingressroute
  namespace: {namespace}
  labels:
    auto-generated: "true"
spec:
  entryPoints:
    - websecure
  routes:
    - match: "Host(`{name}.{domain}`)"
      kind: Rule
      services:
        - name: {name}
          port: {port}
  tls:
    certResolver: beatapostapita_duckdns_letsencrypt

"""


def main():
    """Main function to generate IngressRoutes"""
    namespace: str = "my-media-stack"
    domain: str = "beatapostapita.duckdns.org"
    output_file: str = "auto-generated-ingressroutes.yaml"

    print(f"🔍 Discovering services in namespace: {namespace}")

    # Get all services in the namespace
    services: list[tuple[str, str, str, str, str, str]] = get_services(namespace)

    # Clear the output file
    output_path = Path(output_file)
    output_path.write_text("")

    valid_services: list[str] = []

    # Generate IngressRoutes for each service
    with open(output_file, "w") as f:
        for name, service_type, cluster_ip, external_ip, ports, age in services:
            # Skip headless services and services without ports
            if cluster_ip == "None" or ports == "<none>":
                print(f"   Skipping headless/portless service: {name}")
                continue

            port: str = extract_port(ports)
            if not port:
                print(f"   Skipping service without valid port: {name}")
                continue

            print(f"   Creating IngressRoute for: {name} (port {port})")

            ingressroute_yaml: str = generate_ingressroute(name, port, namespace, domain)
            f.write(ingressroute_yaml)

            valid_services.append(name)

    print(f"✅ Generated IngressRoutes in: {output_file}")
    print("")
    print("🚀 To apply the auto-generated routes, run:")
    print(f"   k3s kubectl apply -f {output_file}")
    print("")
    print("🌍 Your services will be accessible at:")
    for service_name in valid_services:
        print(f"   - https://{service_name}.{domain}")


if __name__ == "__main__":
    main()
