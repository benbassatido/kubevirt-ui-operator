package status

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OperandDeploymentsWaiting returns a non-empty message when plugin or proxy Deployments are not available.
func OperandDeploymentsWaiting(ctx context.Context, c client.Client, namespace, pluginName, proxyName string) string {
	var parts []string
	for _, name := range []string{pluginName, proxyName} {
		if msg := deploymentWaitingMessage(ctx, c, namespace, name); msg != "" {
			parts = append(parts, msg)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return fmt.Sprintf("Waiting for operand Deployments: %s", strings.Join(parts, "; "))
}

func deploymentWaitingMessage(ctx context.Context, c client.Client, namespace, name string) string {
	dep := &appsv1.Deployment{}
	key := client.ObjectKey{Namespace: namespace, Name: name}
	if err := c.Get(ctx, key, dep); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Sprintf("%s not found", name)
		}
		return fmt.Sprintf("%s lookup failed: %v", name, err)
	}
	ok, detail := deploymentAvailable(dep)
	if ok {
		return ""
	}
	if detail != "" {
		return fmt.Sprintf("%s (%s)", name, detail)
	}
	return fmt.Sprintf("%s not available", name)
}

func deploymentAvailable(dep *appsv1.Deployment) (bool, string) {
	desired := int32(1)
	if dep.Spec.Replicas != nil {
		desired = *dep.Spec.Replicas
	}
	if desired == 0 {
		return true, ""
	}
	for _, cond := range dep.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			if cond.Status == corev1.ConditionTrue {
				return true, ""
			}
			if cond.Message != "" {
				return false, cond.Message
			}
			return false, string(cond.Reason)
		}
	}
	if dep.Status.AvailableReplicas >= desired {
		return true, ""
	}
	return false, fmt.Sprintf("%d/%d pods available", dep.Status.AvailableReplicas, desired)
}
