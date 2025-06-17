# Deployment Resource Guidelines

For a single-node cluster with 4 CPU cores and 24GB memory:

1. Critical media services (Plex, Jellyfin):
   - CPU: 200m-500m requests, 1000m-2000m limits
   - Memory: 256Mi-512Mi requests, 1Gi-2Gi limits

2. *arr services (Sonarr, Radarr, etc):
   - CPU: 10m-50m requests, 200m limits
   - Memory: 64Mi requests, 256Mi limits

3. Helper services (homer, wizarr, etc):
   - CPU: 10m requests, 100m limits
   - Memory: 64Mi requests, 128Mi limits

4. Avoid running duplicate services (e.g., don't run both Jellyfin and Plex if one will suffice)

5. Set reasonable replica counts - typically 1 for each service
