package status

import (
	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetCondition upserts a condition by type.
func SetCondition(conditions *[]metav1.Condition, conditionType string, status metav1.ConditionStatus, reason, message string) {
	now := metav1.Now()
	for i := range *conditions {
		if (*conditions)[i].Type == conditionType {
			if (*conditions)[i].Status == status && (*conditions)[i].Reason == reason && (*conditions)[i].Message == message {
				return
			}
			(*conditions)[i].Status = status
			(*conditions)[i].Reason = reason
			(*conditions)[i].Message = message
			(*conditions)[i].LastTransitionTime = now
			return
		}
	}
	*conditions = append(*conditions, metav1.Condition{
		Type:               conditionType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: now,
	})
}

// ReadyFromOperands sets Ready/Progressing/Degraded from operand reconciliation and Deployment availability.
func ReadyFromOperands(conditions *[]metav1.Condition, reconcileErr error, deploymentsWaiting string) {
	if reconcileErr != nil {
		SetCondition(conditions, consoleuiv1alpha1.ConditionTypeProgressing, metav1.ConditionTrue, "Reconciling", "Applying console UI operands")
		SetCondition(conditions, consoleuiv1alpha1.ConditionTypeDegraded, metav1.ConditionTrue, "ReconcileFailed", reconcileErr.Error())
		SetCondition(conditions, consoleuiv1alpha1.ConditionTypeReady, metav1.ConditionFalse, "ReconcileFailed", reconcileErr.Error())
		return
	}
	if deploymentsWaiting != "" {
		SetCondition(conditions, consoleuiv1alpha1.ConditionTypeProgressing, metav1.ConditionTrue, "WaitingForOperands", deploymentsWaiting)
		SetCondition(conditions, consoleuiv1alpha1.ConditionTypeDegraded, metav1.ConditionFalse, "Reconciled", "Operands applied; waiting for Deployments")
		SetCondition(conditions, consoleuiv1alpha1.ConditionTypeReady, metav1.ConditionFalse, "DeploymentsNotAvailable", deploymentsWaiting)
		return
	}
	SetCondition(conditions, consoleuiv1alpha1.ConditionTypeProgressing, metav1.ConditionFalse, "Reconciled", "Operands are up to date")
	SetCondition(conditions, consoleuiv1alpha1.ConditionTypeDegraded, metav1.ConditionFalse, "Reconciled", "No degraded operands")
	SetCondition(conditions, consoleuiv1alpha1.ConditionTypeReady, metav1.ConditionTrue, "Reconciled", "Console UI stack reconciled")
}
