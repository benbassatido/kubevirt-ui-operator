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

func TestHCONamespaceFromEnv(t *testing.T) {
	t.Setenv(EnvHCONamespace, "")
	if got := hcoNamespaceFromEnv(); got != "" {
		t.Fatalf("explicit empty: got %q want empty", got)
	}

	t.Setenv(EnvHCONamespace, "custom-hco")
	if got := hcoNamespaceFromEnv(); got != "custom-hco" {
		t.Fatalf("explicit namespace: got %q", got)
	}
}

func TestLoadFromEnv_hcoNamespaceExplicitEmpty(t *testing.T) {
	t.Setenv(EnvHCONamespace, "")
	cfg := LoadFromEnv()
	if cfg.HCONamespace != "" {
		t.Fatalf("explicit empty HCO namespace: got %q", cfg.HCONamespace)
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
