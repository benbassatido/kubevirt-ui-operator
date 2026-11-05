package operconfig

import (
	"os"
	"strings"
)

const (
	EnvOperandNamespace = "OPERAND_NAMESPACE"

	NamespaceUpstream = "kubevirt-ui"
)

// Config holds operator runtime settings used to scope the manager cache and controller watches.
type Config struct {
	OperandNamespace string
}

// LoadFromEnv reads operator configuration from the environment.
func LoadFromEnv() Config {
	return Config{
		OperandNamespace: operandNamespaceFromEnv(),
	}
}

func operandNamespaceFromEnv() string {
	if ns := strings.TrimSpace(os.Getenv(EnvOperandNamespace)); ns != "" {
		return ns
	}
	return NamespaceUpstream
}
