package builders

import (
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
)

func ProxyServiceAccount(m *Context) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		Name:      names.ProxyDeployment,
		Namespace: m.OperandNamespace,
		Labels:    labels.Proxy(m.CR),
	}
}

func ProxyDeployment(m *Context) *appsv1.Deployment {
	replicas := int32(1)
	tlsArgs := m.Spec.Proxy.TLSArgs
	container := corev1.Container{
		Name:            names.ProxyDeployment,
		Image:           m.ProxyImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports:           []corev1.ContainerPort{{Name: "https", ContainerPort: names.ProxyPort, Protocol: corev1.ProtocolTCP}},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resourceQuantity("10m"),
				corev1.ResourceMemory: resourceQuantity("100Mi"),
			},
		},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: new(false),
			Capabilities:             &corev1.Capabilities{Drop: []corev1.Capability{"ALL"}},
			ReadOnlyRootFilesystem:   new(true),
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: names.ProxyServingCert, MountPath: "/app/cert", ReadOnly: true},
		},
	}
	if len(tlsArgs) > 0 {
		container.Args = append([]string{}, tlsArgs...)
	}
	return &appsv1.Deployment{
		Name:      names.ProxyDeployment,
		Namespace: m.OperandNamespace,
		Labels:    labels.Proxy(m.CR),
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels.ProxySelector(m.CR)},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
					MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
				},
			},
			Template: corev1.PodTemplateSpec{
				Labels: labels.Proxy(m.CR),
				Annotations: map[string]string{
					"openshift.io/required-scc": "restricted-v2",
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: names.ProxyDeployment,
					PriorityClassName:  names.PriorityClass,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: new(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{container},
					Volumes: []corev1.Volume{{
						Name:   names.ProxyServingCert,
						Secret: &corev1.SecretVolumeSource{SecretName: names.ProxyServingCert, DefaultMode: new(int32(420))},
					}},
				},
			},
		},
	}
}

func ProxyService(m *Context) *corev1.Service {
	return &corev1.Service{
		Name:      names.ProxyService,
		Namespace: m.OperandNamespace,
		Labels:    labels.Proxy(m.CR),
		Annotations: map[string]string{
			"service.beta.openshift.io/serving-cert-secret-name": names.ProxyServingCert,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels.ProxySelector(m.CR),
			Ports: []corev1.ServicePort{{
				Name:       names.ProxyDeployment + "-port",
				Port:       names.ProxyPort,
				TargetPort: intstrFromInt32(names.ProxyPort),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
}

func ProxyNetworkPolicy(m *Context) *networkingv1.NetworkPolicy {
	proxyPort := intstrFromInt32(names.ProxyPort)
	udp := corev1.ProtocolUDP
	dns := intstr.FromInt32(53)
	https := intstr.FromInt32(443)
	kubeAPI := intstr.FromInt32(6443)
	npLabels := labels.Proxy(m.CR)
	npLabels["np.kubevirt.io/allow-access-cluster-services"] = "true"
	return &networkingv1.NetworkPolicy{
		Name:      names.ProxyNetworkPolicy,
		Namespace: m.OperandNamespace,
		Labels:    npLabels,
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{MatchLabels: labels.ProxySelector(m.CR)},
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress, networkingv1.PolicyTypeEgress},
			Ingress: []networkingv1.NetworkPolicyIngressRule{{
				Ports: []networkingv1.NetworkPolicyPort{{Protocol: tcpProtocol(), Port: &proxyPort}},
			}},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{Protocol: &udp, Port: &dns},
						{Protocol: tcpProtocol(), Port: &dns},
					},
				},
				{
					To: []networkingv1.NetworkPolicyPeer{{
						NamespaceSelector: &metav1.LabelSelector{},
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"component": "apiserver"},
						},
					}},
					Ports: []networkingv1.NetworkPolicyPort{
						{Protocol: tcpProtocol(), Port: &https},
						{Protocol: tcpProtocol(), Port: &kubeAPI},
					},
				},
			},
		},
	}
}

func ProxyRoute(m *Context) *routev1.Route {
	weight := int32(100)
	targetPort := intstrFromInt32(names.ProxyPort)
	obj := &routev1.Route{
		Name:      names.ProxyRoute,
		Namespace: m.OperandNamespace,
		Labels:    labels.Proxy(m.CR),
		Annotations: map[string]string{
			"haproxy.router.openshift.io/hsts_header": "max-age=31536000;includeSubDomains;preload",
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   names.ProxyService,
				Weight: &weight,
			},
			Port:           &routev1.RoutePort{TargetPort: targetPort},
			TLS:            &routev1.TLSConfig{Termination: routev1.TLSTerminationReencrypt},
			WildcardPolicy: routev1.WildcardPolicyNone,
		},
	}
	return obj
}
