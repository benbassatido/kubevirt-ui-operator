package manifests

import (
	"fmt"

	consolev1 "github.com/openshift/api/console/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
)

func nginxConfig(port int32) string {
	return fmt.Sprintf(`error_log /dev/stdout info;
events {}
http {
  access_log         /dev/stdout;
  include            /etc/nginx/mime.types;
  default_type       application/octet-stream;
  keepalive_timeout  65;
  add_header X-Content-Type-Options nosniff;
  server {
    listen              %d ssl;
    listen              [::]:%d ssl;
    ssl_certificate     /var/serving-cert/tls.crt;
    ssl_certificate_key /var/serving-cert/tls.key;
    ssl_protocols       TLSv1.2 TLSv1.3;
    root                /usr/share/nginx/html;

    location = /plugin-manifest.json {
      add_header Cache-Control 'no-cache, no-store, must-revalidate, proxy-revalidate, max-age=0';
      add_header Pragma 'no-cache';
      add_header Expires '0';
    }
    location = /plugin-entry.js {
      add_header Cache-Control 'no-cache, no-store, must-revalidate, proxy-revalidate, max-age=0';
      add_header Pragma 'no-cache';
      add_header Expires '0';
    }
  }
}
`, port, port)
}

func featuresData(spec *consoleuiv1alpha1.KubeVirtConsoleUISpec) map[string]string {
	f := spec.Config.Features
	data := map[string]string{
		"automaticSubscriptionActivationKey":  f.AutomaticSubscriptionActivationKey,
		"automaticSubscriptionOrganizationId": f.AutomaticSubscriptionOrganizationId,
		"disabledGuestSystemLogsAccess":       f.DisabledGuestSystemLogsAccess,
		"kubevirtApiserverProxy":              f.KubevirtApiserverProxy,
		"loadBalancerEnabled":                 f.LoadBalancerEnabled,
		"nodePortAddress":                     f.NodePortAddress,
		"nodePortEnabled":                     f.NodePortEnabled,
	}
	if f.IPStackType != "" {
		data["ipStackType"] = f.IPStackType
	}
	if f.HCOVersion != "" {
		data["hcoVersion"] = f.HCOVersion
	}
	return data
}

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

func pluginServiceAccount(m *Context) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.Plugin.Name,
			Namespace: m.OperandNamespace,
			Labels:    labels.Plugin(m.CR),
		},
	}
}

func proxyServiceAccount(m *Context) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.Proxy.Name,
			Namespace: m.OperandNamespace,
			Labels:    labels.Proxy(m.CR),
		},
	}
}

func nginxConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.Config.NginxConfigMapName,
			Namespace: m.OperandNamespace,
			Labels:    labels.Plugin(m.CR),
		},
		Data: map[string]string{
			"nginx.conf": nginxConfig(m.Spec.Plugin.Port),
		},
	}
}

func userSettingsConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.Config.UserSettingsConfigMapName,
			Namespace: m.OperandNamespace,
			Labels:    labels.UIConfig(m.CR),
		},
		Data: map[string]string{},
	}
}

func featuresConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.Config.FeaturesConfigMapName,
			Namespace: m.OperandNamespace,
			Labels:    labels.UIConfig(m.CR),
		},
		Data: featuresData(&m.Spec),
	}
}

func pluginDeployment(m *Context) *appsv1.Deployment {
	replicas := consoleuiv1alpha1.Int32OrDefault(m.Spec.Plugin.Replicas, 1)
	plugin := m.Spec.Plugin
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plugin.Name,
			Namespace: m.OperandNamespace,
			Labels:    labels.Plugin(m.CR),
		},
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
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels.Plugin(m.CR),
					Annotations: map[string]string{
						"openshift.io/required-scc": "restricted-v2",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: plugin.Name,
					PriorityClassName:  plugin.PriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: boolPtr(true),
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
							AllowPrivilegeEscalation: boolPtr(false),
							Capabilities:             &corev1.Capabilities{Drop: []corev1.Capability{"ALL"}},
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: plugin.ServingCertSecret, MountPath: "/var/serving-cert", ReadOnly: true},
							{Name: m.Spec.Config.NginxConfigMapName, MountPath: "/etc/nginx/nginx.conf", SubPath: "nginx.conf", ReadOnly: true},
						},
					}},
					Volumes: []corev1.Volume{
						{Name: plugin.ServingCertSecret, VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{SecretName: plugin.ServingCertSecret, DefaultMode: int32Ptr(420)},
						}},
						{Name: m.Spec.Config.NginxConfigMapName, VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: m.Spec.Config.NginxConfigMapName},
							},
						}},
					},
				},
			},
		},
	}
}

func proxyDeployment(m *Context) *appsv1.Deployment {
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

func pluginService(m *Context) *corev1.Service {
	plugin := m.Spec.Plugin
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plugin.ServiceName,
			Namespace: m.OperandNamespace,
			Labels:    labels.Plugin(m.CR),
			Annotations: map[string]string{
				"service.beta.openshift.io/serving-cert-secret-name": plugin.ServingCertSecret,
			},
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

func proxyService(m *Context) *corev1.Service {
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

func configReaderRole(m *Context) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.RBAC.ConfigReaderRoleName,
			Namespace: m.OperandNamespace,
			Labels:    labels.Plugin(m.CR),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{m.Spec.Config.UserSettingsConfigMapName},
				Verbs:         []string{"get", "update", "patch"},
			},
			{
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{m.Spec.Config.FeaturesConfigMapName},
				Verbs:         []string{"get"},
			},
		},
	}
}

func configReaderRoleBinding(m *Context) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.RBAC.ConfigReaderRoleBindingName,
			Namespace: m.OperandNamespace,
			Labels:    labels.Plugin(m.CR),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     m.Spec.RBAC.ConfigReaderRoleName,
		},
		Subjects: []rbacv1.Subject{{
			APIGroup: rbacv1.GroupName,
			Kind:     "Group",
			Name:     "system:authenticated",
		}},
	}
}

func pluginNetworkPolicy(m *Context) *networkingv1.NetworkPolicy {
	port := intstrFromInt32(m.Spec.Plugin.Port)
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubevirt-console-plugin-np",
			Namespace: m.OperandNamespace,
			Labels:    labels.Plugin(m.CR),
		},
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

func proxyNetworkPolicy(m *Context) *networkingv1.NetworkPolicy {
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

func proxyRoute(m *Context) *routev1.Route {
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

func consolePlugin(m *Context) *consolev1.ConsolePlugin {
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
