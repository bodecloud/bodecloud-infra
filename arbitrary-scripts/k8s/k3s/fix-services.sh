#!/bin/bash

# Fix all service files by adding missing selectors
for file in k8s/my-media-stack/*-expose.yaml; do
    if [ -f "$file" ]; then
        echo "Fixing $file..."
        
        # Check if selector already exists
        if ! grep -q "selector:" "$file"; then
            # Extract service name from the file
            service_name=$(basename "$file" | sed 's/-expose\.yaml$//')
            
            # Add selector after the spec: line
            sed -i '/^spec:/a\    selector:\n        com.docker.compose.service: '"$service_name" "$file"
            
            echo "Added selector for $service_name"
        else
            echo "Selector already exists in $file"
        fi
    fi
done

echo "Applying all fixed services..."
kubectl apply -f k8s/my-media-stack/*-expose.yaml 