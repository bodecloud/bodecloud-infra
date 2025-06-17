#!/bin/bash

declare -A NS_MAP=(
  ["bolabaden"]="bolabaden"
  ["my-media-stack"]="bolabaden"
  ["cert-manager"]="cert-manager"
  ["default"]="default"
  ["kube-node-lease"]="kube-node-lease"
  ["kube-public"]="kube-public"
  ["kube-system"]="kube-system"
  ["portainer"]="portainer"
  ["vpn-gateway"]="vpn-gateway"
)

for old_ns in "${!NS_MAP[@]}"; do
  new_ns="${NS_MAP[$old_ns]}"
  if [ $old_ns == $new_ns ]; then
    echo "Skipping $old_ns -> $new_ns (same namespace)"
    continue
  fi

  echo "Migrating from $old_ns to $new_ns..."

  # Get all resource types you care about
  for kind in deployments services configmaps secrets statefulsets daemonsets jobs cronjobs pvc; do
    resources=$(kubectl get $kind -n "$old_ns" --no-headers -o custom-columns=":metadata.name" 2>/dev/null)
    if [ -z "$resources" ]; then
      echo "No $kind resources found in $old_ns"
      continue
    fi
    echo "Found $resources $kind resources in $old_ns"
    for res in $resources; do
      kubectl get $kind "$res" -n "$old_ns" -o yaml |
        yq eval "del(.metadata.uid, .metadata.resourceVersion, .metadata.creationTimestamp, .metadata.selfLink, .status)" - |
        yq eval ".metadata.namespace = \"$new_ns\"" - |
        kubectl apply -f -

      echo "✓ Migrated $kind/$res to $new_ns"
    done
  done
done
