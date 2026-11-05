package labels

import (
	"maps"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
)

const chartName = "kubevirt-console-ui"

const (
	appName   = "kubevirt-console-ui"
	managedBy = "kubevirt-ui-operator"
	appPartOf = "kubevirt-ui"
)

// Base returns labels shared by all operands reconciled by kubevirt-ui-operator.
func Base(cr *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	return map[string]string{
		"helm.sh/chart":                chartName,
		"app.kubernetes.io/instance":   cr.Name,
		"app.kubernetes.io/managed-by": managedBy,
		"app.kubernetes.io/part-of":    appPartOf,
		"app":                          appName,
	}
}

func merge(base map[string]string, extra map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(extra))
	maps.Copy(out, base)
	maps.Copy(out, extra)
	return out
}

// PluginSelector returns Deployment/Service selector labels for the console plugin.
func PluginSelector(_ *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	return map[string]string{
		"app.kubernetes.io/component": "kubevirt-console-plugin",
		"app":                         appName,
	}
}

// ProxySelector returns Deployment/Service selector labels for the apiserver proxy.
func ProxySelector(_ *consoleuiv1alpha1.KubeVirtConsoleUI) map[string]string {
	return map[string]string{
		"app.kubernetes.io/component": "kubevirt-apiserver-proxy",
		"app":                         appName,
	}
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
