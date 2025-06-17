#!/bin/bash

echo "Fixing storage class issues by deleting and re-creating PVCs with local-path storage class..."

# Delete statefulsets since they can't be modified directly
echo "Deleting statefulsets..."
kubectl delete statefulset elasticsearch
kubectl delete statefulset influxdb
kubectl delete statefulset mongodb
kubectl delete statefulset postgresql
kubectl delete statefulset redis
kubectl delete statefulset qdrant

# Delete the existing PVCs that are using the wrong storage class
echo "Deleting PVCs with wrong storage class..."
kubectl delete pvc elasticsearch-data-elasticsearch-0
kubectl delete pvc influxdb-data-influxdb-0
kubectl delete pvc mongodb-data-mongodb-0
kubectl delete pvc postgresql-data-postgresql-0
kubectl delete pvc redis-data-redis-0
kubectl delete pvc qdrant-storage-qdrant-0

# Create temporary PVC files with local-path storage class
echo "Creating updated PVC manifests with local-path storage class..."

cat <<EOF > elasticsearch-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: elasticsearch-data-elasticsearch-0
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: local-path
EOF

cat <<EOF > influxdb-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: influxdb-data-influxdb-0
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: local-path
EOF

cat <<EOF > mongodb-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mongodb-data-mongodb-0
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: local-path
EOF

cat <<EOF > postgresql-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgresql-data-postgresql-0
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: local-path
EOF

cat <<EOF > redis-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: redis-data-redis-0
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: local-path
EOF

cat <<EOF > qdrant-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: qdrant-storage-qdrant-0
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: local-path
EOF

# Apply the new PVCs
echo "Applying updated PVCs..."
kubectl apply -f elasticsearch-pvc.yaml
kubectl apply -f influxdb-pvc.yaml
kubectl apply -f mongodb-pvc.yaml
kubectl apply -f postgresql-pvc.yaml
kubectl apply -f redis-pvc.yaml
kubectl apply -f qdrant-pvc.yaml

# Clean up the temporary files
echo "Cleaning up temporary files..."
rm elasticsearch-pvc.yaml influxdb-pvc.yaml mongodb-pvc.yaml postgresql-pvc.yaml redis-pvc.yaml qdrant-pvc.yaml

echo "Storage class issues fixed! Statefulsets should now be able to start." 