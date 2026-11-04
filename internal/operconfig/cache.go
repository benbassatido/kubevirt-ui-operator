package operconfig

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/fields"

	consolev1 "github.com/openshift/api/console/v1"
	operatorv1 "github.com/openshift/api/operator/v1"

	schedulingv1 "k8s.io/api/scheduling/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	priorityClassName     = "kubevirt-cluster-critical"
	consolePluginName     = "kubevirt-plugin"
	consolePluginCRDName  = "consoleplugins.console.openshift.io"
)

// CacheOptions returns a scoped informer cache: operand namespace plus fixed cluster-scoped objects.
func CacheOptions(cfg Config) cache.Options {
	return cache.Options{
		DefaultNamespaces: map[string]cache.Config{
			cfg.OperandNamespace: {},
		},
		ByObject: map[client.Object]cache.ByObject{
			&schedulingv1.PriorityClass{}: {
				Field: fields.OneTermEqualSelector("metadata.name", priorityClassName),
			},
			&consolev1.ConsolePlugin{}: {
				Field: fields.OneTermEqualSelector("metadata.name", consolePluginName),
			},
			&operatorv1.Console{}: {
				Field: fields.OneTermEqualSelector("metadata.name", "cluster"),
			},
			&apiextensionsv1.CustomResourceDefinition{}: {
				Field: fields.OneTermEqualSelector("metadata.name", consolePluginCRDName),
			},
		},
	}
}
