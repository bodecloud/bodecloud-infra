# Storage Class Fix Plan

All PVCs should use the 'local-path' storage class which is what's available in the cluster.

1. In myprecious-pvcs.ts:
   - Replace all instances of 'ceph-block-ssd' with 'local-path'
   - Set reasonable storage sizes (10Gi is sufficient for most services)
   - Add volumeMode: Filesystem if not already present

2. In myprecious-chart.ts:
   - Ensure all StatefulSet volumeClaimTemplates use storageClassName: 'local-path'

3. For stateful applications like:
   - Elasticsearch
   - MongoDB
   - PostgreSQL
   - Redis
   
   Make sure they have proper init container configurations if needed, and 
   their storage requirements are reasonable for your node's capacity.
