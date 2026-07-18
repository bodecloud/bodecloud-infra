package paas

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Platform represents the source platform of a service
type Platform string

const (
	PlatformDockerCompose Platform = "docker-compose"
	PlatformDockerSwarm   Platform = "docker-swarm"
	PlatformNomad         Platform = "nomad"
	PlatformKubernetes    Platform = "kubernetes"
	PlatformHelm          Platform = "helm"
)

// Service represents a containerized service that can be deployed
type Service struct {
	// Basic identification
	Name          string `json:"name" yaml:"name"`
	Image         string `json:"image" yaml:"image"`
	ContainerName string `json:"container_name,omitempty" yaml:"container_name,omitempty"`
	Hostname      string `json:"hostname,omitempty" yaml:"hostname,omitempty"`

	// Networking
	Ports              []PortMapping                 `json:"ports,omitempty" yaml:"ports,omitempty"`
	Networks           []string                      `json:"networks,omitempty" yaml:"networks,omitempty"`
	NetworkAttachments map[string]*NetworkAttachment `json:"network_attachments,omitempty" yaml:"network_attachments,omitempty"`
	Expose             []string                      `json:"expose,omitempty" yaml:"expose,omitempty"`
	DNS                []string                      `json:"dns,omitempty" yaml:"dns,omitempty"`
	DNSSearch          []string                      `json:"dns_search,omitempty" yaml:"dns_search,omitempty"`
	DNSOptions         []string                      `json:"dns_opt,omitempty" yaml:"dns_opt,omitempty"`
	ExtraHosts         []string                      `json:"extra_hosts,omitempty" yaml:"extra_hosts,omitempty"`
	HostAliases        []HostAlias                   `json:"host_aliases,omitempty" yaml:"host_aliases,omitempty"`

	// Environment and configuration
	Environment                  map[string]string              `json:"environment,omitempty" yaml:"environment,omitempty"`
	EnvFile                      []string                       `json:"env_file,omitempty" yaml:"env_file,omitempty"`
	EnvFileRefs                  []EnvFileRef                   `json:"env_file_refs,omitempty" yaml:"env_file_refs,omitempty"`
	EnvSources                   []EnvSource                    `json:"env_sources,omitempty" yaml:"env_sources,omitempty"`
	EnvFrom                      []EnvFromSource                `json:"env_from,omitempty" yaml:"env_from,omitempty"`
	ImagePullSecrets             []string                       `json:"image_pull_secrets,omitempty" yaml:"image_pull_secrets,omitempty"`
	ImagePullPolicy              string                         `json:"image_pull_policy,omitempty" yaml:"image_pull_policy,omitempty"`
	TerminationMessagePath       string                         `json:"termination_message_path,omitempty" yaml:"termination_message_path,omitempty"`
	TerminationMessagePolicy     string                         `json:"termination_message_policy,omitempty" yaml:"termination_message_policy,omitempty"`
	Tolerations                  []Toleration                   `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	HostNetwork                  bool                           `json:"host_network,omitempty" yaml:"host_network,omitempty"`
	HostNetworkSet               bool                           `json:"-" yaml:"-"`
	HostPID                      *bool                          `json:"host_pid,omitempty" yaml:"host_pid,omitempty"`
	HostIPC                      *bool                          `json:"host_ipc,omitempty" yaml:"host_ipc,omitempty"`
	PIDMode                      string                         `json:"pid_mode,omitempty" yaml:"pid_mode,omitempty"`
	IPCMode                      string                         `json:"ipc_mode,omitempty" yaml:"ipc_mode,omitempty"`
	HostUsers                    *bool                          `json:"host_users,omitempty" yaml:"host_users,omitempty"`
	DNSPolicy                    string                         `json:"dns_policy,omitempty" yaml:"dns_policy,omitempty"`
	OSName                       string                         `json:"os_name,omitempty" yaml:"os_name,omitempty"`
	SchedulerName                string                         `json:"scheduler_name,omitempty" yaml:"scheduler_name,omitempty"`
	PriorityClassName            string                         `json:"priority_class_name,omitempty" yaml:"priority_class_name,omitempty"`
	RuntimeClassName             string                         `json:"runtime_class_name,omitempty" yaml:"runtime_class_name,omitempty"`
	NodeName                     string                         `json:"node_name,omitempty" yaml:"node_name,omitempty"`
	NodeSelector                 map[string]string              `json:"node_selector,omitempty" yaml:"node_selector,omitempty"`
	Subdomain                    string                         `json:"subdomain,omitempty" yaml:"subdomain,omitempty"`
	SetHostnameAsFQDN            *bool                          `json:"set_hostname_as_fqdn,omitempty" yaml:"set_hostname_as_fqdn,omitempty"`
	ShareProcessNamespace        *bool                          `json:"share_process_namespace,omitempty" yaml:"share_process_namespace,omitempty"`
	EnableServiceLinks           *bool                          `json:"enable_service_links,omitempty" yaml:"enable_service_links,omitempty"`
	FSGroup                      *int64                         `json:"fs_group,omitempty" yaml:"fs_group,omitempty"`
	SELinuxOptions               *SELinuxOptions                `json:"se_linux_options,omitempty" yaml:"se_linux_options,omitempty"`
	WindowsOptions               *WindowsSecurityContextOptions `json:"windows_options,omitempty" yaml:"windows_options,omitempty"`
	RunAsNonRoot                 *bool                          `json:"run_as_non_root,omitempty" yaml:"run_as_non_root,omitempty"`
	SupplementalGroups           []int64                        `json:"supplemental_groups,omitempty" yaml:"supplemental_groups,omitempty"`
	ActiveDeadlineSeconds        *int64                         `json:"active_deadline_seconds,omitempty" yaml:"active_deadline_seconds,omitempty"`
	PodRestartPolicy             string                         `json:"pod_restart_policy,omitempty" yaml:"pod_restart_policy,omitempty"`
	ServiceAccountName           string                         `json:"service_account_name,omitempty" yaml:"service_account_name,omitempty"`
	AutomountServiceAccountToken *bool                          `json:"automount_service_account_token,omitempty" yaml:"automount_service_account_token,omitempty"`
	Command                      []string                       `json:"command,omitempty" yaml:"command,omitempty"`
	Entrypoint                   []string                       `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	WorkingDir                   string                         `json:"working_dir,omitempty" yaml:"working_dir,omitempty"`
	GroupAdd                     []string                       `json:"group_add,omitempty" yaml:"group_add,omitempty"`
	Sysctls                      map[string]string              `json:"sysctls,omitempty" yaml:"sysctls,omitempty"`
	Build                        *BuildConfig                   `json:"build,omitempty" yaml:"build,omitempty"`
	Develop                      *DevelopConfig                 `json:"develop,omitempty" yaml:"develop,omitempty"`
	Lifecycle                    *LifecycleHooks                `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
	Devices                      []string                       `json:"devices,omitempty" yaml:"devices,omitempty"`
	DeviceMappings               []DeviceMappingSpec            `json:"device_mappings,omitempty" yaml:"device_mappings,omitempty"`
	Profiles                     []string                       `json:"profiles,omitempty" yaml:"profiles,omitempty"`

	// Volumes and mounts
	Volumes []VolumeMount `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Configs []FileRef     `json:"configs,omitempty" yaml:"configs,omitempty"`
	Secrets []FileRef     `json:"secrets,omitempty" yaml:"secrets,omitempty"`

	// Dependencies and ordering
	DependsOn    []string         `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
	Dependencies []DependencySpec `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Links        []string         `json:"links,omitempty" yaml:"links,omitempty"`

	// Runtime behavior
	Restart           string                 `json:"restart,omitempty" yaml:"restart,omitempty"`
	Privileged        bool                   `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	PrivilegedSet     bool                   `json:"-" yaml:"-"`
	User              string                 `json:"user,omitempty" yaml:"user,omitempty"`
	Group             string                 `json:"group,omitempty" yaml:"group,omitempty"`
	CapAdd            []string               `json:"cap_add,omitempty" yaml:"cap_add,omitempty"`
	CapDrop           []string               `json:"cap_drop,omitempty" yaml:"cap_drop,omitempty"`
	SecurityOpt       []string               `json:"security_opt,omitempty" yaml:"security_opt,omitempty"`
	ReadOnlyRootFS    bool                   `json:"read_only,omitempty" yaml:"read_only,omitempty"`
	ReadOnlyRootFSSet bool                   `json:"-" yaml:"-"`
	Init              *bool                  `json:"init,omitempty" yaml:"init,omitempty"`
	Tty               bool                   `json:"tty,omitempty" yaml:"tty,omitempty"`
	TtySet            bool                   `json:"-" yaml:"-"`
	StdinOpen         bool                   `json:"stdin_open,omitempty" yaml:"stdin_open,omitempty"`
	StdinOpenSet      bool                   `json:"-" yaml:"-"`
	StopSignal        string                 `json:"stop_signal,omitempty" yaml:"stop_signal,omitempty"`
	StopGracePeriod   string                 `json:"stop_grace_period,omitempty" yaml:"stop_grace_period,omitempty"`
	Runtime           string                 `json:"runtime,omitempty" yaml:"runtime,omitempty"`
	LogDriver         string                 `json:"log_driver,omitempty" yaml:"log_driver,omitempty"`
	LogOpt            map[string]string      `json:"log_opt,omitempty" yaml:"log_opt,omitempty"`
	LogExtensions     map[string]interface{} `json:"log_extensions,omitempty" yaml:"log_extensions,omitempty"`
	ComposeCompat     *ComposeCompat         `json:"compose_compat,omitempty" yaml:"compose_compat,omitempty"`
	Connect           *NomadConnectSpec      `json:"connect,omitempty" yaml:"connect,omitempty"`
	Spreads           []NomadSpreadSpec      `json:"spreads,omitempty" yaml:"spreads,omitempty"`
	PidsLimit         int64                  `json:"pids_limit,omitempty" yaml:"pids_limit,omitempty"`
	pidsLimitSet      bool                   `json:"-" yaml:"-"`
	ShmSize           int64                  `json:"shm_size,omitempty" yaml:"shm_size,omitempty"`
	shmSizeSet        bool                   `json:"-" yaml:"-"`
	Replicas          int                    `json:"replicas,omitempty" yaml:"replicas,omitempty"`

	// Resource limits
	CPUShares      int      `json:"cpu_shares,omitempty" yaml:"cpu_shares,omitempty"`
	CPUQuota       int      `json:"cpu_quota,omitempty" yaml:"cpu_quota,omitempty"`
	MemoryLimit    string   `json:"memory_limit,omitempty" yaml:"memory_limit,omitempty"`
	MemorySwap     string   `json:"memory_swap,omitempty" yaml:"memory_swap,omitempty"`
	MemLimit       string   `json:"mem_limit,omitempty" yaml:"mem_limit,omitempty"`
	MemReservation string   `json:"mem_reservation,omitempty" yaml:"mem_reservation,omitempty"`
	CPUs           string   `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	Ulimits        *Ulimits `json:"ulimits,omitempty" yaml:"ulimits,omitempty"`
	UserNSMode     string   `json:"userns_mode,omitempty" yaml:"userns_mode,omitempty"`
	PullPolicy     string   `json:"pull_policy,omitempty" yaml:"pull_policy,omitempty"`

	// Orchestrator-neutral deployment policy. Docker Swarm, Kubernetes, Nomad,
	// and Helm all express these differently, so keep the shared intent here and
	// preserve lossier source details in Extensions.
	Deploy *DeploySpec `json:"deploy,omitempty" yaml:"deploy,omitempty"`

	// Peer-aware failover intent used by the Bolabaden mesh/Headscale plan.
	Failover *FailoverSpec `json:"failover,omitempty" yaml:"failover,omitempty"`

	// Health checks
	HealthCheck               *HealthCheck             `json:"healthcheck,omitempty" yaml:"healthcheck,omitempty"`
	StartupProbe              *HealthCheck             `json:"startup_probe,omitempty" yaml:"startup_probe,omitempty"`
	SeccompProfile            *SeccompProfile          `json:"seccomp_profile,omitempty" yaml:"seccomp_profile,omitempty"`
	AllowPrivilegeEscalation  *bool                    `json:"allow_privilege_escalation,omitempty" yaml:"allow_privilege_escalation,omitempty"`
	ProcMount                 string                   `json:"proc_mount,omitempty" yaml:"proc_mount,omitempty"`
	FSGroupChangePolicy       string                   `json:"fs_group_change_policy,omitempty" yaml:"fs_group_change_policy,omitempty"`
	Affinity                  map[string]interface{}   `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	ReadinessGates            []map[string]interface{} `json:"readiness_gates,omitempty" yaml:"readiness_gates,omitempty"`
	InitContainers            []map[string]interface{} `json:"init_containers,omitempty" yaml:"init_containers,omitempty"`
	ResourceClaims            []map[string]interface{} `json:"resource_claims,omitempty" yaml:"resource_claims,omitempty"`
	EphemeralContainers       []map[string]interface{} `json:"ephemeral_containers,omitempty" yaml:"ephemeral_containers,omitempty"`
	SchedulingGates           []map[string]interface{} `json:"scheduling_gates,omitempty" yaml:"scheduling_gates,omitempty"`
	TopologySpreadConstraints []map[string]interface{} `json:"topology_spread_constraints,omitempty" yaml:"topology_spread_constraints,omitempty"`
	SupplementalGroupsPolicy  string                   `json:"supplemental_groups_policy,omitempty" yaml:"supplemental_groups_policy,omitempty"`

	// Labels and metadata
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// Platform-specific extensions
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`

	// Source platform metadata
	Platform   Platform `json:"platform" yaml:"platform"`
	SourceFile string   `json:"source_file,omitempty" yaml:"source_file,omitempty"`
}

