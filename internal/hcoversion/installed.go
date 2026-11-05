package hcoversion

import (
	"context"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CRDName is the HyperConverged CRD installed by HCO.
const CRDName = "hyperconvergeds.hco.kubevirt.io"

// Installed reports whether the HyperConverged CRD is registered on the cluster.
func Installed(ctx context.Context, c client.Client) bool {
	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := c.Get(ctx, client.ObjectKey{Name: CRDName}, crd); err != nil {
		if apierrors.IsNotFound(err) || meta.IsNoMatchError(err) {
			return false
		}
		return false
	}
	return true
}
