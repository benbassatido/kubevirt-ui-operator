package operconfig

import "testing"

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