// ComposeModel captures a top-level compose-spec model definition.
type ComposeModel struct {
	Name         string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Model        string                 `json:"model,omitempty" yaml:"model,omitempty"`
	ContextSize  int                    `json:"context_size,omitempty" yaml:"context_size,omitempty"`
	RuntimeFlags []string               `json:"runtime_flags,omitempty" yaml:"runtime_flags,omitempty"`
	Extensions   map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// DeviceMappingSpec captures compose-spec service-level device mappings,
// including object-form extension fields that string syntax cannot represent.
type DeviceMappingSpec struct {
	Source      string                 `json:"source,omitempty" yaml:"source,omitempty"`
	Target      string                 `json:"target,omitempty" yaml:"target,omitempty"`
	Permissions string                 `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// DeploySpec captures rollout and placement intent shared by orchestrators.
type DeploySpec struct {
	Mode             string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	EndpointMode     string                 `json:"endpoint_mode,omitempty" yaml:"endpoint_mode,omitempty"`
	Replicas         int                    `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Job              *SwarmJobSpec          `json:"job,omitempty" yaml:"job,omitempty"`
	Placement        *PlacementSpec         `json:"placement,omitempty" yaml:"placement,omitempty"`
	Resources        *ResourceSpec          `json:"resources,omitempty" yaml:"resources,omitempty"`
	UpdateConfig     *UpdatePolicy          `json:"update_config,omitempty" yaml:"update_config,omitempty"`
	MigrateConfig    *MigratePolicy         `json:"migrate_config,omitempty" yaml:"migrate_config,omitempty"`
	RescheduleConfig *ReschedulePolicy      `json:"reschedule_config,omitempty" yaml:"reschedule_config,omitempty"`
	RollbackConfig   *UpdatePolicy          `json:"rollback_config,omitempty" yaml:"rollback_config,omitempty"`
	RestartPolicy    *RestartPolicy         `json:"restart_policy,omitempty" yaml:"restart_policy,omitempty"`
	Labels           map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions       map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// SwarmJobSpec captures Docker/Moby job service counters that Compose exposes
// only as extension data and other orchestrators cannot represent natively.
type SwarmJobSpec struct {
	MaxConcurrent              int                    `json:"max_concurrent,omitempty" yaml:"max_concurrent,omitempty"`
	maxConcurrentSet           bool                   `json:"-" yaml:"-"`
	TotalCompletions           int                    `json:"total_completions,omitempty" yaml:"total_completions,omitempty"`
	totalCompletionsSet        bool                   `json:"-" yaml:"-"`
	CompletionMode             string                 `json:"completion_mode,omitempty" yaml:"completion_mode,omitempty"`
	completionModeSet          bool                   `json:"-" yaml:"-"`
	Suspend                    *bool                  `json:"suspend,omitempty" yaml:"suspend,omitempty"`
	BackoffLimit               int                    `json:"backoff_limit,omitempty" yaml:"backoff_limit,omitempty"`
	backoffLimitSet            bool                   `json:"-" yaml:"-"`
	BackoffLimitPerIndex       int                    `json:"backoff_limit_per_index,omitempty" yaml:"backoff_limit_per_index,omitempty"`
	backoffLimitPerIndexSet    bool                   `json:"-" yaml:"-"`
	TTLSecondsAfterFinished    int                    `json:"ttl_seconds_after_finished,omitempty" yaml:"ttl_seconds_after_finished,omitempty"`
	ttlSecondsAfterFinishedSet bool                   `json:"-" yaml:"-"`
	Extensions                 map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func (job *SwarmJobSpec) MarshalJSON() ([]byte, error) {
	if job == nil {
		return []byte("null"), nil
	}
	return json.Marshal(serializeSwarmJobSpec(job))
}

func (job *SwarmJobSpec) UnmarshalJSON(data []byte) error {
	if job == nil {
		return fmt.Errorf("nil SwarmJobSpec")
	}
	parsed, err := parseSwarmJobSpecJSON(data)
	if err != nil {
		return err
	}
	*job = *parsed
	return nil
}

// FailoverSpec captures peer-forward and singleton-election intent that is not
// natively portable across Compose, Kubernetes, Nomad, Helm, or Swarm.
type FailoverSpec struct {
	Enabled             bool                     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	EnabledSet          bool                     `json:"-" yaml:"-"`
	Type                string                   `json:"type,omitempty" yaml:"type,omitempty"`
	Port                string                   `json:"port,omitempty" yaml:"port,omitempty"`
	HealthcheckPath     string                   `json:"healthcheck_path,omitempty" yaml:"healthcheck_path,omitempty"`
	HealthcheckInterval string                   `json:"healthcheck_interval,omitempty" yaml:"healthcheck_interval,omitempty"`
	MaxRetries          int                      `json:"max_retries,omitempty" yaml:"max_retries,omitempty"`
	RedeployOnPeer      bool                     `json:"redeploy_on_peer,omitempty" yaml:"redeploy_on_peer,omitempty"`
	Singleton           bool                     `json:"singleton,omitempty" yaml:"singleton,omitempty"`
	SingletonElection   string                   `json:"singleton_election,omitempty" yaml:"singleton_election,omitempty"`
	Strategy            string                   `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	PreferLocal         bool                     `json:"prefer_local,omitempty" yaml:"prefer_local,omitempty"`
	Nodes               map[string]*FailoverNode `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Extensions          map[string]interface{}   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// FailoverNode captures one peer candidate in the failover registry.
type FailoverNode struct {
	Status     string                 `json:"status,omitempty" yaml:"status,omitempty"`
	LastSeen   string                 `json:"last_seen,omitempty" yaml:"last_seen,omitempty"`
	Priority   int                    `json:"priority,omitempty" yaml:"priority,omitempty"`
	URL        string                 `json:"url,omitempty" yaml:"url,omitempty"`
	Weight     int                    `json:"weight,omitempty" yaml:"weight,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// MeshSpec captures mesh-network intent such as Headscale control-plane
// configuration and MagicDNS settings.
type MeshSpec struct {
	Enabled          bool                   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	EnabledSet       bool                   `json:"-" yaml:"-"`
	Provider         string                 `json:"provider,omitempty" yaml:"provider,omitempty"`
	ControlPlaneURL  string                 `json:"control_plane_url,omitempty" yaml:"control_plane_url,omitempty"`
	ControlPlaneURLs []string               `json:"control_plane_urls,omitempty" yaml:"control_plane_urls,omitempty"`
	DNSZone          string                 `json:"dns_zone,omitempty" yaml:"dns_zone,omitempty"`
	MagicDNS         bool                   `json:"magic_dns,omitempty" yaml:"magic_dns,omitempty"`
	MagicDNSSet      bool                   `json:"-" yaml:"-"`
	Nodes            map[string]*MeshNode   `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Extensions       map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// MeshNode captures one peer node in the mesh registry.
type MeshNode struct {
	Name       string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Hostname   string                 `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Address    string                 `json:"address,omitempty" yaml:"address,omitempty"`
	Status     string                 `json:"status,omitempty" yaml:"status,omitempty"`
	Role       string                 `json:"role,omitempty" yaml:"role,omitempty"`
	Priority   int                    `json:"priority,omitempty" yaml:"priority,omitempty"`
	URL        string                 `json:"url,omitempty" yaml:"url,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// KubernetesServiceSpec captures a Kubernetes Service manifest in typed form
// while still preserving the original raw object for loss-aware round-trips.
type KubernetesServiceSpec struct {
	Name                          string                  `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace                     string                  `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Type                          string                  `json:"type,omitempty" yaml:"type,omitempty"`
	Selector                      map[string]string       `json:"selector,omitempty" yaml:"selector,omitempty"`
	Ports                         []KubernetesServicePort `json:"ports,omitempty" yaml:"ports,omitempty"`
	ExternalName                  string                  `json:"external_name,omitempty" yaml:"external_name,omitempty"`
	SessionAffinity               string                  `json:"session_affinity,omitempty" yaml:"session_affinity,omitempty"`
	LoadBalancerIP                string                  `json:"load_balancer_ip,omitempty" yaml:"load_balancer_ip,omitempty"`
	LoadBalancerClass             string                  `json:"load_balancer_class,omitempty" yaml:"load_balancer_class,omitempty"`
	LoadBalancerSourceRanges      []string                `json:"load_balancer_source_ranges,omitempty" yaml:"load_balancer_source_ranges,omitempty"`
	ExternalIPs                   []string                `json:"external_ips,omitempty" yaml:"external_ips,omitempty"`
	IPFamilies                    []string                `json:"ip_families,omitempty" yaml:"ip_families,omitempty"`
	IPFamilyPolicy                string                  `json:"ip_family_policy,omitempty" yaml:"ip_family_policy,omitempty"`
	ExternalTrafficPolicy         string                  `json:"external_traffic_policy,omitempty" yaml:"external_traffic_policy,omitempty"`
	InternalTrafficPolicy         string                  `json:"internal_traffic_policy,omitempty" yaml:"internal_traffic_policy,omitempty"`
	TrafficDistribution           string                  `json:"traffic_distribution,omitempty" yaml:"traffic_distribution,omitempty"`
	PublishNotReadyAddresses      *bool                   `json:"publish_not_ready_addresses,omitempty" yaml:"publish_not_ready_addresses,omitempty"`
	AllocateLoadBalancerNodePorts *bool                   `json:"allocate_load_balancer_node_ports,omitempty" yaml:"allocate_load_balancer_node_ports,omitempty"`
	HealthCheckNodePort           int                     `json:"health_check_node_port,omitempty" yaml:"health_check_node_port,omitempty"`
	Annotations                   map[string]string       `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels                        map[string]string       `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions                    map[string]interface{}  `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw                           map[string]interface{}  `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type KubernetesServicePort struct {
	Name        string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Port        int                    `json:"port,omitempty" yaml:"port,omitempty"`
	TargetPort  string                 `json:"target_port,omitempty" yaml:"target_port,omitempty"`
	Protocol    string                 `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	NodePort    int                    `json:"node_port,omitempty" yaml:"node_port,omitempty"`
	AppProtocol string                 `json:"app_protocol,omitempty" yaml:"app_protocol,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// KubernetesServiceAccountSpec captures a Kubernetes ServiceAccount manifest in typed form.
type KubernetesServiceAccountSpec struct {
	Name                         string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace                    string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels                       map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations                  map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Secrets                      []string               `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	ImagePullSecrets             []string               `json:"image_pull_secrets,omitempty" yaml:"image_pull_secrets,omitempty"`
	AutomountServiceAccountToken *bool                  `json:"automount_service_account_token,omitempty" yaml:"automount_service_account_token,omitempty"`
	Extensions                   map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw                          map[string]interface{} `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesHorizontalPodAutoscalerSpec captures a Kubernetes HPA manifest in typed form.
type KubernetesHorizontalPodAutoscalerSpec struct {
	Name        string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace   string                   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	MinReplicas *int                     `json:"min_replicas,omitempty" yaml:"min_replicas,omitempty"`
	MaxReplicas int                      `json:"max_replicas,omitempty" yaml:"max_replicas,omitempty"`
	ScaleTarget map[string]string        `json:"scale_target,omitempty" yaml:"scale_target,omitempty"`
	Metrics     []map[string]interface{} `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	Behavior    map[string]interface{}   `json:"behavior,omitempty" yaml:"behavior,omitempty"`
	Extensions  map[string]interface{}   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw         map[string]interface{}   `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesPodDisruptionBudgetSpec captures a Kubernetes PDB manifest in typed form.
type KubernetesPodDisruptionBudgetSpec struct {
	Name                       string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace                  string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	MinAvailable               string                 `json:"min_available,omitempty" yaml:"min_available,omitempty"`
	MaxUnavailable             string                 `json:"max_unavailable,omitempty" yaml:"max_unavailable,omitempty"`
	UnhealthyPodEvictionPolicy string                 `json:"unhealthy_pod_eviction_policy,omitempty" yaml:"unhealthy_pod_eviction_policy,omitempty"`
	Selector                   map[string]string      `json:"selector,omitempty" yaml:"selector,omitempty"`
	Annotations                map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels                     map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions                 map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw                        map[string]interface{} `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesResourceQuotaSpec captures a Kubernetes ResourceQuota manifest in typed form.
type KubernetesResourceQuotaSpec struct {
	Name          string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace     string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Scopes        []string               `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	ScopeSelector map[string]interface{} `json:"scope_selector,omitempty" yaml:"scope_selector,omitempty"`
	Hard          map[string]interface{} `json:"hard,omitempty" yaml:"hard,omitempty"`
	Annotations   map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels        map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions    map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw           map[string]interface{} `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesLimitRangeSpec captures a Kubernetes LimitRange manifest in typed form.
type KubernetesLimitRangeSpec struct {
	Name        string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace   string                   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Limits      []map[string]interface{} `json:"limits,omitempty" yaml:"limits,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels      map[string]string        `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions  map[string]interface{}   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw         map[string]interface{}   `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesStorageClassSpec captures a Kubernetes StorageClass manifest in typed form.
type KubernetesStorageClassSpec struct {
	Name                 string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace            string                   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Provisioner          string                   `json:"provisioner,omitempty" yaml:"provisioner,omitempty"`
	Parameters           map[string]string        `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	ReclaimPolicy        string                   `json:"reclaim_policy,omitempty" yaml:"reclaim_policy,omitempty"`
	AllowVolumeExpansion *bool                    `json:"allow_volume_expansion,omitempty" yaml:"allow_volume_expansion,omitempty"`
	VolumeBindingMode    string                   `json:"volume_binding_mode,omitempty" yaml:"volume_binding_mode,omitempty"`
	MountOptions         []string                 `json:"mount_options,omitempty" yaml:"mount_options,omitempty"`
	AllowedTopologies    []map[string]interface{} `json:"allowed_topologies,omitempty" yaml:"allowed_topologies,omitempty"`
	Annotations          map[string]string        `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels               map[string]string        `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions           map[string]interface{}   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw                  map[string]interface{}   `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesIngressClassSpec captures a Kubernetes IngressClass manifest in typed form.
type KubernetesIngressClassSpec struct {
	Name        string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace   string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Controller  string                 `json:"controller,omitempty" yaml:"controller,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw         map[string]interface{} `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesWebhookConfigurationSpec captures a Kubernetes webhook configuration in typed form.
type KubernetesWebhookConfigurationSpec struct {
	Name        string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace   string                   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Webhooks    []map[string]interface{} `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels      map[string]string        `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions  map[string]interface{}   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw         map[string]interface{}   `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesCustomResourceDefinitionSpec captures a Kubernetes CRD manifest in typed form.
type KubernetesCustomResourceDefinitionSpec struct {
	Name                     string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace                string                   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Group                    string                   `json:"group,omitempty" yaml:"group,omitempty"`
	Scope                    string                   `json:"scope,omitempty" yaml:"scope,omitempty"`
	Names                    map[string]interface{}   `json:"names,omitempty" yaml:"names,omitempty"`
	Versions                 []map[string]interface{} `json:"versions,omitempty" yaml:"versions,omitempty"`
	Conversion               map[string]interface{}   `json:"conversion,omitempty" yaml:"conversion,omitempty"`
	Validation               map[string]interface{}   `json:"validation,omitempty" yaml:"validation,omitempty"`
	AdditionalPrinterColumns []map[string]interface{} `json:"additional_printer_columns,omitempty" yaml:"additional_printer_columns,omitempty"`
	Extensions               map[string]interface{}   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw                      map[string]interface{}   `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesPriorityClassSpec captures a Kubernetes PriorityClass manifest in typed form.
type KubernetesPriorityClassSpec struct {
	Name             string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace        string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Value            int32                  `json:"value,omitempty" yaml:"value,omitempty"`
	GlobalDefault    bool                   `json:"global_default,omitempty" yaml:"global_default,omitempty"`
	Description      string                 `json:"description,omitempty" yaml:"description,omitempty"`
	PreemptionPolicy string                 `json:"preemption_policy,omitempty" yaml:"preemption_policy,omitempty"`
	Extensions       map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw              map[string]interface{} `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesRuntimeClassSpec captures a Kubernetes RuntimeClass manifest in typed form.
type KubernetesRuntimeClassSpec struct {
	Name       string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Handler    string                 `json:"handler,omitempty" yaml:"handler,omitempty"`
	Overhead   map[string]interface{} `json:"overhead,omitempty" yaml:"overhead,omitempty"`
	Scheduling map[string]interface{} `json:"scheduling,omitempty" yaml:"scheduling,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw        map[string]interface{} `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// KubernetesOpaqueManifestSpec captures a Kubernetes manifest that we bridge
// losslessly without needing a dedicated semantic type yet.
type KubernetesOpaqueManifestSpec struct {
	APIVersion string                 `json:"api_version,omitempty" yaml:"api_version,omitempty"`
	Kind       string                 `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name       string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       map[string]interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Raw        map[string]interface{} `json:"raw,omitempty" yaml:"raw,omitempty"`
}

// NomadConnectSpec captures service mesh connect intent from Nomad and
// compatible portable carriers.
type NomadConnectSpec struct {
	Native         bool                        `json:"native,omitempty" yaml:"native,omitempty"`
	NativeSet      bool                        `json:"-" yaml:"-"`
	SidecarService *NomadConnectSidecarService `json:"sidecar_service,omitempty" yaml:"sidecar_service,omitempty"`
	Gateway        map[string]interface{}      `json:"gateway,omitempty" yaml:"gateway,omitempty"`
	Extensions     map[string]interface{}      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func (connect *NomadConnectSpec) MarshalJSON() ([]byte, error) {
	if connect == nil {
		return []byte("null"), nil
	}
	return json.Marshal(serializeNomadConnectSpec(connect))
}

func (connect *NomadConnectSpec) UnmarshalJSON(data []byte) error {
	if connect == nil {
		return fmt.Errorf("nil NomadConnectSpec")
	}
	parsed, err := parseNomadConnectSpecMapBytes(data)
	if err != nil {
		return err
	}
	*connect = *parsed
	return nil
}

type NomadConnectSidecarService struct {
	Tags       []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Proxy      *NomadConnectProxy     `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	Check      map[string]interface{} `json:"check,omitempty" yaml:"check,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

type NomadConnectProxy struct {
	Upstreams  []NomadConnectUpstream `json:"upstreams,omitempty" yaml:"upstreams,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

type NomadConnectUpstream struct {
	DestinationName  string                 `json:"destination_name,omitempty" yaml:"destination_name,omitempty"`
	LocalBindPort    int                    `json:"local_bind_port,omitempty" yaml:"local_bind_port,omitempty"`
	LocalBindAddress string                 `json:"local_bind_address,omitempty" yaml:"local_bind_address,omitempty"`
	Extensions       map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// NomadSpreadSpec captures Nomad spread blocks and their weighted targets.
type NomadSpreadSpec struct {
	Attribute  string                 `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Portable   string                 `json:"portable,omitempty" yaml:"portable,omitempty"`
	Weight     int                    `json:"weight,omitempty" yaml:"weight,omitempty"`
	Targets    []NomadSpreadTarget    `json:"targets,omitempty" yaml:"targets,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

type NomadSpreadTarget struct {
	Value      string                 `json:"value,omitempty" yaml:"value,omitempty"`
	Percent    int                    `json:"percent,omitempty" yaml:"percent,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// HostAlias captures a Kubernetes hostAliases entry.
type HostAlias struct {
	IP         string                 `json:"ip" yaml:"ip"`
	Hostnames  []string               `json:"hostnames,omitempty" yaml:"hostnames,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// PlacementSpec captures scheduler constraints and preferences.
type PlacementSpec struct {
	Constraints          []string                 `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	Preferences          []string                 `json:"preferences,omitempty" yaml:"preferences,omitempty"`
	MaxReplicasPerNode   int                      `json:"max_replicas_per_node,omitempty" yaml:"max_replicas_per_node,omitempty"`
	Extensions           map[string]interface{}   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	PreferenceExtensions []map[string]interface{} `json:"preference_extensions,omitempty" yaml:"preference_extensions,omitempty"`
}

// ResourceSpec captures portable resource reservations/limits.
type ResourceSpec struct {
	CPUReservation              string                 `json:"cpu_reservation,omitempty" yaml:"cpu_reservation,omitempty"`
	CPULimit                    string                 `json:"cpu_limit,omitempty" yaml:"cpu_limit,omitempty"`
	MemoryReservation           string                 `json:"memory_reservation,omitempty" yaml:"memory_reservation,omitempty"`
	MemoryLimit                 string                 `json:"memory_limit,omitempty" yaml:"memory_limit,omitempty"`
	EphemeralStorageReservation string                 `json:"ephemeral_storage_reservation,omitempty" yaml:"ephemeral_storage_reservation,omitempty"`
	EphemeralStorageLimit       string                 `json:"ephemeral_storage_limit,omitempty" yaml:"ephemeral_storage_limit,omitempty"`
	PidsLimit                   int64                  `json:"pids_limit,omitempty" yaml:"pids_limit,omitempty"`
	pidsLimitSet                bool                   `json:"-" yaml:"-"`
	PidsReservation             int64                  `json:"pids_reservation,omitempty" yaml:"pids_reservation,omitempty"`
	pidsReservationSet          bool                   `json:"-" yaml:"-"`
	Devices                     []ResourceDevice       `json:"devices,omitempty" yaml:"devices,omitempty"`
	GenericResources            []GenericResource      `json:"generic_resources,omitempty" yaml:"generic_resources,omitempty"`
	Extensions                  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	LimitExtensions             map[string]interface{} `json:"limit_extensions,omitempty" yaml:"limit_extensions,omitempty"`
	ReservationExtensions       map[string]interface{} `json:"reservation_extensions,omitempty" yaml:"reservation_extensions,omitempty"`
}

func (resources *ResourceSpec) MarshalJSON() ([]byte, error) {
	if resources == nil {
		return []byte("null"), nil
	}
	return json.Marshal(serializeResourceSpec(resources))
}

func (resources *ResourceSpec) UnmarshalJSON(data []byte) error {
	if resources == nil {
		return fmt.Errorf("nil ResourceSpec")
	}
	parsed, err := parseResourceSpecJSON(data)
	if err != nil {
		return err
	}
	*resources = *parsed
	return nil
}

// ResourceDevice captures Compose/Swarm deploy resource device reservations.
type ResourceDevice struct {
	Capabilities []string               `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Driver       string                 `json:"driver,omitempty" yaml:"driver,omitempty"`
	Count        string                 `json:"count,omitempty" yaml:"count,omitempty"`
	DeviceIDs    []string               `json:"device_ids,omitempty" yaml:"device_ids,omitempty"`
	Options      map[string]string      `json:"options,omitempty" yaml:"options,omitempty"`
	Extensions   map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GenericResource captures Compose/Swarm discrete generic resources.
type GenericResource struct {
	Kind               string                 `json:"kind,omitempty" yaml:"kind,omitempty"`
	Value              string                 `json:"value,omitempty" yaml:"value,omitempty"`
	Extensions         map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	DiscreteExtensions map[string]interface{} `json:"discrete_extensions,omitempty" yaml:"discrete_extensions,omitempty"`
}

// ComposeCompat captures Compose-native fields that do not have a direct
// portable representation in the shared service model yet.
type ComposeCompat struct {
	Attach            *bool                             `json:"attach,omitempty" yaml:"attach,omitempty"`
	Annotations       map[string]string                 `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	BlkioConfig       map[string]interface{}            `json:"blkio_config,omitempty" yaml:"blkio_config,omitempty"`
	CredentialSpec    map[string]interface{}            `json:"credential_spec,omitempty" yaml:"credential_spec,omitempty"`
	Provider          map[string]interface{}            `json:"provider,omitempty" yaml:"provider,omitempty"`
	Extends           map[string]interface{}            `json:"extends,omitempty" yaml:"extends,omitempty"`
	Platform          string                            `json:"platform,omitempty" yaml:"platform,omitempty"`
	PullRefreshAfter  string                            `json:"pull_refresh_after,omitempty" yaml:"pull_refresh_after,omitempty"`
	MemSwappiness     string                            `json:"mem_swappiness,omitempty" yaml:"mem_swappiness,omitempty"`
	MacAddress        string                            `json:"mac_address,omitempty" yaml:"mac_address,omitempty"`
	DomainName        string                            `json:"domain_name,omitempty" yaml:"domain_name,omitempty"`
	CgroupParent      string                            `json:"cgroup_parent,omitempty" yaml:"cgroup_parent,omitempty"`
	Cgroup            string                            `json:"cgroup,omitempty" yaml:"cgroup,omitempty"`
	CPUCount          int64                             `json:"cpu_count,omitempty" yaml:"cpu_count,omitempty"`
	CPUCountSet       bool                              `json:"-" yaml:"-"`
	CPUPercent        float32                           `json:"cpu_percent,omitempty" yaml:"cpu_percent,omitempty"`
	CPUPercentSet     bool                              `json:"-" yaml:"-"`
	CPUPeriod         int64                             `json:"cpu_period,omitempty" yaml:"cpu_period,omitempty"`
	CPUPeriodSet      bool                              `json:"-" yaml:"-"`
	CPURTPeriod       int64                             `json:"cpu_rt_period,omitempty" yaml:"cpu_rt_period,omitempty"`
	CPURTPeriodSet    bool                              `json:"-" yaml:"-"`
	CPURTRuntime      int64                             `json:"cpu_rt_runtime,omitempty" yaml:"cpu_rt_runtime,omitempty"`
	CPURTRuntimeSet   bool                              `json:"-" yaml:"-"`
	CPUSet            string                            `json:"cpuset,omitempty" yaml:"cpuset,omitempty"`
	DeviceCgroupRules []string                          `json:"device_cgroup_rules,omitempty" yaml:"device_cgroup_rules,omitempty"`
	Gpus              []map[string]interface{}          `json:"gpus,omitempty" yaml:"gpus,omitempty"`
	NetworkMode       string                            `json:"network_mode,omitempty" yaml:"network_mode,omitempty"`
	OomKillDisable    bool                              `json:"oom_kill_disable,omitempty" yaml:"oom_kill_disable,omitempty"`
	OomKillDisableSet bool                              `json:"-" yaml:"-"`
	OomScoreAdj       int64                             `json:"oom_score_adj,omitempty" yaml:"oom_score_adj,omitempty"`
	OomScoreAdjSet    bool                              `json:"-" yaml:"-"`
	Scale             *int                              `json:"scale,omitempty" yaml:"scale,omitempty"`
	Models            map[string]map[string]interface{} `json:"models,omitempty" yaml:"models,omitempty"`
	ExternalLinks     []string                          `json:"external_links,omitempty" yaml:"external_links,omitempty"`
	LabelFiles        []string                          `json:"label_files,omitempty" yaml:"label_files,omitempty"`
	StorageOpt        map[string]string                 `json:"storage_opt,omitempty" yaml:"storage_opt,omitempty"`
	UseAPISocket      bool                              `json:"use_api_socket,omitempty" yaml:"use_api_socket,omitempty"`
	UseAPISocketSet   bool                              `json:"-" yaml:"-"`
	Isolation         string                            `json:"isolation,omitempty" yaml:"isolation,omitempty"`
	Tmpfs             []string                          `json:"tmpfs,omitempty" yaml:"tmpfs,omitempty"`
	Uts               string                            `json:"uts,omitempty" yaml:"uts,omitempty"`
	VolumeDriver      string                            `json:"volume_driver,omitempty" yaml:"volume_driver,omitempty"`
	VolumesFrom       []string                          `json:"volumes_from,omitempty" yaml:"volumes_from,omitempty"`
	Extensions        map[string]interface{}            `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// DevelopConfig captures compose-spec develop/watch intent in a portable form.
// Runtime orchestrators do not execute this directly, but preserving it avoids
// losing local dev-loop behavior when Compose is bridged through other formats.
type DevelopConfig struct {
	Watch      []DevelopWatch         `json:"watch,omitempty" yaml:"watch,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// DevelopWatch captures one compose-spec develop.watch trigger.
type DevelopWatch struct {
	Path        string                 `json:"path,omitempty" yaml:"path,omitempty"`
	Action      string                 `json:"action,omitempty" yaml:"action,omitempty"`
	Target      string                 `json:"target,omitempty" yaml:"target,omitempty"`
	Exec        *ServiceHook           `json:"exec,omitempty" yaml:"exec,omitempty"`
	Include     []string               `json:"include,omitempty" yaml:"include,omitempty"`
	Ignore      []string               `json:"ignore,omitempty" yaml:"ignore,omitempty"`
	InitialSync bool                   `json:"initial_sync,omitempty" yaml:"initial_sync,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// LifecycleHooks captures compose-spec service lifecycle hooks.
type LifecycleHooks struct {
	PreStart   []ServiceHook          `json:"pre_start,omitempty" yaml:"pre_start,omitempty"`
	PostStart  []ServiceHook          `json:"post_start,omitempty" yaml:"post_start,omitempty"`
	PreStop    []ServiceHook          `json:"pre_stop,omitempty" yaml:"pre_stop,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func isEmptyLifecycleHooks(lifecycle *LifecycleHooks) bool {
	return lifecycle == nil ||
		(len(lifecycle.PreStart) == 0 &&
			len(lifecycle.PostStart) == 0 &&
			len(lifecycle.PreStop) == 0 &&
			len(lifecycle.Extensions) == 0)
}

// ServiceHook captures compose-spec service hooks.
type ServiceHook struct {
	Command     []string               `json:"command,omitempty" yaml:"command,omitempty"`
	Image       string                 `json:"image,omitempty" yaml:"image,omitempty"`
	User        string                 `json:"user,omitempty" yaml:"user,omitempty"`
	Privileged  bool                   `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	WorkingDir  string                 `json:"working_dir,omitempty" yaml:"working_dir,omitempty"`
	Environment map[string]*string     `json:"environment,omitempty" yaml:"environment,omitempty"`
	PerReplica  bool                   `json:"per_replica,omitempty" yaml:"per_replica,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Toleration captures a Kubernetes scheduling toleration in portable form.
type Toleration struct {
	Key               string                 `json:"key,omitempty" yaml:"key,omitempty"`
	Operator          string                 `json:"operator,omitempty" yaml:"operator,omitempty"`
	Value             string                 `json:"value,omitempty" yaml:"value,omitempty"`
	Effect            string                 `json:"effect,omitempty" yaml:"effect,omitempty"`
	TolerationSeconds *int64                 `json:"toleration_seconds,omitempty" yaml:"toleration_seconds,omitempty"`
	Extensions        map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func isEmptyDevelopConfig(develop *DevelopConfig) bool {
	return develop == nil ||
		(len(develop.Watch) == 0 && len(develop.Extensions) == 0)
}

func isEmptyResourceSpec(resources *ResourceSpec) bool {
	return resources == nil ||
		(resources.CPUReservation == "" &&
			resources.CPULimit == "" &&
			resources.MemoryReservation == "" &&
			resources.MemoryLimit == "" &&
			resources.EphemeralStorageReservation == "" &&
			resources.EphemeralStorageLimit == "" &&
			!resources.pidsLimitSet &&
			resources.PidsLimit == 0 &&
			!resources.pidsReservationSet &&
			resources.PidsReservation == 0 &&
			len(resources.Devices) == 0 &&
			len(resources.GenericResources) == 0 &&
			len(resources.Extensions) == 0 &&
			len(resources.LimitExtensions) == 0 &&
			len(resources.ReservationExtensions) == 0)
}

func isEmptyRestartPolicy(policy *RestartPolicy) bool {
	return policy == nil ||
		(policy.Condition == "" &&
			policy.Delay == "" &&
			policy.MaxAttempts == 0 &&
			policy.Window == "" &&
			len(policy.Extensions) == 0)
}

func isEmptyMigratePolicy(policy *MigratePolicy) bool {
	return policy == nil ||
		(policy.MaxParallel == 0 &&
			policy.HealthCheck == "" &&
			policy.MinHealthyTime == "" &&
			policy.HealthyDeadline == "" &&
			len(policy.Extensions) == 0)
}

func isEmptyReschedulePolicy(policy *ReschedulePolicy) bool {
	return policy == nil ||
		(policy.Attempts == 0 &&
			policy.Interval == "" &&
			policy.Delay == "" &&
			policy.DelayFunction == "" &&
			policy.MaxDelay == "" &&
			!policy.Unlimited &&
			len(policy.Extensions) == 0)
}

func isEmptyComposeCompat(compat *ComposeCompat) bool {
	return compat == nil ||
		(compat.MacAddress == "" &&
			len(compat.Annotations) == 0 &&
			len(compat.BlkioConfig) == 0 &&
			len(compat.CredentialSpec) == 0 &&
			len(compat.Provider) == 0 &&
			len(compat.Extends) == 0 &&
			compat.Platform == "" &&
			compat.PullRefreshAfter == "" &&
			compat.MemSwappiness == "" &&
			compat.DomainName == "" &&
			compat.CgroupParent == "" &&
			compat.Cgroup == "" &&
			compat.Attach == nil &&
			!(compat.CPUCountSet || compat.CPUCount != 0) &&
			!(compat.CPUPercentSet || compat.CPUPercent != 0) &&
			!(compat.CPUPeriodSet || compat.CPUPeriod != 0) &&
			!(compat.CPURTPeriodSet || compat.CPURTPeriod != 0) &&
			!(compat.CPURTRuntimeSet || compat.CPURTRuntime != 0) &&
			compat.CPUSet == "" &&
			len(compat.DeviceCgroupRules) == 0 &&
			len(compat.Gpus) == 0 &&
			compat.NetworkMode == "" &&
			!(compat.OomKillDisableSet || compat.OomKillDisable) &&
			!(compat.OomScoreAdjSet || compat.OomScoreAdj != 0) &&
			compat.Scale == nil &&
			len(compat.Models) == 0 &&
			len(compat.ExternalLinks) == 0 &&
			len(compat.LabelFiles) == 0 &&
			len(compat.StorageOpt) == 0 &&
			!(compat.UseAPISocketSet || compat.UseAPISocket) &&
			compat.Isolation == "" &&
			len(compat.Extensions) == 0 &&
			len(compat.Tmpfs) == 0 &&
			compat.Uts == "" &&
			compat.VolumeDriver == "" &&
			len(compat.VolumesFrom) == 0)
}

func mergeResourceSpec(base, overlay *ResourceSpec) *ResourceSpec {
	if base == nil && overlay == nil {
		return nil
	}
	result := &ResourceSpec{}
	if base != nil {
		*result = *base
		result.Devices = cloneResourceDevices(base.Devices)
		result.GenericResources = cloneGenericResources(base.GenericResources)
		result.Extensions = copyStringInterfaceMap(base.Extensions)
		result.LimitExtensions = copyStringInterfaceMap(base.LimitExtensions)
		result.ReservationExtensions = copyStringInterfaceMap(base.ReservationExtensions)
	}
	if overlay != nil {
		if overlay.CPUReservation != "" {
			result.CPUReservation = overlay.CPUReservation
		}
		if overlay.CPULimit != "" {
			result.CPULimit = overlay.CPULimit
		}
		if overlay.MemoryReservation != "" {
			result.MemoryReservation = overlay.MemoryReservation
		}
		if overlay.MemoryLimit != "" {
			result.MemoryLimit = overlay.MemoryLimit
		}
		if overlay.pidsLimitSet || overlay.PidsLimit > 0 {
			result.PidsLimit = overlay.PidsLimit
			result.pidsLimitSet = true
		}
		if overlay.pidsReservationSet || overlay.PidsReservation > 0 {
			result.PidsReservation = overlay.PidsReservation
			result.pidsReservationSet = true
		}
		if len(overlay.Devices) > 0 {
			result.Devices = cloneResourceDevices(overlay.Devices)
		}
		if len(overlay.GenericResources) > 0 {
			result.GenericResources = cloneGenericResources(overlay.GenericResources)
		}
		if len(overlay.Extensions) > 0 {
			result.Extensions = copyStringInterfaceMap(overlay.Extensions)
		}
		if len(overlay.LimitExtensions) > 0 {
			result.LimitExtensions = copyStringInterfaceMap(overlay.LimitExtensions)
		}
		if len(overlay.ReservationExtensions) > 0 {
			result.ReservationExtensions = copyStringInterfaceMap(overlay.ReservationExtensions)
		}
	}
	if isEmptyResourceSpec(result) {
		return nil
	}
	return result
}

func serializeResourceSpec(resources *ResourceSpec) map[string]interface{} {
	if resources == nil {
		return nil
	}
	result := map[string]interface{}{}
	limits := map[string]interface{}{}
	reservations := map[string]interface{}{}
	if resources.CPUReservation != "" {
		reservations["cpus"] = resources.CPUReservation
	}
	if resources.CPULimit != "" {
		limits["cpus"] = resources.CPULimit
	}
	if resources.MemoryReservation != "" {
		reservations["memory"] = resources.MemoryReservation
	}
	if resources.MemoryLimit != "" {
		limits["memory"] = resources.MemoryLimit
	}
	if resources.EphemeralStorageReservation != "" {
		reservations["ephemeral-storage"] = resources.EphemeralStorageReservation
	}
	if resources.EphemeralStorageLimit != "" {
		limits["ephemeral-storage"] = resources.EphemeralStorageLimit
	}
	if resources.pidsLimitSet || resources.PidsLimit > 0 {
		limits["pids"] = resources.PidsLimit
	}
	if resources.pidsReservationSet || resources.PidsReservation > 0 {
		reservations["pids"] = resources.PidsReservation
	}
	if len(resources.Devices) > 0 {
		devices := make([]map[string]interface{}, 0, len(resources.Devices))
		for _, device := range resources.Devices {
			if len(device.Capabilities) == 0 && device.Driver == "" && device.Count == "" && len(device.DeviceIDs) == 0 && len(device.Options) == 0 && len(device.Extensions) == 0 {
				continue
			}
			item := map[string]interface{}{}
			if len(device.Capabilities) > 0 {
				item["capabilities"] = append([]string{}, device.Capabilities...)
			}
			if device.Driver != "" {
				item["driver"] = device.Driver
			}
			if device.Count != "" {
				item["count"] = device.Count
			}
			if len(device.DeviceIDs) > 0 {
				item["device_ids"] = append([]string{}, device.DeviceIDs...)
			}
			if len(device.Options) > 0 {
				item["options"] = copyStringMap(device.Options)
			}
			for key, value := range device.Extensions {
				item[key] = deepCopyValue(value)
			}
			devices = append(devices, item)
		}
		if len(devices) > 0 {
			reservations["devices"] = devices
		}
	}
	if len(resources.GenericResources) > 0 {
		genericResources := make([]map[string]interface{}, 0, len(resources.GenericResources))
		for _, resource := range resources.GenericResources {
			if resource.Kind == "" && resource.Value == "" && len(resource.Extensions) == 0 && len(resource.DiscreteExtensions) == 0 {
				continue
			}
			item := map[string]interface{}{}
			for key, value := range resource.Extensions {
				item[key] = deepCopyValue(value)
			}
			discrete := map[string]interface{}{}
			if resource.Kind != "" {
				discrete["kind"] = resource.Kind
			}
			if resource.Value != "" {
				discrete["value"] = resource.Value
			}
			for key, value := range resource.DiscreteExtensions {
				discrete[key] = deepCopyValue(value)
			}
			if len(discrete) > 0 {
				item["discrete_resource_spec"] = discrete
			}
			genericResources = append(genericResources, item)
		}
		if len(genericResources) > 0 {
			reservations["generic_resources"] = genericResources
		}
	}
	if len(resources.LimitExtensions) > 0 {
		for key, value := range resources.LimitExtensions {
			limits[key] = deepCopyValue(value)
		}
	}
	if len(resources.ReservationExtensions) > 0 {
		for key, value := range resources.ReservationExtensions {
			reservations[key] = deepCopyValue(value)
		}
	}
	if len(limits) > 0 {
		result["limits"] = limits
	}
	if len(reservations) > 0 {
		result["reservations"] = reservations
	}
	for key, value := range resources.Extensions {
		result[key] = deepCopyValue(value)
	}
	return result
}

// UpdatePolicy captures rollout behavior in portable terms.
type UpdatePolicy struct {
	Parallelism      int                    `json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
	ParallelismSet   bool                   `json:"parallelism_set,omitempty" yaml:"parallelism_set,omitempty"`
	Delay            string                 `json:"delay,omitempty" yaml:"delay,omitempty"`
	Monitor          string                 `json:"monitor,omitempty" yaml:"monitor,omitempty"`
	MaxFailureRatio  string                 `json:"max_failure_ratio,omitempty" yaml:"max_failure_ratio,omitempty"`
	Order            string                 `json:"order,omitempty" yaml:"order,omitempty"`
	OnFailure        string                 `json:"on_failure,omitempty" yaml:"on_failure,omitempty"`
	HealthCheck      string                 `json:"health_check,omitempty" yaml:"health_check,omitempty"`
	MinHealthyTime   string                 `json:"min_healthy_time,omitempty" yaml:"min_healthy_time,omitempty"`
	HealthyDeadline  string                 `json:"healthy_deadline,omitempty" yaml:"healthy_deadline,omitempty"`
	ProgressDeadline string                 `json:"progress_deadline,omitempty" yaml:"progress_deadline,omitempty"`
	AutoRevert       bool                   `json:"auto_revert,omitempty" yaml:"auto_revert,omitempty"`
	AutoRevertSet    bool                   `json:"-" yaml:"-"`
	AutoPromote      bool                   `json:"auto_promote,omitempty" yaml:"auto_promote,omitempty"`
	AutoPromoteSet   bool                   `json:"-" yaml:"-"`
	Canary           int                    `json:"canary,omitempty" yaml:"canary,omitempty"`
	CanarySet        bool                   `json:"-" yaml:"-"`
	Stagger          string                 `json:"stagger,omitempty" yaml:"stagger,omitempty"`
	Extensions       map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// MigratePolicy captures Nomad-style migration behavior in portable terms.
type MigratePolicy struct {
	MaxParallel     int                    `json:"max_parallel,omitempty" yaml:"max_parallel,omitempty"`
	HealthCheck     string                 `json:"health_check,omitempty" yaml:"health_check,omitempty"`
	MinHealthyTime  string                 `json:"min_healthy_time,omitempty" yaml:"min_healthy_time,omitempty"`
	HealthyDeadline string                 `json:"healthy_deadline,omitempty" yaml:"healthy_deadline,omitempty"`
	Extensions      map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// ReschedulePolicy captures Nomad-style rescheduling behavior in portable terms.
type ReschedulePolicy struct {
	Attempts      int                    `json:"attempts,omitempty" yaml:"attempts,omitempty"`
	Interval      string                 `json:"interval,omitempty" yaml:"interval,omitempty"`
	Delay         string                 `json:"delay,omitempty" yaml:"delay,omitempty"`
	DelayFunction string                 `json:"delay_function,omitempty" yaml:"delay_function,omitempty"`
	MaxDelay      string                 `json:"max_delay,omitempty" yaml:"max_delay,omitempty"`
	Unlimited     bool                   `json:"unlimited,omitempty" yaml:"unlimited,omitempty"`
	Extensions    map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// RestartPolicy captures restart behavior in portable terms.
type RestartPolicy struct {
	Condition   string                 `json:"condition,omitempty" yaml:"condition,omitempty"`
	Delay       string                 `json:"delay,omitempty" yaml:"delay,omitempty"`
	MaxAttempts int                    `json:"max_attempts,omitempty" yaml:"max_attempts,omitempty"`
	Window      string                 `json:"window,omitempty" yaml:"window,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// DependencySpec captures startup ordering and readiness intent across platforms.
type DependencySpec struct {
	Name       string                 `json:"name" yaml:"name"`
	Condition  string                 `json:"condition,omitempty" yaml:"condition,omitempty"` // service_started, service_healthy, service_completed_successfully
	Restart    bool                   `json:"restart,omitempty" yaml:"restart,omitempty"`
	Required   *bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// PortMapping represents a port mapping between host and container
type PortMapping struct {
	Name          string                 `json:"name,omitempty" yaml:"name,omitempty"`
	TargetName    string                 `json:"target_name,omitempty" yaml:"target_name,omitempty"`
	HostIP        string                 `json:"host_ip,omitempty" yaml:"host_ip,omitempty"`
	HostPort      string                 `json:"host_port" yaml:"host_port"`
	ContainerPort string                 `json:"container_port" yaml:"container_port"`
	NodePort      string                 `json:"node_port,omitempty" yaml:"node_port,omitempty"`
	Protocol      string                 `json:"protocol,omitempty" yaml:"protocol,omitempty"` // tcp, udp
	AppProtocol   string                 `json:"app_protocol,omitempty" yaml:"app_protocol,omitempty"`
	Mode          string                 `json:"mode,omitempty" yaml:"mode,omitempty"` // ingress, host
	Extensions    map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// VolumeMount represents a volume mount
type VolumeMount struct {
	Source            string                 `json:"source" yaml:"source"`
	Target            string                 `json:"target" yaml:"target"`
	Type              string                 `json:"type,omitempty" yaml:"type,omitempty"` // bind, volume, tmpfs
	ReadOnly          bool                   `json:"read_only,omitempty" yaml:"read_only,omitempty"`
	Consistency       string                 `json:"consistency,omitempty" yaml:"consistency,omitempty"`
	Mode              string                 `json:"mode,omitempty" yaml:"mode,omitempty"` // Z, z, shared, private, slave
	Propagation       string                 `json:"propagation,omitempty" yaml:"propagation,omitempty"`
	MountPropagation  string                 `json:"mount_propagation,omitempty" yaml:"mount_propagation,omitempty"`
	RecursiveReadOnly string                 `json:"recursive_read_only,omitempty" yaml:"recursive_read_only,omitempty"`
	SubPath           string                 `json:"sub_path,omitempty" yaml:"sub_path,omitempty"`
	SubPathExpr       string                 `json:"sub_path_expr,omitempty" yaml:"sub_path_expr,omitempty"`
	CreateHostPath    *bool                  `json:"create_host_path,omitempty" yaml:"create_host_path,omitempty"`
	NoCopy            bool                   `json:"nocopy,omitempty" yaml:"nocopy,omitempty"`
	VolumeLabels      map[string]string      `json:"volume_labels,omitempty" yaml:"volume_labels,omitempty"`
	TmpfsMode         string                 `json:"tmpfs_mode,omitempty" yaml:"tmpfs_mode,omitempty"`
	ImageSubpath      string                 `json:"image_subpath,omitempty" yaml:"image_subpath,omitempty"`
	Options           map[string]string      `json:"options,omitempty" yaml:"options,omitempty"`
	Extensions        map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	BindExtensions    map[string]interface{} `json:"bind_extensions,omitempty" yaml:"bind_extensions,omitempty"`
	VolumeExtensions  map[string]interface{} `json:"volume_extensions,omitempty" yaml:"volume_extensions,omitempty"`
	TmpfsExtensions   map[string]interface{} `json:"tmpfs_extensions,omitempty" yaml:"tmpfs_extensions,omitempty"`
	ImageExtensions   map[string]interface{} `json:"image_extensions,omitempty" yaml:"image_extensions,omitempty"`
}

// FileRef represents a service-level reference to a config or secret object.
type FileRef struct {
	Source     string                 `json:"source" yaml:"source"`
	Key        string                 `json:"key,omitempty" yaml:"key,omitempty"`
	Target     string                 `json:"target,omitempty" yaml:"target,omitempty"`
	UID        string                 `json:"uid,omitempty" yaml:"uid,omitempty"`
	GID        string                 `json:"gid,omitempty" yaml:"gid,omitempty"`
	Mode       string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	ReadOnly   bool                   `json:"read_only,omitempty" yaml:"read_only,omitempty"`
	Optional   *bool                  `json:"optional,omitempty" yaml:"optional,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func mergePortableFileRefs(existing, preserved []FileRef, appendUnmatched bool) []FileRef {
	if len(existing) == 0 {
		if !appendUnmatched {
			return existing
		}
		return append([]FileRef{}, preserved...)
	}
	if len(preserved) == 0 {
		return existing
	}
	result := append([]FileRef{}, existing...)
	for _, ref := range preserved {
		if idx := matchingPortableFileRef(result, ref); idx >= 0 {
			result[idx] = mergePortableFileRef(result[idx], ref)
			continue
		}
		if appendUnmatched {
			result = append(result, ref)
		}
	}
	return result
}

func matchingPortableFileRef(refs []FileRef, candidate FileRef) int {
	best := -1
	bestScore := -1
	for i, ref := range refs {
		if ref.Source != candidate.Source || ref.Key != candidate.Key {
			continue
		}
		score := 1
		switch {
		case ref.Target != "" && candidate.Target != "" && ref.Target == candidate.Target:
			score = 3
		case ref.Target == "" || candidate.Target == "":
			score = 2
		}
		if score > bestScore {
			best = i
			bestScore = score
		}
	}
	return best
}

func mergePortableFileRef(existing, preserved FileRef) FileRef {
	if existing.Source == "" {
		existing.Source = preserved.Source
	}
	if existing.Key == "" {
		existing.Key = preserved.Key
	}
	if existing.Target == "" {
		existing.Target = preserved.Target
	}
	if existing.UID == "" {
		existing.UID = preserved.UID
	}
	if existing.GID == "" {
		existing.GID = preserved.GID
	}
	if existing.Mode == "" {
		existing.Mode = preserved.Mode
	}
	if !existing.ReadOnly {
		existing.ReadOnly = preserved.ReadOnly
	}
	if existing.Optional == nil && preserved.Optional != nil {
		value := *preserved.Optional
		existing.Optional = &value
	}
	if len(preserved.Extensions) > 0 {
		if existing.Extensions == nil {
			existing.Extensions = map[string]interface{}{}
		}
		for key, value := range preserved.Extensions {
			if _, ok := existing.Extensions[key]; !ok {
				existing.Extensions[key] = value
			}
		}
	}
	return existing
}

// EnvFileRef represents a Compose env_file entry, including long syntax metadata.
type EnvFileRef struct {
	Path       string                 `json:"path" yaml:"path"`
	Required   *bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Format     string                 `json:"format,omitempty" yaml:"format,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// EnvSource represents a single environment variable sourced from a config or secret key.
type EnvSource struct {
	Name       string                 `json:"name" yaml:"name"`
	SourceType string                 `json:"source_type" yaml:"source_type"` // config, secret, field, resource
	Source     string                 `json:"source" yaml:"source"`
	Key        string                 `json:"key,omitempty" yaml:"key,omitempty"`
	Optional   bool                   `json:"optional,omitempty" yaml:"optional,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// EnvFromSource represents importing all environment variables from a config or secret.
type EnvFromSource struct {
	SourceType string                 `json:"source_type" yaml:"source_type"` // config, secret
	Source     string                 `json:"source" yaml:"source"`
	Prefix     string                 `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Optional   bool                   `json:"optional,omitempty" yaml:"optional,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// HealthCheck represents a container health check
type HealthCheck struct {
	Test          []string               `json:"test" yaml:"test"`
	Type          string                 `json:"type,omitempty" yaml:"type,omitempty"` // http, tcp, exec
	Path          string                 `json:"path,omitempty" yaml:"path,omitempty"`
	Port          string                 `json:"port,omitempty" yaml:"port,omitempty"`
	Interval      string                 `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout       string                 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Retries       int                    `json:"retries,omitempty" yaml:"retries,omitempty"`
	StartPeriod   string                 `json:"start_period,omitempty" yaml:"start_period,omitempty"`
	StartInterval string                 `json:"start_interval,omitempty" yaml:"start_interval,omitempty"`
	Disable       bool                   `json:"disable,omitempty" yaml:"disable,omitempty"`
	DisableSet    bool                   `json:"-" yaml:"-"`
	Extensions    map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func (health *HealthCheck) MarshalJSON() ([]byte, error) {
	if health == nil {
		return []byte("null"), nil
	}
	return json.Marshal(serializeHealthCheck(health))
}

func (health *HealthCheck) UnmarshalJSON(data []byte) error {
	if health == nil {
		return fmt.Errorf("nil HealthCheck")
	}
	parsed, err := parseHealthCheckJSON(data)
	if err != nil {
		return err
	}
	*health = *parsed
	return nil
}

// SeccompProfile captures Kubernetes seccomp profile settings.
type SeccompProfile struct {
	Type             string                 `json:"type,omitempty" yaml:"type,omitempty"`
	LocalhostProfile string                 `json:"localhost_profile,omitempty" yaml:"localhost_profile,omitempty"`
	Extensions       map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// SELinuxOptions captures Kubernetes SELinux settings for pod/container security contexts.
type SELinuxOptions struct {
	User       string                 `json:"user,omitempty" yaml:"user,omitempty"`
	Role       string                 `json:"role,omitempty" yaml:"role,omitempty"`
	Type       string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Level      string                 `json:"level,omitempty" yaml:"level,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// WindowsSecurityContextOptions captures Kubernetes Windows security settings
// for pod/container security contexts.
type WindowsSecurityContextOptions struct {
	GMSACredentialSpecName *string                `json:"gmsa_credential_spec_name,omitempty" yaml:"gmsa_credential_spec_name,omitempty"`
	GMSACredentialSpec     *string                `json:"gmsa_credential_spec,omitempty" yaml:"gmsa_credential_spec,omitempty"`
	RunAsUserName          *string                `json:"run_as_user_name,omitempty" yaml:"run_as_user_name,omitempty"`
	HostProcess            *bool                  `json:"host_process,omitempty" yaml:"host_process,omitempty"`
	Extensions             map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// BuildConfig represents container build configuration.
type BuildConfig struct {
	Context    string                 `json:"context,omitempty" yaml:"context,omitempty"`
	Dockerfile string                 `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Ulimits represents container ulimit configuration.
type Ulimits struct {
	Nofile *NofileLimit          `json:"nofile,omitempty" yaml:"nofile,omitempty"`
	Limits map[string]UlimitSpec `json:"limits,omitempty" yaml:"limits,omitempty"`
}

// NofileLimit represents soft/hard nofile limits.
type NofileLimit struct {
	Soft       int                    `json:"soft,omitempty" yaml:"soft,omitempty"`
	Hard       int                    `json:"hard,omitempty" yaml:"hard,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// UlimitSpec captures a compose-spec ulimit entry, including object-form
// extension fields and non-nofile limit names such as nproc.
type UlimitSpec struct {
	Single     int                    `json:"single,omitempty" yaml:"single,omitempty"`
	Soft       int                    `json:"soft,omitempty" yaml:"soft,omitempty"`
	Hard       int                    `json:"hard,omitempty" yaml:"hard,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Config represents a configuration file
type Config struct {
	Name               string                 `json:"name" yaml:"name"`
	PlatformName       string                 `json:"platform_name,omitempty" yaml:"platform_name,omitempty"`
	Content            string                 `json:"content,omitempty" yaml:"content,omitempty"`
	Environment        string                 `json:"environment,omitempty" yaml:"environment,omitempty"`
	File               string                 `json:"file,omitempty" yaml:"file,omitempty"`
	Template           string                 `json:"template,omitempty" yaml:"template,omitempty"`
	Mode               string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	External           bool                   `json:"external,omitempty" yaml:"external,omitempty"`
	ExternalSet        bool                   `json:"-" yaml:"-"`
	ExternalExtensions map[string]interface{} `json:"external_extensions,omitempty" yaml:"external_extensions,omitempty"`
	Labels             map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions         map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Secret represents a secret
type Secret struct {
	Name               string                 `json:"name" yaml:"name"`
	PlatformName       string                 `json:"platform_name,omitempty" yaml:"platform_name,omitempty"`
	File               string                 `json:"file,omitempty" yaml:"file,omitempty"`
	Environment        string                 `json:"environment,omitempty" yaml:"environment,omitempty"`
	Template           string                 `json:"template,omitempty" yaml:"template,omitempty"`
	External           bool                   `json:"external,omitempty" yaml:"external,omitempty"`
	ExternalSet        bool                   `json:"-" yaml:"-"`
	ExternalExtensions map[string]interface{} `json:"external_extensions,omitempty" yaml:"external_extensions,omitempty"`
	Driver             string                 `json:"driver,omitempty" yaml:"driver,omitempty"`
	DriverOpts         map[string]string      `json:"driver_opts,omitempty" yaml:"driver_opts,omitempty"`
	Labels             map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions         map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Network represents a Docker network
type Network struct {
	Name               string                 `json:"name" yaml:"name"`
	PlatformName       string                 `json:"platform_name,omitempty" yaml:"platform_name,omitempty"`
	Driver             string                 `json:"driver,omitempty" yaml:"driver,omitempty"`
	DriverOpts         map[string]string      `json:"driver_opts,omitempty" yaml:"driver_opts,omitempty"`
	Attachable         bool                   `json:"attachable,omitempty" yaml:"attachable,omitempty"`
	AttachableSet      bool                   `json:"-" yaml:"-"`
	External           bool                   `json:"external,omitempty" yaml:"external,omitempty"`
	ExternalSet        bool                   `json:"-" yaml:"-"`
	ExternalExtensions map[string]interface{} `json:"external_extensions,omitempty" yaml:"external_extensions,omitempty"`
	Internal           bool                   `json:"internal,omitempty" yaml:"internal,omitempty"`
	InternalSet        bool                   `json:"-" yaml:"-"`
	EnableIPv4         *bool                  `json:"enable_ipv4,omitempty" yaml:"enable_ipv4,omitempty"`
	EnableIPv6         *bool                  `json:"enable_ipv6,omitempty" yaml:"enable_ipv6,omitempty"`
	IPAM               *IPAMConfig            `json:"ipam,omitempty" yaml:"ipam,omitempty"`
	Labels             map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions         map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// NetworkAttachment captures service-scoped network options from compose-spec.
type NetworkAttachment struct {
	Name          string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Aliases       []string               `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	InterfaceName string                 `json:"interface_name,omitempty" yaml:"interface_name,omitempty"`
	IPv4Address   string                 `json:"ipv4_address,omitempty" yaml:"ipv4_address,omitempty"`
	IPv6Address   string                 `json:"ipv6_address,omitempty" yaml:"ipv6_address,omitempty"`
	LinkLocalIPs  []string               `json:"link_local_ips,omitempty" yaml:"link_local_ips,omitempty"`
	MacAddress    string                 `json:"mac_address,omitempty" yaml:"mac_address,omitempty"`
	DriverOpts    map[string]string      `json:"driver_opts,omitempty" yaml:"driver_opts,omitempty"`
	GWPriority    int                    `json:"gw_priority,omitempty" yaml:"gw_priority,omitempty"`
	Priority      int                    `json:"priority,omitempty" yaml:"priority,omitempty"`
	Extensions    map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// IPAMConfig represents IP address management configuration
type IPAMConfig struct {
	Driver     string                 `json:"driver,omitempty" yaml:"driver,omitempty"`
	Config     []IPAMSubnet           `json:"config,omitempty" yaml:"config,omitempty"`
	Options    map[string]string      `json:"options,omitempty" yaml:"options,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// IPAMSubnet represents an IP subnet configuration
type IPAMSubnet struct {
	Subnet       string                 `json:"subnet,omitempty" yaml:"subnet,omitempty"`
	Gateway      string                 `json:"gateway,omitempty" yaml:"gateway,omitempty"`
	IPRange      string                 `json:"ip_range,omitempty" yaml:"ip_range,omitempty"`
	AuxAddresses map[string]string      `json:"aux_addresses,omitempty" yaml:"aux_addresses,omitempty"`
	Extensions   map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Volume represents a named volume
type Volume struct {
	Name               string                 `json:"name" yaml:"name"`
	PlatformName       string                 `json:"platform_name,omitempty" yaml:"platform_name,omitempty"`
	Driver             string                 `json:"driver,omitempty" yaml:"driver,omitempty"`
	DriverOpts         map[string]string      `json:"driver_opts,omitempty" yaml:"driver_opts,omitempty"`
	External           bool                   `json:"external,omitempty" yaml:"external,omitempty"`
	ExternalSet        bool                   `json:"-" yaml:"-"`
	ExternalExtensions map[string]interface{} `json:"external_extensions,omitempty" yaml:"external_extensions,omitempty"`
	Labels             map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Extensions         map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Application represents a complete application with all its components
type Application struct {
	Name                                      string                                             `json:"name,omitempty" yaml:"name,omitempty"`
	Version                                   string                                             `json:"version,omitempty" yaml:"version,omitempty"`
	Services                                  map[string]*Service                                `json:"services" yaml:"services"`
	Networks                                  map[string]*Network                                `json:"networks,omitempty" yaml:"networks,omitempty"`
	Volumes                                   map[string]*Volume                                 `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Configs                                   map[string]*Config                                 `json:"configs,omitempty" yaml:"configs,omitempty"`
	Secrets                                   map[string]*Secret                                 `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Models                                    map[string]*ComposeModel                           `json:"models,omitempty" yaml:"models,omitempty"`
	KubernetesServices                        map[string]*KubernetesServiceSpec                  `json:"kubernetes_services,omitempty" yaml:"kubernetes_services,omitempty"`
	KubernetesServiceAccounts                 map[string]*KubernetesServiceAccountSpec           `json:"kubernetes_service_accounts,omitempty" yaml:"kubernetes_service_accounts,omitempty"`
	KubernetesHPAs                            map[string]*KubernetesHorizontalPodAutoscalerSpec  `json:"kubernetes_hpas,omitempty" yaml:"kubernetes_hpas,omitempty"`
	KubernetesPDBs                            map[string]*KubernetesPodDisruptionBudgetSpec      `json:"kubernetes_pdbs,omitempty" yaml:"kubernetes_pdbs,omitempty"`
	KubernetesResourceQuotas                  map[string]*KubernetesResourceQuotaSpec            `json:"kubernetes_resource_quotas,omitempty" yaml:"kubernetes_resource_quotas,omitempty"`
	KubernetesLimitRanges                     map[string]*KubernetesLimitRangeSpec               `json:"kubernetes_limit_ranges,omitempty" yaml:"kubernetes_limit_ranges,omitempty"`
	KubernetesStorageClasses                  map[string]*KubernetesStorageClassSpec             `json:"kubernetes_storage_classes,omitempty" yaml:"kubernetes_storage_classes,omitempty"`
	KubernetesIngressClasses                  map[string]*KubernetesIngressClassSpec             `json:"kubernetes_ingress_classes,omitempty" yaml:"kubernetes_ingress_classes,omitempty"`
	KubernetesMutatingWebhookConfigurations   map[string]*KubernetesWebhookConfigurationSpec     `json:"kubernetes_mutating_webhook_configurations,omitempty" yaml:"kubernetes_mutating_webhook_configurations,omitempty"`
	KubernetesValidatingWebhookConfigurations map[string]*KubernetesWebhookConfigurationSpec     `json:"kubernetes_validating_webhook_configurations,omitempty" yaml:"kubernetes_validating_webhook_configurations,omitempty"`
	KubernetesCustomResourceDefinitions       map[string]*KubernetesCustomResourceDefinitionSpec `json:"kubernetes_custom_resource_definitions,omitempty" yaml:"kubernetes_custom_resource_definitions,omitempty"`
	KubernetesPriorityClasses                 map[string]*KubernetesPriorityClassSpec            `json:"kubernetes_priority_classes,omitempty" yaml:"kubernetes_priority_classes,omitempty"`
	KubernetesRuntimeClasses                  map[string]*KubernetesRuntimeClassSpec             `json:"kubernetes_runtime_classes,omitempty" yaml:"kubernetes_runtime_classes,omitempty"`
	KubernetesOpaqueManifests                 map[string]*KubernetesOpaqueManifestSpec           `json:"kubernetes_opaque_manifests,omitempty" yaml:"kubernetes_opaque_manifests,omitempty"`

	// Includes for multi-file compositions
	Includes       []string      `json:"include,omitempty" yaml:"include,omitempty"`
	IncludeEntries []interface{} `json:"include_entries,omitempty" yaml:"include_entries,omitempty"`
	Namespace      string        `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// Mesh captures Headscale/Tailscale intent that should survive
	// cross-orchestrator bridging instead of living only in ad hoc extensions.
	Mesh *MeshSpec `json:"mesh,omitempty" yaml:"mesh,omitempty"`

	// Platform-specific extensions
	Extensions map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`

	// Metadata
	Platform    Platform `json:"platform" yaml:"platform"`
	SourceFiles []string `json:"source_files,omitempty" yaml:"source_files,omitempty"`

	// Canonical is the universal cross-orchestrator view derived from the
	// source-specific model. It is rebuilt by parsers and conversions.
	Canonical *CanonicalApplication `json:"canonical,omitempty" yaml:"canonical,omitempty"`
}

// Validate performs basic validation on the application
func (app *Application) Validate() error {
	if app.Services == nil || len(app.Services) == 0 {
		return fmt.Errorf("application must have at least one service")
	}

	for name, service := range app.Services {
		if service.Name == "" {
			service.Name = name
		}
		if service.Image == "" {
			return fmt.Errorf("service %s must have an image", name)
		}
	}

	return nil
}

// String returns a string representation of the application
func (app *Application) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Application (%s):\n", app.Platform))
	sb.WriteString(fmt.Sprintf("  Services: %d\n", len(app.Services)))
	sb.WriteString(fmt.Sprintf("  Networks: %d\n", len(app.Networks)))
	sb.WriteString(fmt.Sprintf("  Volumes: %d\n", len(app.Volumes)))
	sb.WriteString(fmt.Sprintf("  Configs: %d\n", len(app.Configs)))
	sb.WriteString(fmt.Sprintf("  Secrets: %d\n", len(app.Secrets)))
	return sb.String()
}
