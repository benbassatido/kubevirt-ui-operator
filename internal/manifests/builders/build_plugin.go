package builders

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
)

func PluginServiceAccount(m *Context) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		Name:      names.PluginDeployment,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
	}
}

func PluginDeployment(m *Context) *appsv1.Deployment {
	replicas := int32(1)
	return &appsv1.Deployment{
		Name:      names.PluginDeployment,
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
					ServiceAccountName: names.PluginDeployment,
					PriorityClassName:  names.PriorityClass,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: new(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{{
						Name:            names.PluginDeployment,
						Image:           m.PluginImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports:           []corev1.ContainerPort{{Name: "https", ContainerPort: names.PluginPort, Protocol: corev1.ProtocolTCP}},
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
							{Name: names.PluginServingCert, MountPath: "/var/serving-cert", ReadOnly: true},
							{Name: names.NginxConfigMap, MountPath: "/etc/nginx/nginx.conf", SubPath: "nginx.conf", ReadOnly: true},
						},
					}},
					Volumes: []corev1.Volume{
						{Name: names.PluginServingCert,
							Secret: &corev1.SecretVolumeSource{SecretName: names.PluginServingCert, DefaultMode: new(int32(420))}},
						{Name: names.NginxConfigMap,
							ConfigMap: &corev1.ConfigMapVolumeSource{
								Name: names.NginxConfigMap,
							}},
					},
				},
			},
		},
	}
}

func PluginService(m *Context) *corev1.Service {
	return &corev1.Service{
		Name:      names.PluginService,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		Annotations: map[string]string{
			"service.beta.openshift.io/serving-cert-secret-name": names.PluginServingCert,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels.PluginSelector(m.CR),
			Ports: []corev1.ServicePort{{
				Name:       names.PluginDeployment + "-port",
				Port:       names.PluginPort,
				TargetPort: intstrFromInt32(names.PluginPort),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
}

func PluginNetworkPolicy(m *Context) *networkingv1.NetworkPolicy {
	port := intstrFromInt32(names.PluginPort)
	return &networkingv1.NetworkPolicy{
		Name:      names.PluginNetworkPolicy,
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
