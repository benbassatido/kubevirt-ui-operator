package builders

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

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

func NginxConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Name:      m.Spec.Config.NginxConfigMapName,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		Data: map[string]string{
			"nginx.conf": nginxConfig(m.Spec.Plugin.Port),
		},
	}
}

func UserSettingsConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Name:      m.Spec.Config.UserSettingsConfigMapName,
		Namespace: m.OperandNamespace,
		Labels:    labels.UIConfig(m.CR),
		Data:      map[string]string{},
	}
}

func FeaturesConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Name:      m.Spec.Config.FeaturesConfigMapName,
		Namespace: m.OperandNamespace,
		Labels:    labels.UIConfig(m.CR),
		Data:      featuresData(&m.Spec),
	}
}
