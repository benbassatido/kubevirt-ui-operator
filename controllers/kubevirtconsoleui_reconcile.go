package controllers

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/manifests"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/platform"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/status"
)

func (r *KubeVirtConsoleUIReconciler) reconcile(
	ctx context.Context,
	cr *consoleuiv1alpha1.KubeVirtConsoleUI,
	mctx *manifests.Context,
) error {
	met, reason, message := platform.ConsoleRequirements(ctx, r.Client)
	status.SetCondition(&cr.Status.Conditions, consoleuiv1alpha1.ConditionTypePlatformRequirementsMet,
		conditionStatus(met), reason, message)

	var reconcileErr error
	if met {
		reconcileErr = manifests.ReconcileStack(ctx, r.Client, r.Scheme, mctx)
	} else {
		reconcileErr = fmt.Errorf("%s: %s", reason, message)
	}

	if reconcileErr != nil {
		return reconcileErr
	}

	registered, regErr := r.reconcileConsoleRegistration(ctx, mctx)
	cr.Status.ConsolePluginRegistered = registered
	if regErr != nil {
		log.FromContext(ctx).Error(regErr, "console plugin registration failed")
		return regErr
	}
	return nil
}
