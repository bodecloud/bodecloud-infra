# Garden.io Deployment Status

## Reality check

This document is useful as a deployment-progress snapshot for the `garden.io/`
exploration branch.
It is not authoritative current repo status, and it does not override the
Compose-first root runtime as the priority implementation.

For the current evidence-led interpretation of this branch, see:

- [`../knowledgebase/research/garden-k3s-exploration-evidence.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/garden-k3s-exploration-evidence.md)

## Current Deployment Phase: Docker Compose Testing

**Status:** historical branch snapshot

### Deployment Progress

#### ✅ Phase 1: Core Infrastructure (COMPLETE)
- ✅ dockerproxy-ro - Healthy
- ✅ redis - Healthy  
- ✅ mongodb - Healthy

#### ✅ Phase 2: Reverse Proxy (COMPLETE)
- ✅ crowdsec - Healthy
- ✅ nginx-traefik-extensions - Healthy
- ✅ searxng - Healthy
- ✅ tinyauth - Healthy
- ✅ traefik - Healthy

#### ✅ Phase 3: Infrastructure Services (IN PROGRESS)
- ✅ homepage - Healthy
- ⏳ dockerproxy-rw - Starting
- ⏳ dozzle - Starting
- ⏳ portainer - Restarting
- ✅ watchtower - Running

#### ⏳ Phase 4: Application Services (PENDING)
- ⏳ bolabaden-nextjs - Pending
- ⏳ session-manager - Pending
- ⏳ telemetry-auth - Fixed secret, deploying

### Health Summary

**Current Status:**
- Healthy: 9 services
- Starting: 2-3 services
- Unhealthy: 0 services
- Health Percentage: ~75%

### Next Steps

1. ✅ Stop Nomad services - COMPLETE
2. ✅ Deploy core services - COMPLETE
3. ⏳ Deploy remaining services - IN PROGRESS
4. ⏳ Verify all services healthy - PENDING
5. ⏳ Deploy to Kubernetes - PENDING

### Service Dependencies

Services are being deployed in dependency order:
1. Core (dockerproxy-ro, redis, mongodb)
2. Reverse Proxy (traefik, crowdsec, nginx, tinyauth)
3. Infrastructure (homepage, dozzle, portainer, watchtower)
4. Applications (bolabaden-nextjs, session-manager, telemetry-auth)
5. Additional services (Firecrawl, LLM, Stremio, Metrics, WARP)

### Notes

- All services have comprehensive healthchecks
- Services are verified healthy before proceeding
- Kubernetes deployment will only occur after full Docker Compose validation
