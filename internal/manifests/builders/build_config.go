package builders

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
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

func NginxConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Name:      names.NginxConfigMap,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		Data: map[string]string{
			"nginx.conf": nginxConfig(names.PluginPort),
		},
	}
}

func UserSettingsConfigMap(m *Context) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Name:      names.UserSettingsConfigMap,
		Namespace: m.OperandNamespace,
		Labels:    labels.UIConfig(m.CR),
		Data:      map[string]string{},
	}
}

func FeaturesConfigMapShell(m *Context, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Name:      names.FeaturesConfigMap,
		Namespace: m.OperandNamespace,
		Labels:    labels.UIConfig(m.CR),
		Data:      data,
	}
}
