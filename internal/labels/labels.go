package labels

import (
	"maps"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"k8s.io/utils/ptr"
)

const chartName = "kubevirt-console-ui"

// Base returns labels shared by all operands (Helm baseLabels parity).
func Base(cr *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	out := map[string]string{
		"helm.sh/chart":                chartName,
		"app.kubernetes.io/instance":   cr.Name,
		"app.kubernetes.io/managed-by": "kubevirt-ui-operator",
	}
	if ptr.Deref(cr.Spec.HCOLabels.Enabled, true) {
		out["app"] = cr.Spec.HCOLabels.App
		out["app.kubernetes.io/managed-by"] = cr.Spec.HCOLabels.ManagedBy
		out["app.kubernetes.io/part-of"] = cr.Spec.HCOLabels.PartOf
		if cr.Spec.HCOLabels.Version != "" {
			out["app.kubernetes.io/version"] = cr.Spec.HCOLabels.Version
		}
	}
	return out
}

func merge(base map[string]string, extra map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(extra))
	maps.Copy(out, base)
	maps.Copy(out, extra)
	return out
}

// PluginSelector returns Deployment/Service selector labels for the console plugin.
func PluginSelector(cr *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	out := map[string]string{
		"app.kubernetes.io/component": "kubevirt-console-plugin",
	}
	if ptr.Deref(cr.Spec.HCOLabels.Enabled, true) {
		out["app"] = cr.Spec.HCOLabels.App
	}
	return out
}

// ProxySelector returns Deployment/Service selector labels for the apiserver proxy.
func ProxySelector(cr *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	out := map[string]string{
		"app.kubernetes.io/component": "kubevirt-apiserver-proxy",
	}
	if ptr.Deref(cr.Spec.HCOLabels.Enabled, true) {
		out["app"] = cr.Spec.HCOLabels.App
	}
	return out
}

// Plugin returns metadata labels for plugin operands.
func Plugin(cr *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	return merge(Base(cr), map[string]string{"app.kubernetes.io/component": "kubevirt-console-plugin"})
}

// Proxy returns metadata labels for proxy operands.
func Proxy(cr *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	return merge(Base(cr), map[string]string{"app.kubernetes.io/component": "kubevirt-apiserver-proxy"})
}

// UIConfig returns metadata labels for UI ConfigMaps.
func UIConfig(cr *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	return merge(Base(cr), map[string]string{"app.kubernetes.io/component": "kubevirt-ui-config"})
}
