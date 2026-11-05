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

// KubeVirtConsoleUISpec defines optional day-2 settings for the console UI stack.
// Operand names, images, and namespace are fixed by the operator (see OPERAND_NAMESPACE and image env vars).
// Feature toggles and hcoVersion live in the kubevirt-ui-features ConfigMap (operator-managed).
type KubeVirtConsoleUISpec struct {
	// Proxy holds optional proxy container flags (for example TLS profile args on OpenShift).
	Proxy ProxySpec `json:"proxy,omitempty"`
}

type ProxySpec struct {
	TLSArgs []string `json:"tlsArgs,omitempty"`
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
// +kubebuilder:resource:scope=Cluster,shortName=kvui
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`

// KubeVirtConsoleUI deploys the KubeVirt OpenShift console plugin stack.
// The singleton instance must be named "cluster"; operand namespace comes from OPERAND_NAMESPACE.
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
