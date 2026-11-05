package operconfig

import (
	"os"
	"strings"
)

const (
	EnvOperandNamespace = "OPERAND_NAMESPACE"
	EnvHCONamespace     = "HCO_NAMESPACE"
	EnvPluginImage      = "KUBEVIRT_CONSOLE_PLUGIN_IMAGE"
	EnvProxyImage       = "KUBEVIRT_APISERVER_PROXY_IMAGE"

	NamespaceUpstream = "kubevirt-ui"

	// DefaultHCONamespace is where the HyperConverged CR is installed on OCP/CNV clusters.
	DefaultHCONamespace = "openshift-cnv"

	DefaultPluginImage = "quay.io/kubevirt-ui/kubevirt-plugin:latest"
	DefaultProxyImage  = "quay.io/kubevirt-ui/kubevirt-apiserver-proxy:latest"
)

// Config holds operator runtime settings (env), not CR spec.
type Config struct {
	OperandNamespace string
	HCONamespace     string
	PluginImage      string
	ProxyImage       string
	// HCOInstalled is set at startup from the API server (HyperConverged CRD present).
	HCOInstalled bool
}

// LoadFromEnv reads operator configuration from the environment.
func LoadFromEnv() Config {
	return Config{
		OperandNamespace: operandNamespaceFromEnv(),
		HCONamespace:     hcoNamespaceFromEnv(),
		PluginImage:      envOrDefault(EnvPluginImage, DefaultPluginImage),
		ProxyImage:       envOrDefault(EnvProxyImage, DefaultProxyImage),
	}
}

func hcoNamespaceFromEnv() string {
	v, ok := os.LookupEnv(EnvHCONamespace)
	if !ok {
		return DefaultHCONamespace
	}
	return strings.TrimSpace(v)
}

func operandNamespaceFromEnv() string {
	if ns := strings.TrimSpace(os.Getenv(EnvOperandNamespace)); ns != "" {
		return ns
	}
	return NamespaceUpstream
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
