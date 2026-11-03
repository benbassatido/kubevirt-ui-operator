package controllers

import (
	"context"
	"slices"

	operatorv1 "github.com/openshift/api/operator/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/manifests"
)

func (r *KubeVirtConsoleUIReconciler) reconcileConsoleRegistration(ctx context.Context, m *manifests.Context) (bool, error) {
	pluginName := m.Spec.ConsolePlugin.Name
	console := &operatorv1.Console{}
	if err := r.Get(ctx, client.ObjectKey{Name: "cluster"}, console); err != nil {
		return false, err
	}
	if slices.Contains(console.Spec.Plugins, pluginName) {
		return true, nil
	}
	patch := console.DeepCopy()
	patch.Spec.Plugins = append(slices.Clone(console.Spec.Plugins), pluginName)
	if err := r.Patch(ctx, patch, client.MergeFromWithOptions(console, client.MergeFromWithOptimisticLock{})); err != nil {
		return false, err
	}
	return true, nil
}
