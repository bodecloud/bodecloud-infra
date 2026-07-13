# PaaS - Platform as a Service Converter

A Go platform abstraction layer for loading, normalizing, preserving, and converting container orchestration formats used by the Bolabaden infrastructure plan.

The package has two layers:

- `Application`: the existing ergonomic service/network/volume/config/secret model used by converters and infra integration.
- `CanonicalApplication`: a universal resource graph that keeps normalized resources, routes, policies, and raw source-platform objects side by side so adapters can bridge formats without silently throwing away unsupported fields.

## Features

- **Multi-format Support**: Load and convert between Docker Compose YAML, Docker Swarm stack files, Nomad HCL, Kubernetes YAML, and Helm chart directories
- **Compose Spec Parser**: Uses `github.com/compose-spec/compose-go/v2` as the primary Docker Compose/Swarm parser, with the legacy YAML parser retained as a fallback; preserved source YAML can be restored across conversion hops from the canonical graph
- **Canonical Resource Graph**: Preserves services, networks, volumes, configs, secrets, routes, policies, deploy policy, raw compose-go projects, raw Helm chart metadata, and unknown Kubernetes resources, including preservation and re-emission across conversion hops inside the canonical layer
- **Routing and Policy Abstraction**: Normalizes Traefik Compose labels, Kubernetes Ingress, and Kubernetes NetworkPolicy into portable route/policy resources for failover and ACL tooling
- **Route/Policy Emitters**: Emits canonical routes back to Traefik labels for Compose/Swarm and to `networking.k8s.io/v1` Ingress resources for Kubernetes; emits canonical network policies to Kubernetes `NetworkPolicy`
- **Deploy Intent Bridging**: Maps Swarm placement, rollout, restart, replica, and resource intent into Kubernetes-native fields where possible, with loss-aware annotations for fields without direct Kubernetes equivalents
- **Runtime Security Bridging**: Carries user/group identity, sysctls, process/IPC namespace modes, shm size, privileged mode, Linux capabilities, read-only root filesystems, TTY/stdin, init, stop behavior, working directory, entrypoint, and command args across Compose/Swarm, Kubernetes, and Nomad where native equivalents exist
- **DNS and Host Alias Bridging**: Carries Compose/Swarm DNS servers, search domains, resolver options, and extra hosts into Kubernetes pod DNS/hostAliases and Nomad Docker task DNS/extra host fields
- **Volume Mount Bridging**: Carries Compose/Swarm bind, named volume, and tmpfs mounts into Nomad Docker `mount` blocks, including read-only flags, bind propagation, volume `nocopy`, and tmpfs size
- **Helm Chart Support**: Uses the Helm v3 SDK to load chart directories, render templates with Helm/Sprig behavior, parse the resulting Kubernetes manifests, and preserve raw chart files for restoration across conversion hops
- **Roundtrip Testing**: Comprehensive tests ensure fidelity across format conversions
- **Infra Integration**: Direct deployment to Go-based infrastructure code and loading generated services back into the portable model
- **CLI Tool**: Command-line interface for all operations
- **Application Merging**: Combine multiple applications into unified deployments
- **Validation**: Built-in validation for all supported formats

## Supported Formats

- **Docker Compose Spec**: Compose YAML loaded through compose-go, including services, networks, volumes, configs, secrets, ports, environment, health checks, and deploy policy
- **Docker Swarm / Mirantis Stack Files**: Compose-spec files with Swarm `deploy` fields; detect with `x-platform: docker-swarm`, `swarm`, or `stack` naming
- **Nomad HCL**: HashiCorp Nomad job specifications parsed through the HashiCorp HCL syntax AST, with raw HCL preserved for restoration across conversion hops
- **Kubernetes YAML**: Multi-document manifests for Deployments, StatefulSets, DaemonSets, Jobs, CronJobs, Services, ConfigMaps, Secrets, PVCs, Ingresses, and preserved raw/unknown resources that can be restored across conversion hops. Service ports are reconciled to workloads by same-name match or `spec.selector` label matching, independent of manifest order, including named `targetPort` references to named container ports.
- **Helm Charts**: Directory-based chart loading through the Helm v3 SDK, including `.helmignore`, helpers, Sprig functions, release metadata, values coalescing, and dependency-aware rendering

## Installation

```bash
cd paas
go build -o paas cmd/paas/main.go
```

## Usage

### Basic Conversion

```bash
# Convert Docker Compose to Nomad HCL
./paas -input docker-compose.yml -output nomad.hcl -from docker-compose -to nomad

# Convert Nomad HCL to Kubernetes YAML
./paas -input nomad.hcl -output k8s.yaml -from nomad -to kubernetes

# Convert Kubernetes to Docker Compose
./paas -input k8s.yaml -output docker-compose.yml -from kubernetes -to docker-compose
```

### Validation and Analysis

