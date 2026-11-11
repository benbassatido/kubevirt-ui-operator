package builders

import (
	consolev1 "github.com/openshift/api/console/v1"
	corev1 "k8s.io/api/core/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
)

// PriorityClass returns the operator-owned PriorityClass for the UI stack.
func PriorityClass(m *Context) *schedulingv1.PriorityClass {
	policy := corev1.PreemptLowerPriority
	return &schedulingv1.PriorityClass{
		Name:             names.PriorityClass,
		Labels:           labels.Base(m.CR),
		Value:            names.PriorityClassValue,
		GlobalDefault:    false,
		Description:      names.PriorityClassDescription,
		PreemptionPolicy: &policy,
	}
}

func ConsolePlugin(m *Context) *consolev1.ConsolePlugin {
	return &consolev1.ConsolePlugin{
		Name:   names.ConsolePluginName,
		Labels: labels.Plugin(m.CR),
		Spec: consolev1.ConsolePluginSpec{
			DisplayName: names.ConsolePluginDisplayName,
			I18n:        consolev1.ConsolePluginI18n{LoadType: consolev1.Lazy},
			Backend: consolev1.ConsolePluginBackend{
				Type: consolev1.Service,
				Service: &consolev1.ConsolePluginService{
					Name:      names.PluginService,
					Namespace: m.OperandNamespace,
					Port:      names.PluginPort,
					BasePath:  "/",
				},
			},
			Proxy: []consolev1.ConsolePluginProxy{{
				Alias:         names.ConsolePluginProxyAlias,
				Authorization: consolev1.UserToken,
				Endpoint: consolev1.ConsolePluginProxyEndpoint{
					Type: consolev1.ProxyTypeService,
					Service: &consolev1.ConsolePluginProxyServiceConfig{
						Name:      names.ProxyService,
						Namespace: m.OperandNamespace,
						Port:      names.ProxyPort,
					},
				},
			}},
		},
	}
}
