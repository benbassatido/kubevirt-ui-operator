package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTypeReady                   = "Ready"
	ConditionTypeProgressing             = "Progressing"
	ConditionTypeDegraded                = "Degraded"
	ConditionTypePlatformRequirementsMet = "PlatformRequirementsMet"
)

// KubeVirtConsoleUISpec defines the desired console UI stack (Helm parity). Operands are deployed in the CR's metadata.namespace.
type KubeVirtConsoleUISpec struct {
	HCOLabels HCOLabelsSpec `json:"hcoLabels,omitempty"`

	PriorityClass PriorityClassSpec `json:"priorityClass,omitempty"`

	Plugin PluginSpec `json:"plugin,omitempty"`

	Proxy ProxySpec `json:"proxy,omitempty"`

	ConsolePlugin ConsolePluginSpec `json:"consolePlugin,omitempty"`

	Config ConfigSpec `json:"config,omitempty"`

	Console ConsoleSpec `json:"console,omitempty"`

	NetworkPolicy NetworkPolicySpec `json:"networkPolicy,omitempty"`

	Route RouteSpec `json:"route,omitempty"`

	RBAC RBACSpec `json:"rbac,omitempty"`
}

type HCOLabelsSpec struct {
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
	// +kubebuilder:default=kubevirt-hyperconverged
	App string `json:"app,omitempty"`
	// +kubebuilder:default=hco-operator
	ManagedBy string `json:"managedBy,omitempty"`
	// +kubebuilder:default=hyperconverged-cluster
	PartOf  string `json:"partOf,omitempty"`
	Version string `json:"version,omitempty"`
}

type PriorityClassSpec struct {
	// +kubebuilder:default=true
	Create *bool `json:"create,omitempty"`
	// +kubebuilder:default=kubevirt-cluster-critical
	Name string `json:"name,omitempty"`
	// +kubebuilder:default=1000000000
	Value       *int32 `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	// +kubebuilder:default=PreemptLowerPriority
	PreemptionPolicy string `json:"preemptionPolicy,omitempty"`
}

type PluginSpec struct {
	// +kubebuilder:default=kubevirt-console-plugin
	Name string `json:"name,omitempty"`
	// +kubebuilder:default="quay.io/kubevirtui/kubevirt-plugin:latest"
	Image string `json:"image,omitempty"`
	// +kubebuilder:default=9443
	Port int32 `json:"port,omitempty"`
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`
	// +kubebuilder:default=kubevirt-console-plugin-service
	ServiceName string `json:"serviceName,omitempty"`
	// +kubebuilder:default=plugin-serving-cert
	ServingCertSecret string `json:"servingCertSecret,omitempty"`
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

type ProxySpec struct {
	// +kubebuilder:default=kubevirt-apiserver-proxy
	Name string `json:"name,omitempty"`
	// +kubebuilder:default="quay.io/kubevirtui/kubevirt-apiserver-proxy:latest"
	Image string `json:"image,omitempty"`
	// +kubebuilder:default=8080
	Port int32 `json:"port,omitempty"`
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`
	// +kubebuilder:default=kubevirt-apiserver-proxy-service
	ServiceName string `json:"serviceName,omitempty"`
	// +kubebuilder:default=console-proxy-serving-cert
	ServingCertSecret string   `json:"servingCertSecret,omitempty"`
	PriorityClassName string   `json:"priorityClassName,omitempty"`
	TLSArgs           []string `json:"tlsArgs,omitempty"`
}

type ConsolePluginSpec struct {
	// +kubebuilder:default=kubevirt-plugin
	Name string `json:"name,omitempty"`
	// +kubebuilder:default=Kubevirt Console Plugin
	DisplayName string `json:"displayName,omitempty"`
	// +kubebuilder:default=kubevirt-apiserver-proxy
	ProxyAlias string `json:"proxyAlias,omitempty"`
}

type ConfigSpec struct {
	// +kubebuilder:default=nginx-conf
	NginxConfigMapName string `json:"nginxConfigMapName,omitempty"`
	// +kubebuilder:default=kubevirt-user-settings
	UserSettingsConfigMapName string `json:"userSettingsConfigMapName,omitempty"`
	// +kubebuilder:default=kubevirt-ui-features
	FeaturesConfigMapName string       `json:"featuresConfigMapName,omitempty"`
	Features              FeaturesSpec `json:"features,omitempty"`
}

type FeaturesSpec struct {
	AutomaticSubscriptionActivationKey  string `json:"automaticSubscriptionActivationKey,omitempty"`
	AutomaticSubscriptionOrganizationId string `json:"automaticSubscriptionOrganizationId,omitempty"`
	// +kubebuilder:default="false"
	DisabledGuestSystemLogsAccess string `json:"disabledGuestSystemLogsAccess,omitempty"`
	// +kubebuilder:default="true"
	KubevirtApiserverProxy string `json:"kubevirtApiserverProxy,omitempty"`
	// +kubebuilder:default="true"
	LoadBalancerEnabled string `json:"loadBalancerEnabled,omitempty"`
	NodePortAddress     string `json:"nodePortAddress,omitempty"`
	// +kubebuilder:default="false"
	NodePortEnabled string `json:"nodePortEnabled,omitempty"`
	IPStackType     string `json:"ipStackType,omitempty"`
	HCOVersion      string `json:"hcoVersion,omitempty"`
}

type ConsoleSpec struct {
	// +kubebuilder:default=true
	RegisterPlugin *bool `json:"registerPlugin,omitempty"`
}

type NetworkPolicySpec struct {
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
}

type RouteSpec struct {
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
	// +kubebuilder:default=kubevirt-apiserver-proxy
	Name string `json:"name,omitempty"`
	Host string `json:"host,omitempty"`
}

type RBACSpec struct {
	// +kubebuilder:default=kubevirt-ui-config-reader
	ConfigReaderRoleName string `json:"configReaderRoleName,omitempty"`
	// +kubebuilder:default=kubevirt-ui-config-reader-rolebinding
	ConfigReaderRoleBindingName string `json:"configReaderRoleBindingName,omitempty"`
}

// KubeVirtConsoleUIStatus defines observed state.
type KubeVirtConsoleUIStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +optional
	ConsolePluginRegistered bool `json:"consolePluginRegistered,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=kvui
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Namespace",type=string,JSONPath=`.metadata.namespace`

// KubeVirtConsoleUI deploys the KubeVirt OpenShift console plugin stack.
type KubeVirtConsoleUI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeVirtConsoleUISpec   `json:"spec,omitempty"`
	Status KubeVirtConsoleUIStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type KubeVirtConsoleUIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeVirtConsoleUI `json:"items"`
}
