package operconfig

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/fields"

	consolev1 "github.com/openshift/api/console/v1"
	operatorv1 "github.com/openshift/api/operator/v1"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/hcoversion"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/platform"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"

	schedulingv1 "k8s.io/api/scheduling/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CacheOptions returns a scoped informer cache: operand namespace plus fixed cluster-scoped objects.
func CacheOptions(cfg Config) cache.Options {
	byObject := map[client.Object]cache.ByObject{
		&consoleuiv1alpha1.KubeVirtConsoleUI{}: {
			Field: fields.OneTermEqualSelector("metadata.name", names.ResourceName),
		},
		&schedulingv1.PriorityClass{}: {
			Field: fields.OneTermEqualSelector("metadata.name", names.PriorityClass),
		},
		&consolev1.ConsolePlugin{}: {
			Field: fields.OneTermEqualSelector("metadata.name", names.ConsolePluginName),
		},
		&operatorv1.Console{}: {
			Field: fields.OneTermEqualSelector("metadata.name", "cluster"),
		},
		&apiextensionsv1.CustomResourceDefinition{}: {
			Field: fields.OneTermEqualSelector("metadata.name", platform.ConsolePluginCRDName),
		},
	}
	if cfg.HCOInstalled && cfg.HCONamespace != "" {
		byObject[&hcoversion.HyperConverged{}] = cache.ByObject{
			Namespaces: map[string]cache.Config{
				cfg.HCONamespace: {},
			},
		}
	}
	return cache.Options{
		DefaultNamespaces: map[string]cache.Config{
			cfg.OperandNamespace: {},
		},
		ByObject: byObject,
	}
}
