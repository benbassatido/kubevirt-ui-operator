package operconfig

import "testing"

func TestOperandNamespaceFromEnv(t *testing.T) {
	t.Setenv(EnvOperandNamespace, "")
	t.Setenv(EnvDownstream, "")

	if got := operandNamespaceFromEnv(); got != NamespaceUpstream {
		t.Fatalf("default upstream: got %q want %q", got, NamespaceUpstream)
	}

	t.Setenv(EnvDownstream, "true")
	if got := operandNamespaceFromEnv(); got != NamespaceDownstream {
		t.Fatalf("downstream: got %q want %q", got, NamespaceDownstream)
	}

	t.Setenv(EnvOperandNamespace, "custom-ns")
	if got := operandNamespaceFromEnv(); got != "custom-ns" {
		t.Fatalf("explicit namespace: got %q", got)
	}
}
