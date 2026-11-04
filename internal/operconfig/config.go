package operconfig

import (
	"os"
	"strconv"
	"strings"
)

const (
	EnvOperandNamespace = "OPERAND_NAMESPACE"
	EnvDownstream       = "DOWNSTREAM"

	NamespaceUpstream   = "kubevirt-ui"
	NamespaceDownstream = "openshift-virtualization-ui"
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
	if downstream, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv(EnvDownstream))); downstream {
		return NamespaceDownstream
	}
	return NamespaceUpstream
}
