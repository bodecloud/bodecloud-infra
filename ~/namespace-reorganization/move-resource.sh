#!/bin/bash

# Usage: ./move-resource.sh <resource_type> <resource_name> <source_namespace> <target_namespace>
# Example: ./move-resource.sh deployment nginx-auth my-media-stack infrastructure

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <resource_type> <resource_name> <source_namespace> <target_namespace>"
    exit 1
fi

RESOURCE_TYPE=$1
RESOURCE_NAME=$2
SOURCE_NS=$3
TARGET_NS=$4
TEMP_FILE="$HOME/namespace-reorganization/${RESOURCE_TYPE}-${RESOURCE_NAME}.yaml"

# Check if the target namespace exists, create it if not
kubectl get namespace $TARGET_NS &>/dev/null
if [ $? -ne 0 ]; then
    echo "Creating namespace $TARGET_NS..."
    kubectl create namespace $TARGET_NS
fi

# Export the resource definition
echo "Exporting $RESOURCE_TYPE/$RESOURCE_NAME from $SOURCE_NS namespace..."
kubectl get $RESOURCE_TYPE $RESOURCE_NAME -n $SOURCE_NS -o yaml >$TEMP_FILE

# Modify the resource definition for the new namespace
echo "Modifying resource definition for $TARGET_NS namespace..."
sed -i "s/namespace: $SOURCE_NS/namespace: $TARGET_NS/g" $TEMP_FILE

# Remove cluster-specific fields
sed -i '/resourceVersion:/d' $TEMP_FILE
sed -i '/uid:/d' $TEMP_FILE
sed -i '/creationTimestamp:/d' $TEMP_FILE
sed -i '/generation:/d' $TEMP_FILE
sed -i '/status:/,$d' $TEMP_FILE
sed -i '/selfLink:/d' $TEMP_FILE

# Apply the modified resource to the target namespace
echo "Creating $RESOURCE_TYPE/$RESOURCE_NAME in $TARGET_NS namespace..."
kubectl apply -f $TEMP_FILE

# If successful, delete the resource from the source namespace
if [ $? -eq 0 ]; then
    echo "Deleting $RESOURCE_TYPE/$RESOURCE_NAME from $SOURCE_NS namespace..."
    kubectl delete $RESOURCE_TYPE $RESOURCE_NAME -n $SOURCE_NS
    echo "Successfully moved $RESOURCE_TYPE/$RESOURCE_NAME from $SOURCE_NS to $TARGET_NS namespace!"
else
    echo "Failed to create resource in target namespace. Resource not deleted from source namespace."
    exit 1
fi
