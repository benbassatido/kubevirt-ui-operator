package platform

import (
	"context"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorv1 "github.com/openshift/api/operator/v1"
)

const consolePluginCRDName = "consoleplugins.console.openshift.io"

// ConsoleRequirements reports whether the cluster exposes OpenShift/OKD console APIs.
func ConsoleRequirements(ctx context.Context, c client.Client) (met bool, reason, message string) {
	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := c.Get(ctx, types.NamespacedName{Name: consolePluginCRDName}, crd); err != nil {
		if errors.IsNotFound(err) {
			return false, "ConsolePluginCRDMissing", fmt.Sprintf("CRD %s is not installed (OpenShift/OKD console required)", consolePluginCRDName)
		}
		return false, "ConsolePluginCRDLookupFailed", err.Error()
	}

	consoles := &operatorv1.ConsoleList{}
	if err := c.List(ctx, consoles); err != nil {
		return false, "ConsoleOperatorLookupFailed", err.Error()
	}
	if len(consoles.Items) == 0 {
		return false, "ConsoleOperatorMissing", "no consoles.operator.openshift.io instance found"
	}

	return true, "ConsolePresent", "ConsolePlugin CRD and at least one consoles.operator.openshift.io instance are present"
}