```bash
# Validate a Docker Compose file
./paas -input docker-compose.yml -validate

# List all services in an application
./paas -input nomad.hcl -list-services

# Restore preserved Kubernetes source back to a file
./paas -input app.yml -restore-format kubernetes -output restored.yaml

# Restore preserved source through the generic restore entrypoint
./paas -input app.yml -restore-format nomad -output restored.hcl
./paas -input app.yml -restore-format kubernetes
```

### Application Merging

```bash
# Merge multiple compose files
./paas -merge app.yml,db.yml,monitoring.yml -output merged.yml
```

### Infra Deployment

```bash
# Deploy Docker Compose to Go infra code
./paas -input docker-compose.yml -deploy -infra-path ../infra
```

## Architecture

### Core Components

- **`model.go`**: Unified data model representing services, networks, volumes, configs, and secrets
- **`canonical.go`**: Cross-orchestrator canonical resource graph and raw resource preservation
- **`paas.go`**: Main PaaS engine with conversion and validation logic
- **`docker_compose.go`**: Docker Compose parser and serializer
- **`compose_go.go`**: compose-go adapter from the Compose Specification project model into `Application`
- **`nomad.go`**: Nomad HCL parser and serializer
- **`kubernetes.go`**: Kubernetes YAML parser and serializer
- **`helm.go`**: Helm SDK-backed chart loading/rendering plus generated chart output
- **`infra_integration.go`**: Integration with existing Go infrastructure

### Data Model

The unified `Application` struct supports:

```go
type Application struct {
    Version   string
    Platform  Platform
    Services  map[string]*Service
    Networks  map[string]*Network
    Volumes   map[string]*Volume
    Configs   map[string]*Config
    Secrets   map[string]*Secret
    Extensions map[string]interface{} // Platform-specific extensions
    Canonical *CanonicalApplication
}
```

The canonical model is attached automatically by parsers and conversions:

```go
type CanonicalApplication struct {
    Source    Platform
    Services  map[string]*Service
    Networks  map[string]*Network
    Volumes   map[string]*Volume
    Configs   map[string]*Config
    Secrets   map[string]*Secret
    Routes    map[string]*RouteSpec
    Policies  map[string]*PolicySpec
    Resources map[string]*CanonicalResource
}
```

Each `CanonicalResource` has a stable ID such as `docker-swarm:service:api`, `docker-swarm:route:api`, `helm:raw:chart`, or `kubernetes:unknown:custom-object`. Raw resources are intentionally retained so unsupported platform fields remain inspectable instead of being silently flattened.

`RouteSpec` and `PolicySpec` are the first portable inputs for the master plan's failover, DNS routing, ACL, and rate-limit work. They currently extract:

- Traefik router hosts, paths, entrypoints, TLS, and middleware chains from Compose/Swarm labels.
- Auth/rate-limit/middleware policy hints from Traefik labels.
- Kubernetes Ingress hosts, paths, backend service, backend port, and TLS.
- Kubernetes NetworkPolicy identity and selector labels.
- Kubernetes Service objects as raw canonical resources, with numeric and named ports merged into matching workloads by name or selector.

They currently emit:

- Traefik router/service labels from canonical routes when serializing Compose/Swarm.
- Kubernetes `Ingress` resources from canonical routes when serializing Kubernetes YAML.
- Kubernetes `NetworkPolicy` resources from canonical `networkpolicy` policies.

`DeploySpec` carries portable deployment intent across Compose, Swarm, Kubernetes, Nomad, and rendered Helm manifests. Kubernetes serialization maps the direct fields to native manifests:

- `Replicas` -> `Deployment.spec.replicas`
- resource limits/reservations -> container `resources.limits` and `resources.requests`
- simple `node.labels.<key> == <value>` placement constraints -> pod `nodeSelector`
- Swarm `update_config.order: start-first` -> `strategy.type: RollingUpdate` with `maxUnavailable: 0` and `maxSurge: 1`

Fields without direct Kubernetes equivalents are retained as `bolabaden.dev/*` annotations on the generated Deployment, including placement constraints, placement preferences, update delay/failure action, restart-policy condition/delay/attempt/window, and runtime fields such as Compose `security_opt`, `init`, and `stop_signal`. Kubernetes parsing rehydrates those annotations back into portable specs.

Nomad serialization maps portable placement into group-level scheduler intent:

- `node.labels.<key> == <value>` -> `constraint { attribute = "${meta.<key>}" operator = "=" value = "<value>" }`
- `node.hostname == <value>` -> `constraint { attribute = "${node.unique.name}" ... }`
- `nomad.attr.<path> == <value>` -> `constraint { attribute = "${attr.<path>}" ... }`
- parsed Nomad constraints and affinities rehydrate into `DeploySpec.Placement`

Runtime fields are mapped to native platform constructs when possible:

- Compose/Swarm `entrypoint` -> Kubernetes container `command`; Compose/Swarm `command` -> Kubernetes container `args`
- Compose/Swarm numeric `user`/`group` -> Kubernetes `securityContext.runAsUser` and `runAsGroup`
- Compose/Swarm `privileged`, `cap_add`, `cap_drop`, and `read_only` -> Kubernetes container `securityContext`
- Compose/Swarm working directory, stdin, TTY, and stop grace period -> Kubernetes container/pod fields
- Compose/Swarm command, args, working directory, user, privileged mode, capabilities, security options, and read-only root filesystems -> Nomad Docker task HCL

