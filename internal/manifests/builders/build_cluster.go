package builders

import (
	consolev1 "github.com/openshift/api/console/v1"
	corev1 "k8s.io/api/core/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
)

// PriorityClass returns the cluster PriorityClass when create is enabled.
func PriorityClass(m *Context) *schedulingv1.PriorityClass {
	if !consoleuiv1alpha1.BoolOrDefault(m.Spec.PriorityClass.Create, true) {
		return nil
	}
	pc := m.Spec.PriorityClass
	policy := corev1.PreemptionPolicy(pc.PreemptionPolicy)
	return &schedulingv1.PriorityClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pc.Name,
			Labels: labels.Base(m.CR),
		},
		Value:            *pc.Value,
		GlobalDefault:    false,
		Description:      pc.Description,
		PreemptionPolicy: &policy,
	}
}

func ConsolePlugin(m *Context) *consolev1.ConsolePlugin {
	cp := m.Spec.ConsolePlugin
	port := m.Spec.Plugin.Port
	proxyPort := m.Spec.Proxy.Port
	return &consolev1.ConsolePlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cp.Name,
			Labels: labels.Plugin(m.CR),
		},
		Spec: consolev1.ConsolePluginSpec{
			DisplayName: cp.DisplayName,
			I18n:        consolev1.ConsolePluginI18n{LoadType: consolev1.Lazy},
			Backend: consolev1.ConsolePluginBackend{
				Type: consolev1.Service,
				Service: &consolev1.ConsolePluginService{
					Name:      m.Spec.Plugin.ServiceName,
					Namespace: m.OperandNamespace,
					Port:      port,
					BasePath:  "/",
				},
			},
			Proxy: []consolev1.ConsolePluginProxy{{
				Alias:         cp.ProxyAlias,
				Authorization: consolev1.UserToken,
				Endpoint: consolev1.ConsolePluginProxyEndpoint{
					Type: consolev1.ProxyTypeService,
					Service: &consolev1.ConsolePluginProxyServiceConfig{
						Name:      m.Spec.Proxy.ServiceName,
						Namespace: m.OperandNamespace,
						Port:      proxyPort,
					},
				},
			}},
		},
	}
}
