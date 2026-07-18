package paas

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	kubernetesNamespacesExtensionKey                   = "kubernetes.namespaces"
	composeKubernetesNamespacesExtensionKey            = "x-kubernetes-namespaces"
	kubernetesWorkloadsExtensionKey                    = "kubernetes.workloads"
	composeKubernetesWorkloadsExtensionKey             = "x-kubernetes-workloads"
	kubernetesConfigMapsExtensionKey                   = "kubernetes.configMaps"
	composeKubernetesConfigMapsExtensionKey            = "x-kubernetes-config-maps"
	kubernetesSecretsExtensionKey                      = "kubernetes.secretResources"
	composeKubernetesSecretsExtensionKey               = "x-kubernetes-secret-resources"
	kubernetesResourcesExtensionKey                    = "kubernetes.resources"
	composeKubernetesResourcesExtensionKey             = "x-kubernetes-resources"
	kubernetesHPAExtensionKey                          = "kubernetes.hpa"
	composeKubernetesHPAExtensionKey                   = "x-kubernetes-hpa"
	kubernetesHPAsExtensionKey                         = "kubernetes.horizontalPodAutoscalers"
	composeKubernetesHPAsExtensionKey                  = "x-kubernetes-horizontal-pod-autoscalers"
	kubernetesPDBExtensionKey                          = "kubernetes.pdb"
	composeKubernetesPDBExtensionKey                   = "x-kubernetes-pdb"
	kubernetesPDBsExtensionKey                         = "kubernetes.podDisruptionBudgets"
	composeKubernetesPDBsExtensionKey                  = "x-kubernetes-pod-disruption-budgets"
	kubernetesServiceAccountExtensionKey               = "kubernetes.serviceAccount"
	composeKubernetesServiceAccountExtensionKey        = "x-kubernetes-service-account"
	kubernetesServiceAccountsExtensionKey              = "kubernetes.serviceAccounts"
	composeKubernetesServiceAccountsExtensionKey       = "x-kubernetes-service-accounts"
	kubernetesServicesExtensionKey                     = "kubernetes.serviceResources"
	composeKubernetesServicesExtensionKey              = "x-kubernetes-services"
	kubernetesIngressesExtensionKey                    = "kubernetes.ingresses"
	composeKubernetesIngressesExtensionKey             = "x-kubernetes-ingresses"
	kubernetesNetworkPoliciesExtensionKey              = "kubernetes.networkPolicies"
	composeKubernetesNetworkPoliciesExtensionKey       = "x-kubernetes-network-policies"
	kubernetesPersistentVolumesExtensionKey            = "kubernetes.persistentVolumes"
	composeKubernetesPersistentVolumesExtensionKey     = "x-kubernetes-persistent-volumes"
	kubernetesPVCsExtensionKey                         = "kubernetes.persistentVolumeClaims"
	composeKubernetesPVCsExtensionKey                  = "x-kubernetes-persistent-volume-claims"
	kubernetesRBACExtensionKey                         = "kubernetes.rbac"
	composeKubernetesRBACExtensionKey                  = "x-kubernetes-rbac"
	kubernetesRBACResourcesExtensionKey                = "kubernetes.rbacResources"
	composeKubernetesRBACResourcesExtensionKey         = "x-kubernetes-rbac-resources"
	kubernetesResourceQuotasExtensionKey               = "kubernetes.resourceQuotas"
	composeKubernetesResourceQuotasExtensionKey        = "x-kubernetes-resource-quotas"
	kubernetesLimitRangesExtensionKey                  = "kubernetes.limitRanges"
	composeKubernetesLimitRangesExtensionKey           = "x-kubernetes-limit-ranges"
	kubernetesPriorityClassesExtensionKey              = "kubernetes.priorityClasses"
	composeKubernetesPriorityClassesExtensionKey       = "x-kubernetes-priority-classes"
	kubernetesRuntimeClassesExtensionKey               = "kubernetes.runtimeClasses"
	composeKubernetesRuntimeClassesExtensionKey        = "x-kubernetes-runtime-classes"
	kubernetesStorageClassesExtensionKey               = "kubernetes.storageClasses"
	composeKubernetesStorageClassesExtensionKey        = "x-kubernetes-storage-classes"
	kubernetesIngressClassesExtensionKey               = "kubernetes.ingressClasses"
	composeKubernetesIngressClassesExtensionKey        = "x-kubernetes-ingress-classes"
	kubernetesMutatingWebhooksExtensionKey             = "kubernetes.mutatingWebhookConfigurations"
	composeKubernetesMutatingWebhooksExtensionKey      = "x-kubernetes-mutating-webhook-configurations"
	kubernetesValidatingWebhooksExtensionKey           = "kubernetes.validatingWebhookConfigurations"
	composeKubernetesValidatingWebhooksExtensionKey    = "x-kubernetes-validating-webhook-configurations"
	kubernetesCRDsExtensionKey                         = "kubernetes.customResourceDefinitions"
	composeKubernetesCRDsExtensionKey                  = "x-kubernetes-custom-resource-definitions"
	kubernetesCustomResourcesExtensionKey              = "kubernetes.customResources"
	composeKubernetesCustomResourcesExtensionKey       = "x-kubernetes-custom-resources"
	kubernetesSourceDocumentsExtensionKey              = "kubernetes.sourceDocuments"
	kubernetesAnnotationDeployMode                     = "bolabaden.dev/deploy-mode"
	kubernetesAnnotationDeploySpec                     = "bolabaden.dev/deploy-spec"
	kubernetesAnnotationDeployEndpointMode             = "bolabaden.dev/deploy-endpoint-mode"
	kubernetesAnnotationPortableServiceExtensions      = "bolabaden.dev/portable-service-extensions"
	kubernetesAnnotationDeployLabels                   = "bolabaden.dev/deploy-labels"
	kubernetesAnnotationDeployResources                = "bolabaden.dev/deploy-resources"
	kubernetesAnnotationDeployPlacement                = "bolabaden.dev/deploy-placement-constraints"
	kubernetesAnnotationDeployPreferences              = "bolabaden.dev/deploy-placement-preferences"
	kubernetesAnnotationDeployMaxReplicasPerNode       = "bolabaden.dev/deploy-max-replicas-per-node"
	kubernetesAnnotationDeployUpdateParallelism        = "bolabaden.dev/deploy-update-parallelism"
	kubernetesAnnotationDeployUpdateDelay              = "bolabaden.dev/deploy-update-delay"
	kubernetesAnnotationDeployUpdateMonitor            = "bolabaden.dev/deploy-update-monitor"
	kubernetesAnnotationDeployUpdateFailureRate        = "bolabaden.dev/deploy-update-max-failure-ratio"
	kubernetesAnnotationDeployUpdateOrder              = "bolabaden.dev/deploy-update-order"
	kubernetesAnnotationDeployUpdateOnFailure          = "bolabaden.dev/deploy-update-failure-action"
	kubernetesAnnotationDeployUpdateHealthCheck        = "bolabaden.dev/deploy-update-health-check"
	kubernetesAnnotationDeployUpdateMinHealthyTime     = "bolabaden.dev/deploy-update-min-healthy-time"
	kubernetesAnnotationDeployUpdateHealthyDeadline    = "bolabaden.dev/deploy-update-healthy-deadline"
	kubernetesAnnotationDeployUpdateProgressDeadline   = "bolabaden.dev/deploy-update-progress-deadline"
	kubernetesAnnotationDeployUpdateAutoRevert         = "bolabaden.dev/deploy-update-auto-revert"
	kubernetesAnnotationDeployUpdateAutoPromote        = "bolabaden.dev/deploy-update-auto-promote"
	kubernetesAnnotationDeployUpdateCanary             = "bolabaden.dev/deploy-update-canary"
	kubernetesAnnotationDeployUpdateStagger            = "bolabaden.dev/deploy-update-stagger"
	kubernetesAnnotationDeployRollbackParallelism      = "bolabaden.dev/deploy-rollback-parallelism"
	kubernetesAnnotationDeployRollbackDelay            = "bolabaden.dev/deploy-rollback-delay"
	kubernetesAnnotationDeployRollbackMonitor          = "bolabaden.dev/deploy-rollback-monitor"
	kubernetesAnnotationDeployRollbackFailureRate      = "bolabaden.dev/deploy-rollback-max-failure-ratio"
	kubernetesAnnotationDeployRollbackOrder            = "bolabaden.dev/deploy-rollback-order"
	kubernetesAnnotationDeployRollbackOnFailure        = "bolabaden.dev/deploy-rollback-failure-action"
	kubernetesAnnotationDeployRollbackHealthCheck      = "bolabaden.dev/deploy-rollback-health-check"
	kubernetesAnnotationDeployRollbackMinHealthyTime   = "bolabaden.dev/deploy-rollback-min-healthy-time"
	kubernetesAnnotationDeployRollbackHealthyDeadline  = "bolabaden.dev/deploy-rollback-healthy-deadline"
	kubernetesAnnotationDeployRollbackProgressDeadline = "bolabaden.dev/deploy-rollback-progress-deadline"
	kubernetesAnnotationDeployRollbackAutoRevert       = "bolabaden.dev/deploy-rollback-auto-revert"
	kubernetesAnnotationDeployRollbackAutoPromote      = "bolabaden.dev/deploy-rollback-auto-promote"
	kubernetesAnnotationDeployRollbackCanary           = "bolabaden.dev/deploy-rollback-canary"
	kubernetesAnnotationDeployRollbackStagger          = "bolabaden.dev/deploy-rollback-stagger"
	kubernetesAnnotationDeployRestartCondition         = "bolabaden.dev/deploy-restart-condition"
	kubernetesAnnotationDeployRestartDelay             = "bolabaden.dev/deploy-restart-delay"
	kubernetesAnnotationDeployRestartAttempts          = "bolabaden.dev/deploy-restart-max-attempts"
	kubernetesAnnotationDeployRestartWindow            = "bolabaden.dev/deploy-restart-window"
	kubernetesAnnotationDeployJob                      = "bolabaden.dev/deploy-job"
	kubernetesAnnotationDependencies                   = "bolabaden.dev/dependencies"
	kubernetesAnnotationRuntimeSecurityOpt             = "bolabaden.dev/runtime-security-opt"
	kubernetesAnnotationRuntimeInit                    = "bolabaden.dev/runtime-init"
	kubernetesAnnotationRuntimeStopSignal              = "bolabaden.dev/runtime-stop-signal"
	kubernetesAnnotationPortableBuild                  = "bolabaden.dev/portable-build"
	kubernetesAnnotationPortableDevices                = "bolabaden.dev/portable-devices"
	kubernetesAnnotationPortableDeviceMappings         = "bolabaden.dev/portable-device-mappings"
	kubernetesAnnotationPortableExpose                 = "bolabaden.dev/portable-expose"
	kubernetesAnnotationPortablePorts                  = "bolabaden.dev/portable-ports"
	kubernetesAnnotationPortableEnvFiles               = "bolabaden.dev/portable-env-files"
	kubernetesAnnotationPortableHealthcheck            = "bolabaden.dev/portable-healthcheck"
	kubernetesAnnotationPortableStartupProbe           = "bolabaden.dev/portable-startup-probe"
	kubernetesAnnotationPortableTolerations            = "bolabaden.dev/portable-tolerations"
	kubernetesAnnotationPortableDevelop                = "bolabaden.dev/portable-develop"
	kubernetesAnnotationPortableLifecycle              = "bolabaden.dev/portable-lifecycle"
	kubernetesAnnotationPortableGroupAdd               = "bolabaden.dev/portable-group-add"
	kubernetesAnnotationPortableRuntime                = "bolabaden.dev/portable-runtime"
	kubernetesAnnotationPortableLogging                = "bolabaden.dev/portable-logging"
	kubernetesAnnotationPortableComposeCompat          = "bolabaden.dev/portable-compose-compat"
	kubernetesAnnotationPortableLinks                  = "bolabaden.dev/portable-links"
	kubernetesAnnotationPortableMesh                   = "bolabaden.dev/portable-mesh"
	kubernetesAnnotationPortableComposeName            = "bolabaden.dev/portable-compose-name"
	kubernetesAnnotationPortableAppExtensions          = "bolabaden.dev/portable-app-extensions"
	kubernetesAnnotationPortableModels                 = "bolabaden.dev/portable-models"
	kubernetesAnnotationPortableIncludes               = "bolabaden.dev/portable-compose-includes"
	kubernetesAnnotationPortablePIDMode                = "bolabaden.dev/portable-pid-mode"
	kubernetesAnnotationPortableIPCMode                = "bolabaden.dev/portable-ipc-mode"
	kubernetesAnnotationPortablePidsLimit              = "bolabaden.dev/portable-pids-limit"
	kubernetesAnnotationPortableShmSize                = "bolabaden.dev/portable-shm-size"
	kubernetesAnnotationPortableCPUShares              = "bolabaden.dev/portable-cpu-shares"
	kubernetesAnnotationPortableCPUQuota               = "bolabaden.dev/portable-cpu-quota"
	kubernetesAnnotationPortableMemLimit               = "bolabaden.dev/portable-mem-limit"
	kubernetesAnnotationPortableMemorySwap             = "bolabaden.dev/portable-memory-swap"
	kubernetesAnnotationPortableMemReservation         = "bolabaden.dev/portable-mem-reservation"
	kubernetesAnnotationPortableCPUs                   = "bolabaden.dev/portable-cpus"
	kubernetesAnnotationPortableUlimits                = "bolabaden.dev/portable-ulimits"
	kubernetesAnnotationPortableUserNSMode             = "bolabaden.dev/portable-userns-mode"
	kubernetesAnnotationPortablePullPolicy             = "bolabaden.dev/portable-pull-policy"
	kubernetesAnnotationPortableProfiles               = "bolabaden.dev/portable-profiles"
	kubernetesAnnotationPortableConfigs                = "bolabaden.dev/portable-config-refs"
	kubernetesAnnotationPortableSecrets                = "bolabaden.dev/portable-secret-refs"
	kubernetesAnnotationPortableVolumes                = "bolabaden.dev/portable-volumes"
	kubernetesAnnotationPortableNetworks               = "bolabaden.dev/portable-networks"
	kubernetesAnnotationPortableAppVolumes             = "bolabaden.dev/portable-app-volumes"
	kubernetesAnnotationPortableAppConfigs             = "bolabaden.dev/portable-app-configs"
	kubernetesAnnotationPortableAppSecrets             = "bolabaden.dev/portable-app-secrets"
	kubernetesAnnotationPortableNetworkAttachments     = "bolabaden.dev/portable-network-attachments"
	kubernetesAnnotationPortableFailover               = "bolabaden.dev/portable-failover"
	kubernetesAnnotationPortableNomadSpread            = "bolabaden.dev/portable-nomad-spread"
	kubernetesAnnotationPortableNomadConnect           = "bolabaden.dev/portable-nomad-connect"
	kubernetesAnnotationPortableNomadRestart           = "bolabaden.dev/portable-nomad-restart"
	kubernetesAnnotationPortableNomadUpdate            = "bolabaden.dev/portable-nomad-update"
	kubernetesAnnotationPortableNomadMigrate           = "bolabaden.dev/portable-nomad-migrate"
	kubernetesAnnotationPortableNomadReschedule        = "bolabaden.dev/portable-nomad-reschedule"
)

// ParseKubernetesYAML parses Kubernetes YAML manifests
func ParseKubernetesYAML(content string) (*Application, error) {
	app := &Application{
		Platform:   PlatformKubernetes,
		Services:   make(map[string]*Service),
		Networks:   make(map[string]*Network),
		Volumes:    make(map[string]*Volume),
		Configs:    make(map[string]*Config),
		Secrets:    make(map[string]*Secret),
		Models:     make(map[string]*ComposeModel),
		Extensions: make(map[string]interface{}),
	}
	if docs := splitKubernetesSourceDocuments(content); len(docs) > 0 {
		app.Extensions[kubernetesSourceDocumentsExtensionKey] = docs
	}

	decoder := k8syaml.NewYAMLOrJSONDecoder(strings.NewReader(content), 4096)
	for {
		var object unstructured.Unstructured
		err := decoder.Decode(&object)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to decode Kubernetes manifest: %w", err)
		}
		if err := parseKubernetesUnstructuredObject(app, &object); err != nil {
			return nil, err
		}
	}

	reconcileKubernetesServices(app)
	reconcileKubernetesHorizontalPodAutoscalers(app)
	reconcileKubernetesPodDisruptionBudgets(app)
	reconcileKubernetesServiceAccounts(app)
	reconcileKubernetesRBACResources(app)
	rehydrateComposeApplicationExtensions(app)
	syncPortableApplicationState(app)
	app.AttachCanonical()
	return app, nil
}

func parseKubernetesUnstructuredObject(app *Application, object *unstructured.Unstructured) error {
	if object == nil {
		return nil
	}
	if object.IsList() {
		list, err := object.ToList()
		if err != nil {
			return fmt.Errorf("failed to expand Kubernetes list: %w", err)
		}
		for index := range list.Items {
			item := &list.Items[index]
			if err := parseKubernetesUnstructuredObject(app, item); err != nil {
				return err
			}
		}
		return nil
	}
	resource := kruntime.DeepCopyJSON(object.UnstructuredContent())
	if len(resource) == 0 {
		return nil
	}

	appendKubernetesRawManifest(app, resource)
	storeKubernetesOpaqueManifest(app, resource)
	kubernetesInferNamespace(app, resource)
	if err := parseKubernetesResource(app, resource); err != nil {
		return fmt.Errorf("failed to parse Kubernetes resource: %w", err)
	}
	return nil
}

func kubernetesInferNamespace(app *Application, resource map[string]interface{}) {
	if app == nil || app.Namespace != "" || len(resource) == 0 {
		return
	}
	metadata, _ := asMap(resource["metadata"])
	namespace := toString(metadata["namespace"])
	if namespace == "" {
		return
	}
	app.Namespace = namespace
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	app.Extensions["kubernetes.namespace"] = namespace
}

func appendKubernetesRawManifest(app *Application, resource map[string]interface{}) {
	if app == nil || len(resource) == 0 {
		return
	}
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	app.Extensions["kubernetes.raw"] = appendExtensionSlice(app.Extensions["kubernetes.raw"], kruntime.DeepCopyJSON(resource))
}

func parseKubernetesResource(app *Application, resource map[string]interface{}) error {
	kind, _ := resource["kind"].(string)

	switch kind {
	case "Pod":
		return parseKubernetesPod(app, resource)
	case "Deployment":
		return parseKubernetesDeployment(app, resource)
	case "StatefulSet", "DaemonSet", "ReplicaSet", "ReplicationController", "Job", "CronJob":
		return parseKubernetesPodController(app, resource, kind)
	case "Service":
		return parseKubernetesService(app, resource)
	case "ConfigMap":
		return parseKubernetesConfigMap(app, resource)
	case "Secret":
		return parseKubernetesSecret(app, resource)
	case "ServiceAccount":
		return parseKubernetesServiceAccount(app, resource)
	case "Namespace":
		return parseKubernetesNamespace(app, resource)
	case "PersistentVolume":
		return parseKubernetesPersistentVolume(app, resource)
	case "PersistentVolumeClaim":
		return parseKubernetesPVC(app, resource)
	case "Ingress":
		return parseKubernetesIngress(app, resource)
	case "HTTPRoute":
		return parseKubernetesHTTPRoute(app, resource)
	case "NetworkPolicy":
		return parseKubernetesPolicy(app, resource)
	case "HorizontalPodAutoscaler":
		return parseKubernetesHorizontalPodAutoscaler(app, resource)
	case "PodDisruptionBudget":
		return parseKubernetesPodDisruptionBudget(app, resource)
	case "Role", "RoleBinding", "ClusterRole", "ClusterRoleBinding":
		return parseKubernetesRBACResource(app, resource)
	case "ResourceQuota":
		return parseKubernetesResourceQuota(app, resource)
	case "LimitRange":
		return parseKubernetesLimitRange(app, resource)
	case "PriorityClass":
		return parseKubernetesPriorityClass(app, resource)
	case "RuntimeClass":
		return parseKubernetesRuntimeClass(app, resource)
	case "StorageClass":
		return parseKubernetesStorageClass(app, resource)
	case "IngressClass":
		return parseKubernetesIngressClass(app, resource)
	case "MutatingWebhookConfiguration":
		return parseKubernetesMutatingWebhookConfiguration(app, resource)
	case "ValidatingWebhookConfiguration":
		return parseKubernetesValidatingWebhookConfiguration(app, resource)
	case "CustomResourceDefinition":
		return parseKubernetesCustomResourceDefinition(app, resource)
	}

	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[composeKubernetesResourcesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesResourcesExtensionKey], deepCopyValue(resource))
	if isKubernetesCustomResource(resource) {
		app.Extensions[kubernetesCustomResourcesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesCustomResourcesExtensionKey], deepCopyValue(resource))
		app.Extensions[composeKubernetesCustomResourcesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesCustomResourcesExtensionKey], deepCopyValue(resource))
	}
	return nil
}

func isKubernetesCustomResource(resource map[string]interface{}) bool {
	apiVersion := toString(resource["apiVersion"])
	kind := toString(resource["kind"])
	return apiVersion != "" && kind != "" && strings.Contains(apiVersion, "/")
}

func parseKubernetesHorizontalPodAutoscaler(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesHPAs == nil {
		app.KubernetesHPAs = map[string]*KubernetesHorizontalPodAutoscalerSpec{}
	}
	if spec := kubernetesHPASpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesHPAs[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions["kubernetes.hpas"] = appendExtensionSlice(app.Extensions["kubernetes.hpas"], kruntime.DeepCopyJSON(resource))
	app.Extensions[kubernetesHPAsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesHPAsExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesHPAsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesHPAsExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesPodDisruptionBudget(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesPDBs == nil {
		app.KubernetesPDBs = map[string]*KubernetesPodDisruptionBudgetSpec{}
	}
	if spec := kubernetesPDBSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesPDBs[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions["kubernetes.pdbs"] = appendExtensionSlice(app.Extensions["kubernetes.pdbs"], deepCopyValue(resource))
	app.Extensions[kubernetesPDBsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesPDBsExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesPDBsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesPDBsExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesServiceAccount(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesServiceAccounts == nil {
		app.KubernetesServiceAccounts = map[string]*KubernetesServiceAccountSpec{}
	}
	if spec := kubernetesServiceAccountSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesServiceAccounts[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesServiceAccountsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesServiceAccountsExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesServiceAccountsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesServiceAccountsExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesRBACResource(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	storeKubernetesOpaqueManifest(app, resource)
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesRBACResourcesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesRBACResourcesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesRBACResourcesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesRBACResourcesExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesResourceQuota(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesResourceQuotas == nil {
		app.KubernetesResourceQuotas = map[string]*KubernetesResourceQuotaSpec{}
	}
	if spec := kubernetesResourceQuotaSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesResourceQuotas[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesResourceQuotasExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesResourceQuotasExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesResourceQuotasExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesResourceQuotasExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesLimitRange(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesLimitRanges == nil {
		app.KubernetesLimitRanges = map[string]*KubernetesLimitRangeSpec{}
	}
	if spec := kubernetesLimitRangeSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesLimitRanges[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesLimitRangesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesLimitRangesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesLimitRangesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesLimitRangesExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesPriorityClass(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesPriorityClasses == nil {
		app.KubernetesPriorityClasses = map[string]*KubernetesPriorityClassSpec{}
	}
	if spec := kubernetesPriorityClassSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesPriorityClasses[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesPriorityClassesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesPriorityClassesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesPriorityClassesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesPriorityClassesExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesRuntimeClass(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesRuntimeClasses == nil {
		app.KubernetesRuntimeClasses = map[string]*KubernetesRuntimeClassSpec{}
	}
	if spec := kubernetesRuntimeClassSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesRuntimeClasses[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesRuntimeClassesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesRuntimeClassesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesRuntimeClassesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesRuntimeClassesExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesStorageClass(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesStorageClasses == nil {
		app.KubernetesStorageClasses = map[string]*KubernetesStorageClassSpec{}
	}
	if spec := kubernetesStorageClassSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesStorageClasses[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesStorageClassesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesStorageClassesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesStorageClassesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesStorageClassesExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesIngressClass(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesIngressClasses == nil {
		app.KubernetesIngressClasses = map[string]*KubernetesIngressClassSpec{}
	}
	if spec := kubernetesIngressClassSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesIngressClasses[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesIngressClassesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesIngressClassesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesIngressClassesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesIngressClassesExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesMutatingWebhookConfiguration(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesMutatingWebhookConfigurations == nil {
		app.KubernetesMutatingWebhookConfigurations = map[string]*KubernetesWebhookConfigurationSpec{}
	}
	if spec := kubernetesWebhookConfigurationSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesMutatingWebhookConfigurations[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesMutatingWebhooksExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesMutatingWebhooksExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesMutatingWebhooksExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesMutatingWebhooksExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesValidatingWebhookConfiguration(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesValidatingWebhookConfigurations == nil {
		app.KubernetesValidatingWebhookConfigurations = map[string]*KubernetesWebhookConfigurationSpec{}
	}
	if spec := kubernetesWebhookConfigurationSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesValidatingWebhookConfigurations[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesValidatingWebhooksExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesValidatingWebhooksExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesValidatingWebhooksExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesValidatingWebhooksExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesCustomResourceDefinition(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	if app.KubernetesCustomResourceDefinitions == nil {
		app.KubernetesCustomResourceDefinitions = map[string]*KubernetesCustomResourceDefinitionSpec{}
	}
	if spec := kubernetesCustomResourceDefinitionSpecFromMap(resource); spec != nil && spec.Name != "" {
		app.KubernetesCustomResourceDefinitions[spec.Name] = spec
	}
	app.Extensions["kubernetes.resources"] = appendExtensionSlice(app.Extensions["kubernetes.resources"], resource)
	app.Extensions[kubernetesCRDsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesCRDsExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesCRDsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesCRDsExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesDeployment(app *Application, resource map[string]interface{}) error {
	return parseKubernetesPodController(app, resource, "Deployment")
}

func parseKubernetesPod(app *Application, resource map[string]interface{}) error {
	return parseKubernetesPodController(app, resource, "Pod")
}

func parseKubernetesPodController(app *Application, resource map[string]interface{}, kind string) error {
	if app != nil && app.Extensions != nil {
		app.Extensions[kubernetesWorkloadsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesWorkloadsExtensionKey], deepCopyValue(resource))
		app.Extensions[composeKubernetesWorkloadsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesWorkloadsExtensionKey], deepCopyValue(resource))
	}
	metadata, _ := asMap(resource["metadata"])
	spec, _ := asMap(resource["spec"])

	name, _ := metadata["name"].(string)
	if name == "" {
		return fmt.Errorf("%s missing name", strings.ToLower(kind))
	}

	var podMetadata, podSpec map[string]interface{}
	if kind == "Pod" {
		podMetadata = metadata
		podSpec = spec
	} else {
		template, ok := workloadTemplate(kind, spec)
		if !ok {
			return fmt.Errorf("%s %s missing pod template", strings.ToLower(kind), name)
		}
		podMetadata, _ = asMap(template["metadata"])
		podSpec, _ = asMap(template["spec"])
	}
	containers, _ := podSpec["containers"].([]interface{})

	if len(containers) == 0 {
		return fmt.Errorf("%s %s has no containers", strings.ToLower(kind), name)
	}

	// Use first container as main service
	container, ok := asMap(containers[0])
	if !ok {
		return fmt.Errorf("%s %s first container is not a map", strings.ToLower(kind), name)
	}

	service := &Service{
		Name:        name,
		Platform:    PlatformKubernetes,
		Environment: make(map[string]string),
		Labels:      make(map[string]string),
		Extensions:  make(map[string]interface{}),
	}

	// Parse container spec
	if image, ok := container["image"].(string); ok {
		service.Image = image
	}
	service.Replicas = workloadReplicas(kind, spec)
	if kind == "Pod" && service.Replicas == 0 {
		service.Replicas = 1
	}
	if service.Replicas > 0 {
		service.Deploy = &DeploySpec{Replicas: service.Replicas}
	}
	for key, value := range toStringMapLoose(metadata["labels"]) {
		service.Labels[key] = value
	}
	if labels := toStringMapLoose(metadata["labels"]); len(labels) > 0 {
		service.Extensions["kubernetes.workload.labels"] = copyStringMap(labels)
		service.Extensions["x-kubernetes-workload-labels"] = copyStringMap(labels)
	}
	for key, value := range toStringMapLoose(podMetadata["labels"]) {
		service.Labels[key] = value
	}
	if labels := toStringMapLoose(podMetadata["labels"]); len(labels) > 0 {
		service.Extensions["kubernetes.template.labels"] = copyStringMap(labels)
		service.Extensions["x-kubernetes-template-labels"] = copyStringMap(labels)
	}
	workloadAnnotations := toStringMapLoose(metadata["annotations"])
	templateAnnotations := toStringMapLoose(podMetadata["annotations"])
	if len(workloadAnnotations) > 0 {
		service.Extensions["kubernetes.workload.annotations"] = copyStringMap(workloadAnnotations)
		service.Extensions["x-kubernetes-workload-annotations"] = copyStringMap(workloadAnnotations)
	}
	if len(templateAnnotations) > 0 {
		service.Extensions["kubernetes.template.annotations"] = copyStringMap(templateAnnotations)
		service.Extensions["x-kubernetes-template-annotations"] = copyStringMap(templateAnnotations)
	}
	annotations := mergeStringMaps(workloadAnnotations, templateAnnotations)
	if len(annotations) > 0 {
		service.Extensions["kubernetes.annotations"] = annotations
	}
	service.Extensions["kubernetes.kind"] = kind
	parseKubernetesWorkloadSpecExtensions(service, kind, spec)
	if kind == "StatefulSet" {
		if serviceName := toString(spec["serviceName"]); serviceName != "" {
			service.Extensions["kubernetes.statefulSet.serviceName"] = serviceName
		}
		if podManagementPolicy := toString(spec["podManagementPolicy"]); podManagementPolicy != "" {
			service.Extensions["kubernetes.statefulSet.podManagementPolicy"] = podManagementPolicy
		}
		if ordinals, ok := asMap(spec["ordinals"]); ok && len(ordinals) > 0 {
			service.Extensions["kubernetes.statefulSet.ordinals"] = cloneMap(ordinals)
		}
		if retentionPolicy, ok := asMap(spec["persistentVolumeClaimRetentionPolicy"]); ok && len(retentionPolicy) > 0 {
			service.Extensions["kubernetes.statefulSet.persistentVolumeClaimRetentionPolicy"] = cloneMap(retentionPolicy)
		}
		if templates, ok := spec["volumeClaimTemplates"].([]interface{}); ok && len(templates) > 0 {
			claimTemplates := make([]map[string]interface{}, 0, len(templates))
			for _, template := range templates {
				if templateMap, ok := asMap(template); ok {
					claimTemplates = append(claimTemplates, cloneMap(templateMap))
				}
			}
			if len(claimTemplates) > 0 {
				service.Extensions["kubernetes.statefulSet.volumeClaimTemplates"] = cloneMapSlice(claimTemplates)
			}
		}
	}
	if kind == "CronJob" {
		if schedule := toString(spec["schedule"]); schedule != "" {
			service.Extensions["kubernetes.cron.schedule"] = schedule
		}
		if concurrencyPolicy := toString(spec["concurrencyPolicy"]); concurrencyPolicy != "" {
			service.Extensions["kubernetes.cron.concurrencyPolicy"] = concurrencyPolicy
		}
		if suspend, ok := spec["suspend"].(bool); ok {
			service.Extensions["kubernetes.cron.suspend"] = fmt.Sprintf("%t", suspend)
		}
		parseKubernetesCronJobSpecExtensions(service, spec)
	}
	if kind == "Job" {
		parseKubernetesJobSpecExtensions(service, spec)
	}
	parseKubernetesDependencies(service, annotations, podSpec)
	parseKubernetesRuntimeAnnotations(service, annotations)
	parseKubernetesLifecycle(service, container)
	parseKubernetesPortableAnnotations(app, service, annotations)

	if command, ok := container["command"].([]interface{}); ok {
		service.Entrypoint = interfaceSliceToStringSlice(command)
	}

	if args, ok := container["args"].([]interface{}); ok {
		service.Command = interfaceSliceToStringSlice(args)
	}
	if workingDir, ok := container["workingDir"].(string); ok {
		service.WorkingDir = workingDir
	}
	if stdin, ok := container["stdin"].(bool); ok {
		service.StdinOpen = stdin
		service.StdinOpenSet = true
	}
	if tty, ok := container["tty"].(bool); ok {
		service.Tty = tty
		service.TtySet = true
	}
	if imagePullPolicy := toString(container["imagePullPolicy"]); imagePullPolicy != "" {
		service.ImagePullPolicy = imagePullPolicy
	}
	if path := toString(container["terminationMessagePath"]); path != "" {
		service.TerminationMessagePath = path
	}
	if policy := toString(container["terminationMessagePolicy"]); policy != "" {
		service.TerminationMessagePolicy = policy
	}
	parseKubernetesRuntime(service, podSpec, container)
	if seconds := toInt(podSpec["terminationGracePeriodSeconds"]); seconds > 0 {
		service.StopGracePeriod = fmt.Sprintf("%ds", seconds)
	}
	parseKubernetesDNS(service, podSpec)

	if env, ok := container["env"].([]interface{}); ok {
		for _, envVar := range env {
			if envMap, ok := asMap(envVar); ok {
				if name, ok := envMap["name"].(string); ok {
					if value, ok := envMap["value"].(string); ok {
						service.Environment[name] = value
					} else if valueFrom, ok := asMap(envMap["valueFrom"]); ok {
						if source := kubernetesEnvSource(name, valueFrom); source.Name != "" {
							if source.Extensions == nil {
								source.Extensions = map[string]interface{}{}
							}
							source.Extensions["kubernetes.valueFrom"] = deepCopyValue(valueFrom)
							for key, value := range envMap {
								switch key {
								case "name", "value", "valueFrom":
									continue
								default:
									source.Extensions[key] = deepCopyValue(value)
								}
							}
							if len(source.Extensions) == 0 {
								source.Extensions = nil
							}
							service.EnvSources = append(service.EnvSources, source)
						}
					}
				}
			}
		}
	}

	if envFrom, ok := container["envFrom"].([]interface{}); ok {
		for _, envFromItem := range envFrom {
			if envFromMap, ok := asMap(envFromItem); ok {
				if source := kubernetesEnvFromSource(envFromMap); source.Source != "" {
					if source.Extensions == nil {
						source.Extensions = map[string]interface{}{}
					}
					for key, value := range envFromMap {
						switch key {
						case "prefix", "configMapRef", "secretRef":
							continue
						default:
							source.Extensions[key] = deepCopyValue(value)
						}
					}
					if len(source.Extensions) == 0 {
						source.Extensions = nil
					}
					service.EnvFrom = append(service.EnvFrom, source)
				}
			}
		}
	}
	if imagePullSecrets, ok := podSpec["imagePullSecrets"].([]interface{}); ok {
		for _, secret := range imagePullSecrets {
			if secretMap, ok := asMap(secret); ok {
				if name, ok := secretMap["name"].(string); ok && name != "" {
					service.ImagePullSecrets = append(service.ImagePullSecrets, name)
				}
			}
		}
	}
	if tolerations, ok := podSpec["tolerations"].([]interface{}); ok {
		for _, tolerationValue := range tolerations {
			if tolerationMap, ok := asMap(tolerationValue); ok {
				toleration := Toleration{
					Key:        toString(tolerationMap["key"]),
					Operator:   toString(tolerationMap["operator"]),
					Value:      toString(tolerationMap["value"]),
					Effect:     toString(tolerationMap["effect"]),
					Extensions: map[string]interface{}{},
				}
				switch v := tolerationMap["tolerationSeconds"].(type) {
				case int64:
					toleration.TolerationSeconds = &v
				case int:
					value := int64(v)
					toleration.TolerationSeconds = &value
				case float64:
					value := int64(v)
					toleration.TolerationSeconds = &value
				}
				for key, value := range tolerationMap {
					switch key {
					case "key", "operator", "value", "effect", "tolerationSeconds":
					default:
						toleration.Extensions[key] = deepCopyValue(value)
					}
				}
				if len(toleration.Extensions) == 0 {
					toleration.Extensions = nil
				}
				if toleration.Key != "" || toleration.Operator != "" || toleration.Value != "" || toleration.Effect != "" || toleration.TolerationSeconds != nil {
					service.Tolerations = append(service.Tolerations, toleration)
				}
			}
		}
	}
	if hostNetwork, ok := podSpec["hostNetwork"].(bool); ok {
		service.HostNetwork = hostNetwork
		service.HostNetworkSet = true
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.hostNetwork"] = fmt.Sprintf("%t", hostNetwork)
	}
	if value, ok := podSpec["hostPID"].(bool); ok {
		service.HostPID = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.hostPID"] = fmt.Sprintf("%t", value)
		if value {
			service.PIDMode = "host"
			service.Extensions["kubernetes.pidMode"] = "host"
		}
	}
	if value, ok := podSpec["hostIPC"].(bool); ok {
		service.HostIPC = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.hostIPC"] = fmt.Sprintf("%t", value)
		if value {
			service.IPCMode = "host"
			service.Extensions["kubernetes.ipcMode"] = "host"
		}
	}
	if value := toString(podSpec["dnsPolicy"]); value != "" {
		service.DNSPolicy = value
	}
	if osSpec, ok := asMap(podSpec["os"]); ok {
		if name := toString(osSpec["name"]); name != "" {
			service.OSName = name
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["kubernetes.os"] = name
		}
	}
	if value := toString(podSpec["schedulerName"]); value != "" {
		service.SchedulerName = value
	}
	if value := toString(podSpec["priorityClassName"]); value != "" {
		service.PriorityClassName = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.priorityClassName"] = value
	}
	if value := toString(podSpec["runtimeClassName"]); value != "" {
		service.RuntimeClassName = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.runtimeClassName"] = value
	}
	if value := toString(podSpec["nodeName"]); value != "" {
		service.NodeName = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.nodeName"] = value
	}
	if nodeSelector, ok := asMap(podSpec["nodeSelector"]); ok && len(nodeSelector) > 0 {
		if selector, err := toStringMap(nodeSelector); err == nil {
			service.NodeSelector = copyStringMap(selector)
		}
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.nodeSelector"] = copyStringMap(service.NodeSelector)
		service.Extensions["x-kubernetes-node-selector"] = copyStringMap(service.NodeSelector)
	}
	if value := toString(podSpec["subdomain"]); value != "" {
		service.Subdomain = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.subdomain"] = value
	}
	if value, ok := podSpec["setHostnameAsFQDN"].(bool); ok {
		service.SetHostnameAsFQDN = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.setHostnameAsFQDN"] = fmt.Sprintf("%t", value)
	}
	if value, ok := podSpec["hostUsers"].(bool); ok {
		service.HostUsers = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.hostUsers"] = fmt.Sprintf("%t", value)
	}
	if value, ok := podSpec["shareProcessNamespace"].(bool); ok {
		service.ShareProcessNamespace = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.shareProcessNamespace"] = fmt.Sprintf("%t", value)
	}
	if value, ok := podSpec["enableServiceLinks"].(bool); ok {
		service.EnableServiceLinks = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.enableServiceLinks"] = fmt.Sprintf("%t", value)
	}
	if podSecurityContext, ok := asMap(podSpec["securityContext"]); ok {
		if value, ok := podSecurityContext["fsGroup"]; ok {
			fsGroup := int64(toInt(value))
			if fsGroup > 0 {
				service.FSGroup = &fsGroup
			}
		}
		if value := toString(podSecurityContext["fsGroupChangePolicy"]); value != "" {
			service.FSGroupChangePolicy = value
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-fsGroupChangePolicy"] = value
		}
		if value, ok := podSecurityContext["runAsNonRoot"].(bool); ok {
			service.RunAsNonRoot = boolPtr(value)
		}
		if values, ok := podSecurityContext["supplementalGroups"].([]interface{}); ok {
			for _, value := range values {
				if group := int64(toInt(value)); group > 0 {
					service.SupplementalGroups = append(service.SupplementalGroups, group)
				}
			}
		}
		if value := toString(podSecurityContext["supplementalGroupsPolicy"]); value != "" {
			service.SupplementalGroupsPolicy = value
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-supplementalGroupsPolicy"] = value
		}
		if selinux, ok := asMap(podSecurityContext["seLinuxOptions"]); ok {
			service.SELinuxOptions = parseKubernetesSELinuxOptions(selinux)
			if service.SELinuxOptions != nil {
				if service.Extensions == nil {
					service.Extensions = map[string]interface{}{}
				}
				service.Extensions["x-kubernetes-seLinuxOptions"] = cloneMap(selinux)
			}
		}
		if windows, ok := asMap(podSecurityContext["windowsOptions"]); ok {
			service.WindowsOptions = mergeWindowsSecurityContextOptions(service.WindowsOptions, parseKubernetesWindowsSecurityContextOptions(windows))
			applyWindowsHostProcessDefaults(service)
			if service.WindowsOptions != nil {
				if service.Extensions == nil {
					service.Extensions = map[string]interface{}{}
				}
				service.Extensions["x-kubernetes-windowsOptions"] = cloneMap(windows)
			}
		}
	}
	if value := toInt(podSpec["activeDeadlineSeconds"]); value > 0 {
		activeDeadlineSeconds := int64(value)
		service.ActiveDeadlineSeconds = &activeDeadlineSeconds
	}
	if value := toString(podSpec["restartPolicy"]); value != "" {
		service.PodRestartPolicy = value
	}
	if value := toString(podSpec["hostname"]); value != "" {
		service.Hostname = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.hostname"] = value
	}
	if value := toString(podSpec["serviceAccountName"]); value != "" {
		service.ServiceAccountName = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.serviceAccountName"] = value
	}
	if value, ok := podSpec["automountServiceAccountToken"].(bool); ok {
		service.AutomountServiceAccountToken = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.automountServiceAccountToken"] = fmt.Sprintf("%t", value)
	}

	if ports, ok := container["ports"].([]interface{}); ok {
		for _, port := range ports {
			if portMap, ok := asMap(port); ok {
				portMapping := PortMapping{Extensions: map[string]interface{}{}}
				portMapping.Name = toString(portMap["name"])
				portMapping.ContainerPort = strconv.Itoa(toInt(portMap["containerPort"]))
				if protocol, ok := portMap["protocol"].(string); ok {
					portMapping.Protocol = strings.ToLower(protocol)
				}
				portMapping.AppProtocol = toString(portMap["appProtocol"])
				for key, value := range portMap {
					switch key {
					case "name", "containerPort", "protocol", "appProtocol":
						continue
					default:
						portMapping.Extensions[key] = deepCopyValue(value)
					}
				}
				if len(portMapping.Extensions) == 0 {
					portMapping.Extensions = nil
				}
				mergeKubernetesContainerPort(service, portMapping)
			}
		}
	}

	configVolumeSources, secretVolumeSources, tmpfsVolumeSources := kubernetesVolumeSources(podSpec)
	if volumeMounts, ok := container["volumeMounts"].([]interface{}); ok {
		for _, mount := range volumeMounts {
			if mountMap, ok := asMap(mount); ok {
				volumeMount := VolumeMount{}
				if name, ok := mountMap["name"].(string); ok {
					volumeMount.Source = name
				}
				if mountPath, ok := mountMap["mountPath"].(string); ok {
					volumeMount.Target = mountPath
				}
				if readOnly, ok := mountMap["readOnly"].(bool); ok {
					volumeMount.ReadOnly = readOnly
				}
				if propagation, ok := mountMap["mountPropagation"].(string); ok {
					volumeMount.MountPropagation = propagation
				}
				if recursiveReadOnly, ok := mountMap["recursiveReadOnly"].(string); ok {
					volumeMount.RecursiveReadOnly = recursiveReadOnly
				}
				if subPath, ok := mountMap["subPath"].(string); ok {
					volumeMount.SubPath = subPath
				}
				if subPathExpr, ok := mountMap["subPathExpr"].(string); ok {
					volumeMount.SubPathExpr = subPathExpr
				}
				volumeMount.Type = "volume"
				if tmpfs, ok := tmpfsVolumeSources[volumeMount.Source]; ok {
					volumeMount.Type = "tmpfs"
					if size := toString(tmpfs["sizeLimit"]); size != "" {
						volumeMount.Options = ensureStringMap(volumeMount.Options)
						volumeMount.Options["size"] = size
						if volumeMount.TmpfsExtensions == nil {
							volumeMount.TmpfsExtensions = map[string]interface{}{}
						}
						if _, exists := volumeMount.TmpfsExtensions["size"]; !exists {
							volumeMount.TmpfsExtensions["size"] = size
						}
					}
				}
				service.Volumes = append(service.Volumes, volumeMount)
				if refs, ok := configVolumeSources[volumeMount.Source]; ok {
					for _, ref := range refs {
						ref.Target = kubernetesMountedFileTarget(volumeMount.Target, ref.Target)
						ref.ReadOnly = volumeMount.ReadOnly
						service.Configs = append(service.Configs, ref)
					}
				}
				if refs, ok := secretVolumeSources[volumeMount.Source]; ok {
					for _, ref := range refs {
						ref.Target = kubernetesMountedFileTarget(volumeMount.Target, ref.Target)
						ref.ReadOnly = volumeMount.ReadOnly
						service.Secrets = append(service.Secrets, ref)
					}
				}
			}
		}
	}
	parseKubernetesPortableFileRefs(service, annotations)
	parseKubernetesPortableVolumes(service, annotations)

	// Parse volumes at pod level
	if volumes, ok := podSpec["volumes"].([]interface{}); ok {
		for _, vol := range volumes {
			if volMap, ok := asMap(vol); ok {
				if name, ok := volMap["name"].(string); ok {
					volume := &Volume{Name: name}

					if configMap, ok := asMap(volMap["configMap"]); ok {
						volume.Driver = "configMap"
						if cmName, ok := configMap["name"].(string); ok {
							volume.DriverOpts = map[string]string{"configMap": cmName}
						}
					} else if secret, ok := asMap(volMap["secret"]); ok {
						volume.Driver = "secret"
						if secretName, ok := secret["secretName"].(string); ok {
							volume.DriverOpts = map[string]string{"secretName": secretName}
						}
					} else if hostPath, ok := asMap(volMap["hostPath"]); ok {
						volume.Driver = "hostPath"
						if path, ok := hostPath["path"].(string); ok {
							volume.DriverOpts = map[string]string{"path": path}
						}
					} else if emptyDir, ok := asMap(volMap["emptyDir"]); ok {
						volume.Driver = "emptyDir"
						if medium := toString(emptyDir["medium"]); strings.EqualFold(medium, "Memory") {
							volume.DriverOpts = map[string]string{"medium": "Memory"}
							if sizeLimit := toString(emptyDir["sizeLimit"]); sizeLimit != "" {
								volume.DriverOpts["sizeLimit"] = sizeLimit
							}
						}
					}

					app.Volumes[name] = volume
				}
			}
		}
	}
	parseKubernetesPortableAnnotations(app, service, annotations)
	parseKubernetesPortableVolumes(service, annotations)

	if resources, ok := asMap(container["resources"]); ok {
		service.Deploy = ensureServiceDeploy(service)
		service.Deploy.Resources = parseKubernetesResources(resources)
		if service.Deploy.Resources != nil {
			if service.Deploy.Resources.CPULimit != "" {
				service.CPUs = service.Deploy.Resources.CPULimit
			} else if service.Deploy.Resources.CPUReservation != "" {
				service.CPUs = service.Deploy.Resources.CPUReservation
			}
			if service.Deploy.Resources.MemoryLimit != "" {
				service.MemLimit = service.Deploy.Resources.MemoryLimit
				service.MemoryLimit = service.Deploy.Resources.MemoryLimit
			}
			if service.Deploy.Resources.MemoryReservation != "" {
				service.MemReservation = service.Deploy.Resources.MemoryReservation
			}
		}
	}
	parseKubernetesDeployIntent(service, kind, spec, podSpec, annotations)
	if readinessProbe, ok := asMap(container["readinessProbe"]); ok {
		service.HealthCheck = mergeHealthCheckSpec(parseKubernetesProbe(readinessProbe), service.HealthCheck)
	} else if livenessProbe, ok := asMap(container["livenessProbe"]); ok {
		service.HealthCheck = mergeHealthCheckSpec(parseKubernetesProbe(livenessProbe), service.HealthCheck)
	}
	if startupProbe, ok := asMap(container["startupProbe"]); ok {
		service.StartupProbe = mergeHealthCheckSpec(parseKubernetesProbe(startupProbe), service.StartupProbe)
		if service.StartupProbe != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-startup-probe"] = cloneMap(startupProbe)
		}
	}
	if readinessGates, ok := podSpec["readinessGates"].([]interface{}); ok {
		for _, gate := range readinessGates {
			if gateMap, ok := asMap(gate); ok {
				service.ReadinessGates = append(service.ReadinessGates, cloneMap(gateMap))
			}
		}
		if len(service.ReadinessGates) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-readiness-gates"] = cloneMapSlice(service.ReadinessGates)
		}
	}
	if initContainers, ok := podSpec["initContainers"].([]interface{}); ok {
		containers := make([]map[string]interface{}, 0, len(initContainers))
		for _, container := range initContainers {
			if containerMap, ok := asMap(container); ok {
				containers = append(containers, cloneMap(containerMap))
			}
		}
		if len(containers) > 0 {
			service.InitContainers = cloneMapSlice(containers)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-init-containers"] = cloneMapSlice(containers)
		}
	}
	if resourceClaims, ok := podSpec["resourceClaims"].([]interface{}); ok {
		claims := make([]map[string]interface{}, 0, len(resourceClaims))
		for _, claim := range resourceClaims {
			if claimMap, ok := asMap(claim); ok {
				claims = append(claims, cloneMap(claimMap))
			}
		}
		if len(claims) > 0 {
			service.ResourceClaims = cloneMapSlice(claims)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-resource-claims"] = cloneMapSlice(claims)
		}
	}
	if ephemeralContainers, ok := podSpec["ephemeralContainers"].([]interface{}); ok {
		containers := make([]map[string]interface{}, 0, len(ephemeralContainers))
		for _, container := range ephemeralContainers {
			if containerMap, ok := asMap(container); ok {
				containers = append(containers, cloneMap(containerMap))
			}
		}
		if len(containers) > 0 {
			service.EphemeralContainers = cloneMapSlice(containers)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-ephemeral-containers"] = cloneMapSlice(containers)
		}
	}
	if schedulingGates, ok := podSpec["schedulingGates"].([]interface{}); ok {
		for _, gate := range schedulingGates {
			if gateMap, ok := asMap(gate); ok {
				service.SchedulingGates = append(service.SchedulingGates, cloneMap(gateMap))
			}
		}
		if len(service.SchedulingGates) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-scheduling-gates"] = cloneMapSlice(service.SchedulingGates)
		}
	}

	app.Services[name] = service
	workloadServiceNames := []string{name}
	for index := 1; index < len(containers); index++ {
		sidecarContainer, ok := asMap(containers[index])
		if !ok {
			continue
		}
		sidecarName := kubernetesContainerServiceName(name, sidecarContainer, index)
		if sidecarName == "" || sidecarName == name {
			sidecarName = fmt.Sprintf("%s-container-%d", name, index+1)
		}
		if app.Services[sidecarName] != nil {
			sidecarName = fmt.Sprintf("%s-container-%d", name, index+1)
		}
		app.Services[sidecarName] = parseKubernetesAdditionalContainerService(app, name, sidecarName, kind, spec, podMetadata, podSpec, sidecarContainer, annotations)
		workloadServiceNames = append(workloadServiceNames, sidecarName)
	}
	for _, initServiceName := range parseKubernetesInitContainerServices(app, name, kind, spec, podMetadata, podSpec, annotations) {
		dependency := DependencySpec{Name: initServiceName, Condition: "service_completed_successfully"}
		for _, serviceName := range workloadServiceNames {
			if workloadService := app.Services[serviceName]; workloadService != nil {
				workloadService.Dependencies = appendUniqueDependency(workloadService.Dependencies, dependency)
				workloadService.DependsOn = appendUniqueName(workloadService.DependsOn, initServiceName)
			}
		}
	}
	return nil
}

func kubernetesContainerServiceName(workloadName string, container map[string]interface{}, index int) string {
	containerName := toString(container["name"])
	if containerName == "" {
		return fmt.Sprintf("%s-container-%d", workloadName, index+1)
	}
	if containerName == workloadName {
		return workloadName
	}
	return workloadName + "-" + containerName
}

func parseKubernetesAdditionalContainerService(app *Application, workloadName, serviceName, kind string, spec, podMetadata, podSpec, container map[string]interface{}, annotations map[string]string) *Service {
	service := &Service{
		Name:        serviceName,
		Platform:    PlatformKubernetes,
		Environment: make(map[string]string),
		Labels:      make(map[string]string),
		Extensions: map[string]interface{}{
			"kubernetes.kind":      kind,
			"kubernetes.workload":  workloadName,
			"kubernetes.container": toString(container["name"]),
			"kubernetes.sidecar":   "true",
		},
	}
	if image := toString(container["image"]); image != "" {
		service.Image = image
	}
	service.Replicas = workloadReplicas(kind, spec)
	if service.Replicas > 0 {
		service.Deploy = &DeploySpec{Replicas: service.Replicas}
	}
	for key, value := range toStringMapLoose(podMetadata["labels"]) {
		service.Labels[key] = value
	}
	service.Labels["bolabaden.dev/workload"] = workloadName
	service.Labels["bolabaden.dev/container"] = toString(container["name"])
	if len(annotations) > 0 {
		service.Extensions["kubernetes.annotations"] = annotations
	}
	parseKubernetesDependencies(service, annotations, podSpec)
	parseKubernetesRuntimeAnnotations(service, annotations)
	parseKubernetesLifecycle(service, container)
	parseKubernetesPortableAnnotations(app, service, annotations)

	if command, ok := container["command"].([]interface{}); ok {
		service.Entrypoint = interfaceSliceToStringSlice(command)
	}
	if args, ok := container["args"].([]interface{}); ok {
		service.Command = interfaceSliceToStringSlice(args)
	}
	if workingDir := toString(container["workingDir"]); workingDir != "" {
		service.WorkingDir = workingDir
	}
	if stdin, ok := container["stdin"].(bool); ok {
		service.StdinOpen = stdin
		service.StdinOpenSet = true
	}
	if tty, ok := container["tty"].(bool); ok {
		service.Tty = tty
		service.TtySet = true
	}
	parseKubernetesRuntime(service, podSpec, container)
	if seconds := toInt(podSpec["terminationGracePeriodSeconds"]); seconds > 0 {
		service.StopGracePeriod = fmt.Sprintf("%ds", seconds)
	}
	parseKubernetesDNS(service, podSpec)
	parseKubernetesContainerEnv(service, container)
	parseKubernetesContainerPorts(service, container)
	parseKubernetesContainerVolumeMounts(service, podSpec, container)
	parseKubernetesPortableFileRefs(service, annotations)
	parseKubernetesPortableAnnotations(app, service, annotations)
	parseKubernetesPortableVolumes(service, annotations)
	if resources, ok := asMap(container["resources"]); ok {
		service.Deploy = ensureServiceDeploy(service)
		service.Deploy.Resources = parseKubernetesResources(resources)
		if service.Deploy.Resources != nil {
			if service.Deploy.Resources.CPULimit != "" {
				service.CPUs = service.Deploy.Resources.CPULimit
			} else if service.Deploy.Resources.CPUReservation != "" {
				service.CPUs = service.Deploy.Resources.CPUReservation
			}
			if service.Deploy.Resources.MemoryLimit != "" {
				service.MemLimit = service.Deploy.Resources.MemoryLimit
				service.MemoryLimit = service.Deploy.Resources.MemoryLimit
			}
			if service.Deploy.Resources.MemoryReservation != "" {
				service.MemReservation = service.Deploy.Resources.MemoryReservation
			}
		}
	}
	parseKubernetesDeployIntent(service, kind, spec, podSpec, annotations)
	if readinessProbe, ok := asMap(container["readinessProbe"]); ok {
		service.HealthCheck = parseKubernetesProbe(readinessProbe)
	} else if livenessProbe, ok := asMap(container["livenessProbe"]); ok {
		service.HealthCheck = parseKubernetesProbe(livenessProbe)
	}
	if startupProbe, ok := asMap(container["startupProbe"]); ok {
		service.StartupProbe = parseKubernetesProbe(startupProbe)
		if service.StartupProbe != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-startup-probe"] = cloneMap(startupProbe)
		}
	}
	if readinessGates, ok := podSpec["readinessGates"].([]interface{}); ok {
		for _, gate := range readinessGates {
			if gateMap, ok := asMap(gate); ok {
				service.ReadinessGates = append(service.ReadinessGates, cloneMap(gateMap))
			}
		}
		if len(service.ReadinessGates) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-readiness-gates"] = cloneMapSlice(service.ReadinessGates)
		}
	}
	if resourceClaims, ok := podSpec["resourceClaims"].([]interface{}); ok {
		claims := make([]map[string]interface{}, 0, len(resourceClaims))
		for _, claim := range resourceClaims {
			if claimMap, ok := asMap(claim); ok {
				claims = append(claims, cloneMap(claimMap))
			}
		}
		if len(claims) > 0 {
			service.ResourceClaims = cloneMapSlice(claims)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-resource-claims"] = cloneMapSlice(claims)
		}
	}
	if ephemeralContainers, ok := podSpec["ephemeralContainers"].([]interface{}); ok {
		containers := make([]map[string]interface{}, 0, len(ephemeralContainers))
		for _, container := range ephemeralContainers {
			if containerMap, ok := asMap(container); ok {
				containers = append(containers, cloneMap(containerMap))
			}
		}
		if len(containers) > 0 {
			service.EphemeralContainers = cloneMapSlice(containers)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-ephemeral-containers"] = cloneMapSlice(containers)
		}
	}
	if schedulingGates, ok := podSpec["schedulingGates"].([]interface{}); ok {
		for _, gate := range schedulingGates {
			if gateMap, ok := asMap(gate); ok {
				service.SchedulingGates = append(service.SchedulingGates, cloneMap(gateMap))
			}
		}
		if len(service.SchedulingGates) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-scheduling-gates"] = cloneMapSlice(service.SchedulingGates)
		}
	}
	if constraints, ok := podSpec["topologySpreadConstraints"].([]interface{}); ok {
		for _, constraint := range constraints {
			if constraintMap, ok := asMap(constraint); ok {
				service.TopologySpreadConstraints = append(service.TopologySpreadConstraints, cloneMap(constraintMap))
			}
		}
		if len(service.TopologySpreadConstraints) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-topology-spread-constraints"] = cloneMapSlice(service.TopologySpreadConstraints)
		}
	}
	return service
}

func parseKubernetesInitContainerServices(app *Application, workloadName, kind string, spec, podMetadata, podSpec map[string]interface{}, annotations map[string]string) []string {
	initContainers, _ := podSpec["initContainers"].([]interface{})
	var serviceNames []string
	for index, initContainerValue := range initContainers {
		initContainer, ok := asMap(initContainerValue)
		if !ok {
			continue
		}
		if strings.HasPrefix(toString(initContainer["name"]), "wait-for-") {
			continue
		}
		serviceName := kubernetesContainerServiceName(workloadName, initContainer, index)
		if serviceName == "" || serviceName == workloadName {
			serviceName = fmt.Sprintf("%s-init-%d", workloadName, index+1)
		}
		if app.Services[serviceName] != nil {
			serviceName = fmt.Sprintf("%s-init-%d", workloadName, index+1)
		}
		service := parseKubernetesAdditionalContainerService(app, workloadName, serviceName, kind, spec, podMetadata, podSpec, initContainer, annotations)
		service.Replicas = 0
		service.Deploy = nil
		service.Restart = "no"
		delete(service.Extensions, "kubernetes.sidecar")
		service.Extensions["kubernetes.initContainer"] = "true"
		service.Labels["bolabaden.dev/init-container"] = toString(initContainer["name"])
		app.Services[serviceName] = service
		serviceNames = append(serviceNames, serviceName)
	}
	return serviceNames
}

func parseKubernetesContainerEnv(service *Service, container map[string]interface{}) {
	if env, ok := container["env"].([]interface{}); ok {
		for _, envVar := range env {
			if envMap, ok := asMap(envVar); ok {
				if name, ok := envMap["name"].(string); ok {
					if value, ok := envMap["value"].(string); ok {
						service.Environment[name] = value
					} else if valueFrom, ok := asMap(envMap["valueFrom"]); ok {
						if source := kubernetesEnvSource(name, valueFrom); source.Name != "" {
							service.EnvSources = append(service.EnvSources, source)
						}
					}
				}
			}
		}
	}
	if envFrom, ok := container["envFrom"].([]interface{}); ok {
		for _, envFromItem := range envFrom {
			if envFromMap, ok := asMap(envFromItem); ok {
				if source := kubernetesEnvFromSource(envFromMap); source.Source != "" {
					service.EnvFrom = append(service.EnvFrom, source)
				}
			}
		}
	}
}

func parseKubernetesContainerPorts(service *Service, container map[string]interface{}) {
	if ports, ok := container["ports"].([]interface{}); ok {
		for _, port := range ports {
			if portMap, ok := asMap(port); ok {
				portMapping := PortMapping{Extensions: map[string]interface{}{}}
				portMapping.Name = toString(portMap["name"])
				portMapping.ContainerPort = strconv.Itoa(toInt(portMap["containerPort"]))
				if protocol, ok := portMap["protocol"].(string); ok {
					portMapping.Protocol = strings.ToLower(protocol)
				}
				portMapping.AppProtocol = toString(portMap["appProtocol"])
				for key, value := range portMap {
					switch key {
					case "name", "containerPort", "protocol", "appProtocol":
						continue
					default:
						portMapping.Extensions[key] = deepCopyValue(value)
					}
				}
				if len(portMapping.Extensions) == 0 {
					portMapping.Extensions = nil
				}
				mergeKubernetesContainerPort(service, portMapping)
			}
		}
	}
}

func mergeKubernetesContainerPort(service *Service, port PortMapping) {
	if service == nil {
		return
	}
	index := matchingPortIndex(service.Ports, port)
	if index < 0 {
		service.Ports = append(service.Ports, port)
		return
	}
	service.Ports[index] = mergePortMapping(service.Ports[index], port)
}

func parseKubernetesContainerVolumeMounts(service *Service, podSpec, container map[string]interface{}) {
	configVolumeSources, secretVolumeSources, tmpfsVolumeSources := kubernetesVolumeSources(podSpec)
	if volumeMounts, ok := container["volumeMounts"].([]interface{}); ok {
		for _, mount := range volumeMounts {
			if mountMap, ok := asMap(mount); ok {
				volumeMount := VolumeMount{}
				if name, ok := mountMap["name"].(string); ok {
					volumeMount.Source = name
				}
				if mountPath, ok := mountMap["mountPath"].(string); ok {
					volumeMount.Target = mountPath
				}
				if readOnly, ok := mountMap["readOnly"].(bool); ok {
					volumeMount.ReadOnly = readOnly
				}
				volumeMount.Type = "volume"
				if tmpfs, ok := tmpfsVolumeSources[volumeMount.Source]; ok {
					volumeMount.Type = "tmpfs"
					if size := toString(tmpfs["sizeLimit"]); size != "" {
						volumeMount.Options = ensureStringMap(volumeMount.Options)
						volumeMount.Options["size"] = size
						if volumeMount.TmpfsExtensions == nil {
							volumeMount.TmpfsExtensions = map[string]interface{}{}
						}
						if _, exists := volumeMount.TmpfsExtensions["size"]; !exists {
							volumeMount.TmpfsExtensions["size"] = size
						}
					}
				}
				service.Volumes = append(service.Volumes, volumeMount)
				if refs, ok := configVolumeSources[volumeMount.Source]; ok {
					for _, ref := range refs {
						ref.Target = kubernetesMountedFileTarget(volumeMount.Target, ref.Target)
						ref.ReadOnly = volumeMount.ReadOnly
						service.Configs = append(service.Configs, ref)
					}
				}
				if refs, ok := secretVolumeSources[volumeMount.Source]; ok {
					for _, ref := range refs {
						ref.Target = kubernetesMountedFileTarget(volumeMount.Target, ref.Target)
						ref.ReadOnly = volumeMount.ReadOnly
						service.Secrets = append(service.Secrets, ref)
					}
				}
			}
		}
	}
}

func parseKubernetesRuntime(service *Service, podSpec, container map[string]interface{}) {
	podSecurityContext, _ := asMap(podSpec["securityContext"])
	containerSecurityContext, _ := asMap(container["securityContext"])
	securityContext := mergeLooseMaps(podSecurityContext, containerSecurityContext)
	if privileged, ok := securityContext["privileged"].(bool); ok {
		service.Privileged = privileged
		service.PrivilegedSet = true
	}
	if readOnly, ok := securityContext["readOnlyRootFilesystem"].(bool); ok {
		service.ReadOnlyRootFS = readOnly
		service.ReadOnlyRootFSSet = true
	}
	if value, ok := securityContext["allowPrivilegeEscalation"].(bool); ok {
		service.AllowPrivilegeEscalation = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["x-kubernetes-allowPrivilegeEscalation"] = value
	}
	if value := toString(securityContext["procMount"]); value != "" {
		service.ProcMount = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["x-kubernetes-procMount"] = value
	}
	if runAsUser := toString(securityContext["runAsUser"]); runAsUser != "" {
		service.User = runAsUser
	}
	if runAsGroup := toString(securityContext["runAsGroup"]); runAsGroup != "" {
		service.Group = runAsGroup
	}
	if capabilities, ok := asMap(securityContext["capabilities"]); ok {
		if add, ok := capabilities["add"].([]interface{}); ok {
			service.CapAdd = interfaceSliceToStringSlice(add)
		}
		if drop, ok := capabilities["drop"].([]interface{}); ok {
			service.CapDrop = interfaceSliceToStringSlice(drop)
		}
	}
	if selinux, ok := asMap(securityContext["seLinuxOptions"]); ok {
		service.SELinuxOptions = parseKubernetesSELinuxOptions(selinux)
		if service.SELinuxOptions != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-seLinuxOptions"] = cloneMap(selinux)
		}
	}
	if windows, ok := asMap(securityContext["windowsOptions"]); ok {
		service.WindowsOptions = mergeWindowsSecurityContextOptions(service.WindowsOptions, parseKubernetesWindowsSecurityContextOptions(windows))
		applyWindowsHostProcessDefaults(service)
		if service.WindowsOptions != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-windowsOptions"] = cloneMap(windows)
		}
	}
	if sysctls, ok := podSecurityContext["sysctls"].([]interface{}); ok {
		service.Sysctls = kubernetesSysctlsToMap(sysctls)
	}
	if seccompProfile, ok := asMap(securityContext["seccompProfile"]); ok {
		service.SeccompProfile = parseKubernetesSeccompProfile(seccompProfile)
		if service.SeccompProfile != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-seccomp-profile"] = cloneMap(seccompProfile)
		}
	}
	if affinity, ok := asMap(podSpec["affinity"]); ok {
		service.Affinity = cloneMap(affinity)
		if service.Affinity != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-affinity"] = cloneMap(affinity)
		}
	}
	if constraints, ok := podSpec["topologySpreadConstraints"].([]interface{}); ok {
		for _, constraint := range constraints {
			if constraintMap, ok := asMap(constraint); ok {
				service.TopologySpreadConstraints = append(service.TopologySpreadConstraints, cloneMap(constraintMap))
			}
		}
		if len(service.TopologySpreadConstraints) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-topology-spread-constraints"] = cloneMapSlice(service.TopologySpreadConstraints)
		}
	}
	if value := toString(podSpec["serviceAccountName"]); value != "" {
		service.ServiceAccountName = value
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.serviceAccountName"] = value
	}
	if value, ok := podSpec["automountServiceAccountToken"].(bool); ok {
		service.AutomountServiceAccountToken = boolPtr(value)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.automountServiceAccountToken"] = fmt.Sprintf("%t", value)
	}
}

func parseKubernetesRuntimeAnnotations(service *Service, annotations map[string]string) {
	for _, option := range splitAnnotationList(annotations[kubernetesAnnotationRuntimeSecurityOpt]) {
		appendUniqueString(&service.SecurityOpt, option)
	}
	if value := annotations[kubernetesAnnotationRuntimeInit]; value != "" {
		init := strings.EqualFold(value, "true")
		service.Init = &init
	}
	if value := annotations[kubernetesAnnotationRuntimeStopSignal]; value != "" {
		service.StopSignal = value
	}
	if value := annotations[kubernetesAnnotationPortablePIDMode]; value != "" {
		service.PIDMode = value
		if strings.EqualFold(value, "host") {
			service.HostPID = boolPtr(true)
		}
	}
	if value := annotations[kubernetesAnnotationPortableIPCMode]; value != "" {
		service.IPCMode = value
		if strings.EqualFold(value, "host") {
			service.HostIPC = boolPtr(true)
		}
	}
	if value := annotations[kubernetesAnnotationPortablePidsLimit]; value != "" {
		service.PidsLimit = int64(parseInt(value))
		service.pidsLimitSet = true
	}
	if value := annotations[kubernetesAnnotationPortableShmSize]; value != "" {
		service.ShmSize = int64(parseInt(value))
		service.shmSizeSet = true
	}
}

func parseKubernetesPortableAnnotations(app *Application, service *Service, annotations map[string]string) {
	if raw := annotations[kubernetesAnnotationPortableAppExtensions]; raw != "" {
		var extensionMap map[string]interface{}
		if json.Unmarshal([]byte(raw), &extensionMap) == nil && len(extensionMap) > 0 {
			if app.Extensions == nil {
				app.Extensions = map[string]interface{}{}
			}
			for key, value := range extensionMap {
				if _, exists := app.Extensions[key]; !exists {
					app.Extensions[key] = deepCopyValue(value)
				}
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableMesh]; raw != "" {
		var mapped map[string]interface{}
		if json.Unmarshal([]byte(raw), &mapped) == nil && len(mapped) > 0 {
			if mesh := meshSpecFromAny(mapped); mesh != nil {
				app.Mesh = mesh
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableServiceExtensions]; raw != "" {
		var extensionMap map[string]interface{}
		if json.Unmarshal([]byte(raw), &extensionMap) == nil && len(extensionMap) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			for key, value := range extensionMap {
				if _, exists := service.Extensions[key]; !exists {
					service.Extensions[key] = deepCopyValue(value)
				}
			}
		}
	}
	if name := annotations[kubernetesAnnotationPortableComposeName]; name != "" && app.Name == "" {
		app.Name = name
	}
	if raw := annotations[kubernetesAnnotationPortableModels]; raw != "" {
		var decoded map[string]map[string]interface{}
		if json.Unmarshal([]byte(raw), &decoded) == nil && len(decoded) > 0 {
			models := composeModelsFromExtensionMap(decoded)
			if len(models) > 0 {
				if app.Models == nil {
					app.Models = map[string]*ComposeModel{}
				}
				for name, model := range models {
					if _, exists := app.Models[name]; !exists {
						app.Models[name] = cloneComposeModel(model)
					}
				}
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableIncludes]; raw != "" {
		var includes []interface{}
		if json.Unmarshal([]byte(raw), &includes) == nil && len(includes) > 0 {
			app.IncludeEntries = mergeIncludeEntries(app.IncludeEntries, includes)
			app.Includes = mergeUniqueStrings(app.Includes, composeIncludePaths(includes))
		}
	}
	if raw := annotations[kubernetesAnnotationPortableNetworks]; raw != "" {
		var networks map[string]*Network
		if json.Unmarshal([]byte(raw), &networks) == nil && len(networks) > 0 {
			if app.Networks == nil {
				app.Networks = map[string]*Network{}
			}
			for name, network := range networks {
				if network == nil {
					continue
				}
				if _, exists := app.Networks[name]; !exists {
					app.Networks[name] = cloneNetwork(network)
				}
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableAppVolumes]; raw != "" {
		var volumes map[string]*Volume
		if json.Unmarshal([]byte(raw), &volumes) == nil && len(volumes) > 0 {
			if app.Volumes == nil {
				app.Volumes = map[string]*Volume{}
			}
			for name, volume := range volumes {
				app.Volumes[name] = cloneVolume(volume)
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableAppConfigs]; raw != "" {
		var configs map[string]*Config
		if json.Unmarshal([]byte(raw), &configs) == nil && len(configs) > 0 {
			if app.Configs == nil {
				app.Configs = map[string]*Config{}
			}
			for name, config := range configs {
				app.Configs[name] = cloneConfig(config)
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableAppSecrets]; raw != "" {
		var secrets map[string]*Secret
		if json.Unmarshal([]byte(raw), &secrets) == nil && len(secrets) > 0 {
			if app.Secrets == nil {
				app.Secrets = map[string]*Secret{}
			}
			for name, secret := range secrets {
				app.Secrets[name] = cloneSecret(secret)
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableBuild]; raw != "" {
		var build BuildConfig
		if json.Unmarshal([]byte(raw), &build) == nil && buildConfigHasData(&build) {
			service.Build = &build
		}
	}
	if raw := annotations[kubernetesAnnotationPortableDevices]; raw != "" {
		if devices := parseAnnotationStringSlice(raw); len(devices) > 0 {
			service.Devices = devices
		}
	}
	if raw := annotations[kubernetesAnnotationPortableDeviceMappings]; raw != "" {
		var mappings []DeviceMappingSpec
		if json.Unmarshal([]byte(raw), &mappings) == nil && len(mappings) > 0 {
			service.DeviceMappings = cloneDeviceMappings(mappings)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableExpose]; raw != "" {
		if expose := parseAnnotationStringSlice(raw); len(expose) > 0 {
			service.Expose = expose
		}
	}
	if raw := annotations[kubernetesAnnotationPortableHealthcheck]; raw != "" {
		var health HealthCheck
		if json.Unmarshal([]byte(raw), &health) == nil && !isEmptyHealthCheck(&health) {
			service.HealthCheck = mergeHealthCheckSpec(service.HealthCheck, &health)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableStartupProbe]; raw != "" {
		var health HealthCheck
		if json.Unmarshal([]byte(raw), &health) == nil && !isEmptyHealthCheck(&health) {
			service.StartupProbe = mergeHealthCheckSpec(service.StartupProbe, &health)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableLifecycle]; raw != "" {
		var lifecycle LifecycleHooks
		if json.Unmarshal([]byte(raw), &lifecycle) == nil && !isEmptyLifecycleHooks(&lifecycle) {
			service.Lifecycle = mergeLifecycleHooks(service.Lifecycle, &lifecycle)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableTolerations]; raw != "" {
		var tolerations []Toleration
		if json.Unmarshal([]byte(raw), &tolerations) == nil && len(tolerations) > 0 {
			service.Tolerations = cloneTolerations(tolerations)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableDevelop]; raw != "" {
		var develop DevelopConfig
		if json.Unmarshal([]byte(raw), &develop) == nil && !isEmptyDevelopConfig(&develop) {
			service.Develop = &develop
		}
	}
	if raw := annotations[kubernetesAnnotationPortableGroupAdd]; raw != "" {
		if groupAdd := parseAnnotationStringSlice(raw); len(groupAdd) > 0 {
			service.GroupAdd = groupAdd
		}
	}
	if value := annotations[kubernetesAnnotationPortableRuntime]; value != "" {
		service.Runtime = value
	}
	if raw := annotations[kubernetesAnnotationPortableLogging]; raw != "" {
		var logging struct {
			Driver     string                 `json:"driver"`
			Options    map[string]string      `json:"options"`
			Extensions map[string]interface{} `json:"extensions"`
		}
		if json.Unmarshal([]byte(raw), &logging) == nil {
			service.LogDriver = logging.Driver
			service.LogOpt = copyStringMap(logging.Options)
			service.LogExtensions = copyStringInterfaceMap(logging.Extensions)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableComposeCompat]; raw != "" {
		var mapped map[string]interface{}
		if json.Unmarshal([]byte(raw), &mapped) == nil {
			if compat := composeCompatFromExtensionMap(mapped); compat != nil {
				service.ComposeCompat = compat
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableLinks]; raw != "" {
		if links := parseAnnotationStringSlice(raw); len(links) > 0 {
			service.Links = links
		}
	}
	parseKubernetesPortableVolumes(service, annotations)
	if raw := annotations[kubernetesAnnotationPortablePorts]; raw != "" {
		var ports []PortMapping
		if json.Unmarshal([]byte(raw), &ports) == nil && len(ports) > 0 {
			service.Ports = mergePortablePorts(service.Ports, ports)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableEnvFiles]; raw != "" {
		var refs []EnvFileRef
		if json.Unmarshal([]byte(raw), &refs) == nil && len(refs) > 0 {
			service.EnvFileRefs = cloneEnvFileRefs(refs)
			service.EnvFile = envFilePaths(refs)
		}
	}
	if value := annotations[kubernetesAnnotationPortableCPUShares]; value != "" {
		service.CPUShares = parseInt(value)
	}
	if value := annotations[kubernetesAnnotationPortableCPUQuota]; value != "" {
		service.CPUQuota = parseInt(value)
	}
	if value := annotations[kubernetesAnnotationPortableMemLimit]; value != "" {
		service.MemLimit = value
		service.MemoryLimit = value
	}
	if value := annotations[kubernetesAnnotationPortableMemorySwap]; value != "" {
		service.MemorySwap = value
	}
	if value := annotations[kubernetesAnnotationPortableMemReservation]; value != "" {
		service.MemReservation = value
	}
	if value := annotations[kubernetesAnnotationPortableCPUs]; value != "" {
		service.CPUs = value
	}
	if raw := annotations[kubernetesAnnotationPortableUlimits]; raw != "" {
		var limits Ulimits
		if json.Unmarshal([]byte(raw), &limits) == nil && limits.Nofile != nil {
			service.Ulimits = &limits
		}
	}
	if value := annotations[kubernetesAnnotationPortableUserNSMode]; value != "" {
		service.UserNSMode = value
	}
	if value := annotations[kubernetesAnnotationPortablePullPolicy]; value != "" {
		service.PullPolicy = value
	}
	if raw := annotations[kubernetesAnnotationPortableProfiles]; raw != "" {
		if profiles := parseAnnotationStringSlice(raw); len(profiles) > 0 {
			service.Profiles = profiles
		}
	}
	if raw := annotations[kubernetesAnnotationPortableNetworkAttachments]; raw != "" {
		var attachments map[string]*NetworkAttachment
		if json.Unmarshal([]byte(raw), &attachments) == nil && len(attachments) > 0 {
			service.NetworkAttachments = cloneNetworkAttachments(attachments)
			for name := range service.NetworkAttachments {
				appendUniqueString(&service.Networks, name)
			}
			sort.Strings(service.Networks)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableFailover]; raw != "" {
		var mapped map[string]interface{}
		if json.Unmarshal([]byte(raw), &mapped) == nil {
			if failover, err := failoverSpecFromMap(mapped); err == nil {
				service.Failover = cloneFailoverSpec(failover)
			}
		}
	}
	if raw := annotations[kubernetesAnnotationPortableNomadSpread]; raw != "" {
		var spreads []map[string]interface{}
		if json.Unmarshal([]byte(raw), &spreads) == nil && len(spreads) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			cloned := cloneMapSlice(spreads)
			service.Extensions[nomadSpreadExtensionKey] = cloned
			service.Extensions["x-nomad-spread"] = cloneMapSlice(spreads)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableNomadConnect]; raw != "" {
		var connect map[string]interface{}
		if json.Unmarshal([]byte(raw), &connect) == nil && len(connect) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			cloned := copyStringInterfaceMap(connect)
			service.Extensions[nomadConnectExtensionKey] = cloned
			service.Extensions["x-nomad-connect"] = copyStringInterfaceMap(connect)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableNomadRestart]; raw != "" {
		var restart map[string]interface{}
		if json.Unmarshal([]byte(raw), &restart) == nil && len(restart) > 0 {
			applyNomadRestartExtensions(service, restart)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions[nomadRestartExtensionKey] = cloneMap(restart)
			service.Extensions["x-nomad-restart"] = cloneMap(restart)
		}
	}
	for _, item := range []struct {
		annotation string
		key        string
	}{
		{kubernetesAnnotationPortableNomadUpdate, nomadUpdateExtensionKey},
		{kubernetesAnnotationPortableNomadMigrate, nomadMigrateExtensionKey},
		{kubernetesAnnotationPortableNomadReschedule, nomadRescheduleExtensionKey},
	} {
		if raw := annotations[item.annotation]; raw != "" {
			var mapped map[string]interface{}
			if json.Unmarshal([]byte(raw), &mapped) == nil && len(mapped) > 0 {
				if service.Extensions == nil {
					service.Extensions = map[string]interface{}{}
				}
				service.Extensions[item.key] = cloneMap(mapped)
				service.Extensions[composeServiceExtensionKey(item.key)] = cloneMap(mapped)
			}
		}
	}
}

func parseKubernetesPortableFileRefs(service *Service, annotations map[string]string) {
	if service == nil || len(annotations) == 0 {
		return
	}
	if raw := annotations[kubernetesAnnotationPortableConfigs]; raw != "" {
		var refs []FileRef
		if json.Unmarshal([]byte(raw), &refs) == nil && len(refs) > 0 {
			service.Configs = mergePortableFileRefs(service.Configs, refs, false)
		}
	}
	if raw := annotations[kubernetesAnnotationPortableSecrets]; raw != "" {
		var refs []FileRef
		if json.Unmarshal([]byte(raw), &refs) == nil && len(refs) > 0 {
			service.Secrets = mergePortableFileRefs(service.Secrets, refs, false)
		}
	}
}

func parseKubernetesPortableVolumes(service *Service, annotations map[string]string) {
	if service == nil || len(annotations) == 0 {
		return
	}
	if raw := annotations[kubernetesAnnotationPortableVolumes]; raw != "" {
		var volumes []VolumeMount
		if json.Unmarshal([]byte(raw), &volumes) == nil && len(volumes) > 0 {
			service.Volumes = mergePortableVolumeMounts(service.Volumes, volumes)
		}
	}
}

func kubernetesServiceResources(service *Service) *ResourceSpec {
	if service == nil {
		return nil
	}
	if service.Deploy != nil && service.Deploy.Resources != nil {
		resources := cloneResourceSpec(service.Deploy.Resources)
		if claims, ok := kubernetesServiceResourceClaims(service); ok {
			if resources.Extensions == nil {
				resources.Extensions = map[string]interface{}{}
			}
			if _, exists := resources.Extensions["kubernetes.claims"]; !exists {
				resources.Extensions["kubernetes.claims"] = claims
			}
		}
		return resources
	}
	if service.CPUs == "" && service.MemLimit == "" && service.MemReservation == "" && service.MemoryLimit == "" && service.MemorySwap == "" {
		if claims, ok := kubernetesServiceResourceClaims(service); ok {
			return &ResourceSpec{Extensions: map[string]interface{}{"kubernetes.claims": claims}}
		}
		return nil
	}
	result := &ResourceSpec{}
	if service.CPUs != "" {
		result.CPULimit = service.CPUs
		result.CPUReservation = service.CPUs
	}
	if service.MemLimit != "" {
		result.MemoryLimit = service.MemLimit
	} else if service.MemoryLimit != "" {
		result.MemoryLimit = service.MemoryLimit
	}
	if service.MemReservation != "" {
		result.MemoryReservation = service.MemReservation
	}
	if claims, ok := kubernetesServiceResourceClaims(service); ok {
		result.Extensions = map[string]interface{}{"kubernetes.claims": claims}
	}
	return result
}

func kubernetesServiceResourceClaims(service *Service) ([]interface{}, bool) {
	if service == nil {
		return nil, false
	}
	if service.Deploy != nil && service.Deploy.Resources != nil && len(service.Deploy.Resources.Extensions) > 0 {
		if claims, ok := service.Deploy.Resources.Extensions["kubernetes.claims"]; ok {
			if items, ok := claims.([]interface{}); ok && len(items) > 0 {
				return cloneInterfaceSlice(items), true
			}
		}
		if claims, ok := service.Deploy.Resources.Extensions["x-kubernetes-claims"]; ok {
			if items, ok := claims.([]interface{}); ok && len(items) > 0 {
				return cloneInterfaceSlice(items), true
			}
		}
	}
	if claims, ok := service.Extensions["x-kubernetes-claims"]; ok {
		if items, ok := claims.([]interface{}); ok && len(items) > 0 {
			return cloneInterfaceSlice(items), true
		}
	}
	if claims, ok := service.Extensions["kubernetes.claims"]; ok {
		if items, ok := claims.([]interface{}); ok && len(items) > 0 {
			return cloneInterfaceSlice(items), true
		}
	}
	return nil, false
}

func parseKubernetesDNS(service *Service, podSpec map[string]interface{}) {
	if dnsConfig, ok := asMap(podSpec["dnsConfig"]); ok {
		if nameservers, ok := dnsConfig["nameservers"].([]interface{}); ok {
			service.DNS = interfaceSliceToStringSlice(nameservers)
		}
		if searches, ok := dnsConfig["searches"].([]interface{}); ok {
			service.DNSSearch = interfaceSliceToStringSlice(searches)
		}
		if options, ok := dnsConfig["options"].([]interface{}); ok {
			for _, optionValue := range options {
				option, ok := asMap(optionValue)
				if !ok {
					continue
				}
				name := toString(option["name"])
				if name == "" {
					continue
				}
				if value := toString(option["value"]); value != "" {
					service.DNSOptions = append(service.DNSOptions, name+":"+value)
				} else {
					service.DNSOptions = append(service.DNSOptions, name)
				}
			}
		}
	}
	if hostAliases, ok := podSpec["hostAliases"].([]interface{}); ok {
		for _, aliasValue := range hostAliases {
			alias, ok := asMap(aliasValue)
			if !ok {
				continue
			}
			ip := toString(alias["ip"])
			if ip == "" {
				continue
			}
			extensions := map[string]interface{}{}
			for key, value := range alias {
				if key == "ip" || key == "hostnames" {
					continue
				}
				extensions[key] = deepCopyValue(value)
			}
			if hostnames, ok := alias["hostnames"].([]interface{}); ok {
				hostAlias := HostAlias{IP: ip}
				for _, hostname := range interfaceSliceToStringSlice(hostnames) {
					service.ExtraHosts = append(service.ExtraHosts, hostname+"="+ip)
					hostAlias.Hostnames = append(hostAlias.Hostnames, hostname)
				}
				if len(hostAlias.Hostnames) > 0 {
					if len(extensions) > 0 {
						hostAlias.Extensions = extensions
					}
					service.HostAliases = append(service.HostAliases, hostAlias)
				}
			}
		}
	}
}

func workloadTemplate(kind string, spec map[string]interface{}) (map[string]interface{}, bool) {
	if kind == "Pod" {
		return map[string]interface{}{
			"metadata": map[string]interface{}{},
			"spec":     spec,
		}, true
	}
	if kind == "CronJob" {
		jobTemplate, ok := asMap(spec["jobTemplate"])
		if !ok {
			return nil, false
		}
		jobSpec, ok := asMap(jobTemplate["spec"])
		if !ok {
			return nil, false
		}
		template, ok := asMap(jobSpec["template"])
		return template, ok
	}
	template, ok := asMap(spec["template"])
	return template, ok
}

func workloadReplicas(kind string, spec map[string]interface{}) int {
	switch kind {
	case "Pod":
		return 1
	case "DaemonSet":
		return 0
	case "Job", "CronJob":
		return 1
	default:
		return toInt(spec["replicas"])
	}
}

func parseKubernetesWorkloadSpecExtensions(service *Service, kind string, spec map[string]interface{}) {
	if service == nil || len(spec) == 0 {
		return
	}
	if service.Extensions == nil {
		service.Extensions = map[string]interface{}{}
	}
	for _, key := range []string{"revisionHistoryLimit", "minReadySeconds"} {
		if value, ok := spec[key]; ok {
			service.Extensions["kubernetes.workload."+key] = value
		}
	}
	switch kind {
	case "Deployment", "ReplicaSet":
		if value, ok := spec["paused"]; ok {
			service.Extensions["kubernetes.deployment.paused"] = value
		}
		if value, ok := spec["progressDeadlineSeconds"]; ok {
			service.Extensions["kubernetes.deployment.progressDeadlineSeconds"] = value
		}
		if strategy, ok := asMap(spec["strategy"]); ok && len(strategy) > 0 {
			service.Extensions["kubernetes.deployment.strategy"] = cloneMap(strategy)
		}
	case "StatefulSet", "DaemonSet":
		if strategy, ok := asMap(spec["updateStrategy"]); ok && len(strategy) > 0 {
			service.Extensions["kubernetes.workload.updateStrategy"] = cloneMap(strategy)
			if service.Deploy == nil {
				service.Deploy = &DeploySpec{}
			}
			if service.Deploy.UpdateConfig == nil {
				service.Deploy.UpdateConfig = &UpdatePolicy{}
			}
			if service.Deploy.UpdateConfig.Extensions == nil {
				service.Deploy.UpdateConfig.Extensions = map[string]interface{}{}
			}
			service.Deploy.UpdateConfig.Extensions["x-kubernetes-workload-updateStrategy"] = cloneMap(strategy)
		}
		if kind == "StatefulSet" {
			if ordinals, ok := asMap(spec["ordinals"]); ok && len(ordinals) > 0 {
				service.Extensions["kubernetes.statefulSet.ordinals"] = cloneMap(ordinals)
			}
		}
	}
}

func parseKubernetesJobSpecExtensions(service *Service, spec map[string]interface{}) {
	if service == nil || len(spec) == 0 {
		return
	}
	if service.Extensions == nil {
		service.Extensions = map[string]interface{}{}
	}
	if service.Deploy == nil {
		service.Deploy = &DeploySpec{}
	}
	if service.Deploy.Job == nil {
		service.Deploy.Job = &SwarmJobSpec{}
	}
	for _, key := range []string{"parallelism", "completions", "backoffLimit", "backoffLimitPerIndex", "ttlSecondsAfterFinished"} {
		if value, ok := spec[key]; ok {
			service.Extensions["kubernetes.job."+key] = value
			switch key {
			case "parallelism":
				service.Deploy.Job.MaxConcurrent = toInt(value)
				service.Deploy.Job.maxConcurrentSet = true
			case "completions":
				service.Deploy.Job.TotalCompletions = toInt(value)
				service.Deploy.Job.totalCompletionsSet = true
			case "backoffLimit":
				service.Deploy.Job.BackoffLimit = toInt(value)
				service.Deploy.Job.backoffLimitSet = true
			case "backoffLimitPerIndex":
				service.Deploy.Job.BackoffLimitPerIndex = toInt(value)
				service.Deploy.Job.backoffLimitPerIndexSet = true
			case "ttlSecondsAfterFinished":
				service.Deploy.Job.TTLSecondsAfterFinished = toInt(value)
				service.Deploy.Job.ttlSecondsAfterFinishedSet = true
			}
		}
	}
	if value, ok := asMap(spec["podFailurePolicy"]); ok && len(value) > 0 {
		cloned := cloneMap(value)
		service.Extensions["kubernetes.job.podFailurePolicy"] = cloned
		if service.Deploy.Job.Extensions == nil {
			service.Deploy.Job.Extensions = map[string]interface{}{}
		}
		service.Deploy.Job.Extensions["x-kubernetes-job-podFailurePolicy"] = cloned
	}
	if value, ok := asMap(spec["successPolicy"]); ok && len(value) > 0 {
		cloned := cloneMap(value)
		service.Extensions["kubernetes.job.successPolicy"] = cloned
		if service.Deploy.Job.Extensions == nil {
			service.Deploy.Job.Extensions = map[string]interface{}{}
		}
		service.Deploy.Job.Extensions["x-kubernetes-job-successPolicy"] = cloned
	}
	if value, ok := spec["completionMode"]; ok {
		completionMode := toString(value)
		service.Extensions["kubernetes.job.completionMode"] = completionMode
		service.Deploy.Job.CompletionMode = completionMode
		service.Deploy.Job.completionModeSet = true
	}
	if value, ok := spec["suspend"]; ok {
		service.Extensions["kubernetes.job.suspend"] = value
		flag := toBool(value)
		service.Deploy.Job.Suspend = &flag
	}
}

func parseKubernetesCronJobSpecExtensions(service *Service, spec map[string]interface{}) {
	if service == nil || len(spec) == 0 {
		return
	}
	if service.Extensions == nil {
		service.Extensions = map[string]interface{}{}
	}
	for _, key := range []string{"startingDeadlineSeconds", "successfulJobsHistoryLimit", "failedJobsHistoryLimit"} {
		if value, ok := spec[key]; ok {
			service.Extensions["kubernetes.cron."+key] = value
		}
	}
	if jobTemplate, ok := asMap(spec["jobTemplate"]); ok {
		if jobSpec, ok := asMap(jobTemplate["spec"]); ok {
			parseKubernetesJobSpecExtensions(service, jobSpec)
		}
	}
}

func parseKubernetesService(app *Application, resource map[string]interface{}) error {
	metadata, _ := asMap(resource["metadata"])

	name, _ := metadata["name"].(string)
	if name == "" {
		return fmt.Errorf("service missing name")
	}
	if app.KubernetesServices == nil {
		app.KubernetesServices = map[string]*KubernetesServiceSpec{}
	}
	if spec := kubernetesServiceSpecFromMap(resource); spec != nil {
		app.KubernetesServices[name] = spec
	}
	app.Extensions["kubernetes.services"] = appendExtensionSlice(app.Extensions["kubernetes.services"], resource)
	app.Extensions[kubernetesServicesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesServicesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesServicesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesServicesExtensionKey], deepCopyValue(resource))

	return nil
}

func reconcileKubernetesServices(app *Application) {
	value, ok := applicationExtensionValueForKey(app, "kubernetes.services")
	if !ok {
		return
	}
	rawServices, ok := value.([]interface{})
	if !ok {
		return
	}
	for _, raw := range rawServices {
		serviceResource, ok := asMap(raw)
		if !ok {
			continue
		}
		target := targetWorkloadForKubernetesService(app, serviceResource)
		if target == nil {
			continue
		}
		recordKubernetesServiceTarget(app, serviceResource, target.Name)
		recordKubernetesServiceMetadata(target, serviceResource)
		for _, port := range portsFromKubernetesService(serviceResource) {
			mergeServicePort(target, port)
			recordKubernetesServicePortTarget(app, serviceResource, port, target)
		}
	}
}

func reconcileKubernetesHorizontalPodAutoscalers(app *Application) {
	if app == nil {
		return
	}
	attach := func(rawHPAs []interface{}) {
		for _, raw := range rawHPAs {
			hpa, ok := asMap(raw)
			if !ok {
				continue
			}
			spec, _ := asMap(hpa["spec"])
			scaleTargetRef, _ := asMap(spec["scaleTargetRef"])
			targetName := toString(scaleTargetRef["name"])
			if targetName == "" {
				continue
			}
			service := app.Services[targetName]
			if service == nil {
				continue
			}
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			copied, _ := deepCopyValue(hpa).(map[string]interface{})
			service.Extensions[kubernetesHPAExtensionKey] = copied
			service.Extensions[composeKubernetesHPAExtensionKey] = deepCopyValue(hpa)
		}
	}
	if len(app.KubernetesHPAs) > 0 {
		rawHPAs := make([]interface{}, 0, len(app.KubernetesHPAs))
		for _, spec := range app.KubernetesHPAs {
			if spec == nil || len(spec.Raw) == 0 {
				continue
			}
			rawHPAs = append(rawHPAs, deepCopyValue(spec.Raw))
		}
		if len(rawHPAs) > 0 {
			attach(rawHPAs)
			return
		}
	}
	if app.Extensions == nil {
		return
	}
	value, ok := applicationExtensionValueForKey(app, "kubernetes.hpas")
	if !ok {
		return
	}
	rawHPAs, ok := value.([]interface{})
	if !ok {
		return
	}
	attach(rawHPAs)
}

func reconcileKubernetesPodDisruptionBudgets(app *Application) {
	if app == nil {
		return
	}
	attach := func(rawPDBs []interface{}) {
		for _, raw := range rawPDBs {
			pdb, ok := asMap(raw)
			if !ok {
				continue
			}
			service := targetWorkloadForKubernetesPodDisruptionBudget(app, pdb)
			if service == nil {
				continue
			}
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			copied, _ := deepCopyValue(pdb).(map[string]interface{})
			service.Extensions[kubernetesPDBExtensionKey] = copied
			service.Extensions[composeKubernetesPDBExtensionKey] = deepCopyValue(pdb)
		}
	}
	if len(app.KubernetesPDBs) > 0 {
		rawPDBs := make([]interface{}, 0, len(app.KubernetesPDBs))
		for _, spec := range app.KubernetesPDBs {
			if spec == nil || len(spec.Raw) == 0 {
				continue
			}
			rawPDBs = append(rawPDBs, deepCopyValue(spec.Raw))
		}
		if len(rawPDBs) > 0 {
			attach(rawPDBs)
			return
		}
	}
	if app.Extensions == nil {
		return
	}
	value, ok := applicationExtensionValueForKey(app, "kubernetes.pdbs")
	if !ok {
		return
	}
	rawPDBs, ok := value.([]interface{})
	if !ok {
		return
	}
	attach(rawPDBs)
}

func targetWorkloadForKubernetesPodDisruptionBudget(app *Application, pdb map[string]interface{}) *Service {
	if app == nil || len(app.Services) == 0 || len(pdb) == 0 {
		return nil
	}
	spec, _ := asMap(pdb["spec"])
	selector, _ := asMap(spec["selector"])
	matchLabels := toStringMapLoose(selector["matchLabels"])
	if len(matchLabels) == 0 {
		return nil
	}
	var matched *Service
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if !labelsMatch(matchLabels, service.Labels) {
			continue
		}
		if matched != nil {
			return nil
		}
		matched = service
	}
	return matched
}

func reconcileKubernetesServiceAccounts(app *Application) {
	if app == nil || app.Extensions == nil {
		return
	}
	rawServiceAccounts, ok := app.Extensions[kubernetesServiceAccountsExtensionKey].([]interface{})
	if !ok {
		return
	}
	byName := map[string]map[string]interface{}{}
	for _, raw := range rawServiceAccounts {
		serviceAccount, ok := asMap(raw)
		if !ok {
			continue
		}
		metadata, _ := asMap(serviceAccount["metadata"])
		name := toString(metadata["name"])
		if name == "" {
			continue
		}
		byName[name] = serviceAccount
	}
	if len(byName) == 0 {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		serviceAccountName := service.ServiceAccountName
		if serviceAccountName == "" {
			serviceAccountName = extensionStringValue(service, "kubernetes.serviceAccountName")
		}
		if serviceAccountName == "" && service.Extensions != nil {
			serviceAccountName = toString(service.Extensions["x-kubernetes-serviceAccountName"])
		}
		if serviceAccountName == "" {
			continue
		}
		serviceAccount := byName[serviceAccountName]
		if serviceAccount == nil {
			continue
		}
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		copied, _ := deepCopyValue(serviceAccount).(map[string]interface{})
		service.Extensions[kubernetesServiceAccountExtensionKey] = copied
		service.Extensions[composeKubernetesServiceAccountExtensionKey] = deepCopyValue(serviceAccount)
	}
}

func reconcileKubernetesRBACResources(app *Application) {
	if app == nil || app.Extensions == nil || len(app.Services) == 0 {
		return
	}
	rawResources, ok := app.Extensions[kubernetesRBACResourcesExtensionKey].([]interface{})
	if !ok {
		return
	}
	rolesByRef := map[string]map[string]interface{}{}
	var bindings []map[string]interface{}
	for _, raw := range rawResources {
		resource, ok := asMap(raw)
		if !ok {
			continue
		}
		kind := toString(resource["kind"])
		switch kind {
		case "Role", "ClusterRole":
			rolesByRef[kubernetesRBACRoleLookupKey(resource)] = resource
		case "RoleBinding", "ClusterRoleBinding":
			bindings = append(bindings, resource)
		}
	}
	if len(bindings) == 0 {
		return
	}
	for _, binding := range bindings {
		for _, service := range servicesForKubernetesRBACBinding(app, binding) {
			if service == nil {
				continue
			}
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			bundle := kubernetesRBACBundleForService(service)
			bundle["bindings"] = appendUniqueKubernetesResource(bundle["bindings"], binding)
			if role := rolesByRef[kubernetesRBACRoleRefLookupKey(binding)]; role != nil {
				bundle["roles"] = appendUniqueKubernetesResource(bundle["roles"], role)
			}
			service.Extensions[kubernetesRBACExtensionKey] = bundle
			service.Extensions[composeKubernetesRBACExtensionKey] = deepCopyValue(bundle)
		}
	}
}

func servicesForKubernetesRBACBinding(app *Application, binding map[string]interface{}) []*Service {
	if app == nil || len(binding) == 0 {
		return nil
	}
	subjects, _ := binding["subjects"].([]interface{})
	if len(subjects) == 0 {
		return nil
	}
	seen := map[*Service]struct{}{}
	var services []*Service
	for _, subjectValue := range subjects {
		subject, ok := asMap(subjectValue)
		if !ok || toString(subject["kind"]) != "ServiceAccount" {
			continue
		}
		subjectName := toString(subject["name"])
		if subjectName == "" {
			continue
		}
		for _, service := range app.Services {
			if service == nil {
				continue
			}
			if serviceKubernetesServiceAccountName(service) != subjectName {
				continue
			}
			if _, ok := seen[service]; ok {
				continue
			}
			seen[service] = struct{}{}
			services = append(services, service)
		}
	}
	return services
}

func serviceKubernetesServiceAccountName(service *Service) string {
	if service == nil {
		return ""
	}
	if service.ServiceAccountName != "" {
		return service.ServiceAccountName
	}
	if value := extensionStringValue(service, "kubernetes.serviceAccountName"); value != "" {
		return value
	}
	if service.Extensions != nil {
		return toString(service.Extensions["x-kubernetes-serviceAccountName"])
	}
	return ""
}

func kubernetesRBACBundleForService(service *Service) map[string]interface{} {
	bundle := map[string]interface{}{}
	if existing := kubernetesRBACFromService(service); len(existing) > 0 {
		bundle, _ = deepCopyValue(existing).(map[string]interface{})
	}
	if bundle["roles"] == nil {
		bundle["roles"] = []interface{}{}
	}
	if bundle["bindings"] == nil {
		bundle["bindings"] = []interface{}{}
	}
	return bundle
}

func appendUniqueKubernetesResource(existing interface{}, resource map[string]interface{}) []interface{} {
	items, _ := existing.([]interface{})
	key := kubernetesDocumentKeyFromMap(resource)
	for _, item := range items {
		if mapped, ok := asMap(item); ok && key != "" && kubernetesDocumentKeyFromMap(mapped) == key {
			return items
		}
	}
	return append(items, deepCopyValue(resource))
}

func kubernetesRBACRoleLookupKey(resource map[string]interface{}) string {
	metadata, _ := asMap(resource["metadata"])
	return toString(resource["kind"]) + "/" + toString(metadata["name"])
}

func kubernetesRBACRoleRefLookupKey(binding map[string]interface{}) string {
	roleRef, _ := asMap(binding["roleRef"])
	return toString(roleRef["kind"]) + "/" + toString(roleRef["name"])
}

func recordKubernetesServiceMetadata(service *Service, serviceResource map[string]interface{}) {
	if service == nil {
		return
	}
	metadata, _ := asMap(serviceResource["metadata"])
	spec, _ := asMap(serviceResource["spec"])
	serviceName := toString(metadata["name"])
	if service.Extensions == nil {
		service.Extensions = map[string]interface{}{}
	}
	if serviceName != "" {
		service.Extensions["kubernetes.service.name"] = serviceName
	}
	if labels := toStringMapLoose(metadata["labels"]); len(labels) > 0 {
		service.Extensions["kubernetes.service.labels"] = copyStringMap(labels)
	}
	if annotations := toStringMapLoose(metadata["annotations"]); len(annotations) > 0 {
		service.Extensions["kubernetes.service.annotations"] = copyStringMap(annotations)
	}
	for _, key := range []string{
		"type",
		"clusterIP",
		"externalName",
		"sessionAffinity",
		"loadBalancerIP",
		"loadBalancerClass",
		"ipFamilyPolicy",
		"externalTrafficPolicy",
		"internalTrafficPolicy",
		"trafficDistribution",
	} {
		if value := toString(spec[key]); value != "" {
			service.Extensions["kubernetes.service."+key] = value
		}
	}
	for _, key := range []string{"externalIPs", "ipFamilies", "clusterIPs", "loadBalancerSourceRanges"} {
		if values := interfaceSliceToStringSliceLoose(spec[key]); len(values) > 0 {
			service.Extensions["kubernetes.service."+key] = values
		}
	}
	for _, key := range []string{"allocateLoadBalancerNodePorts", "publishNotReadyAddresses"} {
		if value, ok := spec[key].(bool); ok {
			service.Extensions["kubernetes.service."+key] = value
		}
	}
	if value := spec["healthCheckNodePort"]; value != nil {
		service.Extensions["kubernetes.service.healthCheckNodePort"] = value
	}
	if selector := toStringMapLoose(spec["selector"]); len(selector) > 0 {
		service.Extensions["kubernetes.service.selector"] = copyStringMap(selector)
	}
	if sessionAffinityConfig, ok := asMap(spec["sessionAffinityConfig"]); ok && len(sessionAffinityConfig) > 0 {
		service.Extensions["kubernetes.service.sessionAffinityConfig"] = cloneMap(sessionAffinityConfig)
	}
}

func recordKubernetesServiceTarget(app *Application, serviceResource map[string]interface{}, targetName string) {
	metadata, _ := asMap(serviceResource["metadata"])
	serviceName := toString(metadata["name"])
	if app == nil || serviceName == "" || targetName == "" {
		return
	}
	targets := extensionStringMap(app, "kubernetes.serviceTargets")
	targets[serviceName] = targetName
}

func recordKubernetesServicePortTarget(app *Application, serviceResource map[string]interface{}, port PortMapping, target *Service) {
	metadata, _ := asMap(serviceResource["metadata"])
	serviceName := toString(metadata["name"])
	targetPort := resolveKubernetesWorkloadPort(target, port)
	if app == nil || serviceName == "" || targetPort == "" {
		return
	}
	targets := extensionStringMap(app, "kubernetes.servicePortTargets")
	for _, key := range kubernetesServicePortLookupKeys(serviceName, port) {
		targets[key] = targetPort
	}
}

func extensionStringMap(app *Application, key string) map[string]interface{} {
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	if existing, ok := asMap(app.Extensions[key]); ok {
		app.Extensions[key] = existing
		return existing
	}
	created := map[string]interface{}{}
	app.Extensions[key] = created
	return created
}

func kubernetesServicePortLookupKeys(serviceName string, port PortMapping) []string {
	var keys []string
	for _, value := range []string{port.Name, port.HostPort} {
		if value != "" && value != "0" {
			keys = append(keys, serviceName+":"+value)
		}
	}
	return keys
}

func resolveKubernetesWorkloadPort(service *Service, servicePort PortMapping) string {
	if service == nil {
		return ""
	}
	for _, port := range service.Ports {
		if servicePort.TargetName != "" && (port.Name == servicePort.TargetName || port.TargetName == servicePort.TargetName) {
			return port.ContainerPort
		}
		if servicePort.ContainerPort != "" && port.ContainerPort == servicePort.ContainerPort {
			return port.ContainerPort
		}
		if servicePort.Name != "" && port.Name == servicePort.Name {
			return port.ContainerPort
		}
	}
	if servicePort.ContainerPort != "" && parseInt(servicePort.ContainerPort) > 0 {
		return servicePort.ContainerPort
	}
	return ""
}

func targetWorkloadForKubernetesService(app *Application, serviceResource map[string]interface{}) *Service {
	metadata, _ := asMap(serviceResource["metadata"])
	spec, _ := asMap(serviceResource["spec"])
	name := toString(metadata["name"])
	if name != "" {
		if service := app.Services[name]; service != nil {
			return service
		}
	}
	selector := toStringMapLoose(spec["selector"])
	if len(selector) == 0 {
		return nil
	}
	var matches []*Service
	for _, service := range app.Services {
		if labelsMatch(selector, service.Labels) {
			matches = append(matches, service)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	servicePorts := portsFromKubernetesService(serviceResource)
	var matched *Service
	for _, service := range matches {
		if serviceMatchesAnyKubernetesServicePort(service, servicePorts) {
			if matched != nil {
				return nil
			}
			matched = service
		}
	}
	return matched
}

func serviceMatchesAnyKubernetesServicePort(service *Service, servicePorts []PortMapping) bool {
	if service == nil || len(servicePorts) == 0 {
		return false
	}
	for _, servicePort := range servicePorts {
		if resolveKubernetesWorkloadPort(service, servicePort) != "" {
			return true
		}
	}
	return false
}

func kubernetesVolumeSources(podSpec map[string]interface{}) (map[string][]FileRef, map[string][]FileRef, map[string]map[string]interface{}) {
	configs := map[string][]FileRef{}
	secrets := map[string][]FileRef{}
	tmpfs := map[string]map[string]interface{}{}
	volumes, _ := podSpec["volumes"].([]interface{})
	for _, volumeValue := range volumes {
		volume, ok := asMap(volumeValue)
		if !ok {
			continue
		}
		volumeName := toString(volume["name"])
		if volumeName == "" {
			continue
		}
		if configMap, ok := asMap(volume["configMap"]); ok {
			source := toString(configMap["name"])
			if source == "" {
				source = volumeName
			}
			configs[volumeName] = kubernetesProjectedFileRefs(source, configMap, false)
		}
		if secret, ok := asMap(volume["secret"]); ok {
			source := toString(secret["secretName"])
			if source == "" {
				source = volumeName
			}
			secrets[volumeName] = kubernetesProjectedFileRefs(source, secret, true)
		}
		if emptyDir, ok := asMap(volume["emptyDir"]); ok {
			if medium := toString(emptyDir["medium"]); strings.EqualFold(medium, "Memory") {
				info := map[string]interface{}{"medium": "Memory"}
				if sizeLimit := toString(emptyDir["sizeLimit"]); sizeLimit != "" {
					info["sizeLimit"] = sizeLimit
				}
				tmpfs[volumeName] = info
			}
		}
	}
	return configs, secrets, tmpfs
}

func kubernetesProjectedFileRefs(source string, sourceMap map[string]interface{}, readOnly bool) []FileRef {
	optional := boolPtrFromInterface(sourceMap["optional"])
	defaultMode := kubernetesModeString(sourceMap["defaultMode"])
	items, _ := sourceMap["items"].([]interface{})
	if len(items) == 0 {
		return []FileRef{{
			Source:   source,
			Mode:     defaultMode,
			ReadOnly: readOnly,
			Optional: optional,
		}}
	}
	refs := make([]FileRef, 0, len(items))
	for _, itemValue := range items {
		item, ok := asMap(itemValue)
		if !ok {
			continue
		}
		mode := kubernetesModeString(item["mode"])
		if mode == "" {
			mode = defaultMode
		}
		refs = append(refs, FileRef{
			Source:   source,
			Key:      toString(item["key"]),
			Target:   toString(item["path"]),
			Mode:     mode,
			ReadOnly: readOnly,
			Optional: optional,
		})
	}
	if len(refs) == 0 {
		return []FileRef{{Source: source, Mode: defaultMode, ReadOnly: readOnly, Optional: optional}}
	}
	return refs
}

func kubernetesMountedFileTarget(mountPathValue, itemPath string) string {
	if mountPathValue == "" || itemPath == "" {
		return mountPathValue
	}
	if strings.HasPrefix(itemPath, "/") {
		return itemPath
	}
	return path.Join(mountPathValue, itemPath)
}

func kubernetesModeString(value interface{}) string {
	if value == nil {
		return ""
	}
	if str := strings.TrimSpace(toString(value)); str != "" {
		if parsed, err := strconv.ParseInt(str, 0, 64); err == nil {
			return fmt.Sprintf("0%o", parsed)
		}
		return str
	}
	return ""
}

func kubernetesEnvSource(name string, valueFrom map[string]interface{}) EnvSource {
	if ref, ok := asMap(valueFrom["configMapKeyRef"]); ok {
		return EnvSource{
			Name:       name,
			SourceType: "config",
			Source:     toString(ref["name"]),
			Key:        toString(ref["key"]),
			Optional:   toBool(ref["optional"]),
			Extensions: map[string]interface{}{"kubernetes.valueFrom": deepCopyValue(valueFrom)},
		}
	}
	if ref, ok := asMap(valueFrom["secretKeyRef"]); ok {
		return EnvSource{
			Name:       name,
			SourceType: "secret",
			Source:     toString(ref["name"]),
			Key:        toString(ref["key"]),
			Optional:   toBool(ref["optional"]),
			Extensions: map[string]interface{}{"kubernetes.valueFrom": deepCopyValue(valueFrom)},
		}
	}
	if ref, ok := asMap(valueFrom["fieldRef"]); ok {
		return EnvSource{
			Name:       name,
			SourceType: "field",
			Source:     toString(ref["fieldPath"]),
			Extensions: map[string]interface{}{"kubernetes.valueFrom": deepCopyValue(valueFrom)},
		}
	}
	if ref, ok := asMap(valueFrom["resourceFieldRef"]); ok {
		return EnvSource{
			Name:       name,
			SourceType: "resource",
			Source:     toString(ref["resource"]),
			Extensions: map[string]interface{}{"kubernetes.valueFrom": deepCopyValue(valueFrom)},
		}
	}
	return EnvSource{}
}

func kubernetesEnvFromSource(envFrom map[string]interface{}) EnvFromSource {
	prefix := toString(envFrom["prefix"])
	if ref, ok := asMap(envFrom["configMapRef"]); ok {
		return EnvFromSource{
			SourceType: "config",
			Source:     toString(ref["name"]),
			Prefix:     prefix,
			Optional:   toBool(ref["optional"]),
			Extensions: map[string]interface{}{"kubernetes.envFrom": deepCopyValue(envFrom)},
		}
	}
	if ref, ok := asMap(envFrom["secretRef"]); ok {
		return EnvFromSource{
			SourceType: "secret",
			Source:     toString(ref["name"]),
			Prefix:     prefix,
			Optional:   toBool(ref["optional"]),
			Extensions: map[string]interface{}{"kubernetes.envFrom": deepCopyValue(envFrom)},
		}
	}
	return EnvFromSource{}
}

func parseKubernetesDependencies(service *Service, annotations map[string]string, podSpec map[string]interface{}) {
	if encoded := annotations[kubernetesAnnotationDependencies]; encoded != "" {
		var dependencies []DependencySpec
		if err := json.Unmarshal([]byte(encoded), &dependencies); err == nil {
			for _, dependency := range dependencies {
				service.Dependencies = appendUniqueDependency(service.Dependencies, dependency)
				service.DependsOn = appendUniqueName(service.DependsOn, dependency.Name)
			}
		}
	}
	initContainers, _ := podSpec["initContainers"].([]interface{})
	for _, initContainerValue := range initContainers {
		initContainer, ok := asMap(initContainerValue)
		if !ok {
			continue
		}
		name := toString(initContainer["name"])
		if !strings.HasPrefix(name, "wait-for-") {
			continue
		}
		dependencyName := strings.TrimPrefix(name, "wait-for-")
		if dependencyName == "" {
			continue
		}
		service.Dependencies = appendUniqueDependency(service.Dependencies, DependencySpec{
			Name:      dependencyName,
			Condition: "service_started",
		})
		service.DependsOn = appendUniqueName(service.DependsOn, dependencyName)
	}
}

func labelsMatch(selector, labels map[string]string) bool {
	if len(selector) == 0 {
		return false
	}
	for key, expected := range selector {
		if labels[key] != expected {
			return false
		}
	}
	return true
}

func portsFromKubernetesService(serviceResource map[string]interface{}) []PortMapping {
	spec, _ := asMap(serviceResource["spec"])
	ports, _ := spec["ports"].([]interface{})
	var mappings []PortMapping
	for _, portValue := range ports {
		port, ok := asMap(portValue)
		if !ok {
			continue
		}
		hostPort := toString(port["port"])
		targetPort := toString(port["targetPort"])
		targetName := ""
		if parseInt(targetPort) == 0 && targetPort != "" {
			targetName = targetPort
		}
		if targetPort == "" || targetPort == "0" {
			targetPort = hostPort
		}
		if hostPort == "" || hostPort == "0" || targetPort == "" || targetPort == "0" {
			continue
		}
		protocol := strings.ToLower(toString(port["protocol"]))
		if protocol == "" {
			protocol = "tcp"
		}
		mappings = append(mappings, PortMapping{
			Name:          toString(port["name"]),
			TargetName:    targetName,
			HostPort:      hostPort,
			ContainerPort: targetPort,
			NodePort:      toString(port["nodePort"]),
			Protocol:      protocol,
			AppProtocol:   toString(port["appProtocol"]),
			Extensions:    kubernetesPortExtensions(port, []string{"name", "port", "targetPort", "nodePort", "protocol", "appProtocol"}),
		})
	}
	return mappings
}

func portablePortsNeedKubernetesAnnotation(ports []PortMapping) bool {
	for _, port := range ports {
		if port.HostIP != "" || port.TargetName != "" || port.Mode != "" || len(port.Extensions) > 0 {
			return true
		}
	}
	return false
}

func mergeServicePort(service *Service, port PortMapping) {
	if service == nil || port.HostPort == "" || port.ContainerPort == "" {
		return
	}
	for index, existing := range service.Ports {
		if portsReferToSameTarget(existing, port) {
			if service.Ports[index].TargetName == "" && port.TargetName != "" && service.Ports[index].Name == port.TargetName {
				service.Ports[index].TargetName = service.Ports[index].Name
			}
			if port.Name != "" {
				service.Ports[index].Name = port.Name
			}
			if port.AppProtocol != "" {
				service.Ports[index].AppProtocol = port.AppProtocol
			}
			if port.NodePort != "" {
				service.Ports[index].NodePort = port.NodePort
			}
			if len(port.Extensions) > 0 {
				service.Ports[index].Extensions = mergeInterfaceMaps(service.Ports[index].Extensions, port.Extensions)
			}
			if service.Ports[index].TargetName == "" {
				service.Ports[index].TargetName = port.TargetName
			}
			if service.Ports[index].HostPort == "" {
				service.Ports[index].HostPort = port.HostPort
			}
			if service.Ports[index].ContainerPort == "" || service.Ports[index].ContainerPort == service.Ports[index].TargetName {
				service.Ports[index].ContainerPort = port.ContainerPort
			}
			if service.Ports[index].Protocol == "" {
				service.Ports[index].Protocol = port.Protocol
			}
			return
		}
	}
	service.Ports = append(service.Ports, port)
}

func portsReferToSameTarget(existing, candidate PortMapping) bool {
	if candidate.TargetName != "" && existing.Name == candidate.TargetName {
		return true
	}
	if candidate.Name != "" && existing.Name == candidate.Name {
		return true
	}
	if candidate.ContainerPort != "" && existing.ContainerPort == candidate.ContainerPort {
		return true
	}
	return candidate.HostPort != "" && existing.HostPort == candidate.HostPort
}

func storeKubernetesOpaqueManifest(app *Application, resource map[string]interface{}) {
	if app == nil || len(resource) == 0 {
		return
	}
	if app.KubernetesOpaqueManifests == nil {
		app.KubernetesOpaqueManifests = map[string]*KubernetesOpaqueManifestSpec{}
	}
	if spec := kubernetesOpaqueManifestSpecFromMap(resource); spec != nil && spec.Name != "" {
		key := kubernetesDocumentKeyFromMap(spec.Raw)
		if key == "" {
			key = strings.Join([]string{spec.Kind, spec.Namespace, spec.Name}, "/")
		}
		if key != "" {
			app.KubernetesOpaqueManifests[key] = spec
		}
	}
}

func parseKubernetesConfigMap(app *Application, resource map[string]interface{}) error {
	if app != nil && app.Extensions != nil {
		app.Extensions[kubernetesConfigMapsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesConfigMapsExtensionKey], deepCopyValue(resource))
		app.Extensions[composeKubernetesConfigMapsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesConfigMapsExtensionKey], deepCopyValue(resource))
	}
	storeKubernetesOpaqueManifest(app, resource)
	metadata, _ := asMap(resource["metadata"])
	data, _ := asMap(resource["data"])

	name, _ := metadata["name"].(string)
	if name == "" {
		return fmt.Errorf("configmap missing name")
	}
	if parseKubernetesPortablePolicyConfigMap(app, metadata, data) {
		return nil
	}

	config := &Config{
		Name:       name,
		Extensions: map[string]interface{}{},
	}
	if labels := toStringMapLoose(metadata["labels"]); len(labels) > 0 {
		config.Extensions["kubernetes.labels"] = copyStringMap(labels)
		config.Extensions["x-kubernetes-labels"] = copyStringMap(labels)
	}
	if annotations := toStringMapLoose(metadata["annotations"]); len(annotations) > 0 {
		config.Extensions["kubernetes.annotations"] = copyStringMap(annotations)
		config.Extensions["x-kubernetes-annotations"] = copyStringMap(annotations)
	}
	if immutable, ok := resource["immutable"].(bool); ok {
		config.Extensions["kubernetes.immutable"] = immutable
		config.Extensions["x-kubernetes-immutable"] = immutable
	}

	if len(data) > 0 {
		configData := toStringMapLoose(data)
		config.Extensions["kubernetes.data"] = copyStringMap(configData)
		config.Extensions["x-kubernetes-data"] = copyStringMap(configData)
		var content strings.Builder
		for _, key := range sortedMapKeys(configData) {
			content.WriteString(fmt.Sprintf("%s: %s\n", key, configData[key]))
		}
		config.Content = content.String()
	}
	if binaryData := toStringMapLoose(resource["binaryData"]); len(binaryData) > 0 {
		config.Extensions["kubernetes.binaryData"] = copyStringMap(binaryData)
		config.Extensions["x-kubernetes-binaryData"] = copyStringMap(binaryData)
	}

	app.Configs[name] = config
	return nil
}

func parseKubernetesPortablePolicyConfigMap(app *Application, metadata, data map[string]interface{}) bool {
	labels := toStringMapLoose(metadata["labels"])
	policyName := labels["bolabaden.dev/portable-policy"]
	if policyName == "" {
		return false
	}
	raw := toString(data["policy.json"])
	if raw == "" {
		return true
	}
	var policy PolicySpec
	if err := json.Unmarshal([]byte(raw), &policy); err != nil {
		return true
	}
	if policy.Name == "" {
		policy.Name = policyName
	}
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	policies, _ := app.Extensions[composeAppPoliciesExtension].([]interface{})
	app.Extensions[composeAppPoliciesExtension] = append(policies, clonePortablePolicySpec(&policy))
	return true
}

func parseKubernetesSecret(app *Application, resource map[string]interface{}) error {
	if app != nil && app.Extensions != nil {
		app.Extensions[kubernetesSecretsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesSecretsExtensionKey], deepCopyValue(resource))
		app.Extensions[composeKubernetesSecretsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesSecretsExtensionKey], deepCopyValue(resource))
	}
	storeKubernetesOpaqueManifest(app, resource)
	metadata, _ := asMap(resource["metadata"])
	data, _ := asMap(resource["data"])

	name, _ := metadata["name"].(string)
	if name == "" {
		return fmt.Errorf("secret missing name")
	}

	secret := &Secret{
		Name:       name,
		Extensions: map[string]interface{}{},
	}
	if labels := toStringMapLoose(metadata["labels"]); len(labels) > 0 {
		secret.Extensions["kubernetes.labels"] = copyStringMap(labels)
		secret.Extensions["x-kubernetes-labels"] = copyStringMap(labels)
	}
	if annotations := toStringMapLoose(metadata["annotations"]); len(annotations) > 0 {
		secret.Extensions["kubernetes.annotations"] = copyStringMap(annotations)
		secret.Extensions["x-kubernetes-annotations"] = copyStringMap(annotations)
	}
	if immutable, ok := resource["immutable"].(bool); ok {
		secret.Extensions["kubernetes.immutable"] = immutable
		secret.Extensions["x-kubernetes-immutable"] = immutable
	}
	if secretType := toString(resource["type"]); secretType != "" {
		secret.Extensions["kubernetes.type"] = secretType
		secret.Extensions["x-kubernetes-type"] = secretType
	}

	if len(data) > 0 {
		secretData := toStringMapLoose(data)
		secret.Extensions["kubernetes.data"] = copyStringMap(secretData)
		secret.Extensions["x-kubernetes-data"] = copyStringMap(secretData)
		keys := sortedMapKeys(secretData)
		if len(keys) > 0 {
			secret.Environment = strings.Join(keys, ",")
		}
	}
	if stringData := toStringMapLoose(resource["stringData"]); len(stringData) > 0 {
		secret.Extensions["kubernetes.stringData"] = copyStringMap(stringData)
		secret.Extensions["x-kubernetes-stringData"] = copyStringMap(stringData)
	}

	app.Secrets[name] = secret
	return nil
}

func parseKubernetesNamespace(app *Application, resource map[string]interface{}) error {
	storeKubernetesOpaqueManifest(app, resource)
	metadata, _ := asMap(resource["metadata"])
	name, _ := metadata["name"].(string)
	if name == "" {
		return fmt.Errorf("namespace missing name")
	}
	app.Namespace = name
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	app.Extensions["kubernetes.namespace"] = name
	app.Extensions[kubernetesNamespacesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesNamespacesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesNamespacesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesNamespacesExtensionKey], deepCopyValue(resource))
	return nil
}

func parseKubernetesPVC(app *Application, resource map[string]interface{}) error {
	storeKubernetesOpaqueManifest(app, resource)
	metadata, _ := asMap(resource["metadata"])
	spec, _ := asMap(resource["spec"])

	name, _ := metadata["name"].(string)
	if name == "" {
		return fmt.Errorf("pvc missing name")
	}
	app.Extensions[kubernetesPVCsExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesPVCsExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesPVCsExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesPVCsExtensionKey], deepCopyValue(resource))

	volume := &Volume{
		Name:       name,
		Driver:     "persistentVolumeClaim",
		DriverOpts: map[string]string{},
		Extensions: map[string]interface{}{
			"kubernetes.kind":   "PersistentVolumeClaim",
			"x-kubernetes-kind": "PersistentVolumeClaim",
		},
	}
	if labels := toStringMapLoose(metadata["labels"]); len(labels) > 0 {
		volume.Extensions["kubernetes.labels"] = copyStringMap(labels)
		volume.Extensions["x-kubernetes-labels"] = copyStringMap(labels)
	}
	if annotations := toStringMapLoose(metadata["annotations"]); len(annotations) > 0 {
		volume.Extensions["kubernetes.annotations"] = copyStringMap(annotations)
		volume.Extensions["x-kubernetes-annotations"] = copyStringMap(annotations)
	}

	if accessModes := interfaceSliceToStringSliceLoose(spec["accessModes"]); len(accessModes) > 0 {
		volume.DriverOpts["accessModes"] = strings.Join(accessModes, ",")
	}
	if resources, ok := asMap(spec["resources"]); ok {
		if requests, ok := asMap(resources["requests"]); ok {
			if storage, ok := requests["storage"].(string); ok {
				volume.DriverOpts["storage"] = storage
			}
		}
	}
	for _, key := range []string{"storageClassName", "volumeName", "volumeMode"} {
		if value := toString(spec[key]); value != "" {
			volume.DriverOpts[key] = value
		}
	}
	if selector, ok := asMap(spec["selector"]); ok && len(selector) > 0 {
		volume.Extensions["kubernetes.selector"] = cloneMap(selector)
		volume.Extensions["x-kubernetes-selector"] = cloneMap(selector)
	}
	if dataSource, ok := asMap(spec["dataSource"]); ok && len(dataSource) > 0 {
		volume.Extensions["kubernetes.dataSource"] = cloneMap(dataSource)
		volume.Extensions["x-kubernetes-dataSource"] = cloneMap(dataSource)
	}
	if dataSourceRef, ok := asMap(spec["dataSourceRef"]); ok && len(dataSourceRef) > 0 {
		volume.Extensions["kubernetes.dataSourceRef"] = cloneMap(dataSourceRef)
		volume.Extensions["x-kubernetes-dataSourceRef"] = cloneMap(dataSourceRef)
	}

	app.Volumes[name] = volume
	return nil
}

func parseKubernetesPersistentVolume(app *Application, resource map[string]interface{}) error {
	storeKubernetesOpaqueManifest(app, resource)
	metadata, _ := asMap(resource["metadata"])
	spec, _ := asMap(resource["spec"])

	name, _ := metadata["name"].(string)
	if name == "" {
		return fmt.Errorf("persistentvolume missing name")
	}
	app.Extensions[kubernetesPersistentVolumesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesPersistentVolumesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesPersistentVolumesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesPersistentVolumesExtensionKey], deepCopyValue(resource))

	volume := &Volume{
		Name:       name,
		Driver:     "persistentVolume",
		Extensions: map[string]interface{}{"kubernetes.kind": "PersistentVolume", "x-kubernetes-kind": "PersistentVolume"},
	}
	volume.DriverOpts = map[string]string{}
	if labels := toStringMapLoose(metadata["labels"]); len(labels) > 0 {
		volume.Extensions["kubernetes.labels"] = copyStringMap(labels)
		volume.Extensions["x-kubernetes-labels"] = copyStringMap(labels)
	}
	if annotations := toStringMapLoose(metadata["annotations"]); len(annotations) > 0 {
		volume.Extensions["kubernetes.annotations"] = copyStringMap(annotations)
		volume.Extensions["x-kubernetes-annotations"] = copyStringMap(annotations)
	}
	if claimRef, ok := asMap(spec["claimRef"]); ok {
		if claimName := toString(claimRef["name"]); claimName != "" {
			volume.DriverOpts["claimRef"] = claimName
		}
		if claimNamespace := toString(claimRef["namespace"]); claimNamespace != "" {
			volume.DriverOpts["claimRefNamespace"] = claimNamespace
		}
		if claimKind := toString(claimRef["kind"]); claimKind != "" {
			volume.DriverOpts["claimRefKind"] = claimKind
		}
		if claimAPIVersion := toString(claimRef["apiVersion"]); claimAPIVersion != "" {
			volume.DriverOpts["claimRefAPIVersion"] = claimAPIVersion
		}
		if claimUID := toString(claimRef["uid"]); claimUID != "" {
			volume.DriverOpts["claimRefUID"] = claimUID
		}
		if claimResourceVersion := toString(claimRef["resourceVersion"]); claimResourceVersion != "" {
			volume.DriverOpts["claimRefResourceVersion"] = claimResourceVersion
		}
	}
	if accessModes := interfaceSliceToStringSliceLoose(spec["accessModes"]); len(accessModes) > 0 {
		volume.DriverOpts["accessModes"] = strings.Join(accessModes, ",")
	}
	if capacity, ok := asMap(spec["capacity"]); ok {
		if storage := toString(capacity["storage"]); storage != "" {
			volume.DriverOpts["capacity.storage"] = storage
		}
	}
	for _, key := range []string{"storageClassName", "volumeMode", "persistentVolumeReclaimPolicy", "mountOptions"} {
		if values := interfaceSliceToStringSliceLoose(spec[key]); len(values) > 0 {
			volume.DriverOpts[key] = strings.Join(values, ",")
			continue
		}
		if value := toString(spec[key]); value != "" {
			volume.DriverOpts[key] = value
		}
	}
	switch {
	case func() bool { _, ok := asMap(spec["hostPath"]); return ok }():
		hostPath, _ := asMap(spec["hostPath"])
		volume.Driver = "hostPath"
		if path := toString(hostPath["path"]); path != "" {
			volume.DriverOpts["path"] = path
		}
	case func() bool { _, ok := asMap(spec["nfs"]); return ok }():
		nfs, _ := asMap(spec["nfs"])
		volume.Driver = "nfs"
		if server := toString(nfs["server"]); server != "" {
			volume.DriverOpts["server"] = server
		}
		if path := toString(nfs["path"]); path != "" {
			volume.DriverOpts["path"] = path
		}
	case func() bool { _, ok := asMap(spec["csi"]); return ok }():
		csi, _ := asMap(spec["csi"])
		volume.Driver = "csi"
		if driver := toString(csi["driver"]); driver != "" {
			volume.DriverOpts["driver"] = driver
		}
		if handle := toString(csi["volumeHandle"]); handle != "" {
			volume.DriverOpts["volumeHandle"] = handle
		}
		for _, key := range []string{"fsType", "readOnly"} {
			if value := toString(csi[key]); value != "" {
				volume.DriverOpts["csi."+key] = value
			}
		}
		if attributes := toStringMapLoose(csi["volumeAttributes"]); len(attributes) > 0 {
			volume.Extensions["kubernetes.csi.volumeAttributes"] = copyStringMap(attributes)
			volume.Extensions["x-kubernetes-csi-volumeAttributes"] = copyStringMap(attributes)
		}
	case func() bool { _, ok := asMap(spec["local"]); return ok }():
		local, _ := asMap(spec["local"])
		volume.Driver = "local"
		if path := toString(local["path"]); path != "" {
			volume.DriverOpts["path"] = path
		}
	}
	if nodeAffinity, ok := asMap(spec["nodeAffinity"]); ok && len(nodeAffinity) > 0 {
		volume.Extensions["kubernetes.nodeAffinity"] = cloneMap(nodeAffinity)
		volume.Extensions["x-kubernetes-nodeAffinity"] = cloneMap(nodeAffinity)
	}

	app.Volumes[name] = volume
	return nil
}

func parseKubernetesIngress(app *Application, resource map[string]interface{}) error {
	storeKubernetesOpaqueManifest(app, resource)
	spec, _ := asMap(resource["spec"])
	app.Extensions["kubernetes.routes"] = appendExtensionSlice(app.Extensions["kubernetes.routes"], resource)
	app.Extensions[kubernetesIngressesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesIngressesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesIngressesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesIngressesExtensionKey], deepCopyValue(resource))

	if rules, ok := spec["rules"].([]interface{}); ok {
		for _, rule := range rules {
			if ruleMap, ok := asMap(rule); ok {
				if host, ok := ruleMap["host"].(string); ok {
					// Store ingress info in extensions
					if app.Extensions == nil {
						app.Extensions = make(map[string]interface{})
					}
					if ingress, ok := app.Extensions["ingress"].([]string); ok {
						app.Extensions["ingress"] = append(ingress, host)
					} else {
						app.Extensions["ingress"] = []string{host}
					}
				}
			}
		}
	}

	return nil
}

func parseKubernetesHTTPRoute(app *Application, resource map[string]interface{}) error {
	if app == nil || len(resource) == 0 {
		return nil
	}
	spec, _ := asMap(resource["spec"])
	metadata, _ := asMap(resource["metadata"])
	name := toString(metadata["name"])
	if name == "" {
		return nil
	}

	app.Extensions["kubernetes.routes"] = appendExtensionSlice(app.Extensions["kubernetes.routes"], resource)

	route := &RouteSpec{
		Name:        name,
		Protocol:    "http",
		Source:      app.Platform,
		Raw:         resource,
		Metadata:    toStringMapLoose(metadata["labels"]),
		Annotations: toStringMapLoose(metadata["annotations"]),
	}
	if route.Metadata == nil {
		route.Metadata = map[string]string{}
	}
	route.Metadata["kubernetes.kind"] = "HTTPRoute"

	if portableRoute := portableRouteFromKubernetesAnnotation(route.Annotations[kubernetesPortableRouteAnnotation]); portableRoute != nil {
		mergePortableRouteIntoKubernetesRoute(route, portableRoute)
	}

	if hostnames, ok := spec["hostnames"].([]interface{}); ok {
		for _, host := range hostnames {
			if value := toString(host); value != "" {
				route.Hosts = appendUnique(route.Hosts, value)
			}
		}
	}
	if rules, ok := spec["rules"].([]interface{}); ok {
		for _, ruleValue := range rules {
			rule, ok := asMap(ruleValue)
			if !ok {
				continue
			}
			matches, _ := rule["matches"].([]interface{})
			for _, matchValue := range matches {
				match, ok := asMap(matchValue)
				if !ok {
					continue
				}
				if pathSpec, ok := asMap(match["path"]); ok {
					if path := toString(pathSpec["value"]); path != "" {
						route.Paths = appendUnique(route.Paths, path)
					}
				}
			}
			backendRefs, _ := rule["backendRefs"].([]interface{})
			for _, backendValue := range backendRefs {
				backend, ok := asMap(backendValue)
				if !ok {
					continue
				}
				if service := toString(backend["name"]); service != "" && route.Service == "" {
					route.Service = service
					route.Metadata["kubernetes.service"] = service
				}
				if port := toString(backend["port"]); port != "" && route.Port == "" {
					route.Port = port
					route.Metadata["kubernetes.servicePort"] = port
				}
			}
		}
	}

	resolveKubernetesRouteBackend(app, route)
	canonical := canonicalForApplication(app)
	if canonical != nil {
		canonical.AddRoute(route)
	}
	return nil
}

func parseKubernetesPolicy(app *Application, resource map[string]interface{}) error {
	storeKubernetesOpaqueManifest(app, resource)
	app.Extensions["kubernetes.policies"] = appendExtensionSlice(app.Extensions["kubernetes.policies"], resource)
	app.Extensions[kubernetesNetworkPoliciesExtensionKey] = appendExtensionSlice(app.Extensions[kubernetesNetworkPoliciesExtensionKey], deepCopyValue(resource))
	app.Extensions[composeKubernetesNetworkPoliciesExtensionKey] = appendExtensionSlice(app.Extensions[composeKubernetesNetworkPoliciesExtensionKey], deepCopyValue(resource))
	return nil
}

// SerializeKubernetesYAML converts an Application to Kubernetes YAML
func SerializeKubernetesYAML(app *Application) (string, error) {
	emitApp := cloneApplication(app)
	syncPortableApplicationState(emitApp)
	var documents []string

	if namespaceDocs, err := serializeKubernetesNamespaces(emitApp); err != nil {
		return "", err
	} else if len(namespaceDocs) > 0 {
		documents = append(documents, namespaceDocs...)
	}

	workloadDocs, err := serializeKubernetesWorkloadResources(emitApp)
	if err != nil {
		return "", err
	}
	documents = append(documents, workloadDocs...)

	configMapDocs, err := serializeKubernetesConfigMapResources(emitApp)
	if err != nil {
		return "", err
	}
	documents = append(documents, configMapDocs...)

	secretDocs, err := serializeKubernetesSecretResources(emitApp)
	if err != nil {
		return "", err
	}
	documents = append(documents, secretDocs...)

	// Create ConfigMaps
	for name, config := range emitApp.Configs {
		if hasKubernetesTypedResource(emitApp, "ConfigMap", name) {
			continue
		}
		doc, err := serializeKubernetesConfigMap(name, config, emitApp.Namespace)
		if err != nil {
			return "", fmt.Errorf("failed to serialize configmap %s: %w", name, err)
		}
		documents = append(documents, doc)
	}

	// Create Secrets
	for name, secret := range emitApp.Secrets {
		if hasKubernetesTypedResource(emitApp, "Secret", name) {
			continue
		}
		doc, err := serializeKubernetesSecret(name, secret, emitApp.Namespace)
		if err != nil {
			return "", fmt.Errorf("failed to serialize secret %s: %w", name, err)
		}
		documents = append(documents, doc)
	}

	serviceAccountDocs, err := serializeKubernetesServiceAccounts(emitApp)
	if err != nil {
		return "", err
	}
	documents = append(documents, serviceAccountDocs...)

	rbacDocs, err := serializeKubernetesRBACResources(emitApp)
	if err != nil {
		return "", err
	}
	documents = append(documents, rbacDocs...)

	namespacePolicyDocs, err := serializeKubernetesNamespacePolicyResources(emitApp)
	if err != nil {
		return "", err
	}
	documents = append(documents, namespacePolicyDocs...)

	// Create PersistentVolumeClaims
	for name, volume := range emitApp.Volumes {
		if volume == nil {
			continue
		}
		switch {
		case volumeKubernetesKind(volume) == "PersistentVolume" || volume.Driver == "persistentVolume":
			if hasKubernetesTypedResource(emitApp, "PersistentVolume", name) {
				continue
			}
			doc, err := serializeKubernetesPersistentVolume(name, volume)
			if err != nil {
				return "", fmt.Errorf("failed to serialize persistentvolume %s: %w", name, err)
			}
			documents = append(documents, doc)
		case volumeKubernetesKind(volume) == "PersistentVolumeClaim" || volume.Driver == "persistentVolumeClaim":
			if hasKubernetesTypedResource(emitApp, "PersistentVolumeClaim", name) {
				continue
			}
			doc, err := serializeKubernetesPVC(name, volume, emitApp.Namespace)
			if err != nil {
				return "", fmt.Errorf("failed to serialize pvc %s: %w", name, err)
			}
			documents = append(documents, doc)
		}
	}

	// Create Deployments and Services
	for name, service := range emitApp.Services {
		if isKubernetesInitContainerService(service) {
			continue
		}
		// Deployment
		if !hasKubernetesTypedWorkloadResource(emitApp, name) {
			deploymentDoc, err := serializeKubernetesDeployment(emitApp, name, service, emitApp.Namespace)
			if err != nil {
				return "", fmt.Errorf("failed to serialize deployment %s: %w", name, err)
			}
			documents = append(documents, deploymentDoc)
		}

		// Service (if it has ports)
		serviceType := extensionStringValue(service, "kubernetes.service.type")
		if len(service.Ports) > 0 || serviceType == "ExternalName" {
			serviceName := extensionStringValue(service, "kubernetes.service.name")
			if serviceName == "" {
				serviceName = name
			}
			if !hasKubernetesTypedResource(emitApp, "Service", serviceName) {
				serviceDoc, err := serializeKubernetesService(name, service, emitApp.Namespace)
				if err != nil {
					return "", fmt.Errorf("failed to serialize service %s: %w", name, err)
				}
				documents = append(documents, serviceDoc)
			}
		}

		hpa := kubernetesHorizontalPodAutoscalerFromService(service)
		if len(hpa) > 0 {
			hpaName := kubernetesResourceNameOrFallback(hpa, name+"-hpa")
			if !hasKubernetesTypedResource(emitApp, "HorizontalPodAutoscaler", hpaName) {
				hpaDoc, err := serializeKubernetesHorizontalPodAutoscaler(name, service, emitApp.Namespace)
				if err != nil {
					return "", fmt.Errorf("failed to serialize horizontalpodautoscaler %s: %w", name, err)
				}
				if strings.TrimSpace(hpaDoc) != "" {
					documents = append(documents, hpaDoc)
				}
			}
		}

		pdb := kubernetesPodDisruptionBudgetFromService(service)
		if len(pdb) > 0 {
			pdbName := kubernetesResourceNameOrFallback(pdb, name+"-pdb")
			if !hasKubernetesTypedResource(emitApp, "PodDisruptionBudget", pdbName) {
				pdbDoc, err := serializeKubernetesPodDisruptionBudget(name, service, emitApp.Namespace)
				if err != nil {
					return "", fmt.Errorf("failed to serialize poddisruptionbudget %s: %w", name, err)
				}
				if strings.TrimSpace(pdbDoc) != "" {
					documents = append(documents, pdbDoc)
				}
			}
		}
	}

	canonical := canonicalForApplication(emitApp)
	if canonical != nil {
		for _, route := range canonical.Routes {
			if hasKubernetesTypedResource(emitApp, "Ingress", route.Name) {
				continue
			}
			doc, err := serializeKubernetesRoute(route, emitApp)
			if err != nil {
				return "", fmt.Errorf("failed to serialize ingress %s: %w", route.Name, err)
			}
			if strings.TrimSpace(doc) != "" {
				documents = append(documents, doc)
			}
		}
		for _, policy := range canonical.Policies {
			if hasKubernetesTypedResource(emitApp, "NetworkPolicy", policy.Name) {
				continue
			}
			doc, err := serializeKubernetesPolicy(policy, emitApp.Namespace)
			if err != nil {
				return "", fmt.Errorf("failed to serialize network policy %s: %w", policy.Name, err)
			}
			if strings.TrimSpace(doc) != "" {
				documents = append(documents, doc)
			}
		}
	}

	documents, err = appendRawKubernetesResources(documents, emitApp)
	if err != nil {
		return "", err
	}

	return strings.Join(documents, "\n---\n"), nil
}

// SerializeKubernetesJSON converts an Application to Kubernetes JSON.
func SerializeKubernetesJSON(app *Application) (string, error) {
	yamlContent, err := SerializeKubernetesYAML(app)
	if err != nil {
		return "", err
	}
	docs := splitKubernetesSourceDocuments(yamlContent)
	if len(docs) == 0 {
		return "", fmt.Errorf("application does not contain renderable Kubernetes manifests")
	}
	values := make([]interface{}, 0, len(docs))
	for _, doc := range docs {
		var value interface{}
		if err := yaml.Unmarshal([]byte(doc), &value); err != nil {
			return "", fmt.Errorf("failed to convert Kubernetes YAML to JSON: %w", err)
		}
		values = append(values, value)
	}
	if len(values) == 1 {
		jsonBytes, err := json.Marshal(values[0])
		if err != nil {
			return "", fmt.Errorf("failed to marshal Kubernetes JSON: %w", err)
		}
		return string(jsonBytes), nil
	}
	jsonBytes, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Kubernetes JSON: %w", err)
	}
	return string(jsonBytes), nil
}

func serializeKubernetesNamespaces(app *Application) ([]string, error) {
	if app == nil {
		return nil, nil
	}
	resources := kubernetesExtensionResourceSlice(app.Extensions[kubernetesNamespacesExtensionKey], app.Extensions[composeKubernetesNamespacesExtensionKey])
	var documents []string
	emitted := map[string]struct{}{}
	for _, resource := range resources {
		normalized := normalizeKubernetesNamespaceResource(resource)
		key := kubernetesDocumentKeyFromMap(normalized)
		if key != "" {
			if _, ok := emitted[key]; ok {
				continue
			}
			emitted[key] = struct{}{}
		}
		doc, err := marshalYAMLString(normalized)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize namespace resource %s: %w", key, err)
		}
		documents = append(documents, doc)
	}
	namespace := strings.TrimSpace(app.Namespace)
	if namespace == "" {
		if value, ok := applicationExtensionValueForKey(app, "kubernetes.namespace"); ok {
			namespace = strings.TrimSpace(toString(value))
		}
	}
	if namespace == "" {
		namespace = strings.TrimSpace(toString(app.Extensions[composeKubernetesNamespaceExtension]))
	}
	if namespace == "" || namespaceResourceEmitted(emitted, namespace) {
		return documents, nil
	}
	resource := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name": namespace,
		},
	}
	data, err := yaml.Marshal(resource)
	if err != nil {
		return documents, fmt.Errorf("failed to serialize namespace %s: %w", namespace, err)
	}
	return append(documents, string(data)), nil
}

func normalizeKubernetesNamespaceResource(resource map[string]interface{}) map[string]interface{} {
	normalized, _ := deepCopyValue(resource).(map[string]interface{})
	if toString(normalized["apiVersion"]) == "" {
		normalized["apiVersion"] = "v1"
	}
	normalized["kind"] = "Namespace"
	if metadata, _ := asMap(normalized["metadata"]); metadata == nil {
		normalized["metadata"] = map[string]interface{}{}
	}
	return normalized
}

func namespaceResourceEmitted(emitted map[string]struct{}, namespace string) bool {
	for key := range emitted {
		parts := strings.Split(key, "/")
		if len(parts) == 4 && parts[1] == "Namespace" && parts[3] == namespace {
			return true
		}
	}
	return false
}

func serializeKubernetesConfigMapResources(app *Application) ([]string, error) {
	return serializeKubernetesExactResourceSlice(
		app,
		kubernetesConfigMapsExtensionKey,
		composeKubernetesConfigMapsExtensionKey,
		normalizeKubernetesConfigMapResource,
		"configmap",
	)
}

func serializeKubernetesSecretResources(app *Application) ([]string, error) {
	return serializeKubernetesExactResourceSlice(
		app,
		kubernetesSecretsExtensionKey,
		composeKubernetesSecretsExtensionKey,
		normalizeKubernetesSecretResource,
		"secret",
	)
}

func serializeKubernetesExactResourceSlice(app *Application, primaryKey, composeKey string, normalize func(map[string]interface{}, string) map[string]interface{}, label string) ([]string, error) {
	if app == nil {
		return nil, nil
	}
	resources := kubernetesExtensionResourceSlice(app.Extensions[primaryKey], app.Extensions[composeKey])
	var documents []string
	emitted := map[string]struct{}{}
	for _, resource := range resources {
		normalized := normalize(resource, app.Namespace)
		key := kubernetesDocumentKeyFromMap(normalized)
		if key != "" {
			if _, ok := emitted[key]; ok {
				continue
			}
			emitted[key] = struct{}{}
		}
		doc, err := marshalYAMLString(normalized)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize Kubernetes %s resource %s: %w", label, key, err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func normalizeKubernetesConfigMapResource(resource map[string]interface{}, namespace string) map[string]interface{} {
	normalized, _ := deepCopyValue(resource).(map[string]interface{})
	if toString(normalized["apiVersion"]) == "" {
		normalized["apiVersion"] = "v1"
	}
	normalized["kind"] = "ConfigMap"
	metadata, _ := asMap(normalized["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	if namespace != "" && toString(metadata["namespace"]) == "" {
		metadata["namespace"] = namespace
	}
	normalized["metadata"] = metadata
	return normalized
}

func normalizeKubernetesSecretResource(resource map[string]interface{}, namespace string) map[string]interface{} {
	normalized, _ := deepCopyValue(resource).(map[string]interface{})
	if toString(normalized["apiVersion"]) == "" {
		normalized["apiVersion"] = "v1"
	}
	normalized["kind"] = "Secret"
	metadata, _ := asMap(normalized["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	if namespace != "" && toString(metadata["namespace"]) == "" {
		metadata["namespace"] = namespace
	}
	normalized["metadata"] = metadata
	return normalized
}

func serializeKubernetesWorkloadResources(app *Application) ([]string, error) {
	if app == nil {
		return nil, nil
	}
	resources := kubernetesExtensionResourceSlice(app.Extensions[kubernetesWorkloadsExtensionKey], app.Extensions[composeKubernetesWorkloadsExtensionKey])
	var documents []string
	emitted := map[string]struct{}{}
	for _, resource := range resources {
		normalized := normalizeKubernetesWorkloadResource(resource, app.Namespace)
		normalized = kubernetesPortableWorkloadResourceForEmit(app, normalized)
		key := kubernetesDocumentKeyFromMap(normalized)
		if key != "" {
			if _, ok := emitted[key]; ok {
				continue
			}
			emitted[key] = struct{}{}
		}
		doc, err := marshalYAMLString(normalized)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize Kubernetes workload resource %s: %w", key, err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func kubernetesPortableWorkloadResourceForEmit(app *Application, resource map[string]interface{}) map[string]interface{} {
	if app == nil || len(resource) == 0 {
		return resource
	}
	switch toString(resource["kind"]) {
	case "Pod", "Deployment", "ReplicaSet", "ReplicationController", "StatefulSet", "DaemonSet", "Job", "CronJob":
	default:
		return resource
	}
	metadata, _ := asMap(resource["metadata"])
	if metadata == nil {
		return resource
	}
	serviceName := toString(metadata["name"])
	if serviceName == "" {
		return resource
	}
	service := app.Services[serviceName]
	if service == nil {
		return resource
	}
	cloned, _ := deepCopyValue(resource).(map[string]interface{})
	if len(cloned) == 0 {
		return resource
	}
	annotations := map[string]string{}
	if service != nil && len(service.Extensions) > 0 {
		extensions := jsonSerializableExtensionMap(service.Extensions)
		delete(extensions, kubernetesAnnotationPortableServiceExtensions)
		if len(extensions) > 0 {
			if raw, err := json.Marshal(extensions); err == nil {
				annotations[kubernetesAnnotationPortableServiceExtensions] = string(raw)
			}
		}
	}
	if app.Mesh != nil {
		if raw, err := json.Marshal(meshSpecToMap(app.Mesh)); err == nil {
			annotations[kubernetesAnnotationPortableMesh] = string(raw)
		}
	}
	if service.Failover != nil {
		if raw, err := json.Marshal(service.Failover); err == nil {
			annotations[kubernetesAnnotationPortableFailover] = string(raw)
		}
	}
	if len(annotations) > 0 {
		setKubernetesDeploymentAnnotations(cloned, annotations)
	}
	return cloned
}

func normalizeKubernetesWorkloadResource(resource map[string]interface{}, namespace string) map[string]interface{} {
	normalized, _ := deepCopyValue(resource).(map[string]interface{})
	kind := toString(normalized["kind"])
	switch kind {
	case "Pod":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "v1"
		}
	case "Deployment", "ReplicaSet", "StatefulSet", "DaemonSet":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "apps/v1"
		}
	case "ReplicationController":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "v1"
		}
	case "Job", "CronJob":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "batch/v1"
		}
	default:
		return normalized
	}
	metadata, _ := asMap(normalized["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	if namespace != "" && toString(metadata["namespace"]) == "" {
		metadata["namespace"] = namespace
	}
	normalized["metadata"] = metadata
	return normalized
}

func serializeKubernetesServiceAccounts(app *Application) ([]string, error) {
	if app == nil {
		return nil, nil
	}
	var documents []string
	emitted := map[string]struct{}{}

	for _, resource := range kubernetesExtensionResourceSlice(app.Extensions[kubernetesServiceAccountsExtensionKey], app.Extensions[composeKubernetesServiceAccountsExtensionKey]) {
		normalized := normalizeKubernetesServiceAccount(resource, "", nil, app.Namespace)
		key := kubernetesDocumentKeyFromMap(normalized)
		if key != "" {
			if _, ok := emitted[key]; ok {
				continue
			}
			emitted[key] = struct{}{}
		}
		doc, err := marshalYAMLString(normalized)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize serviceaccount resource %s: %w", key, err)
		}
		documents = append(documents, doc)
	}
	if len(app.Services) == 0 {
		return documents, nil
	}
	serviceNames := make([]string, 0, len(app.Services))
	for name, service := range app.Services {
		if service != nil {
			serviceNames = append(serviceNames, name)
		}
	}
	sort.Strings(serviceNames)
	for _, serviceName := range serviceNames {
		service := app.Services[serviceName]
		resource := kubernetesServiceAccountFromService(service)
		if len(resource) == 0 {
			continue
		}
		resource = normalizeKubernetesServiceAccount(resource, serviceName, service, app.Namespace)
		key := kubernetesDocumentKeyFromMap(resource)
		if key != "" {
			if _, ok := emitted[key]; ok {
				continue
			}
			emitted[key] = struct{}{}
		}
		doc, err := marshalYAMLString(resource)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize serviceaccount %s: %w", serviceName, err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func kubernetesServiceAccountFromService(service *Service) map[string]interface{} {
	if service == nil || service.Extensions == nil {
		return nil
	}
	for _, key := range []string{kubernetesServiceAccountExtensionKey, composeKubernetesServiceAccountExtensionKey} {
		if serviceAccount, ok := extensionMapValues(service.Extensions[key]); ok && len(serviceAccount) > 0 {
			return serviceAccount
		}
	}
	return nil
}

func normalizeKubernetesServiceAccount(resource map[string]interface{}, serviceName string, service *Service, namespace string) map[string]interface{} {
	normalized, _ := deepCopyValue(resource).(map[string]interface{})
	if toString(normalized["apiVersion"]) == "" {
		normalized["apiVersion"] = "v1"
	}
	normalized["kind"] = "ServiceAccount"
	metadata, _ := asMap(normalized["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	if toString(metadata["name"]) == "" {
		if service != nil && service.ServiceAccountName != "" {
			metadata["name"] = service.ServiceAccountName
		} else {
			metadata["name"] = serviceName + "-sa"
		}
	}
	if namespace != "" && toString(metadata["namespace"]) == "" {
		metadata["namespace"] = namespace
	}
	normalized["metadata"] = metadata
	return normalized
}

func serializeKubernetesRBACResources(app *Application) ([]string, error) {
	if app == nil {
		return nil, nil
	}
	resources := kubernetesExtensionResourceSlice(app.Extensions[kubernetesRBACResourcesExtensionKey], app.Extensions[composeKubernetesRBACResourcesExtensionKey])
	if len(app.Services) > 0 {
		serviceNames := make([]string, 0, len(app.Services))
		for name, service := range app.Services {
			if service != nil {
				serviceNames = append(serviceNames, name)
			}
		}
		sort.Strings(serviceNames)
		for _, serviceName := range serviceNames {
			bundle := kubernetesRBACFromService(app.Services[serviceName])
			if len(bundle) == 0 {
				continue
			}
			resources = append(resources, kubernetesRBACResourceSlice(bundle["roles"])...)
			resources = append(resources, kubernetesRBACResourceSlice(bundle["bindings"])...)
		}
	}
	var documents []string
	emitted := map[string]struct{}{}
	for _, resource := range resources {
		normalized := normalizeKubernetesRBACResource(resource, app.Namespace)
		key := kubernetesDocumentKeyFromMap(normalized)
		if key != "" {
			if _, ok := emitted[key]; ok {
				continue
			}
			emitted[key] = struct{}{}
		}
		doc, err := marshalYAMLString(normalized)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize rbac resource %s: %w", key, err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func kubernetesRBACFromService(service *Service) map[string]interface{} {
	if service == nil || service.Extensions == nil {
		return nil
	}
	for _, key := range []string{kubernetesRBACExtensionKey, composeKubernetesRBACExtensionKey} {
		if bundle, ok := extensionMapValues(service.Extensions[key]); ok && len(bundle) > 0 {
			return bundle
		}
	}
	return nil
}

func kubernetesRBACResourceSlice(value interface{}) []map[string]interface{} {
	values, _ := value.([]interface{})
	result := make([]map[string]interface{}, 0, len(values))
	for _, item := range values {
		if mapped, ok := asMap(item); ok && len(mapped) > 0 {
			result = append(result, mapped)
		}
	}
	return result
}

func normalizeKubernetesRBACResource(resource map[string]interface{}, namespace string) map[string]interface{} {
	normalized, _ := deepCopyValue(resource).(map[string]interface{})
	if toString(normalized["apiVersion"]) == "" {
		normalized["apiVersion"] = "rbac.authorization.k8s.io/v1"
	}
	metadata, _ := asMap(normalized["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	switch toString(normalized["kind"]) {
	case "Role", "RoleBinding":
		if namespace != "" && toString(metadata["namespace"]) == "" {
			metadata["namespace"] = namespace
		}
	}
	normalized["metadata"] = metadata
	return normalized
}

func serializeKubernetesNamespacePolicyResources(app *Application) ([]string, error) {
	if app == nil {
		return nil, nil
	}
	resources := append(
		kubernetesExtensionResourceSlice(app.Extensions[kubernetesResourceQuotasExtensionKey], app.Extensions[composeKubernetesResourceQuotasExtensionKey]),
		kubernetesExtensionResourceSlice(app.Extensions[kubernetesLimitRangesExtensionKey], app.Extensions[composeKubernetesLimitRangesExtensionKey])...,
	)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesIngressesExtensionKey], app.Extensions[composeKubernetesIngressesExtensionKey])...)
	if rawServices, ok := firstExtensionValue(app.Extensions, kubernetesServicesExtensionKey, "kubernetes.services", composeKubernetesServicesExtensionKey); ok {
		resources = append(resources, kubernetesExtensionResourceSlice(rawServices, nil)...)
	}
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesNetworkPoliciesExtensionKey], app.Extensions[composeKubernetesNetworkPoliciesExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesHPAsExtensionKey], app.Extensions[composeKubernetesHPAsExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesPDBsExtensionKey], app.Extensions[composeKubernetesPDBsExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesPersistentVolumesExtensionKey], app.Extensions[composeKubernetesPersistentVolumesExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesPVCsExtensionKey], app.Extensions[composeKubernetesPVCsExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesPriorityClassesExtensionKey], app.Extensions[composeKubernetesPriorityClassesExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesRuntimeClassesExtensionKey], app.Extensions[composeKubernetesRuntimeClassesExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesStorageClassesExtensionKey], app.Extensions[composeKubernetesStorageClassesExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesIngressClassesExtensionKey], app.Extensions[composeKubernetesIngressClassesExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesMutatingWebhooksExtensionKey], app.Extensions[composeKubernetesMutatingWebhooksExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesValidatingWebhooksExtensionKey], app.Extensions[composeKubernetesValidatingWebhooksExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesCRDsExtensionKey], app.Extensions[composeKubernetesCRDsExtensionKey])...)
	resources = append(resources, kubernetesExtensionResourceSlice(app.Extensions[kubernetesCustomResourcesExtensionKey], app.Extensions[composeKubernetesCustomResourcesExtensionKey])...)
	var documents []string
	emitted := map[string]struct{}{}
	for _, resource := range resources {
		normalized := normalizeKubernetesNamespacePolicyResource(resource, app.Namespace)
		key := kubernetesDocumentKeyFromMap(normalized)
		if key != "" {
			if _, ok := emitted[key]; ok {
				continue
			}
			emitted[key] = struct{}{}
		}
		doc, err := marshalYAMLString(normalized)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize Kubernetes namespace policy resource %s: %w", key, err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func kubernetesExtensionResourceSlice(values ...interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	seen := map[string]struct{}{}
	for _, value := range values {
		items, _ := value.([]interface{})
		for _, item := range items {
			mapped, ok := asMap(item)
			if !ok || len(mapped) == 0 {
				continue
			}
			key := kubernetesDocumentKeyFromMap(mapped)
			if key != "" {
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
			}
			result = append(result, mapped)
		}
	}
	return result
}

func normalizeKubernetesNamespacePolicyResource(resource map[string]interface{}, namespace string) map[string]interface{} {
	normalized, _ := deepCopyValue(resource).(map[string]interface{})
	kind := toString(normalized["kind"])
	switch kind {
	case "Service":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "v1"
		}
	case "Ingress":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "networking.k8s.io/v1"
		}
	case "NetworkPolicy":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "networking.k8s.io/v1"
		}
	case "HorizontalPodAutoscaler":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "autoscaling/v2"
		}
	case "PodDisruptionBudget":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "policy/v1"
		}
	case "PersistentVolume", "PersistentVolumeClaim":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "v1"
		}
	case "ResourceQuota", "LimitRange":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "v1"
		}
	case "PriorityClass":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "scheduling.k8s.io/v1"
		}
	case "RuntimeClass":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "node.k8s.io/v1"
		}
	case "StorageClass":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "storage.k8s.io/v1"
		}
	case "IngressClass":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "networking.k8s.io/v1"
		}
	case "MutatingWebhookConfiguration", "ValidatingWebhookConfiguration":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "admissionregistration.k8s.io/v1"
		}
	case "CustomResourceDefinition":
		if toString(normalized["apiVersion"]) == "" {
			normalized["apiVersion"] = "apiextensions.k8s.io/v1"
		}
	default:
		return normalized
	}
	metadata, _ := asMap(normalized["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	if namespace != "" && toString(metadata["namespace"]) == "" && (kind == "ResourceQuota" || kind == "LimitRange" || kind == "Service" || kind == "Ingress" || kind == "NetworkPolicy" || kind == "HorizontalPodAutoscaler" || kind == "PodDisruptionBudget" || kind == "PersistentVolumeClaim") {
		metadata["namespace"] = namespace
	}
	normalized["metadata"] = metadata
	return normalized
}

func hasKubernetesTypedResource(app *Application, kind, name string) bool {
	namespace := ""
	if app != nil && kubernetesKindIsNamespaced(kind) {
		namespace = app.Namespace
	}
	return hasKubernetesTypedResourceInNamespace(app, kind, name, namespace)
}

func hasKubernetesTypedResourceInNamespace(app *Application, kind, name, namespace string) bool {
	if app == nil || name == "" {
		return false
	}
	keys := []string{}
	switch kind {
	case "Namespace":
		keys = []string{kubernetesNamespacesExtensionKey, composeKubernetesNamespacesExtensionKey}
	case "Pod", "Deployment", "ReplicaSet", "ReplicationController", "StatefulSet", "DaemonSet", "Job", "CronJob":
		keys = []string{kubernetesWorkloadsExtensionKey, composeKubernetesWorkloadsExtensionKey}
	case "ConfigMap":
		keys = []string{kubernetesConfigMapsExtensionKey, composeKubernetesConfigMapsExtensionKey}
	case "Secret":
		keys = []string{kubernetesSecretsExtensionKey, composeKubernetesSecretsExtensionKey}
	case "ServiceAccount":
		keys = []string{kubernetesServiceAccountsExtensionKey, composeKubernetesServiceAccountsExtensionKey}
	case "Service":
		keys = []string{kubernetesServicesExtensionKey, "kubernetes.services", composeKubernetesServicesExtensionKey}
	case "Ingress":
		keys = []string{kubernetesIngressesExtensionKey, composeKubernetesIngressesExtensionKey, "kubernetes.routes"}
	case "NetworkPolicy":
		keys = []string{kubernetesNetworkPoliciesExtensionKey, composeKubernetesNetworkPoliciesExtensionKey, "kubernetes.policies"}
	case "Role", "RoleBinding", "ClusterRole", "ClusterRoleBinding":
		keys = []string{kubernetesRBACResourcesExtensionKey, composeKubernetesRBACResourcesExtensionKey}
	case "HorizontalPodAutoscaler":
		keys = []string{kubernetesHPAsExtensionKey, composeKubernetesHPAsExtensionKey}
	case "PodDisruptionBudget":
		keys = []string{kubernetesPDBsExtensionKey, composeKubernetesPDBsExtensionKey}
	case "PersistentVolume":
		keys = []string{kubernetesPersistentVolumesExtensionKey, composeKubernetesPersistentVolumesExtensionKey}
	case "PersistentVolumeClaim":
		keys = []string{kubernetesPVCsExtensionKey, composeKubernetesPVCsExtensionKey}
	case "ResourceQuota":
		keys = []string{kubernetesResourceQuotasExtensionKey, composeKubernetesResourceQuotasExtensionKey}
	case "LimitRange":
		keys = []string{kubernetesLimitRangesExtensionKey, composeKubernetesLimitRangesExtensionKey}
	case "PriorityClass":
		keys = []string{kubernetesPriorityClassesExtensionKey, composeKubernetesPriorityClassesExtensionKey}
	case "RuntimeClass":
		keys = []string{kubernetesRuntimeClassesExtensionKey, composeKubernetesRuntimeClassesExtensionKey}
	case "StorageClass":
		keys = []string{kubernetesStorageClassesExtensionKey, composeKubernetesStorageClassesExtensionKey}
	case "IngressClass":
		keys = []string{kubernetesIngressClassesExtensionKey, composeKubernetesIngressClassesExtensionKey}
	case "MutatingWebhookConfiguration":
		keys = []string{kubernetesMutatingWebhooksExtensionKey, composeKubernetesMutatingWebhooksExtensionKey}
	case "ValidatingWebhookConfiguration":
		keys = []string{kubernetesValidatingWebhooksExtensionKey, composeKubernetesValidatingWebhooksExtensionKey}
	case "CustomResourceDefinition":
		keys = []string{kubernetesCRDsExtensionKey, composeKubernetesCRDsExtensionKey}
	default:
		return false
	}
	keys = appendKubernetesGenericResourceKeys(keys)
	keys = appendKubernetesLegacyResourceKeys(kind, keys)
	for _, key := range keys {
		for _, resource := range kubernetesExtensionResourceSlice(app.Extensions[key]) {
			if toString(resource["kind"]) != kind {
				continue
			}
			metadata, _ := asMap(resource["metadata"])
			if toString(metadata["name"]) == name && kubernetesResourceNamespaceMatches(kind, metadata, namespace) {
				return true
			}
		}
	}
	return false
}

func appendKubernetesLegacyResourceKeys(kind string, keys []string) []string {
	switch kind {
	case "Service":
		return append(keys, "kubernetes.services")
	default:
		return keys
	}
}

func appendKubernetesGenericResourceKeys(keys []string) []string {
	return append(keys, kubernetesResourcesExtensionKey, composeKubernetesResourcesExtensionKey)
}

func kubernetesResourceNamespaceMatches(kind string, metadata map[string]interface{}, namespace string) bool {
	if !kubernetesKindIsNamespaced(kind) {
		return true
	}
	resourceNamespace := toString(metadata["namespace"])
	if namespace != "" && resourceNamespace == "" {
		return true
	}
	return resourceNamespace == namespace
}

func kubernetesKindIsNamespaced(kind string) bool {
	switch kind {
	case "Namespace",
		"PersistentVolume",
		"PriorityClass",
		"RuntimeClass",
		"StorageClass",
		"IngressClass",
		"ClusterRole",
		"ClusterRoleBinding",
		"MutatingWebhookConfiguration",
		"ValidatingWebhookConfiguration",
		"CustomResourceDefinition":
		return false
	default:
		return true
	}
}

func hasKubernetesTypedWorkloadResource(app *Application, name string) bool {
	if app == nil || name == "" {
		return false
	}
	namespace := app.Namespace
	for _, kind := range []string{"Pod", "Deployment", "ReplicaSet", "ReplicationController", "StatefulSet", "DaemonSet", "Job", "CronJob"} {
		if hasKubernetesTypedResourceInNamespace(app, kind, name, namespace) {
			return true
		}
	}
	return false
}

func serializeKubernetesHorizontalPodAutoscaler(serviceName string, service *Service, namespace string) (string, error) {
	hpa := kubernetesHorizontalPodAutoscalerFromService(service)
	if len(hpa) == 0 {
		return "", nil
	}
	resource, _ := deepCopyValue(hpa).(map[string]interface{})
	if toString(resource["apiVersion"]) == "" {
		resource["apiVersion"] = "autoscaling/v2"
	}
	resource["kind"] = "HorizontalPodAutoscaler"
	metadata, _ := asMap(resource["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	if toString(metadata["name"]) == "" {
		metadata["name"] = serviceName + "-hpa"
	}
	if namespace != "" && toString(metadata["namespace"]) == "" {
		metadata["namespace"] = namespace
	}
	resource["metadata"] = metadata
	spec, _ := asMap(resource["spec"])
	if spec == nil {
		spec = map[string]interface{}{}
	}
	scaleTargetRef, _ := asMap(spec["scaleTargetRef"])
	if scaleTargetRef == nil {
		scaleTargetRef = map[string]interface{}{}
	}
	if toString(scaleTargetRef["apiVersion"]) == "" {
		scaleTargetRef["apiVersion"] = "apps/v1"
	}
	if toString(scaleTargetRef["kind"]) == "" {
		scaleTargetRef["kind"] = "Deployment"
	}
	scaleTargetRef["name"] = serviceName
	spec["scaleTargetRef"] = scaleTargetRef
	resource["spec"] = spec
	return marshalYAMLString(resource)
}

func kubernetesHorizontalPodAutoscalerFromService(service *Service) map[string]interface{} {
	if service == nil || service.Extensions == nil {
		return nil
	}
	for _, key := range []string{kubernetesHPAExtensionKey, composeKubernetesHPAExtensionKey} {
		if hpa, ok := extensionMapValues(service.Extensions[key]); ok && len(hpa) > 0 {
			return hpa
		}
	}
	return nil
}

func serializeKubernetesPodDisruptionBudget(serviceName string, service *Service, namespace string) (string, error) {
	pdb := kubernetesPodDisruptionBudgetFromService(service)
	if len(pdb) == 0 {
		return "", nil
	}
	resource, _ := deepCopyValue(pdb).(map[string]interface{})
	if toString(resource["apiVersion"]) == "" {
		resource["apiVersion"] = "policy/v1"
	}
	resource["kind"] = "PodDisruptionBudget"
	metadata, _ := asMap(resource["metadata"])
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	if toString(metadata["name"]) == "" {
		metadata["name"] = serviceName + "-pdb"
	}
	if namespace != "" && toString(metadata["namespace"]) == "" {
		metadata["namespace"] = namespace
	}
	resource["metadata"] = metadata
	spec, _ := asMap(resource["spec"])
	if spec == nil {
		spec = map[string]interface{}{}
	}
	if selector, _ := asMap(spec["selector"]); len(selector) == 0 {
		if labels := serviceExtensionStringMap(service, "kubernetes.template.labels", "x-kubernetes-template-labels"); len(labels) > 0 {
			spec["selector"] = map[string]interface{}{"matchLabels": labels}
		} else if len(service.Labels) > 0 {
			spec["selector"] = map[string]interface{}{"matchLabels": copyStringMap(service.Labels)}
		} else {
			spec["selector"] = map[string]interface{}{"matchLabels": map[string]string{"app": serviceName}}
		}
	}
	resource["spec"] = spec
	return marshalYAMLString(resource)
}

func kubernetesPodDisruptionBudgetFromService(service *Service) map[string]interface{} {
	if service == nil || service.Extensions == nil {
		return nil
	}
	for _, key := range []string{kubernetesPDBExtensionKey, composeKubernetesPDBExtensionKey} {
		if pdb, ok := extensionMapValues(service.Extensions[key]); ok && len(pdb) > 0 {
			return pdb
		}
	}
	return nil
}

func kubernetesResourceNameOrFallback(resource map[string]interface{}, fallback string) string {
	metadata, _ := asMap(resource["metadata"])
	if name := toString(metadata["name"]); name != "" {
		return name
	}
	return fallback
}

// RestoreKubernetesSource writes the preserved Kubernetes manifests back to
// disk. It prefers the original source extensions when present and falls back
// to the canonical raw resource graph after cross-format conversion hops.
func RestoreKubernetesSource(app *Application, filename string) error {
	content, err := kubernetesSourceContent(app)
	if err != nil {
		return err
	}
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("filename is required")
	}
	if dir := filepath.Dir(filename); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create Kubernetes output directory: %w", err)
		}
	}
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to restore Kubernetes source: %w", err)
	}
	return nil
}

// RestoreKubernetesJSONSource writes the preserved Kubernetes manifests as JSON.
func RestoreKubernetesJSONSource(app *Application, filename string) error {
	content, err := kubernetesJSONSourceContent(app)
	if err != nil {
		return err
	}
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("filename is required")
	}
	if dir := filepath.Dir(filename); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create Kubernetes JSON output directory: %w", err)
		}
	}
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to restore Kubernetes JSON: %w", err)
	}
	return nil
}

func kubernetesSourceContent(app *Application) (string, error) {
	rawDocs := kubernetesRawManifestDocs(app)
	if len(rawDocs) == 0 {
		return "", fmt.Errorf("application does not contain raw Kubernetes manifests")
	}
	return strings.Join(rawDocs, "\n---\n"), nil
}

func kubernetesJSONSourceContent(app *Application) (string, error) {
	jsonContent, err := SerializeKubernetesJSON(app)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(jsonContent) == "" {
		return "", fmt.Errorf("application does not contain raw Kubernetes manifests")
	}
	return jsonContent, nil
}

func kubernetesJSONToYAMLContent(content string) string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}
	var decoded interface{}
	if err := json.Unmarshal([]byte(trimmed), &decoded); err != nil {
		return ""
	}
	switch typed := decoded.(type) {
	case []interface{}:
		documents := make([]string, 0, len(typed))
		for _, item := range typed {
			doc, err := yaml.Marshal(item)
			if err != nil {
				return ""
			}
			if trimmedDoc := strings.TrimSpace(string(doc)); trimmedDoc != "" {
				documents = append(documents, trimmedDoc)
			}
		}
		return strings.Join(documents, "\n---\n")
	case map[string]interface{}:
		doc, err := yaml.Marshal(typed)
		if err != nil {
			return ""
		}
		return string(doc)
	default:
		return ""
	}
}

func kubernetesRawManifestDocs(app *Application) []string {
	if docs := kubernetesSourceDocsFromCanonical(app); len(docs) > 0 {
		return docs
	}
	if docs := kubernetesSourceDocsFromExtension(app); len(docs) > 0 {
		return docs
	}
	if app == nil || app.Platform != PlatformKubernetes {
		if docs := kubernetesRawManifestDocsFromCanonicalResources(app); len(docs) > 0 {
			return docs
		}
	}
	docs := kubernetesRawDocsFromExtension(app)
	if len(docs) > 0 {
		return docs
	}
	return kubernetesRawManifestDocsFromCanonicalResources(app)
}

func kubernetesRawManifestDocsFromCanonicalResources(app *Application) []string {
	canonical := canonicalForApplication(app)
	if canonical == nil {
		return nil
	}
	type rawKubernetesResource struct {
		ordinal int
		key     string
		doc     string
		kind    ResourceKind
	}
	rawDocsByKey := map[string]rawKubernetesResource{}
	for _, resource := range canonical.Resources {
		if resource == nil || (resource.Platform != PlatformKubernetes && !kubernetesRawResourceLike(resource.Raw)) {
			continue
		}
		if resource.Kind != ResourceKindRaw && resource.Kind != ResourceKindUnknown {
			continue
		}
		raw, ok := resource.Raw.(map[string]interface{})
		if !ok {
			continue
		}
		doc, err := marshalYAMLString(raw)
		if err != nil {
			continue
		}
		candidate := rawKubernetesResource{
			ordinal: resource.Ordinal,
			key:     kubernetesDocumentKeyFromMap(raw),
			doc:     doc,
			kind:    resource.Kind,
		}
		if existing, ok := rawDocsByKey[candidate.key]; ok {
			if existing.kind == ResourceKindUnknown && candidate.kind == ResourceKindRaw {
				rawDocsByKey[candidate.key] = candidate
			}
			continue
		}
		rawDocsByKey[candidate.key] = candidate
	}
	rawDocs := make([]rawKubernetesResource, 0, len(rawDocsByKey))
	for _, item := range rawDocsByKey {
		rawDocs = append(rawDocs, item)
	}
	sort.Slice(rawDocs, func(i, j int) bool {
		if rawDocs[i].ordinal != rawDocs[j].ordinal {
			return rawDocs[i].ordinal < rawDocs[j].ordinal
		}
		return rawDocs[i].key < rawDocs[j].key
	})
	restoredDocs := make([]string, 0, len(rawDocs))
	for _, item := range rawDocs {
		restoredDocs = append(restoredDocs, item.doc)
	}
	return restoredDocs
}

func splitKubernetesSourceDocuments(content string) []string {
	var docs []string
	var current strings.Builder
	lines := strings.Split(content, "\n")
	flush := func() {
		doc := strings.TrimSpace(current.String())
		current.Reset()
		if doc != "" {
			docs = append(docs, doc)
		}
	}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" || strings.HasPrefix(trimmed, "--- #") {
			flush()
			continue
		}
		if current.Len() > 0 {
			current.WriteByte('\n')
		}
		current.WriteString(line)
	}
	flush()
	return docs
}

func kubernetesSourceDocsFromCanonical(app *Application) []string {
	canonical := canonicalForApplication(app)
	if canonical == nil || len(canonical.Extensions) == 0 {
		return nil
	}
	return kubernetesSourceDocsFromValue(canonical.Extensions[kubernetesSourceDocumentsExtensionKey])
}

func kubernetesSourceDocsFromExtension(app *Application) []string {
	if app == nil || len(app.Extensions) == 0 {
		return nil
	}
	return kubernetesSourceDocsFromValue(app.Extensions[kubernetesSourceDocumentsExtensionKey])
}

func kubernetesSourceDocsFromValue(value interface{}) []string {
	switch typed := value.(type) {
	case []string:
		docs := make([]string, 0, len(typed))
		for _, doc := range typed {
			if trimmed := strings.TrimSpace(doc); trimmed != "" {
				docs = append(docs, trimmed)
			}
		}
		return docs
	case []interface{}:
		docs := make([]string, 0, len(typed))
		for _, item := range typed {
			if trimmed := strings.TrimSpace(toString(item)); trimmed != "" {
				docs = append(docs, trimmed)
			}
		}
		return docs
	default:
		return nil
	}
}

func kubernetesRawDocsFromExtension(app *Application) []string {
	if app == nil || app.Extensions == nil {
		return nil
	}
	value, ok := applicationExtensionValueForKey(app, "kubernetes.raw")
	if !ok {
		return nil
	}
	switch typed := value.(type) {
	case []interface{}:
		docs := make([]string, 0, len(typed))
		for _, item := range typed {
			if mapped, ok := asMap(item); ok {
				if doc, err := marshalYAMLString(mapped); err == nil {
					docs = append(docs, doc)
				}
			}
		}
		return docs
	case []map[string]interface{}:
		docs := make([]string, 0, len(typed))
		for _, item := range typed {
			if doc, err := marshalYAMLString(item); err == nil {
				docs = append(docs, doc)
			}
		}
		return docs
	default:
		return nil
	}
}

func appendRawKubernetesResources(documents []string, app *Application) ([]string, error) {
	emitted := map[string]struct{}{}
	documentIndexes := map[string]int{}
	for index, doc := range documents {
		if key := kubernetesDocumentKeyFromYAML(doc); key != "" {
			emitted[key] = struct{}{}
			documentIndexes[key] = index
		}
	}
	for _, extensionKey := range []string{
		"kubernetes.raw",
		"kubernetes.services",
		"kubernetes.hpas",
		"kubernetes.pdbs",
		"kubernetes.routes",
		"kubernetes.policies",
		"kubernetes.resources",
		composeKubernetesResourcesExtensionKey,
	} {
		var extensionValue interface{}
		if strings.HasPrefix(extensionKey, "kubernetes.") {
			if value, ok := applicationExtensionValueForKey(app, extensionKey); ok {
				extensionValue = value
			}
		} else if app != nil && app.Extensions != nil {
			extensionValue = app.Extensions[extensionKey]
		}
		for _, raw := range kubernetesRawExtensionResources(extensionValue) {
			resource := kubernetesResourceForRawAppend(app, extensionKey, raw)
			key := kubernetesDocumentKeyFromMap(kubernetesResourceKeyForRawAppend(app, extensionKey, raw, resource))
			if key == "" {
				continue
			}
			doc, err := marshalYAMLString(resource)
			if err != nil {
				return documents, fmt.Errorf("failed to serialize raw Kubernetes resource %s: %w", key, err)
			}
			if existingIndex, exists := documentIndexes[key]; exists {
				if rawKubernetesExtensionShouldOverwrite(extensionKey) && !kubernetesDocumentHasPortableAnnotations(documents[existingIndex]) {
					documents[existingIndex] = doc
				}
				continue
			}
			documents = append(documents, doc)
			emitted[key] = struct{}{}
			documentIndexes[key] = len(documents) - 1
		}
	}
	for _, raw := range kubernetesCanonicalRawResources(app) {
		key := kubernetesDocumentKeyFromMap(normalizeKubernetesGenericResourceForEmit(raw, applicationNamespace(app)))
		if key == "" {
			continue
		}
		if _, exists := emitted[key]; exists {
			continue
		}
		doc, err := marshalYAMLString(raw)
		if err != nil {
			return documents, fmt.Errorf("failed to serialize canonical raw Kubernetes resource %s: %w", key, err)
		}
		documents = append(documents, doc)
		emitted[key] = struct{}{}
		documentIndexes[key] = len(documents) - 1
	}
	return documents, nil
}

func rawKubernetesExtensionShouldOverwrite(extensionKey string) bool {
	switch extensionKey {
	case "kubernetes.raw",
		"kubernetes.services",
		"kubernetes.hpas",
		"kubernetes.pdbs",
		"kubernetes.routes",
		"kubernetes.policies",
		"kubernetes.resources",
		composeKubernetesResourcesExtensionKey:
		return true
	default:
		return false
	}
}

func kubernetesDocumentHasPortableAnnotations(doc string) bool {
	return strings.Contains(doc, kubernetesAnnotationPortableMesh) ||
		strings.Contains(doc, kubernetesAnnotationPortableFailover) ||
		strings.Contains(doc, kubernetesAnnotationPortableServiceExtensions) ||
		strings.Contains(doc, kubernetesAnnotationPortableAppExtensions)
}

func kubernetesResourceKeyForRawAppend(app *Application, extensionKey string, raw map[string]interface{}, resource map[string]interface{}) map[string]interface{} {
	if extensionKey != "kubernetes.raw" {
		return resource
	}
	return normalizeKubernetesGenericResourceForEmit(raw, applicationNamespace(app))
}

func kubernetesResourceForRawAppend(app *Application, extensionKey string, resource map[string]interface{}) map[string]interface{} {
	namespace := ""
	if app != nil {
		namespace = app.Namespace
	}
	if extensionKey == "kubernetes.raw" {
		return kubernetesPortableWorkloadResourceForEmit(app, normalizeKubernetesGenericResourceForEmit(resource, namespace))
	}
	return kubernetesPortableWorkloadResourceForEmit(app, normalizeKubernetesGenericResourceForEmit(resource, namespace))
}

func normalizeKubernetesGenericResourceForEmit(resource map[string]interface{}, namespace string) map[string]interface{} {
	kind := toString(resource["kind"])
	switch kind {
	case "ConfigMap":
		return normalizeKubernetesConfigMapResource(resource, namespace)
	case "Secret":
		return normalizeKubernetesSecretResource(resource, namespace)
	case "Pod", "Deployment", "ReplicaSet", "ReplicationController", "StatefulSet", "DaemonSet", "Job", "CronJob":
		return normalizeKubernetesWorkloadResource(resource, namespace)
	case "ServiceAccount":
		return normalizeKubernetesServiceAccount(resource, "", nil, namespace)
	case "Role", "RoleBinding", "ClusterRole", "ClusterRoleBinding":
		return normalizeKubernetesRBACResource(resource, namespace)
	case "Service",
		"Ingress",
		"NetworkPolicy",
		"HorizontalPodAutoscaler",
		"PodDisruptionBudget",
		"PersistentVolume",
		"PersistentVolumeClaim",
		"ResourceQuota",
		"LimitRange",
		"PriorityClass",
		"RuntimeClass",
		"StorageClass",
		"IngressClass",
		"MutatingWebhookConfiguration",
		"ValidatingWebhookConfiguration",
		"CustomResourceDefinition":
		return normalizeKubernetesNamespacePolicyResource(resource, namespace)
	default:
		normalized, _ := deepCopyValue(resource).(map[string]interface{})
		return normalized
	}
}

func kubernetesCanonicalRawResources(app *Application) []map[string]interface{} {
	canonical := canonicalForApplication(app)
	if canonical == nil {
		return nil
	}
	var resources []map[string]interface{}
	for _, resource := range canonical.Resources {
		if resource == nil || (resource.Platform != PlatformKubernetes && !kubernetesRawResourceLike(resource.Raw)) {
			continue
		}
		if resource.Kind != ResourceKindRaw && resource.Kind != ResourceKindUnknown {
			continue
		}
		raw, ok := resource.Raw.(map[string]interface{})
		if !ok {
			continue
		}
		resources = append(resources, raw)
	}
	sort.Slice(resources, func(i, j int) bool {
		return kubernetesDocumentKeyFromMap(resources[i]) < kubernetesDocumentKeyFromMap(resources[j])
	})
	return resources
}

func kubernetesRawExtensionResources(value interface{}) []map[string]interface{} {
	var resources []map[string]interface{}
	switch typed := value.(type) {
	case []interface{}:
		for _, item := range typed {
			if mapped, ok := asMap(item); ok {
				resources = append(resources, mapped)
			}
		}
	case []map[string]interface{}:
		resources = append(resources, typed...)
	}
	return resources
}

func kubernetesDocumentKeyFromYAML(doc string) string {
	var resource map[string]interface{}
	if err := yaml.Unmarshal([]byte(doc), &resource); err != nil {
		return ""
	}
	return kubernetesDocumentKeyFromMap(resource)
}

func kubernetesDocumentKeyFromMap(resource map[string]interface{}) string {
	if len(resource) == 0 {
		return ""
	}
	apiVersion := toString(resource["apiVersion"])
	kind := toString(resource["kind"])
	if apiVersion == "" {
		apiVersion = kubernetesDefaultAPIVersionForKind(kind)
	}
	metadata, _ := asMap(resource["metadata"])
	name := toString(metadata["name"])
	namespace := toString(metadata["namespace"])
	if kind == "" || name == "" {
		return ""
	}
	return strings.Join([]string{apiVersion, kind, namespace, name}, "/")
}

func kubernetesDefaultAPIVersionForKind(kind string) string {
	switch kind {
	case "Pod",
		"Service",
		"ConfigMap",
		"Secret",
		"ServiceAccount",
		"Namespace",
		"PersistentVolume",
		"PersistentVolumeClaim",
		"ResourceQuota",
		"LimitRange":
		return "v1"
	case "Deployment", "ReplicaSet", "StatefulSet", "DaemonSet":
		return "apps/v1"
	case "ReplicationController":
		return "v1"
	case "Job", "CronJob":
		return "batch/v1"
	case "Ingress", "NetworkPolicy", "IngressClass":
		return "networking.k8s.io/v1"
	case "HorizontalPodAutoscaler":
		return "autoscaling/v2"
	case "PodDisruptionBudget":
		return "policy/v1"
	case "Role", "RoleBinding", "ClusterRole", "ClusterRoleBinding":
		return "rbac.authorization.k8s.io/v1"
	case "PriorityClass":
		return "scheduling.k8s.io/v1"
	case "RuntimeClass":
		return "node.k8s.io/v1"
	case "StorageClass":
		return "storage.k8s.io/v1"
	case "MutatingWebhookConfiguration", "ValidatingWebhookConfiguration":
		return "admissionregistration.k8s.io/v1"
	case "CustomResourceDefinition":
		return "apiextensions.k8s.io/v1"
	default:
		return ""
	}
}

func serializeKubernetesDeployment(app *Application, name string, service *Service, namespace string) (string, error) {
	replicas := service.Replicas
	if service.Deploy != nil && service.Deploy.Replicas > 0 {
		replicas = service.Deploy.Replicas
	}
	if replicas == 0 {
		replicas = 1
	}

	deployment := kubernetesWorkloadSkeleton(name, service, replicas, namespace)
	applyKubernetesWorkloadMetadata(deployment, service)

	applyKubernetesWorkloadSpecExtensions(deployment, service)
	applyKubernetesDeployIntent(deployment, service)
	applyKubernetesDependencies(deployment, service, app)
	applyKubernetesRuntimeAnnotations(deployment, service)
	applyKubernetesPortableAnnotations(deployment, service, app)

	podTemplate := kubernetesSerializedPodTemplate(deployment)
	podSpec := podTemplate["spec"].(map[string]interface{})
	container := podSpec["containers"].([]map[string]interface{})[0]

	// Add command
	if len(service.Entrypoint) > 0 {
		container["command"] = service.Entrypoint
	}
	if len(service.Command) > 0 {
		container["args"] = service.Command
	}
	if service.WorkingDir != "" {
		container["workingDir"] = service.WorkingDir
	}
	if service.StdinOpenSet {
		container["stdin"] = service.StdinOpen
	} else if service.StdinOpen {
		container["stdin"] = true
	}
	if service.TtySet {
		container["tty"] = service.Tty
	} else if service.Tty {
		container["tty"] = true
	}
	if securityContext := kubernetesSecurityContext(service); len(securityContext) > 0 {
		container["securityContext"] = securityContext
	}
	if podSecurityContext := kubernetesPodSecurityContext(service); len(podSecurityContext) > 0 {
		podSpec["securityContext"] = podSecurityContext
	}
	if service.ActiveDeadlineSeconds != nil && *service.ActiveDeadlineSeconds > 0 {
		podSpec["activeDeadlineSeconds"] = *service.ActiveDeadlineSeconds
	}
	if service.PodRestartPolicy != "" {
		podSpec["restartPolicy"] = service.PodRestartPolicy
	}
	if seconds := durationSeconds(service.StopGracePeriod); seconds > 0 {
		podSpec["terminationGracePeriodSeconds"] = seconds
	}
	if service.ServiceAccountName != "" {
		podSpec["serviceAccountName"] = service.ServiceAccountName
	} else if value := extensionStringValue(service, "kubernetes.serviceAccountName"); value != "" {
		podSpec["serviceAccountName"] = value
	}
	if service.PriorityClassName != "" {
		podSpec["priorityClassName"] = service.PriorityClassName
	} else if value := extensionStringValue(service, "kubernetes.priorityClassName"); value != "" {
		podSpec["priorityClassName"] = value
	}
	if service.RuntimeClassName != "" {
		podSpec["runtimeClassName"] = service.RuntimeClassName
	} else if value := extensionStringValue(service, "kubernetes.runtimeClassName"); value != "" {
		podSpec["runtimeClassName"] = value
	}
	if service.NodeName != "" {
		podSpec["nodeName"] = service.NodeName
	} else if value := extensionStringValue(service, "kubernetes.nodeName"); value != "" {
		podSpec["nodeName"] = value
	}
	if service.Subdomain != "" {
		podSpec["subdomain"] = service.Subdomain
	} else if value := extensionStringValue(service, "kubernetes.subdomain"); value != "" {
		podSpec["subdomain"] = value
	}
	if service.SetHostnameAsFQDN != nil {
		podSpec["setHostnameAsFQDN"] = *service.SetHostnameAsFQDN
	} else if value := extensionStringValue(service, "kubernetes.setHostnameAsFQDN"); value != "" {
		podSpec["setHostnameAsFQDN"] = strings.EqualFold(value, "true")
	}
	if service.HostUsers != nil {
		podSpec["hostUsers"] = *service.HostUsers
	} else if value := extensionStringValue(service, "kubernetes.hostUsers"); value != "" {
		podSpec["hostUsers"] = strings.EqualFold(value, "true")
	}
	if service.ShareProcessNamespace != nil {
		podSpec["shareProcessNamespace"] = *service.ShareProcessNamespace
	} else if value := extensionStringValue(service, "kubernetes.shareProcessNamespace"); value != "" {
		podSpec["shareProcessNamespace"] = strings.EqualFold(value, "true")
	}
	if service.EnableServiceLinks != nil {
		podSpec["enableServiceLinks"] = *service.EnableServiceLinks
	} else if value := extensionStringValue(service, "kubernetes.enableServiceLinks"); value != "" {
		podSpec["enableServiceLinks"] = strings.EqualFold(value, "true")
	}
	if service.HostPID != nil {
		podSpec["hostPID"] = *service.HostPID
	} else if value := extensionStringValue(service, "kubernetes.hostPID"); value != "" {
		podSpec["hostPID"] = strings.EqualFold(value, "true")
	}
	if service.PIDMode != "" {
		if strings.EqualFold(service.PIDMode, "host") {
			podSpec["hostPID"] = true
		}
	} else if value := extensionStringValue(service, "kubernetes.pidMode"); value != "" {
		if strings.EqualFold(value, "host") {
			podSpec["hostPID"] = true
		}
	}
	if service.HostIPC != nil {
		podSpec["hostIPC"] = *service.HostIPC
	} else if value := extensionStringValue(service, "kubernetes.hostIPC"); value != "" {
		podSpec["hostIPC"] = strings.EqualFold(value, "true")
	}
	if service.IPCMode != "" {
		if strings.EqualFold(service.IPCMode, "host") {
			podSpec["hostIPC"] = true
		}
	} else if value := extensionStringValue(service, "kubernetes.ipcMode"); value != "" {
		if strings.EqualFold(value, "host") {
			podSpec["hostIPC"] = true
		}
	}
	if service.Hostname != "" {
		podSpec["hostname"] = service.Hostname
	} else if value := extensionStringValue(service, "kubernetes.hostname"); value != "" {
		podSpec["hostname"] = value
	}
	if service.AutomountServiceAccountToken != nil {
		podSpec["automountServiceAccountToken"] = *service.AutomountServiceAccountToken
	} else if value := extensionStringValue(service, "kubernetes.automountServiceAccountToken"); value != "" {
		podSpec["automountServiceAccountToken"] = strings.EqualFold(value, "true")
	}
	if len(service.ImagePullSecrets) > 0 {
		var imagePullSecrets []map[string]interface{}
		for _, name := range service.ImagePullSecrets {
			if strings.TrimSpace(name) == "" {
				continue
			}
			imagePullSecrets = append(imagePullSecrets, map[string]interface{}{"name": name})
		}
		if len(imagePullSecrets) > 0 {
			podSpec["imagePullSecrets"] = imagePullSecrets
		}
	}
	if service.ImagePullPolicy != "" {
		container["imagePullPolicy"] = service.ImagePullPolicy
	} else if value := extensionStringValue(service, "kubernetes.imagePullPolicy"); value != "" {
		container["imagePullPolicy"] = value
	}
	if service.TerminationMessagePath != "" {
		container["terminationMessagePath"] = service.TerminationMessagePath
	} else if value := extensionStringValue(service, "kubernetes.terminationMessagePath"); value != "" {
		container["terminationMessagePath"] = value
	}
	if service.TerminationMessagePolicy != "" {
		container["terminationMessagePolicy"] = service.TerminationMessagePolicy
	} else if value := extensionStringValue(service, "kubernetes.terminationMessagePolicy"); value != "" {
		container["terminationMessagePolicy"] = value
	}
	if len(service.Tolerations) > 0 {
		tolerations := serializeKubernetesTolerationsNative(service.Tolerations)
		if len(tolerations) > 0 {
			podSpec["tolerations"] = tolerations
		}
	}
	hostNetwork := service.HostNetwork
	hostNetworkSet := service.HostNetworkSet
	if service.WindowsOptions != nil && service.WindowsOptions.HostProcess != nil && *service.WindowsOptions.HostProcess {
		hostNetwork = true
		hostNetworkSet = true
	}
	if hostNetworkSet {
		podSpec["hostNetwork"] = hostNetwork
	} else if hostNetwork {
		podSpec["hostNetwork"] = true
	} else if value, ok := extensionBoolValue(service, "kubernetes.hostNetwork"); ok {
		podSpec["hostNetwork"] = value
	}
	if len(service.Affinity) > 0 {
		podSpec["affinity"] = cloneMap(service.Affinity)
	} else if value, ok := extensionValueForKey(service, "kubernetes.affinity"); ok {
		if affinity, ok := extensionMapValues(value); ok && len(affinity) > 0 {
			podSpec["affinity"] = affinity
		}
	}
	if len(service.TopologySpreadConstraints) > 0 {
		podSpec["topologySpreadConstraints"] = cloneMapSlice(service.TopologySpreadConstraints)
	} else if value, ok := extensionValueForKey(service, "kubernetes.topologySpreadConstraints"); ok {
		if constraints, ok := kubernetesTopologySpreadConstraintsFromExtension(value); len(constraints) > 0 && ok {
			podSpec["topologySpreadConstraints"] = cloneMapSlice(constraints)
		}
	}
	if len(service.ReadinessGates) > 0 {
		podSpec["readinessGates"] = cloneMapSlice(service.ReadinessGates)
	} else if value, ok := extensionValueForKey(service, "kubernetes.readinessGates"); ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); len(gates) > 0 && ok {
			podSpec["readinessGates"] = cloneMapSlice(gates)
		}
	}
	if value, ok := extensionValueForKey(service, "kubernetes.resourceClaims"); ok {
		if claims, ok := kubernetesMapSliceFromExtension(value); ok && len(claims) > 0 {
			podSpec["resourceClaims"] = cloneMapSlice(claims)
		}
	}
	if value, ok := extensionValueForKey(service, "x-kubernetes-resource-claims"); ok {
		if claims, ok := kubernetesMapSliceFromExtension(value); ok && len(claims) > 0 {
			podSpec["resourceClaims"] = cloneMapSlice(claims)
		}
	}
	if len(service.SchedulingGates) > 0 {
		podSpec["schedulingGates"] = cloneMapSlice(service.SchedulingGates)
	} else if value, ok := extensionValueForKey(service, "kubernetes.schedulingGates"); ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); len(gates) > 0 && ok {
			podSpec["schedulingGates"] = cloneMapSlice(gates)
		}
	}
	if value, ok := extensionValueForKey(service, "kubernetes.ephemeralContainers"); ok {
		if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
			podSpec["ephemeralContainers"] = cloneMapSlice(containers)
		}
	}
	if value, ok := extensionValueForKey(service, "x-kubernetes-ephemeral-containers"); ok {
		if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
			podSpec["ephemeralContainers"] = cloneMapSlice(containers)
		}
	}
	if service.SchedulerName != "" {
		podSpec["schedulerName"] = service.SchedulerName
	}
	applyKubernetesDNS(deployment, service)

	// Add environment variables
	if env := kubernetesEnvVars(service); len(env) > 0 {
		container["env"] = env
	}
	if envFrom := kubernetesEnvFromVars(service); len(envFrom) > 0 {
		container["envFrom"] = envFrom
	}

	// Add ports
	if len(service.Ports) > 0 {
		var ports []map[string]interface{}
		for _, port := range service.Ports {
			portSpec := map[string]interface{}{
				"containerPort": parseInt(port.ContainerPort),
			}
			if port.Name != "" {
				portSpec["name"] = port.Name
			}
			if port.Protocol != "" && port.Protocol != "tcp" {
				portSpec["protocol"] = strings.ToUpper(port.Protocol)
			}
			if port.AppProtocol != "" {
				portSpec["appProtocol"] = port.AppProtocol
			}
			for key, value := range port.Extensions {
				portSpec[key] = deepCopyValue(value)
			}
			ports = append(ports, portSpec)
		}
		container["ports"] = ports
	}

	if resources := kubernetesServiceResources(service); resources != nil {
		container["resources"] = serializeKubernetesResources(resources)
	}
	if lifecycle := serializeKubernetesLifecycle(service.Lifecycle); len(lifecycle) > 0 {
		container["lifecycle"] = lifecycle
	}
	if probe := serializeKubernetesProbe(service.StartupProbe); probe != nil {
		container["startupProbe"] = probe
	}
	if probe := serializeKubernetesProbe(service.HealthCheck); probe != nil {
		container["readinessProbe"] = probe
		container["livenessProbe"] = cloneMap(probe)
	}

	volumeMounts, volumes := kubernetesFileAndVolumeMounts(service)
	if len(volumeMounts) > 0 {
		container["volumeMounts"] = volumeMounts
	}
	if len(volumes) > 0 {
		podSpec["volumes"] = volumes
	}

	data, err := yaml.Marshal(deployment)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func kubernetesWorkloadSkeleton(name string, service *Service, replicas int, namespace string) map[string]interface{} {
	kind := kubernetesWorkloadKind(service)
	template := map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"app": name,
			},
		},
		"spec": map[string]interface{}{
			"containers": []map[string]interface{}{
				{
					"name":  name,
					"image": service.Image,
				},
			},
		},
	}
	switch kind {
	case "DaemonSet":
		return map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": name,
					},
				},
				"template": template,
			},
		}
	case "StatefulSet":
		serviceName := extensionStringValue(service, "kubernetes.statefulSet.serviceName")
		if serviceName == "" {
			serviceName = name
		}
		spec := map[string]interface{}{
			"replicas":    replicas,
			"serviceName": serviceName,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app": name,
				},
			},
			"template": template,
		}
		if podManagementPolicy := extensionStringValue(service, "kubernetes.statefulSet.podManagementPolicy"); podManagementPolicy != "" {
			spec["podManagementPolicy"] = podManagementPolicy
		}
		if ordinals, ok := extensionValueForKey(service, "kubernetes.statefulSet.ordinals"); ok {
			if ordinalSpec, ok := extensionMapValues(ordinals); ok && len(ordinalSpec) > 0 {
				spec["ordinals"] = ordinalSpec
			}
		}
		if retentionPolicy, ok := extensionValueForKey(service, "kubernetes.statefulSet.persistentVolumeClaimRetentionPolicy"); ok {
			if policy, ok := asMap(retentionPolicy); ok && len(policy) > 0 {
				spec["persistentVolumeClaimRetentionPolicy"] = cloneMap(policy)
			}
		}
		if templates, ok := extensionValueForKey(service, "kubernetes.statefulSet.volumeClaimTemplates"); ok {
			if claimTemplates, ok := kubernetesMapSliceFromExtension(templates); ok && len(claimTemplates) > 0 {
				spec["volumeClaimTemplates"] = cloneMapSlice(claimTemplates)
			}
		}
		return map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": spec,
		}
	case "Job":
		spec := map[string]interface{}{
			"template": template,
		}
		applyKubernetesJobSpecExtensions(spec, service)
		return map[string]interface{}{
			"apiVersion": "batch/v1",
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": spec,
		}
	case "CronJob":
		jobSpec := map[string]interface{}{
			"template": template,
		}
		applyKubernetesJobSpecExtensions(jobSpec, service)
		spec := map[string]interface{}{
			"schedule": kubernetesCronSchedule(service),
			"jobTemplate": map[string]interface{}{
				"spec": jobSpec,
			},
		}
		if concurrencyPolicy := extensionStringValue(service, "kubernetes.cron.concurrencyPolicy"); concurrencyPolicy != "" {
			spec["concurrencyPolicy"] = concurrencyPolicy
		}
		if suspend := extensionStringValue(service, "kubernetes.cron.suspend"); suspend != "" {
			spec["suspend"] = strings.EqualFold(suspend, "true")
		}
		applyKubernetesCronJobSpecExtensions(spec, service)
		return map[string]interface{}{
			"apiVersion": "batch/v1",
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": spec,
		}
	case "Pod":
		return map[string]interface{}{
			"apiVersion": "v1",
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": template["spec"],
		}
	default:
		return map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"replicas": replicas,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": name,
					},
				},
				"template": template,
			},
		}
	}
}

func applyKubernetesWorkloadSpecExtensions(workload map[string]interface{}, service *Service) {
	if service == nil || len(service.Extensions) == 0 {
		return
	}
	kind := toString(workload["kind"])
	spec, ok := asMap(workload["spec"])
	if !ok {
		return
	}
	if value, ok := extensionIntValue(service, "kubernetes.workload.revisionHistoryLimit"); ok {
		spec["revisionHistoryLimit"] = value
	}
	if value, ok := extensionIntValue(service, "kubernetes.workload.minReadySeconds"); ok {
		spec["minReadySeconds"] = value
	}
	switch kind {
	case "Deployment":
		if value, ok := extensionBoolValue(service, "kubernetes.deployment.paused"); ok {
			spec["paused"] = value
		}
		if value, ok := extensionIntValue(service, "kubernetes.deployment.progressDeadlineSeconds"); ok {
			spec["progressDeadlineSeconds"] = value
		}
		if value, ok := extensionValueForKey(service, "kubernetes.deployment.strategy"); ok {
			if strategy, ok := extensionMapValues(value); ok && len(strategy) > 0 {
				spec["strategy"] = strategy
			}
		}
	case "StatefulSet", "DaemonSet":
		if service.Deploy != nil {
			if value, ok := updatePolicyExtensionValue(service.Deploy.UpdateConfig, "x-kubernetes-workload-updateStrategy", "kubernetes.workload.updateStrategy"); ok {
				if strategy, ok := extensionMapValues(value); ok && len(strategy) > 0 {
					spec["updateStrategy"] = strategy
				}
			} else if value, ok := extensionValueForKey(service, "kubernetes.workload.updateStrategy"); ok {
				if strategy, ok := extensionMapValues(value); ok && len(strategy) > 0 {
					spec["updateStrategy"] = strategy
				}
			}
		} else if value, ok := extensionValueForKey(service, "kubernetes.workload.updateStrategy"); ok {
			if strategy, ok := extensionMapValues(value); ok && len(strategy) > 0 {
				spec["updateStrategy"] = strategy
			}
		}
		if kind == "StatefulSet" {
			if value, ok := extensionValueForKey(service, "kubernetes.statefulSet.ordinals"); ok {
				if ordinals, ok := extensionMapValues(value); ok && len(ordinals) > 0 {
					spec["ordinals"] = ordinals
				}
			}
			if value, ok := extensionValueForKey(service, "kubernetes.statefulSet.persistentVolumeClaimRetentionPolicy"); ok {
				if policy, ok := extensionMapValues(value); ok && len(policy) > 0 {
					spec["persistentVolumeClaimRetentionPolicy"] = policy
				}
			}
			if value, ok := extensionValueForKey(service, "kubernetes.statefulSet.volumeClaimTemplates"); ok {
				if claimTemplates, ok := kubernetesMapSliceFromExtension(value); ok && len(claimTemplates) > 0 {
					spec["volumeClaimTemplates"] = cloneMapSlice(claimTemplates)
				}
			}
		}
	}
}

func applyKubernetesWorkloadMetadata(workload map[string]interface{}, service *Service) {
	if workload == nil || service == nil {
		return
	}
	if metadata, ok := asMap(workload["metadata"]); ok {
		if labels := serviceExtensionStringMap(service, "kubernetes.workload.labels", "x-kubernetes-workload-labels"); len(labels) > 0 {
			metadata["labels"] = labels
		}
		if annotations := serviceExtensionStringMap(service, "kubernetes.workload.annotations", "x-kubernetes-workload-annotations"); len(annotations) > 0 {
			metadata["annotations"] = annotations
		}
	}
	if template := kubernetesSerializedPodTemplate(workload); template != nil {
		metadata, _ := asMap(template["metadata"])
		if metadata == nil {
			metadata = map[string]interface{}{}
			template["metadata"] = metadata
		}
		if labels := serviceExtensionStringMap(service, "kubernetes.template.labels", "x-kubernetes-template-labels"); len(labels) > 0 {
			metadata["labels"] = labels
		}
		if annotations := serviceExtensionStringMap(service, "kubernetes.template.annotations", "x-kubernetes-template-annotations"); len(annotations) > 0 {
			metadata["annotations"] = annotations
		}
	}
}

func serviceExtensionStringMap(service *Service, keys ...string) map[string]string {
	if service == nil || service.Extensions == nil {
		return nil
	}
	for _, key := range keys {
		if mapped := toStringMapLoose(service.Extensions[key]); len(mapped) > 0 {
			return mapped
		}
	}
	return nil
}

func applyKubernetesJobSpecExtensions(spec map[string]interface{}, service *Service) {
	if service == nil || len(service.Extensions) == 0 || spec == nil {
		if service == nil || spec == nil {
			return
		}
	} else {
		for _, key := range []string{"parallelism", "completions", "backoffLimit", "backoffLimitPerIndex", "ttlSecondsAfterFinished"} {
			if value, ok := extensionIntValue(service, "kubernetes.job."+key); ok {
				spec[key] = value
			}
		}
		if completionMode := extensionStringValue(service, "kubernetes.job.completionMode"); completionMode != "" {
			spec["completionMode"] = completionMode
		}
		if value, ok := extensionValueForKey(service, "kubernetes.job.podFailurePolicy"); ok {
			if policy, ok := extensionMapValues(value); ok && len(policy) > 0 {
				spec["podFailurePolicy"] = policy
			}
		}
		if _, ok := extensionValueForKey(service, "kubernetes.job.podFailurePolicy"); !ok {
			if value, ok := firstExtensionValue(service.Extensions, "x-kubernetes-job-podFailurePolicy", "x-kubernetes-job-pod-failure-policy"); ok {
				if policy, ok := extensionMapValues(value); ok && len(policy) > 0 {
					spec["podFailurePolicy"] = policy
				}
			}
		}
		if value, ok := extensionValueForKey(service, "kubernetes.job.successPolicy"); ok {
			if policy, ok := extensionMapValues(value); ok && len(policy) > 0 {
				spec["successPolicy"] = policy
			}
		}
		if _, ok := extensionValueForKey(service, "kubernetes.job.successPolicy"); !ok {
			if value, ok := firstExtensionValue(service.Extensions, "x-kubernetes-job-successPolicy", "x-kubernetes-job-success-policy"); ok {
				if policy, ok := extensionMapValues(value); ok && len(policy) > 0 {
					spec["successPolicy"] = policy
				}
			}
		}
		if value, ok := extensionBoolValue(service, "kubernetes.job.suspend"); ok {
			spec["suspend"] = value
		}
	}
	if service.Deploy != nil && !isEmptySwarmJobSpec(service.Deploy.Job) {
		if service.Deploy.Job.maxConcurrentSet || service.Deploy.Job.MaxConcurrent > 0 {
			spec["parallelism"] = service.Deploy.Job.MaxConcurrent
		}
		if service.Deploy.Job.totalCompletionsSet || service.Deploy.Job.TotalCompletions > 0 {
			spec["completions"] = service.Deploy.Job.TotalCompletions
		}
		if service.Deploy.Job.completionModeSet || service.Deploy.Job.CompletionMode != "" {
			spec["completionMode"] = service.Deploy.Job.CompletionMode
		}
		if value, ok := swarmJobExtensionValue(service, "kubernetes.job.podFailurePolicy", "x-kubernetes-job-podFailurePolicy"); ok {
			if policy, ok := extensionMapValues(value); ok && len(policy) > 0 {
				spec["podFailurePolicy"] = policy
			}
		}
		if value, ok := swarmJobExtensionValue(service, "kubernetes.job.successPolicy", "x-kubernetes-job-successPolicy"); ok {
			if policy, ok := extensionMapValues(value); ok && len(policy) > 0 {
				spec["successPolicy"] = policy
			}
		}
		if service.Deploy.Job.Suspend != nil {
			spec["suspend"] = *service.Deploy.Job.Suspend
		}
		if service.Deploy.Job.backoffLimitSet || service.Deploy.Job.BackoffLimit > 0 {
			spec["backoffLimit"] = service.Deploy.Job.BackoffLimit
		}
		if service.Deploy.Job.backoffLimitPerIndexSet || service.Deploy.Job.BackoffLimitPerIndex > 0 {
			spec["backoffLimitPerIndex"] = service.Deploy.Job.BackoffLimitPerIndex
		}
		if service.Deploy.Job.ttlSecondsAfterFinishedSet || service.Deploy.Job.TTLSecondsAfterFinished > 0 {
			spec["ttlSecondsAfterFinished"] = service.Deploy.Job.TTLSecondsAfterFinished
		}
	}
}

func swarmJobExtensionValue(service *Service, keys ...string) (interface{}, bool) {
	if service == nil || service.Deploy == nil || service.Deploy.Job == nil {
		return nil, false
	}
	if value, ok := firstExtensionValue(service.Deploy.Job.Extensions, keys...); ok {
		return value, true
	}
	return firstExtensionValue(service.Extensions, keys...)
}

func applyKubernetesCronJobSpecExtensions(spec map[string]interface{}, service *Service) {
	if service == nil || len(service.Extensions) == 0 || spec == nil {
		return
	}
	for _, key := range []string{"startingDeadlineSeconds", "successfulJobsHistoryLimit", "failedJobsHistoryLimit"} {
		if value, ok := extensionIntValue(service, "kubernetes.cron."+key); ok {
			spec[key] = value
		}
	}
}

func kubernetesWorkloadKind(service *Service) string {
	switch extensionStringValue(service, "kubernetes.kind") {
	case "Pod", "Deployment", "ReplicaSet", "ReplicationController", "StatefulSet", "DaemonSet", "Job", "CronJob":
		return extensionStringValue(service, "kubernetes.kind")
	default:
		if service != nil && service.Deploy != nil && isSwarmJobMode(service.Deploy.Mode) {
			return "Job"
		}
		if service != nil && service.Deploy != nil && strings.EqualFold(service.Deploy.Mode, "global") {
			return "DaemonSet"
		}
		return "Deployment"
	}
}

func kubernetesCronSchedule(service *Service) string {
	if schedule := extensionStringValue(service, "kubernetes.cron.schedule"); schedule != "" {
		return schedule
	}
	return "* * * * *"
}

func kubernetesSerializedPodTemplate(workload map[string]interface{}) map[string]interface{} {
	if workload["kind"] == "Pod" {
		spec, _ := asMap(workload["spec"])
		return map[string]interface{}{"spec": spec}
	}
	spec, _ := asMap(workload["spec"])
	if workload["kind"] == "CronJob" {
		jobTemplate, _ := asMap(spec["jobTemplate"])
		jobSpec, _ := asMap(jobTemplate["spec"])
		template, _ := asMap(jobSpec["template"])
		return template
	}
	template, _ := asMap(spec["template"])
	return template
}

func kubernetesEnvVars(service *Service) []map[string]interface{} {
	var env []map[string]interface{}
	for key, value := range service.Environment {
		env = append(env, map[string]interface{}{
			"name":  key,
			"value": value,
		})
	}
	for _, source := range service.EnvSources {
		if source.Name == "" || source.SourceType == "" || source.Source == "" {
			continue
		}
		item := map[string]interface{}{"name": source.Name}
		if raw, ok := source.Extensions["kubernetes.valueFrom"].(map[string]interface{}); ok {
			valueFrom := cloneMap(raw)
			switch source.SourceType {
			case "config":
				if ref, ok := asMap(valueFrom["configMapKeyRef"]); ok {
					ref["name"] = source.Source
					if source.Key != "" {
						ref["key"] = source.Key
					}
					if source.Optional {
						ref["optional"] = true
					}
					valueFrom["configMapKeyRef"] = ref
				}
			case "secret":
				if ref, ok := asMap(valueFrom["secretKeyRef"]); ok {
					ref["name"] = source.Source
					if source.Key != "" {
						ref["key"] = source.Key
					}
					if source.Optional {
						ref["optional"] = true
					}
					valueFrom["secretKeyRef"] = ref
				}
			case "field":
				if ref, ok := asMap(valueFrom["fieldRef"]); ok {
					ref["fieldPath"] = source.Source
					valueFrom["fieldRef"] = ref
				}
			case "resource":
				if ref, ok := asMap(valueFrom["resourceFieldRef"]); ok {
					ref["resource"] = source.Source
					valueFrom["resourceFieldRef"] = ref
				}
			default:
				continue
			}
			item["valueFrom"] = valueFrom
			for key, value := range source.Extensions {
				if key == "kubernetes.valueFrom" {
					continue
				}
				item[key] = deepCopyValue(value)
			}
			env = append(env, item)
			continue
		}
		valueFrom := map[string]interface{}{}
		ref := map[string]interface{}{}
		switch source.SourceType {
		case "config":
			ref["name"] = source.Source
			if source.Key != "" {
				ref["key"] = source.Key
			}
			if source.Optional {
				ref["optional"] = true
			}
			valueFrom["configMapKeyRef"] = ref
		case "secret":
			ref["name"] = source.Source
			if source.Key != "" {
				ref["key"] = source.Key
			}
			if source.Optional {
				ref["optional"] = true
			}
			valueFrom["secretKeyRef"] = ref
		case "field":
			ref["fieldPath"] = source.Source
			valueFrom["fieldRef"] = ref
		case "resource":
			ref["resource"] = source.Source
			valueFrom["resourceFieldRef"] = ref
		default:
			continue
		}
		item["valueFrom"] = valueFrom
		for key, value := range source.Extensions {
			if key == "kubernetes.valueFrom" {
				continue
			}
			item[key] = deepCopyValue(value)
		}
		env = append(env, item)
	}
	return env
}

func kubernetesSecurityContext(service *Service) map[string]interface{} {
	securityContext := map[string]interface{}{}
	if service.PrivilegedSet {
		securityContext["privileged"] = service.Privileged
	} else if service.Privileged {
		securityContext["privileged"] = true
	}
	if service.ReadOnlyRootFSSet {
		securityContext["readOnlyRootFilesystem"] = service.ReadOnlyRootFS
	} else if service.ReadOnlyRootFS {
		securityContext["readOnlyRootFilesystem"] = true
	}
	if service.AllowPrivilegeEscalation != nil {
		securityContext["allowPrivilegeEscalation"] = *service.AllowPrivilegeEscalation
	} else if value, ok := extensionBoolValue(service, "kubernetes.allowPrivilegeEscalation"); ok {
		securityContext["allowPrivilegeEscalation"] = value
	}
	if service.ProcMount != "" {
		securityContext["procMount"] = service.ProcMount
	} else if value, ok := extensionValueForKey(service, "kubernetes.procMount"); ok {
		if text := toString(value); text != "" {
			securityContext["procMount"] = text
		}
	}
	userID, groupID := parseRuntimeUserGroup(service.User, service.Group)
	if userID > 0 {
		securityContext["runAsUser"] = userID
	}
	if groupID > 0 {
		securityContext["runAsGroup"] = groupID
	}
	capabilities := map[string]interface{}{}
	if len(service.CapAdd) > 0 {
		capabilities["add"] = service.CapAdd
	}
	if len(service.CapDrop) > 0 {
		capabilities["drop"] = service.CapDrop
	}
	if len(capabilities) > 0 {
		securityContext["capabilities"] = capabilities
	}
	if selinux := serializeKubernetesSELinuxOptions(service.SELinuxOptions); len(selinux) > 0 {
		securityContext["seLinuxOptions"] = selinux
	} else if value, ok := extensionValueForKey(service, "kubernetes.seLinuxOptions"); ok {
		if mapped, ok := extensionMapValues(value); ok && len(mapped) > 0 {
			securityContext["seLinuxOptions"] = mapped
		}
	}
	if windows := serializeKubernetesWindowsSecurityContextOptions(service.WindowsOptions); len(windows) > 0 {
		securityContext["windowsOptions"] = windows
	} else if value, ok := extensionValueForKey(service, "kubernetes.windowsOptions"); ok {
		if mapped, ok := extensionMapValues(value); ok && len(mapped) > 0 {
			securityContext["windowsOptions"] = mapped
		}
	}
	if seccompProfile := serializeKubernetesSeccompProfile(service.SeccompProfile); len(seccompProfile) > 0 {
		securityContext["seccompProfile"] = seccompProfile
	} else if value, ok := extensionValueForKey(service, "kubernetes.seccompProfile"); ok {
		if profile, ok := extensionMapValues(value); ok && len(profile) > 0 {
			securityContext["seccompProfile"] = profile
		}
	}
	return securityContext
}

func kubernetesPodSecurityContext(service *Service) map[string]interface{} {
	securityContext := map[string]interface{}{}
	if service.FSGroup != nil && *service.FSGroup > 0 {
		securityContext["fsGroup"] = *service.FSGroup
	}
	if selinux := serializeKubernetesSELinuxOptions(service.SELinuxOptions); len(selinux) > 0 {
		securityContext["seLinuxOptions"] = selinux
	} else if value, ok := extensionValueForKey(service, "kubernetes.seLinuxOptions"); ok {
		if mapped, ok := extensionMapValues(value); ok && len(mapped) > 0 {
			securityContext["seLinuxOptions"] = mapped
		}
	}
	if windows := serializeKubernetesWindowsSecurityContextOptions(service.WindowsOptions); len(windows) > 0 {
		securityContext["windowsOptions"] = windows
	} else if value, ok := extensionValueForKey(service, "kubernetes.windowsOptions"); ok {
		if mapped, ok := extensionMapValues(value); ok && len(mapped) > 0 {
			securityContext["windowsOptions"] = mapped
		}
	}
	if service.FSGroupChangePolicy != "" {
		securityContext["fsGroupChangePolicy"] = service.FSGroupChangePolicy
	} else if value := extensionStringValue(service, "kubernetes.fsGroupChangePolicy"); value != "" {
		securityContext["fsGroupChangePolicy"] = value
	}
	if service.RunAsNonRoot != nil {
		securityContext["runAsNonRoot"] = *service.RunAsNonRoot
	}
	if len(service.SupplementalGroups) > 0 {
		var supplementalGroups []int64
		for _, group := range service.SupplementalGroups {
			if group > 0 {
				supplementalGroups = append(supplementalGroups, group)
			}
		}
		if len(supplementalGroups) > 0 {
			securityContext["supplementalGroups"] = supplementalGroups
		}
	}
	if service.SupplementalGroupsPolicy != "" {
		securityContext["supplementalGroupsPolicy"] = service.SupplementalGroupsPolicy
	} else if value := extensionStringValue(service, "kubernetes.supplementalGroupsPolicy"); value != "" {
		securityContext["supplementalGroupsPolicy"] = value
	}
	if len(service.Sysctls) > 0 {
		securityContext["sysctls"] = kubernetesSysctlsFromMap(service.Sysctls)
	}
	return securityContext
}

func kubernetesSysctlsToMap(items []interface{}) map[string]string {
	if len(items) == 0 {
		return nil
	}
	result := map[string]string{}
	for _, item := range items {
		entry, ok := asMap(item)
		if !ok {
			continue
		}
		name := toString(entry["name"])
		value := toString(entry["value"])
		if name == "" || value == "" {
			continue
		}
		result[name] = value
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func kubernetesSysctlsFromMap(values map[string]string) []map[string]interface{} {
	if len(values) == 0 {
		return nil
	}
	items := make([]map[string]interface{}, 0, len(values))
	for _, key := range sortedMapKeys(values) {
		value := values[key]
		if key == "" || value == "" {
			continue
		}
		items = append(items, map[string]interface{}{
			"name":  key,
			"value": value,
		})
	}
	if len(items) == 0 {
		return nil
	}
	return items
}

func parseRuntimeUserGroup(user, group string) (int, int) {
	user = strings.TrimSpace(user)
	group = strings.TrimSpace(group)
	if before, after, ok := strings.Cut(user, ":"); ok {
		user = before
		if group == "" {
			group = after
		}
	}
	return parseInt(user), parseInt(group)
}

func applyKubernetesDNS(deployment map[string]interface{}, service *Service) {
	podSpec := kubernetesSerializedPodTemplate(deployment)["spec"].(map[string]interface{})
	dnsConfig := map[string]interface{}{}
	if service.DNSPolicy != "" {
		podSpec["dnsPolicy"] = service.DNSPolicy
	} else if len(service.DNS) > 0 {
		podSpec["dnsPolicy"] = "None"
	}
	if service.OSName != "" {
		podSpec["os"] = map[string]interface{}{"name": service.OSName}
	} else if value := extensionStringValue(service, "kubernetes.os"); value != "" {
		podSpec["os"] = map[string]interface{}{"name": value}
	}
	if len(service.DNS) > 0 {
		dnsConfig["nameservers"] = service.DNS
	}
	if len(service.DNSSearch) > 0 {
		dnsConfig["searches"] = service.DNSSearch
	}
	if len(service.DNSOptions) > 0 {
		var options []map[string]interface{}
		for _, option := range service.DNSOptions {
			name, value := splitDNSOption(option)
			if name == "" {
				continue
			}
			item := map[string]interface{}{"name": name}
			if value != "" {
				item["value"] = value
			}
			options = append(options, item)
		}
		if len(options) > 0 {
			dnsConfig["options"] = options
		}
	}
	if len(dnsConfig) > 0 {
		podSpec["dnsConfig"] = dnsConfig
	}
	if aliases := kubernetesHostAliasesFromService(service); len(aliases) > 0 {
		podSpec["hostAliases"] = aliases
	}
}

func kubernetesHostAliasesFromService(service *Service) []map[string]interface{} {
	if service == nil {
		return nil
	}
	if len(service.HostAliases) > 0 {
		aliases := make([]map[string]interface{}, 0, len(service.HostAliases))
		for _, alias := range service.HostAliases {
			ip := strings.TrimSpace(alias.IP)
			if ip == "" || len(alias.Hostnames) == 0 {
				continue
			}
			item := map[string]interface{}{
				"ip":        ip,
				"hostnames": append([]string{}, alias.Hostnames...),
			}
			for key, value := range alias.Extensions {
				item[key] = deepCopyValue(value)
			}
			aliases = append(aliases, item)
		}
		if len(aliases) > 0 {
			return aliases
		}
	}
	return kubernetesHostAliases(service.ExtraHosts)
}

func kubernetesHostAliasesFromExtension(value interface{}) []HostAlias {
	items, ok := value.([]interface{})
	if !ok {
		return nil
	}
	var aliases []HostAlias
	for _, item := range items {
		aliasMap, ok := asMap(item)
		if !ok {
			continue
		}
		ip := toString(aliasMap["ip"])
		if ip == "" {
			continue
		}
		alias := HostAlias{IP: ip}
		if hostnames, ok := aliasMap["hostnames"].([]interface{}); ok {
			alias.Hostnames = interfaceSliceToStringSlice(hostnames)
		}
		for key, value := range aliasMap {
			if key == "ip" || key == "hostnames" {
				continue
			}
			if alias.Extensions == nil {
				alias.Extensions = map[string]interface{}{}
			}
			alias.Extensions[key] = deepCopyValue(value)
		}
		if len(alias.Hostnames) > 0 {
			aliases = append(aliases, alias)
		}
	}
	return aliases
}

func splitDNSOption(option string) (string, string) {
	option = strings.TrimSpace(option)
	if option == "" {
		return "", ""
	}
	for _, sep := range []string{":", "="} {
		if name, value, ok := strings.Cut(option, sep); ok {
			return strings.TrimSpace(name), strings.TrimSpace(value)
		}
	}
	return option, ""
}

func kubernetesHostAliases(extraHosts []string) []map[string]interface{} {
	byIP := map[string][]string{}
	for _, entry := range extraHosts {
		hostnames, ip := parseExtraHostEntry(entry)
		if ip == "" || len(hostnames) == 0 {
			continue
		}
		for _, hostname := range hostnames {
			values := byIP[ip]
			appendUniqueString(&values, hostname)
			byIP[ip] = values
		}
	}
	var aliases []map[string]interface{}
	for _, ip := range sortedHostAliasIPs(byIP) {
		aliases = append(aliases, map[string]interface{}{
			"ip":        ip,
			"hostnames": byIP[ip],
		})
	}
	return aliases
}

func parseExtraHostEntry(entry string) ([]string, string) {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return nil, ""
	}
	if host, ip, ok := strings.Cut(entry, "="); ok {
		return strings.Fields(strings.TrimSpace(host)), strings.TrimSpace(ip)
	}
	if host, ip, ok := strings.Cut(entry, ":"); ok && strings.Count(ip, ":") == 0 {
		return strings.Fields(strings.TrimSpace(host)), strings.TrimSpace(ip)
	}
	parts := strings.Fields(entry)
	if len(parts) < 2 {
		return nil, ""
	}
	ip := parts[0]
	return parts[1:], ip
}

func sortedHostAliasIPs(input map[string][]string) []string {
	ips := make([]string, 0, len(input))
	for ip := range input {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	return ips
}

func kubernetesEnvFromVars(service *Service) []map[string]interface{} {
	var envFrom []map[string]interface{}
	for _, source := range service.EnvFrom {
		if source.SourceType == "" || source.Source == "" {
			continue
		}
		var item map[string]interface{}
		if raw, ok := source.Extensions["kubernetes.envFrom"].(map[string]interface{}); ok {
			item = cloneMap(raw)
		} else {
			item = map[string]interface{}{}
		}
		if source.Prefix != "" {
			item["prefix"] = source.Prefix
		}
		ref := map[string]interface{}{"name": source.Source}
		if source.Optional {
			ref["optional"] = true
		}
		switch source.SourceType {
		case "config":
			item["configMapRef"] = ref
		case "secret":
			item["secretRef"] = ref
		default:
			continue
		}
		for key, value := range source.Extensions {
			if key == "kubernetes.envFrom" {
				continue
			}
			item[key] = deepCopyValue(value)
		}
		envFrom = append(envFrom, item)
	}
	return envFrom
}

func applyKubernetesDependencies(deployment map[string]interface{}, service *Service, app *Application) {
	dependencies := serviceDependencies(service)
	if len(dependencies) == 0 && len(service.InitContainers) == 0 {
		return
	}
	podTemplate := kubernetesSerializedPodTemplate(deployment)
	podMetadata := podTemplate["metadata"].(map[string]interface{})
	annotations, _ := podMetadata["annotations"].(map[string]interface{})
	if annotations == nil {
		annotations = map[string]interface{}{}
		podMetadata["annotations"] = annotations
	}
	if len(dependencies) > 0 {
		if encoded, err := json.Marshal(dependencies); err == nil {
			annotations[kubernetesAnnotationDependencies] = string(encoded)
		}
	}

	podSpec := podTemplate["spec"].(map[string]interface{})
	initContainers := cloneMapSlice(service.InitContainers)
	seenInitContainerNames := map[string]struct{}{}
	for _, initContainer := range initContainers {
		name := toString(initContainer["name"])
		if name != "" {
			seenInitContainerNames[name] = struct{}{}
		}
	}
	for _, dependency := range dependencies {
		if dependency.Name == "" {
			continue
		}
		var container map[string]interface{}
		if dependencyService := kubernetesDependencyService(app, service, dependency.Name); dependencyService != nil {
			container = kubernetesInitContainerFromService(dependency.Name, dependencyService)
		} else {
			container = kubernetesDependencyInitContainer(dependency)
		}
		if container == nil {
			continue
		}
		name := toString(container["name"])
		if name != "" {
			if _, ok := seenInitContainerNames[name]; ok {
				continue
			}
			seenInitContainerNames[name] = struct{}{}
		}
		initContainers = append(initContainers, container)
	}
	if len(initContainers) > 0 {
		podSpec["initContainers"] = initContainers
	}
}

func kubernetesDependencyService(app *Application, owner *Service, dependencyName string) *Service {
	if app == nil || owner == nil || dependencyName == "" {
		return nil
	}
	dependency := app.Services[dependencyName]
	if !isKubernetesInitContainerService(dependency) {
		return nil
	}
	ownerWorkload := owner.Name
	if workload := extensionStringValue(owner, "kubernetes.workload"); workload != "" {
		ownerWorkload = workload
	}
	if dependencyWorkload := extensionStringValue(dependency, "kubernetes.workload"); dependencyWorkload != "" && dependencyWorkload != ownerWorkload {
		return nil
	}
	return dependency
}

func isKubernetesInitContainerService(service *Service) bool {
	return strings.EqualFold(extensionStringValue(service, "kubernetes.initContainer"), "true")
}

func extensionStringValue(service *Service, key string) string {
	if service == nil || service.Extensions == nil {
		return ""
	}
	if value, ok := extensionValueForKey(service, key); ok {
		return toString(value)
	}
	return ""
}

func extensionStringSliceValue(service *Service, key string) []string {
	if service == nil || service.Extensions == nil {
		return nil
	}
	if value, ok := extensionValueForKey(service, key); ok {
		return interfaceSliceToStringSliceLoose(value)
	}
	return nil
}

func extensionIntValue(service *Service, key string) (int, bool) {
	if service == nil || service.Extensions == nil {
		return 0, false
	}
	value, ok := extensionValueForKey(service, key)
	if !ok {
		return 0, false
	}
	switch typed := value.(type) {
	case int:
		return typed, true
	case int32:
		return int(typed), true
	case int64:
		return int(typed), true
	case float64:
		return int(typed), true
	case string:
		if strings.TrimSpace(typed) == "" {
			return 0, false
		}
		return parseInt(typed), true
	default:
		return toInt(value), true
	}
}

func extensionBoolValue(service *Service, key string) (bool, bool) {
	if service == nil || service.Extensions == nil {
		return false, false
	}
	value, ok := extensionValueForKey(service, key)
	if !ok {
		return false, false
	}
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		if strings.TrimSpace(typed) == "" {
			return false, false
		}
		return strings.EqualFold(typed, "true"), true
	default:
		return toBool(value), true
	}
}

func extensionValueForKey(service *Service, key string) (interface{}, bool) {
	if service == nil || service.Extensions == nil {
		return nil, false
	}
	if value, ok := service.Extensions[key]; ok {
		return value, true
	}
	switch key {
	case "kubernetes.nodeSelector":
		if value, ok := service.Extensions["x-kubernetes-node-selector"]; ok {
			return value, true
		}
		if value, ok := service.Extensions["x-kubernetes-nodeSelector"]; ok {
			return value, true
		}
	case "kubernetes.resourceClaims":
		if value, ok := service.Extensions["x-kubernetes-resource-claims"]; ok {
			return value, true
		}
		if value, ok := service.Extensions["x-kubernetes-resourceClaims"]; ok {
			return value, true
		}
	case "kubernetes.seccompProfile":
		if value, ok := service.Extensions["x-kubernetes-seccomp-profile"]; ok {
			return value, true
		}
		if value, ok := service.Extensions["x-kubernetes-seccompProfile"]; ok {
			return value, true
		}
	case "kubernetes.readinessGates":
		if value, ok := service.Extensions["x-kubernetes-readiness-gates"]; ok {
			return value, true
		}
		if value, ok := service.Extensions["x-kubernetes-readinessGates"]; ok {
			return value, true
		}
	case "kubernetes.initContainers":
		if value, ok := service.Extensions["x-kubernetes-init-containers"]; ok {
			return value, true
		}
		if value, ok := service.Extensions["x-kubernetes-initContainers"]; ok {
			return value, true
		}
	case "kubernetes.schedulingGates":
		if value, ok := service.Extensions["x-kubernetes-scheduling-gates"]; ok {
			return value, true
		}
		if value, ok := service.Extensions["x-kubernetes-schedulingGates"]; ok {
			return value, true
		}
	case "kubernetes.hostAliases":
		if value, ok := service.Extensions["x-kubernetes-hostAliases"]; ok {
			return value, true
		}
		if value, ok := service.Extensions["x-kubernetes-host-aliases"]; ok {
			return value, true
		}
	}
	if alias := composeServiceExtensionKey(key); alias != key {
		if value, ok := service.Extensions[alias]; ok {
			return value, true
		}
	}
	return nil, false
}

func kubernetesInitContainerFromService(name string, service *Service) map[string]interface{} {
	containerName := extensionStringValue(service, "kubernetes.container")
	if containerName == "" {
		containerName = name
	}
	container := map[string]interface{}{
		"name":  sanitizeKubernetesName(containerName),
		"image": service.Image,
	}
	if len(service.Entrypoint) > 0 {
		container["command"] = service.Entrypoint
	}
	if len(service.Command) > 0 {
		container["args"] = service.Command
	}
	if service.WorkingDir != "" {
		container["workingDir"] = service.WorkingDir
	}
	if env := kubernetesEnvVars(service); len(env) > 0 {
		container["env"] = env
	}
	if envFrom := kubernetesEnvFromVars(service); len(envFrom) > 0 {
		container["envFrom"] = envFrom
	}
	if securityContext := kubernetesSecurityContext(service); len(securityContext) > 0 {
		container["securityContext"] = securityContext
	}
	if resources := kubernetesServiceResources(service); resources != nil {
		container["resources"] = serializeKubernetesResources(resources)
	}
	return container
}

func kubernetesDependencyInitContainer(dependency DependencySpec) map[string]interface{} {
	name := "wait-for-" + sanitizeKubernetesName(dependency.Name)
	host := dependency.Name
	command := fmt.Sprintf("until nslookup %s; do echo waiting for %s; sleep 2; done", host, host)
	return map[string]interface{}{
		"name":    name,
		"image":   "busybox:1.36.1",
		"command": []string{"sh", "-c", command},
	}
}

func kubernetesFileAndVolumeMounts(service *Service) ([]map[string]interface{}, []map[string]interface{}) {
	var volumeMounts []map[string]interface{}
	var volumes []map[string]interface{}
	seen := map[string]struct{}{}

	for _, ref := range service.Configs {
		name := kubernetesRefVolumeName("config", kubernetesRefVolumeNameSource(ref), seen)
		volumeMounts = append(volumeMounts, kubernetesVolumeMount(name, kubernetesProjectedMountPath(ref), ref.ReadOnly, "", "", "", ""))
		volumes = append(volumes, map[string]interface{}{
			"name":      name,
			"configMap": kubernetesProjectedVolumeSource(ref, false),
		})
	}
	for _, ref := range service.Secrets {
		name := kubernetesRefVolumeName("secret", kubernetesRefVolumeNameSource(ref), seen)
		volumeMounts = append(volumeMounts, kubernetesVolumeMount(name, kubernetesProjectedMountPath(ref), true, "", "", "", ""))
		volumes = append(volumes, map[string]interface{}{
			"name":   name,
			"secret": kubernetesProjectedVolumeSource(ref, true),
		})
	}
	for _, vol := range service.Volumes {
		name := vol.Source
		if name == "" {
			name = kubernetesRefVolumeName("tmpfs", vol.Target, seen)
		} else {
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
		}
		volumeMounts = append(volumeMounts, kubernetesVolumeMount(name, vol.Target, vol.ReadOnly, vol.SubPath, vol.SubPathExpr, vol.MountPropagation, vol.RecursiveReadOnly))
		emptyDir := map[string]interface{}{}
		if vol.Type == "tmpfs" {
			emptyDir["medium"] = "Memory"
			if size := vol.Options["size"]; size != "" {
				emptyDir["sizeLimit"] = size
			}
		}
		volumes = append(volumes, map[string]interface{}{
			"name":     name,
			"emptyDir": emptyDir,
		})
	}
	return volumeMounts, volumes
}

func kubernetesRefVolumeNameSource(ref FileRef) string {
	if ref.Key == "" {
		return ref.Source
	}
	return ref.Source + "-" + ref.Key
}

func kubernetesProjectedVolumeSource(ref FileRef, secret bool) map[string]interface{} {
	source := map[string]interface{}{}
	if secret {
		source["secretName"] = ref.Source
	} else {
		source["name"] = ref.Source
	}
	if ref.Optional != nil {
		source["optional"] = *ref.Optional
	}
	if ref.Key == "" {
		if mode := kubernetesModeValue(ref.Mode); mode != nil {
			source["defaultMode"] = mode
		}
		return source
	}
	item := map[string]interface{}{
		"key":  ref.Key,
		"path": kubernetesProjectedItemPath(ref),
	}
	if mode := kubernetesModeValue(ref.Mode); mode != nil {
		item["mode"] = mode
	}
	source["items"] = []map[string]interface{}{item}
	return source
}

func kubernetesProjectedMountPath(ref FileRef) string {
	if ref.Key == "" || ref.Target == "" {
		return ref.Target
	}
	dir := path.Dir(ref.Target)
	if dir == "." {
		return ""
	}
	return dir
}

func kubernetesProjectedItemPath(ref FileRef) string {
	if ref.Target == "" {
		return ref.Key
	}
	base := path.Base(ref.Target)
	if base == "." || base == "/" {
		return ref.Key
	}
	return base
}

func kubernetesModeValue(mode string) interface{} {
	if strings.TrimSpace(mode) == "" {
		return nil
	}
	if parsed, err := strconv.ParseInt(strings.TrimSpace(mode), 0, 32); err == nil {
		return int(parsed)
	}
	return mode
}

func kubernetesVolumeMount(name, target string, readOnly bool, subPath, subPathExpr, mountPropagation, recursiveReadOnly string) map[string]interface{} {
	mount := map[string]interface{}{"name": name}
	if target != "" {
		mount["mountPath"] = target
	}
	if readOnly {
		mount["readOnly"] = true
	}
	if mountPropagation != "" {
		mount["mountPropagation"] = mountPropagation
	}
	if recursiveReadOnly != "" {
		mount["recursiveReadOnly"] = recursiveReadOnly
	}
	if subPath != "" {
		mount["subPath"] = subPath
	}
	if subPathExpr != "" {
		mount["subPathExpr"] = subPathExpr
	}
	return mount
}

func kubernetesRefVolumeName(prefix, source string, seen map[string]struct{}) string {
	base := sanitizeKubernetesName(prefix + "-" + source)
	name := base
	for index := 1; ; index++ {
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			return name
		}
		name = fmt.Sprintf("%s-%d", base, index)
	}
}

func sanitizeKubernetesName(value string) string {
	value = strings.ToLower(value)
	var builder strings.Builder
	lastDash := false
	for _, r := range value {
		valid := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if valid {
			builder.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "ref"
	}
	if len(result) > 63 {
		result = strings.Trim(result[:63], "-")
	}
	if result == "" {
		return "ref"
	}
	return result
}

func serializeKubernetesService(name string, service *Service, namespace string) (string, error) {
	serviceName := extensionStringValue(service, "kubernetes.service.name")
	if serviceName == "" {
		serviceName = name
	}
	k8sService := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"name":      serviceName,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{},
	}
	spec := k8sService["spec"].(map[string]interface{})
	metadata := k8sService["metadata"].(map[string]interface{})
	applyKubernetesServiceObjectMetadata(metadata, service)
	applyKubernetesServiceMetadata(spec, service)
	if _, ok := spec["selector"]; !ok {
		if extensionStringValue(service, "kubernetes.service.type") != "ExternalName" {
			spec["selector"] = map[string]interface{}{
				"app": name,
			}
		}
	}

	if len(service.Ports) > 0 || extensionStringValue(service, "kubernetes.service.type") != "ExternalName" {
		var ports []map[string]interface{}
		for _, port := range service.Ports {
			portSpec := map[string]interface{}{
				"port":       parseInt(port.HostPort),
				"targetPort": kubernetesIntOrString(kubernetesTargetPort(port)),
			}
			if port.Name != "" {
				portSpec["name"] = port.Name
			}
			if port.Protocol != "" && port.Protocol != "tcp" {
				portSpec["protocol"] = strings.ToUpper(port.Protocol)
			}
			if port.NodePort != "" {
				portSpec["nodePort"] = parseInt(port.NodePort)
			}
			if port.AppProtocol != "" {
				portSpec["appProtocol"] = port.AppProtocol
			}
			for key, value := range port.Extensions {
				portSpec[key] = deepCopyValue(value)
			}
			ports = append(ports, portSpec)
		}
		if len(ports) > 0 {
			spec["ports"] = ports
		}
	}

	data, err := yaml.Marshal(k8sService)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func applyKubernetesServiceObjectMetadata(metadata map[string]interface{}, service *Service) {
	if labels := extensionStringMapValue(service, "kubernetes.service.labels"); len(labels) > 0 {
		metadata["labels"] = labels
	}
	if annotations := extensionStringMapValue(service, "kubernetes.service.annotations"); len(annotations) > 0 {
		metadata["annotations"] = annotations
	}
}

func applyKubernetesServiceMetadata(spec map[string]interface{}, service *Service) {
	if selector := extensionStringMapValue(service, "kubernetes.service.selector"); len(selector) > 0 {
		spec["selector"] = selector
	}
	for _, key := range []string{
		"type",
		"clusterIP",
		"externalName",
		"sessionAffinity",
		"loadBalancerIP",
		"loadBalancerClass",
		"ipFamilyPolicy",
		"externalTrafficPolicy",
		"internalTrafficPolicy",
		"trafficDistribution",
	} {
		if value := extensionStringValue(service, "kubernetes.service."+key); value != "" {
			spec[key] = value
		}
	}
	for _, key := range []string{"externalIPs", "ipFamilies", "clusterIPs", "loadBalancerSourceRanges"} {
		if values := extensionStringSliceValue(service, "kubernetes.service."+key); len(values) > 0 {
			spec[key] = values
		}
	}
	for _, key := range []string{"allocateLoadBalancerNodePorts", "publishNotReadyAddresses"} {
		if value, ok := extensionBoolValue(service, "kubernetes.service."+key); ok {
			spec[key] = value
		}
	}
	if value, ok := extensionIntValue(service, "kubernetes.service.healthCheckNodePort"); ok {
		spec["healthCheckNodePort"] = value
	}
	if value, ok := extensionValueForKey(service, "kubernetes.service.sessionAffinityConfig"); ok {
		if sessionAffinityConfig, ok := extensionMapValues(value); ok && len(sessionAffinityConfig) > 0 {
			spec["sessionAffinityConfig"] = sessionAffinityConfig
		}
	}
}

func extensionStringMapValue(service *Service, key string) map[string]string {
	if service == nil || service.Extensions == nil {
		return nil
	}
	return toStringMapLoose(service.Extensions[key])
}

func serializeKubernetesConfigMap(name string, config *Config, namespace string) (string, error) {
	configData := configKubernetesData(config)
	if len(configData) == 0 {
		configData = map[string]string{name + ".yaml": config.Content}
	}
	configMap := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
		"data": configData,
	}
	metadata := configMap["metadata"].(map[string]interface{})
	if labels := configKubernetesLabels(config); len(labels) > 0 {
		metadata["labels"] = labels
	}
	if annotations := configKubernetesAnnotations(config); len(annotations) > 0 {
		metadata["annotations"] = annotations
	}
	if immutable := configKubernetesImmutable(config); immutable != nil {
		configMap["immutable"] = *immutable
	}
	if binaryData := configKubernetesBinaryData(config); len(binaryData) > 0 {
		configMap["binaryData"] = binaryData
	}

	data, err := yaml.Marshal(configMap)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func serializeKubernetesSecret(name string, secret *Secret, namespace string) (string, error) {
	secretData := secretKubernetesData(secret)
	if len(secretData) == 0 {
		secretData = map[string]string{"value": "dGVzdA=="} // base64 encoded "test"
	}
	k8sSecret := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
		"type": secretKubernetesType(secret),
		"data": secretData,
	}
	metadata := k8sSecret["metadata"].(map[string]interface{})
	if labels := secretKubernetesLabels(secret); len(labels) > 0 {
		metadata["labels"] = labels
	}
	if annotations := secretKubernetesAnnotations(secret); len(annotations) > 0 {
		metadata["annotations"] = annotations
	}
	if immutable := secretKubernetesImmutable(secret); immutable != nil {
		k8sSecret["immutable"] = *immutable
	}
	if stringData := secretKubernetesStringData(secret); len(stringData) > 0 {
		k8sSecret["stringData"] = stringData
	}

	data, err := yaml.Marshal(k8sSecret)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func configKubernetesLabels(config *Config) map[string]string {
	if config == nil || config.Extensions == nil {
		return nil
	}
	if labels := toStringMapLoose(config.Extensions["kubernetes.labels"]); len(labels) > 0 {
		return labels
	}
	if labels := toStringMapLoose(config.Extensions["x-kubernetes-labels"]); len(labels) > 0 {
		return labels
	}
	return nil
}

func configKubernetesAnnotations(config *Config) map[string]string {
	if config == nil || config.Extensions == nil {
		return nil
	}
	if annotations := toStringMapLoose(config.Extensions["kubernetes.annotations"]); len(annotations) > 0 {
		return annotations
	}
	if annotations := toStringMapLoose(config.Extensions["x-kubernetes-annotations"]); len(annotations) > 0 {
		return annotations
	}
	return nil
}

func configKubernetesImmutable(config *Config) *bool {
	if config == nil || config.Extensions == nil {
		return nil
	}
	if value, ok := config.Extensions["kubernetes.immutable"]; ok {
		flag := toBool(value)
		return &flag
	}
	if value, ok := config.Extensions["x-kubernetes-immutable"]; ok {
		flag := toBool(value)
		return &flag
	}
	return nil
}

func configKubernetesData(config *Config) map[string]string {
	if config == nil || config.Extensions == nil {
		return nil
	}
	if data := toStringMapLoose(config.Extensions["kubernetes.data"]); len(data) > 0 {
		return data
	}
	if data := toStringMapLoose(config.Extensions["x-kubernetes-data"]); len(data) > 0 {
		return data
	}
	return nil
}

func configKubernetesBinaryData(config *Config) map[string]string {
	if config == nil || config.Extensions == nil {
		return nil
	}
	if data := toStringMapLoose(config.Extensions["kubernetes.binaryData"]); len(data) > 0 {
		return data
	}
	if data := toStringMapLoose(config.Extensions["x-kubernetes-binaryData"]); len(data) > 0 {
		return data
	}
	return nil
}

func secretKubernetesLabels(secret *Secret) map[string]string {
	if secret == nil || secret.Extensions == nil {
		return nil
	}
	if labels := toStringMapLoose(secret.Extensions["kubernetes.labels"]); len(labels) > 0 {
		return labels
	}
	if labels := toStringMapLoose(secret.Extensions["x-kubernetes-labels"]); len(labels) > 0 {
		return labels
	}
	return nil
}

func secretKubernetesAnnotations(secret *Secret) map[string]string {
	if secret == nil || secret.Extensions == nil {
		return nil
	}
	if annotations := toStringMapLoose(secret.Extensions["kubernetes.annotations"]); len(annotations) > 0 {
		return annotations
	}
	if annotations := toStringMapLoose(secret.Extensions["x-kubernetes-annotations"]); len(annotations) > 0 {
		return annotations
	}
	return nil
}

func secretKubernetesImmutable(secret *Secret) *bool {
	if secret == nil || secret.Extensions == nil {
		return nil
	}
	if value, ok := secret.Extensions["kubernetes.immutable"]; ok {
		flag := toBool(value)
		return &flag
	}
	if value, ok := secret.Extensions["x-kubernetes-immutable"]; ok {
		flag := toBool(value)
		return &flag
	}
	return nil
}

func secretKubernetesData(secret *Secret) map[string]string {
	if secret == nil || secret.Extensions == nil {
		return nil
	}
	if data := toStringMapLoose(secret.Extensions["kubernetes.data"]); len(data) > 0 {
		return data
	}
	if data := toStringMapLoose(secret.Extensions["x-kubernetes-data"]); len(data) > 0 {
		return data
	}
	return nil
}

func secretKubernetesStringData(secret *Secret) map[string]string {
	if secret == nil || secret.Extensions == nil {
		return nil
	}
	if data := toStringMapLoose(secret.Extensions["kubernetes.stringData"]); len(data) > 0 {
		return data
	}
	if data := toStringMapLoose(secret.Extensions["x-kubernetes-stringData"]); len(data) > 0 {
		return data
	}
	return nil
}

func secretKubernetesType(secret *Secret) string {
	if secret == nil || secret.Extensions == nil {
		return "Opaque"
	}
	if secretType := toString(secret.Extensions["kubernetes.type"]); secretType != "" {
		return secretType
	}
	if secretType := toString(secret.Extensions["x-kubernetes-type"]); secretType != "" {
		return secretType
	}
	return "Opaque"
}

func serializeKubernetesPVC(name string, volume *Volume, namespace string) (string, error) {
	accessModes := []string{"ReadWriteOnce"}
	if volume.DriverOpts != nil {
		if value := volume.DriverOpts["accessModes"]; value != "" {
			accessModes = strings.Split(value, ",")
		}
	}
	storage := "1Gi"
	if volume.DriverOpts != nil {
		if value := volume.DriverOpts["storage"]; value != "" {
			storage = value
		}
	}
	pvc := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "PersistentVolumeClaim",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"accessModes": accessModes,
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"storage": storage,
				},
			},
		},
	}
	metadata := pvc["metadata"].(map[string]interface{})
	if labels := volumeKubernetesLabels(volume); len(labels) > 0 {
		metadata["labels"] = labels
	}
	if annotations := volumeKubernetesAnnotations(volume); len(annotations) > 0 {
		metadata["annotations"] = annotations
	}
	spec := pvc["spec"].(map[string]interface{})
	if volume.DriverOpts != nil {
		for _, key := range []string{"storageClassName", "volumeName", "volumeMode"} {
			if value := volume.DriverOpts[key]; value != "" {
				spec[key] = value
			}
		}
	}
	if selector := volumeKubernetesMap(volume, "kubernetes.selector", "x-kubernetes-selector"); len(selector) > 0 {
		spec["selector"] = selector
	}
	if dataSource := volumeKubernetesMap(volume, "kubernetes.dataSource", "x-kubernetes-dataSource"); len(dataSource) > 0 {
		spec["dataSource"] = dataSource
	}
	if dataSourceRef := volumeKubernetesMap(volume, "kubernetes.dataSourceRef", "x-kubernetes-dataSourceRef"); len(dataSourceRef) > 0 {
		spec["dataSourceRef"] = dataSourceRef
	}

	data, err := yaml.Marshal(pvc)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func volumeKubernetesKind(volume *Volume) string {
	if volume == nil || volume.Extensions == nil {
		return ""
	}
	if kind := toString(volume.Extensions["kubernetes.kind"]); kind != "" {
		return kind
	}
	return toString(volume.Extensions["x-kubernetes-kind"])
}

func volumeKubernetesLabels(volume *Volume) map[string]string {
	if volume == nil || volume.Extensions == nil {
		return nil
	}
	if labels := toStringMapLoose(volume.Extensions["kubernetes.labels"]); len(labels) > 0 {
		return labels
	}
	if labels := toStringMapLoose(volume.Extensions["x-kubernetes-labels"]); len(labels) > 0 {
		return labels
	}
	return nil
}

func volumeKubernetesAnnotations(volume *Volume) map[string]string {
	if volume == nil || volume.Extensions == nil {
		return nil
	}
	if annotations := toStringMapLoose(volume.Extensions["kubernetes.annotations"]); len(annotations) > 0 {
		return annotations
	}
	if annotations := toStringMapLoose(volume.Extensions["x-kubernetes-annotations"]); len(annotations) > 0 {
		return annotations
	}
	return nil
}

func volumeKubernetesMap(volume *Volume, keys ...string) map[string]interface{} {
	if volume == nil || volume.Extensions == nil {
		return nil
	}
	for _, key := range keys {
		if mapped, ok := asMap(volume.Extensions[key]); ok && len(mapped) > 0 {
			return cloneMap(mapped)
		}
	}
	return nil
}

func volumeKubernetesStringMap(volume *Volume, keys ...string) map[string]string {
	if volume == nil || volume.Extensions == nil {
		return nil
	}
	for _, key := range keys {
		if mapped := toStringMapLoose(volume.Extensions[key]); len(mapped) > 0 {
			return mapped
		}
	}
	return nil
}

func serializeKubernetesPersistentVolume(name string, volume *Volume) (string, error) {
	pv := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "PersistentVolume",
		"metadata": map[string]interface{}{
			"name": name,
		},
		"spec": map[string]interface{}{},
	}
	metadata := pv["metadata"].(map[string]interface{})
	if labels := volumeKubernetesLabels(volume); len(labels) > 0 {
		metadata["labels"] = labels
	}
	if annotations := volumeKubernetesAnnotations(volume); len(annotations) > 0 {
		metadata["annotations"] = annotations
	}
	spec := pv["spec"].(map[string]interface{})
	if volume.DriverOpts != nil {
		if value := volume.DriverOpts["capacity.storage"]; value != "" {
			spec["capacity"] = map[string]interface{}{"storage": value}
		}
		if value := volume.DriverOpts["accessModes"]; value != "" {
			spec["accessModes"] = strings.Split(value, ",")
		}
		for _, key := range []string{"storageClassName", "volumeMode", "persistentVolumeReclaimPolicy"} {
			if value := volume.DriverOpts[key]; value != "" {
				spec[key] = value
			}
		}
		if value := volume.DriverOpts["mountOptions"]; value != "" {
			spec["mountOptions"] = strings.Split(value, ",")
		}
		if claimRef := volume.DriverOpts["claimRef"]; claimRef != "" {
			ref := map[string]interface{}{"name": claimRef}
			if namespace := volume.DriverOpts["claimRefNamespace"]; namespace != "" {
				ref["namespace"] = namespace
			}
			for optKey, refKey := range map[string]string{
				"claimRefKind":            "kind",
				"claimRefAPIVersion":      "apiVersion",
				"claimRefUID":             "uid",
				"claimRefResourceVersion": "resourceVersion",
			} {
				if value := volume.DriverOpts[optKey]; value != "" {
					ref[refKey] = value
				}
			}
			spec["claimRef"] = ref
		}
		switch volume.Driver {
		case "hostPath":
			if path := volume.DriverOpts["path"]; path != "" {
				spec["hostPath"] = map[string]interface{}{"path": path}
			}
		case "nfs":
			nfs := map[string]interface{}{}
			if server := volume.DriverOpts["server"]; server != "" {
				nfs["server"] = server
			}
			if path := volume.DriverOpts["path"]; path != "" {
				nfs["path"] = path
			}
			if len(nfs) > 0 {
				spec["nfs"] = nfs
			}
		case "csi":
			csi := map[string]interface{}{}
			if driver := volume.DriverOpts["driver"]; driver != "" {
				csi["driver"] = driver
			}
			if handle := volume.DriverOpts["volumeHandle"]; handle != "" {
				csi["volumeHandle"] = handle
			}
			if fsType := volume.DriverOpts["csi.fsType"]; fsType != "" {
				csi["fsType"] = fsType
			}
			if readOnly := volume.DriverOpts["csi.readOnly"]; readOnly != "" {
				csi["readOnly"] = strings.EqualFold(readOnly, "true")
			}
			if attributes := volumeKubernetesStringMap(volume, "kubernetes.csi.volumeAttributes", "x-kubernetes-csi-volumeAttributes"); len(attributes) > 0 {
				csi["volumeAttributes"] = attributes
			}
			if len(csi) > 0 {
				spec["csi"] = csi
			}
		case "local":
			if path := volume.DriverOpts["path"]; path != "" {
				spec["local"] = map[string]interface{}{"path": path}
			}
		}
	}
	if nodeAffinity := volumeKubernetesMap(volume, "kubernetes.nodeAffinity", "x-kubernetes-nodeAffinity"); len(nodeAffinity) > 0 {
		spec["nodeAffinity"] = nodeAffinity
	}

	data, err := yaml.Marshal(pv)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Helper functions
func interfaceSliceToStringSlice(slice []interface{}) []string {
	var result []string
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

func interfaceSliceToStringSliceLoose(value interface{}) []string {
	switch typed := value.(type) {
	case []string:
		return append([]string{}, typed...)
	case []interface{}:
		return interfaceSliceToStringSlice(typed)
	default:
		return nil
	}
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int(f)
	}
	return 0
}

func kubernetesTargetPort(port PortMapping) string {
	if port.TargetName != "" {
		return port.TargetName
	}
	return port.ContainerPort
}

func kubernetesPortExtensions(port map[string]interface{}, skip []string) map[string]interface{} {
	if len(port) == 0 {
		return nil
	}
	skipSet := map[string]struct{}{}
	for _, key := range skip {
		skipSet[key] = struct{}{}
	}
	extensions := map[string]interface{}{}
	for key, value := range port {
		if _, ok := skipSet[key]; ok {
			continue
		}
		extensions[key] = deepCopyValue(value)
	}
	if len(extensions) == 0 {
		return nil
	}
	return extensions
}

func kubernetesIntOrString(value string) interface{} {
	if number := parseInt(value); number > 0 {
		return number
	}
	return value
}

func appendExtensionSlice(existing interface{}, value interface{}) []interface{} {
	if slice, ok := existing.([]interface{}); ok {
		return append(slice, value)
	}
	return []interface{}{value}
}

func ensureServiceDeploy(service *Service) *DeploySpec {
	if service.Deploy == nil {
		service.Deploy = &DeploySpec{}
	}
	return service.Deploy
}

func parseKubernetesResources(resources map[string]interface{}) *ResourceSpec {
	spec := &ResourceSpec{}
	if limits, ok := asMap(resources["limits"]); ok {
		if cpu := toString(limits["cpu"]); cpu != "" {
			if normalized, preserveRaw := normalizeKubernetesCPUQuantity(cpu); normalized != "" {
				spec.CPULimit = normalized
				if preserveRaw {
					if spec.Extensions == nil {
						spec.Extensions = map[string]interface{}{}
					}
					spec.Extensions["x-kubernetes-cpu-limit"] = cpu
				}
			} else {
				spec.CPULimit = cpu
			}
		}
		if memory := toString(limits["memory"]); memory != "" {
			if normalized, preserveRaw := normalizeKubernetesMemoryQuantity(memory); normalized != "" {
				spec.MemoryLimit = normalized
				if preserveRaw {
					if spec.Extensions == nil {
						spec.Extensions = map[string]interface{}{}
					}
					spec.Extensions["x-kubernetes-memory-limit"] = memory
				}
			} else {
				spec.MemoryLimit = memory
			}
		}
		if storage := toString(limits["ephemeral-storage"]); storage != "" {
			spec.EphemeralStorageLimit = storage
			if spec.LimitExtensions == nil {
				spec.LimitExtensions = map[string]interface{}{}
			}
			spec.LimitExtensions["x-kubernetes-ephemeral-storage-limit"] = storage
		}
	}
	if requests, ok := asMap(resources["requests"]); ok {
		if cpu := toString(requests["cpu"]); cpu != "" {
			if normalized, preserveRaw := normalizeKubernetesCPUQuantity(cpu); normalized != "" {
				spec.CPUReservation = normalized
				if preserveRaw {
					if spec.Extensions == nil {
						spec.Extensions = map[string]interface{}{}
					}
					spec.Extensions["x-kubernetes-cpu-reservation"] = cpu
				}
			} else {
				spec.CPUReservation = cpu
			}
		}
		if memory := toString(requests["memory"]); memory != "" {
			if normalized, preserveRaw := normalizeKubernetesMemoryQuantity(memory); normalized != "" {
				spec.MemoryReservation = normalized
				if preserveRaw {
					if spec.Extensions == nil {
						spec.Extensions = map[string]interface{}{}
					}
					spec.Extensions["x-kubernetes-memory-reservation"] = memory
				}
			} else {
				spec.MemoryReservation = memory
			}
		}
		if storage := toString(requests["ephemeral-storage"]); storage != "" {
			spec.EphemeralStorageReservation = storage
			if spec.ReservationExtensions == nil {
				spec.ReservationExtensions = map[string]interface{}{}
			}
			spec.ReservationExtensions["x-kubernetes-ephemeral-storage-reservation"] = storage
		}
	}
	if claims, ok := resources["claims"].([]interface{}); ok {
		if len(claims) > 0 {
			if spec.Extensions == nil {
				spec.Extensions = map[string]interface{}{}
			}
			spec.Extensions["kubernetes.claims"] = cloneInterfaceSlice(claims)
		}
	}
	if claims, ok := resources["x-kubernetes-claims"].([]interface{}); ok {
		if len(claims) > 0 {
			if spec.Extensions == nil {
				spec.Extensions = map[string]interface{}{}
			}
			spec.Extensions["kubernetes.claims"] = cloneInterfaceSlice(claims)
		}
	}
	if isEmptyResourceSpec(spec) {
		return nil
	}
	return spec
}

func parseKubernetesDeployIntent(service *Service, kind string, spec, podSpec map[string]interface{}, annotations map[string]string) {
	deploy := ensureServiceDeploy(service)
	if raw := annotations[kubernetesAnnotationDeploySpec]; raw != "" {
		var portable DeploySpec
		if json.Unmarshal([]byte(raw), &portable) == nil && !isEmptyDeploySpec(&portable) {
			service.Deploy = cloneDeploySpec(&portable)
			deploy = ensureServiceDeploy(service)
		}
	}
	if service.Replicas > 0 {
		deploy.Replicas = service.Replicas
	}
	if kind == "DaemonSet" {
		deploy.Mode = "global"
	}
	if mode := annotations[kubernetesAnnotationDeployMode]; mode != "" {
		deploy.Mode = mode
	}
	if endpointMode := annotations[kubernetesAnnotationDeployEndpointMode]; endpointMode != "" {
		deploy.EndpointMode = endpointMode
	}
	if raw := annotations[kubernetesAnnotationDeployLabels]; raw != "" {
		var labels map[string]string
		if json.Unmarshal([]byte(raw), &labels) == nil && len(labels) > 0 {
			deploy.Labels = copyStringMap(labels)
		}
	}
	if raw := annotations[kubernetesAnnotationDeployResources]; raw != "" {
		var resources ResourceSpec
		if json.Unmarshal([]byte(raw), &resources) == nil {
			deploy.Resources = mergeResourceSpec(deploy.Resources, &resources)
		}
	}
	if raw := annotations[kubernetesAnnotationDeployJob]; raw != "" {
		var job SwarmJobSpec
		if json.Unmarshal([]byte(raw), &job) == nil && !isEmptySwarmJobSpec(&job) {
			deploy.Job = &job
		}
	}
	if kind == "Job" {
		if deploy.Mode == "" {
			deploy.Mode = "replicated-job"
		}
		if deploy.Job == nil {
			deploy.Job = &SwarmJobSpec{}
		}
		if value, ok := spec["parallelism"]; ok && deploy.Job.MaxConcurrent == 0 {
			deploy.Job.MaxConcurrent = toInt(value)
		}
		if value, ok := spec["completions"]; ok && deploy.Job.TotalCompletions == 0 {
			deploy.Job.TotalCompletions = toInt(value)
		}
		if isEmptySwarmJobSpec(deploy.Job) {
			deploy.Job = nil
		}
	}

	placement := &PlacementSpec{}
	if nodeSelector, ok := asMap(podSpec["nodeSelector"]); ok {
		if selector, err := toStringMap(nodeSelector); err == nil {
			service.NodeSelector = copyStringMap(selector)
		}
		for key, value := range nodeSelector {
			appendUniqueString(&placement.Constraints, fmt.Sprintf("node.labels.%s == %s", key, toString(value)))
		}
	}
	for _, constraint := range splitAnnotationList(annotations[kubernetesAnnotationDeployPlacement]) {
		appendUniqueString(&placement.Constraints, constraint)
	}
	for _, preference := range splitAnnotationList(annotations[kubernetesAnnotationDeployPreferences]) {
		appendUniqueString(&placement.Preferences, preference)
	}
	if value := annotations[kubernetesAnnotationDeployMaxReplicasPerNode]; value != "" {
		placement.MaxReplicasPerNode = parseInt(value)
	}
	if len(placement.Constraints) > 0 || len(placement.Preferences) > 0 || placement.MaxReplicasPerNode > 0 {
		deploy.Placement = mergeKubernetesPlacement(deploy.Placement, placement)
	}

	if update := parseKubernetesUpdatePolicy(spec, annotations, deploy.UpdateConfig); update != nil {
		deploy.UpdateConfig = update
	}
	if rollback := parseKubernetesRollbackPolicy(annotations, deploy.RollbackConfig); rollback != nil {
		deploy.RollbackConfig = rollback
	}
	if restart := parseKubernetesRestartPolicy(annotations); restart != nil {
		deploy.RestartPolicy = mergeKubernetesRestartPolicy(deploy.RestartPolicy, restart)
	}
	if isEmptyDeploySpec(deploy) {
		service.Deploy = nil
	}
}

func mergeKubernetesPlacement(existing, parsed *PlacementSpec) *PlacementSpec {
	if parsed == nil {
		return clonePlacementSpec(existing)
	}
	merged := clonePlacementSpec(parsed)
	if existing == nil {
		return merged
	}
	if len(merged.Extensions) == 0 && len(existing.Extensions) > 0 {
		merged.Extensions = copyStringInterfaceMap(existing.Extensions)
	}
	if len(merged.PreferenceExtensions) == 0 && len(existing.PreferenceExtensions) > 0 {
		merged.PreferenceExtensions = make([]map[string]interface{}, len(existing.PreferenceExtensions))
		for i, extension := range existing.PreferenceExtensions {
			merged.PreferenceExtensions[i] = copyStringInterfaceMap(extension)
		}
	}
	return merged
}

func mergeKubernetesRestartPolicy(existing, parsed *RestartPolicy) *RestartPolicy {
	if parsed == nil {
		return cloneRestartPolicy(existing)
	}
	merged := cloneRestartPolicy(parsed)
	if len(merged.Extensions) == 0 && existing != nil && len(existing.Extensions) > 0 {
		merged.Extensions = copyStringInterfaceMap(existing.Extensions)
	}
	return merged
}

func parseKubernetesUpdatePolicy(spec map[string]interface{}, annotations map[string]string, existing *UpdatePolicy) *UpdatePolicy {
	update := cloneUpdatePolicy(existing)
	if update == nil {
		update = &UpdatePolicy{}
	}
	if strategy, ok := asMap(spec["strategy"]); ok {
		if update.Extensions == nil {
			update.Extensions = map[string]interface{}{}
		}
		update.Extensions["kubernetes-deployment-strategy"] = cloneMap(strategy)
		if strings.EqualFold(toString(strategy["type"]), "RollingUpdate") {
			if rolling, ok := asMap(strategy["rollingUpdate"]); ok {
				maxUnavailable := toString(rolling["maxUnavailable"])
				maxSurge := toString(rolling["maxSurge"])
				if maxUnavailable == "0" && maxSurge == "1" {
					update.Order = "start-first"
				} else if maxSurge == "0" && maxUnavailable != "" && maxUnavailable != "0" {
					update.Order = "stop-first"
					update.Parallelism = toInt(maxUnavailable)
					update.ParallelismSet = true
				}
			}
		}
	}
	if value := annotations[kubernetesAnnotationDeployUpdateParallelism]; value != "" {
		update.Parallelism = parseInt(value)
		update.ParallelismSet = true
	}
	if value := annotations[kubernetesAnnotationDeployUpdateDelay]; value != "" {
		update.Delay = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateMonitor]; value != "" {
		update.Monitor = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateFailureRate]; value != "" {
		update.MaxFailureRatio = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateOrder]; value != "" {
		update.Order = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateOnFailure]; value != "" {
		update.OnFailure = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateHealthCheck]; value != "" {
		update.HealthCheck = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateMinHealthyTime]; value != "" {
		update.MinHealthyTime = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateHealthyDeadline]; value != "" {
		update.HealthyDeadline = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateProgressDeadline]; value != "" {
		update.ProgressDeadline = value
	}
	if value := annotations[kubernetesAnnotationDeployUpdateAutoRevert]; value != "" {
		update.AutoRevert = toBool(value)
		update.AutoRevertSet = true
	}
	if value := annotations[kubernetesAnnotationDeployUpdateAutoPromote]; value != "" {
		update.AutoPromote = toBool(value)
		update.AutoPromoteSet = true
	}
	if value := annotations[kubernetesAnnotationDeployUpdateCanary]; value != "" {
		update.Canary = parseInt(value)
		update.CanarySet = true
	}
	if value := annotations[kubernetesAnnotationDeployUpdateStagger]; value != "" {
		update.Stagger = value
	}
	if isEmptyUpdatePolicy(update) {
		return nil
	}
	return update
}

func parseKubernetesRollbackPolicy(annotations map[string]string, existing *UpdatePolicy) *UpdatePolicy {
	rollback := cloneUpdatePolicy(existing)
	if rollback == nil {
		rollback = &UpdatePolicy{}
	}
	if value := annotations[kubernetesAnnotationDeployRollbackParallelism]; value != "" {
		rollback.Parallelism = parseInt(value)
		rollback.ParallelismSet = true
	}
	if value := annotations[kubernetesAnnotationDeployRollbackDelay]; value != "" {
		rollback.Delay = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackMonitor]; value != "" {
		rollback.Monitor = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackFailureRate]; value != "" {
		rollback.MaxFailureRatio = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackOrder]; value != "" {
		rollback.Order = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackOnFailure]; value != "" {
		rollback.OnFailure = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackHealthCheck]; value != "" {
		rollback.HealthCheck = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackMinHealthyTime]; value != "" {
		rollback.MinHealthyTime = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackHealthyDeadline]; value != "" {
		rollback.HealthyDeadline = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackProgressDeadline]; value != "" {
		rollback.ProgressDeadline = value
	}
	if value := annotations[kubernetesAnnotationDeployRollbackAutoRevert]; value != "" {
		rollback.AutoRevert = toBool(value)
		rollback.AutoRevertSet = true
	}
	if value := annotations[kubernetesAnnotationDeployRollbackAutoPromote]; value != "" {
		rollback.AutoPromote = toBool(value)
		rollback.AutoPromoteSet = true
	}
	if value := annotations[kubernetesAnnotationDeployRollbackCanary]; value != "" {
		rollback.Canary = parseInt(value)
		rollback.CanarySet = true
	}
	if value := annotations[kubernetesAnnotationDeployRollbackStagger]; value != "" {
		rollback.Stagger = value
	}
	if isEmptyUpdatePolicy(rollback) {
		return nil
	}
	return rollback
}

func isEmptyUpdatePolicy(policy *UpdatePolicy) bool {
	return policy == nil ||
		(!policy.ParallelismSet &&
			policy.Parallelism == 0 &&
			policy.Delay == "" &&
			policy.Monitor == "" &&
			policy.MaxFailureRatio == "" &&
			policy.Order == "" &&
			policy.OnFailure == "" &&
			policy.HealthCheck == "" &&
			policy.MinHealthyTime == "" &&
			policy.HealthyDeadline == "" &&
			policy.ProgressDeadline == "" &&
			!policy.AutoRevert &&
			!policy.AutoRevertSet &&
			!policy.AutoPromote &&
			!policy.AutoPromoteSet &&
			policy.Canary == 0 &&
			!policy.CanarySet &&
			policy.Stagger == "" &&
			len(policy.Extensions) == 0)
}

func parseKubernetesRestartPolicy(annotations map[string]string) *RestartPolicy {
	restart := &RestartPolicy{
		Condition:   annotations[kubernetesAnnotationDeployRestartCondition],
		Delay:       annotations[kubernetesAnnotationDeployRestartDelay],
		MaxAttempts: parseInt(annotations[kubernetesAnnotationDeployRestartAttempts]),
		Window:      annotations[kubernetesAnnotationDeployRestartWindow],
	}
	if isEmptyRestartPolicy(restart) {
		return nil
	}
	return restart
}

func applyKubernetesDeployIntent(deployment map[string]interface{}, service *Service) {
	if service == nil || service.Deploy == nil {
		return
	}
	deploy := service.Deploy
	spec := deployment["spec"].(map[string]interface{})
	template := kubernetesSerializedPodTemplate(deployment)
	templateSpec := template["spec"].(map[string]interface{})

	annotations := kubernetesDeployAnnotations(deploy)
	if len(annotations) > 0 {
		setKubernetesDeploymentAnnotations(deployment, annotations)
	}

	if len(service.NodeSelector) > 0 {
		templateSpec["nodeSelector"] = copyStringMap(service.NodeSelector)
	} else if nodeSelector := nodeSelectorFromPlacement(deploy.Placement); len(nodeSelector) > 0 {
		templateSpec["nodeSelector"] = nodeSelector
	}
	if deployment["kind"] == "Deployment" {
		if strategy := kubernetesStrategyFromUpdatePolicy(deploy.UpdateConfig); len(strategy) > 0 {
			spec["strategy"] = strategy
		}
	}
}

func applyKubernetesRuntimeAnnotations(deployment map[string]interface{}, service *Service) {
	annotations := map[string]string{}
	if len(service.SecurityOpt) > 0 {
		annotations[kubernetesAnnotationRuntimeSecurityOpt] = strings.Join(service.SecurityOpt, "\n")
	}
	if service.Init != nil {
		annotations[kubernetesAnnotationRuntimeInit] = fmt.Sprintf("%t", *service.Init)
	}
	if service.StopSignal != "" {
		annotations[kubernetesAnnotationRuntimeStopSignal] = service.StopSignal
	}
	if len(annotations) > 0 {
		setKubernetesDeploymentAnnotations(deployment, annotations)
	}
}

func applyKubernetesPortableAnnotations(deployment map[string]interface{}, service *Service, app *Application) {
	annotations := map[string]string{}
	if service != nil && len(service.Extensions) > 0 {
		extensions := jsonSerializableExtensionMap(service.Extensions)
		delete(extensions, kubernetesAnnotationPortableServiceExtensions)
		if len(extensions) > 0 {
			if raw, err := json.Marshal(extensions); err == nil {
				annotations[kubernetesAnnotationPortableServiceExtensions] = string(raw)
			}
		}
	}
	if len(app.Extensions) > 0 {
		extensions := jsonSerializableExtensionMap(app.Extensions)
		for _, key := range []string{
			kubernetesSourceDocumentsExtensionKey,
			"kubernetes.namespace",
			"kubernetes.raw",
			"kubernetes.resources",
			"kubernetes.workloads",
			"kubernetes.serviceResources",
			"kubernetes.services",
			"kubernetes.hpas",
			"kubernetes.pdbs",
			"kubernetes.networkPolicies",
			"kubernetes.ingresses",
			"kubernetes.serviceAccounts",
			"kubernetes.rbacResources",
			"kubernetes.resourceQuotas",
			"kubernetes.limitRanges",
			"kubernetes.priorityClasses",
			"kubernetes.runtimeClasses",
			"kubernetes.storageClasses",
			"kubernetes.ingressClasses",
			"kubernetes.mutatingWebhookConfigurations",
			"kubernetes.validatingWebhookConfigurations",
			"kubernetes.customResourceDefinitions",
			"kubernetes.customResources",
			composeKubernetesResourcesExtensionKey,
			composeKubernetesWorkloadsExtensionKey,
			composeKubernetesServicesExtensionKey,
			composeKubernetesHPAsExtensionKey,
			composeKubernetesPDBsExtensionKey,
			composeKubernetesNetworkPoliciesExtensionKey,
			composeKubernetesIngressesExtensionKey,
			composeKubernetesServiceAccountsExtensionKey,
			composeKubernetesRBACResourcesExtensionKey,
			composeKubernetesResourceQuotasExtensionKey,
			composeKubernetesLimitRangesExtensionKey,
			composeKubernetesPriorityClassesExtensionKey,
			composeKubernetesRuntimeClassesExtensionKey,
			composeKubernetesStorageClassesExtensionKey,
			composeKubernetesIngressClassesExtensionKey,
			composeKubernetesMutatingWebhooksExtensionKey,
			composeKubernetesValidatingWebhooksExtensionKey,
			composeKubernetesCRDsExtensionKey,
			composeKubernetesCustomResourcesExtensionKey,
		} {
			delete(extensions, key)
		}
		delete(extensions, kubernetesAnnotationPortableAppExtensions)
		delete(extensions, kubernetesAnnotationPortableServiceExtensions)
		delete(extensions, kubernetesAnnotationPortableComposeName)
		delete(extensions, kubernetesAnnotationPortableModels)
		delete(extensions, kubernetesAnnotationPortableIncludes)
		delete(extensions, kubernetesAnnotationPortableNetworks)
		delete(extensions, kubernetesAnnotationPortableAppVolumes)
		delete(extensions, kubernetesAnnotationPortableAppConfigs)
		delete(extensions, kubernetesAnnotationPortableAppSecrets)
		if len(extensions) > 0 {
			if raw, err := json.Marshal(extensions); err == nil {
				annotations[kubernetesAnnotationPortableAppExtensions] = string(raw)
			}
		}
	}
	if app.Name != "" {
		annotations[kubernetesAnnotationPortableComposeName] = app.Name
	}
	if app.Mesh != nil {
		if raw, err := json.Marshal(meshSpecToMap(app.Mesh)); err == nil {
			annotations[kubernetesAnnotationPortableMesh] = string(raw)
		}
	}
	if models := applicationModelsForEmit(app); len(models) > 0 {
		if raw, err := json.Marshal(models); err == nil {
			annotations[kubernetesAnnotationPortableModels] = string(raw)
		}
	}
	if len(app.IncludeEntries) > 0 {
		if raw, err := json.Marshal(app.IncludeEntries); err == nil {
			annotations[kubernetesAnnotationPortableIncludes] = string(raw)
		}
	}
	if networks := applicationNetworksForEmit(app); len(networks) > 0 {
		if raw, err := json.Marshal(networks); err == nil {
			annotations[kubernetesAnnotationPortableNetworks] = string(raw)
		}
	}
	if volumes := applicationVolumesForEmit(app); len(volumes) > 0 {
		if raw, err := json.Marshal(volumes); err == nil {
			annotations[kubernetesAnnotationPortableAppVolumes] = string(raw)
		}
	}
	if configs := applicationConfigsForEmit(app); len(configs) > 0 {
		if raw, err := json.Marshal(configs); err == nil {
			annotations[kubernetesAnnotationPortableAppConfigs] = string(raw)
		}
	}
	if secrets := applicationSecretsForEmit(app); len(secrets) > 0 {
		if raw, err := json.Marshal(secrets); err == nil {
			annotations[kubernetesAnnotationPortableAppSecrets] = string(raw)
		}
	}
	if buildConfigHasData(service.Build) {
		if raw, err := json.Marshal(service.Build); err == nil {
			annotations[kubernetesAnnotationPortableBuild] = string(raw)
		}
	}
	if len(service.Devices) > 0 {
		if raw, err := json.Marshal(service.Devices); err == nil {
			annotations[kubernetesAnnotationPortableDevices] = string(raw)
		}
	}
	if len(service.DeviceMappings) > 0 {
		if raw, err := json.Marshal(service.DeviceMappings); err == nil {
			annotations[kubernetesAnnotationPortableDeviceMappings] = string(raw)
		}
	}
	if len(service.Expose) > 0 {
		if raw, err := json.Marshal(service.Expose); err == nil {
			annotations[kubernetesAnnotationPortableExpose] = string(raw)
		}
	}
	if portablePortsNeedKubernetesAnnotation(service.Ports) {
		if raw, err := json.Marshal(service.Ports); err == nil {
			annotations[kubernetesAnnotationPortablePorts] = string(raw)
		}
	}
	if len(service.EnvFileRefs) > 0 {
		if raw, err := json.Marshal(service.EnvFileRefs); err == nil {
			annotations[kubernetesAnnotationPortableEnvFiles] = string(raw)
		}
	} else if len(service.EnvFile) > 0 {
		refs := make([]EnvFileRef, 0, len(service.EnvFile))
		for _, path := range service.EnvFile {
			refs = append(refs, EnvFileRef{Path: path})
		}
		if raw, err := json.Marshal(refs); err == nil {
			annotations[kubernetesAnnotationPortableEnvFiles] = string(raw)
		}
	}
	if !isEmptyHealthCheck(service.HealthCheck) {
		if raw, err := json.Marshal(service.HealthCheck); err == nil {
			annotations[kubernetesAnnotationPortableHealthcheck] = string(raw)
		}
	}
	if !isEmptyHealthCheck(service.StartupProbe) {
		if raw, err := json.Marshal(service.StartupProbe); err == nil {
			annotations[kubernetesAnnotationPortableStartupProbe] = string(raw)
		}
	}
	if len(service.Tolerations) > 0 {
		needsAnnotation := false
		for _, toleration := range service.Tolerations {
			if len(toleration.Extensions) > 0 {
				needsAnnotation = true
				break
			}
		}
		if needsAnnotation {
			if raw, err := json.Marshal(service.Tolerations); err == nil {
				annotations[kubernetesAnnotationPortableTolerations] = string(raw)
			}
		}
	}
	if !isEmptyDevelopConfig(service.Develop) {
		if raw, err := json.Marshal(service.Develop); err == nil {
			annotations[kubernetesAnnotationPortableDevelop] = string(raw)
		}
	}
	if !isEmptyLifecycleHooks(service.Lifecycle) {
		if raw, err := json.Marshal(service.Lifecycle); err == nil {
			annotations[kubernetesAnnotationPortableLifecycle] = string(raw)
		}
	}
	if service.CPUShares > 0 {
		annotations[kubernetesAnnotationPortableCPUShares] = strconv.Itoa(service.CPUShares)
	}
	if service.CPUQuota > 0 {
		annotations[kubernetesAnnotationPortableCPUQuota] = strconv.Itoa(service.CPUQuota)
	}
	if service.MemLimit != "" {
		annotations[kubernetesAnnotationPortableMemLimit] = service.MemLimit
	} else if service.MemoryLimit != "" {
		annotations[kubernetesAnnotationPortableMemLimit] = service.MemoryLimit
	}
	if service.MemorySwap != "" {
		annotations[kubernetesAnnotationPortableMemorySwap] = service.MemorySwap
	}
	if service.MemReservation != "" {
		annotations[kubernetesAnnotationPortableMemReservation] = service.MemReservation
	}
	if service.CPUs != "" {
		annotations[kubernetesAnnotationPortableCPUs] = service.CPUs
	}
	if !isEmptyUlimits(service.Ulimits) {
		if raw, err := json.Marshal(service.Ulimits); err == nil {
			annotations[kubernetesAnnotationPortableUlimits] = string(raw)
		}
	}
	if service.UserNSMode != "" {
		annotations[kubernetesAnnotationPortableUserNSMode] = service.UserNSMode
	}
	if service.PullPolicy != "" {
		annotations[kubernetesAnnotationPortablePullPolicy] = service.PullPolicy
	}
	if len(service.Profiles) > 0 {
		if raw, err := json.Marshal(service.Profiles); err == nil {
			annotations[kubernetesAnnotationPortableProfiles] = string(raw)
		}
	}
	if len(service.Configs) > 0 {
		if raw, err := json.Marshal(service.Configs); err == nil {
			annotations[kubernetesAnnotationPortableConfigs] = string(raw)
		}
	}
	if len(service.Secrets) > 0 {
		if raw, err := json.Marshal(service.Secrets); err == nil {
			annotations[kubernetesAnnotationPortableSecrets] = string(raw)
		}
	}
	if len(service.NetworkAttachments) > 0 {
		if raw, err := json.Marshal(service.NetworkAttachments); err == nil {
			annotations[kubernetesAnnotationPortableNetworkAttachments] = string(raw)
		}
	}
	if service.Failover != nil {
		if raw, err := json.Marshal(service.Failover); err == nil {
			annotations[kubernetesAnnotationPortableFailover] = string(raw)
		}
	}
	if restart := nomadRestartBlockForService(service); len(restart) > 0 {
		if raw, err := json.Marshal(restart); err == nil {
			annotations[kubernetesAnnotationPortableNomadRestart] = string(raw)
		}
	}
	if update := nomadSchedulerExtensionMap(service, nomadUpdateExtensionKey, "x-nomad-update"); len(update) > 0 {
		if raw, err := json.Marshal(update); err == nil {
			annotations[kubernetesAnnotationPortableNomadUpdate] = string(raw)
		}
	}
	if migrate := nomadSchedulerExtensionMap(service, nomadMigrateExtensionKey, "x-nomad-migrate"); len(migrate) > 0 {
		if raw, err := json.Marshal(migrate); err == nil {
			annotations[kubernetesAnnotationPortableNomadMigrate] = string(raw)
		}
	}
	if reschedule := nomadSchedulerExtensionMap(service, nomadRescheduleExtensionKey, "x-nomad-reschedule"); len(reschedule) > 0 {
		if raw, err := json.Marshal(reschedule); err == nil {
			annotations[kubernetesAnnotationPortableNomadReschedule] = string(raw)
		}
	}
	if value, ok := extensionValueForKey(service, nomadSpreadExtensionKey); ok {
		if spreads, ok := value.([]map[string]interface{}); ok && len(spreads) > 0 {
			if raw, err := json.Marshal(spreads); err == nil {
				annotations[kubernetesAnnotationPortableNomadSpread] = string(raw)
			}
		} else if spreads, ok := value.([]interface{}); ok && len(spreads) > 0 {
			if raw, err := json.Marshal(spreads); err == nil {
				annotations[kubernetesAnnotationPortableNomadSpread] = string(raw)
			}
		}
	}
	if value, ok := extensionValueForKey(service, nomadConnectExtensionKey); ok {
		if connect, ok := asMap(value); ok && len(connect) > 0 {
			if raw, err := json.Marshal(connect); err == nil {
				annotations[kubernetesAnnotationPortableNomadConnect] = string(raw)
			}
		}
	}
	if len(service.GroupAdd) > 0 {
		if raw, err := json.Marshal(service.GroupAdd); err == nil {
			annotations[kubernetesAnnotationPortableGroupAdd] = string(raw)
		}
	}
	if service.Runtime != "" {
		annotations[kubernetesAnnotationPortableRuntime] = service.Runtime
	}
	if service.LogDriver != "" || len(service.LogOpt) > 0 || len(service.LogExtensions) > 0 {
		if raw, err := json.Marshal(map[string]interface{}{
			"driver":     service.LogDriver,
			"options":    service.LogOpt,
			"extensions": service.LogExtensions,
		}); err == nil {
			annotations[kubernetesAnnotationPortableLogging] = string(raw)
		}
	}
	if !isEmptyComposeCompat(service.ComposeCompat) {
		if raw, err := json.Marshal(service.ComposeCompat); err == nil {
			annotations[kubernetesAnnotationPortableComposeCompat] = string(raw)
		}
	}
	if len(service.Links) > 0 {
		if raw, err := json.Marshal(service.Links); err == nil {
			annotations[kubernetesAnnotationPortableLinks] = string(raw)
		}
	}
	if len(service.Volumes) > 0 {
		if raw, err := json.Marshal(service.Volumes); err == nil {
			annotations[kubernetesAnnotationPortableVolumes] = string(raw)
		}
	}
	if service.PIDMode != "" {
		annotations[kubernetesAnnotationPortablePIDMode] = service.PIDMode
	}
	if service.IPCMode != "" {
		annotations[kubernetesAnnotationPortableIPCMode] = service.IPCMode
	}
	if service.pidsLimitSet || service.PidsLimit > 0 {
		annotations[kubernetesAnnotationPortablePidsLimit] = fmt.Sprintf("%d", service.PidsLimit)
	}
	if service.shmSizeSet || service.ShmSize > 0 {
		annotations[kubernetesAnnotationPortableShmSize] = fmt.Sprintf("%d", service.ShmSize)
	}
	if len(annotations) > 0 {
		setKubernetesDeploymentAnnotations(deployment, annotations)
	}
}

func setKubernetesDeploymentAnnotations(deployment map[string]interface{}, annotations map[string]string) {
	metadata := deployment["metadata"].(map[string]interface{})
	current := toStringMapLoose(metadata["annotations"])
	merged := mergeStringMaps(current, annotations)
	if len(merged) > 0 {
		metadata["annotations"] = merged
	}
}

func kubernetesDeployAnnotations(deploy *DeploySpec) map[string]string {
	annotations := map[string]string{}
	if !isEmptyDeploySpec(deploy) {
		if raw, err := json.Marshal(deploy); err == nil {
			annotations[kubernetesAnnotationDeploySpec] = string(raw)
		}
	}
	if deploy.Mode != "" {
		annotations[kubernetesAnnotationDeployMode] = deploy.Mode
	}
	if deploy.EndpointMode != "" {
		annotations[kubernetesAnnotationDeployEndpointMode] = deploy.EndpointMode
	}
	if len(deploy.Labels) > 0 {
		if raw, err := json.Marshal(deploy.Labels); err == nil {
			annotations[kubernetesAnnotationDeployLabels] = string(raw)
		}
	}
	if hasKubernetesPortableResourceExtras(deploy.Resources) {
		if raw, err := json.Marshal(deploy.Resources); err == nil {
			annotations[kubernetesAnnotationDeployResources] = string(raw)
		}
	}
	if !isEmptySwarmJobSpec(deploy.Job) {
		if raw, err := json.Marshal(deploy.Job); err == nil {
			annotations[kubernetesAnnotationDeployJob] = string(raw)
		}
	}
	if deploy.Placement != nil {
		if len(deploy.Placement.Constraints) > 0 {
			annotations[kubernetesAnnotationDeployPlacement] = strings.Join(deploy.Placement.Constraints, "\n")
		}
		if len(deploy.Placement.Preferences) > 0 {
			annotations[kubernetesAnnotationDeployPreferences] = strings.Join(deploy.Placement.Preferences, "\n")
		}
		if deploy.Placement.MaxReplicasPerNode > 0 {
			annotations[kubernetesAnnotationDeployMaxReplicasPerNode] = fmt.Sprintf("%d", deploy.Placement.MaxReplicasPerNode)
		}
	}
	if deploy.UpdateConfig != nil {
		addKubernetesUpdatePolicyAnnotations(annotations, deploy.UpdateConfig,
			kubernetesAnnotationDeployUpdateParallelism,
			kubernetesAnnotationDeployUpdateDelay,
			kubernetesAnnotationDeployUpdateMonitor,
			kubernetesAnnotationDeployUpdateFailureRate,
			kubernetesAnnotationDeployUpdateOrder,
			kubernetesAnnotationDeployUpdateOnFailure,
			kubernetesAnnotationDeployUpdateHealthCheck,
			kubernetesAnnotationDeployUpdateMinHealthyTime,
			kubernetesAnnotationDeployUpdateHealthyDeadline,
			kubernetesAnnotationDeployUpdateProgressDeadline,
			kubernetesAnnotationDeployUpdateAutoRevert,
			kubernetesAnnotationDeployUpdateAutoPromote,
			kubernetesAnnotationDeployUpdateCanary,
			kubernetesAnnotationDeployUpdateStagger,
		)
	}
	if deploy.RollbackConfig != nil {
		addKubernetesUpdatePolicyAnnotations(annotations, deploy.RollbackConfig,
			kubernetesAnnotationDeployRollbackParallelism,
			kubernetesAnnotationDeployRollbackDelay,
			kubernetesAnnotationDeployRollbackMonitor,
			kubernetesAnnotationDeployRollbackFailureRate,
			kubernetesAnnotationDeployRollbackOrder,
			kubernetesAnnotationDeployRollbackOnFailure,
			kubernetesAnnotationDeployRollbackHealthCheck,
			kubernetesAnnotationDeployRollbackMinHealthyTime,
			kubernetesAnnotationDeployRollbackHealthyDeadline,
			kubernetesAnnotationDeployRollbackProgressDeadline,
			kubernetesAnnotationDeployRollbackAutoRevert,
			kubernetesAnnotationDeployRollbackAutoPromote,
			kubernetesAnnotationDeployRollbackCanary,
			kubernetesAnnotationDeployRollbackStagger,
		)
	}
	if deploy.RestartPolicy != nil {
		if deploy.RestartPolicy.Condition != "" {
			annotations[kubernetesAnnotationDeployRestartCondition] = deploy.RestartPolicy.Condition
		}
		if deploy.RestartPolicy.Delay != "" {
			annotations[kubernetesAnnotationDeployRestartDelay] = deploy.RestartPolicy.Delay
		}
		if deploy.RestartPolicy.MaxAttempts > 0 {
			annotations[kubernetesAnnotationDeployRestartAttempts] = strconv.Itoa(deploy.RestartPolicy.MaxAttempts)
		}
		if deploy.RestartPolicy.Window != "" {
			annotations[kubernetesAnnotationDeployRestartWindow] = deploy.RestartPolicy.Window
		}
	}
	if len(annotations) == 0 {
		return nil
	}
	return annotations
}

func hasKubernetesPortableResourceExtras(resources *ResourceSpec) bool {
	return resources != nil &&
		(resources.pidsLimitSet || resources.PidsLimit > 0 ||
			resources.pidsReservationSet || resources.PidsReservation > 0 ||
			len(resources.Devices) > 0 ||
			len(resources.GenericResources) > 0 ||
			len(resources.Extensions) > 0 ||
			len(resources.LimitExtensions) > 0 ||
			len(resources.ReservationExtensions) > 0)
}

func addKubernetesUpdatePolicyAnnotations(annotations map[string]string, policy *UpdatePolicy, parallelismKey, delayKey, monitorKey, failureRateKey, orderKey, onFailureKey, healthCheckKey, minHealthyTimeKey, healthyDeadlineKey, progressDeadlineKey, autoRevertKey, autoPromoteKey, canaryKey, staggerKey string) {
	if policy == nil {
		return
	}
	if policy.ParallelismSet || policy.Parallelism > 0 {
		annotations[parallelismKey] = strconv.Itoa(policy.Parallelism)
	}
	if policy.Delay != "" {
		annotations[delayKey] = policy.Delay
	}
	if policy.Monitor != "" {
		annotations[monitorKey] = policy.Monitor
	}
	if policy.MaxFailureRatio != "" {
		annotations[failureRateKey] = policy.MaxFailureRatio
	}
	if policy.Order != "" {
		annotations[orderKey] = policy.Order
	}
	if policy.OnFailure != "" {
		annotations[onFailureKey] = policy.OnFailure
	}
	if policy.HealthCheck != "" {
		annotations[healthCheckKey] = policy.HealthCheck
	}
	if policy.MinHealthyTime != "" {
		annotations[minHealthyTimeKey] = policy.MinHealthyTime
	}
	if policy.HealthyDeadline != "" {
		annotations[healthyDeadlineKey] = policy.HealthyDeadline
	}
	if policy.ProgressDeadline != "" {
		annotations[progressDeadlineKey] = policy.ProgressDeadline
	}
	if policy.AutoRevertSet || policy.AutoRevert {
		annotations[autoRevertKey] = fmt.Sprintf("%t", policy.AutoRevert)
	}
	if policy.AutoPromoteSet || policy.AutoPromote {
		annotations[autoPromoteKey] = fmt.Sprintf("%t", policy.AutoPromote)
	}
	if policy.CanarySet || policy.Canary > 0 {
		annotations[canaryKey] = strconv.Itoa(policy.Canary)
	}
	if policy.Stagger != "" {
		annotations[staggerKey] = policy.Stagger
	}
}

func nodeSelectorFromPlacement(placement *PlacementSpec) map[string]string {
	if placement == nil {
		return nil
	}
	selector := map[string]string{}
	for _, constraint := range placement.Constraints {
		key, value, ok := parseNodeLabelConstraint(constraint)
		if ok {
			selector[key] = value
		}
	}
	if len(selector) == 0 {
		return nil
	}
	return selector
}

func parseNodeLabelConstraint(constraint string) (string, string, bool) {
	constraint = strings.TrimSpace(constraint)
	var parts []string
	switch {
	case strings.Contains(constraint, "=="):
		parts = strings.SplitN(constraint, "==", 2)
	case strings.Contains(constraint, "="):
		parts = strings.SplitN(constraint, "=", 2)
	default:
		return "", "", false
	}
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	key = strings.TrimPrefix(key, "node.labels.")
	if key == "" || value == "" || strings.ContainsAny(key, " !<>") {
		return "", "", false
	}
	return key, value, true
}

func kubernetesStrategyFromUpdatePolicy(update *UpdatePolicy) map[string]interface{} {
	if update == nil {
		return nil
	}
	if value, ok := updatePolicyExtensionValue(update, "kubernetes.deployment.strategy", "kubernetes-deployment-strategy", "x-kubernetes-deployment-strategy"); ok {
		if strategy, ok := extensionMapValues(value); ok && len(strategy) > 0 {
			return strategy
		}
	}
	rollingUpdate := map[string]interface{}{}
	switch strings.ToLower(update.Order) {
	case "start-first":
		rollingUpdate["maxUnavailable"] = 0
		rollingUpdate["maxSurge"] = 1
	case "stop-first":
		rollingUpdate["maxSurge"] = 0
		if update.Parallelism > 0 {
			rollingUpdate["maxUnavailable"] = update.Parallelism
		} else if update.ParallelismSet {
			rollingUpdate["maxUnavailable"] = "100%"
		} else {
			rollingUpdate["maxUnavailable"] = 1
		}
	}
	if len(rollingUpdate) == 0 {
		return nil
	}
	return map[string]interface{}{
		"type":          "RollingUpdate",
		"rollingUpdate": rollingUpdate,
	}
}

func updatePolicyExtensionValue(policy *UpdatePolicy, keys ...string) (interface{}, bool) {
	if policy == nil || len(policy.Extensions) == 0 {
		return nil, false
	}
	for _, key := range keys {
		if value, ok := policy.Extensions[key]; ok {
			return value, true
		}
	}
	return nil, false
}

func splitAnnotationList(value string) []string {
	var result []string
	for _, part := range strings.FieldsFunc(value, func(r rune) bool { return r == '\n' || r == ',' }) {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseAnnotationStringSlice(value string) []string {
	var result []string
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(value), &result); err == nil {
		return result
	}
	return splitAnnotationList(value)
}

func appendUniqueString(values *[]string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	for _, existing := range *values {
		if existing == value {
			return
		}
	}
	*values = append(*values, value)
}

func mergeStringMaps(maps ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, values := range maps {
		for key, value := range values {
			result[key] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func mergeLooseMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for _, values := range maps {
		for key, value := range values {
			result[key] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func isEmptyDeploySpec(deploy *DeploySpec) bool {
	if deploy == nil {
		return true
	}
	return deploy.Mode == "" &&
		deploy.EndpointMode == "" &&
		deploy.Replicas == 0 &&
		deploy.Job == nil &&
		deploy.Placement == nil &&
		deploy.Resources == nil &&
		deploy.UpdateConfig == nil &&
		deploy.MigrateConfig == nil &&
		deploy.RescheduleConfig == nil &&
		deploy.RollbackConfig == nil &&
		deploy.RestartPolicy == nil &&
		len(deploy.Labels) == 0 &&
		len(deploy.Extensions) == 0
}

func serializeKubernetesResources(resources *ResourceSpec) map[string]interface{} {
	result := map[string]interface{}{}
	limits := map[string]string{}
	requests := map[string]string{}
	if resources.CPULimit != "" {
		if raw := kubernetesCPUQuantityExtension(resources.Extensions, resources.LimitExtensions, "x-kubernetes-cpu-limit", "kubernetes.cpu.limit", "kubernetes.cpuLimit"); raw != "" {
			limits["cpu"] = raw
		} else {
			limits["cpu"] = resources.CPULimit
		}
	}
	if resources.MemoryLimit != "" {
		if raw := kubernetesMemoryQuantityExtension(resources.Extensions, resources.LimitExtensions, "x-kubernetes-memory-limit", "kubernetes.memory.limit", "kubernetes.memoryLimit"); raw != "" {
			limits["memory"] = raw
		} else {
			limits["memory"] = resources.MemoryLimit
		}
	}
	if resources.EphemeralStorageLimit != "" {
		if raw := kubernetesCPUQuantityExtension(resources.Extensions, resources.LimitExtensions, "x-kubernetes-ephemeral-storage-limit", "kubernetes.ephemeralStorage.limit", "kubernetes.ephemeralStorageLimit"); raw != "" {
			limits["ephemeral-storage"] = raw
		} else {
			limits["ephemeral-storage"] = resources.EphemeralStorageLimit
		}
	}
	if resources.CPUReservation != "" {
		if raw := kubernetesCPUQuantityExtension(resources.Extensions, resources.ReservationExtensions, "x-kubernetes-cpu-reservation", "kubernetes.cpu.reservation", "kubernetes.cpuReservation"); raw != "" {
			requests["cpu"] = raw
		} else {
			requests["cpu"] = resources.CPUReservation
		}
	}
	if resources.MemoryReservation != "" {
		if raw := kubernetesMemoryQuantityExtension(resources.Extensions, resources.ReservationExtensions, "x-kubernetes-memory-reservation", "kubernetes.memory.reservation", "kubernetes.memoryReservation"); raw != "" {
			requests["memory"] = raw
		} else {
			requests["memory"] = resources.MemoryReservation
		}
	}
	if resources.EphemeralStorageReservation != "" {
		if raw := kubernetesCPUQuantityExtension(resources.Extensions, resources.ReservationExtensions, "x-kubernetes-ephemeral-storage-reservation", "kubernetes.ephemeralStorage.reservation", "kubernetes.ephemeralStorageReservation"); raw != "" {
			requests["ephemeral-storage"] = raw
		} else {
			requests["ephemeral-storage"] = resources.EphemeralStorageReservation
		}
	}
	if len(limits) > 0 {
		result["limits"] = limits
	}
	if len(requests) > 0 {
		result["requests"] = requests
	}
	if len(resources.Extensions) > 0 {
		if claims, ok := resources.Extensions["kubernetes.claims"]; ok {
			if items, ok := claims.([]interface{}); ok && len(items) > 0 {
				result["claims"] = cloneInterfaceSlice(items)
			}
		}
		if claims, ok := resources.Extensions["x-kubernetes-claims"]; ok {
			if items, ok := claims.([]interface{}); ok && len(items) > 0 {
				result["claims"] = cloneInterfaceSlice(items)
			}
		}
	}
	return result
}

func normalizeKubernetesCPUQuantity(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	quantity, err := resource.ParseQuantity(value)
	if err != nil {
		return value, false
	}
	normalized := strconv.FormatFloat(quantity.AsApproximateFloat64(), 'f', -1, 64)
	return normalized, strings.HasSuffix(strings.ToLower(value), "m")
}

func normalizeKubernetesMemoryQuantity(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	quantity, err := resource.ParseQuantity(value)
	if err != nil {
		return value, false
	}
	return strconv.FormatInt(quantity.Value(), 10), true
}

func kubernetesCPUQuantityExtension(primary map[string]interface{}, secondary map[string]interface{}, keys ...string) string {
	for _, extensions := range []map[string]interface{}{primary, secondary} {
		for _, key := range keys {
			if value, ok := extensions[key]; ok {
				if text := toString(value); text != "" {
					return text
				}
			}
		}
	}
	return ""
}

func kubernetesMemoryQuantityExtension(primary map[string]interface{}, secondary map[string]interface{}, keys ...string) string {
	return kubernetesCPUQuantityExtension(primary, secondary, keys...)
}

func parseKubernetesProbe(probe map[string]interface{}) *HealthCheck {
	health := &HealthCheck{}
	if len(probe) > 0 {
		health.Extensions = map[string]interface{}{}
	}
	if httpGet, ok := asMap(probe["httpGet"]); ok {
		path := toString(httpGet["path"])
		port := toString(httpGet["port"])
		health.Type = "http"
		health.Path = path
		health.Port = port
		if path != "" || port != "" {
			health.Test = []string{"CMD-SHELL", fmt.Sprintf("wget --spider -q http://127.0.0.1:%s%s || exit 1", port, path)}
		}
	}
	if tcpSocket, ok := asMap(probe["tcpSocket"]); ok {
		health.Type = "tcp"
		health.Port = toString(tcpSocket["port"])
	}
	if execProbe, ok := asMap(probe["exec"]); ok {
		health.Type = "exec"
		if command, err := toStringSlice(execProbe["command"]); err == nil {
			health.Test = append([]string{"CMD"}, command...)
		}
	}
	if grpcProbe, ok := asMap(probe["grpc"]); ok {
		health.Type = "grpc"
		grpc := cloneMap(grpcProbe)
		health.Port = toString(grpcProbe["port"])
		if len(grpc) > 0 {
			health.Extensions["x-kubernetes-grpc"] = grpc
		}
	}
	if interval := toInt(probe["periodSeconds"]); interval > 0 {
		health.Interval = fmt.Sprintf("%ds", interval)
	}
	if timeout := toInt(probe["timeoutSeconds"]); timeout > 0 {
		health.Timeout = fmt.Sprintf("%ds", timeout)
	}
	if retries := toInt(probe["failureThreshold"]); retries > 0 {
		health.Retries = retries
	}
	if startPeriod := toInt(probe["initialDelaySeconds"]); startPeriod > 0 {
		health.StartPeriod = fmt.Sprintf("%ds", startPeriod)
	}
	for key, value := range probe {
		switch key {
		case "httpGet", "tcpSocket", "grpc", "exec", "periodSeconds", "timeoutSeconds", "failureThreshold", "initialDelaySeconds", "successThreshold", "terminationGracePeriodSeconds":
		default:
			health.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(health.Extensions) == 0 {
		health.Extensions = nil
	}
	if health.Test == nil && health.Type == "" && health.Interval == "" && health.Timeout == "" && health.Retries == 0 && health.StartPeriod == "" {
		if len(health.Extensions) == 0 {
			return nil
		}
	}
	return normalizeHealthCheck(health)
}

func parseKubernetesLifecycle(service *Service, container map[string]interface{}) {
	lifecycle, ok := asMap(container["lifecycle"])
	if !ok {
		return
	}
	parsed := &LifecycleHooks{
		Extensions: map[string]interface{}{
			"kubernetes.lifecycle": cloneMap(lifecycle),
		},
	}
	if postStart, ok := asMap(lifecycle["postStart"]); ok {
		if hook := parseKubernetesLifecycleHook(postStart); hook != nil {
			parsed.PostStart = append(parsed.PostStart, *hook)
		}
	}
	if preStop, ok := asMap(lifecycle["preStop"]); ok {
		if hook := parseKubernetesLifecycleHook(preStop); hook != nil {
			parsed.PreStop = append(parsed.PreStop, *hook)
		}
	}
	if isEmptyLifecycleHooks(parsed) {
		return
	}
	if service.Lifecycle == nil {
		service.Lifecycle = parsed
		return
	}
	if len(parsed.Extensions) > 0 {
		if service.Lifecycle.Extensions == nil {
			service.Lifecycle.Extensions = map[string]interface{}{}
		}
		for key, value := range parsed.Extensions {
			if _, exists := service.Lifecycle.Extensions[key]; !exists {
				service.Lifecycle.Extensions[key] = value
			}
		}
	}
	if len(parsed.PreStart) > 0 {
		service.Lifecycle.PreStart = append(service.Lifecycle.PreStart, parsed.PreStart...)
	}
	if len(parsed.PostStart) > 0 {
		service.Lifecycle.PostStart = append(service.Lifecycle.PostStart, parsed.PostStart...)
	}
	if len(parsed.PreStop) > 0 {
		service.Lifecycle.PreStop = append(service.Lifecycle.PreStop, parsed.PreStop...)
	}
}

func parseKubernetesLifecycleHook(raw map[string]interface{}) *ServiceHook {
	if len(raw) == 0 {
		return nil
	}
	hook := &ServiceHook{
		Extensions: map[string]interface{}{
			"kubernetes.lifecycle": cloneMap(raw),
		},
	}
	if execMap, ok := asMap(raw["exec"]); ok {
		if command, err := toStringSlice(execMap["command"]); err == nil && len(command) > 0 {
			hook.Command = command
		}
	}
	return hook
}

func serializeKubernetesLifecycle(lifecycle *LifecycleHooks) map[string]interface{} {
	if isEmptyLifecycleHooks(lifecycle) {
		return nil
	}
	result := map[string]interface{}{}
	if raw, ok := lifecycle.Extensions["kubernetes.lifecycle"].(map[string]interface{}); ok {
		result = cloneMap(raw)
	}
	if len(lifecycle.PostStart) > 0 {
		if hook := serializeKubernetesLifecycleHook(&lifecycle.PostStart[0]); hook != nil {
			result["postStart"] = hook
		}
	}
	if len(lifecycle.PreStop) > 0 {
		if hook := serializeKubernetesLifecycleHook(&lifecycle.PreStop[0]); hook != nil {
			result["preStop"] = hook
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func serializeKubernetesLifecycleHook(hook *ServiceHook) map[string]interface{} {
	if hook == nil {
		return nil
	}
	if raw, ok := hook.Extensions["kubernetes.lifecycle"].(map[string]interface{}); ok {
		cloned := cloneMap(raw)
		if len(hook.Command) > 0 {
			execMap, _ := asMap(cloned["exec"])
			if execMap == nil {
				execMap = map[string]interface{}{}
			}
			execMap["command"] = append([]string{}, hook.Command...)
			cloned["exec"] = execMap
		}
		return cloned
	}
	if len(hook.Command) == 0 {
		return nil
	}
	return map[string]interface{}{
		"exec": map[string]interface{}{
			"command": append([]string{}, hook.Command...),
		},
	}
}

func parseKubernetesSeccompProfile(raw map[string]interface{}) *SeccompProfile {
	if len(raw) == 0 {
		return nil
	}
	profile := &SeccompProfile{
		Type:             toString(raw["type"]),
		LocalhostProfile: toString(raw["localhostProfile"]),
		Extensions:       map[string]interface{}{},
	}
	for key, value := range raw {
		switch key {
		case "type", "localhostProfile":
		default:
			profile.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(profile.Extensions) == 0 {
		profile.Extensions = nil
	}
	if profile.Type == "" && profile.LocalhostProfile == "" {
		if len(profile.Extensions) == 0 {
			return nil
		}
	}
	return profile
}

func serializeKubernetesSeccompProfile(profile *SeccompProfile) map[string]interface{} {
	if profile == nil {
		return nil
	}
	data := map[string]interface{}{}
	if profile.Type != "" {
		data["type"] = profile.Type
	}
	if profile.LocalhostProfile != "" {
		data["localhostProfile"] = profile.LocalhostProfile
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

func serializeKubernetesSeccompProfileExtension(profile *SeccompProfile) map[string]interface{} {
	if profile == nil {
		return nil
	}
	data := serializeKubernetesSeccompProfile(profile)
	if len(profile.Extensions) == 0 {
		return data
	}
	if data == nil {
		data = map[string]interface{}{}
	}
	for key, value := range profile.Extensions {
		if _, exists := data[key]; exists {
			continue
		}
		data[key] = deepCopyValue(value)
	}
	return data
}

func kubernetesTopologySpreadConstraintsFromExtension(value interface{}) ([]map[string]interface{}, bool) {
	items, ok := extensionSliceValues(value)
	if !ok {
		return nil, false
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if itemMap, ok := asMap(item); ok {
			result = append(result, cloneMap(itemMap))
		}
	}
	if len(result) == 0 {
		return nil, false
	}
	return result, true
}

func kubernetesReadinessGatesFromExtension(value interface{}) ([]map[string]interface{}, bool) {
	items, ok := extensionSliceValues(value)
	if !ok {
		return nil, false
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if itemMap, ok := asMap(item); ok {
			result = append(result, cloneMap(itemMap))
		}
	}
	if len(result) == 0 {
		return nil, false
	}
	return result, true
}

func kubernetesMapSliceFromExtension(value interface{}) ([]map[string]interface{}, bool) {
	items, ok := extensionSliceValues(value)
	if !ok {
		return nil, false
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if itemMap, ok := asMap(item); ok {
			result = append(result, cloneMap(itemMap))
		}
	}
	if len(result) == 0 {
		return nil, false
	}
	return result, true
}

func serializeKubernetesProbe(health *HealthCheck) map[string]interface{} {
	health = normalizeHealthCheck(health)
	if health == nil {
		return nil
	}
	if health.Disable {
		return nil
	}
	probe := map[string]interface{}{}
	switch health.Type {
	case "http":
		if health.Port == "" && health.Path == "" {
			return nil
		}
		httpGet := map[string]interface{}{}
		if health.Path != "" {
			httpGet["path"] = health.Path
		} else {
			httpGet["path"] = "/"
		}
		if health.Port != "" {
			httpGet["port"] = kubernetesIntOrString(health.Port)
		}
		probe["httpGet"] = httpGet
	case "tcp":
		if health.Port == "" {
			return nil
		}
		probe["tcpSocket"] = map[string]interface{}{"port": kubernetesIntOrString(health.Port)}
	case "grpc":
		grpc := map[string]interface{}{}
		if value, ok := health.Extensions["x-kubernetes-grpc"]; ok {
			if mapped, ok := asMap(value); ok {
				grpc = cloneMap(mapped)
			}
		}
		if health.Port != "" {
			grpc["port"] = kubernetesIntOrString(health.Port)
		}
		if len(grpc) == 0 {
			return nil
		}
		probe["grpc"] = grpc
	default:
		command := healthCheckCommand(health)
		if len(command) == 0 {
			return nil
		}
		probe["exec"] = map[string]interface{}{"command": command}
	}
	if seconds := durationSeconds(health.Interval); seconds > 0 {
		probe["periodSeconds"] = seconds
	}
	if seconds := durationSeconds(health.Timeout); seconds > 0 {
		probe["timeoutSeconds"] = seconds
	}
	if health.Retries > 0 {
		probe["failureThreshold"] = health.Retries
	}
	if seconds := durationSeconds(health.StartPeriod); seconds > 0 {
		probe["initialDelaySeconds"] = seconds
	}
	return probe
}

func serializeKubernetesProbeExtension(health *HealthCheck) map[string]interface{} {
	if health == nil {
		return nil
	}
	probe := serializeKubernetesProbe(health)
	if len(health.Extensions) == 0 {
		return probe
	}
	if probe == nil {
		probe = map[string]interface{}{}
	}
	for key, value := range health.Extensions {
		if _, exists := probe[key]; exists {
			continue
		}
		probe[key] = deepCopyValue(value)
	}
	return probe
}

func parseKubernetesSELinuxOptions(raw map[string]interface{}) *SELinuxOptions {
	if len(raw) == 0 {
		return nil
	}
	options := &SELinuxOptions{
		User:       toString(raw["user"]),
		Role:       toString(raw["role"]),
		Type:       toString(raw["type"]),
		Level:      toString(raw["level"]),
		Extensions: map[string]interface{}{},
	}
	for key, value := range raw {
		switch key {
		case "user", "role", "type", "level":
		default:
			options.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(options.Extensions) == 0 {
		options.Extensions = nil
	}
	if options.User == "" && options.Role == "" && options.Type == "" && options.Level == "" && len(options.Extensions) == 0 {
		return nil
	}
	return options
}

func serializeKubernetesSELinuxOptions(options *SELinuxOptions) map[string]interface{} {
	if options == nil {
		return nil
	}
	data := map[string]interface{}{}
	if options.User != "" {
		data["user"] = options.User
	}
	if options.Role != "" {
		data["role"] = options.Role
	}
	if options.Type != "" {
		data["type"] = options.Type
	}
	if options.Level != "" {
		data["level"] = options.Level
	}
	if len(options.Extensions) > 0 {
		for key, value := range options.Extensions {
			if _, exists := data[key]; exists {
				continue
			}
			data[key] = deepCopyValue(value)
		}
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

func mergeWindowsSecurityContextOptions(native, portable *WindowsSecurityContextOptions) *WindowsSecurityContextOptions {
	if native == nil {
		return cloneWindowsSecurityContextOptions(portable)
	}
	if portable == nil {
		return cloneWindowsSecurityContextOptions(native)
	}
	merged := cloneWindowsSecurityContextOptions(native)
	if merged.GMSACredentialSpecName == nil && portable.GMSACredentialSpecName != nil {
		value := *portable.GMSACredentialSpecName
		merged.GMSACredentialSpecName = &value
	}
	if merged.GMSACredentialSpec == nil && portable.GMSACredentialSpec != nil {
		value := *portable.GMSACredentialSpec
		merged.GMSACredentialSpec = &value
	}
	if merged.RunAsUserName == nil && portable.RunAsUserName != nil {
		value := *portable.RunAsUserName
		merged.RunAsUserName = &value
	}
	if merged.HostProcess == nil && portable.HostProcess != nil {
		value := *portable.HostProcess
		merged.HostProcess = &value
	}
	if len(portable.Extensions) > 0 {
		if merged.Extensions == nil {
			merged.Extensions = map[string]interface{}{}
		}
		for key, value := range portable.Extensions {
			if _, exists := merged.Extensions[key]; !exists {
				merged.Extensions[key] = deepCopyValue(value)
			}
		}
	}
	return merged
}

func parseKubernetesWindowsSecurityContextOptions(raw map[string]interface{}) *WindowsSecurityContextOptions {
	if len(raw) == 0 {
		return nil
	}
	options := &WindowsSecurityContextOptions{Extensions: map[string]interface{}{}}
	if value := toString(raw["gmsaCredentialSpecName"]); value != "" {
		options.GMSACredentialSpecName = &value
	}
	if value := toString(raw["gmsaCredentialSpec"]); value != "" {
		options.GMSACredentialSpec = &value
	}
	if value := toString(raw["runAsUserName"]); value != "" {
		options.RunAsUserName = &value
	}
	if value, ok := raw["hostProcess"].(bool); ok {
		options.HostProcess = boolPtr(value)
	}
	for key, value := range raw {
		switch key {
		case "gmsaCredentialSpecName", "gmsaCredentialSpec", "runAsUserName", "hostProcess":
			continue
		default:
			options.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(options.Extensions) == 0 {
		options.Extensions = nil
	}
	if options.GMSACredentialSpecName == nil && options.GMSACredentialSpec == nil && options.RunAsUserName == nil && options.HostProcess == nil && len(options.Extensions) == 0 {
		return nil
	}
	return options
}

func serializeKubernetesWindowsSecurityContextOptions(options *WindowsSecurityContextOptions) map[string]interface{} {
	if options == nil {
		return nil
	}
	data := map[string]interface{}{}
	if options.GMSACredentialSpecName != nil && *options.GMSACredentialSpecName != "" {
		data["gmsaCredentialSpecName"] = *options.GMSACredentialSpecName
	}
	if options.GMSACredentialSpec != nil && *options.GMSACredentialSpec != "" {
		data["gmsaCredentialSpec"] = *options.GMSACredentialSpec
	}
	if options.RunAsUserName != nil && *options.RunAsUserName != "" {
		data["runAsUserName"] = *options.RunAsUserName
	}
	if options.HostProcess != nil {
		data["hostProcess"] = *options.HostProcess
	}
	if len(options.Extensions) > 0 {
		for key, value := range options.Extensions {
			if _, exists := data[key]; exists {
				continue
			}
			data[key] = deepCopyValue(value)
		}
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

func applyWindowsHostProcessDefaults(service *Service) {
	if service == nil || service.WindowsOptions == nil || service.WindowsOptions.HostProcess == nil || !*service.WindowsOptions.HostProcess {
		return
	}
	if !service.HostNetwork {
		service.HostNetwork = true
		service.HostNetworkSet = true
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.hostNetwork"] = "true"
	}
}

func cloneMap(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(input))
	for key, value := range input {
		if nested, ok := value.(map[string]interface{}); ok {
			output[key] = cloneMap(nested)
			continue
		}
		output[key] = value
	}
	return output
}

func toStringMapLoose(value interface{}) map[string]string {
	result := map[string]string{}
	if typed, ok := value.(map[string]string); ok {
		for key, val := range typed {
			result[key] = val
		}
		return result
	}
	if mapped, ok := asMap(value); ok {
		for key, val := range mapped {
			result[key] = toString(val)
		}
	}
	return result
}
