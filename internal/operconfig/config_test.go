package operconfig

import (
	"testing"
)

func TestOperandNamespaceFromEnv(t *testing.T) {
	t.Setenv(EnvOperandNamespace, "")

	if got := operandNamespaceFromEnv(); got != NamespaceUpstream {
		t.Fatalf("default: got %q want %q", got, NamespaceUpstream)
	}

	t.Setenv(EnvOperandNamespace, "custom-ns")
	if got := operandNamespaceFromEnv(); got != "custom-ns" {
		t.Fatalf("explicit namespace: got %q", got)
	}
}

func TestLoadFromEnv_images(t *testing.T) {
	t.Setenv(EnvPluginImage, "plugin:tag")
	t.Setenv(EnvProxyImage, "proxy:tag")
	t.Setenv(EnvOperandNamespace, "ns")

	cfg := LoadFromEnv()
	if cfg.PluginImage != "plugin:tag" || cfg.ProxyImage != "proxy:tag" || cfg.OperandNamespace != "ns" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
