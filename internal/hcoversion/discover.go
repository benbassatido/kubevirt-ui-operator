package hcoversion

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OperatorVersion reads status.versions[name=operator] from HyperConverged in namespace.
// Returns empty string when HCO is not installed or version is not reported yet.
func OperatorVersion(ctx context.Context, c client.Client, namespace string) (string, error) {
	if namespace == "" {
		return "", nil
	}
	list := &HyperConvergedList{}
	if err := c.List(ctx, list, client.InNamespace(namespace)); err != nil {
		if isHyperConvergedUnavailable(err) {
			return "", nil
		}
		return "", err
	}
	for i := range list.Items {
		if v := operatorVersionFromStatus(&list.Items[i]); v != "" {
			return v, nil
		}
	}
	return "", nil
}

func operatorVersionFromStatus(hc *HyperConverged) string {
	for _, item := range hc.Status.Versions {
		if item.Name == OperatorVersionName && item.Version != "" {
			return item.Version
		}
	}
	return ""
}

func isHyperConvergedUnavailable(err error) bool {
	if meta.IsNoMatchError(err) {
		return true
	}
	if apierrors.IsNotFound(err) {
		return true
	}
	if apierrors.IsForbidden(err) {
		return true
	}
	return false
}
