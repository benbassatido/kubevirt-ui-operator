package builders

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
)

func PluginServiceAccount(m *Context) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		Name:      m.Spec.Plugin.Name,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
	}
}

func PluginDeployment(m *Context) *appsv1.Deployment {
	replicas := ptr.Deref(m.Spec.Plugin.Replicas, 1)
	plugin := m.Spec.Plugin
	return &appsv1.Deployment{
		Name:      plugin.Name,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels.PluginSelector(m.CR)},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
					MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
				},
			},
			Template: corev1.PodTemplateSpec{
				Labels: labels.Plugin(m.CR),
				Annotations: map[string]string{
					"openshift.io/required-scc": "restricted-v2",
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: plugin.Name,
					PriorityClassName:  plugin.PriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: new(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{{
						Name:            plugin.Name,
						Image:           plugin.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports:           []corev1.ContainerPort{{Name: "https", ContainerPort: plugin.Port, Protocol: corev1.ProtocolTCP}},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resourceQuantity("10m"),
								corev1.ResourceMemory: resourceQuantity("100Mi"),
							},
						},
						SecurityContext: &corev1.SecurityContext{
							AllowPrivilegeEscalation: new(false),
							Capabilities:             &corev1.Capabilities{Drop: []corev1.Capability{"ALL"}},
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: plugin.ServingCertSecret, MountPath: "/var/serving-cert", ReadOnly: true},
							{Name: m.Spec.Config.NginxConfigMapName, MountPath: "/etc/nginx/nginx.conf", SubPath: "nginx.conf", ReadOnly: true},
						},
					}},
					Volumes: []corev1.Volume{
						{Name: plugin.ServingCertSecret,
							Secret: &corev1.SecretVolumeSource{SecretName: plugin.ServingCertSecret, DefaultMode: new(int32(420))}},
						{Name: m.Spec.Config.NginxConfigMapName,
							ConfigMap: &corev1.ConfigMapVolumeSource{
								Name: m.Spec.Config.NginxConfigMapName,
							}},
					},
				},
			},
		},
	}
}

func PluginService(m *Context) *corev1.Service {
	plugin := m.Spec.Plugin
	return &corev1.Service{
		Name:      plugin.ServiceName,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		Annotations: map[string]string{
			"service.beta.openshift.io/serving-cert-secret-name": plugin.ServingCertSecret,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels.PluginSelector(m.CR),
			Ports: []corev1.ServicePort{{
				Name:       plugin.Name + "-port",
				Port:       plugin.Port,
				TargetPort: intstrFromInt32(plugin.Port),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
}

func PluginNetworkPolicy(m *Context) *networkingv1.NetworkPolicy {
	port := intstrFromInt32(m.Spec.Plugin.Port)
	return &networkingv1.NetworkPolicy{
		Name:      "kubevirt-console-plugin-np",
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{MatchLabels: labels.PluginSelector(m.CR)},
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
			Ingress: []networkingv1.NetworkPolicyIngressRule{{
				From: []networkingv1.NetworkPolicyPeer{{
					NamespaceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"kubernetes.io/metadata.name": "openshift-console"},
					},
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": "console", "component": "ui"},
					},
				}},
				Ports: []networkingv1.NetworkPolicyPort{{Protocol: tcpProtocol(), Port: &port}},
			}},
		},
	}
}
