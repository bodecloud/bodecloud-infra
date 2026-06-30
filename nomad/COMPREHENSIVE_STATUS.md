# Comprehensive Nomad Cluster Status Report
Generated: $(date)

## Reality check

This file is valuable because it preserves concrete Nomad and Consul cluster
observations.
It should not be read as proof that the Nomad branch has already earned full
promotion as the repo's settled answer.

The same document still records critical SPOFs and partial-node availability.

For the assimilated reading of the Nomad branch, see:

- [`../knowledgebase/research/nomad-exploration-evidence.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/nomad-exploration-evidence.md)

## Executive Summary

### ✅ Operational Status
- **Nomad Servers**: 2 active (healthy quorum) ✅
- **Nomad Clients**: 2 ready, 2 down ⚠️
- **Consul Servers**: 1 active ⚠️ **CRITICAL SINGLE POINT OF FAILURE**
- **Consul Services**: 31 registered ✅
- **Running Services**: Most critical services operational ✅

### 🎯 Requirements Status

#### 1. 1:1 Parity with Docker Compose
- ⚠️ **PARTIAL**: Most services match, but some discrepancies found:
  - firecrawl-group: Resources fixed (was reduced, now matches)
  - nuq-postgres: Uses local image instead of build (Nomad limitation)
  - Some environment variables may need verification

#### 2. Fully Healthy
- ✅ **MOSTLY**: Most services running
- ⚠️ **ISSUES**:
  - dozzle-group: Failing (55 failed attempts)
  - firecrawl-group: Queued (waiting for dependencies)
  - playwright-service-group: Starting
  - nuq-postgres-group: Starting

#### 3. Complete Fallback/Failover
- ⚠️ **PARTIAL**: 
  - 7 services have HA (count > 1) ✅
  - Many services have count=1 (single points of failure) ⚠️
  - Spread constraints configured for HA services ✅

#### 4. Zero Single Points of Failure
- ❌ **NOT ACHIEVED**: Multiple single points of failure identified:
  - Consul: Only 1 server (CRITICAL)
  - Many services: count=1
  - 2 nodes down (reduces failover capacity)

#### 5. All Nodes Functional
- ⚠️ **PARTIAL**: 
  - 2 nodes ready (micklethefickle, cloudserver1) ✅
  - 2 nodes down (beatapostapita, cloudserver2) ❌

## Detailed Service Status

### ✅ Running Services (Critical)
- mongodb-group: 1 running
- redis-group: 1 running
- searxng-group: 2 running (HA) ✅
- homepage-group: 2 running (HA) ✅
- bolabaden-nextjs-group: 2 running (HA) ✅
- stremio-group: 1 running (HA configured)
- traefik-group: 1 running (HA configured, count=3)
- aiostreams-group: 1 running (HA configured)
- litellm-group: 2 running (HA) ✅
- litellm-postgres-group: 1 running
- infrastructure-services: 1 running (dockerproxy-ro available)

### ⚠️ Services with Issues
- **dozzle-group**: 0 running, 55 failed
  - Latest allocation (28cfe520) failed
  - Needs investigation
  
- **firecrawl-group**: 1 queued, 0 running
  - Waiting for dependencies (playwright-service, nuq-postgres)
  - Constrained to micklethefickle node
  
- **playwright-service-group**: 1 starting, 0 running
  - Constrained to micklethefickle node
  - Should start soon
  
- **nuq-postgres-group**: 1 starting, 0 running
  - Constrained to micklethefickle node
  - Should start soon

### Services with HA (count > 1)
1. ✅ searxng-group: count=2, spread configured
2. ✅ homepage-group: count=2, spread configured
3. ✅ bolabaden-nextjs-group: count=2, spread configured
4. ✅ aiostreams-group: count=2, spread configured
5. ✅ stremio-group: count=2, spread configured
6. ✅ traefik-group: count=3, spread configured
7. ✅ litellm-group: count=2, spread configured

### Services WITHOUT HA (Single Points of Failure)
- mongodb-group (OK - MongoDB handles replication)
- redis-group ⚠️
- nuq-postgres-group ⚠️
- litellm-postgres-group ⚠️
- playwright-service-group ⚠️
- firecrawl-group ⚠️
- dozzle-group ⚠️
- And 20+ other services

## Critical Issues

### 1. Consul Single Point of Failure - CRITICAL
- **Current**: 1 server (micklethefickle)
- **Required**: 3+ servers for HA
- **Impact**: If Consul fails, all service discovery fails
- **Priority**: HIGH
- **Action**: Deploy 2+ additional Consul servers

### 2. Node Availability
- **Down Nodes**: beatapostapita, cloudserver2
- **Reason**: Node heartbeat missed
- **Impact**: Reduced capacity, limited failover
- **Priority**: MEDIUM
- **Action**: Investigate and restore or remove from cluster

### 3. Dozzle Service Failing
- **Status**: 55 failed attempts, latest allocation failed
- **Impact**: Log viewer unavailable
- **Priority**: LOW (non-critical service)
- **Action**: Investigate failure cause

### 4. Firecrawl Dependencies
- **Status**: firecrawl queued, waiting for playwright-service and nuq-postgres
- **Impact**: Firecrawl service unavailable
- **Priority**: MEDIUM
- **Action**: Monitor - services are starting

## 1:1 Parity Issues Found

### Fixed
- ✅ firecrawl-group resources: Updated to match docker-compose (4000 CPU, 4096 memory)

### Remaining
- ⚠️ nuq-postgres: Uses `my-media-stack-nuq-postgres:local` instead of build
  - **Reason**: Nomad doesn't support docker-compose build syntax
  - **Solution**: Pre-build image or use registry image
  - **Status**: Acceptable workaround if image exists

### Needs Verification
- Environment variables: Compare all env vars with docker-compose
- Volume mounts: Verify all volumes match
- Health checks: Verify all health checks match
- Network configuration: Verify network settings

## Recommendations

### Immediate (Priority: HIGH)
1. **Deploy Additional Consul Servers**
   ```bash
   # Deploy Consul on cloudserver1 and another node
   # Ensure 3+ servers for HA
   ```

2. **Monitor Firecrawl Dependencies**
   - Wait for playwright-service and nuq-postgres to start
   - Verify firecrawl starts once dependencies are ready

3. **Investigate Dozzle Failures**
   - Check allocation logs for root cause
   - Verify dockerproxy-ro service is accessible
   - Check resource constraints

### Short-term (Priority: MEDIUM)
1. **Restore Down Nodes**
   - Investigate why beatapostapita and cloudserver2 are down
   - Check network connectivity
   - Restore or remove from cluster

2. **Increase HA for Critical Services**
   - Evaluate Redis Sentinel/Cluster for redis-group
   - Consider read replicas for postgres services where possible
   - Add HA for stateless services

3. **Complete 1:1 Parity Verification**
   - Systematic comparison of all services
   - Document any necessary differences
   - Update Nomad HCL to match docker-compose exactly

### Long-term (Priority: LOW)
1. **Build Process for nuq-postgres**
   - Set up CI/CD to build and push image
   - Or use official image if available
   - Update Nomad HCL to use registry image

2. **Comprehensive Testing**
   - Test all service endpoints
   - Verify failover scenarios
   - Load testing for HA services

## Status Summary

| Requirement | Status | Notes |
|------------|--------|-------|
| 1:1 Parity | ⚠️ Partial | Most match, some differences |
| Fully Healthy | ⚠️ Mostly | Most running, some issues |
| Complete Failover | ⚠️ Partial | 7 services HA, many single points |
| Zero SPOF | ❌ No | Consul SPOF, many services count=1 |
| All Nodes Operational | ⚠️ Partial | 2 ready, 2 down |

## Overall Status: 🟡 PARTIALLY HEALTHY

**Cluster is functional with most services operational, but critical infrastructure gaps (Consul HA, node recovery) need addressing for production readiness and zero single points of failure.**
