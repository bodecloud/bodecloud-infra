# Kubernetes Cluster Export

This directory contains scripts to export all resources from a Kubernetes cluster to YAML manifests. The exported manifests can be used for backup, migration, or GitOps purposes.

## Available Export Scripts

Three different export methods are provided:

1. **export-manifests.sh**: Custom script that exports all resources, Helm releases, and creates a Kustomize structure.
2. **export-with-kubedump.sh**: Uses the [kube-dump](https://github.com/WoozyMasta/kube-dump) tool for exporting resources.
3. **export-with-manifest-exporter.sh**: Uses an approach similar to [kubernetes-manifest-exporter](https://github.com/jonchen727/kubernetes-manifest-exporter) with kubectl-neat for clean output.

## Usage

Choose one of the export methods and run the script from this directory:

```bash
# Method 1: Custom export script
cd ./k8s/my-media-stack/clean_manifests_v2
./export-manifests.sh

# Method 2: Using kube-dump
cd ./k8s/my-media-stack/clean_manifests_v2
./export-with-kubedump.sh

# Method 3: Using manifest-exporter approach
cd ./k8s/my-media-stack/clean_manifests_v2
./export-with-manifest-exporter.sh
```

## Output Structure

Each script creates a slightly different output structure:

### export-manifests.sh

```shell
./
├── namespaced/           # Namespaced resources organized by namespace and resource type
├── cluster/              # Cluster-wide resources organized by resource type
├── kubeconfig/           # Sanitized kubeconfig
├── helm-releases/        # Helm releases with manifests and values
└── kustomize/            # Kustomize structure referencing all exported resources
```

### export-with-kubedump.sh

```shell
./kube-dump/              # All resources exported by kube-dump
```

### export-with-manifest-exporter.sh

```shell
./manifest-exporter/
├── namespaces/           # Namespaced resources organized by namespace and resource type
└── cluster/              # Cluster-wide resources organized by resource type
```

## Requirements

- kubectl configured with access to your cluster
- helm (optional, for exporting Helm releases)
- kubectl-neat (optional, installed automatically by export-with-manifest-exporter.sh)

## Notes

- Some resources like events, nodes, and componentstatuses are skipped by default to reduce clutter.
- The exported manifests have cluster-specific metadata removed (UIDs, resourceVersions, etc.).
- For large clusters, the export process may take several minutes to complete.