Service DNS fields are mapped to native platform constructs where possible:

- Compose/Swarm `dns`, `dns_search`, and `dns_opt` -> Kubernetes pod `dnsPolicy: None` plus `dnsConfig`
- Compose/Swarm `extra_hosts` -> Kubernetes pod `hostAliases`
- Compose/Swarm `dns`, `dns_search`, `dns_opt`, and `extra_hosts` -> Nomad Docker task `dns_servers`, `dns_search_domains`, `dns_options`, and `extra_hosts`

Service volume mounts are mapped to native platform constructs where possible:

- Compose/Swarm long-form bind, named volume, and tmpfs mounts -> Nomad Docker `mount` blocks
- Bind `read_only` and `bind.propagation` -> Nomad `readonly` and `bind_options.propagation`
- Named volume `volume.nocopy` -> Nomad `volume_options.no_copy`
- Tmpfs `tmpfs.size` -> Nomad `tmpfs_options.size`

`PortMapping` also preserves port names across platforms. Kubernetes container port names, Service port names, and named Service `targetPort` values are tracked separately enough to round-trip a Service such as `name: web` with `targetPort: http-api` back to a container port named `http-api`. Nomad network port labels and Compose long-form port names use the same portable field where possible.

### Roundtrip Testing

The system includes comprehensive roundtrip tests that verify:

1. **Docker Compose → Nomad → Docker Compose**
2. **Docker Compose → Kubernetes → Docker Compose**
3. **Nomad → Kubernetes → Nomad**
4. **Serialize/Deserialize fidelity**
5. **Swarm deploy intent -> Kubernetes native fields and annotations -> `DeploySpec`**
6. **Compose runtime security/process fields -> Kubernetes/Nomad native fields -> portable model**
7. **Compose DNS and host alias fields -> Kubernetes/Nomad native fields -> portable model**
8. **Compose service volume mounts -> Nomad Docker mount blocks -> Compose long syntax**

Run tests with:

```bash
go test -v
```

## API Usage

### Programmatic Usage

```go
package main

import (
    "github.com/your-org/my-media-stack/paas"
)

func main() {
    // Create PaaS instance
    paas := paas.New(&paas.PaaSConfig{
        WorkDir: "/tmp/paas",
    })

    // Load application
    app, err := paas.LoadFile("docker-compose.yml")
    if err != nil {
        panic(err)
    }

    // Convert to different format
    nomadApp, err := paas.Convert(app, paas.PlatformDockerCompose, paas.PlatformNomad)
    if err != nil {
        panic(err)
    }

    // Save converted application
    err = paas.SaveFile(nomadApp, "nomad.hcl")
    if err != nil {
        panic(err)
    }
}
```

### Infra Integration

```go
// Deploy application to Go infrastructure
integration, err := paas.NewInfraIntegration("../infra")
if err != nil {
    panic(err)
}

err = integration.DeployToInfra(app, "services_generated")
if err != nil {
    panic(err)
}
```

## Limitations

### Current Limitations

1. **Generated Helm charts are still partial**: parsing uses the Helm SDK renderer, raw chart files can be restored across conversion hops, and generated charts now emit a basic `values.yaml` plus templated core deployment fields, but `SerializeHelmChart` still does not reconstruct a fully idiomatic chart for every object and helper pattern.
2. **Kubernetes relationship modeling is partial**: common workload, Service selector/port relationships, ingress, policy, deploy intent, and support objects are normalized; unknown objects are preserved canonically, not fully interpreted.
3. **Nomad parsing covers common job/group/task Docker service fields**: constraints and simple affinities are mapped into portable placement, while advanced scheduler constructs, spread blocks, and Consul Connect need deeper canonical mapping.
4. **Round-trip fidelity is loss-aware, not lossless**: source-specific raw resources are preserved for inspection and across conversion hops in the canonical layer, but emitters still serialize through the simplified target model unless a target-specific adapter is added.
5. **Runtime unification is not implemented here**: sync agents, failover agents, Headscale HA, DDNS, and rate-limit lobby services are separate implementation modules from the infrastructure master plan.

### Future Enhancements

- Target-specific emitters that can rehydrate preserved raw platform fields
- Full Kubernetes resource relationship graph with owner references, selectors, policies, and ingress/routing objects
- Service mesh integration
- CI/CD pipeline integration
- Multi-cluster deployment support

## Testing

### Unit Tests

```bash
go test ./...
```

### Roundtrip Tests

```bash
go test -run TestRoundTrip
```

### Benchmarking

```bash
go test -bench=.
```

## Contributing

1. Follow Go best practices and conventions
2. Add tests for new functionality
3. Ensure roundtrip tests pass
4. Update documentation
5. Use conventional commit messages

## License

This project is part of the my-media-stack infrastructure and follows the same licensing terms.
