package controllers

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/status"
)

const defaultRequeue = 30 * time.Second

func conditionStatus(met bool) metav1.ConditionStatus {
	if met {
		return metav1.ConditionTrue
	}
	return metav1.ConditionFalse
}

func (r *KubeVirtConsoleUIReconciler) updateStatusAndRequeue(
	ctx context.Context,
	cr *consoleuiv1alpha1.KubeVirtConsoleUI,
	reconcileErr error,
) (ctrl.Result, error) {
	status.ReadyFromOperands(&cr.Status.Conditions, reconcileErr)
	cr.Status.ObservedGeneration = cr.Generation

	if err := r.Status().Update(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if reconcileErr != nil {
		return ctrl.Result{RequeueAfter: defaultRequeue}, nil //nolint:nilerr // failure is in status; backoff via RequeueAfter
	}
	return ctrl.Result{}, nil
}
