package builders

import (
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
)

func ProxyServiceAccount(m *Context) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.Proxy.Name,
			Namespace: m.OperandNamespace,
			Labels:    labels.Proxy(m.CR),
		},
	}
}

func ProxyDeployment(m *Context) *appsv1.Deployment {
	replicas := consoleuiv1alpha1.Int32OrDefault(m.Spec.Proxy.Replicas, 1)
	proxy := m.Spec.Proxy
	container := corev1.Container{
		Name:            proxy.Name,
		Image:           proxy.Image,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports:           []corev1.ContainerPort{{Name: "https", ContainerPort: proxy.Port, Protocol: corev1.ProtocolTCP}},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resourceQuantity("10m"),
				corev1.ResourceMemory: resourceQuantity("100Mi"),
			},
		},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: boolPtr(false),
			Capabilities:             &corev1.Capabilities{Drop: []corev1.Capability{"ALL"}},
			ReadOnlyRootFilesystem:   boolPtr(true),
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: proxy.ServingCertSecret, MountPath: "/app/cert", ReadOnly: true},
		},
	}
	if len(proxy.TLSArgs) > 0 {
		container.Args = append([]string{}, proxy.TLSArgs...)
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proxy.Name,
			Namespace: m.OperandNamespace,
			Labels:    labels.Proxy(m.CR),
		},
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
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels.Proxy(m.CR),
					Annotations: map[string]string{
						"openshift.io/required-scc": "restricted-v2",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: proxy.Name,
					PriorityClassName:  proxy.PriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: boolPtr(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{container},
					Volumes: []corev1.Volume{{
						Name: proxy.ServingCertSecret,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{SecretName: proxy.ServingCertSecret, DefaultMode: int32Ptr(420)},
						},
					}},
				},
			},
		},
	}
}

func ProxyService(m *Context) *corev1.Service {
	proxy := m.Spec.Proxy
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proxy.ServiceName,
			Namespace: m.OperandNamespace,
			Labels:    labels.Proxy(m.CR),
			Annotations: map[string]string{
				"service.beta.openshift.io/serving-cert-secret-name": proxy.ServingCertSecret,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels.ProxySelector(m.CR),
			Ports: []corev1.ServicePort{{
				Name:       proxy.Name + "-port",
				Port:       proxy.Port,
				TargetPort: intstrFromInt32(proxy.Port),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
}

func ProxyNetworkPolicy(m *Context) *networkingv1.NetworkPolicy {
	proxyPort := intstrFromInt32(m.Spec.Proxy.Port)
	udp := corev1.ProtocolUDP
	dns := intstr.FromInt32(53)
	https := intstr.FromInt32(443)
	kubeAPI := intstr.FromInt32(6443)
	npLabels := labels.Proxy(m.CR)
	npLabels["np.kubevirt.io/allow-access-cluster-services"] = "true"
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubevirt-apiserver-proxy-np",
			Namespace: m.OperandNamespace,
			Labels:    npLabels,
		},
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
	route := m.Spec.Route
	weight := int32(100)
	targetPort := intstrFromInt32(m.Spec.Proxy.Port)
	obj := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      route.Name,
			Namespace: m.OperandNamespace,
			Labels:    labels.Proxy(m.CR),
			Annotations: map[string]string{
				"haproxy.router.openshift.io/hsts_header": "max-age=31536000;includeSubDomains;preload",
			},
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   m.Spec.Proxy.ServiceName,
				Weight: &weight,
			},
			Port:           &routev1.RoutePort{TargetPort: targetPort},
			TLS:            &routev1.TLSConfig{Termination: routev1.TLSTerminationReencrypt},
			WildcardPolicy: routev1.WildcardPolicyNone,
		},
	}
	if route.Host != "" {
		obj.Spec.Host = route.Host
	}
	return obj
}
